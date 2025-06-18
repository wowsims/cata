package core

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type healthBar struct {
	unit *Unit

	currentHealth float64

	DamageTakenHealthMetrics *ResourceMetrics
}

func (unit *Unit) EnableHealthBar() {
	unit.healthBar = healthBar{
		unit:                     unit,
		DamageTakenHealthMetrics: unit.NewHealthMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken}),
	}
}

func (unit *Unit) HasHealthBar() bool {
	return unit.healthBar.unit != nil
}

func (hb *healthBar) reset(_ *Simulation) {
	if hb.unit == nil {
		return
	}
	hb.currentHealth = hb.MaxHealth()
}

func (hb *healthBar) MaxHealth() float64 {
	return hb.unit.stats[stats.Health]
}

func (hb *healthBar) CurrentHealth() float64 {
	return hb.currentHealth
}

func (hb *healthBar) CurrentHealthPercent() float64 {
	return hb.currentHealth / hb.unit.stats[stats.Health]
}

func (hb *healthBar) GainHealth(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to gain negative health!")
	}

	oldHealth := hb.currentHealth
	newHealth := min(oldHealth+amount, hb.unit.MaxHealth())
	metrics.AddEvent(amount, newHealth-oldHealth)

	if sim.Log != nil {
		hb.unit.Log(sim, "Gained %0.3f health from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, oldHealth, newHealth, hb.MaxHealth())
	}

	hb.currentHealth = newHealth
}

func (hb *healthBar) RemoveHealth(sim *Simulation, amount float64) {
	if amount < 0 {
		panic("Trying to remove negative health!")
	}

	oldHealth := hb.currentHealth
	newHealth := max(oldHealth-amount, 0)
	metrics := hb.DamageTakenHealthMetrics
	metrics.AddEvent(-amount, newHealth-oldHealth)

	// TMI calculations need timestamps and Max HP information for each damage taken event
	if hb.unit.Metrics.isTanking {
		entry := tmiListItem{
			Timestamp:      sim.CurrentTime,
			WeightedDamage: amount / hb.MaxHealth(),
		}
		hb.unit.Metrics.tmiList = append(hb.unit.Metrics.tmiList, entry)
	}

	if sim.Log != nil {
		hb.unit.Log(sim, "Spent %0.3f health from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, oldHealth, newHealth, hb.MaxHealth())
	}

	hb.currentHealth = newHealth
}

// Used for dynamic updates to maximum health from "Last Stand" effects
func (hb *healthBar) UpdateMaxHealth(sim *Simulation, bonusHealth float64, metrics *ResourceMetrics) {
	hb.unit.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})

	if bonusHealth >= 0 {
		hb.GainHealth(sim, bonusHealth, metrics)
	} else {
		hb.RemoveHealth(sim, max(0, min(-bonusHealth, hb.currentHealth-1))) // Last Stand effects always leave the player with at least 1 HP when they expire
	}
}

var ChanceOfDeathAuraLabel = "Chance of Death"

func (character *Character) trackChanceOfDeath(healingModel *proto.HealingModel) {
	character.Unit.Metrics.isTanking = false
	for _, target := range character.Env.Encounter.TargetUnits {
		if (target.CurrentTarget == &character.Unit) || (target.SecondaryTarget == &character.Unit) {
			character.Unit.Metrics.isTanking = true
		}
	}

	character.RegisterAura(Aura{
		Label:    ChanceOfDeathAuraLabel,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Damage > 0 {
				aura.Unit.RemoveHealth(sim, result.Damage)

				if aura.Unit.CurrentHealth() <= 0 && !aura.Unit.Metrics.Died {
					// Queue a pending action to let shield effects give health
					StartDelayedAction(sim, DelayedActionOptions{
						DoAt: sim.CurrentTime,
						OnAction: func(s *Simulation) {
							if aura.Unit.CurrentHealth() <= 0 && !aura.Unit.Metrics.Died {
								aura.Unit.Metrics.Died = true
								if sim.Log != nil {
									character.Log(sim, "Dead")
								}
							}
						},
					})
				}
			}
		},
		OnPeriodicDamageTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Damage > 0 {
				aura.Unit.RemoveHealth(sim, result.Damage)

				if aura.Unit.CurrentHealth() <= 0 && !aura.Unit.Metrics.Died {
					// Queue a pending action to let shield effects give health
					StartDelayedAction(sim, DelayedActionOptions{
						DoAt: sim.CurrentTime,
						OnAction: func(s *Simulation) {
							if aura.Unit.CurrentHealth() <= 0 && !aura.Unit.Metrics.Died {
								aura.Unit.Metrics.Died = true
								if sim.Log != nil {
									character.Log(sim, "Dead")
								}
							}
						},
					})
				}
			}
		},
	})

	if !character.Unit.Metrics.isTanking {
		return
	}

	if healingModel == nil {
		return
	}

	character.Unit.Metrics.tmiBin = healingModel.BurstWindow

	if healingModel.Hps != 0 {
		character.applyHealingModel(healingModel)
	}
}

