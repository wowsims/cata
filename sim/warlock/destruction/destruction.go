package destruction

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func RegisterDestructionWarlock() {
	core.RegisterAgentFactory(
		proto.Player_DestructionWarlock{},
		proto.Spec_SpecDestructionWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDestructionWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DestructionWarlock)
			if !ok {
				panic("Invalid spec value for Destruction Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDestructionWarlock(character *core.Character, options *proto.Player) *DestructionWarlock {
	destroOptions := options.GetDestructionWarlock().Options
	destroLock := &DestructionWarlock{
		Warlock: warlock.NewWarlock(character, options, destroOptions.ClassOptions),
	}

	return destroLock
}

type DestructionWarlock struct {
	*warlock.Warlock
}

func (destroLock *DestructionWarlock) GetWarlock() *warlock.Warlock {
	return destroLock.Warlock
}

func (destroLock *DestructionWarlock) Reset(sim *core.Simulation) {
	destroLock.Warlock.Reset(sim)
}
