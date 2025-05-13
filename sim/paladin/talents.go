package paladin

import (
	"slices"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) ApplyTalents() {
	if paladin.Level >= 15 {
		paladin.registerSpeedOfLight()
		paladin.registerLongArmOfTheLaw()
		paladin.registerPursuitOfJustice()
	}

	// Level 30 talents are just CC

	if paladin.Level >= 45 {
		paladin.registerSelflessHealer()
		// Eternal Flame handled in word_of_glory.go
		paladin.registerSacredShield()
	}

	if paladin.Level >= 60 {
		paladin.registerHandOfPurity()
		paladin.registerUnbreakableSpirit()
		// Skipping Clemecy
	}

	if paladin.Level >= 75 {
		paladin.registerHolyAvenger()
		paladin.registerSanctifiedWrath()
		paladin.registerDivinePurpose()
	}

	if paladin.Level >= 90 {
		paladin.registerHolyPrism()
		paladin.registerLightsHammer()
		paladin.registerExecutionSentence()
	}
}

func (paladin *Paladin) registerSpeedOfLight() {
	if !paladin.Talents.SpeedOfLight {
		return
	}

	actionID := core.ActionID{SpellID: 85499}
	speedOfLightAura := paladin.RegisterAura(core.Aura{
		Label:    "Speed of Light" + paladin.Label,
		ActionID: actionID,
		Duration: time.Second * 8,
	})
	speedOfLightAura.NewMovementSpeedEffect(0.7)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3.5,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 45,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: speedOfLightAura,
	})
}

func (paladin *Paladin) registerLongArmOfTheLaw() {
	if !paladin.Talents.LongArmOfTheLaw {
		return
	}

	longArmOfTheLawAura := paladin.RegisterAura(core.Aura{
		Label:    "Long Arm of the Law" + paladin.Label,
		ActionID: core.ActionID{SpellID: 87173},
		Duration: time.Second * 3,
	})
	longArmOfTheLawAura.NewMovementSpeedEffect(0.45)

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Long Arm of the Law Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 87172},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskJudgment,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			longArmOfTheLawAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) registerPursuitOfJustice() {
	if !paladin.Talents.PursuitOfJustice {
		return
	}

	speedLevels := []float64{0.0, 0.15, 0.20, 0.25, 0.30}

	var movementSpeedEffect *core.ExclusiveEffect
	pursuitOfJusticeAura := paladin.RegisterAura(core.Aura{
		Label:     "Pursuit of Justice" + paladin.Label,
		ActionID:  core.ActionID{SpellID: 114695},
		Duration:  core.NeverExpires,
		MaxStacks: 4,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
			aura.SetStacks(sim, 1)
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			paladin.MultiplyMovementSpeed(sim, 1.0/(1+speedLevels[oldStacks]))

			newSpeed := speedLevels[newStacks]
			paladin.MultiplyMovementSpeed(sim, 1+newSpeed)
			movementSpeedEffect.SetPriority(sim, newSpeed)
		},
	})

	movementSpeedEffect = pursuitOfJusticeAura.NewExclusiveEffect("MovementSpeed", true, core.ExclusiveEffect{
		Priority: speedLevels[1],
	})

	paladin.HolyPower.RegisterOnGain(func(gain, realGain int32, actionID core.ActionID, sim *core.Simulation) {
		pursuitOfJusticeAura.Activate(sim)
		pursuitOfJusticeAura.SetStacks(sim, paladin.SpendableHolyPower()+1)
	})
	paladin.HolyPower.RegisterOnSpend(func(amount int32, actionID core.ActionID, sim *core.Simulation) {
		pursuitOfJusticeAura.Activate(sim)
		pursuitOfJusticeAura.SetStacks(sim, paladin.SpendableHolyPower()+1)
	})
}

