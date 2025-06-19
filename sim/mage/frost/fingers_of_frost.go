package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (frost *FrostMage) registerFingersOfFrost() {
	/*
		Ice Lance does 4x Damage against Frozen enemies, FoF adds a bonus 25%.
		This effect affects Deep Freeze as well but does no damage so it's been ommitted. \
		https://www.wowhead.com/mop-classic/spell=30455/ice-lance and https://www.wowhead.com/mop-classic/spell=112965/fingers-of-frost for more information.
	*/

	frost.FingersOfFrostAura = frost.RegisterAura(core.Aura{
		Label:     "Fingers of Frost",
		ActionID:  core.ActionID{SpellID: 112965},
		Duration:  time.Second * 15,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			frost.iceLanceFrozenCritBuffMod.UpdateFloatValue(frost.GetFrozenCritPercentage())
			frost.iceLanceFrozenCritBuffMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			frost.iceLanceFrozenCritBuffMod.Deactivate()
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 3.0,
		ClassMask:  mage.MageSpellIceLance,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.25,
		ClassMask:  mage.MageSpellIceLance,
	})
}
