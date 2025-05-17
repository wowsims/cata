package druid

import (
	"github.com/wowsims/mop/sim/core"
)

// T14 Balance
var ItemSetRegaliaOfTheEternalBloosom = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Eternal Blossom",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Starfall deals 20% additional damage.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  DruidSpellStarfall,
				FloatValue: 10,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the duration of your Moonfire and Sunfire spells by 2 sec.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotNumberOfTicks_Flat,
				ClassMask: DruidSpellMoonfireDoT | DruidSpellSunfireDoT,
				IntValue:  1,
			})
		},
	},
})

func init() {
}
