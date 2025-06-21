package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerGlyphs() {
	// Majors MOP

	// Glyph of Frostfire Bolt
	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostfireBolt) {
		mage.AddDynamicMod(core.SpellModConfig{
			ClassMask: MageSpellFrostfireBolt,
			TimeValue: time.Millisecond * -500,
			Kind:      core.SpellMod_CastTime_Flat,
		})
	}

	// Glyph of Cone of Cold
	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfConeOfCold) {
		mage.AddDynamicMod(core.SpellModConfig{
			ClassMask:  MageSpellConeOfCold,
			FloatValue: 2.0,
			Kind:       core.SpellMod_DamageDone_Pct,
		})
	}

	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfWaterElemental) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_AllowCastWhileMoving,
			ClassMask: MageWaterElementalSpellWaterBolt,
		})
	}

}
