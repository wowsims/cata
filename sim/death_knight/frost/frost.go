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
			UseAMS:             false,
		}, player.TalentsString, 0),
	}

	return fdk
}

func (fdk *FrostDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return fdk.DeathKnight
}

func (fdk FrostDeathKnight) getMasteryFrostBonus() float64 {
	return 0.16 + 0.02*fdk.GetMasteryPoints()
}

func (fdk *FrostDeathKnight) Initialize() {
	fdk.DeathKnight.Initialize()

	// fdk.registerFrostStrikeSpell()
}

func (fdk *FrostDeathKnight) ApplyTalents() {
	fdk.DeathKnight.ApplyTalents()
	fdk.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86524)

	masteryMod := fdk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: fdk.getMasteryFrostBonus(),
		School:     core.SpellSchoolFrost,
	})

	fdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		masteryMod.UpdateFloatValue(fdk.getMasteryFrostBonus())
	})

	core.MakePermanent(fdk.GetOrRegisterAura(core.Aura{
		Label:    "Frozen Heart",
		ActionID: core.ActionID{SpellID: 77514},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

	// Icy Talons
	fdk.PseudoStats.MeleeSpeedMultiplier *= 1.2
	core.MakePermanent(fdk.GetOrRegisterAura(core.Aura{
		Label:    "Icy Talons" + fdk.Label,
		ActionID: core.ActionID{SpellID: 50887},
	}))

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
