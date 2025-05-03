package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (asnRogue *AssassinationRogue) registerAllPassives() {
	asnRogue.registerBlindsidePassive()
}

func (asnRogue *AssassinationRogue) registerBlindsidePassive() {
	blindsideProc := asnRogue.RegisterAura(core.Aura{
		Label:    "Blindside",
		ActionID: core.ActionID{SpellID: 121153},
		Duration: time.Second * 10,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && spell.SpellID == 111240 {
				// Dispatch casted, consume aura
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(core.MakeProcTriggerAura(&asnRogue.Unit, core.ProcTrigger{
		Name:           "Blindside Proc Trigger",
		ActionID:       core.ActionID{ItemID: 121152},
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: rogue.RogueSpellMutilate,
		ProcChance:     0.3,
		Outcome:        core.OutcomeHit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			blindsideProc.Activate(sim)
		},
	}))
}
