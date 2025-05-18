package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const incinerateVariance = 0.1
const incinerateScale = 1.54 * 1.15 // Hotfix
const incinerateCoeff = 1.54 * 1.15

func (destro *DestructionWarlock) registerIncinerate() {
	destro.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 29722},
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
		CritMultiplier:           destro.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         incinerateCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if destro.FABAura.IsActive() {
				destro.FABAura.Deactivate(sim)
			}

			baseDamage := destro.CalcAndRollDamageRange(sim, incinerateScale, incinerateVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.DidCrit() {
				destro.BurningEmbers.Gain(2, spell.ActionID, sim)
			} else {
				destro.BurningEmbers.Gain(1, spell.ActionID, sim)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
