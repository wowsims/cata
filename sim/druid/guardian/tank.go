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
		Druid:   druid.New(character, druid.Bear, selfBuffs, options.TalentsString),
		Options: tankOptions.Options,
	}

	bear.registerTreants()

	bear.EnableRageBar(core.RageBarOptions{
		StartingRage:       bear.Options.StartingRage,
		BaseRageMultiplier: 2.5,
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

	Options *proto.GuardianDruid_Options

	Treants GuardianTreants

	// Aura references
	EnrageAura          *core.Aura
	SavageDefenseAura   *core.Aura
	SonOfUrsocAura      *core.Aura
	ToothAndClawBuff    *core.Aura
	ToothAndClawDebuffs core.AuraArray

	// Spell references
	Enrage        *druid.DruidSpell
	ForceOfNature *druid.DruidSpell
	SavageDefense *druid.DruidSpell
	SonOfUrsoc    *druid.DruidSpell
}

func (bear *GuardianDruid) GetDruid() *druid.Druid {
	return bear.Druid
}

func (bear *GuardianDruid) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.LeaderOfThePack = true
}

func (bear *GuardianDruid) ApplyTalents() {
	bear.Druid.ApplyTalents()
	bear.applyMastery()
	bear.applyThickHide()
	bear.applyLeatherSpecialization()
	bear.RegisterVengeance(84840, bear.BearFormAura)
	bear.registerIncarnation()
	bear.registerForceOfNature()
}

func (bear *GuardianDruid) applyMastery() {
	const baseMasteryMod = 1.16
	const masteryModPerPoint = 0.02

	armorMultiplierDep := bear.NewDynamicMultiplyStat(stats.Armor, baseMasteryMod+masteryModPerPoint*bear.GetMasteryPoints())

	bear.AddOnMasteryStatChanged(func(sim *core.Simulation, _ float64, newMasteryRating float64) {
		bear.UpdateDynamicStatDep(sim, armorMultiplierDep, baseMasteryMod+masteryModPerPoint*core.MasteryRatingToMasteryPoints(newMasteryRating))
	})

	bear.BearFormAura.AttachStatDependency(armorMultiplierDep)
}

func (bear *GuardianDruid) applyThickHide() {
	// Back out the additional multiplier needed to reach 4.3x total (+330%)
	const thickHideBearMulti = 4.3 / druid.BaseBearArmorMulti
	bear.BearFormAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		bear.ApplyDynamicEquipScaling(sim, stats.Armor, thickHideBearMulti)
	})
	bear.BearFormAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		bear.RemoveDynamicEquipScaling(sim, stats.Armor, thickHideBearMulti)
	})
	bear.ApplyEquipScaling(stats.Armor, thickHideBearMulti)

	// Magical DR
	bear.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 0.75
	bear.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 0.75
	bear.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 0.75
	bear.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 0.75
	bear.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 0.75
	bear.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 0.75

	// Physical DR
	bear.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= 0.88

	// Crit immunity
	bear.PseudoStats.ReducedCritTakenChance += 0.06
}

func (bear *GuardianDruid) applyLeatherSpecialization() {
	bear.GuardianLeatherSpecTracker = bear.RegisterArmorSpecializationTracker(proto.ArmorType_ArmorTypeLeather, 86096)
	bear.GuardianLeatherSpecDep = bear.NewDynamicMultiplyStat(stats.Stamina, 1.05)

	// Need redundant enabling/disabling of the dep both here and in forms.go because we
	// don't know whether the leather spec tracker or Bear Form will activate first.
	bear.GuardianLeatherSpecTracker.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		if bear.InForm(druid.Bear) {
			bear.EnableBuildPhaseStatDep(sim, bear.GuardianLeatherSpecDep)
		}
	})

	bear.GuardianLeatherSpecTracker.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		if bear.InForm(druid.Bear) {
			bear.DisableBuildPhaseStatDep(sim, bear.GuardianLeatherSpecDep)
		}
	})
}

func (bear *GuardianDruid) Initialize() {
	bear.Druid.Initialize()
	bear.RegisterFeralTankSpells()
	bear.registerEnrageSpell()
	bear.registerSavageDefenseSpell()
	bear.registerToothAndClawPassive()
	bear.ApplyPrimalFury()
	bear.ApplyLeaderOfThePack()
	bear.ApplyNurturingInstinct()
}

func (bear *GuardianDruid) Reset(sim *core.Simulation) {
	bear.Druid.Reset(sim)
	bear.Druid.ClearForm(sim)
	bear.BearFormAura.Activate(sim)
	bear.Druid.PseudoStats.Stunned = false
}
