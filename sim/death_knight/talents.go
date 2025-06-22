package death_knight

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (dk *DeathKnight) ApplyTalents() {
	if dk.Level >= 15 {
		dk.registerRoilingBlood()
		dk.registerPlagueLeech()
		dk.registerUnholyBlight()
	}

	if dk.Level >= 30 {
		dk.registerLichborne()
		dk.registerAntiMagicZone()
		dk.registerPurgatory()
	}

	if dk.Level >= 45 {
		dk.registerDeathsAdvance()
		dk.registerChillblains()
		dk.registerAsphyxiate()
	}

	if dk.Level >= 60 {
		dk.registerDeathPact()
		dk.registerDeathSiphon()
		dk.registerConversion()
	}

	if dk.Level >= 75 {
		dk.registerBloodTap()
		dk.registerRunicCorruption()
		dk.registerRunicEmpowerment()
	}

	if dk.Level >= 90 {
		dk.registerGorefiendsGrasp()
		dk.registerRemorselessWinter()
		dk.registerDesecratedGround()
	}
}

// Your Blood Boil ability now also triggers Pestilence if it strikes a diseased target.
func (dk *DeathKnight) registerRoilingBlood() {
	if !dk.Talents.RoilingBlood {
		return
	}

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Roiling Blood" + dk.Label,
		ActionID:       core.ActionID{SpellID: 108170},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: DeathKnightSpellBloodBoil,
		Outcome:        core.OutcomeLanded,
		ICD:            core.SpellBatchWindow,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return dk.BloodPlagueSpell.Dot(result.Target).IsActive() ||
				dk.FrostFeverSpell.Dot(result.Target).IsActive()
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dk.PestilenceSpell.Cost.PercentModifier -= 100
			dk.PestilenceSpell.Cast(sim, result.Target)
			dk.PestilenceSpell.Cost.PercentModifier += 100
		},
	})
}

// Draw forth the infection from an enemy, consuming your Blood Plague and Frost Fever diseases on the target to activate two random fully-depleted runes as Death Runes.
func (dk *DeathKnight) registerPlagueLeech() {
	if !dk.Talents.PlagueLeech {
		return
	}

	actionID := core.ActionID{SpellID: 123693}
	runeMetrics := []*core.ResourceMetrics{
		dk.NewBloodRuneMetrics(actionID),
		dk.NewFrostRuneMetrics(actionID),
		dk.NewUnholyRuneMetrics(actionID),
		dk.NewDeathRuneMetrics(actionID),
	}

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellPlagueLeech,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 25,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return dk.BloodPlagueSpell.Dot(target).IsActive() &&
				dk.FrostFeverSpell.Dot(target).IsActive() &&
				dk.AnyDepletedRunes()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if !result.Landed() {
				return
			}

			dk.BloodPlagueSpell.Dot(target).Deactivate(sim)
			dk.FrostFeverSpell.Dot(target).Deactivate(sim)
			dk.ConvertAndRegenPlagueLeechRunes(sim, spell, runeMetrics)
		},
	})
}

// Surrounds the Death Knight with a vile swarm of unholy insects for 10 sec, stinging all enemies within 10 yards every 1 sec, infecting them with Blood Plague and Frost Fever.
func (dk *DeathKnight) registerUnholyBlight() {
	if !dk.Talents.UnholyBlight {
		return
	}

	unholyBlight := dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 115994},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagPassiveSpell,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, target := range sim.Encounter.TargetUnits {
				dk.BloodPlagueSpell.Cast(sim, target)
				dk.FrostFeverSpell.Cast(sim, target)
			}
		},
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 115989},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: DeathKnightSpellUnholyBlight,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Unholy Blight" + dk.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				unholyBlight.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Hot(&dk.Unit).Apply(sim)
		},
	})
}

