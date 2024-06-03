package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
	"time"
)

// Tier 11 ret
var ItemSetReinforcedSapphiriumBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battleplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  SpellMaskTemplarsVerdict,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent) {
			// Handled in inquisition.go
		},
	},
})

// PvP sets
func (paladin *Paladin) registerPvpSets() {
	var setNames = []string{
		"Gladiator's Vindication",
		"Bloodthirsty Gladiator's Vindication",
		"Vicious Gladiator's Vindication",
		"Ruthless Gladiator's Vindication",
		"Cataclysmic Gladiator's Vindication",
	}

	pvpSetBonuses := map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStat(stats.Strength, 70)
		},
		4: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStat(stats.Strength, 90)
			paladin.AddStaticMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: SpellMaskJudgement,
				TimeValue: -1 * time.Second,
			})
		},
	}

	for _, setName := range setNames {
		core.NewItemSet(core.ItemSet{
			Name:    setName,
			Bonuses: pvpSetBonuses,
		})
	}

	switch paladin.Hands().ID {
	case 64844, 70649, 60414, 65591, 72379, 70250, 70488, 73707, 73570:
		paladin.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: 0.05,
		})
	default:
		break
	}
}

// Tier 11 prot
var ItemSetReinforcedSapphiriumBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battlearmor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  SpellMaskCrusaderStrike,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent) {
			// TODO: Handle in guardian_of_ancient_kings.go
		},
	},
})
