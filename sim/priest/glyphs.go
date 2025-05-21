package priest

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (priest *Priest) ApplyGlyphs() {
	// Glyph of Dispersion
	// Glyph of Mindspike
	// Glyph of Shadow Word Death
	if priest.HasMinorGlyph(proto.PriestMinorGlyph_GlyphOfTheSha) {
		priest.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_GlobalCooldown_Flat,
			TimeValue: -core.GCDDefault,
			ClassMask: PriestSpellMindBender | PriestSpellShadowFiend,
		})
	}
}
