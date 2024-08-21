package druid

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (druid *Druid) NewTreant() *Treant {
	treant := &Treant{
		Pet:        core.NewPet("Treant", &druid.Character, treantBaseStats, druid.makeStatInheritance(), false, false),
		druidOwner: druid,
	}

	treant.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	treant.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance/83.3)

	treant.PseudoStats.DamageDealtMultiplier = 1

	treant.EnableAutoAttacks(treant, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  252,
			BaseDamageMax:  357,
			SwingSpeed:     1.65,
			CritMultiplier: druid.BalanceCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	treant.Pet.OnPetEnable = treant.enable
	treant.Pet.OnPetDisable = treant.disable

	druid.AddPet(treant)
	return treant
}

func (treant *Treant) GetPet() *core.Pet {
	return &treant.Pet
}

func (treant *Treant) enable(sim *core.Simulation) {
	treant.snapshotStat = stats.Stats{
		stats.AttackPower: treant.druidOwner.GetStat(stats.SpellPower) * 0.65,

		stats.ArcaneResistance: treant.druidOwner.GetStat(stats.ArcaneResistance) * 0.4,
		stats.FireResistance:   treant.druidOwner.GetStat(stats.FireResistance) * 0.4,
		stats.FrostResistance:  treant.druidOwner.GetStat(stats.FrostResistance) * 0.4,
		stats.NatureResistance: treant.druidOwner.GetStat(stats.NatureResistance) * 0.4,
		stats.ShadowResistance: treant.druidOwner.GetStat(stats.ShadowResistance) * 0.4,
	}

	treant.AddStatsDynamic(sim, treant.snapshotStat)
}

func (treant *Treant) disable(sim *core.Simulation) {
	treant.AddStatsDynamic(sim, treant.snapshotStat.Invert())
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

// ///////////////////////////////////////
type Treant struct {
	core.Pet
	druidOwner *Druid

	snapshotStat stats.Stats
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

	stats.MeleeCrit: 0.05,
}

const PetExpertiseScale = 3.25

func (druid *Druid) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		treantHitChance := (ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance) * 0.08 / 0.17
		//hitRatingFromOwner := math.Floor(ownerHitChance) * core.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3189,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.AttackPower] * 0.65,

			stats.MeleeHit:  math.Floor(treantHitChance) * core.MeleeHitRatingPerHitChance,
			stats.Expertise: math.Floor(math.Floor(treantHitChance)*PetExpertiseScale) * core.ExpertisePerQuarterPercentReduction,
		}
	}
}

func (treant *Treant) Initialize() {
	// Nothing
}

func (treant *Treant) ExecuteCustomRotation(_ *core.Simulation) {
}
