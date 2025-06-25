package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerScorchSpell() {

	scorchVariance := 0.17   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Variance"
	scorchScaling := .98     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Coefficient"
	scorchCoefficient := .84 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "BonusCoefficient"

	fire.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2948},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellScorch,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 2,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           fire.DefaultCritMultiplier(),
		BonusCoefficient:         scorchCoefficient,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := fire.CalcAndRollDamageRange(sim, scorchScaling, scorchVariance)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			fire.HandleHeatingUp(sim, spell, result)
		},
	})
}
