package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type EarthElemental struct {
	core.Pet

	Pulverize *core.Spell

	shamanOwner *Shaman
}

var EarthElementalSpellPowerScaling = 1.3 // Estimated from beta testing

func (shaman *Shaman) NewEarthElemental(isGuardian bool) *EarthElemental {
	earthElemental := &EarthElemental{
		Pet: core.NewPet(core.PetConfig{
			Name:                            core.Ternary(isGuardian, "Greater Earth Elemental", "Primal Earth Elemental"),
			Owner:                           &shaman.Character,
			BaseStats:                       shaman.earthElementalBaseStats(isGuardian),
			StatInheritance:                 shaman.earthElementalStatInheritance(isGuardian),
			EnabledOnStart:                  false,
			IsGuardian:                      isGuardian,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
		shamanOwner: shaman,
	}
	scalingDamage := shaman.CalcScalingSpellDmg(1.3)
	baseMeleeDamage := core.TernaryFloat64(isGuardian, scalingDamage, scalingDamage*1.8)
	earthElemental.EnableAutoAttacks(earthElemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  baseMeleeDamage,
			BaseDamageMax:  baseMeleeDamage,
			SwingSpeed:     2,
			CritMultiplier: earthElemental.DefaultCritMultiplier(),
			SpellSchool:    core.SpellSchoolPhysical,
		},
		AutoSwingMelee: true,
	})

	earthElemental.OnPetEnable = earthElemental.enable(isGuardian)
	earthElemental.OnPetDisable = earthElemental.disable

	shaman.AddPet(earthElemental)

	return earthElemental
}

func (earthElemental *EarthElemental) enable(isGuardian bool) func(*core.Simulation) {
	return func(sim *core.Simulation) {
		earthElemental.EnableDynamicStats(earthElemental.shamanOwner.earthElementalStatInheritance(isGuardian))
	}
}

func (earthElemental *EarthElemental) disable(sim *core.Simulation) {

}

func (earthElemental *EarthElemental) GetPet() *core.Pet {
	return &earthElemental.Pet
}

func (earthElemental *EarthElemental) Initialize() {
	earthElemental.registerPulverize()
}

func (earthElemental *EarthElemental) Reset(_ *core.Simulation) {

}

func (earthElemental *EarthElemental) ExecuteCustomRotation(sim *core.Simulation) {
	/*
		Pulverize on cd
	*/
	target := earthElemental.CurrentTarget

	earthElemental.TryCast(sim, target, earthElemental.Pulverize)

	if !earthElemental.GCD.IsReady(sim) {
		return
	}

	minCd := earthElemental.Pulverize.CD.ReadyAt()
	earthElemental.ExtendGCDUntil(sim, max(minCd, sim.CurrentTime+time.Second))
}

func (earthElemental *EarthElemental) TryCast(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	if !spell.Cast(sim, target) {
		return false
	}
	// all spell casts reset the elemental's swing timer
	earthElemental.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+spell.CurCast.CastTime, false)
	return true
}

func (shaman *Shaman) earthElementalBaseStats(isGuardian bool) stats.Stats {
	return stats.Stats{
		stats.Stamina: core.TernaryFloat64(isGuardian, 10457, 10457*1.5),
	}
}

func (shaman *Shaman) earthElementalStatInheritance(isGuardian bool) core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerHitRating := ownerStats[stats.HitRating]
		ownerExpertiseRating := ownerStats[stats.ExpertiseRating]
		ownerSpellCritPercent := ownerStats[stats.SpellCritPercent]
		ownerPhysicalCritPercent := ownerStats[stats.PhysicalCritPercent]
		ownerHasteRating := ownerStats[stats.HasteRating]
		hitExpRating := (ownerHitRating + ownerExpertiseRating) / 2
		critPercent := max(ownerPhysicalCritPercent, ownerSpellCritPercent)

		power := core.TernaryFloat64(shaman.Spec == proto.Spec_SpecEnhancementShaman, ownerStats[stats.AttackPower]*0.65, ownerStats[stats.SpellPower])

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * core.TernaryFloat64(isGuardian, 1, 1.5),
			stats.AttackPower: power * core.TernaryFloat64(isGuardian, EarthElementalSpellPowerScaling, EarthElementalSpellPowerScaling*1.8),

			stats.HitRating:           hitExpRating,
			stats.ExpertiseRating:     hitExpRating,
			stats.SpellCritPercent:    critPercent,
			stats.PhysicalCritPercent: critPercent,
			stats.HasteRating:         ownerHasteRating,
		}
	}
}
