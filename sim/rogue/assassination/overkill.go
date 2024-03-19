package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

var OverkillActionID = core.ActionID{SpellID: 58427}

func (sinRogue *AssassinationRogue) registerOverkill() {
	if !sinRogue.Talents.Overkill {
		return
	}

	effectDuration := time.Second * 20
	if sinRogue.StealthAura.IsActive() {
		effectDuration = core.NeverExpires
	}

	sinRogue.OverkillAura = sinRogue.RegisterAura(core.Aura{
		Label:    "Overkill",
		ActionID: OverkillActionID,
		Duration: effectDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			sinRogue.ApplyEnergyTickMultiplier(0.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			sinRogue.ApplyEnergyTickMultiplier(-0.3)
		},
	})
}
