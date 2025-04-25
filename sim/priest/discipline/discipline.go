package discipline

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/priest"
)

func RegisterDisciplinePriest() {
	core.RegisterAgentFactory(
		proto.Player_DisciplinePriest{},
		proto.Spec_SpecDisciplinePriest,
		func(character *core.Character, options *proto.Player) core.Agent {
			return newDisciplinePriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DisciplinePriest)
			if !ok {
				panic("Invalid spec value for Discipline Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

type DisciplinePriest struct {
	*priest.Priest

	Options *proto.DisciplinePriest_Options
}

func newDisciplinePriest(character *core.Character, options *proto.Player) *DisciplinePriest {
	discOptions := options.GetDisciplinePriest()

	selfBuffs := priest.SelfBuffs{
		UseInnerFire:   discOptions.Options.ClassOptions.Armor == proto.PriestOptions_InnerFire,
		UseShadowfiend: discOptions.Options.ClassOptions.UseShadowfiend,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	discPriest := &DisciplinePriest{
		Priest:  basePriest,
		Options: discOptions.Options,
	}

	discPriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	// TODO: Fix this to work with the new talent system.
	// if discPriest.Talents.PowerInfusion && discPriest.Options.PowerInfusionTarget != nil {
	// 	discPriest.SelfBuffs.PowerInfusionTarget = discPriest.Options.PowerInfusionTarget
	// }

	return discPriest
}

func (discPriest *DisciplinePriest) GetPriest() *priest.Priest {
	return discPriest.Priest
}

func (discPriest *DisciplinePriest) GetMainTarget() *core.Unit {
	target := discPriest.Env.Raid.GetFirstTargetDummy()
	if target == nil {
		return &discPriest.Unit
	} else {
		return &target.Unit
	}
}

func (discPriest *DisciplinePriest) Initialize() {
	discPriest.CurrentTarget = discPriest.GetMainTarget()
	discPriest.Priest.Initialize()
	// discPriest.Priest.RegisterHealingSpells()

	// // discPriest.ApplyRapture(discPriest.Options.RapturesPerMinute)
	// discPriest.RegisterHymnOfHopeCD()
}

func (holyPriest *DisciplinePriest) ApplyTalents() {
}

func (discPriest *DisciplinePriest) Reset(sim *core.Simulation) {
	//discPriest.Priest.Reset(sim)
}
