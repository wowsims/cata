package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warlock"
)

const soulfireScale = 0.854
const soulfireCoeff = 0.854
const soulfireVariance = 0.2

func (demonology *DemonologyWarlock) registerSoulfire() {
	demonology.Soulfire = demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6353},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellSoulFire,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 15,
			PercentModifier: 1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 4 * time.Second,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           demonology.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         soulfireCoeff,
		BonusCritPercent:         100,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := demonology.CalcAndRollDamageRange(sim, soulfireScale, soulfireVariance)

			// Damage is increased by crit chance
			spell.DamageMultiplier *= (1 + demonology.GetStat(stats.SpellCritPercent)/100)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier /= (1 + demonology.GetStat(stats.SpellCritPercent)/100)

			if !demonology.IsInMeta() {
				demonology.DemonicFury.Gain(sim, 30, spell.ActionID)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
