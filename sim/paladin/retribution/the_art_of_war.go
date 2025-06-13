package retribution

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

/*
Your autoattacks have a 20% chance of resetting the cooldown of your Exorcism.
(Proc chance: 20%)
*/
func (ret *RetributionPaladin) registerArtOfWar() {
	ret.TheArtOfWarAura = ret.RegisterAura(core.Aura{
		Label:    "The Art Of War" + ret.Label,
		ActionID: core.ActionID{SpellID: 59578},
		Duration: time.Second * 6,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ret.Exorcism.CD.Reset()
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:          core.CallbackOnCastComplete,
		ClassSpellMask:    paladin.SpellMaskExorcism,
		SpellFlagsExclude: core.SpellFlagPassiveSpell,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.TheArtOfWarAura.Deactivate(sim)
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
			ret.TheArtOfWarAura.Activate(sim)
		},
	})
}
