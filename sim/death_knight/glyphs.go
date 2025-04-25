package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
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
			ClassMask:  DeathKnightSpellDeathCoil | DeathKnightSpellDeathCoilHeal,
			FloatValue: 0.15,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfScourgeStrike) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellScourgeStrikeShadow,
			FloatValue: 0.3,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfHeartStrike) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellHeartStrike,
			FloatValue: 0.3,
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfDeathStrike) {
		dsMod := dk.AddDynamicMod(core.SpellModConfig{
			Kind:      core.SpellMod_DamageDone_Pct,
			ClassMask: DeathKnightSpellDeathStrike,
		})

		core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
			Name:           "Death Strike Glyph Activate",
			Callback:       core.CallbackOnApplyEffects,
			ClassSpellMask: DeathKnightSpellDeathStrike,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				dsMod.UpdateFloatValue(min(0.4, (dk.CurrentRunicPower()/5.0)*0.02))
			},
		})

		dk.RegisterResetEffect(func(s *core.Simulation) {
			dsMod.Activate()
		})
	}

	if dk.HasPrimeGlyph(proto.DeathKnightPrimeGlyph_GlyphOfRuneStrike) {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  DeathKnightSpellRuneStrike,
			FloatValue: 10,
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
