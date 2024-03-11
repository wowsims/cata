package affliction

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func RegisterAfflictionWarlock() {
	core.RegisterAgentFactory(
		proto.Player_AfflictionWarlock{},
		proto.Spec_SpecAfflictionWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewAfflictionWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_AfflictionWarlock)
			if !ok {
				panic("Invalid spec value for Affliction Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewAfflictionWarlock(character *core.Character, options *proto.Player) *AfflictionWarlock {
	affLock := &AfflictionWarlock{
		Warlock: warlock.NewWarlock(character, options, options.GetAfflictionWarlock().Options.WarlockOptions),
	}

	return affLock
}

type AfflictionWarlock struct {
	*warlock.Warlock
}

func (affLock *AfflictionWarlock) GetWarlock() *warlock.Warlock {
	return affLock.Warlock
}

func (affLock *AfflictionWarlock) Reset(sim *core.Simulation) {
	affLock.Warlock.Reset(sim)
}
