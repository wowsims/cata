package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (druid *Druid) RazorClawsMultiplier(masteryRating float64) float64 {
	razorClawsMulti := 1.0

	if druid.Spec == proto.Spec_SpecFeralDruid {
		razorClawsMulti += 0.25 + 0.03125*core.MasteryRatingToMasteryPoints(masteryRating)
	}

	return razorClawsMulti
}

func (druid *Druid) ThickHideMultiplier() float64 {
	thickHideMulti := 1.0

	if druid.Talents.ThickHide > 0 {
		thickHideMulti += 0.04 + 0.03*float64(druid.Talents.ThickHide-1)
	}

	return thickHideMulti
}

func (druid *Druid) BearArmorMultiplier() float64 {
	thickHideBearMulti := 1.0 + 0.26*float64(druid.Talents.ThickHide) // This is a bear-specific multiplier that stacks with the generic multiplier calculated above.
	return 2.2 * thickHideBearMulti
}

func (druid *Druid) ApplyTalents() {
	druid.MultiplyStat(stats.Mana, 1.0+0.05*float64(druid.Talents.Furor))
	// druid.AddStat(stats.SpellHit, float64(druid.Talents.BalanceOfPower)*2*core.SpellHitRatingPerHitChance)
	// druid.AddStat(stats.SpellCrit, float64(druid.Talents.NaturalPerfection)*1*core.CritRatingPerCritChance)
	// druid.PseudoStats.CastSpeedMultiplier *= 1 + (float64(druid.Talents.CelestialFocus) * 0.01)
	// druid.PseudoStats.DamageDealtMultiplier *= 1 + (float64(druid.Talents.EarthAndMoon) * 0.02)
	// druid.PseudoStats.SpiritRegenRateCasting = float64(druid.Talents.Intensity) * (0.5 / 3)
	// druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.02*float64(druid.Talents.Naturalist)
	druid.ApplyEquipScaling(stats.Armor, druid.ThickHideMultiplier())
	druid.PseudoStats.ReducedCritTakenChance += 0.02 * float64(druid.Talents.ThickHide)

	// if druid.Talents.LunarGuidance > 0 {
	// 	bonus := 0.04 * float64(druid.Talents.LunarGuidance)
	// 	druid.AddStatDependency(stats.Intellect, stats.SpellPower, bonus)
	// }

	// if druid.Talents.Dreamstate > 0 {
	// 	bonus := 0.04 * float64(druid.Talents.Dreamstate)
	// 	druid.AddStatDependency(stats.Intellect, stats.MP5, bonus)
	// }

	if druid.Talents.HeartOfTheWild > 0 {
		bonus := 0.02 * float64(druid.Talents.HeartOfTheWild)
		druid.MultiplyStat(stats.Intellect, 1.0+bonus)
	}

	// if druid.Talents.ImprovedFaerieFire > 0 && druid.CurrentTarget.HasAuraWithTag(core.FaerieFireAuraTag) {
	// 	druid.AddStat(stats.SpellCrit, float64(druid.Talents.ImprovedFaerieFire)*1*core.CritRatingPerCritChance)
	// }

	// if druid.Talents.ImprovedMarkOfTheWild > 0 {
	// 	bonus := 0.01 * float64(druid.Talents.ImprovedMarkOfTheWild)
	// 	druid.MultiplyStat(stats.Stamina, 1.0+bonus)
	// 	druid.MultiplyStat(stats.Strength, 1.0+bonus)
	// 	druid.MultiplyStat(stats.Agility, 1.0+bonus)
	// 	druid.MultiplyStat(stats.Intellect, 1.0+bonus)
	// 	druid.MultiplyStat(stats.Spirit, 1.0+bonus)
	// }

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

	// druid.setupNaturesGrace()
	// druid.registerNaturesSwiftnessCD()
	// druid.applyEarthAndMoon()
	// druid.applyMoonkinForm()
	druid.applyPrimalFury()
	// druid.applyEclipse()
	druid.applyLotp()
	// druid.applyPredatoryInstincts()
	druid.applyNaturalReaction()
	// druid.applyOwlkinFrenzy()
	// druid.applyInfectedWounds()
	druid.applyFurySwipes()
	druid.applyPrimalMadness()
}

