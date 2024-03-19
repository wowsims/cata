package hunter

import (
	"github.com/wowsims/cata/sim/core"
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
