package blood

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
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

	vengeance *core.VengeanceTracker
}

func NewBloodDeathKnight(character *core.Character, options *proto.Player) *BloodDeathKnight {
	dkOptions := options.GetBloodDeathKnight()

	bdk := &BloodDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			IsDps:              false,
			StartingRunicPower: dkOptions.Options.ClassOptions.StartingRunicPower,
			Spec:               proto.Spec_SpecBloodDeathKnight,
		}, options.TalentsString, 50034),
		vengeance: &core.VengeanceTracker{},
	}

	bdk.EnableAutoAttacks(bdk, core.AutoAttackOptions{
		MainHand:       bdk.WeaponFromMainHand(bdk.DefaultMeleeCritMultiplier()),
		OffHand:        bdk.WeaponFromOffHand(bdk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(&bdk.Unit, healingModel.InspirationUptime)
		}
	}

	return bdk
}

func (bdk *BloodDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return bdk.DeathKnight
}

func (bdk *BloodDeathKnight) Initialize() {
	bdk.DeathKnight.Initialize()

	bdk.registerHeartStrikeSpell()
}

func (bdk BloodDeathKnight) getMasteryBonus() float64 {
	return 0.5 + 0.0625*bdk.GetMasteryPoints()
}

func (bdk *BloodDeathKnight) ApplyTalents() {
	bdk.DeathKnight.ApplyTalents()

	// Veteran of the Third War
	bdk.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: death_knight.DeathKnightSpellOutbreak,
		TimeValue: time.Second * -30,
	})
	bdk.MultiplyStat(stats.Stamina, 1.09)
	bdk.AddStat(stats.Expertise, 6*core.ExpertisePerQuarterPercentReduction)
	core.MakePermanent(bdk.GetOrRegisterAura(core.Aura{
		Label:    "Veteran of the Third War",
		ActionID: core.ActionID{SpellID: 50029},
	}))

	// Vengeance
	core.ApplyVengeanceEffect(&bdk.Character, bdk.vengeance, 93099)

	// Mastery: Blood Shield
	core.MakePermanent(bdk.GetOrRegisterAura(core.Aura{
		Label:    "Blood Shield",
		ActionID: core.ActionID{SpellID: 77513},
	}))

}

func (bdk *BloodDeathKnight) Reset(sim *core.Simulation) {
	bdk.DeathKnight.Reset(sim)
}
