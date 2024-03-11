package fire

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/mage"
)

func RegisterFireMage() {
	core.RegisterAgentFactory(
		proto.Player_FireMage{},
		proto.Spec_SpecFireMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFireMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FireMage)
			if !ok {
				panic("Invalid spec value for Fire Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFireMage(character *core.Character, options *proto.Player) *FireMage {
	fireMage := &FireMage{
		Mage: mage.NewMage(character, options, options.GetFireMage().Options.MageOptions),
	}

	return fireMage
}

type FireMage struct {
	*mage.Mage
}

func (fireMage *FireMage) GetMage() *mage.Mage {
	return fireMage.Mage
}

func (fireMage *FireMage) Reset(sim *core.Simulation) {
	fireMage.Mage.Reset(sim)
}
