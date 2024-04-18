package fury

import "github.com/wowsims/cata/sim/core"

func (war *FuryWarrior) ApplyGlyphs() {
	war.Warrior.ApplyGlyphs()

	war.applyGlyphOfBloodthirst()
	war.applyGlyphOfRagingBlow()
}

func (war *FuryWarrior) applyGlyphOfBloodthirst() {
	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskBloodthirst,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.1,
	})
}

func (war *FuryWarrior) applyGlyphOfRagingBlow() {
	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskRagingBlow,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5 * core.CritRatingPerCritChance,
	})
}
