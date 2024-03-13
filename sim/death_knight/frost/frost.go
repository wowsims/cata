package frost

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/death_knight"
)

func RegisterFrostDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_FrostDeathKnight{},
		proto.Spec_SpecFrostDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFrostDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FrostDeathKnight)
			if !ok {
				panic("Invalid spec value for Frost Death Knight!")
			}
			player.Spec = playerSpec
		},
	)
}

type FrostDeathKnight struct {
	*death_knight.DeathKnight
}

func NewFrostDeathKnight(character *core.Character, player *proto.Player) *FrostDeathKnight {
	frostOptions := player.GetFrostDeathKnight().Options

	fdk := &FrostDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			StartingRunicPower: frostOptions.ClassOptions.StartingRunicPower,
			PetUptime:          frostOptions.ClassOptions.PetUptime,
			IsDps:              true,

			UseAMS:            frostOptions.UseAms,
			AvgAMSSuccessRate: frostOptions.AvgAmsSuccessRate,
			AvgAMSHit:         frostOptions.AvgAmsHit,
		}, player.TalentsString),
	}

	fdk.EnableAutoAttacks(fdk, core.AutoAttackOptions{
		MainHand:       fdk.WeaponFromMainHand(fdk.DefaultMeleeCritMultiplier()),
		OffHand:        fdk.WeaponFromOffHand(fdk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return fdk
}

// func (fdk *FrostDeathKnight) FrostPointsInBlood() int32 {
// 	return fdk.Talents.Butchery + fdk.Talents.Subversion + fdk.Talents.BladeBarrier + fdk.Talents.DarkConviction
// }

// func (fdk *FrostDeathKnight) FrostPointsInFrost() int32 {
// 	return fdk.Talents.ViciousStrikes + fdk.Talents.Virulence + fdk.Talents.Epidemic + fdk.Talents.RavenousDead + fdk.Talents.Necrosis + fdk.Talents.BloodCakedBlade
// }

func (fdk *FrostDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return fdk.DeathKnight
}

func (fdk *FrostDeathKnight) Initialize() {
	fdk.DeathKnight.Initialize()
}

func (fdk *FrostDeathKnight) Reset(sim *core.Simulation) {
	fdk.DeathKnight.Reset(sim)
	//fdk.Presence = death_knight.UnsetPresence
}
