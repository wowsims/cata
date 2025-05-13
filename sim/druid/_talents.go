package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (druid *Druid) RazorClawsMultiplier(masteryRating float64) float64 {
	razorClawsMulti := 1.0

	if druid.Spec == proto.Spec_SpecFeralDruid {
		razorClawsMulti += 0.25 + 0.03125*core.MasteryRatingToMasteryPoints(masteryRating)
	}

	return razorClawsMulti
}

func (druid *Druid) applyThickHide() {
	numPoints := druid.Talents.ThickHide

	if numPoints == 0 {
		return
	}

	druid.ApplyEquipScaling(stats.Armor, []float64{1.0, 1.04, 1.07, 1.1}[numPoints])
	druid.PseudoStats.ReducedCritTakenChance += 0.02 * float64(numPoints)

	if !druid.InForm(Bear) {
		return
	}

	// This is a bear-specific multiplier that stacks multiplicatively with
	// both the generic multiplier above and the baseline Bear Form
	// multiplier of 2.2 .
	thickHideBearMulti := 1.0 + 0.26*float64(numPoints)

	druid.BearFormAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		druid.ApplyDynamicEquipScaling(sim, stats.Armor, thickHideBearMulti)
	})

	druid.BearFormAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		druid.RemoveDynamicEquipScaling(sim, stats.Armor, thickHideBearMulti)
	})

	druid.ApplyEquipScaling(stats.Armor, thickHideBearMulti)
}

func (druid *Druid) ApplyTalents() {
	druid.MultiplyStat(stats.Mana, 1.0+0.05*float64(druid.Talents.Furor))
	druid.applyThickHide()
	druid.applyMasterShapeshifter()
	druid.applyHeartOfTheWild()

	// Balance
	druid.applyNaturesGrace()
	druid.applyStarlightWrath()
	druid.applyBalanceOfPower()
	druid.applyNaturesMajesty()
	druid.applyMoonglow()
	druid.applyGenesis()
	druid.applyEuphoria()
	druid.applyShootingStars()
	druid.applyGaleWinds()
	druid.applyEarthAndMoon()
	druid.applyMoonkinForm()
	druid.applyLunarShower()
	// if druid.Talents.PrimalPrecision > 0 {
	// 	druid.AddStat(stats.Expertise, 5.0*float64(druid.Talents.PrimalPrecision)*core.ExpertisePerQuarterPercentReduction)
	// }

	// if druid.Talents.LivingSpirit > 0 {
	// 	bonus := 0.05 * float64(druid.Talents.LivingSpirit)
	// 	druid.MultiplyStat(stats.Spirit, 1.0+bonus)
	// }

	if druid.Talents.Perseverance > 0 {
		multiplier := 1.0 - 0.02*float64(druid.Talents.Perseverance)
		druid.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= multiplier
		druid.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= multiplier
		druid.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= multiplier
		druid.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= multiplier
		druid.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= multiplier
		druid.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= multiplier
	}

	// druid.registerNaturesSwiftnessCD()
	druid.applyFeralSwiftness()
	druid.applyPrimalFury()
	druid.applyLotp()
	// druid.applyPredatoryInstincts()
	druid.applyNaturalReaction()
	// druid.applyOwlkinFrenzy()
	druid.applyInfectedWounds()
	druid.applyFurySwipes()
	druid.applyPrimalMadness()
	druid.applyStampede()
	druid.ApplyGlyphs()
}

func (druid *Druid) applyHeartOfTheWild() {
	if druid.Talents.HeartOfTheWild == 0 {
		return
	}

	multiplier := 1.0 + 0.02*float64(druid.Talents.HeartOfTheWild)
	druid.MultiplyStat(stats.Intellect, multiplier)

	if druid.CatFormAura != nil {
		druid.HotWCatDep = druid.NewDynamicMultiplyStat(stats.AttackPower, []float64{1.0, 1.03, 1.07, 1.1}[druid.Talents.HeartOfTheWild])

		if druid.InForm(Cat) {
			druid.StatDependencyManager.EnableDynamicStatDep(druid.HotWCatDep)
		}
	}

	if druid.BearFormAura != nil {
		druid.HotWBearDep = druid.NewDynamicMultiplyStat(stats.Stamina, multiplier)

		if druid.InForm(Bear) {
			druid.StatDependencyManager.EnableDynamicStatDep(druid.HotWBearDep)
		}
	}
}

