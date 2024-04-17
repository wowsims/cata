package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (mage *Mage) ApplyTalents() {

	mage.applyIgnite()
	// mage.applyImpact()
	// mage.applyIceFloes()
	// mage.applyPiercingChill()
	// mage.applyPermaFrost

	// mage.applyImprovedFreeze()
	// mage.applyEnduringWinter()
	// mage.applyColdSnap()

	// mage.applyImprovedFlamestrike()
	// mage.applyCriticalMass()
	// mage.applyFrostfireOrb()

	// Stat Buffs
	// mage.PseudoStats.SpiritRegenRateCasting += float64(mage.Talents.ArcaneMeditation) / 6
	// if mage.Talents.StudentOfTheMind > 0 {
	// 	mage.MultiplyStat(stats.Spirit, 1.0+[]float64{0, .04, .07, .10}[mage.Talents.StudentOfTheMind])
	// }
	// if mage.Talents.ArcaneMind > 0 {
	// 	mage.MultiplyStat(stats.Intellect, 1.0+0.03*float64(mage.Talents.ArcaneMind))
	// }
	// if mage.Talents.MindMastery > 0 {
	// 	mage.AddStatDependency(stats.Intellect, stats.SpellPower, 0.03*float64(mage.Talents.MindMastery))
	// }
	// mage.AddStat(stats.SpellCrit, float64(mage.Talents.ArcaneInstability)*1*core.CritRatingPerCritChance)
	// mage.PseudoStats.DamageDealtMultiplier *= 1 + .01*float64(mage.Talents.ArcaneInstability)
	// mage.PseudoStats.DamageDealtMultiplier *= 1 + .01*float64(mage.Talents.PlayingWithFire)

	/* --------------------------------------
				  ARCANE TALENTS
	---------------------------------------*/
	// Arcane Specialization Bonus
	if mage.Spec == proto.Spec_SpecArcaneMage {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion, // | MageSpellArcaneMissiles,
			FloatValue: 0.25,
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	// Cooldowns/Special Implementations
	mage.registerArcanePowerCD()
	mage.registerPresenceOfMindCD()
	mage.applyFocusMagic()
	mage.applyArcanePotency()
	mage.applyArcaneConcentration()

	// Netherwind Presence
	if mage.Talents.NetherwindPresence > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcaneBarrage,
			FloatValue: -0.01 * float64(mage.Talents.NetherwindPresence) * core.HasteRatingPerHastePercent,
			Kind:       core.SpellMod_CastTime_Pct,
		})
	}

	// Torment the Weak
	if mage.Talents.TormentTheWeak > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion, //| MageSpellArcaneMissiles,
			FloatValue: 0.02 * float64(mage.Talents.TormentTheWeak),
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	//Improved Arcane Missiles
	if mage.Talents.ImprovedArcaneMissiles > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneMissilesCast,
			IntValue:  int64(mage.Talents.ImprovedArcaneMissiles),
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
		})
	}

	// Missile Barrage
	if mage.Talents.MissileBarrage > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneMissilesCast,
			TimeValue: time.Millisecond * time.Duration(-100*mage.Talents.MissileBarrage),
			Kind:      core.SpellMod_DotTickLength_Flat,
		})
	}

	// ArcaneFlows
	// Implemented inside relevant spells due to % cooldown reduction

	// Arcane Tactics
	// Raid buff, tbd

	// Improved Arcane Explosion
	if mage.Talents.ImprovedArcaneExplosion > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellArcaneExplosion,
			TimeValue: -1 * time.Duration(0.3*float64(mage.Talents.ImprovedArcaneExplosion)),
			Kind:      core.SpellMod_GlobalCooldown_Flat,
		})
	}
	if mage.Talents.ImprovedArcaneExplosion > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellArcaneExplosion,
			FloatValue: -0.25 * float64(mage.Talents.ImprovedArcaneExplosion),
			Kind:       core.SpellMod_PowerCost_Pct,
		})
	}

	//

	/* --------------------------------------
				  FIRE TALENTS
	---------------------------------------*/
	// Fire  Specialization Bonus
	if mage.Spec == proto.Spec_SpecFireMage {
		mage.AddStaticMod(core.SpellModConfig{
			School:     core.SpellSchoolFire,
			FloatValue: 0.25,
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	// Mastery
	if mage.Spec == proto.Spec_SpecFireMage {
		//
		fireMastery := mage.AddDynamicMod(core.SpellModConfig{
			ClassMask:  MageSpellFireDoT,
			FloatValue: float64(1.22 + 0.28*mage.GetMasteryPoints()),
			Kind:       core.SpellMod_DamageDone_Pct,
		})
		fireMastery.Activate()

		mage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
			fireMastery.UpdateFloatValue(1.22 + 0.28*core.MasteryRatingToMasteryPoints(newMastery))
		})
	}

	// Cooldowns/Special Implementations
	mage.applyHotStreak()
	mage.applyMoltenFury()
	mage.applyMasterOfElements()
	mage.applyPyromaniac()

	// Improved Fire Blast
	if mage.Talents.ImprovedFireBlast > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellFireBlast,
			FloatValue: .04 * float64(mage.Talents.ImprovedFireBlast) * core.CritRatingPerCritChance,
			Kind:       core.SpellMod_BonusCrit_Rating,
		})
	}

	// Fire Power
	if mage.Talents.FirePower > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			School:     core.SpellSchoolFire,
			FloatValue: 0.01 * float64(mage.Talents.FirePower),
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	// Improved Scorch
	if mage.Talents.ImprovedScorch > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellScorch,
			FloatValue: -0.5 * float64(mage.Talents.ImprovedScorch),
			Kind:       core.SpellMod_PowerCost_Pct,
		})
	}

	// Improved Flamestrike
	if mage.Talents.ImprovedFlamestrike > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellFlamestrike,
			FloatValue: -0.5 * float64(mage.Talents.ImprovedFlamestrike),
			Kind:       core.SpellMod_CastTime_Pct,
		})
	}

	// Pyromaniac

	// Critical Mass
	if mage.Talents.CriticalMass > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellLivingBombDot | MageSpellLivingBombExplosion | MageSpellFlameOrb,
			FloatValue: 0.05 * float64(mage.Talents.CriticalMass),
			Kind:       core.SpellMod_DamageDone_Pct,
		})
	}

	/* --------------------------------------
				 FROST TALENTS
	---------------------------------------*/

	// Cooldowns/Special Implementations
	mage.registerIcyVeinsCD()
	mage.registerColdSnapCD()
	mage.applyFingersOfFrost()
	mage.applyEarlyFrost()
	mage.applyBrainFreeze()

	// Piercing Ice
	if mage.Talents.PiercingIce > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellsAll,
			FloatValue: 0.01 * float64(mage.Talents.PiercingIce) * core.CritRatingPerCritChance,
			Kind:       core.SpellMod_BonusCrit_Rating,
		})
	}

	//Enduring Winter
	if mage.Talents.EnduringWinter > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellsAll,
			FloatValue: -1 * []float64{0, 0.03, 0.06, 0.1}[mage.Talents.EnduringWinter],
			Kind:       core.SpellMod_PowerCost_Pct,
		})
	}
}

