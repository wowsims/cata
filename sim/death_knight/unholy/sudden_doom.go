package unholy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Reduces the cost of Death Coil by 20%.

While in Unholy Presence, grants your main-hand autoattacks a chance to make your next Death Coil cost no Runic Power.
*/
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
			return spell.CurCast.Cost <= 0
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if uhdk.T13Dps2pc.IsActive() {
				suddenDoomAura.RemoveStack(sim)
			} else {
				suddenDoomAura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  death_knight.DeathKnightSpellDeathCoil | death_knight.DeathKnightSpellDeathCoilHeal,
		FloatValue: -2.0,
	})

	// Dummy spell to react with triggers
	sdProcSpell := uhdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 81340},
		Flags:          core.SpellFlagNoLogs | core.SpellFlagNoMetrics,
		ClassSpellMask: death_knight.DeathKnightSpellSuddenDoom,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)

			// T13 2pc: Sudden Doom has a 20% chance to grant 2 charges when triggered instead of 1.
			spell.RelatedSelfBuff.MaxStacks = core.TernaryInt32(uhdk.T13Dps2pc.IsActive(), 2, 0)
			if uhdk.T13Dps2pc.IsActive() {
				stacks := core.TernaryInt32(sim.Proc(0.2, "T13 2pc"), 2, 1)
				spell.RelatedSelfBuff.SetStacks(sim, stacks)
			}
		},

		RelatedSelfBuff: suddenDoomAura,
	})

	core.MakeProcTriggerAura(&uhdk.Unit, core.ProcTrigger{
		Name:     "Sudden Doom Trigger" + uhdk.Label,
		ActionID: core.ActionID{SpellID: 49530},
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto,
		Outcome:  core.OutcomeLanded,
		DPM:      uhdk.NewStaticLegacyPPMManager(3.0, core.ProcMaskMeleeMHAuto),

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return uhdk.UnholyPresenceSpell.RelatedSelfBuff.IsActive()
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			sdProcSpell.Cast(sim, &uhdk.Unit)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  death_knight.DeathKnightSpellDeathCoil | death_knight.DeathKnightSpellDeathCoilHeal,
		FloatValue: -0.2,
	})
}
