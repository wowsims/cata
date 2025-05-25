package protection

import (
	"github.com/wowsims/mop/sim/core"
)

func (prot *ProtectionPaladin) registerRighteousFury() {
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:    "Righteous Fury" + prot.Label,
		ActionID: core.ActionID{SpellID: 25780},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.PseudoStats.ThreatMultiplier *= 5.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.PseudoStats.ThreatMultiplier /= 5.0
		},
	}))
}