/*
Draw upon unholy energy to become undead for 10 sec.
While undead, you are immune to Charm, Fear, and Sleep effects, and Death Coil will heal you.
*/
func (dk *DeathKnight) registerLichborne() {
	if !dk.Talents.Lichborne {
		return
	}

	actionID := core.ActionID{SpellID: 49039}
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: DeathKnightSpellLichborne,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: dk.RegisterAura(core.Aura{
			Label:    "Lichborne",
			ActionID: actionID,
			Duration: time.Second * 10,
		}),
	})
}

/*
Places a large, stationary Anti-Magic Zone that reduces spell damage done to party or raid members inside it by 40%.
The Anti-Magic Zone lasts for 3 sec.
*/
func (dk *DeathKnight) registerAntiMagicZone() {
	if !dk.Talents.AntiMagicZone {
		return
	}

	antiMagicZoneAuras := dk.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		spellDamageMultiplier := 0.6

		return unit.RegisterAura(core.Aura{
			Label:    "Anti-Magic Zone" + unit.Label,
			ActionID: core.ActionID{SpellID: 145629},
			Duration: time.Second * 3,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= spellDamageMultiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= spellDamageMultiplier
				unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= spellDamageMultiplier
			},
		})
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51052},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: DeathKnightSpellAntiMagicZone,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, unit := range sim.Raid.AllUnits {
				antiMagicZoneAuras.Get(unit).Activate(sim)
			}
		},

		RelatedAuraArrays: antiMagicZoneAuras.ToMap(),
	})
}

/*
An unholy pact grants you the ability to fight on through damage that would kill mere mortals.
When you would sustain fatal damage, you instead are wrapped in a Shroud of Purgatory, absorbing incoming healing equal to the amount of damage prevented, lasting 3 sec.
If any healing absorption remains when Shroud of Purgatory expires, you die.
Otherwise, you survive.
This effect may only occur every 3 min.
*/
func (dk *DeathKnight) registerPurgatory() {
	if !dk.Talents.Purgatory {
		return
	}

	perditionAura := dk.RegisterAura(core.Aura{
		Label:    "Perdition" + dk.Label,
		ActionID: core.ActionID{SpellID: 123981},
		Duration: time.Minute * 3,
	})

	actionID := core.ActionID{SpellID: 116888}
	healthMetrics := dk.NewHealthMetrics(actionID)

	var currentShield float64
	shroudOfPurgatoryAura := dk.RegisterAura(core.Aura{
		Label:     "Shroud of Purgatory" + dk.Label,
		ActionID:  actionID,
		Duration:  time.Second * 3,
		MaxStacks: math.MaxInt32,

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if currentShield > 0 {
				dk.RemoveHealth(sim, dk.CurrentHealth())
				dk.Died(sim)
			}
		},
	})

	dk.AddDynamicHealingTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !shroudOfPurgatoryAura.IsActive() {
			return
		}

		originalHealing := result.Damage
		if result.Damage >= currentShield {
			result.Damage -= currentShield
			currentShield = 0
			shroudOfPurgatoryAura.Deactivate(sim)
		} else {
			currentShield -= result.Damage
			shroudOfPurgatoryAura.SetStacks(sim, int32(currentShield))
			result.Damage = 0
		}

		if sim.Log != nil {
			dk.Log(sim, "Purgatory absorbed %.1f healing", originalHealing-result.Damage)
		}
	})

	dk.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult, isPeriodic bool) {
		if result.Damage < dk.CurrentHealth() {
			return
		}

		if perditionAura.IsActive() && !shroudOfPurgatoryAura.IsActive() {
			return
		}

		var newDamage float64
		if shroudOfPurgatoryAura.IsActive() {
			newDamage = 0
			dk.GainHealth(sim, 1.0, healthMetrics)
			currentShield += result.Damage - 1.0
		} else {
			newDamage = dk.CurrentHealth() - 1
			currentShield = result.Damage - newDamage

			shroudOfPurgatoryAura.Activate(sim)
			if !perditionAura.IsActive() {
				perditionAura.Activate(sim)
			}
		}

		shroudOfPurgatoryAura.SetStacks(sim, int32(currentShield))

		if sim.Log != nil {
			dk.Log(sim, "Purgatory absorbed %.1f damage", result.Damage-newDamage)
		}

		result.Damage = newDamage
	})
}

