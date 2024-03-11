package holy

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/priest"
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
	smiteOptions := options.GetHolyPriest()

	selfBuffs := priest.SelfBuffs{
		UseInnerFire:   smiteOptions.Options.ClassOptions.Armor == proto.PriestOptions_InnerFire,
		UseShadowfiend: smiteOptions.Options.ClassOptions.UseShadowfiend,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	holyPriest := &HolyPriest{
		Priest: basePriest,
	}

	// holyPriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	// if holyPriest.Talents.PowerInfusion && smiteOptions.Options.PowerInfusionTarget != nil {
	// 	holyPriest.SelfBuffs.PowerInfusionTarget = smiteOptions.Options.PowerInfusionTarget
	// }

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

	holyPriest.RegisterHolyFireSpell()
	holyPriest.RegisterSmiteSpell()
	holyPriest.RegisterPenanceSpell()
	holyPriest.RegisterHymnOfHopeCD()
}

func (holyPriest *HolyPriest) Reset(sim *core.Simulation) {
	holyPriest.Priest.Reset(sim)
}
