package fury

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/warrior"
)

func RegisterFuryWarrior() {
	core.RegisterAgentFactory(
		proto.Player_FuryWarrior{},
		proto.Spec_SpecFuryWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFuryWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FuryWarrior)
			if !ok {
				panic("Invalid spec value for Fury Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type FuryWarrior struct {
	*warrior.Warrior

	Options *proto.FuryWarrior_Options
}

func NewFuryWarrior(character *core.Character, options *proto.Player) *FuryWarrior {
	furyOptions := options.GetFuryWarrior().Options

	war := &FuryWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{
			StanceSnapshot: furyOptions.StanceSnapshot,
		}),
		Options: furyOptions,
	}

	rbo := core.RageBarOptions{
		StartingRage:   furyOptions.ClassOptions.StartingRage,
		RageMultiplier: 1.0,
	}
	if mh := war.GetMHWeapon(); mh != nil {
		rbo.MHSwingSpeed = mh.SwingSpeed
	}
	if oh := war.GetOHWeapon(); oh != nil {
		rbo.OHSwingSpeed = oh.SwingSpeed
	}

	war.EnableRageBar(rbo)
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return war
}

func (war *FuryWarrior) RegisterSpecializationEffects() {

	// Unshackled Fury
	// The actual effects of Unshackled Fury need to be handled by specific spells
	// as it modifies the "benefit" of them (e.g. it both increases Raging Blow's damage
	// and Enrage's damage bonus)
	war.EnrageEffectMultiplier = war.GetMasteryBonusMultiplier()
	war.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		war.EnrageEffectMultiplier = war.GetMasteryBonusMultiplier()
	})

	// Dual Wield specialization
	war.AutoAttacks.OHConfig().DamageMultiplier *= 1.25

	// Precision
	war.AddStat(stats.MeleeHit, 3*core.MeleeHitRatingPerHitChance)
}

func (war *FuryWarrior) GetMasteryBonusMultiplier() float64 {
	return 1 + (11.2+5.6*war.GetMasteryPoints())/100
}

func (war *FuryWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *FuryWarrior) Initialize() {
	war.Warrior.Initialize()
	war.RegisterSpecializationEffects()
	war.RegisterBloodthirst()
}

func (war *FuryWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
}