func (mage *Mage) applyPyromaniac() {
	if mage.Talents.Pyromaniac == 0 {
		return
	}
	/*
		pyromaniacMod := mage.AddDynamicMod(core.SpellModConfig{
			ClassMask:  MageSpellsAll,
			FloatValue: -.05 * float64(mage.Talents.Pyromaniac),
			Kind:       core.SpellMod_CastTime_Pct,
		})
		var activeFireDots []*core.Spell

		mage.PyromaniacAura = mage.RegisterAura(core.Aura{
			Label:    "Pyromaniac Trackers",
			ActionID: core.ActionID{SpellID: 83582},
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				if len(sim.AllUnits) < 3 {
					return
				}
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {

				 for _, aoeTarget := range sim.Encounter.TargetUnits {
					if mage.LivingBomb.Dot(aoeTarget).RemainingDuration(sim) > 0 {

					}
					if spell.ClassSpellMask == MageSpellFireDoT {
						activeFireDots = append(activeFireDots, spell)
						core.StartDelayedAction(sim, core.DelayedActionOptions{
							DoAt: sim.CurrentTime + spell.Dot(mage.CurrentTarget).RemainingDuration(sim),
							OnAction: func(sim *core.Simulation) {
								l := len(activeFireDots)
								activeFireDots = activeFireDots[:l-1]

							},
						})
					}
				}
				fmt.Println("activeFireDots: ", len(activeFireDots))

				if len(activeFireDots) >= 3 {
					pyromaniacMod.Activate()
				} else {
					pyromaniacMod.Deactivate()
				}
			},
		})*/
}

