package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warrior *Warrior) applyPrimeGlyphs() {
	if warrior.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfRevenge) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskRevenge,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.1,
		})
	}

	if warrior.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfSlam) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskSlam,
			Kind:       core.SpellMod_BonusCrit_Rating,
			FloatValue: 5 * core.CritRatingPerCritChance,
		})
	}

	if warrior.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfOverpower) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskOverpower,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.1,
		})
	}
}

func (warrior *Warrior) applyMajorGlyphs() {
	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfShieldWall) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskShieldWall,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: time.Minute * 2,
		})
	}

	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfResonatingPower) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskThunderClap,
			Kind:       core.SpellMod_PowerCost_Flat,
			FloatValue: -5,
		})
	}
}

func (warrior *Warrior) applyMinorGlyphs() {
	// Since they're raid buffs/debuffs, shouts and their glyph effects are handled in buffs.go and debuffs.go

	if warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfFuriousSundering) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskSunderArmor,
			Kind:       core.SpellMod_PowerCost_Pct,
			FloatValue: 0.5,
		})
	}

	if warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfShatteringThrow) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask:  SpellMaskShatteringThrow,
			Kind:       core.SpellMod_CastTime_Pct,
			FloatValue: -1.0,
		})
	}
}

func (warrior *Warrior) ApplyGlyphs() {
	warrior.applyPrimeGlyphs()
	warrior.applyMajorGlyphs()
	warrior.applyMinorGlyphs()
}
