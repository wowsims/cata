package shaman

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type EarthElemental struct {
	core.Pet

	shamanOwner *Shaman
}

var EarthElementalSpellPowerScaling = 0.749 // Estimated from beta testing

func (shaman *Shaman) NewEarthElemental(isGuardian bool) *EarthElemental {
	earthElemental := &EarthElemental{
		Pet: core.NewPet(core.PetConfig{
			Name:                            core.Ternary(isGuardian, "Greater Earth Elemental", "Primal Earth Elemental"),
			Owner:                           &shaman.Character,
			BaseStats:                       earthElementalPetBaseStats,
			StatInheritance:                 shaman.earthElementalStatInheritance(isGuardian),
			EnabledOnStart:                  false,
			IsGuardian:                      isGuardian,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
		shamanOwner: shaman,
	}
	baseMeleeDamage := core.TernaryFloat64(isGuardian, 1222.72, 2114.25) //Estimated from beta testing
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

	earthElemental.OnPetEnable = earthElemental.enable
	earthElemental.OnPetDisable = earthElemental.disable

	shaman.AddPet(earthElemental)

	return earthElemental
}

func (earthElemental *EarthElemental) enable(sim *core.Simulation) {
}

func (earthElemental *EarthElemental) disable(sim *core.Simulation) {

}

func (earthElemental *EarthElemental) GetPet() *core.Pet {
	return &earthElemental.Pet
}

func (earthElemental *EarthElemental) Initialize() {

}

func (earthElemental *EarthElemental) Reset(_ *core.Simulation) {

}

func (earthElemental *EarthElemental) ExecuteCustomRotation(sim *core.Simulation) {

}

var earthElementalPetBaseStats = stats.Stats{}

func (shaman *Shaman) earthElementalStatInheritance(isGuardian bool) core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		ownerSpellHitPercent := ownerStats[stats.SpellHitPercent]
		ownerPhysicalHitPercent := ownerStats[stats.PhysicalHitPercent]
		ownerExpertiseRating := ownerStats[stats.ExpertiseRating]
		ownerSpellCritPercent := ownerStats[stats.SpellCritPercent]
		ownerPhysicalCritPercent := ownerStats[stats.PhysicalCritPercent]
		ownerHasteRating := ownerStats[stats.HasteRating]

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * core.TernaryFloat64(isGuardian, 0.75, 0.9),
			stats.AttackPower: ownerStats[stats.SpellPower] * core.TernaryFloat64(isGuardian, EarthElementalSpellPowerScaling, EarthElementalSpellPowerScaling*1.8),

			stats.PhysicalHitPercent:  max(ownerSpellHitPercent/2, ownerPhysicalHitPercent),
			stats.SpellHitPercent:     max(ownerSpellHitPercent, ownerExpertiseRating/core.ExpertisePerQuarterPercentReduction/4+ownerPhysicalHitPercent),
			stats.ExpertiseRating:     max(ownerSpellHitPercent*core.ExpertisePerQuarterPercentReduction*2, ownerExpertiseRating),
			stats.SpellCritPercent:    ownerSpellCritPercent,
			stats.PhysicalCritPercent: ownerPhysicalCritPercent,
			stats.HasteRating:         ownerHasteRating,
		}
	}
}
