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
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   24,
		ClassSpellMask: warlock.WarlockSpellIncinerate,

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
			return destruction.BurningEmbers.CanSpend(10) && destruction.FABAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			reduction := destruction.getFABReduction()
			spell.DamageMultiplier *= reduction
			destruction.BurningEmbers.Spend(10, spell.ActionID, sim)
			for _, enemy := range sim.Encounter.TargetUnits {
				baseDamage := destruction.CalcAndRollDamageRange(sim, bafIncinerateScale, incinerateVariance)
				result := spell.CalcDamage(sim, enemy, baseDamage, spell.OutcomeMagicHitAndCrit)
				if result.DidCrit() {
					destruction.BurningEmbers.Gain(2, spell.ActionID, sim)
				} else {
					destruction.BurningEmbers.Gain(1, spell.ActionID, sim)
				}

				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}

			spell.DamageMultiplier /= reduction
		},
	})
}
