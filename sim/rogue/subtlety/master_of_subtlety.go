package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (subRogue *SubtletyRogue) registerMasterOfSubtletyCD() {
	var MasterOfSubtletyID = core.ActionID{SpellID: 31223}

	subRogue.MasterOfSubtletyAura = subRogue.RegisterAura(core.Aura{
		Label:    "Master of Subtlety",
		ActionID: MasterOfSubtletyID,
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			subRogue.PseudoStats.DamageDealtMultiplier *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			subRogue.PseudoStats.DamageDealtMultiplier /= 1.1
		},
	})
}
