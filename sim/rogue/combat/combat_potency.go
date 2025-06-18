package combat

import "github.com/wowsims/mop/sim/core"

func (comRogue *CombatRogue) applyCombatPotency() {
	energyBonus := 15.0
	energyMetrics := comRogue.NewEnergyMetrics(core.ActionID{SpellID: 35546})

	comRogue.RegisterAura(core.Aura{
		Label:    "Combat Potency",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && (spell.ProcMask.Matches(core.ProcMaskMeleeOHAuto) || spell.SpellID == 86392) { // 86392 = Main Gauche
				procChance := 0.2
				if spell.ProcMask.Matches(core.ProcMaskMeleeOHAuto) {
					ohSpeed := comRogue.GetOHWeapon().SwingSpeed
					procChance = (20 * ohSpeed / 1.4) / 100
				}

				if sim.RandomFloat("Combat Potency") < procChance {
					comRogue.AddEnergy(sim, energyBonus, energyMetrics)
				}
			}
		},
	})
}