func (mage *Mage) applyHotStreak() {
	if !mage.Talents.HotStreak {
		return
	}

	ImprovedHotStreakProcChance := float64(mage.Talents.ImprovedHotStreak) * 0.5
	BaseHotStreakProcChance := float64(0.25) // Research needed
	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	// Unimproved Hot Streak Proc Aura
	mage.HotStreakAura = mage.RegisterAura(core.Aura{
		Label:    "Hot Streak",
		ActionID: core.ActionID{SpellID: 48108},
		Duration: time.Second * 10,
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if t10ProcAura != nil {
				t10ProcAura.Activate(sim)
			}
		},
	})

	// Improved Hotstreak Crit Stacking Aura
	mage.hotStreakCritAura = mage.RegisterAura(core.Aura{
		Label:     "Hot Streak Proc Aura",
		ActionID:  core.ActionID{SpellID: 44448}, //, Tag: 1}, Removing Tag gets rid of the (??) in Timeline
		MaxStacks: 2,
		Duration:  time.Hour,
	})

	// Aura to allow the character to track crits
	mage.RegisterAura(core.Aura{
		Label:    "Hot Streak Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(HotStreakSpells) {
				return
			}

			// Hot Streak Base Talent Proc
			if result.DidCrit() && spell.Flags.Matches(HotStreakSpells) {
				if sim.Proc(BaseHotStreakProcChance, "Hot Streak") {
					if mage.HotStreakAura.IsActive() {
						mage.HotStreakAura.Refresh(sim)
					} else {
						mage.HotStreakAura.Activate(sim)
					}
				}
			}

			// Improved Hot Streak
			if mage.Talents.ImprovedHotStreak > 0 {

				// If you didn't crit, reset your crit counter
				if !result.DidCrit() {
					mage.hotStreakCritAura.SetStacks(sim, 0)
					mage.hotStreakCritAura.Deactivate(sim)
					return
				}

				// If you did crit, check against talents to see if you proc
				// If you proc and had 1 stack, set crit counter to 0 and give hot streak.
				if mage.hotStreakCritAura.GetStacks() == 1 {
					if sim.Proc(ImprovedHotStreakProcChance, "Improved Hot Streak") {
						mage.hotStreakCritAura.SetStacks(sim, 0)
						mage.hotStreakCritAura.Deactivate(sim)

						mage.HotStreakAura.Activate(sim)
					}

					// If you proc and had 0 stacks of crits, add to your crit counter.
					// No idea if 1 out of 2 talent points means you have a 50% chance to
					// add to the 1st stack of crit, or only the 2nd. Doesn't seem
					// all that important to check since every fire mage in the world
					// will go 2 out of 2 points, but worth researching.
					// If it checks 1st crit as well, can add a proc check to this too
				} else {
					mage.hotStreakCritAura.Activate(sim)
					mage.hotStreakCritAura.AddStack(sim)
				}
			}
		},
	})

}

