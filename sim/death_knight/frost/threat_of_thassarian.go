package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/death_knight"
)

func (fdk *FrostDeathKnight) registerThreatOfThassarian() {
	checkWeaponType := func(sim *core.Simulation, aura *core.Aura) {
		if fdk.HasMHWeapon() && fdk.HasOHWeapon() {
			aura.Activate(sim)
		} else {
			aura.Deactivate(sim)
		}
	}

	fdk.ThreatOfThassarianAura = fdk.RegisterAura(core.Aura{
		Label:    "Threat of Thassarian" + fdk.Label,
		ActionID: core.ActionID{SpellID: 66192},
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			checkWeaponType(sim, aura)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.DeathKnightSpellFrostStrike,
		FloatValue: 0.5,
	})

	fdk.RegisterItemSwapCallback(core.AllWeaponSlots(), func(sim *core.Simulation, _ proto.ItemSlot) {
		checkWeaponType(sim, fdk.ThreatOfThassarianAura)
	})
}
