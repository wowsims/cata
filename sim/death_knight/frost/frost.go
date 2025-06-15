package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/death_knight"
)

func RegisterFrostDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_FrostDeathKnight{},
		proto.Spec_SpecFrostDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFrostDeathKnight(character, options)
		},
		func(player *proto.Player, spec any) {
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
			IsDps:              true,
		}, player.TalentsString, 0),
	}

	return fdk
}

func (fdk *FrostDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return fdk.DeathKnight
}

func (fdk *FrostDeathKnight) Initialize() {
	fdk.DeathKnight.Initialize()

	fdk.registerMastery()

	fdk.registerBrittleBones()
	// fdk.registerFrostStrikeSpell()
	// fdk.registerHowlingBlastSpell()
	fdk.registerIcyTalons()
	fdk.registerImprovedFrostPresence()
	fdk.registerKillingMachine()
	fdk.registerMightOfTheFrozenWastes()
	// fdk.registerObliterateSpell()
	fdk.registerPillarOfFrost()
	fdk.registerRime()
}

func (fdk *FrostDeathKnight) ApplyTalents() {
	fdk.DeathKnight.ApplyTalents()
	fdk.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86113)

	// Blood of the North
	permanentDeathRunes := []int8{0, 1}
	fdk.SetPermanentDeathRunes(permanentDeathRunes)
	core.MakePermanent(fdk.GetOrRegisterAura(core.Aura{
		Label:    "Blood of the North" + fdk.Label,
		ActionID: core.ActionID{SpellID: 54637},
	}))
}

func (fdk *FrostDeathKnight) Reset(sim *core.Simulation) {
	fdk.DeathKnight.Reset(sim)
}
