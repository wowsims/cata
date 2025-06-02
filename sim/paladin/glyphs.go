package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) registerGlyphs() {
	// Major glyphs
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfAvengingWrath) {
		paladin.registerGlyphOfAvengingWrath()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfBurdenOfGuilt) {
		paladin.registerGlyphOfBurdenOfGuilt()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDazingShield) {
		paladin.registerGlyphOfDazingShield()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDenounce) {
		paladin.registerGlyphOfDenounce()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDevotionAura) {
		paladin.registerGlyphOfDevotionAura()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivinePlea) {
		paladin.registerGlyphOfDivinePlea()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivineProtection) {
		// Handled in divine_protection.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivineStorm) {
		paladin.registerGlyphOfDivineStorm()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivinity) {
		paladin.registerGlyphOfDivinity()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDoubleJeopardy) {
		paladin.registerGlyphOfDoubleJeopardy()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFinalWrath) {
		// Handled in protection/holy_wrath.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFlashOfLight) {
		paladin.registerGlyphOfFlashOfLight()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFocusedShield) {
		paladin.registerGlyphOfFocusedShield()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfTheRighteous) {
		// Handled in hammer_of_the_righteous.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHarshWords) {
		paladin.registerGlyphOfHarshWords()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHolyShock) {
		paladin.registerGlyphOfHolyShock()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfIllumination) {
		paladin.registerGlyphOfIllumination()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfImmediateTruth) {
		paladin.registerGlyphOfImmediateTruth()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfLightOfDawn) {
		paladin.registerGlyphOfLightOfDawn()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfMassExorcism) {
		paladin.registerGlyphOfMassExorcism()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfProtectorOfTheInnocent) {
		paladin.registerGlyphOfProtectorOfTheInnocent()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfTemplarsVerdict) {
		paladin.registerGlyphOfTemplarsVerdict()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfTheAlabasterShield) {
		paladin.registerGlyphOfTheAlabasterShield()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfTheBattleHealer) {
		paladin.registerGlyphOfTheBattleHealer()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfWordOfGlory) {
		paladin.registerGlyphOfWordOfGlory()
	}

	// Minor glyphs
	if paladin.HasMinorGlyph(proto.PaladinMinorGlyph_GlyphOfFocusedWrath) {
		// Handled in protection/holy_wrath.go
	}
}

// While Avenging Wrath is active, you are healed for 1% of your maximum health every 2 sec.
func (paladin *Paladin) registerGlyphOfAvengingWrath() {
	actionID := core.ActionID{SpellID: 115547}
	healthMetrics := paladin.NewHealthMetrics(actionID)

	var healPA *core.PendingAction
	glyphAura := paladin.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Glyph of Avenging Wrath" + paladin.Label,
		Duration: core.DurationFromSeconds(core.TernaryFloat64(paladin.Talents.SanctifiedWrath, 30, 20)),

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			healPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 2,
				NumTicks: 10,
				OnAction: func(sim *core.Simulation) {
					paladin.GainHealth(sim, paladin.MaxHealth()*0.01, healthMetrics)
				},
			})
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt: sim.CurrentTime + core.SpellBatchWindow,
				OnAction: func(sim *core.Simulation) {
					if healPA != nil {
						healPA.Cancel(sim)
					}
				},
			})
		},
	})

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(SpellMaskAvengingWrath) {
			paladin.AvengingWrathAura.AttachDependentAura(glyphAura)
		}
	})
}

// Your Judgment hits fill your target with doubt and remorse, reducing movement speed by 50% for 2 sec.
func (paladin *Paladin) registerGlyphOfBurdenOfGuilt() {
	burdenOfGuiltAuras := paladin.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return unit.RegisterAura(core.Aura{
			Label:    "Burden of Guilt" + unit.Label,
			ActionID: core.ActionID{SpellID: 110300},
			Duration: time.Second * 2,
		}).AttachMultiplicativePseudoStatBuff(&unit.PseudoStats.MovementSpeedMultiplier, 0.5)
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Burden of Guilt" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 54931},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskJudgment,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			burdenOfGuiltAuras.Get(result.Target).Activate(sim)
		},
	})
}

// Your Avenger's Shield now also dazes targets for 10 sec.
func (paladin *Paladin) registerGlyphOfDazingShield() {
	if paladin.Spec != proto.Spec_SpecProtectionPaladin {
		return
	}

	dazedAuras := paladin.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return unit.RegisterAura(core.Aura{
			Label:    "Dazed - Avenger's Shield" + unit.Label,
			ActionID: core.ActionID{SpellID: 63529},
			Duration: time.Second * 10,
		}).AttachMultiplicativePseudoStatBuff(&unit.PseudoStats.MovementSpeedMultiplier, 0.5)
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Dazing Shield" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 56414},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskAvengersShield,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dazedAuras.Get(result.Target).Activate(sim)
		},
	})
}

