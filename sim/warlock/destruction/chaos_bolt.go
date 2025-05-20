package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

var chaosBoltVariance = 0.2
var chaosBoltScale = 2.5875
var chaosBoltCoeff = 2.5875

func (destro *DestructionWarlock) registerChaosBolt() {
	destro.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 116858},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellChaosBolt,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 3000 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           destro.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         chaosBoltCoeff,
		BonusCritPercent:         100,
		MissileSpeed:             16,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destro.CalcAndRollDamageRange(sim, chaosBoltScale, chaosBoltVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)
			})

			if result.Landed() {
				destro.BurningEmbers.Spend(10, spell.ActionID, sim)
			}
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destro.BurningEmbers.CanSpend(10)
		},
	})
}