func (mage *Mage) applyArcanePotency() {
	if mage.Talents.ArcanePotency == 0 {
		return
	}

	arcanePotencyMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: []float64{0.0, 0.07, 0.15}[mage.Talents.ArcanePotency] * core.CritRatingPerCritChance,
		Kind:       core.SpellMod_BonusCrit_Rating,
	})

	mage.ArcanePotencyAura = mage.RegisterAura(core.Aura{
		Label:     "Arcane Potency",
		Duration:  core.NeverExpires,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			//prevent the spell that procced it from spending it
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + time.Millisecond*10,
				OnAction: func(sim *core.Simulation) {
					aura.SetStacks(sim, 2)
					arcanePotencyMod.Activate()
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			arcanePotencyMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Only remove a stack if it's an applicable spell
			if spell.ClassSpellMask == arcanePotencyMod.ClassMask {
				aura.RemoveStack(sim)
			}
		},
	},
	)
}

func (mage *Mage) applyArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	// The result that caused the proc. Used to check we don't deactivate from the same proc.
	var proccedAt time.Duration
	var proccedSpell *core.Spell

	// Tracks if Clearcasting should proc
	mage.RegisterAura(core.Aura{
		Label:    "Arcane Concentration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagMage) || spell == mage.ArcaneMissiles {
				return
			}
			if !result.Landed() {
				return
			}

			procChance := []float64{0, 0.03, 0.06, 0.1}[mage.Talents.ArcaneConcentration]
			// Arcane Missile ticks can proc CC, just at a low rate of about 1.5% with 5/5 Arcane Concentration
			if spell == mage.ArcaneMissilesTickSpell {
				procChance *= 0.15
			}
			if !sim.Proc(procChance, "Arcane Concentration") {
				return
			}
			proccedAt = sim.CurrentTime
			proccedSpell = spell

			//mage.ArcanePotencyAura.Activate(sim)
			mage.ClearcastingAura.Activate(sim)
			mage.ArcaneBlastAura.GetStacks()
		},
	})

	/* 	if mage.Talents.ArcanePotency > 0 {
		mage.ArcanePotencyAura = mage.RegisterAura(core.Aura{
			Label:    "Arcane Potency",
			ActionID: core.ActionID{SpellID: 31572},
			Duration: time.Hour,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, bonusCrit)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.SpellCrit, -bonusCrit)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.Flags.Matches(SpellFlagMage) {
					return
				}
				// Don't spend on the spell that procced it
				if proccedAt == sim.CurrentTime {
					return
				}
				aura.Deactivate(sim)
			},
		})
	} */
	clearCastingMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAllDamaging,
		FloatValue: -1,
		Kind:       core.SpellMod_PowerCost_Pct,
	})
	// The Clearcasting proc
	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			//mage.ArcanePotencyAura.Activate(sim)
			clearCastingMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			clearCastingMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}
			if spell.DefaultCast.Cost == 0 {
				return
			}
			if spell == mage.ArcaneMissiles && mage.ArcaneMissilesProcAura.IsActive() {
				return
			}
			if proccedAt == sim.CurrentTime && proccedSpell == spell {
				// Means this is another hit from the same cast that procced CC.
				return
			}
			aura.Deactivate(sim)
		},
	})
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	actionID := core.ActionID{SpellID: 12043}
	var spellToUse *core.Spell
	mage.Env.RegisterPostFinalizeEffect(func() {
		spellToUse = mage.ArcaneBlast
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(120*(1-[]float64{0.0, 0.07, 0.15}[mage.Talents.ArcaneFlows])),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			if !mage.GCD.IsReady(sim) {
				return false
			}
			if mage.ArcanePowerAura.IsActive() {
				return false
			}

			manaCost := spellToUse.DefaultCast.Cost * mage.PseudoStats.CostMultiplier
			if spellToUse == mage.ArcaneBlast {
				manaCost *= float64(mage.ArcaneBlastAura.GetStacks()) * 1.75
			}
			return mage.CurrentMana() >= manaCost
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.ArcanePotencyAura != nil {
				mage.ArcanePotencyAura.Activate(sim)
			}
			normalCastTime := spellToUse.DefaultCast.CastTime
			spellToUse.DefaultCast.CastTime = 0
			spellToUse.Cast(sim, mage.CurrentTarget)
			spellToUse.DefaultCast.CastTime = normalCastTime
		},
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}
	actionID := core.ActionID{SpellID: 12042}

	var affectedSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Flags.Matches(SpellFlagMage) {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	mage.ArcanePowerAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Power",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.2
				spell.CostMultiplier += 0.1
			}
			mage.arcanePowerGCDmod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive -= 0.2
				spell.CostMultiplier -= 0.2
			}
			mage.arcanePowerGCDmod.Deactivate()
		},
	})
	core.RegisterPercentDamageModifierEffect(mage.ArcanePowerAura, 1.2)

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(120*(1-[]float64{0.0, 0.07, 0.15}[mage.Talents.ArcaneFlows])),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.ArcanePowerAura.Activate(sim)
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.ArcanePotencyAura.IsActive()
		},
	})
	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

