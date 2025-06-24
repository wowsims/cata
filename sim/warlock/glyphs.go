package warlock

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warlock *Warlock) registerEternalResolve() {
	if !warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfEternalResolve) {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellAgony | WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDoom,
		Kind:       core.SpellMod_DotBaseDuration_Pct,
		FloatValue: 0.5,
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellAgony | WarlockSpellCorruption | WarlockSpellUnstableAffliction | WarlockSpellDoom,
		Kind:       core.SpellMod_DotDamageDone_Pct,
		FloatValue: -0.2,
	})
}
