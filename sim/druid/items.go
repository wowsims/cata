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
				apDepByStackCount[int32(i)] = druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0 + 0.01*float64(i))
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

func init() {
}
