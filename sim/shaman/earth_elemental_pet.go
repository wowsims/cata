package shaman

import (
	"math"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type EarthElemental struct {
	core.Pet

	shamanOwner *Shaman
}

func (shaman *Shaman) NewEarthElemental(bonusSpellPower float64) *EarthElemental {
	earthElemental := &EarthElemental{
		Pet:         core.NewPet("Greater Earth Elemental", &shaman.Character, earthElementalPetBaseStats, shaman.earthElementalStatInheritance(), false, true),
		shamanOwner: shaman,
	}
	earthElemental.EnableManaBar()
	earthElemental.EnableAutoAttacks(earthElemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  354, //Estimated from beta testing
			BaseDamageMax:  396, //Estimated from beta testing
			SwingSpeed:     2,
			CritMultiplier: 2, //Estimated from beta testing
			SpellSchool:    core.SpellSchoolPhysical,
		},
		AutoSwingMelee: true,
	})

	if bonusSpellPower > 0 {
		earthElemental.AddStat(stats.AttackPower, float64(bonusSpellPower)*0.749) // 0.107 * 7
	}

	if shaman.Race == proto.Race_RaceDraenei {
		earthElemental.AddStats(stats.Stats{
			stats.HitRating:       -core.PhysicalHitRatingPerHitPercent,
			stats.ExpertiseRating: math.Floor(-core.SpellHitRatingPerHitPercent * 27 / 16),
		})
	}

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

var earthElementalPetBaseStats = stats.Stats{
	stats.Health:      7976, //TODO need to be more accurate
	stats.Stamina:     0,
	stats.AttackPower: 0,

	stats.PhysicalCritPercent: 6.8, //TODO need testing
}

func (shaman *Shaman) earthElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		flooredOwnerSpellHitPercent := math.Floor(ownerStats[stats.SpellHitPercent])
		hitRatingFromOwner := flooredOwnerSpellHitPercent * core.SpellHitRatingPerHitPercent

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 1.06,     //TODO need to be more accurate
			stats.AttackPower: ownerStats[stats.SpellPower] * 0.749, // 0.107 * 7 TODO need to be more accurate

			stats.HitRating: hitRatingFromOwner,

			/*
				TODO working on figuring this out, getting close need more trials. will need to remove specific buffs,
				ie does not gain the benefit from draenei buff.
			*/
			stats.ExpertiseRating: math.Floor(hitRatingFromOwner * 27 / 16),
		}
	}
}
