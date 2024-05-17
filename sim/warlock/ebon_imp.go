package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

type EbonImpPet struct {
	core.Pet
}

func (warlock *Warlock) NewEbonImp() *EbonImpPet {
	baseStats := stats.Stats{
		stats.MeleeCrit: 5.0 * core.CritRatingPerCritChance, // rough guess

		// rough guess; definitely some misses and dodges, even if the warlock is hit capped
		// does not seem to scale with gear or if it does then only by a small fraction
		stats.MeleeHit:  7 * core.MeleeHitRatingPerHitChance,
		stats.Expertise: 24 * core.ExpertisePerQuarterPercentReduction,
	}

	statInheritance := func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.MeleeHaste: ownerStats[stats.SpellHaste],
			stats.MeleeCrit:  ownerStats[stats.SpellCrit],
		}
	}

	imp := &EbonImpPet{Pet: core.NewPet("Ebon Imp", &warlock.Character, baseStats, statInheritance, false, true)}
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