/*
You passively move 10% faster, and movement-impairing effects may not reduce you below 70% of normal movement speed.
When activated, you gain 30% movement speed and may not be slowed below 100% of normal movement speed for 6 seconds.
*/
func (dk *DeathKnight) registerDeathsAdvance() {
	if !dk.Talents.DeathsAdvance {
		return
	}

	actionID := core.ActionID{SpellID: 96268}
	dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: dk.RegisterAura(core.Aura{
			Label:    "Death's Advance" + dk.Label,
			ActionID: actionID,
			Duration: time.Second * 6,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				dk.MultiplyMovementSpeed(sim, 1.3)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				dk.MultiplyMovementSpeed(sim, 1.0/1.3)
			},
		}),
	})
}

// Victims of your Frost Fever disease are Chilled, reducing movement speed by 50% for 10 sec, and your Chains of Ice immobilizes targets for 3 sec.
func (dk *DeathKnight) registerChillblains() {
	if !dk.Talents.Chilblains {
		return
	}
}

/*
Lifts an enemy target off the ground and crushes their throat with dark energy, stunning them for 5 sec.
Functions as a silence if the target is immune to stuns.

Replaces Strangulate.
*/
func (dk *DeathKnight) registerAsphyxiate() {
	if !dk.Talents.Asphyxiate {
		return
	}
}

// Drain vitality from an undead minion, healing the Death Knight for 50% of his maximum health and causing the minion to suffer damage equal to 50% of its maximum health.
func (dk *DeathKnight) registerDeathPact() {
	actionID := core.ActionID{SpellID: 48743}

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagAPL | core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,
		ClassSpellMask: DeathKnightSpellDeathPact,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return dk.Ghoul.Pet.IsEnabled()
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healthGain := dk.MaxHealth() * 0.5
			spell.CalcAndDealHealing(sim, spell.Unit, healthGain, spell.OutcomeHealing)
			dk.Ghoul.RemoveHealth(sim, dk.Ghoul.MaxHealth()*0.5)
			if dk.Ghoul.CurrentHealth() <= 0 {
				dk.Ghoul.Disable(sim)
			}
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
		})
	}
}

var DeathSiphonActionID = core.ActionID{SpellID: 108196}

// Deal 6926 Shadowfrost damage to an enemy, healing the Death Knight for 150% of damage dealt.
func (dk *DeathKnight) registerDeathSiphon() {
	if !dk.Talents.DeathSiphon {
		return
	}

	siphonedDamage := 0.0
	deathSiphonHealSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 116783},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

		DamageMultiplier: 1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, siphonedDamage, spell.OutcomeHealing)
		},
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathSiphonActionID,
		SpellSchool:    core.SpellSchoolShadowFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellDeathSiphon,

		MaxRange: 40,

		RuneCost: core.RuneCostOptions{
			DeathRuneCost:  1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcAndRollDamageRange(sim, 6.59999990463, 0.15000000596) + 0.37400001287*spell.MeleeAttackPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				siphonedDamage = result.Damage
				deathSiphonHealSpell.Cast(sim, &dk.Unit)
			}

			spell.DealDamage(sim, result)
		},
	})

	if dk.Spec == proto.Spec_SpecBloodDeathKnight {
		dk.RuneWeapon.AddCopySpell(DeathSiphonActionID, dk.registerDrwDeathSiphon())
	}
}

func (dk *DeathKnight) registerDrwDeathSiphon() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathSiphonActionID,
		SpellSchool: core.SpellSchoolShadowFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		MaxRange: 40,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcAndRollDamageRange(sim, 6.59999990463, 0.15000000596) + 0.37400001287*spell.MeleeAttackPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

