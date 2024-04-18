package mage

import (
	"github.com/wowsims/cata/sim/core"
)

// T11
var ItemSetFirelordsVestments = core.NewItemSet(core.ItemSet{
	Name: "Firelord's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//Increase critical strike chance of Arcane Missiles, Ice Lance, and Pyroblast by 5%
		},
		4: func(agent core.Agent) {
			//Reduces cast time of Arcane Blast, Fireball, FFB, and Frostbolt by 10%
		},
	},
})
