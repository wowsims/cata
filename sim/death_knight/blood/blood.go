package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/death_knight"
)

func RegisterBloodDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_BloodDeathKnight{},
		proto.Spec_SpecBloodDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewBloodDeathKnight(character, options)
		},
		func(player *proto.Player, spec any) {
			playerSpec, ok := spec.(*proto.Player_BloodDeathKnight)
			if !ok {
				panic("Invalid spec value for Blood Death Knight!")
			}
			player.Spec = playerSpec
		},
	)
}

// Threat Done By Caster setup
const (
	TDBC_DarkCommand int = iota

	TDBC_Total
)

type BloodDeathKnight struct {
	*death_knight.DeathKnight
}

func NewBloodDeathKnight(character *core.Character, options *proto.Player) *BloodDeathKnight {
	dkOptions := options.GetBloodDeathKnight()

	bdk := &BloodDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			IsDps:              false,
			StartingRunicPower: dkOptions.Options.ClassOptions.StartingRunicPower,
			Spec:               proto.Spec_SpecBloodDeathKnight,
		}, options.TalentsString, 50034),
	}

	bdk.RuneWeapon = bdk.NewRuneWeapon()

	bdk.Bloodworm = make([]*death_knight.BloodwormPet, 5)
	for i := range 5 {
		bdk.Bloodworm[i] = bdk.NewBloodwormPet(i)
	}

	return bdk
}

func (bdk *BloodDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return bdk.DeathKnight
}

func (bdk *BloodDeathKnight) Initialize() {
	bdk.DeathKnight.Initialize()

	bdk.registerMastery()

	bdk.registerBloodParasite()
	bdk.registerBloodRites()
	bdk.registerBoneShield()
	bdk.registerCrimsonScourge()
	bdk.registerDancingRuneWeapon()
	bdk.registerDarkCommand()
	bdk.registerHeartStrike()
	bdk.registerHotfixPassive()
	bdk.registerImprovedBloodPresence()
	bdk.registerRiposte()
	bdk.registerRuneStrike()
	bdk.registerRuneTap()
	bdk.registerSanguineFortitude()
	bdk.registerScarletFever()
	bdk.registerScentOfBlood()
	bdk.registerVampiricBlood()
	bdk.registerVeteranOfTheThirdWar()
	bdk.registerWillOfTheNecropolis()

	bdk.RuneWeapon.AddCopySpell(HeartStrikeActionID, bdk.registerDrwHeartStrike())
	bdk.RuneWeapon.AddCopySpell(RuneStrikeActionID, bdk.registerDrwRuneStrike())
}

func (bdk *BloodDeathKnight) ApplyTalents() {
	bdk.DeathKnight.ApplyTalents()
	bdk.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypePlate, 86537)

	// Vengeance
	bdk.RegisterVengeance(93099, nil)
}

func (bdk *BloodDeathKnight) Reset(sim *core.Simulation) {
	bdk.DeathKnight.Reset(sim)
}