func (paladin *Paladin) registerSelflessHealer() {
	if !paladin.Talents.SelflessHealer {
		return
	}

	hpGainActionID := core.ActionID{SpellID: 148502}
	classMask := SpellMaskFlashOfLight | SpellMaskDivineLight | SpellMaskHolyRadiance

	castTimePerStack := []float64{0, -0.35, -0.7, -1}
	castTimeMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  classMask,
		FloatValue: castTimePerStack[0],
	})

	costPerStack := []int32{0, -35, -70, -100}
	costMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: classMask,
		IntValue:  costPerStack[0],
	})

	// TODO: Handle effectiveness modifier in the respective spell files since they're target specific

	var selflessHealerAura *core.Aura
	selflessHealerAura = paladin.RegisterAura(core.Aura{
		Label:     "Selfless Healer" + paladin.Label,
		ActionID:  core.ActionID{SpellID: 114250},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
			costMod.Activate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			castTimeMod.UpdateFloatValue(castTimePerStack[newStacks])
			costMod.UpdateIntValue(costPerStack[newStacks])
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
			costMod.Deactivate()
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Selfless Healer Consume Trigger" + paladin.Label,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: classMask,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			selflessHealerAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Selfless Healer Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 85804},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskJudgment,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			selflessHealerAura.Activate(sim)
			selflessHealerAura.AddStack(sim)

			if paladin.Spec == proto.Spec_SpecHolyPaladin {
				paladin.HolyPower.Gain(1, hpGainActionID, sim)
			}
		},
	})
}

func (paladin *Paladin) registerSacredShield() {
	if !paladin.Talents.SacredShield {
		return
	}

	isHoly := paladin.Spec == proto.Spec_SpecHolyPaladin
	actionID := core.ActionID{SpellID: core.TernaryInt32(isHoly, 148039, 20925)}

	castConfig := core.CastConfig{
		DefaultCast: core.Cast{
			GCD: core.GCDDefault,
		},
		IgnoreHaste: true,
	}

	if !isHoly {
		castConfig.CD = core.Cooldown{
			Timer:    paladin.NewTimer(),
			Duration: time.Second * 6,
		}
	}

	scalingCoef := core.TernaryFloat64(isHoly, 0.30000001192, 0.20999999344)
	baseHealing := paladin.CalcScalingSpellDmg(scalingCoef)
	spCoef := core.TernaryFloat64(isHoly, 1.17, 0.819)
	absorbDuration := time.Second * 6

	var absorbAuras core.DamageAbsorptionAuraArray
	sacredShield := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful,
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolHoly,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(isHoly, 16, 0),
		},

		MaxRange: 40,

		Cast:         castConfig,
		Charges:      core.TernaryInt(isHoly, 3, 0),
		RechargeTime: core.TernaryDuration(isHoly, time.Second*10, 0),

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Sacred Shield",
			},

			TickLength:          absorbDuration,
			NumberOfTicks:       5,
			AffectedByCastSpeed: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				aura := absorbAuras.Get(target)
				aura.Duration = dot.TickPeriod()
				aura.Activate(sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !isHoly {
				for _, unit := range paladin.Env.AllUnits {
					if unit.Type == core.EnemyUnit {
						continue
					}

					aura := unit.GetAuraByID(actionID)
					if aura == nil {
						continue
					}

					aura.Deactivate(sim)
				}
			}

			hot := spell.Hot(target)
			hot.Apply(sim)
			hot.TickOnce(sim)
		},
	})

	absorbAuras = paladin.NewAllyDamageAbsorptionAuraArray(func(unit *core.Unit) *core.DamageAbsorptionAura {
		return unit.NewDamageAbsorptionAura(
			"Sacred Shield (Absorb)",
			core.ActionID{SpellID: 65148},
			absorbDuration,
			func(unit *core.Unit) float64 {
				return baseHealing + sacredShield.SpellPower()*spCoef
			})
	})
}