/*
Continuously converts Runic Power to health, restoring 3% of maximum health every 1 sec.
Only base Runic Power generation from spending runes may occur while Conversion is active.
This effect lasts until canceled, or Runic Power is exhausted.
*/
func (dk *DeathKnight) registerConversion() {
	if !dk.Talents.Conversion {
		return
	}

	conversionHealSpell := dk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 119980},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := dk.MaxHealth() * 0.03
			spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealing)
		},
	})

	actionID := core.ActionID{SpellID: 119975}

	var conversionSpell *core.Spell
	var healPa *core.PendingAction
	dk.ConversionAura = dk.RegisterAura(core.Aura{
		Label:    "Conversion" + dk.Label,
		ActionID: actionID,
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			conversionHealSpell.Cast(sim, &dk.Unit)

			healPa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second,
				OnAction: func(sim *core.Simulation) {
					if !conversionSpell.Cost.MeetsRequirement(sim, conversionSpell) {
						dk.ConversionAura.Deactivate(sim)
						return
					}

					conversionSpell.Cost.SpendCost(sim, conversionSpell)
					conversionHealSpell.Cast(sim, &dk.Unit)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			healPa.Cancel(sim)
		},
	})

	conversionSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: DeathKnightSpellConversion,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 5,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: dk.ConversionAura,
	})
}

func (dk *DeathKnight) registerBloodTap() {
	if !dk.Talents.BloodTap {
		return
	}

	bloodChargeAura := dk.RegisterAura(core.Aura{
		Label:     "Blood Charge" + dk.Label,
		ActionID:  core.ActionID{SpellID: 114851},
		Duration:  time.Second * 25,
		MaxStacks: 12,
	})

	actionID := core.ActionID{SpellID: 45529}

	runeMetrics := []*core.ResourceMetrics{
		dk.NewBloodRuneMetrics(actionID),
		dk.NewFrostRuneMetrics(actionID),
		dk.NewUnholyRuneMetrics(actionID),
		dk.NewDeathRuneMetrics(actionID),
		dk.NewRunicPowerMetrics(actionID),
	}

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellBloodTap,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bloodChargeAura.GetStacks() >= 5 && dk.AnyDepletedRunes()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if dk.ConvertAndRegenBloodTapRune(sim, spell, runeMetrics) {
				bloodChargeAura.RemoveStacks(sim, 5)
			}
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Blood Charge Trigger" + dk.Label,
		Callback:       core.CallbackOnSpellHitDealt,
		ProcMask:       core.ProcMaskMeleeMH | core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellDeathCoil | DeathKnightSpellFrostStrike | DeathKnightSpellRuneStrike,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bloodChargeAura.Activate(sim)
			bloodChargeAura.AddStacks(sim, 2)
		},
	})
}

func (dk *DeathKnight) getRunicMasteryAura() *core.StatBuffAura {
	if dk.CouldHaveSetBonus(ItemSetNecroticBoneplateBattlegear, 4) {
		return dk.NewTemporaryStatsAura("Runic Mastery", core.ActionID{SpellID: 105647}, stats.Stats{stats.MasteryRating: 710}, time.Second*12)
	}

	return nil
}

/*
When you land a damaging Death Coil, Frost Strike, or Rune Strike, you have a 45% chance to activate a random fully-depleted rune.
(Proc chance: 45%)
*/
func (dk *DeathKnight) registerRunicEmpowerment() {
	if !dk.Talents.RunicEmpowerment {
		return
	}

	runicMasteryAura := dk.getRunicMasteryAura()

	// Runic Empowerement refreshes random runes on cd
	actionID := core.ActionID{SpellID: 81229}
	runeMetrics := []*core.ResourceMetrics{
		dk.NewBloodRuneMetrics(actionID),
		dk.NewFrostRuneMetrics(actionID),
		dk.NewUnholyRuneMetrics(actionID),
		dk.NewDeathRuneMetrics(actionID),
	}

	core.MakePermanent(dk.RegisterAura(core.Aura{
		Label:    "Runic Empowerement" + dk.Label,
		ActionID: actionID,
	}))

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Runic Empowerement Trigger" + dk.Label,
		Callback:       core.CallbackOnSpellHitDealt,
		ProcMask:       core.ProcMaskMeleeMH | core.ProcMaskSpellDamage,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: DeathKnightSpellDeathCoil | DeathKnightSpellFrostStrike | DeathKnightSpellRuneStrike,
		ProcChance:     0.45,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dk.RegenRunicEmpowermentRune(sim, runeMetrics)

			// T13 4pc: Runic Empowerment has a 25% chance to also grant 710 mastery rating for 12 sec when activated.
			if dk.T13Dps4pc.IsActive() && sim.Proc(0.25, "T13 4pc") {
				runicMasteryAura.Activate(sim)
			}
		},
	})
}

