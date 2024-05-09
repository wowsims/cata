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
	mage.applyPyromaniac()

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
			ClassMask:  MageSpellsAll,
			FloatValue: 0.01 * float64(mage.Talents.FirePower),
			Kind:       core.SpellMod_DamageDone_Pct,
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

	// Critical Mass
	if mage.Talents.CriticalMass > 0 {
		mage.AddStaticMod(core.SpellModConfig{
			ClassMask:  MageSpellLivingBomb | MageSpellFlameOrb,
			FloatValue: 0.05 * float64(mage.Talents.CriticalMass),
			Kind:       core.SpellMod_DamageDone_Flat,
		})

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

	ImprovedHotStreakProcChance := float64(mage.Talents.ImprovedHotStreak) * 0.5
	BaseHotStreakProcChance := float64(-2.7*(mage.GetStat(stats.SpellCrit)/core.CritRatingPerCritChance)/100 + 0.9) // EJ settled on -2.7*critChance+0.9

	hotStreakCostMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -1,
		ClassMask:  MageSpellPyroblast,
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
				if sim.Proc(BaseHotStreakProcChance, "Hot Streak") {
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
					if sim.Proc(ImprovedHotStreakProcChance, "Improved Hot Streak") {
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
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyCastSpeed(hasteBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyCastSpeed(1 / hasteBonus)
		},
	})

	mage.RegisterAura(core.Aura{
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
			dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite, mage.PyroblastDot, mage.Combustion}
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
	})
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
	const IgniteTicksFresh = 2

	// Ignite proc listener
	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Ignite Talent",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			// EJ post says combustion crits do not proc ignite
			// https://web.archive.org/web/20120219014159/http://elitistjerks.com/f75/t110187-cataclysm_mage_simulators_formulators/p3/#post1824829
			if spell.ClassSpellMask&(MageSpellLivingBombDot|MageSpellCombustion) == 0 && result.DidCrit() {
				mage.procIgnite(sim, result)
			}
		},
	}))

	// The ignite dot
	mage.Ignite = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 12846},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskProc,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: MageSpellIgnite,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Ignite",
				Tag:       "IgniteDot",
				MaxStacks: 1000000,
			},
			NumberOfTicks:       IgniteTicksFresh,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: false,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {

			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.Spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Hits++
			spell.Dot(target).Apply(sim)
		},
	})
}

func (mage *Mage) procIgnite(sim *core.Simulation, result *core.SpellResult) {
	const IgniteTicksFresh = 2
	const IgniteTicksRefresh = 3
	var currentMastery float64 = 1.224 + 0.028*mage.GetMasteryPoints()

	igniteDamageMultiplier := []float64{0.0, 0.13, 0.26, 0.40}[mage.Talents.Ignite]
	newDamage := result.Damage * igniteDamageMultiplier * currentMastery
	dot := mage.Ignite.Dot(result.Target)

	// Cata Ignite
	// 1st ignite application = 4s, split into 2 ticks (2s, 0s)
	// Ignite refreshes: Duration = 4s + MODULO(remaining duration, 2), max 6s. Split damage over 3 ticks at 4s, 2s, 0s.
	if dot.IsActive() {
		outstandingDamage := dot.SnapshotBaseDamage * float64(dot.NumTicksRemaining(sim))
		dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(IgniteTicksRefresh)
	} else {
		dot.SnapshotBaseDamage = newDamage / IgniteTicksFresh
	}
	mage.Ignite.Cast(sim, result.Target)
	dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
}

func (mage *Mage) applyImpact() {

	if mage.Talents.Impact == 0 {
		return
	}

	// TODO make this work :)
	// Currently casts a fresh set of DoTs
	// afaik it should spread exact copies of the DoTs
	impactAura := mage.RegisterAura(core.Aura{
		Label:    "Impact",
		ActionID: core.ActionID{SpellID: 64343},
		Duration: time.Second * 10,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {

			if spell.ClassSpellMask == MageSpellFireBlast {
				// TODO:
				// originalTarget := mage.CurrentTarget
				// duplicatableDots := map[*core.Spell]float64{
				// 	mage.LivingBombImpact:   mage.LivingBomb.Dot(originalTarget).SnapshotBaseDamage,
				// 	mage.PyroblastDotImpact: mage.PyroblastDot.Dot(originalTarget).SnapshotBaseDamage,
				// 	mage.Ignite:             mage.Ignite.Dot(originalTarget).SnapshotBaseDamage,
				// 	mage.Combustion:         mage.Combustion.Dot(originalTarget).SnapshotBaseDamage,
				// }
				// for _, aoeTarget := range sim.Encounter.TargetUnits {
				// 	if aoeTarget == originalTarget {
				// 		continue
				// 	}
				// 	for spell, damage := range duplicatableDots {
				// 		spell.Dot(aoeTarget).Snapshot(aoeTarget, damage)
				// 		spell.Dot(aoeTarget).Apply(sim)
				// 	}
				// }
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
