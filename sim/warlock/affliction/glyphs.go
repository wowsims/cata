package affliction

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
)

func (affliction *AfflictionWarlock) registerGlpyhs() {

	if affliction.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfUnstableAffliction) {
		affliction.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_CastTime_Pct,
			ClassMask:  warlock.WarlockSpellUnstableAffliction,
			FloatValue: -0.25,
		})
	}
}
