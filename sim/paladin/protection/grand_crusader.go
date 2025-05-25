package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

func (prot *ProtectionPaladin) registerGrandCrusader() {
	hpActionID := core.ActionID{SpellID: 98057}
	prot.CanTriggerHolyAvengerHpGain(hpActionID)

	var grandCrusaderAura *core.Aura
	grandCrusaderAura = prot.RegisterAura(core.Aura{
		Label:    "Grand Crusader" + prot.Label,
		ActionID: core.ActionID{SpellID: 85416},
		Duration: time.Second * 6,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: paladin.SpellMaskAvengersShield,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			prot.HolyPower.Gain(1, hpActionID, sim)
			grandCrusaderAura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&prot.Unit, core.ProcTrigger{
		Name:       "Grand Crusader Trigger" + prot.Label,
		ActionID:   core.ActionID{SpellID: 85043},
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeDodge | core.OutcomeParry,
		ProcChance: 0.3,
		ICD:        time.Second,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			prot.AvengersShield.CD.Reset()
			grandCrusaderAura.Activate(sim)
		},
	})
}
