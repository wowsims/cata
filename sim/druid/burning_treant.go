package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type BurningTreant struct {
	core.Pet

	owner *Druid

	Fireseed *core.Spell
}

func (druid *Druid) NewBurningTreant() *BurningTreant {
	baseStats := stats.Stats{stats.SpellCritPercent: 0}

	statInheritance := func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent: ownerStats[stats.SpellHitPercent],
		}
	}

	burningTreant := &BurningTreant{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Burning Treant",
			Owner:           &druid.Character,
			BaseStats:       baseStats,
			StatInheritance: statInheritance,
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
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
	if treant.Fireseed.CanCast(sim, treant.CurrentTarget) {
		treant.Fireseed.Cast(sim, treant.CurrentTarget)
		delay := time.Duration(sim.RollWithLabel(250.0, 1000.0, "Fireseed cast delay")) * time.Millisecond
		treant.WaitUntil(sim, treant.NextGCDAt()+delay)
		return
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
		CritMultiplier:   treant.DefaultCritMultiplier(),
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
