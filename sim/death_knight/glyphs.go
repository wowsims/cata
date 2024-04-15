package death_knight

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (dk *DeathKnight) ApplyGlyphs() {
	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfDeathAndDecay) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			ClassMask: DeathKnightSpellDeathAndDecay,
			IntValue:  5,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfDeathCoil) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDeathCoil,
			FloatValue: 0.15,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfScourgeStrike) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellScourgeStrikeShadow,
			FloatValue: 0.30,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfIcyTouch) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellFrostFever,
			FloatValue: 0.2,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfObliterate) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellObliterate,
			FloatValue: 0.2,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfFrostStrike) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_RunicPowerCost_Flat,
			ClassMask:  DeathKnightSpellFrostStrike,
			ProcMask:   core.ProcMaskMeleeMH,
			FloatValue: -8,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfHowlingBlast) {
		core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
			Name:           "Howling Blast Disease",
			Callback:       core.CallbackOnSpellHitDealt,
			Outcome:        core.OutcomeLanded,
			ClassSpellMask: DeathKnightSpellHowlingBlast,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				dk.FrostFeverSpell.Cast(sim, result.Target)
			},
		})
	}
}
