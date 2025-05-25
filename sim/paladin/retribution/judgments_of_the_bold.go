package retribution

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

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
			ret.HolyPower.Gain(1, actionID, sim)

			auraArray.Get(result.Target).Activate(sim)
		},
	})
}
