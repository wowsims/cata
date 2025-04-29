package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) ApplyGlyphs() {
	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfFireElementalTotem) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskFireElementalTotem,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Minute * -5,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfFlameShock) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskFlameShock,
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			IntValue:  3,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfLavaBurst) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskLavaBurst | SpellMaskLavaBurstOverload,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.1,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfLavaLash) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskLavaLash,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.2,
		})
	}

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfLightningBolt) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskLightningBolt,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.04,
		})
	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfThunder) {
		shaman.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskThunderstorm,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -10,
		})
	}
}