func (druid *Druid) applyNaturesGrace() {
	if druid.Talents.NaturesGrace == 0 {
		return
	}

	ngAuraSpellId := []int32{0, 16880, 61345, 61345}[druid.Talents.NaturesGrace]
	ngAuraSpellHastePct := []float64{0, 0.05, 0.1, 0.15}[druid.Talents.NaturesGrace]

	ngAura := druid.RegisterAura(core.Aura{
		Label:    "Natures Grace Proc",
		ActionID: core.ActionID{SpellID: ngAuraSpellId},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyCastSpeed(1 + ngAuraSpellHastePct)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyCastSpeed(1 / (1 + ngAuraSpellHastePct))
		},
	})

	ngTrigger := core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:           "Natures Grace (Talent)",
		Callback:       core.CallbackOnSpellHitDealt,
		ProcMask:       core.ProcMaskSpellDamage,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: DruidSpellMoonfire | DruidSpellSunfire | DruidSpellInsectSwarm,
		ICD:            time.Second * 60,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			ngAura.Activate(sim)
		},
	})

	if druid.HasEclipseBar() {
		druid.AddEclipseCallback(func(_ Eclipse, gained bool, sim *core.Simulation) {
			if gained {
				ngTrigger.Icd.Reset()
			}
		})
	}
}

func (druid *Druid) applyBalanceOfPower() {
	if druid.Talents.BalanceOfPower > 0 {
		druid.AddStaticMod(core.SpellModConfig{
			School:     core.SpellSchoolArcane | core.SpellSchoolNature,
			FloatValue: 0.01 * float64(druid.Talents.BalanceOfPower),
			Kind:       core.SpellMod_DamageDone_Pct,
		})

		druid.AddStat(stats.SpellHitPercent, -0.5*float64(druid.Talents.BalanceOfPower)*druid.GetBaseStats()[stats.Spirit]/core.SpellHitRatingPerHitPercent)
		druid.AddStatDependency(stats.Spirit, stats.SpellHitPercent, 0.5*float64(druid.Talents.BalanceOfPower)/core.SpellHitRatingPerHitPercent)
	}
}

func (druid *Druid) applyStarlightWrath() {
	if druid.Talents.StarlightWrath > 0 {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellStarfire | DruidSpellWrath,
			TimeValue: time.Millisecond * time.Duration([]int{0, -150, -250, -500}[druid.Talents.StarlightWrath]),
			Kind:      core.SpellMod_CastTime_Flat,
		})
	}
}

func (druid *Druid) applyNaturesMajesty() {
	if druid.Talents.NaturesMajesty > 0 {
		druid.AddStat(stats.SpellCritPercent, 2*float64(druid.Talents.NaturesMajesty))
	}
}

func (druid *Druid) applyMoonglow() {
	if druid.Talents.Moonglow > 0 {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidDamagingSpells | DruidHealingSpells,
			IntValue:  -3 * druid.Talents.Moonglow,
			Kind:      core.SpellMod_PowerCost_Pct,
		})
	}
}

func (druid *Druid) applyGenesis() {
	if druid.Talents.Genesis > 0 {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellMoonfireDoT | DruidSpellSunfireDoT | DruidSpellInsectSwarm,
			IntValue:  1 * druid.Talents.Genesis,
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
		})

		// TODO: periodic healing spells
		// TODO: swiftmend
	}
}

func (druid *Druid) applyEuphoria() {
	if druid.Talents.Euphoria == 0 {
		return
	}

	euphoriaSpellId := []int32{0, 81061, 81062}[druid.Talents.Euphoria]
	euphoriaProcChancePct := []float64{0, 0.12, 0.24}[druid.Talents.Euphoria]
	euphoriaManaGainPct := []float64{0, 0.08, 0.16}[druid.Talents.Euphoria]

	// Mana return
	manaMetrics := druid.NewManaMetrics(core.ActionID{
		SpellID: euphoriaSpellId,
	})

	if druid.HasEclipseBar() {
		druid.AddEclipseCallback(func(_ Eclipse, gained bool, sim *core.Simulation) {
			if gained {
				druid.AddMana(sim, druid.MaxMana()*euphoriaManaGainPct, manaMetrics)
			}
		})
	}

	// Double eclipse energy
	euphoriaDummyAura := druid.GetOrRegisterAura(core.Aura{
		Label:    "Euphoria Dummy Aura",
		ActionID: core.ActionID{SpellID: euphoriaSpellId},
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {

			aura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:           "Euphoria",
		Callback:       core.CallbackOnApplyEffects,
		ProcChance:     euphoriaProcChancePct,
		ClassSpellMask: DruidSpellWrath | DruidSpellStarfire,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			euphoriaDummyAura.Activate(sim)
		},
	})
}

