package arcane

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/mage"
)

func RegisterArcaneMage() {
	core.RegisterAgentFactory(
		proto.Player_ArcaneMage{},
		proto.Spec_SpecArcaneMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewArcaneMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ArcaneMage)
			if !ok {
				panic("Invalid spec value for Arcane Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewArcaneMage(character *core.Character, options *proto.Player) *ArcaneMage {
	arcaneMage := &ArcaneMage{
		Mage: mage.NewMage(character, options, options.GetArcaneMage().Options.MageOptions),
	}

	return arcaneMage
}

type ArcaneMage struct {
	*mage.Mage
}

func (arcaneMage *ArcaneMage) GetMage() *mage.Mage {
	return arcaneMage.Mage
}

func (arcaneMage *ArcaneMage) Reset(sim *core.Simulation) {
	arcaneMage.Mage.Reset(sim)
}
