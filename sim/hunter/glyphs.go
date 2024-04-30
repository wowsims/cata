package hunter

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) ApplyGlyphs() {
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfArcaneShot) {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  HunterSpellArcaneShot,
			FloatValue: 0.12,
		})
	}
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfExplosiveShot) {
		hunter.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Rating,
			ClassMask:  HunterSpellExplosiveShot,
			FloatValue: 6 * core.CritRatingPerCritChance,
		})
	}
}