// func (druid *Druid) setupNaturesGrace() {
// 	if druid.Talents.NaturesGrace == 0 {
// 		return
// 	}

// 	druid.NaturesGraceProcAura = druid.RegisterAura(core.Aura{
// 		Label:    "Natures Grace Proc",
// 		ActionID: core.ActionID{SpellID: 16886},
// 		Duration: time.Second * 3,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.MultiplyCastSpeed(1.2)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			druid.MultiplyCastSpeed(1 / 1.2)
// 		},
// 	})

// 	procChance := []float64{0, .33, .66, 1}[druid.Talents.NaturesGrace]

// 	druid.RegisterAura(core.Aura{
// 		Label:    "Natures Grace",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if !result.Outcome.Matches(core.OutcomeCrit) {
// 				return
// 			}
// 			if spell.Flags.Matches(SpellFlagNaturesGrace) && sim.Proc(procChance, "Natures Grace") {
// 				druid.NaturesGraceProcAura.Activate(sim)
// 			}
// 		},
// 	})
// }

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
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
		DamageMultiplier: 3.1,
		CritMultiplier:   druid.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()), spell.OutcomeMeleeSpecialHitAndCrit)
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
				if druid.IsMangle(spell) || druid.Shred.IsEqual(spell) || druid.Rake.IsEqual(spell) {
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

// Modifies the Bleed aura to apply the bonus.
func (druid *Druid) applyRendAndTear(aura core.Aura) core.Aura {
	if druid.FerociousBite == nil || druid.Talents.RendAndTear == 0 || druid.AssumeBleedActive {
		return aura
	}

	bonusCrit := []float64{0.0, 8.0, 17.0, 25.0}[druid.Talents.RendAndTear] * core.CritRatingPerCritChance

	aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
		if druid.BleedsActive == 0 {
			druid.FerociousBite.BonusCritRating += bonusCrit
		}
		druid.BleedsActive++
	})
	aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		druid.BleedsActive--
		if druid.BleedsActive == 0 {
			druid.FerociousBite.BonusCritRating -= bonusCrit
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
// 	solarProcMultiplier := 1.4 + core.TernaryFloat64(druid.HasSetBonus(ItemSetNightsongGarb, 2), 0.07, 0)
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
// 	lunarBonusCrit := (40 + core.TernaryFloat64(druid.HasSetBonus(ItemSetNightsongGarb, 2), 7, 0)) * core.CritRatingPerCritChance
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
	healthRestore := 0.05

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

// 	onGainMod := druid.MeleeCritMultiplier(Cat)
// 	onExpireMod := druid.MeleeCritMultiplier(Humanoid)

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
	if druid.Talents.NaturalReaction == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 59071}
	rageMetrics := druid.NewRageMetrics(actionID)
	rageAdded := 1.0 + 2.0*float64(druid.Talents.NaturalReaction-1)

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
}

// func (druid *Druid) applyInfectedWounds() {
// 	if druid.Talents.InfectedWounds == 0 {
// 		return
// 	}

// 	iwAuras := druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
// 		return core.InfectedWoundsAura(target, druid.Talents.InfectedWounds)
// 	})
// 	druid.Env.RegisterPreFinalizeEffect(func() {
// 		if druid.Shred != nil {
// 			druid.Shred.RelatedAuras = append(druid.Shred.RelatedAuras, iwAuras)
// 		}
// 		if druid.MangleCat != nil {
// 			druid.MangleCat.RelatedAuras = append(druid.MangleCat.RelatedAuras, iwAuras)
// 		}
// 		if druid.MangleBear != nil {
// 			druid.MangleBear.RelatedAuras = append(druid.MangleBear.RelatedAuras, iwAuras)
// 		}
// 		if druid.Maul != nil {
// 			druid.Maul.RelatedAuras = append(druid.Maul.RelatedAuras, iwAuras)
// 		}
// 	})

// 	druid.RegisterAura(core.Aura{
// 		Label:    "Infected Wounds Talent",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.Landed() && (druid.Shred.IsEqual(spell) || druid.Maul.IsEqual(spell) || druid.MangleCat.IsEqual(spell) || druid.MangleBear.IsEqual(spell)) {
// 				iwAuras.Get(result.Target).Activate(sim)
// 			}
// 		},
// 	})
// }
