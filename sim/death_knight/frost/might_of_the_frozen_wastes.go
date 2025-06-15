package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/death_knight"
)

func (fdk *FrostDeathKnight) registerMightOfTheFrozenWastes() {
	checkWeaponType := func(sim *core.Simulation, aura *core.Aura) {
		mhWeapon := fdk.GetMHWeapon()
		if mhWeapon != nil && mhWeapon.HandType == proto.HandType_HandTypeTwoHand {
			aura.Activate(sim)
		} else {
			aura.Deactivate(sim)
		}
	}

	mightOfTheFrozenWastesAura := fdk.RegisterAura(core.Aura{
		Label:    "Might of the Frozen Wastes" + fdk.Label,
		ActionID: core.ActionID{SpellID: 81333},
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			checkWeaponType(sim, aura)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.3,
		ProcMask:   core.ProcMaskMelee,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.4,
		ClassMask:  death_knight.DeathKnightSpellObliterate,
	})

	fdk.RegisterItemSwapCallback(core.MeleeWeaponSlots(), func(sim *core.Simulation, _ proto.ItemSlot) {
		checkWeaponType(sim, mightOfTheFrozenWastesAura)
	})
}