func (druid *Druid) applyMoonkinForm() {
	if !druid.Talents.MoonkinForm {
		return
	}

	druid.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane | core.SpellSchoolNature,
		FloatValue: 0.1,
		Kind:       core.SpellMod_DamageDone_Pct,
	})
}

func (druid *Druid) applyMasterShapeshifter() {
	if !druid.Talents.MasterShapeshifter {
		return
	}

	if druid.InForm(Moonkin) {
		druid.AddStaticMod(core.SpellModConfig{
			School:     core.SpellSchoolArcane | core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolHoly | core.SpellSchoolNature | core.SpellSchoolShadow,
			FloatValue: 0.04,
			Kind:       core.SpellMod_DamageDone_Pct,
		})
	}

	if druid.CatFormAura != nil {
		druid.CatFormAura.AttachStatBuff(stats.PhysicalCritPercent, 4)
	}

	if druid.BearFormAura != nil {
		druid.BearFormAura.AttachMultiplicativePseudoStatBuff(&druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical], 1.04)
	}
}

func (druid *Druid) applyFeralSwiftness() {
	if druid.Talents.FeralSwiftness == 0 {
		return
	}

	dodgeBonus := 0.02 * float64(druid.Talents.FeralSwiftness)

	if druid.CatFormAura != nil {
		druid.CatFormAura.AttachAdditivePseudoStatBuff(&druid.PseudoStats.BaseDodgeChance, dodgeBonus)
	}

	if druid.BearFormAura != nil {
		druid.BearFormAura.AttachAdditivePseudoStatBuff(&druid.PseudoStats.BaseDodgeChance, dodgeBonus)
	}
}

func (druid *Druid) applyShootingStars() {
	if druid.Talents.ShootingStars == 0 {
		return
	}

	ssCastTimeMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask:  DruidSpellStarsurge,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
	})

	ssAuraSpellId := []int32{0, 93398, 93399}[druid.Talents.ShootingStars]

	ssAura := druid.RegisterAura(core.Aura{
		Label:    "Shooting Stars Proc",
		ActionID: core.ActionID{SpellID: ssAuraSpellId},
		Duration: time.Second * 12,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask != DruidSpellStarsurge {
				return
			}

			ssCastTimeMod.Deactivate()
			aura.Deactivate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ssCastTimeMod.Activate()
			druid.Starsurge.CD.Reset()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ssCastTimeMod.Deactivate()
		},
	})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:           "Shooting Stars (Talent)",
		Callback:       core.CallbackOnPeriodicDamageDealt,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.02 * float64(druid.Talents.ShootingStars),
		ICD:            time.Second * 6,
		ClassSpellMask: DruidSpellInsectSwarm | DruidSpellMoonfireDoT,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			ssAura.Activate(sim)
		},
	})
}

func (druid *Druid) applyGaleWinds() {
	if druid.Talents.GaleWinds > 0 {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask:  DruidSpellTyphoon | DruidSpellHurricane,
			FloatValue: float64(0.15 * float64(druid.Talents.GaleWinds)),
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}
}

func (druid *Druid) applyEarthAndMoon() {
	if druid.Talents.EarthAndMoon {
		druid.AddStaticMod(core.SpellModConfig{
			FloatValue: 0.02,
			Kind:       core.SpellMod_DamageDone_Pct,
		})
	}
}

