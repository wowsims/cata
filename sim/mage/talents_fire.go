package mage

import (
	"time"

	"github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

//"github.com/wowsims/mop/sim/core/proto"

func (mage *Mage) ApplyFireTalents() {
	// Cooldowns/Special Implementations
	mage.applyIgnite()
	mage.applyImpact()
	mage.applyHotStreak()
	mage.applyMoltenFury()
	mage.applyMasterOfElements()
	mage.applyPyromaniac()

	// Improved Fire Blast
	if mage.Talents.ImprovedFireBlast > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellFireBlast,
			FloatValue: 4 * float64(mage.Talents.ImprovedFireBlast),
			Kind:       core.SpellMod_BonusCrit_Percent,
		})
	}

	// Fire Power
	if mage.Talents.FirePower > 0 {
		firePowerModConfig := func(classMask int64) core.SpellModConfig {
			return core.SpellModConfig{
				School:     core.SpellSchoolFire,
				ClassMask:  classMask,
				FloatValue: 0.01 * float64(mage.Talents.FirePower),
				Kind:       core.SpellMod_DamageDone_Pct,
			}
		}

		mage.AddStaticMod(firePowerModConfig(MageSpellsAll))
		mage.flameOrb.AddStaticMod(firePowerModConfig(MageSpellFlagNone))
	}

	// Improved Scorch
	if mage.Talents.ImprovedScorch > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellScorch,
			IntValue:  -50 * mage.Talents.ImprovedScorch,
			Kind:      core.SpellMod_PowerCost_Pct,
		})
	}

	// Firestarter
	if mage.Talents.Firestarter {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask: MageSpellScorch,
			Kind:      core.SpellMod_AllowCastWhileMoving,
		})
	}

	// Improved Flamestrike
	if mage.Talents.ImprovedFlamestrike > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellFlamestrike,
			FloatValue: -0.5 * float64(mage.Talents.ImprovedFlamestrike),
			Kind:       core.SpellMod_CastTime_Pct,
		})

		mage.Unit.RegisterSpell(mage.GetFlameStrikeConfig(88148, true))
		core.MakePermanent(core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
			Name:           "Improved Flame Strike - Talent",
			ActionID:       core.ActionID{SpellID: 84674},
			Callback:       core.CallbackOnSpellHitDealt,
			ProcChance:     1,
			Outcome:        core.OutcomeLanded,
			ProcMask:       core.ProcMaskSpellDamage,
			ClassSpellMask: MageSpellBlastWave,
			ICD:            time.Millisecond * 1,
			ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				return len(spell.Unit.Env.Encounter.Targets) >= 2
			},
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				flameStrikeCopy := spell.Unit.GetSpell(core.ActionID{SpellID: 88148, Tag: 1})
				flameStrikeCopy.Cast(sim, result.Target)
			},
		}))
	}

	// Critical Mass
	if mage.Talents.CriticalMass > 0 {
		criticalMassModConfig := func(classMask int64) core.SpellModConfig {
			return core.SpellModConfig{
				ClassMask:  classMask,
				FloatValue: 0.05 * float64(mage.Talents.CriticalMass),
				Kind:       core.SpellMod_DamageDone_Flat,
			}
		}

		mage.AddStaticMod(criticalMassModConfig(MageSpellLivingBomb | MageSpellFlameOrb))
		mage.flameOrb.AddStaticMod(criticalMassModConfig(MageSpellFlameOrb))

		criticalMassDebuff := mage.NewEnemyAuraArray(core.CriticalMassAura)

		core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
			Name:           "Critical Mass Trigger",
			Callback:       core.CallbackOnSpellHitDealt,
			ClassSpellMask: MageSpellPyroblast | MageSpellScorch,
			Outcome:        core.OutcomeLanded,
			ProcChance:     float64(mage.Talents.CriticalMass) / 3.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				criticalMassDebuff.Get(result.Target).Activate(sim)
			},
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

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Master of Elements",
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
	}))
}

