package arms

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func RegisterArmsWarrior() {
	core.RegisterAgentFactory(
		proto.Player_ArmsWarrior{},
		proto.Spec_SpecArmsWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewArmsWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ArmsWarrior)
			if !ok {
				panic("Invalid spec value for Arms Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type ArmsWarrior struct {
	*warrior.Warrior

	Options *proto.ArmsWarrior_Options

	TasteForBloodAura   *core.Aura
	SweepingStrikesAura *core.Aura
	SuddenDeathAura     *core.Aura
}

func NewArmsWarrior(character *core.Character, options *proto.Player) *ArmsWarrior {
	armsOptions := options.GetArmsWarrior().Options

	war := &ArmsWarrior{
		Warrior: warrior.NewWarrior(character, armsOptions.ClassOptions, options.TalentsString, warrior.WarriorInputs{
			StanceSnapshot: armsOptions.StanceSnapshot,
		}),
		Options: armsOptions,
	}

	return war
}

const (
	StrikesOfOpportunityHitID int32 = 76858
)

func (war *ArmsWarrior) GetMasteryProcChance() float64 {
	return (8 + 2.2*war.GetMasteryPoints()) / 100
}

func (war *ArmsWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *ArmsWarrior) Initialize() {
	war.Warrior.Initialize()
	war.registerPassives()

	war.registerMortalStrike()
	war.registerOverpower()
	war.registerSlam()
	war.registerSweepingStrikes()
}

func (war *ArmsWarrior) registerPassives() {
	war.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86526)

	war.registerMastery()
	war.registerSeasonedSoldier()
	war.registerSuddenDeath()
	war.registerTasteForBlood()
}

func (war *ArmsWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
}
