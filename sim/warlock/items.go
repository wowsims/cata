package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// T11
// TODO: untested, since it's currently not working ingame
var ItemSetMaleficRaiment = core.NewItemSet(core.ItemSet{
	Name: "Shadowflame Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(WarlockAgent).GetWarlock().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_CastTime_Pct,
				ClassMask:  WarlockSpellChaosBolt | WarlockSpellHandOfGuldan | WarlockSpellHaunt,
				FloatValue: -0.1,
			})
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			dmgMod := warlock.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  WarlockSpellFelFlame,
				FloatValue: 3.0,
			})

			aura := warlock.GetOrRegisterAura(core.Aura{
				Label:     "Fel Spark",
				ActionID:  core.ActionID{SpellID: 89937},
				Duration:  15 * time.Second,
				MaxStacks: 2,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Deactivate()
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == warlock.FelFlame && result.Landed() {
						aura.RemoveStack(sim)
					}
				},
			})

			core.MakePermanent(warlock.RegisterAura(core.Aura{
				Label:           "Item - Warlock T11 4P Bonus",
				ActionID:        core.ActionID{SpellID: 89935},
				ActionIDForProc: aura.ActionID,
				OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == warlock.ImmolateDot && sim.Proc(0.02, "Warlock 4pT11") {
						aura.Activate(sim)
						aura.SetStacks(sim, 2)
					}
				},
			}))
		},
	},
})
