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

// Increases your movement speed by 70% for 8 sec.
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

// A successful Judgment increases your movement speed by 45% for 3 sec.
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

// You gain 15% movement speed at all times, plus an additional 5% movement speed for each current charge of Holy Power up to 3.
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

	paladin.HolyPower.RegisterOnGain(func(sim *core.Simulation, gain, realGain int32, actionID core.ActionID) {
		pursuitOfJusticeAura.Activate(sim)
		pursuitOfJusticeAura.SetStacks(sim, paladin.SpendableHolyPower()+1)
	})
	paladin.HolyPower.RegisterOnSpend(func(sim *core.Simulation, amount int32, actionID core.ActionID) {
		pursuitOfJusticeAura.Activate(sim)
		pursuitOfJusticeAura.SetStacks(sim, paladin.SpendableHolyPower()+1)
	})
}

/*
Your successful Judgments

-- Holy Insight --
generate a charge of Holy Power and
-- /Holy Insight --

reduce the cast time and mana cost of your next Flash of Light

-- Denounce --
, Divine Light, or Holy Radiance
-- /Denounce --

by 35% per stack and improves its effectiveness by 20% per stack when used to heal others.
Stacks up to 3 times.
(500ms cooldown)
*/
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

	paladin.SelflessHealerAura = paladin.RegisterAura(core.Aura{
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
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: classMask,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.SelflessHealerAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Selfless Healer Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 85804},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskJudgment,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.SelflessHealerAura.Activate(sim)
			paladin.SelflessHealerAura.AddStack(sim)

			if paladin.Spec == proto.Spec_SpecHolyPaladin {
				paladin.HolyPower.Gain(sim, 1, hpGainActionID)
			}
		},
	})
}

/*
Ret & Prot:

Protects the target with a shield of Holy Light for 30 sec.
The shield absorbs up to (240 + 0.819 * <SP>) damage every 6 sec.
Can be active only on one target at a time.

Holy:

Protects the target with a shield of Holy Light for 30 sec.
The shield absorbs up to (343 + 1.17 * <SP>) damage every 6 sec.
Max 3 charges.
*/
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
	}

	if !isHoly {
		castConfig.CD = core.Cooldown{
			Timer:    paladin.NewTimer(),
			Duration: time.Second * 6,
		}
	}

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
			if target.IsOpponent(&paladin.Unit) {
				target = &paladin.Unit
			}

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

	baseHealing := paladin.CalcScalingSpellDmg(core.TernaryFloat64(isHoly, 0.30000001192, 0.20999999344))
	spCoef := core.TernaryFloat64(isHoly, 1.17, 0.819)
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

/*
Places a Hand on the friendly target, reducing damage taken by 10% and damage from harmful periodic effects by an additional 80% (less for some creature attacks) for 6 sec.
Players may only have one Hand on them per Paladin at any one time.
*/
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

// Reduces the cooldown of your Divine Shield, Divine Protection and Lay on Hands by 50%.
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
		FloatValue: 0.5,
	})
}

// Abilities that generate Holy Power will deal 30% additional damage and healing, and generate 3 charges of Holy Power for the next 18 sec.
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

	paladin.HolyPower.RegisterOnGain(func(sim *core.Simulation, gain int32, actualGain int32, triggeredActionID core.ActionID) {
		if !holyAvengerAura.IsActive() {
			return
		}

		if slices.Contains(paladin.HolyAvengerActionIDFilter, triggeredActionID) {
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + core.SpellBatchWindow,
				OnAction: func(sim *core.Simulation) {
					paladin.HolyPower.Gain(sim, 2, actionID)
				},
			})
		}
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		Flags:       core.SpellFlagAPL | core.SpellFlagHelpful,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,

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

