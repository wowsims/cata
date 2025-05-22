package guardian

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

func RegisterGuardianDruid() {
	core.RegisterAgentFactory(
		proto.Player_GuardianDruid{},
		proto.Spec_SpecGuardianDruid,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewGuardianDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_GuardianDruid)
			if !ok {
				panic("Invalid spec value for Guardian Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewGuardianDruid(character *core.Character, options *proto.Player) *GuardianDruid {
	tankOptions := options.GetGuardianDruid()
	selfBuffs := druid.SelfBuffs{}

	bear := &GuardianDruid{
		Druid:     druid.New(character, druid.Bear, selfBuffs, options.TalentsString),
		Options:   tankOptions.Options,
		vengeance: &core.VengeanceTracker{},
	}

	bear.EnableRageBar(core.RageBarOptions{
		StartingRage:   bear.Options.StartingRage,
		RageMultiplier: 1,
		MHSwingSpeed:   2.5,
	})
	bear.EnableAutoAttacks(bear, core.AutoAttackOptions{
		// Base paw weapon.
		MainHand:       bear.GetBearWeapon(),
		AutoSwingMelee: true,
	})

	bear.RegisterBearFormAura()

	return bear
}

type GuardianDruid struct {
	*druid.Druid

	Options   *proto.GuardianDruid_Options
	vengeance *core.VengeanceTracker

	// Aura references
	EnrageAura *core.Aura

	// Spell references
	Enrage *druid.DruidSpell
}

func (bear *GuardianDruid) GetDruid() *druid.Druid {
	return bear.Druid
}

func (bear *GuardianDruid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.LeaderOfThePack = true
}

func (bear *GuardianDruid) ApplyTalents() {
	// bear.Druid.ApplyTalents()
	bear.applyMastery()
	core.ApplyVengeanceEffect(&bear.Character, bear.vengeance, 84840)
}

func (bear *GuardianDruid) applyMastery() {
	const masteryModPerPoint = 0.02
	bear.AddStat(stats.Armor, bear.GetStat(stats.Armor) * (1.0 + masteryModPerPoint * bear.GetMasteryPoints()))
	bear.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
		currentArmor := bear.GetStat(stats.Armor)
		baseArmor := currentArmor / (1.0 + masteryModPerPoint * core.MasteryRatingToMasteryPoints(oldMasteryRating))
		newArmor := baseArmor * (1.0 + masteryModPerPoint * core.MasteryRatingToMasteryPoints(newMasteryRating))
		bear.AddStatDynamic(sim, stats.Armor, newArmor - currentArmor)
	})
}

func (bear *GuardianDruid) Initialize() {
	bear.Druid.Initialize()
	bear.RegisterFeralTankSpells()
}

func (bear *GuardianDruid) Reset(sim *core.Simulation) {
	bear.Druid.Reset(sim)
	bear.Druid.ClearForm(sim)
	bear.BearFormAura.Activate(sim)
	bear.Druid.PseudoStats.Stunned = false
}
