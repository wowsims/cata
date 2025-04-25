package holy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/priest"
)

func RegisterHolyPriest() {
	core.RegisterAgentFactory(
		proto.Player_HolyPriest{},
		proto.Spec_SpecHolyPriest,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewHolyPriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_HolyPriest)
			if !ok {
				panic("Invalid spec value for Holy Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewHolyPriest(character *core.Character, options *proto.Player) *HolyPriest {
	holyOptions := options.GetHolyPriest()

	selfBuffs := priest.SelfBuffs{
		UseInnerFire:   holyOptions.Options.ClassOptions.Armor == proto.PriestOptions_InnerFire,
		UseShadowfiend: holyOptions.Options.ClassOptions.UseShadowfiend,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	holyPriest := &HolyPriest{
		Priest: basePriest,
	}

	return holyPriest
}

type HolyPriest struct {
	*priest.Priest
}

func (holyPriest *HolyPriest) GetPriest() *priest.Priest {
	return holyPriest.Priest
}

func (holyPriest *HolyPriest) Initialize() {
	holyPriest.Priest.Initialize()

	// holyPriest.RegisterHolyFireSpell()
	// holyPriest.RegisterSmiteSpell()
	// holyPriest.RegisterPenanceSpell()
	// holyPriest.RegisterHymnOfHopeCD()
}

func (holyPriest *HolyPriest) ApplyTalents() {
}

func (holyPriest *HolyPriest) Reset(sim *core.Simulation) {
	//holyPriest.Priest.Reset(sim)
}
