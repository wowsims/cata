package mage

import (
	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerIceLanceSpell() {
	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30455},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellIceLance,
		MissileSpeed:   38,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 6,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultCritMultiplier(),
		BonusCoefficient:         0.378,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.432 * mage.ClassSpellScaling
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