/*
When you land a damaging Death Coil, Frost Strike, or Rune Strike, you have a 45% chance to activate Runic Corruption, increasing your rune regeneration rate by 100% for 3 sec.
(Proc chance: 45%)
*/
func (dk *DeathKnight) registerRunicCorruption() {
	if !dk.Talents.RunicCorruption {
		return
	}

	runicMasteryAura := dk.getRunicMasteryAura()

	duration := time.Second * 3
	multi := 2.0
	// Runic Corruption gives rune regen speed
	regenAura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Runic Corruption",
		ActionID: core.ActionID{SpellID: 51460},
		Duration: duration,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyRuneRegenSpeed(sim, multi)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.MultiplyRuneRegenSpeed(sim, 1/multi)
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Runic Corruption Trigger" + dk.Label,
		Callback:       core.CallbackOnSpellHitDealt,
		ProcMask:       core.ProcMaskMeleeMH | core.ProcMaskSpellDamage,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: DeathKnightSpellDeathCoil | DeathKnightSpellFrostStrike | DeathKnightSpellRuneStrike,
		ProcChance:     0.45,

		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			hasteMultiplier := 1.0 + dk.GetStat(stats.HasteRating)/(100*core.HasteRatingPerHastePercent)
			if regenAura.IsActive() {
				totalMultiplier := 1 / (hasteMultiplier * (dk.GetRuneRegenMultiplier() / multi))
				hastedDuration := core.DurationFromSeconds(duration.Seconds() * totalMultiplier)
				regenAura.UpdateExpires(regenAura.ExpiresAt() + hastedDuration)
			} else {
				totalMultiplier := 1 / (hasteMultiplier * dk.GetRuneRegenMultiplier())
				hastedDuration := core.DurationFromSeconds(duration.Seconds() * totalMultiplier)
				regenAura.Duration = hastedDuration
				regenAura.Activate(sim)
			}

			// T13 4pc: Runic Corruption has a 40% chance to also grant 710 mastery rating for 12 sec when activated.
			if dk.T13Dps4pc.IsActive() && sim.Proc(0.4, "T13 4pc") {
				runicMasteryAura.Activate(sim)
			}
		},
	})
}

// Shadowy tendrils coil around all enemies within 20 yards of a target (hostile or friendly), pulling them to the target's location.
func (dk *DeathKnight) registerGorefiendsGrasp() {
	if !dk.Talents.GorefiendsGrasp {
		return
	}
}

/*
Surrounds the Death Knight with a swirling tempest of frigid air for 8 sec, chilling enemies within 8 yards every 1 sec.
Each pulse reduces targets' movement speed by 15% for 3 sec, stacking up to 5 times.
Upon receiving a fifth application, an enemy will be stunned for 6 sec.
*/
func (dk *DeathKnight) registerRemorselessWinter() {
	if !dk.Talents.RemorselessWinter {
		return
	}
}

/*
Corrupts the ground in a 8 yard radius beneath the Death Knight for 10 sec.
While standing in this corruption, the Death Knight is immune to effects that cause loss of control.
This ability instantly removes such effects when activated.
*/
func (dk *DeathKnight) registerDesecratedGround() {
	if !dk.Talents.DesecratedGround {
		return
	}
}
