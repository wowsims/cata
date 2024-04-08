package elemental

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/shaman"
)

func RegisterElementalShaman() {
	core.RegisterAgentFactory(
		proto.Player_ElementalShaman{},
		proto.Spec_SpecElementalShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewElementalShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ElementalShaman)
			if !ok {
				panic("Invalid spec value for Elemental Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewElementalShaman(character *core.Character, options *proto.Player) *ElementalShaman {
	eleOptions := options.GetElementalShaman().Options

	selfBuffs := shaman.SelfBuffs{
		Shield: eleOptions.ClassOptions.Shield,
	}

	totems := &proto.ShamanTotems{}
	if eleOptions.ClassOptions.Totems != nil {
		totems = eleOptions.ClassOptions.Totems
	}

	inRange := eleOptions.ThunderstormRange == proto.ElementalShaman_Options_TSInRange
	ele := &ElementalShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, inRange),
	}

	if mh := ele.GetMHWeapon(); mh != nil {
		ele.ApplyFlametongueImbueToItem(mh)
	}

	// TODO: this would never be a case?
	if oh := ele.GetOHWeapon(); oh != nil {
		ele.ApplyFlametongueImbueToItem(oh)
	}

	return ele
}

type ElementalShaman struct {
	*shaman.Shaman
}

func (eleShaman *ElementalShaman) GetShaman() *shaman.Shaman {
	return eleShaman.Shaman
}

func (eleShaman *ElementalShaman) Reset(sim *core.Simulation) {
	eleShaman.Shaman.Reset(sim)
}
