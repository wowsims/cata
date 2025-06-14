package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

var bafIncinerateScale = 1.568
var bafIncinerateCoeff = 1.568

func (destruction *DestructionWarlock) registerFireAndBrimstoneIncinerate() {
	destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 114654},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		MissileSpeed:   24,
		ClassSpellMask: warlock.WarlockSpellFaBIncinerate,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 5},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2000 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           destruction.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         bafIncinerateCoeff,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destruction.BurningEmbers.CanSpend(10)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !destruction.FABAura.IsActive() {
				destruction.FABAura.Activate(sim)
			}

			reduction := destruction.getFABReduction()
			spell.DamageMultiplier *= reduction
			destruction.BurningEmbers.Spend(sim, 10, spell.ActionID)
			for _, enemy := range sim.Encounter.TargetUnits {
				baseDamage := destruction.CalcAndRollDamageRange(sim, bafIncinerateScale, incinerateVariance)
				result := spell.CalcDamage(sim, enemy, baseDamage, spell.OutcomeMagicHitAndCrit)
				var emberGain int32 = 1
				if destruction.T15_4pc.IsActive() && sim.Proc(0.08, "T15 4p") {
					emberGain += 1
				}

				// ember lottery
				if sim.Proc(0.15, "Ember Lottery") {
					emberGain *= 2
				}

				if result.DidCrit() {
					emberGain += 1
				}

				destruction.BurningEmbers.Gain(sim, emberGain, spell.ActionID)

				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}

			spell.DamageMultiplier /= reduction
		},
	})
}
