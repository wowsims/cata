package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// T11 Feral
var ItemSetStormridersBattlegarb = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Battlegarb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in rake.go and lacerate.go
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			var apDepByStackCount = map[int32]*stats.StatDependency{}

			for i := 1; i <= 3; i++ {
				apDepByStackCount[int32(i)] = druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.01*float64(i))
			}

			druid.StrengthOfThePantherAura = druid.RegisterAura(core.Aura{
				Label:     "Strength of the Panther",
				ActionID:  core.ActionID{SpellID: 90166},
				Duration:  time.Second * 30,
				MaxStacks: 3,

				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					if oldStacks > 0 {
						druid.DisableDynamicStatDep(sim, apDepByStackCount[oldStacks])
					}

					if newStacks > 0 {
						druid.EnableDynamicStatDep(sim, apDepByStackCount[newStacks])
					}
				},
			})
		},
	},
})

// T11 Balance
var ItemSetStormridersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				FloatValue: 5 * core.CritRatingPerCritChance,
				ClassMask:  DruidSpellDoT | DruidSpellMoonfire | DruidSpellSunfire,
			})
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			tierSet4pMod := druid.AddDynamicMod(core.SpellModConfig{
				School: core.SpellSchoolArcane | core.SpellSchoolNature,
				Kind:   core.SpellMod_BonusCrit_Rating,
			})

			tierSet4pAura := druid.RegisterAura(core.Aura{
				ActionID:  core.ActionID{SpellID: 90163},
				Label:     "Druid T11 Balance 4P Bonus",
				Duration:  time.Second * 8,
				MaxStacks: 3,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.SetStacks(sim, aura.MaxStacks)

					tierSet4pMod.UpdateFloatValue(float64(aura.GetStacks()) * 5 * core.CritRatingPerCritChance)
					tierSet4pMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					tierSet4pMod.Deactivate()
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidCrit() && aura.GetStacks() > 0 {
						aura.RemoveStack(sim)
						tierSet4pMod.UpdateFloatValue(float64(aura.GetStacks()) * 5 * core.CritRatingPerCritChance)
					}
				},
			})

			druid.AddEclipseCallback(func(_ Eclipse, gained bool, sim *core.Simulation) {
				if gained {
					tierSet4pAura.Activate(sim)
				} else {
					tierSet4pAura.Deactivate(sim)
				}
			})
		},
	},
})

func init() {
}