// Avenging Wrath lasts 50% longer and grants more frequent access to one of your abilities while it lasts.
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
		hpGainActionID := core.ActionID{SpellID: 53376}

		paladin.HolyPower.RegisterOnGain(func(sim *core.Simulation, gain, realGain int32, actionID core.ActionID) {
			if actionID.SameAction(paladin.JudgmentsOfTheWiseActionID) && paladin.AvengingWrathAura.IsActive() {
				paladin.HolyPower.Gain(sim, 1, hpGainActionID)
			}
		})

		sanctifiedWrathAura.AttachMultiplicativePseudoStatBuff(&paladin.PseudoStats.HealingTakenMultiplier, 1.2)
	} else if paladin.Spec == proto.Spec_SpecRetributionPaladin {
		// Reduces the cooldown of Hammer of Wrath by 50%.
		cdClassMask = SpellMaskHammerOfWrath
	}

	sanctifiedWrathAura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		ClassMask:  cdClassMask,
		FloatValue: 0.5,
	})

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(SpellMaskAvengingWrath) {
			paladin.AvengingWrathAura.AttachDependentAura(sanctifiedWrathAura)
		}
	})
}

type AuraDeactivationCheck func(aura *core.Aura, spell *core.Spell) bool

func (paladin *Paladin) divinePurposeFactory(label string, spellID int32, duration time.Duration, auraDeactivationCheck AuraDeactivationCheck) *core.Aura {
	procChances := []float64{0, 0.25 * (1 / 3), 0.25 * (2 / 3), 0.25}
	aura := paladin.RegisterAura(core.Aura{
		Label:    label + paladin.Label,
		ActionID: core.ActionID{SpellID: spellID},
		Duration: duration,
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           label + " Consume Trigger" + paladin.Label,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskSpender,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			var hpSpent int32
			if aura.IsActive() && (auraDeactivationCheck == nil || auraDeactivationCheck(aura, spell)) {
				aura.Deactivate(sim)
				hpSpent = 3
			} else if spell.Matches(SpellMaskDivineStorm | SpellMaskTemplarsVerdict | SpellMaskShieldOfTheRighteous) {
				hpSpent = 3
			} else if spell.Matches(SpellMaskInquisition | SpellMaskWordOfGlory | SpellMaskHarshWords) {
				hpSpent = paladin.DynamicHolyPowerSpent
			} else {
				return
			}

			if sim.Proc(procChances[hpSpent], label+paladin.Label) {
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + core.SpellBatchWindow,
					OnAction: func(sim *core.Simulation) {
						aura.Activate(sim)
					},
				})
			}
		},
	})

	return aura
}

/*
Abilities that cost Holy Power have a 25% chance to cause the Divine Purpose effect.

Divine Purpose
Your next Holy Power ability will consume no Holy Power and will cast as if 3 Holy Power were consumed.
Lasts 8 sec.
*/
func (paladin *Paladin) registerDivinePurpose() {
	if !paladin.Talents.DivinePurpose {
		return
	}

	paladin.DivinePurposeAura = paladin.divinePurposeFactory("Divine Purpose", 90174, time.Second*8, func(aura *core.Aura, spell *core.Spell) bool {
		return true
	})
}

