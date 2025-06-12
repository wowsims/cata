package fire

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerMastery() {

	fire.ignite = fire.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 413843},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellIgnite,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Ignite",
			},
			NumberOfTicks: 2,
			TickLength:    2 * time.Second,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 1)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			dot := spell.Dot(target)
			return spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

}

func (fire *FireMage) ApplyIgnite(sim *core.Simulation, target *core.Unit, damage float64) {
	existingDot := fire.ignite.Dot(target)

	masteryAdjustedDamage := damage * fire.GetMasteryBonus()

	if existingDot.IsActive() {
		masteryAdjustedDamage += float64(existingDot.RemainingTicks()) * existingDot.Spell.ExpectedTickDamage(sim, target)
	}

	fmt.Println(masteryAdjustedDamage)
	fire.ignite.DamageMultiplier *= masteryAdjustedDamage
	fire.ignite.Dot(target).Apply(sim)
	fire.ignite.DamageMultiplier /= masteryAdjustedDamage
}

func (fire *FireMage) GetMasteryBonus() float64 {
	return (.12 + 0.015*fire.GetMasteryPoints())
}
