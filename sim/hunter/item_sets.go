package hunter

import (
	"github.com/wowsims/mop/sim/core"
)

var YaunGolSlayersBattlegear = core.NewItemSet(core.ItemSet{
	Name:                    "Yaungol Slayer Battlegear",
	ID:                      1129,
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  HunterSpellExplosiveShot,
				FloatValue: 0.05,
			})
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:              core.SpellMod_DamageDone_Pct,
				ClassMask:         HunterSpellKillCommand,
				ShouldApplyToPets: true,
				FloatValue:        0.15,
			})

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  HunterSpellChimeraShot,
				FloatValue: 0.15,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Stub
		},
	},
})

var SaurokStalker = core.NewItemSet(core.ItemSet{
	Name:                    "Battlegear of the Saurok Stalker",
	ID:                      1157,
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Summon Thunderhawk
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			//
		},
	},
})
