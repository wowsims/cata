package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// T11
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

			aura := warlock.RegisterAura(core.Aura{
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
					if spell.Matches(WarlockSpellFelFlame) && result.Landed() {
						aura.RemoveStack(sim)
					}
				},
			})

			core.MakePermanent(warlock.RegisterAura(core.Aura{
				Label:           "Item - Warlock T11 4P Bonus",
				ActionID:        core.ActionID{SpellID: 89935},
				ActionIDForProc: aura.ActionID,
				OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if spell.Matches(WarlockSpellImmolateDot|WarlockSpellUnstableAffliction) &&
						sim.Proc(0.02, "Warlock 4pT11") {
						aura.Activate(sim)
						aura.SetStacks(sim, 2)
					}
				},
			}))
		},
	},
})

var ItemSetGladiatorsFelshroud = core.NewItemSet(core.ItemSet{
	ID:   910,
	Name: "Gladiator's Felshroud",

	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.(WarlockAgent).GetWarlock().AddStats(stats.Stats{
				stats.Intellect: 70,
			})
		},
		4: func(agent core.Agent) {
			lock := agent.(WarlockAgent).GetWarlock()
			lock.AddStats(stats.Stats{
				stats.Intellect: 90,
			})

			// TODO: enable if we ever implement death coil
			// lock.AddStaticMod(core.SpellModConfig{
			// 	Kind:       core.SpellMod_Cooldown_Flat,
			// 	ClassMask:  WarlockSpellDeathCoil,
			// 	FloatValue: -30 * time.Second,
			// })
		},
	},
})
