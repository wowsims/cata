package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

type BurningTreant struct {
	core.Pet

	owner *Druid

	Fireseed *core.Spell
}

func (druid *Druid) NewBurningTreant() *BurningTreant {
	var burningTreantBaseStats = stats.Stats{stats.SpellCritPercent: 5}

	burningTreant := &BurningTreant{
		Pet:   core.NewPet("Burning Treant (Druid T12 Balance 2P Bonus)", &druid.Character, burningTreantBaseStats, createStatInheritance(), false, true),
		owner: druid,
	}

	druid.AddPet(burningTreant)
	return burningTreant
}

func (treant *BurningTreant) GetPet() *core.Pet {
	return &treant.Pet
}

func (treant *BurningTreant) Initialize() {
	treant.registerFireseedSpell()
}

func (treant *BurningTreant) Reset(_ *core.Simulation) {
}

func (treant *BurningTreant) ExecuteCustomRotation(sim *core.Simulation) {
	if success := treant.Fireseed.Cast(sim, treant.CurrentTarget); !success {
		treant.Disable(sim)
	}
}

func createStatInheritance() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent: ownerStats[stats.SpellHitPercent],
		}
	}
}

func (treant *BurningTreant) registerFireseedSpell() {
	treant.Fireseed = treant.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 99026},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 24,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      0,
				CastTime: time.Second * 2,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   treant.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(5192, 6035)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
