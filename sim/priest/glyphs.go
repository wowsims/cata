package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (priest *Priest) ApplyGlyphs() {

	// Glyph of Dispersion
	// Glyph of Mindspike
	// Glyph of Shadow Word Death
	if priest.HasMinorGlyph(proto.PriestMinorGlyph_GlyphOfTheSha) {
		priest.OnSpellRegistered(func(spell *core.Spell) {
			if spell.ClassSpellMask&(PriestSpellMindBender|PriestSpellShadowFiend) > 0 {
				spell.DefaultCast.GCD = time.Duration(0)
			}
		})
	}
}
