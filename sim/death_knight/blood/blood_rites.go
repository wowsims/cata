package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (bdk *BloodDeathKnight) registerBloodRites() {
	core.MakePermanent(bdk.RegisterAura(core.Aura{
		Label:    "Blood Rites" + bdk.Label,
		ActionID: core.ActionID{SpellID: 50034},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.DeathKnightSpellDeathStrike,
		FloatValue: 0.4,
	})
}