func (paladin *Paladin) registerHandOfPurity() {
	if !paladin.Talents.HandOfPurity {
		return
	}

	actionID := core.ActionID{SpellID: 114039}

	handAuras := paladin.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		aura := unit.RegisterAura(core.Aura{
			Label:    "Hand of Purity" + unit.Label,
			ActionID: actionID,
			Duration: time.Second * 6,
		}).AttachMultiplicativePseudoStatBuff(&unit.PseudoStats.DamageTakenMultiplier, 0.9)

		unit.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult, isPeriodic bool) {
			if !isPeriodic || result.Damage == 0 || !result.Landed() || !aura.IsActive() {
				return
			}

			incomingDamage := result.Damage
			result.Damage *= incomingDamage * 0.2

			if sim.Log != nil {
				unit.Log(sim, "Hand of Purity absorbed %.1f damage", incomingDamage-result.Damage)
			}
		})

		return aura
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful,
		SpellSchool: core.SpellSchoolHoly,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 7.0,
		},

		MaxRange: 40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			handAuras.Get(target).Activate(sim)
		},
	})
}

func (paladin *Paladin) registerUnbreakableSpirit() {
	if !paladin.Talents.UnbreakableSpirit {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Unbreakable Spirit" + paladin.Label,
		ActionID: core.ActionID{SpellID: 114154},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		ClassMask:  SpellMaskDivineProtection | SpellMaskLayOnHands | SpellMaskDivineShield,
		FloatValue: -0.5,
	})
}

func (paladin *Paladin) registerSanctifiedWrath() {
	if !paladin.Talents.SanctifiedWrath {
		return
	}

	sanctifiedWrathAura := paladin.RegisterAura(core.Aura{
		Label:    "Sanctified Wrath" + paladin.Label,
		ActionID: core.ActionID{SpellID: 114232},
		Duration: time.Second * 30,
	})

	var cdClassMask int64
	if paladin.Spec == proto.Spec_SpecHolyPaladin {
		// Reduces the cooldown of Holy Shock by 50% and increases the critical strike chance of Holy Shock by 20%.
		cdClassMask = SpellMaskHolyShock

		sanctifiedWrathAura.AttachSpellMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  SpellMaskHolyShock,
			FloatValue: 0.2,
		})
	} else if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		// Reduces the cooldown of Judgment by 50%, and causes Judgment to generate one additional Holy Power.
		// Avenging Wrath also increases healing received by 20%.
		cdClassMask = SpellMaskJudgment

		paladin.HolyPower.RegisterOnGain(func(gain, realGain int32, actionID core.ActionID, sim *core.Simulation) {
			if actionID.SpellID == 105427 && paladin.AvengingWrathAura.IsActive() {
				paladin.HolyPower.Gain(1, core.ActionID{SpellID: 53376}, sim)
			}
		})

		sanctifiedWrathAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.HealingTakenMultiplier *= 1.2
		}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.HealingTakenMultiplier /= 1.2
		})
	} else if paladin.Spec == proto.Spec_SpecRetributionPaladin {
		// Reduces the cooldown of Hammer of Wrath by 50%.
		cdClassMask = SpellMaskHammerOfWrath
	}

	sanctifiedWrathAura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		ClassMask:  cdClassMask,
		FloatValue: -0.5,
	})

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(SpellMaskAvengingWrath) {
			paladin.AvengingWrathAura.AttachDependentAura(sanctifiedWrathAura)
		}
	})
}

