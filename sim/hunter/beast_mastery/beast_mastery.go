package beast_mastery

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/hunter"
)

func RegisterBeastMasteryHunter() {
	core.RegisterAgentFactory(
		proto.Player_BeastMasteryHunter{},
		proto.Spec_SpecBeastMasteryHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewBeastMasteryHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BeastMasteryHunter)
			if !ok {
				panic("Invalid spec value for Beast Mastery Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewBeastMasteryHunter(character *core.Character, options *proto.Player) *BeastMasteryHunter {
	bmHunter := &BeastMasteryHunter{
		Hunter: hunter.NewHunter(character, options, options.GetBeastMasteryHunter().Options),
	}

	return bmHunter
}

type BeastMasteryHunter struct {
	*hunter.Hunter
}

func (bmHunter *BeastMasteryHunter) GetHunter() *hunter.Hunter {
	return bmHunter.Hunter
}

func (bmHunter *BeastMasteryHunter) Reset(sim *core.Simulation) {
	bmHunter.Hunter.Reset(sim)
}
