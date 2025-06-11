package fury

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
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

	BloodsurgeAura  *core.Aura
	MeatCleaverAura *core.Aura
}

func NewFuryWarrior(character *core.Character, options *proto.Player) *FuryWarrior {
	furyOptions := options.GetFuryWarrior().Options

	war := &FuryWarrior{
		Warrior: warrior.NewWarrior(character, furyOptions.ClassOptions, options.TalentsString, warrior.WarriorInputs{
			StanceSnapshot: furyOptions.StanceSnapshot,
		}),
		Options: furyOptions,
	}

	war.ApplySyncType(furyOptions.SyncType)

	return war
}

func (war *FuryWarrior) GetMasteryBonusMultiplier() float64 {
	return (8 + 1.4*war.GetMasteryPoints()) / 100
}

func (war *FuryWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *FuryWarrior) Initialize() {
	war.Warrior.Initialize()
	war.registerPassives()
	war.registerBloodthirst()
	war.registerRagingBlow()
	war.registerWildStrike()
}

func (war *FuryWarrior) registerPassives() {
	war.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86526)

	war.registerCrazedBerserker()
	war.registerFlurry()
	war.registerBloodsurge()
	war.registerMeatCleaver()
	war.registerSingleMindedFuryOrTitansGrip()
	war.registerUnshackledFury()
}

func (war *FuryWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
}

func (war *FuryWarrior) ApplySyncType(syncType proto.WarriorSyncType) {
	if syncType == proto.WarriorSyncType_WarriorSyncMainhandOffhandSwings {
		war.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			aa := &war.AutoAttacks
			if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed(); nextMHSwingAt > aa.OffhandSwingAt() {
				aa.SetOffhandSwingAt(nextMHSwingAt)
			}

			return mhSwingSpell
		})

	} else {
		war.AutoAttacks.SetReplaceMHSwing(nil)
	}
}
