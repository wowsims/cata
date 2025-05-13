package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) applyGlyphs() {
	// Major glyphs
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfAvengingWrath) {
		paladin.registerGlyphOfAvengingWrath()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDevotionAura) {
		paladin.registerGlyphOfDevotionAura()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivineProtection) {
		// Handled in divine_protection.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDivineStorm) {
		paladin.registerGlyphOfDivineStorm()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfDoubleJeopardy) {
		paladin.registerGlyphOfDoubleJeopardy()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFinalWrath) {
		// Handled in protection/holy_wrath.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfFocusedShield) {
		// Handled in protection/avengers_shield.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfTheRighteous) {
		// Handled in hammer_of_the_righteous.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHarshWords) {
		paladin.registerGlyphOfHarshWords()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfImmediateTruth) {
		paladin.registerGlyphOfImmediateTruth()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfMassExorcism) {
		// Handled in retribution/exorcism.go
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfProtectorOfTheInnocent) {
		paladin.registerGlyphOfProtectorOfTheInnocent()
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfTheAlabasterShield) {
		paladin.registerGlyphOfTheAlabasterShield()
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

// Your Divine Storm also heals you for 5% of your maximum health.
func (paladin *Paladin) registerGlyphOfDivineStorm() {
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

// Your Word of Glory can now also be used on enemy targets, causing Holy damage approximately equal to the amount it would have healed.
// Does not work with Eternal Flame.
func (paladin *Paladin) registerGlyphOfHarshWords() {
	if paladin.Talents.EternalFlame {
		return
	}

	isProt := paladin.Spec == proto.Spec_SpecProtectionPaladin
	actionID := core.ActionID{SpellID: 130552}
	scalingCoef := 3.73000001907
	variance := 0.1080000028
	spCoef := 0.37700000405

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

		BonusCoefficient: spCoef,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcAndRollDamageRange(sim, scalingCoef, variance)

			damageMultiplier := spell.DamageMultiplier
			spell.DamageMultiplier *= float64(paladin.DynamicHolyPowerSpent)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier = damageMultiplier

			if result.Landed() {
				paladin.HolyPower.SpendUpTo(paladin.DynamicHolyPowerSpent, actionID, sim)
			}

			spell.DealOutcome(sim, result)
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
