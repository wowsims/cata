package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) applyGlyphs() {
	// Majors MOP

	// Glyph of Frostfire Bolt
	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostfireBolt) {
		mage.arcanePowerGCDmod = mage.AddDynamicMod(core.SpellModConfig{
			ClassMask:  mage.MageSpellFrostfireBolt,
			FloatValue: -500 * time.Millisecond,
			Kind:       core.SpellMod_CastTime_Flat,
		})
	}

}
