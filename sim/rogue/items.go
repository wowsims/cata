package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

var Tier11 = core.NewItemSet(core.ItemSet{
	Name: "Wind Dancer's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// +5% Crit to Backstab, Mutilate, and Sinister Strike
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				FloatValue: 5 * core.CritRatingPerCritChance,
				ClassMask:  RogueSpellBackstab | RogueSpellMutilate | RogueSpellSinisterStrike,
			})
		},
		4: func(agent core.Agent) {
			// 1% Chance on Auto Attack to increase crit of next Evis or Envenom by +100% for 15 seconds
			rogue := agent.(RogueAgent).GetRogue()

			t11Proc := rogue.RegisterAura(core.Aura{
				Label:    "Deadly Scheme Proc",
				ActionID: core.ActionID{SpellID: 90472},
				Duration: time.Second * 15,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.Envenom.BonusCritRating += 100 * core.CritRatingPerCritChance
					rogue.Eviscerate.BonusCritRating += 100 * core.CritRatingPerCritChance
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.Envenom.BonusCritRating -= 100 * core.CritRatingPerCritChance
					rogue.Eviscerate.BonusCritRating -= 100 * core.CritRatingPerCritChance
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == rogue.Envenom || spell == rogue.Eviscerate {
						aura.Deactivate(sim)
					}
				},
			})

			core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
				Name:       "Deadly Scheme Aura",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.01,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					t11Proc.Activate(sim)
				},
			})
		},
	},
})

func MakeArenaSet(setName string) *core.ItemSet {
	return core.NewItemSet(core.ItemSet{
		Name: setName,

		Bonuses: map[int32]core.ApplyEffect{
			2: func(agent core.Agent) {
				agent.GetCharacter().AddStat(stats.Agility, 70)
			},
			4: func(agent core.Agent) {
				agent.GetCharacter().AddStat(stats.Agility, 90)
				// 10 maximum energy added in rogue.go
			},
		},
	})
}

var CataArenaS0 = MakeArenaSet("Bloodthirsty Gladiator's Vestments")
var CataArenaS1 = MakeArenaSet("Vicious Gladiator's Vestments")
var CataArenaS2 = MakeArenaSet("Ruthless Gladiator's Vestments")
var CataArenaS3 = MakeArenaSet("Cataclysmic Gladiator's Vestments")
