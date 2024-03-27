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
		RageMultiplier: 1.25, // Endless Rage is now part of Anger Management, now an Arms specialization passive
	}
	if mh := war.GetMHWeapon(); mh != nil {
		rbo.MHSwingSpeed = mh.SwingSpeed
	}
	war.EnableRageBar(rbo)

	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

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
	war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.12
}

const (
	StrikesOfOpportunityHitID int32 = 76858
)

func (war *ArmsWarrior) GetMasteryProcChance() float64 {
	return (17.6 + 2.2*war.GetMasteryPoints()) / 100
}

func (war *ArmsWarrior) RegisterMastery() {
	// TODO:
	//	test what things the extra attack can proc
	//	does the extra attack use the same hit table
	//  can it proc off of missed/dodged/parried attacks
	//
	// 4.3.3 simcraft implements SoO as a standard autoattack with a 0.5s ICD
	procAttackConfig := *war.AutoAttacks.MHConfig()
	procAttackConfig.ActionID = core.ActionID{SpellID: StrikesOfOpportunityHitID, Tag: procAttackConfig.ActionID.Tag}
	procAttack := war.RegisterSpell(procAttackConfig)

	icd := core.Cooldown{
		Timer:    war.NewTimer(),
		Duration: time.Millisecond * 500,
	}

	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Strikes of Opportunity",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ActionID.IsOtherAction(proto.OtherAction_OtherActionAttack) {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			proc := war.GetMasteryProcChance()
			if sim.RandomFloat(aura.Label) < proc {
				icd.Use(sim)
				procAttack.Cast(sim, result.Target)
			}
		},
	}))
}

func (war *ArmsWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *ArmsWarrior) Initialize() {
	war.Warrior.Initialize()
	war.RegisterSpecializationEffects()

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