// Your Holy Shocks reduce the cast time of your next Denounce by 0.5 sec. This effect stacks up to 3 times.
func (paladin *Paladin) registerGlyphOfDenounce() {
	if paladin.Spec != proto.Spec_SpecHolyPaladin {
		return
	}

	cdMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: SpellMaskDenounce,
		TimeValue: time.Millisecond * -500,
	})

	denounceAura := paladin.RegisterAura(core.Aura{
		Label:     "Glyph of Denounce" + paladin.Label,
		ActionID:  core.ActionID{SpellID: 115654},
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			cdMod.Activate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			cdMod.UpdateTimeValue(time.Millisecond * time.Duration(-500*newStacks))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cdMod.Deactivate()
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Denounce Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 56420},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskHolyShock,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			denounceAura.Activate(sim)
			denounceAura.AddStack(sim)
		},
	})
}

// Devotion Aura no longer affects party or raid members, but the cooldown is reduced by 60 sec.
func (paladin *Paladin) registerGlyphOfDevotionAura() {
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Devotion Aura" + paladin.Label,
		ActionID: core.ActionID{SpellID: 146955},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: SpellMaskDevotionAura,
		TimeValue: time.Second * -60,
	})
}

// Divine Plea returns 50% less mana but has a 50% shorter cooldown.
func (paladin *Paladin) registerGlyphOfDivinePlea() {
	if paladin.Spec != proto.Spec_SpecHolyPaladin {
		return
	}

	// TODO: Handle the mana return part in holy/divine_plea.go
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Divine Plea" + paladin.Label,
		ActionID: core.ActionID{SpellID: 63223},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_Cooldown_Multiplier,
		ClassMask:  SpellMaskDivinePlea,
		FloatValue: 0.5,
	})
}

// Your Divine Storm also heals you for 5% of your maximum health.
func (paladin *Paladin) registerGlyphOfDivineStorm() {
	if paladin.Spec != proto.Spec_SpecRetributionPaladin {
		return
	}

	healthMetrics := paladin.NewHealthMetrics(core.ActionID{SpellID: 115515})
	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Divine Storm" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 63220},
		Callback:       core.CallbackOnCastComplete, // DS doesn't have to hit anything, it still heals
		ClassSpellMask: SpellMaskDivineStorm,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.GainHealth(sim, paladin.MaxHealth()*0.05, healthMetrics)
		},
	})
}

// Increases the cooldown of your Lay on Hands by 2 min but causes it to give you 10% of your maximum mana.
func (paladin *Paladin) registerGlyphOfDivinity() {
	manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: 54986})

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Divinity" + paladin.Label,
		ActionID: core.ActionID{SpellID: 54939},
	})).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskLayOnHands,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AddMana(sim, paladin.MaxMana()*0.10, manaMetrics)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: SpellMaskLayOnHands,
		TimeValue: time.Minute * 2,
	})
}

// Judging a target increases the damage of your next Judgment by 20%, but only if used on a second target.
func (paladin *Paladin) registerGlyphOfDoubleJeopardy() {
	spellMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskJudgment,
		FloatValue: 0.2,
	})

	var triggeredTarget *core.Unit
	doubleJeopardyAura := paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Double Jeopardy" + paladin.Label,
		ActionID: core.ActionID{SpellID: 121027},
		Duration: time.Second * 10,

		OnApplyEffects: func(aura *core.Aura, sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.Matches(SpellMaskJudgment) {
				aura.Deactivate(sim)

				if target != triggeredTarget {
					spellMod.Activate()
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(SpellMaskJudgment) {
				spellMod.Deactivate()
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Double Jeopardy Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 54922},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskJudgment,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(SpellMaskJudgment) && !doubleJeopardyAura.IsActive() {
				triggeredTarget = result.Target
				doubleJeopardyAura.Activate(sim)
			}
		},
	})
}

