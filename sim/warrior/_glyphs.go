package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
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
			Kind:       core.SpellMod_BonusCrit_Percent,
			FloatValue: 5,
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
			ClassMask: SpellMaskThunderClap,
			Kind:      core.SpellMod_PowerCost_Flat,
			IntValue:  -5,
		})
	}

	if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfRapidCharge) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskCharge,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: -time.Second * 1,
		})
	}
}

func (warrior *Warrior) applyMinorGlyphs() {
	// Since they're raid buffs/debuffs, shouts and their glyph effects are handled in buffs.go and debuffs.go

	if warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfFuriousSundering) {
		warrior.AddStaticMod(core.SpellModConfig{
			ClassMask: SpellMaskSunderArmor,
			Kind:      core.SpellMod_PowerCost_Pct,
			IntValue:  -50,
		})
	}
}

func (warrior *Warrior) ApplyGlyphs() {
	warrior.applyPrimeGlyphs()
	warrior.applyMajorGlyphs()
	warrior.applyMinorGlyphs()
}