func (paladin *Paladin) holyPrismFactory(spellID int32, targets []*core.Unit, timer *core.Timer, isHealing bool) {
	numTargets := len(targets)
	actionID := core.ActionID{SpellID: spellID}

	aoeConfig := core.SpellConfig{
		ActionID:    actionID.WithTag(2),
		Flags:       core.SpellFlagPassiveSpell,
		SpellSchool: core.SpellSchoolHoly,

		MaxRange:     40,
		MissileSpeed: 100,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 0.9620000124,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)

			for idx, aoeTarget := range targets {
				base := paladin.CalcAndRollDamageRange(sim, 9.52900028229, 0.20000000298)
				// isHealing = true means the direct spell is a heal and the aoe spell is damage
				if !isHealing {
					results[idx] = spell.CalcHealing(sim, aoeTarget, base, spell.OutcomeHealingCrit)
				} else {
					results[idx] = spell.CalcDamage(sim, aoeTarget, base, spell.OutcomeMagicCrit)
				}
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				for _, result := range results {
					// isHealing = true means the direct spell is a heal and the aoe spell is damage
					if !isHealing {
						spell.DealHealing(sim, result)
					} else {
						spell.DealDamage(sim, result)
					}
				}
			})
		},
	}

	// isHealing = true means the direct spell is a heal and the aoe spell is damage
	if !isHealing {
		aoeConfig.Flags |= core.SpellFlagHelpful
		aoeConfig.ProcMask = core.ProcMaskSpellHealing
	} else {
		aoeConfig.ProcMask = core.ProcMaskSpellDamage
	}

	aoeSpell := paladin.RegisterSpell(aoeConfig)

	directSpellConfig := core.SpellConfig{
		ActionID:    actionID.WithTag(1),
		Flags:       core.SpellFlagAPL,
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
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 20,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 1.4279999733,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if isHealing && target.IsOpponent(&paladin.Unit) {
				target = &paladin.Unit
			}

			base := paladin.CalcAndRollDamageRange(sim, 14.13099956512, 0.20000000298)

			var result *core.SpellResult
			if isHealing {
				result = spell.CalcHealing(sim, target, base, spell.OutcomeHealingCrit)
			} else {
				result = spell.CalcDamage(sim, target, base, spell.OutcomeMagicHitAndCrit)
			}

			if isHealing {
				aoeSpell.Cast(sim, paladin.CurrentTarget)
			} else if result.Landed() {
				aoeSpell.Cast(sim, &paladin.Unit)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if isHealing {
					spell.DealHealing(sim, result)
				} else {
					spell.DealDamage(sim, result)
				}
			})
		},
	}

	if isHealing {
		directSpellConfig.Flags |= core.SpellFlagHelpful
		directSpellConfig.ProcMask = core.ProcMaskSpellHealing
	} else {
		directSpellConfig.ProcMask = core.ProcMaskSpellDamage
	}

	paladin.RegisterSpell(directSpellConfig)
}

/*
Sends a beam of light toward a target, turning them into a prism for Holy energy.
If an enemy is the prism, they take (<14522-17751> + 1.428 * <SP>) Holy damage and radiate (<9793-11970> + 0.962 * <SP>) healing to 5 nearby allies within 15 yards.
If an ally is the prism, they are healed for (<14522-17751> + 1.428 * <SP>) and radiate (<9793-11970> + 0.962 * <SP>) Holy damage to 5 nearby enemies within 15 yards.
*/
func (paladin *Paladin) registerHolyPrism() {
	if !paladin.Talents.HolyPrism {
		return
	}

	onUseTimer := paladin.NewTimer()

	friendlyTargets := paladin.Env.Raid.GetFirstNPlayersOrPets(5)
	paladin.holyPrismFactory(114852, friendlyTargets, onUseTimer, false)

	enemyTargets := core.MapSlice(paladin.Env.Encounter.ActiveTargets[:min(5, int32(len(paladin.Env.Encounter.ActiveTargets)))], func(target *core.Target) *core.Unit {
		return &target.Unit
	})
	paladin.holyPrismFactory(114871, enemyTargets, onUseTimer, true)
}

/*
Hurl a Light-infused hammer into the ground, where it will blast a 10 yard area with Arcing Light for 14 sec.

Arcing Light
Deals (<3267-3994> + 0.321 * <SP>) Holy damage to enemies and reduces their movement speed by 50% for 2 sec.
Heals allies for (<3267-3994> + 0.321 * <SP>) every 2 sec.
*/
func (paladin *Paladin) registerLightsHammer() {
	if !paladin.Talents.LightsHammer {
		return
	}

	enemyTargets := paladin.Env.Encounter.ActiveTargets
	friendlyTargets := paladin.Env.Raid.GetFirstNPlayersOrPets(6)

	tickCount := int32(8)

	arcingLightDamage := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 114919},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagAoE,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Arcing Light (Damage)" + paladin.Label,
			},
			NumberOfTicks: tickCount,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				results := make([]*core.SpellResult, len(enemyTargets))

				for idx, currentTarget := range enemyTargets {
					baseDamage := paladin.CalcAndRollDamageRange(sim, 3.17899990082, 0.20000000298) +
						0.32100000978*dot.Spell.SpellPower()
					results[idx] = dot.Spell.CalcPeriodicDamage(sim, &currentTarget.Unit, baseDamage, dot.OutcomeTickMagicHitAndCrit)
				}

				for _, result := range results {
					dot.Spell.DealPeriodicDamage(sim, result)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.BaseTickCount = tickCount
			dot.Apply(sim)
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
			NumberOfTicks: tickCount,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				results := make([]*core.SpellResult, len(friendlyTargets))

				for idx, aoeTarget := range friendlyTargets {
					baseHealing := paladin.CalcAndRollDamageRange(sim, 3.17899990082, 0.20000000298) +
						0.32100000978*dot.Spell.SpellPower()
					results[idx] = dot.Spell.CalcHealing(sim, aoeTarget, baseHealing, dot.OutcomeTickHealingCrit)
				}

				for _, result := range results {
					dot.Spell.DealPeriodicHealing(sim, result)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hot := spell.AOEHot()
			hot.BaseTickCount = tickCount
			hot.Apply(sim)
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
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.Proc(0.5, "Arcing Light 9 ticks"+paladin.Label) {
				tickCount = 9
			} else {
				tickCount = 8
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				arcingLightDamage.Cast(sim, target)
				arcingLightHealing.Cast(sim, target)
			})
		},
	})
}

