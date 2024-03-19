package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/warrior"
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
}

func NewArmsWarrior(character *core.Character, options *proto.Player) *ArmsWarrior {
	armsOptions := options.GetArmsWarrior().Options

	war := &ArmsWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{
			StanceSnapshot: armsOptions.StanceSnapshot,
		}),
		Options: armsOptions,
	}

	rbo := core.RageBarOptions{
		StartingRage:   armsOptions.ClassOptions.StartingRage,
		RageMultiplier: 1.25, // Endless Rage is now part of Anger Management, now an Arms specialization ability
	}
	if mh := war.GetMHWeapon(); mh != nil {
		rbo.MHSwingSpeed = mh.SwingSpeed
	}
	war.EnableRageBar(rbo)

	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	war.RegisterSpecializationEffects()

	return war
}

func (war *ArmsWarrior) RegisterSpecializationEffects() {
	// Strikes of Opportunity
	war.RegisterMastery()

	// Anger Management (flat rage multiplier is set in the RageBarOptions above) (12296)
	rageMetrics := war.NewRageMetrics(core.ActionID{SpellID: 12296})
	war.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: time.Second * 3,
			OnAction: func(sim *core.Simulation) {
				war.AddRage(sim, 1, rageMetrics)
				war.LastAMTick = sim.CurrentTime
			},
		})
	})

	// Two-Handed Weapon Specialization (12712)
	war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1 + 0.12
}

func (war *ArmsWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *ArmsWarrior) Initialize() {
	war.Warrior.Initialize()

	// if war.Options.UseRecklessness {
	// 	war.RegisterRecklessnessCD()
	// }

	// if war.Options.ClassOptions.UseShatteringThrow {
	// 	war.RegisterShatteringThrowCD()
	// }

	// war.BattleStanceAura.BuildPhase = core.CharacterBuildPhaseTalents
}

// func (war *ArmsWarrior) Reset(sim *core.Simulation) {
// 	war.Warrior.Reset(sim)
// 	war.BattleStanceAura.Activate(sim)
// 	war.Stance = warrior.BattleStance
// }