// When you Flash of Light a target, it increases your next heal done to that target within 7 sec by 10%.
func (paladin *Paladin) registerGlyphOfFlashOfLight() {
	glyphAuras := paladin.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		return unit.RegisterAura(core.Aura{
			Label:    "Glyph of Flash of Light" + unit.Label,
			ActionID: core.ActionID{SpellID: 54957},
			Duration: time.Second * 7,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				paladin.AttackTables[unit.UnitIndex].HealingDealtMultiplier *= 1.1
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				paladin.AttackTables[unit.UnitIndex].HealingDealtMultiplier /= 1.1
			},
			OnHealTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Unit == &paladin.Unit {
					aura.Deactivate(sim)
				}
			},
		})
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Flash of Light Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 57955},
		Callback:       core.CallbackOnHealDealt,
		ClassSpellMask: SpellMaskFlashOfLight,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			glyphAuras.Get(result.Target).Activate(sim)
		},
	})
}

// Your Avenger's Shield hits 2 fewer targets, but for 30% more damage.
func (paladin *Paladin) registerGlyphOfFocusedShield() {
	if paladin.Spec != proto.Spec_SpecProtectionPaladin {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Focused Shield" + paladin.Label,
		ActionID: core.ActionID{SpellID: 54930},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskAvengersShield,
		FloatValue: 0.3,
	})
}

// Your Word of Glory can now also be used on enemy targets, causing Holy damage approximately equal to the amount it would have healed.
// Does not work with Eternal Flame.
func (paladin *Paladin) registerGlyphOfHarshWords() {
	if paladin.Talents.EternalFlame {
		return
	}

	isProt := paladin.Spec == proto.Spec_SpecProtectionPaladin
	actionID := core.ActionID{SpellID: 130552}

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellDamage,
		SpellSchool:    core.SpellSchoolHoly,
		ClassSpellMask: SpellMaskHarshWords,
		MetricSplits:   4,

		MaxRange: 40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.TernaryDuration(isProt, 0, core.GCDDefault),
				NonEmpty: isProt,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				paladin.DynamicHolyPowerSpent = paladin.SpendableHolyPower()
				spell.SetMetricsSplit(paladin.DynamicHolyPowerSpent)
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.HolyPower.CanSpend(1)
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 0.37700000405,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damageMultiplier := spell.DamageMultiplier
			spell.DamageMultiplier *= float64(paladin.DynamicHolyPowerSpent)

			baseDamage := paladin.CalcAndRollDamageRange(sim, 3.73000001907, 0.1080000028)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.DamageMultiplier = damageMultiplier

			if result.Landed() {
				paladin.HolyPower.SpendUpTo(sim, paladin.DynamicHolyPowerSpent, actionID)
			}

			spell.DealDamage(sim, result)
		},
	})
}

// Decreases the healing of Holy Shock by 50% but increases its damage by 50%.
func (paladin *Paladin) registerGlyphOfHolyShock() {
	if paladin.Spec != proto.Spec_SpecHolyPaladin {
		return
	}

	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Holy Shock" + paladin.Label,
		ActionID: core.ActionID{SpellID: 63224},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskHolyShockDamage,
		FloatValue: 0.5,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskHolyShockHeal,
		FloatValue: -0.5,
	})
}

// Your Holy Shock criticals grant 1% mana return, but Holy Insight returns 10% less mana.
// (800ms cooldown)
func (paladin *Paladin) registerGlyphOfIllumination() {
	if paladin.Spec != proto.Spec_SpecHolyPaladin {
		return
	}

	manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: 115314})

	// TODO: Handle the Holy Insight part in holy/holy.go
	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Illumination" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 54937},
		Callback:       core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskHolyShock,
		ICD:            time.Millisecond * 800,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AddMana(sim, paladin.MaxMana()*0.01, manaMetrics)
		},
	})
}

// Increases the instant damage done by Seal of Truth by 40%, but decreases the damage done by Censure by 50%.
func (paladin *Paladin) registerGlyphOfImmediateTruth() {
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Immediate Truth" + paladin.Label,
		ActionID: core.ActionID{SpellID: 115546},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskSealOfTruth,
		FloatValue: 0.4,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskCensure,
		FloatValue: -0.5,
	})
}

// Light of Dawn affects 2 fewer targets, but heals each target for 25% more.
func (paladin *Paladin) registerGlyphOfLightOfDawn() {
	if paladin.Spec != proto.Spec_SpecHolyPaladin {
		return
	}

	// TODO: Handle the target count part in holy/light_of_dawn.go
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Light of Dawn" + paladin.Label,
		ActionID: core.ActionID{SpellID: 54940},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskLightOfDawn,
		FloatValue: 0.25,
	})
}

