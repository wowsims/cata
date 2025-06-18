package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (sinRogue *AssassinationRogue) applySealFate() {
	cpMetrics := sinRogue.NewComboPointMetrics(core.ActionID{SpellID: 14190})

	icd := core.Cooldown{
		Timer:    sinRogue.NewTimer(),
		Duration: 500 * time.Millisecond,
	}

	core.MakePermanent(sinRogue.RegisterAura(core.Aura{
		Label: "Seal Fate",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(rogue.SpellFlagBuilder | rogue.SpellFlagSealFate) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if icd.IsReady(sim) {
				sinRogue.AddComboPointsOrAnticipation(sim, 1, cpMetrics)
				icd.Use(sim)

				if sinRogue.T16EnergyAura != nil {
					sinRogue.T16EnergyAura.Activate(sim)
					sinRogue.T16EnergyAura.AddStack(sim)
				}
			}
		},
	}))
}
