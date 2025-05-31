package protection

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

// Your Judgment hits grant one charge of Holy Power.
func (prot *ProtectionPaladin) registerJudgmentsOfTheWise() {
	prot.JudgmentsOfTheWiseActionID = core.ActionID{SpellID: 105427}
	prot.CanTriggerHolyAvengerHpGain(prot.JudgmentsOfTheWiseActionID)

	core.MakeProcTriggerAura(&prot.Unit, core.ProcTrigger{
		Name:           "Judgments of the Wise" + prot.Label,
		ActionID:       core.ActionID{SpellID: 105424},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskJudgment,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			prot.HolyPower.Gain(sim, 1, prot.JudgmentsOfTheWiseActionID)
		},
	})
}
