package feral

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
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
	cat.CannotShredTarget = feralOptions.Options.CannotShredTarget
	// TODO: Fix this to work with the new talent system.
	// cat.maxRipTicks = cat.MaxRipTicks()
	// cat.primalMadnessBonus = 10.0 * float64(cat.Talents.PrimalMadness)
	cat.maxRipTicks = 0
	cat.primalMadnessBonus = 0

	cat.EnableEnergyBar(core.EnergyBarOptions{
		MaxComboPoints: 5,
		MaxEnergy:      100.0,
		UnitClass:      proto.Class_ClassDruid,
	})
	cat.EnableRageBar(core.RageBarOptions{BaseRageMultiplier: 2.5})

	cat.EnableAutoAttacks(cat, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand:       cat.GetCatWeapon(),
		AutoSwingMelee: true,
	})

	cat.RegisterCatFormAura()
	cat.RegisterBearFormAura()

	return cat
}

type FeralDruid struct {
	*druid.Druid

	// Rotation FeralDruidRotation

	readyToShift       bool
	readyToGift        bool
	waitingForTick     bool
	maxRipTicks        int32
	primalMadnessBonus float64
	berserkUsed        bool
	bleedAura          *core.Aura
	tempSnapshotAura   *core.Aura
	lastShift          time.Duration
	cachedRipEndThresh time.Duration
	nextActionAt       time.Duration
	usingHardcodedAPL  bool
	// pendingPool        *PoolingActions
	// pendingPoolWeaves  *PoolingActions
}

func (cat *FeralDruid) GetDruid() *druid.Druid {
	return cat.Druid
}

func (cat *FeralDruid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.LeaderOfThePack = true
}

func (cat *FeralDruid) Initialize() {
	cat.Druid.Initialize()
	cat.RegisterFeralCatSpells()
	cat.ApplyPrimalFury()
	cat.ApplyLeaderOfThePack()
	cat.ApplyNurturingInstinct()

	snapshotHandler := func(aura *core.Aura, sim *core.Simulation) {
		previousRipSnapshotPower := cat.Rip.NewSnapshotPower
		previousRakeSnapshotPower := cat.Rake.NewSnapshotPower
		cat.UpdateBleedPower(cat.Rip, sim, cat.CurrentTarget, false, true)
		cat.UpdateBleedPower(cat.Rake, sim, cat.CurrentTarget, false, true)

		if cat.Rip.NewSnapshotPower > previousRipSnapshotPower+0.001 {
			if !cat.tempSnapshotAura.IsActive() || (aura.ExpiresAt() < cat.tempSnapshotAura.ExpiresAt()) {
				cat.tempSnapshotAura = aura

				if sim.Log != nil {
					cat.Log(sim, "New bleed snapshot aura found: %s", aura.ActionID)
				}
			}
		} else if cat.tempSnapshotAura.IsActive() {
			cat.Rip.NewSnapshotPower = previousRipSnapshotPower
			cat.Rake.NewSnapshotPower = previousRakeSnapshotPower
		} else {
			cat.tempSnapshotAura = nil
		}
	}

	// cat.TigersFuryAura.ApplyOnGain(snapshotHandler)
	// cat.TigersFuryAura.ApplyOnExpire(snapshotHandler)
	cat.AddOnTemporaryStatsChange(func(sim *core.Simulation, buffAura *core.Aura, _ stats.Stats) {
		snapshotHandler(buffAura, sim)
	})
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
	cat.nextActionAt = -core.NeverExpires

	// Reset snapshot power values until first cast
	cat.Rip.CurrentSnapshotPower = 0
	cat.Rip.NewSnapshotPower = 0
	cat.Rake.CurrentSnapshotPower = 0
	cat.Rake.NewSnapshotPower = 0
	cat.tempSnapshotAura = nil
}
