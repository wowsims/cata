package subtlety

import "github.com/wowsims/mop/sim/core"

func (subRogue *SubtletyRogue) applyInitiative() {
	if subRogue.Talents.Initiative == 0 {
		return
	}

	procChance := 0.5 * float64(subRogue.Talents.Initiative)
	cpMetrics := subRogue.NewComboPointMetrics(core.ActionID{SpellID: 13979})

	subRogue.RegisterAura(core.Aura{
		Label:    "Initiative",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == subRogue.Garrote || spell == subRogue.Ambush {
				if result.Landed() {
					if sim.Proc(procChance, "Initiative") {
						subRogue.AddComboPoints(sim, 1, cpMetrics)
					}
				}
			}
		},
	})
}
