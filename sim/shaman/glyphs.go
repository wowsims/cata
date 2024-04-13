package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (shaman *Shaman) ApplyGlyphs() {

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfFireElementalTotem) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: int64(SpellMaskFireElementalTotem),
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Minute * -5,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfFlameShock) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: int64(SpellMaskFlameShock),
			Kind:      core.SpellMod_Dot_NumberOfTicks_Flat,
			IntValue:  3,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfLavaBurst) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  int64(SpellMaskLavaBurst),
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.1,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfLavaLash) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  int64(SpellMaskLavaLash),
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.2,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfLightningBolt) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  int64(SpellMaskLightningBolt),
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.04,
		})
	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfThunder) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: int64(SpellMaskThunderstorm),
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -10,
		})
	}
}
