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
}
