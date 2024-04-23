package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ArmsWarrior) ApplyGlyphs() {
	war.Warrior.ApplyGlyphs()

	war.applyGlyphOfBladestorm()
	war.applyGlyphOfMortalStrike()
	war.applyGlyphOfSweepingStrikes()
}

func (war *ArmsWarrior) applyGlyphOfBladestorm() {
	if !war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfBladestorm) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskBladestorm,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -15 * time.Second,
	})
}

func (war *ArmsWarrior) applyGlyphOfMortalStrike() {
	if !war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfMortalStrike) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskMortalStrike,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.1,
	})
}

func (war *ArmsWarrior) applyGlyphOfSweepingStrikes() {
	if !war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfSweepingStrikes) {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskSweepingStrikes,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -1.0,
	})
}
