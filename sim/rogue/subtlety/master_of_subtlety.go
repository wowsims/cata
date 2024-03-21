package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (subRogue *SubtletyRogue) registerMasterOfSubtletyCD() {
	var MasterOfSubtletyID = core.ActionID{SpellID: 31223}

	percent := 1.1

	effectDuration := time.Second * 6
	if subRogue.StealthAura.IsActive() {
		effectDuration = core.NeverExpires
	}

	subRogue.MasterOfSubtletyAura = subRogue.RegisterAura(core.Aura{
		Label:    "Master of Subtlety",
		ActionID: MasterOfSubtletyID,
		Duration: effectDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			subRogue.PseudoStats.DamageDealtMultiplier *= percent
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			subRogue.PseudoStats.DamageDealtMultiplier *= 1 / percent
		},
	})
}
