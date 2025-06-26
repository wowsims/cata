package druid

import (
	"github.com/wowsims/mop/sim/core"
)

// T14 Balance
var ItemSetRegaliaOfTheEternalBloosom = core.NewItemSet(core.ItemSet{
	Name:                    "Regalia of the Eternal Blossom",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Your Starfall deals 20% additional damage.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  DruidSpellStarfall,
				FloatValue: 0.2,
			})
		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increases the duration of your Moonfire and Sunfire spells by 2 sec.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotNumberOfTicks_Flat,
				ClassMask: DruidSpellMoonfireDoT | DruidSpellSunfireDoT,
				IntValue:  1,
			})
		},
	},
})

// T15 Balance
var ItemSetRegaliaOfTheHauntedForest = core.NewItemSet(core.ItemSet{
	Name:                    "Regalia of the Haunted Forest",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// Increases the critical strike chance of Starsurge by 10%.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  DruidSpellStarsurge,
				FloatValue: 10,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Nature's Grace now also grants 1000 critical strike and 1000 mastery for its duration.
		},
	},
})

// T16 Balance
var ItemSetRegaliaOfTheShatteredVale = core.NewItemSet(core.ItemSet{
	ID:                      1197,
	DisabledInChallengeMode: true,
	Name:                    "Regalia of the Shattered Vale",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Arcane spells cast while in Lunar Eclipse will shoot a single Lunar Bolt at the target. Nature spells cast while in a Solar Eclipse will shoot a single Solar Bolt at the target.
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your chance to get Shooting Stars from a critical strike from Moonfire or Sunfire is increased by 8%.
		},
	},
})

func init() {
}
