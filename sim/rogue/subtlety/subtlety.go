package subtlety

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func RegisterSubtletyRogue() {
	core.RegisterAgentFactory(
		proto.Player_SubtletyRogue{},
		proto.Spec_SpecSubtletyRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewSubtletyRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SubtletyRogue)
			if !ok {
				panic("Invalid spec value for Subtlety Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewSubtletyRogue(character *core.Character, options *proto.Player) *SubtletyRogue {
	subOptions := options.GetSubtletyRogue().Options

	subRogue := &SubtletyRogue{
		Rogue: rogue.NewRogue(character, options.TalentsString),
	}
	subRogue.SubtletyOptions = subOptions

	return subRogue
}

type SubtletyRogue struct {
	*rogue.Rogue
}

func (subRogue *SubtletyRogue) GetRogue() *rogue.Rogue {
	return subRogue.Rogue
}

func (subRogue *SubtletyRogue) Reset(sim *core.Simulation) {
	subRogue.Rogue.Reset(sim)
}
