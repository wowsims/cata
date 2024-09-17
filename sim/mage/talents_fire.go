package mage

import (
	"math"
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
			FloatValue: 4 * float64(mage.Talents.ImprovedFireBlast),
			Kind:       core.SpellMod_BonusCrit_Percent,
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

	// Simcraft uses a reference from ElitistJerks that's no longer available, but their formula is
	// max(0, -2.73 * player crit + 0.95)
	// https://web.archive.org/web/20120208064232/http://elitistjerks.com/f75/t110326-cataclysm_fire_mage_compendium/p6/#post1831143 or
	// https://web.archive.org/web/20120208064232/http://elitistjerks.com/f75/t110326-cataclysm_fire_mage_compendium/p6/#post1831207
	baseCritPercent := mage.GetStat(stats.SpellCritPercent) + (mage.GetStat(stats.CritRating) / core.CritRatingPerCritPercent) + 1*float64(mage.Talents.PiercingIce)
	mage.hotStreakProcChance = max(0, float64(-2.7*baseCritPercent/100+0.9)) // EJ settled on -2.7*critChance+0.9

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
				if sim.Proc(mage.hotStreakProcChance, "Hot Streak") {
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

	// Ignite proc listener
	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
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

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			mage.procIgnite(sim, result)
		},
	})

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
				MaxStacks: math.MaxInt32,
			},
			NumberOfTicks:       2,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.Spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}

func (mage *Mage) procIgnite(sim *core.Simulation, result *core.SpellResult) {
	var currentMastery float64 = 1 + math.Floor(22.4+2.8*mage.GetMasteryPoints())/100
	igniteDamageMultiplier := []float64{0.0, 0.13, 0.26, 0.40}[mage.Talents.Ignite]
	newDamage := result.Damage * igniteDamageMultiplier * currentMastery
	dot := mage.Ignite.Dot(result.Target)
	outstandingDamage := dot.OutstandingDmg()
	totalDamage := outstandingDamage + newDamage

	// Cata Ignite
	// 1st ignite application = 4s, split into 2 ticks (2s, 0s)
	// Ignite refreshes: Duration = 4s + MODULO(remaining duration, 2), max 6s. Split damage over 3 ticks at 4s, 2s, 0s.
	newTickCount := dot.BaseTickCount + core.TernaryInt32(dot.IsActive(), 1, 0)
	dot.SnapshotBaseDamage = totalDamage / float64(newTickCount)
	mage.Ignite.Cast(sim, result.Target)
	dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
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
