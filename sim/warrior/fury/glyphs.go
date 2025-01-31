package fury

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *FuryWarrior) ApplyGlyphs() {
	war.Warrior.ApplyGlyphs()

	war.applyGlyphOfBloodthirst()
	war.applyGlyphOfRagingBlow()
}

func (war *FuryWarrior) applyGlyphOfBloodthirst() {
	if !war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfBloodthirst) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskBloodthirst,
		Kind:      core.SpellMod_DamageDone_Flat,
		IntValue:  10,
	})
}

func (war *FuryWarrior) applyGlyphOfRagingBlow() {
	if !war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfRagingBlow) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskRagingBlow,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5,
	})
}
