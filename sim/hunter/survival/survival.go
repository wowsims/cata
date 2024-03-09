package survival

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/hunter"
)

func RegisterSurvivalHunter() {
	core.RegisterAgentFactory(
		proto.Player_SurvivalHunter{},
		proto.Spec_SpecSurvivalHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewSurvivalHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SurvivalHunter)
			if !ok {
				panic("Invalid spec value for Survival Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewSurvivalHunter(character *core.Character, options *proto.Player) *SurvivalHunter {
	svHunter := &SurvivalHunter{
		Hunter: hunter.NewHunter(character, options, options.GetSurvivalHunter().Options),
	}

	return svHunter
}

type SurvivalHunter struct {
	*hunter.Hunter
}

func (svHunter *SurvivalHunter) GetHunter() *hunter.Hunter {
	return svHunter.Hunter
}

func (svHunter *SurvivalHunter) Reset(sim *core.Simulation) {
	svHunter.Hunter.Reset(sim)
}
