package retribution

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

/*
Your Judgment hits grant one charge of Holy Power and cause the Physical Vulnerability effect.

Physical Vulnerability
Weakens the constitution of an enemy target, increasing their physical damage taken by 4% for 30 sec.
*/
func (ret *RetributionPaladin) registerJudgmentsOfTheBold() {
	actionID := core.ActionID{SpellID: 111528}
	ret.CanTriggerHolyAvengerHpGain(actionID)

	auraArray := ret.NewEnemyAuraArray(core.PhysVulnerabilityAura)
	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Judgments of the Bold" + ret.Label,
		ActionID:       core.ActionID{SpellID: 111529},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskJudgment,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HolyPower.Gain(sim, 1, actionID)

			auraArray.Get(result.Target).Activate(sim)
		},
	})
}
