package balance

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

// T14 Balance
var ItemSetRegaliaOfTheEternalBloosom = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Eternal Blossom",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Starfall deals 20% additional damage.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  druid.DruidSpellStarfall,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the duration of your Moonfire and Sunfire spells by 2 sec.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotNumberOfTicks_Flat,
				ClassMask: druid.DruidSpellMoonfireDoT | druid.DruidSpellSunfireDoT,
				IntValue:  1,
			})
		},
	},
})

// T15 Balance
var ItemSetRegaliaOfTheHauntedForest = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Haunted Forest",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the critical strike chance of Starsurge by 10%.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  druid.DruidSpellStarsurge,
				FloatValue: 10,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Nature's Grace now also grants 1000 critical strike and 1000 mastery for its duration.
			moonkin := agent.(BalanceDruidAgent).GetBalanceDruid()
			moonkin.NaturesGrace.AttachStatBuff(stats.MasteryRating, 1000)
			moonkin.NaturesGrace.AttachStatBuff(stats.CritRating, 1000)
		},
	},
})

func init() {
}
