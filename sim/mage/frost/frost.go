package frost

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/mage"
)

func RegisterFrostMage() {
	core.RegisterAgentFactory(
		proto.Player_FrostMage{},
		proto.Spec_SpecFrostMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFrostMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FrostMage)
			if !ok {
				panic("Invalid spec value for Frost Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFrostMage(character *core.Character, options *proto.Player) *FrostMage {
	frostOptions := options.GetFrostMage().Options

	frostMage := &FrostMage{
		Mage: mage.NewMage(character, options, frostOptions.ClassOptions),
	}
	frostMage.FrostOptions = frostOptions

	return frostMage
}

type FrostMage struct {
	*mage.Mage
}

func (frostMage *FrostMage) GetMage() *mage.Mage {
	return frostMage.Mage
}

func (frostMage *FrostMage) Reset(sim *core.Simulation) {
	frostMage.Mage.Reset(sim)
}
