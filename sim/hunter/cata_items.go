package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

var ItemSetLightningChargedBattleGear = core.NewItemSet(core.ItemSet{
	Name: "Lightning-Charged Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.SerpentSting.BonusCritRating += 0.05
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.SteadyShot.DefaultCast.CastTime = time.Second + (time.Millisecond + 800)
			hunter.CobraShot.DefaultCast.CastTime = time.Second + (time.Millisecond + 800)
		},
	},
})
