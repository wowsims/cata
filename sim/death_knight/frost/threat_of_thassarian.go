package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (fdk *FrostDeathKnight) registerThreatOfThassarian() {
	core.MakePermanent(fdk.RegisterAura(core.Aura{
		Label:    "Threat of Thassarian" + fdk.Label,
		ActionID: core.ActionID{SpellID: 66192},
	})).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.DeathKnightSpellFrostStrike,
		FloatValue: 0.5,
	})
}
