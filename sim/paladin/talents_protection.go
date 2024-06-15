package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"time"
)

func (paladin *Paladin) applyProtectionTalents() {
	paladin.applySealsOfThePure()
	paladin.applyShieldOfTheTemplar()
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

func (paladin *Paladin) applyShieldOfTheTemplar() {
	if paladin.Talents.ShieldOfTheTemplar == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskGuardianOfAncientKings,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * time.Duration(40*paladin.Talents.ShieldOfTheTemplar)),
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskAvengingWrath,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * time.Duration(20*paladin.Talents.ShieldOfTheTemplar)),
	})
}
