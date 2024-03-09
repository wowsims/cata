package blood

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/death_knight"
)

func RegisterBloodDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_BloodDeathKnight{},
		proto.Spec_SpecBloodDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewBloodDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BloodDeathKnight)
			if !ok {
				panic("Invalid spec value for Blood Death Knight!")
			}
			player.Spec = playerSpec
		},
	)
}

type BloodDeathKnight struct {
	*death_knight.DeathKnight
}

func NewBloodDeathKnight(character *core.Character, player *proto.Player) *BloodDeathKnight {
	bloodOptions := player.GetBloodDeathKnight().Options

	bdk := &BloodDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			StartingRunicPower: bloodOptions.StartingRunicPower,
			PetUptime:          bloodOptions.PetUptime,
			DrwPestiApply:      bloodOptions.DrwPestiApply,
			IsDps:              true,

			UseAMS:            bloodOptions.UseAms,
			AvgAMSSuccessRate: bloodOptions.AvgAmsSuccessRate,
			AvgAMSHit:         bloodOptions.AvgAmsHit,
		}, player.TalentsString),
	}

	bdk.EnableAutoAttacks(bdk, core.AutoAttackOptions{
		MainHand:       bdk.WeaponFromMainHand(bdk.DefaultMeleeCritMultiplier()),
		OffHand:        bdk.WeaponFromOffHand(bdk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return bdk
}

func (bdk *BloodDeathKnight) FrostPointsInBlood() int32 {
	return bdk.Talents.Butchery + bdk.Talents.Subversion + bdk.Talents.BladeBarrier + bdk.Talents.DarkConviction
}

func (bdk *BloodDeathKnight) FrostPointsInFrost() int32 {
	return bdk.Talents.ViciousStrikes + bdk.Talents.Virulence + bdk.Talents.Epidemic + bdk.Talents.RavenousDead + bdk.Talents.Necrosis + bdk.Talents.BloodCakedBlade
}

func (bdk *BloodDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return bdk.DeathKnight
}

func (bdk *BloodDeathKnight) Initialize() {
	bdk.DeathKnight.Initialize()
}

func (bdk *BloodDeathKnight) Reset(sim *core.Simulation) {
	bdk.DeathKnight.Reset(sim)
	bdk.Presence = death_knight.UnsetPresence
}
