package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const conflagrateScale = 1.725
const conflagrateVariance = 0.1
const conflagrateCoeff = 1.725

func (destruction *DestructionWarlock) registerConflagrate() {
	destruction.Conflagrate = destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 17962},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellConflagrate,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 1},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    destruction.NewTimer(),
				Duration: 12 * time.Second,
			},
		},
		DamageMultiplier: 1.0,
		CritMultiplier:   destruction.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: conflagrateCoeff,
		Charges:          2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destruction.CalcAndRollDamageRange(sim, conflagrateScale, conflagrateVariance)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			destruction.BurningEmbers.Gain(core.TernaryInt32(result.DidCrit(), 2, 1), spell.ActionID, sim)
		},
	})
}
