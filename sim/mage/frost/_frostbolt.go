package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

var frostboltVariance    = 0.24 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A116 Field: "Variance"
var frostboltScale       = 1.5 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A116 Field: "Coefficient"
var frostboltCoefficient = 1.5 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A116 Field: "BonusCoefficient"

func (frostMage *FrostMage) registerFrostboltSpell() {
	frostMage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 116},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostbolt,
		MissileSpeed:   28,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           frostMage.DefaultCritMultiplier(),
		BonusCoefficient:         frostboltCoefficient
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := frostMage.CalcAndRollDamageRange(sim, frostboltScale, frostboltVariance)
			
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
