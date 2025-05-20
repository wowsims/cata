package mage

import (
	"github.com/wowsims/mop/sim/core"
)

// Values found at https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=30455
var iceLanceScaling = 0.34
var iceLanceCoefficient = 0.34
var iceLanceVariance = 0.25

func (mage *Mage) registerIceLanceSpell() {
	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30455},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellIceLance,
		MissileSpeed:   38,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultCritMultiplier(),
		BonusCoefficient:         iceLanceCoefficient,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := mage.CalcAndRollDamageRange(sim, iceLanceScaling, iceLanceVariance)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
