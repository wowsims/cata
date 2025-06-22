package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/death_knight"
)

func (fdk *FrostDeathKnight) registerMightOfTheFrozenWastes() {
	fdk.MightOfTheFrozenWastesAura = fdk.RegisterAura(core.Aura{
		Label:    "Might of the Frozen Wastes" + fdk.Label,
		ActionID: core.ActionID{SpellID: 81333},
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			if mh := fdk.GetMHWeapon(); mh != nil && mh.HandType == proto.HandType_HandTypeTwoHand {
				aura.Activate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMelee,
		FloatValue: 0.3,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.DeathKnightSpellObliterate,
		FloatValue: 0.4,
	})
}
