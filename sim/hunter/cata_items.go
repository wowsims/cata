package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

var ItemSetLightningChargedBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Lightning-Charged Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// 5% Crit on SS
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				ClassMask:  HunterSpellSerpentSting,
				FloatValue: 5 * core.CritRatingPerCritChance,
			})
		},
		4: func(agent core.Agent) {
			// Cobra & Steady Shot < 0.2s cast time
			// Cannot be spell modded for now
		},
	},
})
var ItemSetBloodthirstyGladiatorsPursuit = core.NewItemSet(core.ItemSet{
	Name: "Bloodthirsty Gladiator's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.AddStats(stats.Stats{
				stats.Agility: 70,
			})
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			// Multiply focus regen 1.05
			hunter.AddStats(stats.Stats{
				stats.Agility: 90,
			})
		},
	},
})

func (hunter *Hunter) addBloodthirstyGloves() {
	switch hunter.Hands().ID {
	case 64991, 64709, 60424, 65544, 70534, 70260, 70441, 72369, 73717, 73583:
		hunter.AddStaticMod(core.SpellModConfig{
			ClassMask: HunterSpellExplosiveTrap | HunterSpellBlackArrow,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: -time.Second * 2,
		})
	default:
		break
	}
}
