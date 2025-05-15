package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type EbonImpPet struct {
	core.Pet
}

func (warlock *Warlock) NewEbonImp() *EbonImpPet {
	baseStats := stats.Stats{
		stats.PhysicalCritPercent: 5, // rough guess

		// rough guess; definitely some misses and dodges, even if the warlock is hit capped
		// does not seem to scale with gear or if it does then only by a small fraction
		stats.PhysicalHitPercent: 7,
		stats.ExpertiseRating:    24 * core.ExpertisePerQuarterPercentReduction,
	}

	statInheritance := func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.HasteRating:         ownerStats[stats.HasteRating],
			stats.PhysicalCritPercent: ownerStats[stats.SpellCritPercent],
		}
	}

	imp := &EbonImpPet{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Ebon Imp",
			Owner:           &warlock.Character,
			BaseStats:       baseStats,
			StatInheritance: statInheritance,
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
	}
	imp.EnableAutoAttacks(imp, core.AutoAttackOptions{
		MainHand: core.Weapon{
			// determined from logs on target dummies (damage before armor and modifiers)
			// TODO: are they affected by command? If not correct for that
			// TODO: they also sometimes get parried because of their spawn position ..
			BaseDamageMin:  876,
			BaseDamageMax:  1251,
			SwingSpeed:     2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	warlock.AddPet(imp)

	return imp
}

func (imp *EbonImpPet) GetPet() *core.Pet {
	return &imp.Pet
}

func (imp *EbonImpPet) Initialize() {}

func (imp *EbonImpPet) Reset(_ *core.Simulation) {}

func (imp *EbonImpPet) ExecuteCustomRotation(sim *core.Simulation) {
	imp.SetRotationTimer(sim, time.Duration(1<<63-1))
}
