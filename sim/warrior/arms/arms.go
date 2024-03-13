package arms

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
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

	// rbo := core.RageBarOptions{
	// 	StartingRage:   armsOptions.ClassOptions.StartingRage,
	// 	RageMultiplier: core.TernaryFloat64(war.Talents.EndlessRage, 1.25, 1),
	// }
	// if mh := war.GetMHWeapon(); mh != nil {
	// 	rbo.MHSwingSpeed = mh.SwingSpeed
	// }
	// if oh := war.GetOHWeapon(); oh != nil {
	// 	rbo.OHSwingSpeed = oh.SwingSpeed
	// }

	// war.EnableRageBar(rbo)
	// war.EnableAutoAttacks(war, core.AutoAttackOptions{
	// 	MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
	// 	OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
	// 	AutoSwingMelee: true,
	// 	ReplaceMHSwing: war.TryHSOrCleave,
	// })

	return war
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
