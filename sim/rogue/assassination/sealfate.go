package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/rogue"
)

func (sinRogue *AssassinationRogue) applySealFate() {
	if sinRogue.Talents.SealFate == 0 {
		return
	}

	procChance := 0.2 * float64(sinRogue.Talents.SealFate)
	cpMetrics := sinRogue.NewComboPointMetrics(core.ActionID{SpellID: 14190})

	icd := core.Cooldown{
		Timer:    sinRogue.NewTimer(),
		Duration: 500 * time.Millisecond,
	}

	sinRogue.RegisterAura(core.Aura{
		Label:    "Seal Fate",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(rogue.SpellFlagBuilder) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if icd.IsReady(sim) && sim.Proc(procChance, "Seal Fate") {
				sinRogue.AddComboPoints(sim, 1, cpMetrics)
				icd.Use(sim)
			}
		},
	})
}
