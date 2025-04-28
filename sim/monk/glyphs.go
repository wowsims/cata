package monk

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (monk *Monk) applyGlyphs() {
	if monk.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfFistsOfFury) {
		monk.registerGlyphOfFistsOfFury()
	}

	if monk.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfTargetedExpulsion) {
		monk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  MonkSpellExpelHarm,
			FloatValue: -0.5,
		})
	}
}

func (monk *Monk) registerGlyphOfFistsOfFury() {

	parryBuff := monk.RegisterAura(core.Aura{
		Label:    "Glyph of Fists of Fury" + monk.Label,
		ActionID: core.ActionID{SpellID: 125671},
	}).AttachAdditivePseudoStatBuff(&monk.PseudoStats.BaseParryChance, 1)

	core.MakeProcTriggerAura(&monk.Unit, core.ProcTrigger{
		Name:           "Glyph of Fists of Fury Trigger" + monk.Label,
		ClassSpellMask: MonkSpellFistsOfFury,
		Callback:       core.CallbackOnCastComplete,
		SpellFlags:     SpellFlagSpender,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			parryBuff.Duration = spell.AOEDot().RemainingDuration(sim)
			parryBuff.Activate(sim)
		},
	})
}
