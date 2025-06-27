package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (war *Warrior) applyMajorGlyphs() {
	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRecklessness) {
		war.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskRecklessness,
			Kind:      core.SpellMod_Custom,
			ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
				spell.RelatedSelfBuff.AttachSpellMod(core.SpellModConfig{
					ProcMask:   core.ProcMaskMeleeSpecial,
					Kind:       core.SpellMod_BonusCrit_Percent,
					FloatValue: -12,
				})
			},
			RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {},
		})

		war.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskRecklessness,
			Kind:      core.SpellMod_BuffDuration_Flat,
			TimeValue: 6 * time.Second,
		})
	}

	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfDeathFromAbove) {
		war.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskHeroicLeap,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: -15 * time.Second,
		})
	}

	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfSweepingStrikes) {
		actionID := core.ActionID{SpellID: 58384}
		rageMetrics := war.NewRageMetrics(actionID)
		core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
			Name:           "Glyph of Sweeping Strikes",
			ActionID:       actionID,
			ClassSpellMask: SpellMaskSweepingStrikesHit,
			Callback:       core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				war.AddRage(sim, 1, rageMetrics)
			},
		})
	}

	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfResonatingPower) {
		war.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskThunderClap,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.5,
		})

		war.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskThunderClap,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: 3 * time.Second,
		})
	}

	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfIncite) {
		actionID := core.ActionID{SpellID: 122016}

		war.InciteAura = war.RegisterAura(core.Aura{
			Label:     "Incite",
			ActionID:  actionID,
			Duration:  10 * time.Second,
			MaxStacks: 3,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				war.HeroicStrikeCleaveCostMod.Activate()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if !war.UltimatumAura.IsActive() {
					war.HeroicStrikeCleaveCostMod.Deactivate()
				}
			},
		}).AttachProcTrigger(core.ProcTrigger{
			Name:           "Incite - Consume",
			ClassSpellMask: SpellMaskHeroicStrike | SpellMaskCleave,
			Callback:       core.CallbackOnSpellHitDealt,

			ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				return spell.CurCast.Cost <= 0 && !war.UltimatumAura.IsActive()
			},

			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				war.InciteAura.RemoveStack(sim)
			},
		})

		core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
			Name:           "Incite - Trigger",
			ClassSpellMask: SpellMaskDemoralizingShout,
			Callback:       core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				war.InciteAura.Activate(sim)
				war.InciteAura.SetStacks(sim, 3)
			},
		})
	}

	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRagingWind) {
		actionID := core.ActionID{SpellID: 115317}
		ragingWindAura := war.RegisterAura(core.Aura{
			Label:    "Raging Wind",
			ActionID: actionID,
			Duration: 6 * time.Second,
		}).AttachSpellMod(core.SpellModConfig{
			ClassMask:  SpellMaskWhirlwind | SpellMaskWhirlwindOh,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.1,
		})

		core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
			Name:           "Raging Wind - Consume",
			ClassSpellMask: SpellMaskWhirlwind,
			Callback:       core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				ragingWindAura.Deactivate(sim)
			},
		})

		core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
			Name:           "Raging Wind - Trigger",
			ClassSpellMask: SpellMaskRagingBlowMH,
			Callback:       core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				ragingWindAura.Activate(sim)
			},
		})
	}

	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfHoldTheLine) {
		actionID := core.ActionID{SpellID: 84619}

		holdTheLine := war.RegisterAura(core.Aura{
			Label:    "Hold the Line",
			ActionID: actionID,
			Duration: 5 * time.Second,
		}).AttachSpellMod(core.SpellModConfig{
			ClassMask:  SpellMaskRevenge,
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.5,
		})

		core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
			Name:           "Hold the Line - Consume",
			ClassSpellMask: SpellMaskRevenge,
			Callback:       core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				holdTheLine.Deactivate(sim)
			},
		})

		core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
			Name:     "Hold the Line - Trigger",
			Callback: core.CallbackOnSpellHitTaken,
			Outcome:  core.OutcomeParry,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				holdTheLine.Activate(sim)
			},
		})
	}

	if war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfHeavyRepercussions) {
		war.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskShieldBlock,
			Kind:      core.SpellMod_Custom,
			ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
				war.ShieldBlockAura.AttachSpellMod(core.SpellModConfig{
					ClassMask:  SpellMaskShieldSlam,
					Kind:       core.SpellMod_DamageDone_Pct,
					FloatValue: 0.5,
				})
			},
			RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {},
		})
	}

}

func (war *Warrior) applyMinorGlyphs() {

}

func (war *Warrior) ApplyGlyphs() {
	war.applyMajorGlyphs()
	war.applyMinorGlyphs()
}