func (druid *Druid) applyLunarShower() {
	if druid.Talents.LunarShower == 0 {
		return
	}

	lunarShowerDmgMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask: DruidSpellMoonfire | DruidSpellSunfire,
		Kind:      core.SpellMod_DamageDone_Pct,
	})

	lunarShowerResourceMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask: DruidSpellMoonfire | DruidSpellSunfire,
		Kind:      core.SpellMod_PowerCost_Pct,
	})

	var lunarShowerAura = druid.RegisterAura(core.Aura{
		Label:     "Lunar Shower",
		Duration:  time.Second * 3,
		ActionID:  core.ActionID{SpellID: 33603},
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			lunarShowerDmgMod.UpdateFloatValue(float64(aura.GetStacks()) * 0.15)
			lunarShowerDmgMod.Activate()

			lunarShowerResourceMod.UpdateIntValue(aura.GetStacks() * -10)
			lunarShowerResourceMod.Activate()

			// While under the effects of Lunar Shower, Moonfire and Sunfire generate 8 eclipse energy
			druid.SetSpellEclipseEnergy(DruidSpellMoonfire, MoonfireLunarShowerEnergyGain, MoonfireLunarShowerEnergyGain)
			druid.SetSpellEclipseEnergy(DruidSpellSunfire, SunfireLunarShowerEnergyGain, SunfireLunarShowerEnergyGain)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			lunarShowerDmgMod.Deactivate()
			lunarShowerResourceMod.Deactivate()

			druid.SetSpellEclipseEnergy(DruidSpellMoonfire, MoonfireBaseEnergyGain, MoonfireBaseEnergyGain)
			druid.SetSpellEclipseEnergy(DruidSpellSunfire, SunfireBaseEnergyGain, SunfireBaseEnergyGain)
		},
	})

	druid.RegisterAura(core.Aura{
		Label:    "Lunar Shower Handler",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask != DruidSpellMoonfire && spell.ClassSpellMask != DruidSpellSunfire {
				return
			}

			// does not proc off procs
			if spell.ProcMask.Matches(core.ProcMaskProc) {
				return
			}

			if lunarShowerAura.IsActive() {
				if lunarShowerAura.GetStacks() < 3 {
					lunarShowerAura.AddStack(sim)
					lunarShowerAura.Refresh(sim)
				}
			} else {
				lunarShowerAura.Activate(sim)
				lunarShowerAura.SetStacks(sim, 1)
			}
		},
	})
}

// func (druid *Druid) registerNaturesSwiftnessCD() {
// 	if !druid.Talents.NaturesSwiftness {
// 		return
// 	}
// 	actionID := core.ActionID{SpellID: 17116}

// 	var nsAura *core.Aura
// 	nsSpell := druid.RegisterSpell(Humanoid|Moonkin|Tree, core.SpellConfig{
// 		ActionID: actionID,
// 		Flags:    core.SpellFlagNoOnCastComplete,
// 		Cast: core.CastConfig{
// 			CD: core.Cooldown{
// 				Timer:    druid.NewTimer(),
// 				Duration: time.Minute * 3,
// 			},
// 		},
// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			nsAura.Activate(sim)
// 		},
// 	})

// 	nsAura = druid.RegisterAura(core.Aura{
// 		Label:    "Natures Swiftness",
// 		ActionID: actionID,
// 		Duration: core.NeverExpires,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			if druid.Starfire != nil {
// 				druid.Starfire.CastTimeMultiplier -= 1
// 			}
// 			if druid.Wrath != nil {
// 				druid.Wrath.CastTimeMultiplier -= 1
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			if druid.Starfire != nil {
// 				druid.Starfire.CastTimeMultiplier += 1
// 			}
// 			if druid.Wrath != nil {
// 				druid.Wrath.CastTimeMultiplier += 1
// 			}
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if !druid.Wrath.IsEqual(spell) && !druid.Starfire.IsEqual(spell) {
// 				return
// 			}

// 			// Remove the buff and put skill on CD
// 			aura.Deactivate(sim)
// 			nsSpell.CD.Use(sim)
// 			druid.UpdateMajorCooldowns()
// 		},
// 	})

// 	druid.AddMajorCooldown(core.MajorCooldown{
// 		Spell: nsSpell.Spell,
// 		Type:  core.CooldownTypeDPS,
// 		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
// 			// Don't use NS unless we're casting a full-length starfire or wrath.
// 			return !character.HasTemporarySpellCastSpeedIncrease()
// 		},
// 	})
// }

// func (druid *Druid) applyEarthAndMoon() {
// 	if druid.Talents.EarthAndMoon == 0 {
// 		return
// 	}

// 	eamAuras := druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
// 		return core.EarthAndMoonAura(target, druid.Talents.EarthAndMoon)
// 	})
// 	druid.Env.RegisterPreFinalizeEffect(func() {
// 		if druid.Starfire != nil {
// 			druid.Starfire.RelatedAuras = append(druid.Starfire.RelatedAuras, eamAuras)
// 		}
// 		if druid.Wrath != nil {
// 			druid.Wrath.RelatedAuras = append(druid.Wrath.RelatedAuras, eamAuras)
// 		}
// 	})

// 	earthAndMoonSpell := druid.RegisterSpell(Any, core.SpellConfig{
// 		ActionID: core.ActionID{SpellID: 60432},
// 		ProcMask: core.ProcMaskSuppressedProc,
// 		Flags:    core.SpellFlagNoLogs,
// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			eamAuras.Get(target).Activate(sim)
// 		},
// 	})

// 	druid.RegisterAura(core.Aura{
// 		Label:    "Earth And Moon Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.Landed() && (druid.Starfire.IsEqual(spell) || druid.Wrath.IsEqual(spell)) {
// 				earthAndMoonSpell.Cast(sim, result.Target)
// 			}
// 		},
// 	})
// }

func (druid *Druid) applyFurySwipes() {
	if druid.Talents.FurySwipes == 0 {
		return
	}

	furySwipesSpell := druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 80861},
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		Flags:            core.SpellFlagMeleeMetrics,
		DamageMultiplier: 3.1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()), spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:       "Fury Swipes Trigger",
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
		Harmful:    true,
		ProcChance: 0.05 * float64(druid.Talents.FurySwipes),
		ICD:        time.Second * 3,

		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			if druid.InForm(Cat | Bear) {
				furySwipesSpell.Cast(sim, druid.CurrentTarget)
			}
		},
	})
}

