package feral

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/druid"
)

func RegisterFeralDruid() {
	core.RegisterAgentFactory(
		proto.Player_FeralDruid{},
		proto.Spec_SpecFeralDruid,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFeralDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FeralDruid)
			if !ok {
				panic("Invalid spec value for Feral Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewFeralDruid(character *core.Character, options *proto.Player) *FeralDruid {
	feralOptions := options.GetFeralDruid()
	selfBuffs := druid.SelfBuffs{}

	cat := &FeralDruid{
		Druid: druid.New(character, druid.Cat, selfBuffs, options.TalentsString),
	}

	// cat.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	// if feralOptions.Options.ClassOptions.InnervateTarget != nil {
	// 	cat.SelfBuffs.InnervateTarget = feralOptions.Options.ClassOptions.InnervateTarget
	// }

	cat.AssumeBleedActive = feralOptions.Options.AssumeBleedActive
	cat.maxRipTicks = cat.MaxRipTicks()

	cat.EnableEnergyBar(100.0)
	cat.EnableRageBar(core.RageBarOptions{RageMultiplier: 1, MHSwingSpeed: 2.5})

	cat.EnableAutoAttacks(cat, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand:       cat.GetCatWeapon(),
		AutoSwingMelee: true,
	})
	// cat.ReplaceBearMHFunc = func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	// 	return cat.checkReplaceMaul(sim, mhSwingSpell)
	// }

	return cat
}

type FeralDruid struct {
	*druid.Druid

	Rotation FeralDruidRotation

	readyToShift      bool
	readyToGift       bool
	waitingForTick    bool
	maxRipTicks       int32
	berserkUsed       bool
	bleedAura         *core.Aura
	lastShift         time.Duration
	ripRefreshPending bool
	usingHardcodedAPL bool

	rotationAction *core.PendingAction
}

func (cat *FeralDruid) GetDruid() *druid.Druid {
	return cat.Druid
}

func (cat *FeralDruid) Initialize() {
	cat.Druid.Initialize()
	cat.RegisterFeralCatSpells()
}

func (cat *FeralDruid) ApplyTalents() {
	cat.Druid.ApplyTalents()
	cat.MultiplyStat(stats.AttackPower, 1.25) // Aggression passive
}

func (cat *FeralDruid) Reset(sim *core.Simulation) {
	cat.Druid.Reset(sim)
	cat.Druid.ClearForm(sim)
	cat.CatFormAura.Activate(sim)
	cat.readyToShift = false
	cat.waitingForTick = false
	cat.berserkUsed = false
	cat.rotationAction = nil
}
