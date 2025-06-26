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
	var frenzyStacksSnapshot int32
	focusFireAura := bmHunter.RegisterAura(core.Aura{
		Label:    "Focus Fire",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			frenzyStacksSnapshot = bmHunter.Pet.FrenzyAura.GetStacks()
			if frenzyStacksSnapshot >= 1 {
				bmHunter.Pet.FrenzyAura.Deactivate(sim)
				bmHunter.Pet.AddFocus(sim, 6*float64(frenzyStacksSnapshot), petFocusMetrics)
				bmHunter.MultiplyRangedHaste(sim, 1+(float64(frenzyStacksSnapshot)*0.06))
				if sim.Log != nil {
					bmHunter.Pet.Log(sim, "Consumed %d stacks of Frenzy for Focus Fire.", frenzyStacksSnapshot)
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if frenzyStacksSnapshot > 0 {
				bmHunter.MultiplyRangedHaste(sim, 1/(1+(float64(frenzyStacksSnapshot)*0.06)))
			}
			frenzyStacksSnapshot = 0
		},
	})

	focusFireSpell := bmHunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return bmHunter.Pet.FrenzyAura.GetStacks() > 0 && !focusFireAura.IsActive()
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
