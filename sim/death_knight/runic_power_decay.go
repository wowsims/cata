package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (dk *DeathKnight) registerRunicPowerDecay() {
	decayMetrics := dk.NewRunicPowerMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionPrepull})

	var decay *core.PendingAction
	dk.RunicPowerDecayAura = dk.GetOrRegisterAura(core.Aura{
		Label:    "Runic Power Decay" + dk.Label,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if sim.CurrentTime >= 0 || dk.CurrentRunicPower() <= 0 {
				dk.RunicPowerDecayAura.Deactivate(sim)
				return
			}

			dk.SpendRunicPower(sim, 1, decayMetrics)

			decay = &core.PendingAction{
				Priority:     core.ActionPriorityPrePull,
				NextActionAt: sim.CurrentTime + time.Second,
				OnAction: func(sim *core.Simulation) {
					if dk.CurrentRunicPower() <= 0 {
						aura.Deactivate(sim)
						return
					}

					dk.SpendRunicPower(sim, 1, decayMetrics)

					nextTick := sim.CurrentTime + time.Second
					if nextTick >= 0 {
						aura.Deactivate(sim)
						return
					}

					decay.NextActionAt = nextTick
					sim.AddPendingAction(decay)
				},
			}

			sim.AddPendingAction(decay)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if decay != nil {
				decay.Cancel(sim)
			}
		},
	})
}