// Master of Elements
func (mage *Mage) applyMasterOfElements() {
	if mage.Talents.MasterOfElements == 0 {
		return
	}

	refundCoeff := 0.15 * float64(mage.Talents.MasterOfElements)
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 29077})

	mage.RegisterAura(core.Aura{
		Label:    "Master of Elements",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if spell.CurCast.Cost == 0 {
				return
			}
			if result.DidCrit() {
				if refundCoeff < 0 {
					mage.SpendMana(sim, -1*spell.DefaultCast.Cost*refundCoeff, manaMetrics)
				} else {
					mage.AddMana(sim, spell.DefaultCast.Cost*refundCoeff, manaMetrics)
				}
			}
		},
	})
}

func (mage *Mage) registerIcyVeinsCD() {
	if !mage.Talents.IcyVeins {
		return
	}

	icyVeinsMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: -0.2,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	actionID := core.ActionID{SpellID: 12472}
	icyVeinsAura := mage.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			icyVeinsMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			icyVeinsMod.Deactivate()
		},
	})

	mage.IcyVeins = mage.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: MageSpellIcyVeins,
		Flags:          core.SpellFlagNoOnCastComplete,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},

		Cast: core.CastConfig{

			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * time.Duration(180*[]float64{1, .93, .86, .80}[mage.Talents.IceFloes]),
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Need to check for icy veins already active in case Cold Snap is used right after.
			return !icyVeinsAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			icyVeinsAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.IcyVeins,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerColdSnapCD() {
	if !mage.Talents.ColdSnap {
		return
	}

	actionID := core.ActionID{SpellID: 11958}
	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 8,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Don't use if there are no cooldowns to reset.
			return (mage.IcyVeins != nil && !mage.IcyVeins.IsReady(sim)) ||
				(mage.SummonWaterElemental != nil && !mage.SummonWaterElemental.IsReady(sim))
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if mage.IcyVeins != nil {
				mage.IcyVeins.CD.Reset()
			}
			if mage.SummonWaterElemental != nil {
				mage.SummonWaterElemental.CD.Reset()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Ideally wait for both water ele and icy veins so we can reset both.
			if mage.IcyVeins != nil && mage.IcyVeins.IsReady(sim) {
				return false
			}
			if mage.SummonWaterElemental != nil && mage.SummonWaterElemental.IsReady(sim) {
				return false
			}
			return true
		},
	})
}

func (mage *Mage) applyMoltenFury() {
	if mage.Talents.MoltenFury == 0 {
		return
	}

	moltenFuryMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: .04 * float64(mage.Talents.MoltenFury),
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	mage.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 35 {
				moltenFuryMod.Activate()

				// For some reason Molten Fury doesn't apply to living bomb DoT, so cancel it out.
				// 4/15/24 - Comment above was from before. Worth checking this out.
				/*if mage.LivingBomb != nil {
					mage.LivingBomb.DamageMultiplier /= multiplier
				}*/
			}
		})
	})
}

func (mage *Mage) hasChillEffect(spell *core.Spell) bool {
	return spell == mage.Frostbolt || spell == mage.FrostfireBolt || (spell == mage.Blizzard && mage.Talents.IceShards > 0)
}

