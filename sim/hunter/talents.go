package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) applyThrillOfTheHunt() {
	if !hunter.Talents.ThrillOfTheHunt {
		return
	}

	actionID := core.ActionID{SpellID: 109306}
	procChance := 0.30

	tothMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: HunterSpellMultiShot | HunterSpellArcaneShot,
		IntValue:  -20,
	})

	tothAura := hunter.RegisterAura(core.Aura{
		Label:     "Thrill of the Hunt",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			tothMod.Activate()

		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			tothMod.Deactivate()

		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(HunterSpellMultiShot) || spell.Matches(HunterSpellArcaneShot) {
				aura.RemoveStack(sim)

			}
		},
	})

	hunter.RegisterAura(core.Aura{
		Label:    "Thrill of the Hunt Proccer",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Needs to cost Focus to proc
			if spell.CurCast.Cost <= 0 {
				return
			}

			if sim.RandomFloat("Thrill of the Hunt") < procChance {
				tothAura.Activate(sim)
				tothAura.SetStacks(sim, 3)
			}
		},
	})
}
