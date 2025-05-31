package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

/*
When you dodge or parry a melee attack you have a 30% chance of refreshing the cooldown on your next Avenger's Shield and causing it to generate a charge of Holy Power if used within 6 sec.
(Proc chance: 30%, 1s cooldown)
*/
func (prot *ProtectionPaladin) registerGrandCrusader() {
	hpActionID := core.ActionID{SpellID: 98057}
	prot.CanTriggerHolyAvengerHpGain(hpActionID)

	var grandCrusaderAura *core.Aura
	grandCrusaderAura = prot.RegisterAura(core.Aura{
		Label:    "Grand Crusader" + prot.Label,
		ActionID: core.ActionID{SpellID: 85416},
		Duration: time.Second * 6,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.AvengersShield.CD.Reset()
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: paladin.SpellMaskAvengersShield,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			prot.HolyPower.Gain(sim, 1, hpActionID)
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
			grandCrusaderAura.Activate(sim)
		},
	})
}
