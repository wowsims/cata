package tank

import (
	"github.com/wowsims/cata/sim/DeathKnight"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func RegisterTankDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_TankDeathKnight{},
		proto.Spec_SpecTankDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewTankDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_TankDeathKnight)
			if !ok {
				panic("Invalid spec value for Tank DeathKnight!")
			}
			player.Spec = playerSpec
		},
	)
}

type TankDeathKnight struct {
	*DeathKnight.DeathKnight
}

func NewTankDeathKnight(character *core.Character, options *proto.Player) *TankDeathKnight {
	dkOptions := options.GetTankDeathKnight()

	tankDk := &TankDeathKnight{
		DeathKnight: DeathKnight.NewDeathKnight(character, DeathKnight.DeathKnightInputs{
			IsDps:              false,
			StartingRunicPower: dkOptions.Options.StartingRunicPower,
		}, options.TalentsString),
	}

	tankDk.Inputs.UnholyFrenzyTarget = dkOptions.Options.UnholyFrenzyTarget

	tankDk.EnableAutoAttacks(tankDk, core.AutoAttackOptions{
		MainHand:       tankDk.WeaponFromMainHand(tankDk.DefaultMeleeCritMultiplier()),
		OffHand:        tankDk.WeaponFromOffHand(tankDk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if tankDk.RuneStrikeQueued && tankDk.RuneStrike.CanCast(sim, nil) {
				return tankDk.RuneStrike
			} else {
				return mhSwingSpell
			}
		},
	})

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(tankDk.GetCharacter(), healingModel.InspirationUptime)
		}
	}

	return tankDk
}

func (dk *TankDeathKnight) GetDeathKnight() *DeathKnight.DeathKnight {
	return dk.DeathKnight
}

func (dk *TankDeathKnight) Initialize() {
	dk.DeathKnight.Initialize()
}

func (dk *TankDeathKnight) Reset(sim *core.Simulation) {
	dk.DeathKnight.Reset(sim)

	dk.Presence = DeathKnight.UnsetPresence
	dk.DeathKnight.PseudoStats.Stunned = false
}