func (character *Character) applyHealingModel(healingModel *proto.HealingModel) {
	// Store variance parameters for healing cadence. Note that low rolls on
	// cadence are special cased here so that the model is still well-behaved
	// when CadenceVariation exceeds CadenceSeconds.
	medianCadence := healingModel.CadenceSeconds
	if medianCadence == 0 {
		medianCadence = 2.0
	}
	minCadence := max(0.0, medianCadence-healingModel.CadenceVariation)
	cadenceVariationLow := medianCadence - minCadence

	healingModelActionID := ActionID{OtherID: proto.OtherAction_OtherActionHealingModel}
	healthMetrics := character.NewHealthMetrics(healingModelActionID)

	// Register a shield aura on the tank to model the aggregate impact of
	// shield spells that contribute towards the total modeled HPS.
	var healPerTick float64
	var absorbShield *DamageAbsorptionAura

	absorbFrac := Clamp(healingModel.AbsorbFrac, 0, 1)

	if absorbFrac > 0 {
		absorbShield = character.NewDamageAbsorptionAura(AbsorptionAuraConfig{
			Aura: Aura{
				Label:    "Healing Model Absorb Shield" + character.Label,
				ActionID: healingModelActionID,
			},
			ShieldStrengthCalculator: func(_ *Unit) float64 {
				return max(absorbShield.ShieldStrength, healPerTick*absorbFrac)
			},
		})
	}

	character.RegisterResetEffect(func(sim *Simulation) {
		// Hack since we don't have OnHealingReceived aura handlers yet.
		//ardentDefenderAura := character.GetAura("Ardent Defender")
		//willOfTheNecropolisAura := character.GetAura("Will of The Necropolis")

		// Initialize randomized cadence model
		timeToNextHeal := DurationFromSeconds(0.0)
		healPerTick = 0.0
		pa := &PendingAction{
			NextActionAt: timeToNextHeal,
		}

		pa.OnAction = func(sim *Simulation) {
			// Use modeled HPS to scale heal per tick based on random cadence
			healPerTick = healingModel.Hps * (float64(timeToNextHeal) / float64(time.Second)) * character.PseudoStats.HealingTakenMultiplier * character.PseudoStats.ExternalHealingTakenMultiplier

			if healPerTick > 0 {
				// Execute the direct portion of the heal
				character.GainHealth(sim, healPerTick*(1.0-absorbFrac), healthMetrics)

				// Turn the remainder into an absorb shield
				if absorbShield != nil {
					absorbShield.Activate(sim)
				}
			}

			// Might use this again in the future to track "absorb" metrics but currently disabled
			//if ardentDefenderAura != nil && character.CurrentHealthPercent() >= 0.35 {
			//	ardentDefenderAura.Deactivate(sim)
			//}

			// if willOfTheNecropolisAura != nil && character.CurrentHealthPercent() > 0.35 {
			// 	willOfTheNecropolisAura.Deactivate(sim)
			// }

			// Random roll for time to next heal. In the case where CadenceVariation exceeds CadenceSeconds, then
			// CadenceSeconds is treated as the median, with two separate uniform distributions to the left and right
			// of it.
			signRoll := sim.RandomFloat("Healing Cadence Variation Sign")
			magnitudeRoll := sim.RandomFloat("Healing Cadence Variation Magnitude")

			if signRoll < 0.5 {
				timeToNextHeal = DurationFromSeconds(minCadence + magnitudeRoll*cadenceVariationLow)
			} else {
				timeToNextHeal = DurationFromSeconds(medianCadence + magnitudeRoll*healingModel.CadenceVariation)
			}

			// Refresh action
			pa.NextActionAt = sim.CurrentTime + timeToNextHeal
			sim.AddPendingAction(pa)
		}

		sim.AddPendingAction(pa)
	})
}

func (character *Character) GetPresimOptions(playerConfig *proto.Player) *PresimOptions {
	healingModel := playerConfig.HealingModel
	if healingModel == nil || healingModel.Hps != 0 || healingModel.CadenceSeconds == 0 {
		// If Hps is not 0, then we don't need to run the presim.
		// Tank sims should always have nonzero Cadence set, even if disabled
		return nil
	}
	return &PresimOptions{
		SetPresimPlayerOptions: func(player *proto.Player) {
			player.HealingModel = nil
		},
		OnPresimResult: func(presimResult *proto.UnitMetrics, iterations int32, duration time.Duration) bool {
			character.applyHealingModel(&proto.HealingModel{
				Hps:            presimResult.Dtps.Avg * 1.50,
				CadenceSeconds: healingModel.CadenceSeconds,
			})
			return true
		},
	}
}
