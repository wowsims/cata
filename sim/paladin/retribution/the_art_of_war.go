package retribution

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

func (ret *RetributionPaladin) registerArtOfWar() {
	var artOfWarAura *core.Aura
	artOfWarAura = ret.RegisterAura(core.Aura{
		Label:    "The Art Of War" + ret.Label,
		ActionID: core.ActionID{SpellID: 59578},
		Duration: time.Second * 6,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ret.Exorcism.CD.Reset()
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: paladin.SpellMaskExorcism,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			artOfWarAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:       "The Art of War Trigger" + ret.Label,
		ActionID:   core.ActionID{SpellID: 87138},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			artOfWarAura.Activate(sim)
		},
	})
}
