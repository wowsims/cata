package druid

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) ApplyGlyphs() {

	if druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfMoonfire) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask:  DruidSpellMoonfireDoT | DruidSpellSunfireDoT,
			FloatValue: 0.2,
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}
}
