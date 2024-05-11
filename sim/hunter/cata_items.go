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
		},
		4: func(agent core.Agent) {
			// Cobra & Steady Shot < 0.2s cast time
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
	case 64709: //Todo: Add more ids here when needed
		hunter.AddStaticMod(core.SpellModConfig{
			ClassMask: HunterSpellExplosiveTrap,
			Kind:      core.SpellMod_Cooldown_Flat,
			TimeValue: -time.Second * 2,
		})
	default:
		break
	}
}