func (paladin *Paladin) registerDivinePurpose() {
	if !paladin.Talents.DivinePurpose {
		return
	}

	actionID := core.ActionID{SpellID: 90174}
	duration := time.Second * 8
	procChances := []float64{0, 0.08, 0.166, 0.25}
	paladin.DivinePurposeAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Purpose" + paladin.Label,
		ActionID: actionID,
		Duration: duration,
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Divine Purpose Consume Trigger" + paladin.Label,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskSpender,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			var hpSpent int32
			if paladin.DivinePurposeAura.IsActive() {
				paladin.DivinePurposeAura.Deactivate(sim)
				hpSpent = 3
			} else if spell.Matches(SpellMaskDivineStorm | SpellMaskTemplarsVerdict | SpellMaskShieldOfTheRighteous) {
				hpSpent = 3
			} else if spell.Matches(SpellMaskInquisition | SpellMaskWordOfGlory | SpellMaskHarshWords) {
				hpSpent = paladin.DynamicHolyPowerSpent
			} else {
				return
			}

			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + core.SpellBatchWindow,
				OnAction: func(sim *core.Simulation) {
					if sim.Proc(procChances[hpSpent], "Divine Purpose"+paladin.Label) {
						paladin.DivinePurposeAura.Activate(sim)
					}
				},
			})
		},
	})
}

func (paladin *Paladin) registerHolyAvenger() {
	if !paladin.Talents.HolyAvenger {
		return
	}

	var classMask int64
	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		classMask = SpellMaskBuilderProt
	} else if paladin.Spec == proto.Spec_SpecHolyPaladin {
		classMask = SpellMaskBuilderHoly
	} else {
		classMask = SpellMaskBuilderRet
	}

	actionID := core.ActionID{SpellID: 105809}
	holyAvengerAura := paladin.RegisterAura(core.Aura{
		Label:    "Holy Avenger" + paladin.Label,
		ActionID: actionID,
		Duration: time.Second * 18,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  classMask,
		FloatValue: 0.3,
	})

	paladin.HolyPower.RegisterOnGain(func(gain int32, actualGain int32, triggeredActionID core.ActionID, sim *core.Simulation) {
		if !holyAvengerAura.IsActive() {
			return
		}

		if slices.Contains(paladin.HolyAvengerActionIDFilter, &triggeredActionID) {
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + core.SpellBatchWindow,
				OnAction: func(sim *core.Simulation) {
					paladin.HolyPower.Gain(2, actionID, sim)
				},
			})
		}
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		ProcMask: core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 2 * time.Minute,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: holyAvengerAura,
	})
}

func (paladin *Paladin) registerHolyPrism() {
	if !paladin.Talents.HolyPrism {
		return
	}

	numEnemyTargets := min(5, paladin.Env.GetNumTargets())

	damageActionID := core.ActionID{SpellID: 114852}
	healActionID := core.ActionID{SpellID: 114871}

	onUseTimer := paladin.NewTimer()
	onUseCD := time.Second * 20

	targetScalingCoef := 14.13099956512
	targetVariance := 0.20000000298
	targetSpCoef := 1.4279999733

	aoeScalingCoef := 9.52900028229
	aoeVariance := 0.20000000298
	aoeSpCoef := 0.9620000124

	aoeHealSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    damageActionID.WithTag(2),
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolHoly,

		MaxRange:     40,
		MissileSpeed: 100,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := paladin.CalcAndRollDamageRange(sim, aoeScalingCoef, aoeVariance) +
				aoeSpCoef*spell.SpellPower()
			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealOutcome(sim, result)
			})
		},
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    damageActionID.WithTag(1),
		Flags:       core.SpellFlagAPL,
		ProcMask:    core.ProcMaskSpellDamage,
		SpellSchool: core.SpellSchoolHoly,

		MaxRange:     40,
		MissileSpeed: 100,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5.4,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    onUseTimer,
				Duration: onUseCD,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcAndRollDamageRange(sim, targetScalingCoef, targetVariance) +
				targetSpCoef*spell.SpellPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				aoeHealSpell.Cast(sim, &paladin.Unit)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealOutcome(sim, result)
			})
		},
	})

	aoeDamageSpell := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    healActionID.WithTag(2),
		Flags:       core.SpellFlagPassiveSpell,
		ProcMask:    core.ProcMaskSpellDamage,
		SpellSchool: core.SpellSchoolHoly,

		MaxRange:     40,
		MissileSpeed: 100,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numEnemyTargets)

			for i := 0; i < len(results); i++ {
				baseDamage := paladin.CalcAndRollDamageRange(sim, aoeScalingCoef, aoeVariance) +
					aoeSpCoef*spell.SpellPower()
				results[i] = spell.CalcDamage(sim, paladin.Env.Raid.AllPlayerUnits[i], baseDamage, spell.OutcomeMagicCrit)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				for _, result := range results {
					spell.DealOutcome(sim, result)
				}
			})
		},
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    healActionID.WithTag(1),
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful,
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolHoly,

		MaxRange:     40,
		MissileSpeed: 100,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5.4,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    onUseTimer,
				Duration: onUseCD,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := paladin.CalcAndRollDamageRange(sim, targetScalingCoef, targetVariance) +
				targetSpCoef*spell.SpellPower()

			result := spell.CalcHealing(sim, &paladin.Unit, baseHealing, spell.OutcomeHealingCrit)

			aoeDamageSpell.Cast(sim, target)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealOutcome(sim, result)
			})
		},
	})
}