func (druid *Druid) applyPrimalFury() {
	if druid.Talents.PrimalFury == 0 {
		return
	}

	procChance := []float64{0, 0.5, 1}[druid.Talents.PrimalFury]
	actionID := core.ActionID{SpellID: 37117}
	rageMetrics := druid.NewRageMetrics(actionID)
	cpMetrics := druid.NewComboPointMetrics(actionID)

	druid.RegisterAura(core.Aura{
		Label:    "Primal Fury",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if druid.InForm(Bear) {
				if result.Outcome.Matches(core.OutcomeCrit) {
					if sim.Proc(procChance, "Primal Fury") {
						druid.AddRage(sim, 5, rageMetrics)
					}
				}
			} else if druid.InForm(Cat) {
				if druid.IsMangle(spell) || druid.Shred.IsEqual(spell) || druid.Rake.IsEqual(spell) || druid.Ravage.IsEqual(spell) {
					if result.Outcome.Matches(core.OutcomeCrit) {
						if sim.Proc(procChance, "Primal Fury") {
							druid.AddComboPoints(sim, 1, cpMetrics)
						}
					}
				}
			}
		},
	})
}

func (druid *Druid) applyPrimalMadness() {
	if (druid.Talents.PrimalMadness == 0) || !druid.InForm(Cat|Bear) {
		return
	}

	actionID := core.ActionID{SpellID: 80315 + druid.Talents.PrimalMadness}
	druid.PrimalMadnessRageMetrics = druid.NewRageMetrics(actionID)

	if !druid.InForm(Cat) {
		return
	}

	energyMetrics := druid.NewEnergyMetrics(actionID)
	energyGain := 10.0 * float64(druid.Talents.PrimalMadness)

	druid.PrimalMadnessAura = druid.RegisterAura(core.Aura{
		Label:    "Primal Madness",
		ActionID: actionID,
		Duration: core.NeverExpires, // duration is tied to Tiger's Fury / Berserk durations
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.UpdateMaxEnergy(sim, energyGain, energyMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.UpdateMaxEnergy(sim, -energyGain, energyMetrics)
		},
	})
}

func (druid *Druid) applyStampede() {
	if (druid.Talents.Stampede == 0) || !druid.InForm(Cat|Bear) {
		return
	}

	bearHasteMod := 1.0 + 0.15*float64(druid.Talents.Stampede)

	druid.StampedeBearAura = druid.RegisterAura(core.Aura{
		Label:    "Stampede (Bear)",
		ActionID: core.ActionID{SpellID: 81015 + druid.Talents.Stampede},
		Duration: time.Second * 8,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyAttackSpeed(sim, bearHasteMod)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.MultiplyAttackSpeed(sim, 1.0/bearHasteMod)
		},
	})

	if !druid.InForm(Cat) {
		return
	}

	ravageCostMod := 50 * druid.Talents.Stampede

	druid.StampedeCatAura = druid.RegisterAura(core.Aura{
		Label:    "Stampede (Cat)",
		ActionID: core.ActionID{SpellID: 81020 + druid.Talents.Stampede},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.Ravage.Cost.PercentModifier -= ravageCostMod
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.Ravage.Cost.PercentModifier += ravageCostMod
		},
	})
}

