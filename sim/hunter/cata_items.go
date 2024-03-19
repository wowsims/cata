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
			hunter.RegisterAura(core.Aura{
				Label:    "T11 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					hunter.SerpentSting.BonusCritRating += 0.05
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					hunter.SerpentSting.BonusCritRating -= 0.05
				},
			})
		},
		4: func(agent core.Agent) {
			hunter := agent.(HunterAgent).GetHunter()
			hunter.RegisterAura(core.Aura{
				Label:    "T11 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					hunter.SteadyShot.DefaultCast.CastTime = time.Millisecond * 1800
					hunter.CobraShot.DefaultCast.CastTime = time.Millisecond * 1800
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						hunter.SteadyShot.DefaultCast.CastTime = time.Second * 2
					hunter.CobraShot.DefaultCast.CastTime = time.Second * 2
				},
			})
		},
	},
})
