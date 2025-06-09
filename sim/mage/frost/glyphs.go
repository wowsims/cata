package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

func (frost *FrostMage) registerGlyphs() {
	if frost.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfWaterElemental) {
		frost.waterElemental.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_AllowCastWhileMoving,
			ClassMask: mage.MageWaterElementalSpellWaterBolt,
		})
	}
}
