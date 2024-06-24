package paladin

import (
	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) ActivateRighteousFury() {
	paladin.RighteousFuryAura = paladin.RegisterAura(core.Aura{
		Label:    "Righteous Fury",
		ActionID: core.ActionID{SpellID: 25780},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.ThreatMultiplier *= 5.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.ThreatMultiplier /= 5.0
		},
	})
}