func (paladin *Paladin) executionSentenceFactory(spellID int32, label string, cd *core.Timer, tickMultipliers []float64, bonusCoef float64, tickSpCoef float64, isHealing bool) {
	config := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       core.SpellFlagAPL,

		MaxRange: 40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cd,
				Duration: time.Minute,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if isHealing {
				if target.IsOpponent(&paladin.Unit) {
					target = &paladin.Unit
				}

				spell.Hot(target).Apply(sim)
			} else {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			}
		},
	}

	dotConfig := core.DotConfig{
		Aura: core.Aura{
			Label: label + paladin.Label,
		},
		NumberOfTicks: 10,
		TickLength:    time.Second,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.Snapshot(target, dot.Spell.SpellPower())
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			snapshotSpellPower := dot.SnapshotBaseDamage

			tickMultiplier := tickMultipliers[dot.TickCount()]
			dot.SnapshotBaseDamage = tickMultiplier*paladin.CalcScalingSpellDmg(0.42599999905) +
				tickMultiplier*tickSpCoef*snapshotSpellPower

			if isHealing {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeSnapshotCrit)
			} else {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			}

			dot.SnapshotBaseDamage = snapshotSpellPower
		},
	}

	if isHealing {
		config.Hot = dotConfig
		config.Flags |= core.SpellFlagHelpful
		config.ProcMask = core.ProcMaskSpellHealing
	} else {
		config.Dot = dotConfig
		config.ProcMask = core.ProcMaskSpellDamage
	}

	paladin.RegisterSpell(config)
}

/*
Execution Sentence:

A hammer slowly falls from the sky, causing (<SP> * 5936 / 1000 + 26.72716306 * 486) Holy damage over 10 sec.
This damage is dealt slowly at first and increases over time, culminating in a final burst of damage.
Dispelling the effect triggers the final burst.

Stay of Execution:

If used on friendly targets, the falling hammer heals the target for (<SP> * 5936 / 1000 + 26.72716306 * 486) healing over 10 sec.
This healing is dealt slowly at first and increases over time, culminating in a final burst of healing.
Dispelling the effect triggers the final burst.
*/
func (paladin *Paladin) registerExecutionSentence() {
	if !paladin.Talents.ExecutionSentence {
		return
	}

	totalBonusCoef := 0.0

	tickMultipliers := make([]float64, 11)
	tickMultipliers[0] = 1.0
	for i := 1; i < 10; i++ {
		tickMultipliers[i] = tickMultipliers[i-1] * 1.1
		totalBonusCoef += tickMultipliers[i]
	}
	tickMultipliers[10] = tickMultipliers[9] * 5
	totalBonusCoef += tickMultipliers[10]

	tickSpCoef := 5936 / 1000.0 * (1 / totalBonusCoef)

	cd := paladin.NewTimer()

	paladin.executionSentenceFactory(114916, "Execution Sentence", cd, tickMultipliers, totalBonusCoef, tickSpCoef, false)
	paladin.executionSentenceFactory(146586, "Stay of Execution", cd, tickMultipliers, totalBonusCoef, tickSpCoef, true)
}