func (paladin *Paladin) registerLightsHammer() {
	if !paladin.Talents.LightsHammer {
		return
	}

	scalingCoef := 3.17899990082
	variance := 0.20000000298
	spCoef := 0.32100000978

	arcingLightDamage := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 114919},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Arcing Light (Damage)" + paladin.Label,
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					baseDamage := paladin.CalcAndRollDamageRange(sim, scalingCoef, variance) +
						spCoef*dot.Spell.SpellPower()
					dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, baseDamage, dot.OutcomeTickMagicHitAndCrit)
				}
			},
		},
	})

	arcingLightHealing := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 119952},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Hot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Arcing Light (Healing)" + paladin.Label,
			},
			NumberOfTicks: 8,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Raid.AllUnits {
					baseHealing := paladin.CalcAndRollDamageRange(sim, scalingCoef, variance) +
						spCoef*dot.Spell.SpellPower()
					dot.Spell.CalcAndDealPeriodicHealing(sim, aoeTarget, baseHealing, dot.OutcomeTickHealingCrit)
				}
			},
		},
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 114158},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		MaxRange:     30,
		MissileSpeed: 20,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aoeDamageDot := arcingLightDamage.AOEDot()
			aoeHealingDot := arcingLightHealing.AOEDot()

			if sim.Proc(0.5, "Arcing Light 9 ticks"+paladin.Label) {
				aoeDamageDot.BaseTickCount = 9
				aoeHealingDot.BaseTickCount = 9
			} else {
				aoeDamageDot.BaseTickCount = 8
				aoeHealingDot.BaseTickCount = 8
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				aoeDamageDot.Apply(sim)
				aoeHealingDot.Apply(sim)
			})
		},
	})
}

func (paladin *Paladin) registerExecutionSentence() {
	if !paladin.Talents.ExecutionSentence {
		return
	}

	baseTickDamage := paladin.CalcScalingSpellDmg(0.42599999905)
	spCoef := 5936 / 1000.0
	totalBonusCoef := 0.0

	tickMultipliers := make([]float64, 11)
	tickMultipliers[0] = 1.0
	for i := 1; i < 10; i++ {
		tickMultipliers[i] = tickMultipliers[i-1] * 1.1
		totalBonusCoef += tickMultipliers[i]
	}
	tickMultipliers[10] = tickMultipliers[9] * 5
	totalBonusCoef += tickMultipliers[10]

	tickSpCoef := spCoef * (1 / totalBonusCoef)

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 114916},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		MaxRange: 40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Execution Sentence" + paladin.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.SpellPower())
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				snapshotSpellPower := dot.SnapshotBaseDamage

				tickMultiplier := tickMultipliers[dot.TickCount()+1]
				dot.SnapshotBaseDamage = tickMultiplier*baseTickDamage +
					tickMultiplier*tickSpCoef*snapshotSpellPower

				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

				dot.SnapshotBaseDamage = snapshotSpellPower
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