func (mage *Mage) applyFingersOfFrost() {
	if mage.Talents.FingersOfFrost == 0 {
		return
	}

	//Talent gives 7/14/20 percent chance to proc FoF on spell hit
	procChance := []float64{0, 0.07, 0.14, 0.20}[mage.Talents.FingersOfFrost]

	fingersOfFrostIceLanceDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellIceLance,
		FloatValue: 0.25,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	fingersOfFrostFrozenCritMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellIceLance | MageSpellDeepFreeze,
		FloatValue: mage.GetStat(stats.SpellCrit) * 2,
		Kind:       core.SpellMod_BonusCrit_Rating,
	})

	mage.FingersOfFrostAura = mage.RegisterAura(core.Aura{
		Label:     "Fingers of Frost Proc",
		ActionID:  core.ActionID{SpellID: 44545},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			fingersOfFrostFrozenCritMod.Activate()
			fingersOfFrostIceLanceDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			fingersOfFrostFrozenCritMod.Deactivate()
			fingersOfFrostIceLanceDamageMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == mage.IceLance || spell == mage.DeepFreeze {
				aura.RemoveStack(sim)
			}
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Fingers of Frost Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			fmt.Println(mage.hasChillEffect(spell))
			if mage.hasChillEffect(spell) && sim.Proc(procChance, "FingersOfFrostProc") {
				mage.FingersOfFrostAura.Activate(sim)
				mage.FingersOfFrostAura.AddStack(sim)
			}
		},
	})
}

func (mage *Mage) applyBrainFreeze() {
	if mage.Talents.BrainFreeze == 0 {
		return
	}

	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)
	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	brainFreezeCostMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellBrainFreeze,
		FloatValue: -1,
		Kind:       core.SpellMod_PowerCost_Pct,
	})

	brainFreezeCastMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellBrainFreeze,
		FloatValue: -1,
		Kind:       core.SpellMod_CastTime_Pct,
	})

	mage.BrainFreezeAura = mage.GetOrRegisterAura(core.Aura{
		Label:    "Brain Freeze Proc",
		ActionID: core.ActionID{SpellID: 57761},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			brainFreezeCostMod.Activate()
			brainFreezeCastMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			brainFreezeCostMod.Deactivate()
			brainFreezeCastMod.Deactivate()
			if t10ProcAura != nil {
				t10ProcAura.Activate(sim)
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == mage.FrostfireBolt || spell == mage.Fireball {
				if !hasT8_4pc || !sim.Proc(T84PcProcChance, "MageT84PC") {
					aura.Deactivate(sim)
				}
			}
		},
	})

	procChance := .05 * float64(mage.Talents.BrainFreeze)
	mage.RegisterAura(core.Aura{
		Label:    "Brain Freeze Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if mage.hasChillEffect(spell) && sim.Proc(procChance, "Brain Freeze") {
				mage.BrainFreezeAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) applyEarlyFrost() {
	if mage.Talents.EarlyFrost == 0 {
		return
	}
	/* earlyWinterMod := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellFrostbolt,
		FloatValue: 0.3 * float64(mage.Talents.EarlyFrost),
	}) */
}

// func (mage *Mage) applyWintersChill() {
// 	if mage.Talents.WintersChill == 0 {
// 		return
// 	}

// 	procChance := []float64{0, 0.33, 0.66, 1}[mage.Talents.WintersChill]

// 	wcAuras := mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
// 		return core.WintersChillAura(target, 0)
// 	})
// 	mage.Env.RegisterPreFinalizeEffect(func() {
// 		for _, spell := range mage.GetSpellsMatchingSchool(core.SpellSchoolFrost) {
// 			spell.RelatedAuras = append(spell.RelatedAuras, wcAuras)
// 		}
// 	})

// 	mage.RegisterAura(core.Aura{
// 		Label:    "Winters Chill Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Landed() || !spell.SpellSchool.Matches(core.SpellSchoolFrost) {
// 				return
// 			}

// 			if sim.Proc(procChance, "Winters Chill") {
// 				aura := wcAuras.Get(result.Target)
// 				aura.Activate(sim)
// 				if aura.IsActive() {
// 					aura.AddStack(sim)
// 				}
// 			}
// 		},
// 	})
// }
