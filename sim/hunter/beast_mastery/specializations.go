package beast_mastery

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (bmHunter *BeastMasteryHunter) ApplyTalents() {
	bmHunter.applyFrenzy()
	bmHunter.Hunter.ApplyTalents()
}

func (bmHunter *BeastMasteryHunter) applyFrenzy() {
	if bmHunter.Pet == nil {
		return
	}
	actionId := core.ActionID{SpellID: 19623}
	bmHunter.Pet.FrenzyAura = bmHunter.Pet.RegisterAura(core.Aura{
		Label:     "Frenzy",
		Duration:  time.Second * 30,
		ActionID:  actionId,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			aura.Unit.MultiplyMeleeSpeed(sim, 1/(1+0.04*float64(oldStacks)))
			aura.Unit.MultiplyMeleeSpeed(sim, 1+0.04*float64(newStacks))
		},
	})

	bmHunter.Pet.RegisterAura(core.Aura{
		Label:    "FrenzyHandler",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				return
			}
			if bmHunter.Pet.FrenzyAura.IsActive() {
				if bmHunter.Pet.FrenzyAura.GetStacks() != 5 {
					bmHunter.Pet.FrenzyAura.AddStack(sim)
				}
				bmHunter.Pet.FrenzyAura.Refresh(sim)
			} else {
				bmHunter.Pet.FrenzyAura.Activate(sim)
				bmHunter.Pet.FrenzyAura.SetStacks(sim, 1)
			}
		},
	})
}
