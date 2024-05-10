package druid

import (
	"time"

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

	if druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfInsectSwarm) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask:  DruidSpellInsectSwarm,
			FloatValue: 0.3,
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	if druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfWrath) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask:  DruidSpellWrath,
			FloatValue: 0.1,
			Kind:       core.SpellMod_DamageDone_Flat,
		})
	}

	if druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfall) {
		druid.AddStaticMod(core.SpellModConfig{
			ClassMask: DruidSpellStarfall,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Second * -30,
		})
	}
}