func (mage *Mage) applyHotStreak() {
	if !mage.Talents.HotStreak {
		return
	}

	improvedHotStreakProcChance := float64(mage.Talents.ImprovedHotStreak) * 0.5

	// This is the new formula as the old Simcraft / EJ has been debunked by PTR testing.
	calculateHotStreakProcChance := func(x float64) float64 {
		return -2.67*min(x, 0.3402) + 0.9230
	}

	hotStreakCostMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
		ClassMask: MageSpellPyroblast,
	})

	hotStreakCastTimeMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
		ClassMask:  MageSpellPyroblast,
	})

	// Unimproved Hot Streak Proc Aura
	hotStreakAura := mage.RegisterAura(core.Aura{
		Label:    "Hot Streak",
		ActionID: core.ActionID{SpellID: 48108},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hotStreakCostMod.Activate()
			hotStreakCastTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hotStreakCostMod.Deactivate()
			hotStreakCastTimeMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask != MageSpellPyroblast {
				return
			}
			if spell.CurCast.Cost > 0 || spell.CurCast.CastTime > 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	// Improved Hotstreak Crit Stacking Aura
	hotStreakCritAura := mage.RegisterAura(core.Aura{
		Label:     "Hot Streak Proc Aura",
		ActionID:  core.ActionID{SpellID: 44448}, //, Tag: 1}, Removing Tag gets rid of the (??) in Timeline
		MaxStacks: 2,
		Duration:  core.NeverExpires,
	})

	const hotStreakSpells = MageSpellPyroblast | MageSpellFireBlast | MageSpellFireball |
		MageSpellFlameOrb | MageSpellFrostfireBolt | MageSpellScorch

	// Aura to allow the character to track crits
	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Hot Streak Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&hotStreakSpells == 0 {
				return
			}

			// Pyroblast! cannot trigger hot streak
			// TODO can Pyroblast! *reset* hot streak crit streak? This implementation assumes no.
			// If so, will need to envelope it around the hot streak checks
			if spell.ClassSpellMask == MageSpellPyroblast && spell.CurCast.CastTime == 0 {
				return
			}
			// Hot Streak Base Talent Proc
			if result.DidCrit() {
				baseCritPercent := mage.GetStat(stats.SpellCritPercent)
				hotStreakProcChance := mage.baseHotStreakProcChance + calculateHotStreakProcChance(baseCritPercent/100)

				if sim.Proc(hotStreakProcChance, "Hot Streak") {
					hotStreakAura.Activate(sim)
				}
			}

			// Improved Hot Streak
			if mage.Talents.ImprovedHotStreak > 0 {
				// If you didn't crit, reset your crit counter
				if !result.DidCrit() {
					hotStreakCritAura.SetStacks(sim, 0)
					hotStreakCritAura.Deactivate(sim)
					return
				}

				// If you did crit, check against talents to see if you proc
				// If you proc and had 1 stack, set crit counter to 0 and give hot streak.
				if hotStreakCritAura.GetStacks() == 1 {
					if sim.Proc(improvedHotStreakProcChance, "Improved Hot Streak") {
						hotStreakCritAura.SetStacks(sim, 0)
						hotStreakCritAura.Deactivate(sim)

						hotStreakAura.Activate(sim)
					}

					// If you proc and had 0 stacks of crits, add to your crit counter.
					// No idea if 1 out of 2 talent points means you have a 50% chance to
					// add to the 1st stack of crit, or only the 2nd. Doesn't seem
					// all that important to check since every fire mage in the world
					// will go 2 out of 2 points, but worth researching.
					// If it checks 1st crit as well, can add a proc check to this too
				} else {
					hotStreakCritAura.Activate(sim)
					hotStreakCritAura.AddStack(sim)
				}
			}
		},
	}))
}

