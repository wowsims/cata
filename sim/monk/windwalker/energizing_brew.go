package windwalker

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (ww *WindwalkerMonk) registerEnergizingBrew() {
	actionID := core.ActionID{SpellID: 115288}
	energyMetrics := ww.NewEnergyMetrics(actionID)

	energizingBrewAura := ww.RegisterAura(core.Aura{
		Label:    "Energizing Brew" + ww.Label,
		ActionID: actionID,
		Duration: time.Second * 6,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 1,
				NumTicks: int(aura.Duration.Seconds()),
				OnAction: func(sim *core.Simulation) {
					ww.AddEnergy(sim, 10, energyMetrics)
				},
			})
		},
	})

	ww.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: monk.MonkSpellEnergizingBrew,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    ww.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			energizingBrewAura.Activate(sim)
		},
		RelatedSelfBuff: energizingBrewAura,
	})
}
