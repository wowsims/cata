package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (destruction *DestructionWarlock) registerBackdraft() {
	buff := destruction.RegisterAura(core.Aura{
		Label:     "Backdraft",
		ActionID:  core.ActionID{SpellID: 117828},
		Duration:  time.Second * 15,
		MaxStacks: 6,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(warlock.WarlockSpellChaosBolt) && aura.GetStacks() >= 3 {
				aura.SetStacks(sim, aura.GetStacks()-3)
				return
			}

			if spell.Matches(warlock.WarlockSpellIncinerate | warlock.WarlockSpellFaBIncinerate) {
				aura.RemoveStack(sim)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -0.3,
		ClassMask:  warlock.WarlockSpellIncinerate | warlock.WarlockSpellFaBIncinerate,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -0.3,
		ClassMask:  warlock.WarlockSpellIncinerate | warlock.WarlockSpellFaBIncinerate,
	})

	// chaos bolt requries 3 charges
	mod := destruction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -0.3,
		ClassMask:  warlock.WarlockSpellChaosBolt,
	})

	buff.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
		if newStacks >= 3 {
			mod.Activate()
		} else {
			mod.Deactivate()
		}
	}

	buff.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		mod.Deactivate()
	})

	core.MakeProcTriggerAura(&destruction.Unit, core.ProcTrigger{
		Name:           "Backdraft - Trigger",
		ClassSpellMask: warlock.WarlockSpellConflagrate | warlock.WarlockSpellFaBConflagrate,
		Callback:       core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buff.Activate(sim)

			// always grants 3 stacks
			buff.SetStacks(sim, buff.GetStacks()+3)
		},
	})
}
