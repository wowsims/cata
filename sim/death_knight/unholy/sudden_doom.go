package unholy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (uhdk *UnholyDeathKnight) registerSuddenDoom() {
	var suddenDoomAura *core.Aura
	suddenDoomAura = uhdk.GetOrRegisterAura(core.Aura{
		Label:     "Sudden Doom" + uhdk.Label,
		ActionID:  core.ActionID{SpellID: 81340},
		Duration:  time.Second * 10,
		MaxStacks: 0,
	}).AttachProcTrigger(core.ProcTrigger{
		Name:           "Sudden Doom Consume Trigger" + uhdk.Label,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: death_knight.DeathKnightSpellDeathCoil | death_knight.DeathKnightSpellDeathCoilHeal,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.CurCast.Cost == 0
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if uhdk.T13Dps2pc.IsActive() {
				suddenDoomAura.RemoveStack(sim)
			} else {
				suddenDoomAura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: death_knight.DeathKnightSpellDeathCoil | death_knight.DeathKnightSpellDeathCoilHeal,
		IntValue:  -100,
	})

	core.MakeProcTriggerAura(&uhdk.Unit, core.ProcTrigger{
		Name:     "Sudden Doom Trigger" + uhdk.Label,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto,
		Outcome:  core.OutcomeLanded,
		PPM:      3.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			suddenDoomAura.Activate(sim)

			// T13 2pc: Sudden Doom has a 20% chance to grant 2 charges when triggered instead of 1.
			suddenDoomAura.MaxStacks = core.TernaryInt32(uhdk.T13Dps2pc.IsActive(), 2, 0)
			if uhdk.T13Dps2pc.IsActive() {
				stacks := core.TernaryInt32(sim.Proc(0.2, "T13 2pc"), 2, 1)
				suddenDoomAura.SetStacks(sim, stacks)
			}
		},
	})
}
