package protection

import (
	"github.com/wowsims/mop/sim/core"
)

// Increases your threat generation while active, making you a more effective tank.
func (prot *ProtectionPaladin) registerRighteousFury() {
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:    "Righteous Fury" + prot.Label,
		ActionID: core.ActionID{SpellID: 25780},

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.PseudoStats.ThreatMultiplier *= 7.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.PseudoStats.ThreatMultiplier /= 7.0
		},
	}))
}
