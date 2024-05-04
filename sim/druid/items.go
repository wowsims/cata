package druid

import (
	"github.com/wowsims/cata/sim/core"
)

// T11 Feral
var ItemSetStormridersBattlegarb = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Battlegarb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in rake.go and lacerate.go
		},
		4: func(agent core.Agent) {
			// Implemented in mangle.go and survival_instincts.go
		},
	},
})

func init() {
}