// Reduces the range of Exorcism to melee range, but causes 25% damage to all enemies within 8 yards of the primary target.
func (paladin *Paladin) registerGlyphOfMassExorcism() {
	if paladin.Spec != proto.Spec_SpecRetributionPaladin {
		return
	}

	numTargets := paladin.Env.GetNumTargets() - 1

	massExorcism := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 879}.WithTag(2), // Actual 122032
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete | core.SpellFlagAoE,
		ClassSpellMask: SpellMaskExorcism,

		MaxRange: core.MaxMeleeRange,

		DamageMultiplier: 0.25,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)

			currentTarget := sim.Environment.NextTargetUnit(target)
			for idx := range numTargets {
				baseDamage := paladin.CalcAndRollDamageRange(sim, 6.09499979019, 0.1099999994) +
					0.67699998617*spell.MeleeAttackPower()

				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				currentTarget = sim.Environment.NextTargetUnit(currentTarget)
			}

			for idx := range numTargets {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Mass Exorcism" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 122028},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskExorcism,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ActionID.Tag == 2 || numTargets == 0 {
				return
			}

			massExorcism.Cast(sim, result.Target)
		},
	}).ExposeToAPL(122028)
}

// When you use Word of Glory to heal another target, it also heals you for 20% of the amount.
func (paladin *Paladin) registerGlyphOfProtectorOfTheInnocent() {
	var lastHeal float64
	protectorOfTheInnocent := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 115536},
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful | core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods,
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolHoly,

		DamageMultiplier: 0.2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, lastHeal, spell.OutcomeHealing)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Protector of the Innocent" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 93466},
		Callback:       core.CallbackOnHealDealt,
		ClassSpellMask: SpellMaskWordOfGlory,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target == &paladin.Unit {
				return
			}

			lastHeal = result.Damage
			protectorOfTheInnocent.Cast(sim, &paladin.Unit)
		},
	})
}

// You take 10% less damage for 6 sec after dealing damage with Templar's Verdict or Exorcism.
func (paladin *Paladin) registerGlyphOfTemplarsVerdict() {
	glyphOfTemplarVerdictAura := paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Templar's Verdict" + paladin.Label,
		ActionID: core.ActionID{SpellID: 115668},
		Duration: time.Second * 6,
	}).AttachMultiplicativePseudoStatBuff(&paladin.PseudoStats.DamageTakenMultiplier, 0.9)

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Templar's Verdict Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 54926},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskExorcism | SpellMaskTemplarsVerdict,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			glyphOfTemplarVerdictAura.Activate(sim)
		},
	})
}

// Your successful blocks increase the damage of your next Shield of the Righteous by 10%. Stacks up to 3 times.
func (paladin *Paladin) registerGlyphOfTheAlabasterShield() {
	spellMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskShieldOfTheRighteous,
		FloatValue: 0.1,
	})

	alabasterShieldAura := paladin.RegisterAura(core.Aura{
		Label:     "Alabaster Shield" + paladin.Label,
		ActionID:  core.ActionID{SpellID: 121467},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			spellMod.Activate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			spellMod.UpdateFloatValue(0.1 * float64(newStacks))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			spellMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(SpellMaskShieldOfTheRighteous) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:     "Glyph of the Alabaster Shield" + paladin.Label,
		ActionID: core.ActionID{SpellID: 63222},
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeBlock,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			alabasterShieldAura.Activate(sim)
			alabasterShieldAura.AddStack(sim)
		},
	})
}

// Melee attacks from Seal of Insight heal the most wounded member of your raid or party for 30% of the normal heal instead of you.
func (paladin *Paladin) registerGlyphOfTheBattleHealer() {
	// Targeting handled in seal_of_insight.go
	core.MakePermanent(paladin.RegisterAura(core.Aura{
		Label:    "Glyph of the Battle Healer" + paladin.Label,
		ActionID: core.ActionID{SpellID: 119477},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskSealOfInsight,
		FloatValue: -0.7,
	})
}

// Increases your damage by 3% per Holy Power spent after you cast Word of Glory or Eternal Flame on a friendly target. Lasts 6 sec.
func (paladin *Paladin) registerGlyphOfWordOfGlory() {
	glyphAura := paladin.RegisterAura(core.Aura{
		Label:    "Glyph of Word of Glory" + paladin.Label,
		ActionID: core.ActionID{SpellID: 115522},
		Duration: time.Second * 6,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			paladin.PseudoStats.DamageDealtMultiplier *= (1 + 0.03*float64(newStacks)) / (1 + 0.03*float64(oldStacks))
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Glyph of Word of Glory Trigger" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 54936},
		Callback:       core.CallbackOnHealDealt,
		ClassSpellMask: SpellMaskWordOfGlory,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Target == &paladin.Unit {
				return
			}

			glyphAura.Activate(sim)
			glyphAura.SetStacks(sim, paladin.DynamicHolyPowerSpent)
		},
	})
}
