package beast_mastery

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (bmHunter *BeastMasteryHunter) registerFocusFireSpell() {
	if bmHunter.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 82692}
	petFocusMetrics := bmHunter.Pet.NewFocusMetrics(actionID)
	focusFireAura := bmHunter.RegisterAura(core.Aura{
		Label:    "Focus Fire",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bmHunter.Pet.FrenzyStacksSnapshot = bmHunter.Pet.FrenzyAura.GetStacks()
			if bmHunter.Pet.FrenzyStacksSnapshot >= 1 {
				bmHunter.Pet.FrenzyAura.Deactivate(sim)
				bmHunter.Pet.AddFocus(sim, 6, petFocusMetrics)
				bmHunter.MultiplyRangedSpeed(sim, 1+(float64(bmHunter.Pet.FrenzyStacksSnapshot)*0.06))
				if sim.Log != nil {
					bmHunter.Pet.Log(sim, "Consumed %d stacks of Frenzy for Focus Fire.", bmHunter.Pet.FrenzyStacksSnapshot)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if bmHunter.Pet.FrenzyStacksSnapshot > 0 {
				bmHunter.MultiplyRangedSpeed(sim, 1/(1+(float64(bmHunter.Pet.FrenzyStacksSnapshot)*0.06)))
			}
		},
	})

	focusFireSpell := bmHunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 1,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bmHunter.Pet.FrenzyAura.GetStacks() >= bmHunter.Pet.FrenzyStacksSnapshot
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if focusFireAura.IsActive() {
				focusFireAura.Deactivate(sim)
			}
			focusFireAura.Activate(sim)
		},
	})

	bmHunter.AddMajorCooldown(core.MajorCooldown{
		Spell: focusFireSpell,
		Type:  core.CooldownTypeDPS,
	})

}