// Modifies the Bleed aura to apply the bonus.
func (druid *Druid) applyRendAndTear(aura core.Aura) core.Aura {
	if druid.FerociousBite == nil || druid.Talents.RendAndTear == 0 || druid.AssumeBleedActive {
		return aura
	}

	bonusCritPercent := []float64{0.0, 8.0, 17.0, 25.0}[druid.Talents.RendAndTear]

	aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		if druid.BleedsActive == 0 {
			druid.FerociousBite.BonusCritPercent += bonusCritPercent
		}
		druid.BleedsActive++
	})
	aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		druid.BleedsActive--
		if druid.BleedsActive == 0 {
			druid.FerociousBite.BonusCritPercent -= bonusCritPercent
		}
	})

	return aura
}

// func (druid *Druid) applyEclipse() {
// 	druid.SolarICD = core.Cooldown{Timer: druid.NewTimer(), Duration: 0}
// 	druid.LunarICD = core.Cooldown{Timer: druid.NewTimer(), Duration: 0}
// 	if druid.Talents.Eclipse == 0 {
// 		return
// 	}

// 	// Delay between eclipses
// 	eclipseDuration := time.Millisecond * 15000
// 	interEclipseDelay := eclipseDuration - time.Millisecond*500

// 	// Solar
// 	solarProcChance := (1.0 / 3.0) * float64(druid.Talents.Eclipse)
// 	solarProcMultiplier := 1.4 + core.TernaryFloat64(druid.CouldHaveSetBonus(ItemSetNightsongGarb, 2), 0.07, 0)
// 	druid.SolarICD.Duration = time.Millisecond * 30000
// 	druid.SolarEclipseProcAura = druid.RegisterAura(core.Aura{
// 		Icd:      &druid.SolarICD,
// 		Label:    "Solar Eclipse proc",
// 		Duration: eclipseDuration,
// 		ActionID: core.ActionID{SpellID: 48517},
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.Wrath.DamageMultiplier *= solarProcMultiplier
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.Wrath.DamageMultiplier /= solarProcMultiplier
// 		},
// 	})

// 	druid.RegisterAura(core.Aura{
// 		Label:    "Eclipse (Solar)",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Outcome.Matches(core.OutcomeCrit) {
// 				return
// 			}
// 			if !druid.Starfire.IsEqual(spell) {
// 				return
// 			}
// 			if !druid.SolarICD.Timer.IsReady(sim) {
// 				return
// 			}
// 			if druid.LunarICD.Timer.TimeToReady(sim) > interEclipseDelay {
// 				return
// 			}
// 			if sim.RandomFloat("Eclipse (Solar)") < solarProcChance {
// 				druid.SolarICD.Use(sim)
// 				druid.SolarEclipseProcAura.Activate(sim)
// 			}
// 		},
// 	})

// 	// Lunar
// 	lunarProcChance := 0.2 * float64(druid.Talents.Eclipse)
// 	lunarBonusCrit := (40 + core.TernaryFloat64(druid.CouldHaveSetBonus(ItemSetNightsongGarb, 2), 7, 0)) * core.CritRatingPerCritChance
// 	druid.LunarICD.Duration = time.Millisecond * 30000
// 	druid.LunarEclipseProcAura = druid.RegisterAura(core.Aura{
// 		Icd:      &druid.LunarICD,
// 		Label:    "Lunar Eclipse proc",
// 		Duration: eclipseDuration,
// 		ActionID: core.ActionID{SpellID: 48518},
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.Starfire.BonusCritRating += lunarBonusCrit
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.Starfire.BonusCritRating -= lunarBonusCrit
// 		},
// 	})
// 	druid.RegisterAura(core.Aura{
// 		Label:    "Eclipse (Lunar)",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Outcome.Matches(core.OutcomeCrit) {
// 				return
// 			}
// 			if !druid.Wrath.IsEqual(spell) {
// 				return
// 			}
// 			if !druid.LunarICD.Timer.IsReady(sim) {
// 				return
// 			}
// 			if druid.SolarICD.Timer.TimeToReady(sim) > interEclipseDelay {
// 				return
// 			}
// 			if sim.RandomFloat("Eclipse (Lunar)") < lunarProcChance {
// 				druid.LunarICD.Use(sim)
// 				druid.LunarEclipseProcAura.Activate(sim)
// 			}
// 		},
// 	})
// }

