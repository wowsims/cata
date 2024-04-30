package subtlety

import "github.com/wowsims/cata/sim/core"

func (subRogue *SubtletyRogue) applySerratedBlades() {
	chancePerPoint := 0.1 * float64(subRogue.Talents.SerratedBlades)

	subRogue.RegisterAura(core.Aura{
		Label:    "Serrated Blades",
		Duration: core.NeverExpires,

		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Unit == &subRogue.Unit && spell == subRogue.Eviscerate {
				procChance := float64(subRogue.ComboPoints()) * chancePerPoint
				if sim.Proc(procChance, "Serrated Blades") {
					rupAura := subRogue.Rupture.Dot(result.Target)
					if rupAura.IsActive() {
						// println(rupAura.Duration)
						rupAura.Activate(sim)
					}
				}
			}
		},
	})
}
