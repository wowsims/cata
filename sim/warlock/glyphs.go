package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warlock *Warlock) ApplyGlyphs() {
	if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfConflagrate) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellConflagrate,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -2,
		})
	}

	if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfChaosBolt) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellChaosBolt,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -2,
		})
	}

	if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfIncinerate) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellIncinerate,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.05,
		})
	}

	// TODO: only applies to periodic damage
	if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfImmolate) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellImmolate,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.10,
		})
	}

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellLifeTap,
			Kind:       core.SpellMod_GlobalCooldown_Flat,
			FloatValue: 0.5,
		})
	}

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowBolt) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellShadowBolt,
			Kind:       core.SpellMod_PowerCost_Pct,
			FloatValue: -0.15,
		})
	}

	if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfBaneOfAgony) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellBaneOfAgony,
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			IntValue:  2,
		})
	}

	if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfUnstableAffliction) {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellUnstableAffliction,
			Kind:      core.SpellMod_CastTime_Flat,
			TimeValue: time.Millisecond * 200,
		})
	}

	//if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfLashOfPain) {
	// pet.AddStaticMod(core.SpellModConfig{
	// 	ClassMask: WarlockSpellSuccubusLashOfPain,
	// 	Kind:      core.SpellMod_DamageDone_Pct,
	// 	FloatValue: 0.25,
	// })
	//}

	// if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfMetamorphosis) {
	// 	warlock.AddStaticMod(core.SpellModConfig{
	// 		ClassMask: WarlockSpellMetamorphosis,
	// 		Kind:      core.SpellMod_,
	// 		TimeValue: time.Second * 6,
	// 	})
	// }

	// if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfFelguard) {
	// pet.AddStaticMod(core.SpellModConfig{
	// 	ClassMask: WarlockSpellFelGuardLegionStrike,
	// 	Kind:      core.SpellMod_DamageDone_Pct,
	// 	FloatValue: 0.05,
	// })
	// }

	// if warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfImp) {
	// 	pet.AddStaticMod(core.SpellModConfig{
	// 		ClassMask: WarlockSpellImpFireBolt,
	// 		Kind:      core.SpellMod_DamageDone_Pct,
	// 		FloatValue: 0.20,
	// 	})
	// 	}

	//TODO: Soul Swap with spell
}
