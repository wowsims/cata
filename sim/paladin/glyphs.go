package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (paladin *Paladin) ApplyGlyphs() {
	// Prime Glyphs
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfCrusaderStrike) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Rating,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: 5 * core.CritRatingPerCritChance,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfJudgement) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  SpellMaskJudgement,
			FloatValue: 0.1,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfTemplarSVerdict) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  SpellMaskTemplarsVerdict,
			FloatValue: 0.15,
		})
	}
	if paladin.HasPrimeGlyph(proto.PaladinPrimeGlyph_GlyphOfSealOfTruth) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusExpertise_Rating,
			ClassMask:  SpellMaskSealOfTruth,
			FloatValue: 10 * core.ExpertisePerQuarterPercentReduction,
		})
	}

	// Major Glyphs
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfHammerOfWrath) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Pct,
			ClassMask:  SpellMaskHammerOfWrath,
			FloatValue: -1,
		})
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_Cooldown_Multiplier,
			ClassMask:  SpellMaskConsecration,
			FloatValue: 0.2,
		})
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DotNumberOfTicks_Flat,
			ClassMask:  SpellMaskConsecration,
			FloatValue: 2,
		})
	}
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfTheAsceticCrusader) {
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Pct,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: -0.3,
		})
	}
}
