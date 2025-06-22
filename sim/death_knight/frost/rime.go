package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Your Obliterate has a 45% chance to cause your next Howling Blast or Icy Touch to consume no runes.
(Proc chance: 45%)
*/
func (fdk *FrostDeathKnight) registerRime() {
	var freezingFogAura *core.Aura
	freezingFogAura = fdk.GetOrRegisterAura(core.Aura{
		Label:     "Freezing Fog" + fdk.Label,
		ActionID:  core.ActionID{SpellID: 59052},
		Duration:  time.Second * 15,
		MaxStacks: 0,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: death_knight.DeathKnightSpellIcyTouch | death_knight.DeathKnightSpellHowlingBlast,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.CurCast.Cost <= 0
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if fdk.T13Dps2pc.IsActive() {
				freezingFogAura.RemoveStack(sim)
			} else {
				freezingFogAura.Deactivate(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  death_knight.DeathKnightSpellIcyTouch | death_knight.DeathKnightSpellHowlingBlast,
		FloatValue: -2.0,
	})

	core.MakeProcTriggerAura(&fdk.Unit, core.ProcTrigger{
		Name:           "Rime" + fdk.Label,
		ActionID:       core.ActionID{SpellID: 59057},
		Callback:       core.CallbackOnSpellHitDealt,
		ProcMask:       core.ProcMaskMeleeMH,
		ClassSpellMask: death_knight.DeathKnightSpellObliterate,
		Outcome:        core.OutcomeLanded,
		ProcChance:     0.45,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			freezingFogAura.Activate(sim)

			// T13 2pc: Rime has a 60% chance to grant 2 charges when triggered instead of 1.
			freezingFogAura.MaxStacks = core.TernaryInt32(fdk.T13Dps2pc.IsActive(), 2, 0)
			if fdk.T13Dps2pc.IsActive() {
				stacks := core.TernaryInt32(sim.Proc(0.6, "T13 2pc"), 2, 1)
				freezingFogAura.SetStacks(sim, stacks)
			}
		},
	})
}
