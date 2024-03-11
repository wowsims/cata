package demonology

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func RegisterDemonologyWarlock() {
	core.RegisterAgentFactory(
		proto.Player_DemonologyWarlock{},
		proto.Spec_SpecDemonologyWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDemonologyWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DemonologyWarlock)
			if !ok {
				panic("Invalid spec value for Demonology Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDemonologyWarlock(character *core.Character, options *proto.Player) *DemonologyWarlock {
	demoLock := &DemonologyWarlock{
		Warlock: warlock.NewWarlock(character, options, options.GetDemonologyWarlock().Options.WarlockOptions),
	}

	return demoLock
}

type DemonologyWarlock struct {
	*warlock.Warlock
}

func (demoLock *DemonologyWarlock) GetWarlock() *warlock.Warlock {
	return demoLock.Warlock
}

func (demoLock *DemonologyWarlock) Reset(sim *core.Simulation) {
	demoLock.Warlock.Reset(sim)
}
