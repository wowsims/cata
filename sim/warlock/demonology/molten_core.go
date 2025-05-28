package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (demonology *DemonologyWarlock) registerMoltenCore() {
	buff := demonology.RegisterAura(core.Aura{
		Label:     "Demonic Core",
		ActionID:  core.ActionID{SpellID: 122355},
		Duration:  time.Second * 30,
		MaxStacks: 10,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -0.5,
		ClassMask:  warlock.WarlockSpellSoulFire,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -0.5,
		ClassMask:  warlock.WarlockSpellSoulFire,
	})

	// When Shadow Flame or Wild Imp deals damage 8% chance to proc
	// When Chaos Wave -> 100% Proc Chance
	apply := func(unit *core.Unit) {
		core.MakeProcTriggerAura(unit, core.ProcTrigger{
			Name:           "Demonic Core Tracker",
			Outcome:        core.OutcomeLanded,
			ClassSpellMask: warlock.WarlockSpellImpFireBolt | warlock.WarlockSpellShadowflameDot | warlock.WarlockSpellChaosWave | warlock.WarlockSpellShadowBolt | warlock.WarlockSpellSoulFire | warlock.WarlockSpellTouchOfChaos,
			Callback:       core.CallbackOnPeriodicDamageDealt | core.CallbackOnSpellHitDealt | core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Matches(warlock.WarlockSpellSoulFire) && result == nil && buff.IsActive() {
					buff.RemoveStack(sim)
				}

				if spell.Matches(warlock.WarlockSpellShadowflameDot) && sim.Proc(0.08, "Demonic Core Proc") {
					buff.Activate(sim)
					buff.AddStack(sim)
				}

				// proc fire bolt on cast
				if result == nil && spell.Matches(warlock.WarlockSpellImpFireBolt) && sim.Proc(0.08, "Demonic Core Proc") {
					buff.Activate(sim)
					buff.AddStack(sim)
				}

				if spell.Matches(warlock.WarlockSpellChaosWave) && result != nil && result.Landed() {
					buff.Activate(sim)
					buff.AddStack(sim)
				}

				// Decimation Passive effect, proc on cast
				if sim.IsExecutePhase25() && spell.Matches(warlock.WarlockSpellShadowBolt|warlock.WarlockSpellSoulFire) && result == nil {
					buff.Activate(sim)
					buff.AddStack(sim)
				}
			},
		})
	}

	apply(&demonology.Unit)
	for _, pet := range demonology.WildImps {
		apply(&pet.Unit)
	}
}
