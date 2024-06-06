package paladin

import (
	"github.com/wowsims/cata/sim/core"
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

func (paladin *Paladin) addBloodthirstyGloves() {
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
