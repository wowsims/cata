package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const (
	WrathBonusCoeff = 1.338
	WrathCoeff      = 2.676
	WrathVariance   = 0.25
)

func (druid *Druid) registerWrathSpell() {
	druid.Wrath = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5176},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellWrath,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 8.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		BonusCoefficient: WrathBonusCoeff,

		DamageMultiplier: 1,

		CritMultiplier: druid.DefaultCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := druid.CalcAndRollDamageRange(sim, WrathCoeff, WrathVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
