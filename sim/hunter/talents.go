package hunter

import "github.com/wowsims/mop/sim/core"

func (hunter *Hunter) applyThrillOfTheHunt() {
	if !hunter.Talents.ThrillOfTheHunt {
		return
	}

	//procChance := 0.3
	//focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34499})

	hunter.RegisterAura(core.Aura{
		Label:    "Thrill of the Hunt",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// mask 256
			//Todo: Check focus cost
			// if spell == hunter.ArcaneShot || spell.ClassSpellMask == HunterSpellExplosiveShot || spell == hunter.BlackArrow {
			// 	if sim.Proc(procChance, "ThrillOfTheHunt") {
			// 		hunter.AddFocus(sim, spell.DefaultCast.Cost*0.4, focusMetrics)
			// 	}
			// }
		},
	})
}
