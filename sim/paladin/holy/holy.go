package holy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

func RegisterHolyPaladin() {
	core.RegisterAgentFactory(
		proto.Player_HolyPaladin{},
		proto.Spec_SpecHolyPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewHolyPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_HolyPaladin) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Holy Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewHolyPaladin(character *core.Character, options *proto.Player) *HolyPaladin {
	holyOptions := options.GetHolyPaladin()

	holy := &HolyPaladin{
		Paladin: paladin.NewPaladin(character, options.TalentsString, holyOptions.Options.ClassOptions),
	}

	return holy
}

type HolyPaladin struct {
	*paladin.Paladin
}

func (holy *HolyPaladin) GetPaladin() *paladin.Paladin {
	return holy.Paladin
}

func (holy *HolyPaladin) ApplyTalents() {
	holy.Paladin.ApplyTalents()
	holy.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypePlate, 86525)
}

func (holy *HolyPaladin) Initialize() {
	holy.Paladin.Initialize()

	holy.registerHotfixPassive()
}

func (holy *HolyPaladin) Reset(sim *core.Simulation) {
	holy.Paladin.Reset(sim)
}
