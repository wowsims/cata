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
		StartingRage:       furyOptions.ClassOptions.StartingRage,
		BaseRageMultiplier: 1,
	}

	war.PrecisionKnown = true

	war.EnableRageBar(rbo)
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	return war
}

func (war *FuryWarrior) RegisterSpecializationEffects() {

	// Unshackled Fury
	// The actual effects of Unshackled Fury need to be handled by specific spells
	// as it modifies the "benefit" of them (e.g. it both increases Raging Blow's damage
	// and Enrage's damage bonus)
	war.EnrageEffectMultiplier = war.GetMasteryBonusMultiplier(war.GetMasteryPoints())
	war.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		war.EnrageEffectMultiplier = war.GetMasteryBonusMultiplier(war.GetMasteryPoints())
	})

	// Dual Wield specialization
	war.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMeleeOH,
		FloatValue: 0.25,
	})

	// Precision
	war.AddStat(stats.PhysicalHitPercent, 3)
	war.AutoAttacks.MHConfig().DamageMultiplier *= 1.4
	war.AutoAttacks.OHConfig().DamageMultiplier *= 1.4
}

func (war *FuryWarrior) GetMasteryBonusMultiplier(masteryPoints float64) float64 {
	return 1 + (11.2+5.6*masteryPoints)/100
}

func (war *FuryWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *FuryWarrior) Initialize() {
	war.Warrior.Initialize()
	war.RegisterSpecializationEffects()
	war.RegisterBloodthirst()
}

func (war *FuryWarrior) ApplyTalents() {}

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
