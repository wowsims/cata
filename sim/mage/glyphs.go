package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) applyGlyphs() {
	//Primes
	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfArcaneBarrage) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  MageSpellArcaneBarrage,
			FloatValue: 0.04,
		})
	}

	// Arcane Blast handled in spell due to handling stacks

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfArcaneMissiles) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  MageSpellArcaneMissilesTick,
			FloatValue: 5,
		})
	}

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfConeOfCold) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  MageSpellConeOfCold,
			FloatValue: 0.25,
		})
	}

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfDeepFreeze) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  MageSpellDeepFreeze,
			FloatValue: 0.2,
		})
	}

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfFireball) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  MageSpellFireball,
			FloatValue: 5,
		})
	}

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfFrostbolt) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  MageSpellFrostbolt,
			FloatValue: 5,
		})
	}

	//Frostfire bolt Dot handled inside spell due to changing behavior
	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfFrostfire) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  MageSpellFrostfireBolt,
			FloatValue: 0.15,
		})
	}

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfIceLance) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  MageSpellIceLance,
			FloatValue: .05,
		})
	}

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfLivingBomb) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  MageSpellLivingBomb,
			FloatValue: .03,
		})
	}

	//Molten Armor Glyph added inside the armor spell in order to reflect on sheet
	/* 	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfMoltenArmor) && mage.Options.Armor == proto.MageOptions_MoltenArmor {
		mage.moltenArmorMod.UpdateFloatValue(5 * core.CritRatingPerCritChance)
	} */

	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfPyroblast) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_BonusCrit_Percent,
			ClassMask:  MageSpellPyroblast | MageSpellPyroblastDot,
			FloatValue: 5,
		})
	}

	// Majors

	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcanePower) {
		mage.arcanePowerGCDmod = mage.AddDynamicMod(core.SpellModConfig{
			ClassMask: MageSpellMirrorImage,
			TimeValue: time.Millisecond * -1500,
			Kind:      core.SpellMod_GlobalCooldown_Flat,
		})
	}

	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfDragonsBreath) {
		mage.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_Cooldown_Flat,
			ClassMask: MageSpellDragonsBreath,
			TimeValue: -3 * time.Second,
		})
	}

	// Minors

	// Mirror Images added inside pet's rotation

}
