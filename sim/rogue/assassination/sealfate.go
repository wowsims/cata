package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (sinRogue *AssassinationRogue) applySealFate() {
	if sinRogue.Talents.SealFate == 0 {
		return
	}

	procChance := 0.5 * float64(sinRogue.Talents.SealFate)
	cpMetrics := sinRogue.NewComboPointMetrics(core.ActionID{SpellID: 14190})

	icd := core.Cooldown{
		Timer:    sinRogue.NewTimer(),
		Duration: 500 * time.Millisecond,
	}

	core.MakePermanent(sinRogue.RegisterAura(core.Aura{
		Label: "Seal Fate",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(rogue.SpellFlagBuilder) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if icd.IsReady(sim) && (procChance == 1 || sim.Proc(procChance, "Seal Fate")) {
				sinRogue.AddComboPoints(sim, 1, cpMetrics)
				icd.Use(sim)
			}
		},
	}))
}
