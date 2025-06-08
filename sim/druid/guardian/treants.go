package guardian

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

type GuardianTreant struct {
	*druid.DefaultTreantImpl
}

func (bear *GuardianDruid) newTreant() *GuardianTreant {
	treant := &GuardianTreant{
		DefaultTreantImpl: bear.NewDefaultTreant(druid.TreantConfig{
			StatInheritance: func(ownerStats stats.Stats) stats.Stats {
				combinedHitExp := 0.5 * (ownerStats[stats.HitRating] + ownerStats[stats.ExpertiseRating])

				return stats.Stats{
					stats.Health:              0.4 * ownerStats[stats.Health],
					stats.Armor:               4 * ownerStats[stats.Armor],
					stats.AttackPower:         1.2 * ownerStats[stats.AttackPower],
					stats.HitRating:           combinedHitExp,
					stats.ExpertiseRating:     combinedHitExp,
					stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
				}
			},

			EnableAutos:             true,
			WeaponDamageCoefficient: 3.20000004768,
		}),
	}

	treant.PseudoStats.DamageDealtMultiplier *= 0.2
	bear.AddPet(treant)

	return treant
}

func (bear *GuardianDruid) registerTreants() {
	for idx := range bear.Treants {
		bear.Treants[idx] = bear.newTreant()
	}
}

func (treant *GuardianTreant) Enable(sim *core.Simulation) {
	treant.DefaultTreantImpl.Enable(sim)
	treant.ExtendGCDUntil(sim, sim.CurrentTime + time.Second * 15)
}
