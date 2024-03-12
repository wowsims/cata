package marksmanship

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/hunter"
)

func RegisterMarksmanshipHunter() {
	core.RegisterAgentFactory(
		proto.Player_MarksmanshipHunter{},
		proto.Spec_SpecMarksmanshipHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewMarksmanshipHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_MarksmanshipHunter)
			if !ok {
				panic("Invalid spec value for Marksmanship Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewMarksmanshipHunter(character *core.Character, options *proto.Player) *MarksmanshipHunter {
	mmOptions := options.GetMarksmanshipHunter().Options

	mmHunter := &MarksmanshipHunter{
		Hunter: hunter.NewHunter(character, options, mmOptions.ClassOptions),
	}
	mmHunter.MarksmanshipOptions = mmOptions

	return mmHunter
}

type MarksmanshipHunter struct {
	*hunter.Hunter
}

func (mmHunter *MarksmanshipHunter) GetHunter() *hunter.Hunter {
	return mmHunter.Hunter
}

func (mmHunter *MarksmanshipHunter) Reset(sim *core.Simulation) {
	mmHunter.Hunter.Reset(sim)
}