func (mage *Mage) applyPyromaniac() {
	if mage.Talents.Pyromaniac == 0 {
		return
	}

	hasteBonus := 1.0 + .05*float64(mage.Talents.Pyromaniac)
	pyromaniacAura := mage.GetOrRegisterAura(core.Aura{
		Label:    "Pyromaniac",
		ActionID: core.ActionID{SpellID: 83582},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyCastSpeed(hasteBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyCastSpeed(1 / hasteBonus)
		},
	})

	if len(mage.Env.Encounter.TargetUnits) < 3 {
		return
	}

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label:    "Pyromaniac Trackers",
		ActionID: core.ActionID{SpellID: 83582},
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite, mage.Pyroblast, mage.Combustion}
			activeDotTargets := 0
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				for _, spells := range dotSpells {
					if spells.Dot(aoeTarget).IsActive() {
						activeDotTargets++
						break
					}
				}
			}
			if activeDotTargets >= 3 && !pyromaniacAura.IsActive() {
				pyromaniacAura.Activate(sim)
			} else if activeDotTargets < 3 && pyromaniacAura.IsActive() {
				pyromaniacAura.Deactivate(sim)
			}
		},
	}))
}

func (mage *Mage) applyMoltenFury() {
	if mage.Talents.MoltenFury == 0 {
		return
	}

	moltenFuryMulti := 1.0 + .04*float64(mage.Talents.MoltenFury)

	moltenFuryAuras := mage.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return unit.GetOrRegisterAura(core.Aura{
			Label:    "Molten Fury",
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				mage.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= moltenFuryMulti
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				mage.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= moltenFuryMulti
			},
		})
	})

	mage.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 35 {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					moltenFuryAuras.Get(aoeTarget).Activate(sim)
				}
			}
		})
	})
}

func (mage *Mage) applyIgnite() {
	if mage.Talents.Ignite == 0 {
		return
	}

	igniteDamageMultiplier := []float64{0.0, 0.13, 0.26, 0.40}[mage.Talents.Ignite]

	mage.Ignite = cata.RegisterIgniteEffect(&mage.Unit, cata.IgniteConfig{
		ActionID:       core.ActionID{SpellID: 12846},
		ClassSpellMask: MageSpellIgnite,
		DotAuraLabel:   "Ignite",
		DotAuraTag:     "IgniteDot",

		ProcTrigger: core.ProcTrigger{
			Name:     "Ignite Talent",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellDamage,
			Outcome:  core.OutcomeCrit,

			ExtraCondition: func(_ *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
				if !spell.SpellSchool.Matches(core.SpellSchoolFire) {
					return false
				}

				// EJ post says combustion crits do not proc ignite
				// https://web.archive.org/web/20120219014159/http://elitistjerks.com/f75/t110187-cataclysm_mage_simulators_formulators/p3/#post1824829
				return spell.ClassSpellMask&(MageSpellLivingBombDot|MageSpellCombustion|MageSpellLivingBomb) == 0
			},
		},

		DamageCalculator: func(result *core.SpellResult) float64 {
			var masteryMultiplier float64 = 1 + (22.4+2.8*mage.GetMasteryPoints())/100
			return result.Damage * igniteDamageMultiplier * masteryMultiplier
		},
	})

	// This is needed because we want to listen for the spell "cast" event that refreshes the Dot
	mage.Ignite.Flags ^= core.SpellFlagNoOnCastComplete

}

func (mage *Mage) applyImpact() {
	if mage.Talents.Impact == 0 {
		return
	}

	var duplicatableDots []*core.Spell
	impactAura := mage.RegisterAura(core.Aura{
		Label:    "Impact",
		ActionID: core.ActionID{SpellID: 64343},
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			duplicatableDots = []*core.Spell{mage.LivingBomb, mage.Pyroblast.RelatedDotSpell, mage.Ignite, mage.Combustion.RelatedDotSpell}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask == MageSpellFireBlast {

				originalTarget := mage.CurrentTarget

				for _, aoeTarget := range sim.Encounter.TargetUnits {
					if aoeTarget == originalTarget {
						continue
					}
					for _, dotSpell := range duplicatableDots {
						originaldot := dotSpell.Dot(originalTarget)
						if !originaldot.IsActive() {
							continue
						}

						newdot := dotSpell.Dot(aoeTarget)
						if dotSpell != mage.Ignite {
							newdot.CopyDotAndApply(sim, originaldot) // See attached .go file
						} else {
							// TODO Impact Ignite
						}
					}
				}
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Impact Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: MageSpellsAll,
		ProcChance:     0.05 * float64(mage.Talents.Impact),
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			mage.FireBlast.CD.Reset()
			impactAura.Activate(sim)
		},
	})
}
