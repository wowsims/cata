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

func NewBloodDeathKnight(character *core.Character, options *proto.Player) *BloodDeathKnight {
	dkOptions := options.GetBloodDeathKnight()

	bdk := &BloodDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			IsDps:              false,
			StartingRunicPower: dkOptions.Options.StartingRunicPower,
		}, options.TalentsString),
	}

	bdk.EnableAutoAttacks(bdk, core.AutoAttackOptions{
		MainHand:       bdk.WeaponFromMainHand(bdk.DefaultMeleeCritMultiplier()),
		OffHand:        bdk.WeaponFromOffHand(bdk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if bdk.RuneStrikeQueued && bdk.RuneStrike.CanCast(sim, nil) {
				return bdk.RuneStrike
			} else {
				return mhSwingSpell
			}
		},
	})

	healingModel := options.HealingModel
	if healingModel != nil {
		if healingModel.InspirationUptime > 0.0 {
			core.ApplyInspiration(bdk.GetCharacter(), healingModel.InspirationUptime)
		}
	}

	return bdk
}

func (dk *BloodDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return dk.DeathKnight
}

func (dk *BloodDeathKnight) Initialize() {
	dk.DeathKnight.Initialize()
}

func (dk *BloodDeathKnight) Reset(sim *core.Simulation) {
	dk.DeathKnight.Reset(sim)

	dk.Presence = death_knight.UnsetPresence
	dk.DeathKnight.PseudoStats.Stunned = false
}