// func (druid *Druid) applyOwlkinFrenzy() {
// 	if druid.Talents.OwlkinFrenzy == 0 {
// 		return
// 	}

// 	druid.OwlkinFrenzyAura = druid.RegisterAura(core.Aura{
// 		Label:    "Owlkin Frenzy proc",
// 		ActionID: core.ActionID{SpellID: 48393},
// 		Duration: time.Second * 10,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.PseudoStats.DamageDealtMultiplier *= 1.1
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.PseudoStats.DamageDealtMultiplier /= 1.1
// 		},
// 	})
// }

func (druid *Druid) applyLotp() {
	if !druid.Talents.LeaderOfThePack {
		return
	}

	actionID := core.ActionID{SpellID: 17007}
	manaMetrics := druid.NewManaMetrics(actionID)
	healthMetrics := druid.NewHealthMetrics(actionID)
	manaRestore := 0.08
	healthRestore := 0.04 // Tooltip says 5%, but only healing for 4% in-game

	icd := core.Cooldown{
		Timer:    druid.NewTimer(),
		Duration: time.Second * 6,
	}

	druid.RegisterAura(core.Aura{
		Icd:      &icd,
		Label:    "Improved Leader of the Pack",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) || !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if !druid.InForm(Cat | Bear) {
				return
			}
			icd.Use(sim)
			druid.AddMana(sim, druid.MaxMana()*manaRestore, manaMetrics)
			druid.GainHealth(sim, druid.MaxHealth()*healthRestore, healthMetrics)
		},
	})
}

// func (druid *Druid) applyPredatoryInstincts() {
// 	if druid.Talents.PredatoryInstincts == 0 {
// 		return
// 	}

// 	onGainMod := druid.CritMultiplier(Cat)
// 	onExpireMod := druid.CritMultiplier(Humanoid)

// 	druid.PredatoryInstinctsAura = druid.RegisterAura(core.Aura{
// 		Label:    "Predatory Instincts",
// 		Duration: core.NeverExpires,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.Lacerate.CritMultiplier = onGainMod
// 			druid.Rip.CritMultiplier = onGainMod
// 			druid.Rake.CritMultiplier = onGainMod
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.Lacerate.CritMultiplier = onExpireMod
// 			druid.Rip.CritMultiplier = onExpireMod
// 			druid.Rake.CritMultiplier = onExpireMod
// 		},
// 	})
// }

func (druid *Druid) applyNaturalReaction() {
	if (druid.Talents.NaturalReaction == 0) || (druid.BearFormAura == nil) {
		return
	}

	actionID := core.ActionID{SpellID: 59071}
	rageMetrics := druid.NewRageMetrics(actionID)
	numPoints := float64(druid.Talents.NaturalReaction)
	rageAdded := 1.0 + 2.0*(numPoints-1.0)

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:     "Natural Reaction Trigger",
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			if druid.InForm(Bear) && result.Outcome.Matches(core.OutcomeDodge) {
				druid.AddRage(sim, rageAdded, rageMetrics)
			}
		},
	})

	druid.BearFormAura.AttachMultiplicativePseudoStatBuff(&druid.PseudoStats.DamageTakenMultiplier, 1.0-0.09*numPoints)
	druid.BearFormAura.AttachAdditivePseudoStatBuff(&druid.PseudoStats.BaseDodgeChance, 0.03*numPoints)
}

func (druid *Druid) applyInfectedWounds() {
	if druid.Talents.InfectedWounds == 0 {
		return
	}

	iwAuras := druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.InfectedWoundsAura(target, druid.Talents.InfectedWounds)
	})
	druid.Env.RegisterPreFinalizeEffect(func() {
		triggeringSpells := []*DruidSpell{druid.Shred, druid.Ravage, druid.Maul, druid.MangleCat, druid.MangleBear}

		for _, spell := range triggeringSpells {
			if spell != nil {
				spell.RelatedAuraArrays = spell.RelatedAuraArrays.Append(iwAuras)
			}
		}
	})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:           "Infected Wounds Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: DruidSpellShred | DruidSpellRavage | DruidSpellMaul | DruidSpellMangle,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			iwAuras.Get(result.Target).Activate(sim)
		},
	})
}
