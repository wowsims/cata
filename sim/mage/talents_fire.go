package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

//"github.com/wowsims/cata/sim/core/proto"

func (mage *Mage) ApplyFireTalents() {

	// Cooldowns/Special Implementations
	mage.applyIgnite()
	mage.applyImpact()
	mage.applyHotStreak()
	mage.applyMoltenFury()
	mage.applyMasterOfElements()
	//mage.applyPyromaniac()

	// Improved Fire Blast
	if mage.Talents.ImprovedFireBlast > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellFireBlast,
			FloatValue: 4 * float64(mage.Talents.ImprovedFireBlast) * core.CritRatingPerCritChance,
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

func (mage *Mage) applyHotStreak() {
	if !mage.Talents.HotStreak {
		return
	}

	ImprovedHotStreakProcChance := float64(mage.Talents.ImprovedHotStreak) * 0.5
	BaseHotStreakProcChance := float64(-2.7*(mage.GetStat(stats.SpellCrit)/core.CritRatingPerCritChance)/100 + 0.9) // EJ settled on -2.7*critChance+0.9
	// Unimproved Hot Streak Proc Aura
	mage.HotStreakAura = mage.RegisterAura(core.Aura{
		Label:    "Hot Streak",
		ActionID: core.ActionID{SpellID: 48108},
		Duration: time.Second * 10,
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
				// Pyroblast! cannot trigger hot streak
				if spell != mage.Pyroblast && spell.CastTimeMultiplier != 0 {
					if sim.Proc(BaseHotStreakProcChance, "Hot Streak") {
						if mage.HotStreakAura.IsActive() {
							mage.HotStreakAura.Refresh(sim)
						} else {
							mage.HotStreakAura.Activate(sim)
						}
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
				if len(activeFireDots) >= 3 {
					pyromaniacMod.Activate()
				} else {
					pyromaniacMod.Deactivate()
				}
			},
		})*/
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

func (mage *Mage) applyIgnite() {

	if mage.Talents.Ignite == 0 {
		return
	}

	const IgniteTicks = 2
	// Ignite proc listener
	mage.RegisterAura(core.Aura{
		Label:    "Ignite Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			// EJ post says combustion crits do not proc ignite
			// https://web.archive.org/web/20120219014159/http://elitistjerks.com/f75/t110187-cataclysm_mage_simulators_formulators/p3/#post1824829
			if spell.SpellSchool.Matches(core.SpellSchoolFire) && result.DidCrit() && (spell != mage.Combustion || spell != mage.CombustionImpact) {
				mage.procIgnite(sim, result)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if mage.LivingBomb != nil && result.DidCrit() {
				mage.procIgnite(sim, result)
			}
		},
	})

	mage.IgniteDamageTracker = core.MakePermanent(mage.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 413843},
		Label:     "Ignite Damage Tracker",
		Duration:  core.NeverExpires,
		MaxStacks: 1000000,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	}))

	// The ignite dot
	mage.Ignite = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 413843},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskProc,
		Flags:          SpellFlagMage | core.SpellFlagIgnoreModifiers,
		ClassSpellMask: MageSpellIgnite,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Ignite",
				Tag:   "IgniteDot",
			},
			NumberOfTicks: IgniteTicks,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Need to do mastery here
				currentMastery := 1.22 + 0.028*mage.GetMasteryPoints()

				result := dot.Spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage*currentMastery, dot.OutcomeTick)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Hits++
			spell.Dot(target).ApplyOrReset(sim)
		},
	})
}

func (mage *Mage) procIgnite(sim *core.Simulation, result *core.SpellResult) {
	const IgniteTicks = 2
	igniteDamageMultiplier := []float64{0.0, 0.13, 0.26, 0.40}[mage.Talents.Ignite]

	dot := mage.Ignite.Dot(result.Target)

	newDamage := result.Damage * igniteDamageMultiplier

	// if ignite was still active, we store up the remaining damage to be added to the next application
	outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

	dot.SnapshotAttackerMultiplier = 1
	// Add the remaining damage to the new ignite proc, divide it over 2 ticks
	dot.SnapshotBaseDamage = ((outstandingDamage + newDamage) / float64(IgniteTicks)) // * (1.22 + 0.028*mage.GetMasteryPoints()))
	mage.Ignite.Cast(sim, result.Target)
	if mage.IgniteDamageTracker.IsActive() {
		mage.IgniteDamageTracker.SetStacks(sim, int32(dot.SnapshotBaseDamage))
	}
}

func (mage *Mage) applyImpact() {

	if mage.Talents.Impact == 0 {
		return
	}

	mage.ImpactAura = mage.RegisterAura(core.Aura{
		Label:    "Impact",
		ActionID: core.ActionID{SpellID: 64343},
		Duration: time.Second * 10,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {

			if spell == mage.FireBlast {
				originalTarget := mage.CurrentTarget

				duplicatableDots := map[*core.Spell]float64{
					//mage.LivingBombImpact:   mage.LivingBomb.Dot(originalTarget).SnapshotBaseDamage,
					mage.PyroblastDotImpact: mage.PyroblastDot.Dot(originalTarget).SnapshotBaseDamage,
					//mage.Ignite:           mage.Ignite.Dot(originalTarget).SnapshotBaseDamage,
					mage.CombustionImpact: mage.Combustion.Dot(originalTarget).SnapshotBaseDamage,
				}
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					mage.CurrentTarget = aoeTarget
					if mage.CurrentTarget == originalTarget {
						continue
					}
					for spell, damage := range duplicatableDots {
						spell.Dot(aoeTarget).Snapshot(aoeTarget, damage)
						spell.Dot(aoeTarget).Apply(sim)
					}
				}
				mage.CurrentTarget = originalTarget // reset to original target
				aura.Deactivate(sim)
			}
		},
	})

	//TODO make this work :)
	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Impact Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: MageSpellsAll,
		ProcChance:     float64(0.05) * float64(mage.Talents.Impact),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			mage.FireBlast.CD.Reset()
			mage.ImpactAura.Activate(sim)
		},
	})
}
