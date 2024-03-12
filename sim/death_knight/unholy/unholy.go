package unholy

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/death_knight"
)

func RegisterUnholyDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_UnholyDeathKnight{},
		proto.Spec_SpecUnholyDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewUnholyDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
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
}

func NewUnholyDeathKnight(character *core.Character, player *proto.Player) *UnholyDeathKnight {
	unholyOptions := player.GetUnholyDeathKnight().Options

	uhdk := &UnholyDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			StartingRunicPower: unholyOptions.ClassOptions.StartingRunicPower,
			PetUptime:          unholyOptions.ClassOptions.PetUptime,
			IsDps:              true,

			UseAMS:            unholyOptions.UseAms,
			AvgAMSSuccessRate: unholyOptions.AvgAmsSuccessRate,
			AvgAMSHit:         unholyOptions.AvgAmsHit,
		}, player.TalentsString),
	}

	uhdk.Inputs.UnholyFrenzyTarget = unholyOptions.UnholyFrenzyTarget

	uhdk.EnableAutoAttacks(uhdk, core.AutoAttackOptions{
		MainHand:       uhdk.WeaponFromMainHand(uhdk.DefaultMeleeCritMultiplier()),
		OffHand:        uhdk.WeaponFromOffHand(uhdk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return uhdk
}

// func (uhdk *UnholyDeathKnight) FrostPointsInBlood() int32 {
// 	return uhdk.Talents.Butchery + uhdk.Talents.Subversion + uhdk.Talents.BladeBarrier + uhdk.Talents.DarkConviction
// }

// func (uhdk *UnholyDeathKnight) FrostPointsInUnholy() int32 {
// 	return uhdk.Talents.ViciousStrikes + uhdk.Talents.Virulence + uhdk.Talents.Epidemic + uhdk.Talents.RavenousDead + uhdk.Talents.Necrosis + uhdk.Talents.BloodCakedBlade
// }

func (uhdk *UnholyDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return uhdk.DeathKnight
}

func (uhdk *UnholyDeathKnight) Initialize() {
	uhdk.DeathKnight.Initialize()
}

func (uhdk *UnholyDeathKnight) Reset(sim *core.Simulation) {
	uhdk.DeathKnight.Reset(sim)
	//uhdk.Presence = death_knight.UnsetPresence
}
