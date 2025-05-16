package druid

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Treant struct {
	core.Pet

	druidOwner *Druid
}

type Treants struct {
	Treant1 *Treant
	Treant2 *Treant
	Treant3 *Treant
}

func (treants *Treants) EnableWithTimeout(sim *core.Simulation) {
	treants.Treant1.EnableWithTimeout(sim, treants.Treant1, time.Second*30)
	treants.Treant2.EnableWithTimeout(sim, treants.Treant2, time.Second*30)
	treants.Treant3.EnableWithTimeout(sim, treants.Treant3, time.Second*30)
}

func (treants *Treants) CancelGCDTimer(sim *core.Simulation) {
	treants.Treant1.CancelGCDTimer(sim)
	treants.Treant2.CancelGCDTimer(sim)
	treants.Treant3.CancelGCDTimer(sim)
}

var treantBaseStats = stats.Stats{
	stats.Stamina:   422,
	stats.Spirit:    116,
	stats.Intellect: 120,
	stats.Armor:     11092,

	stats.Agility:     1218,
	stats.Strength:    476,
	stats.AttackPower: -20,

	stats.PhysicalCritPercent: 1.1515 + 1.8,
}

func (druid *Druid) NewTreant() *Treant {
	treant := &Treant{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Treant",
			Owner:           &druid.Character,
			BaseStats:       treantBaseStats,
			StatInheritance: druid.makeStatInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      false,
		}),
		druidOwner: druid,
	}

	treant.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	treant.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[proto.Class_ClassWarrior])

	treant.EnableAutoAttacks(treant, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  252,
			BaseDamageMax:  357,
			SwingSpeed:     1.9,
			CritMultiplier: druid.DefaultCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	druid.AddPet(treant)
	return treant
}

func (treant *Treant) GetPet() *core.Pet {
	return &treant.Pet
}

func (treant *Treant) Reset(sim *core.Simulation) {
	treant.Disable(sim)
	if sim.Log != nil {
		treant.Log(sim, "Base Stats: %s", treantBaseStats)
		inheritedStats := treant.druidOwner.makeStatInheritance()(treant.druidOwner.GetStats())
		treant.Log(sim, "Inherited Stats: %s", inheritedStats)
		treant.Log(sim, "Total Stats: %s", treant.GetStats())
	}
}

const PetExpertiseScale = 3.25

func (druid *Druid) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		treantHitChance := ownerStats[stats.SpellHitPercent] * 8 / 17

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3189,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.SpellPower] * 0.65,

			stats.HitRating:       treantHitChance * core.PhysicalHitRatingPerHitPercent,
			stats.ExpertiseRating: math.Floor(treantHitChance*PetExpertiseScale) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}

func (treant *Treant) Initialize() {
}

func (treant *Treant) ExecuteCustomRotation(_ *core.Simulation) {
}
