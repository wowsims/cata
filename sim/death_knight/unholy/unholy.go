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

	lastScourgeStrikeDamage float64
}

func NewUnholyDeathKnight(character *core.Character, player *proto.Player) *UnholyDeathKnight {
	unholyOptions := player.GetUnholyDeathKnight().Options

	uhdk := &UnholyDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			Spec: proto.Spec_SpecUnholyDeathKnight,

			StartingRunicPower: unholyOptions.ClassOptions.StartingRunicPower,
			PetUptime:          unholyOptions.ClassOptions.PetUptime,
			IsDps:              true,

			UseAMS:            unholyOptions.UseAms,
			AvgAMSSuccessRate: unholyOptions.AvgAmsSuccessRate,
			AvgAMSHit:         unholyOptions.AvgAmsHit,
		}, player.TalentsString, 56835),
	}

	uhdk.Inputs.UnholyFrenzyTarget = unholyOptions.UnholyFrenzyTarget

	return uhdk
}

func (uhdk UnholyDeathKnight) getMasteryShadowBonus() float64 {
	return 0.2 + 0.025*uhdk.GetMasteryPoints()
}

func (uhdk *UnholyDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return uhdk.DeathKnight
}

func (uhdk *UnholyDeathKnight) Initialize() {
	uhdk.DeathKnight.Initialize()

	// uhdk.registerScourgeStrikeSpell()
}

func (uhdk *UnholyDeathKnight) ApplyTalents() {
	uhdk.DeathKnight.ApplyTalents()
	uhdk.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86524)

	// Mastery: Dreadblade
	masteryMod := uhdk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		School:    core.SpellSchoolShadow,
		ClassMask: death_knight.DeathKnightSpellScourgeStrikeShadow | death_knight.DeathKnightSpellUnholyBlight,
	})

	uhdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		uhdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *=
			(1.2 + 0.025*core.MasteryRatingToMasteryPoints(newMastery)) / (1.2 + 0.025*core.MasteryRatingToMasteryPoints(oldMastery))
		masteryMod.UpdateFloatValue(uhdk.getMasteryShadowBonus())
	})

	core.MakePermanent(uhdk.GetOrRegisterAura(core.Aura{
		Label:    "Dreadblade",
		ActionID: core.ActionID{SpellID: 77515},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			uhdk.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.2 + 0.025*uhdk.GetMasteryPoints()
			masteryMod.UpdateFloatValue(uhdk.getMasteryShadowBonus())
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

	// Unholy Might
	uhdk.MultiplyStat(stats.Strength, 1.25)
	core.MakePermanent(uhdk.GetOrRegisterAura(core.Aura{
		Label:    "Unholy Might",
		ActionID: core.ActionID{SpellID: 91107},
	}))
}

func (uhdk *UnholyDeathKnight) Reset(sim *core.Simulation) {
	uhdk.DeathKnight.Reset(sim)
}
