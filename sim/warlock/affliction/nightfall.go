package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (affliction *AfflictionWarlock) registerNightfall() {
	buff := affliction.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 17941},
		Label:    "Shadow Trance",
		Duration: time.Second * 6,
	})

	core.MakeProcTriggerAura(&affliction.Unit, core.ProcTrigger{
		Name:           "Nightfall",
		ClassSpellMask: warlock.WarlockSpellCorruption,
		Callback:       core.CallbackOnPeriodicDamageDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dot := spell.Dot(result.Target)
			if dot == nil || dot != affliction.LastCorruption || !sim.Proc(0.1, "Nightfall Proc") {
				return
			}

			affliction.SoulShards.Gain(sim, 1, buff.ActionID)
			buff.Activate(sim)
		},
	})
}
