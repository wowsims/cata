package paladin

import "github.com/wowsims/cata/sim/core"

func (paladin *Paladin) applyProtectionTalents() {
	paladin.applySealsOfThePure()
}

func (paladin *Paladin) applySealsOfThePure() {
	if paladin.Talents.SealsOfThePure == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskSealOfRighteousness | SpellMaskSealOfTruth | SpellMaskSealOfJustice,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.06 * float64(paladin.Talents.SealsOfThePure),
	})
}
