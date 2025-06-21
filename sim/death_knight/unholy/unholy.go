package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/death_knight"
)

func RegisterUnholyDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_UnholyDeathKnight{},
		proto.Spec_SpecUnholyDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewUnholyDeathKnight(character, options)
		},
		func(player *proto.Player, spec any) {
			playerSpec, ok := spec.(*proto.Player_UnholyDeathKnight)
			if !ok {
				panic("Invalid spec value for Unholy Death Knight!")
			}
			player.Spec = playerSpec
		},
	)
}

type UnholyDeathKnight struct {
	*death_knight.DeathKnight

	lastScourgeStrikeDamage float64
}

func NewUnholyDeathKnight(character *core.Character, player *proto.Player) *UnholyDeathKnight {
	unholyOptions := player.GetUnholyDeathKnight().Options

	uhdk := &UnholyDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			Spec: proto.Spec_SpecUnholyDeathKnight,

			StartingRunicPower: unholyOptions.ClassOptions.StartingRunicPower,
			IsDps:              true,
		}, player.TalentsString, 56835),
	}

	uhdk.Gargoyle = uhdk.NewGargoyle()
	uhdk.Inputs.UnholyFrenzyTarget = unholyOptions.UnholyFrenzyTarget

	return uhdk
}

func (uhdk *UnholyDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return uhdk.DeathKnight
}

func (uhdk *UnholyDeathKnight) Initialize() {
	uhdk.DeathKnight.Initialize()

	uhdk.registerMastery()

	uhdk.registerBloodStrike()
	uhdk.registerDarkTransformation()
	uhdk.registerEbonPlaguebringer()
	uhdk.registerFesteringStrike()
	uhdk.registerImprovedUnholyPresence()
	uhdk.registerMasterOfGhouls()
	uhdk.registerScourgeStrike()
	uhdk.registerReaping()
	uhdk.registerShadowInfusion()
	uhdk.registerSuddenDoom()
	uhdk.registerSummonGargoyle()
	uhdk.registerUnholyFrenzy()
	uhdk.registerUnholyMight()
}

func (uhdk *UnholyDeathKnight) ApplyTalents() {
	uhdk.DeathKnight.ApplyTalents()
	uhdk.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86536)
}

func (uhdk *UnholyDeathKnight) Reset(sim *core.Simulation) {
	uhdk.DeathKnight.Reset(sim)
}
