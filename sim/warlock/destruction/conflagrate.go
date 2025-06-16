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
		},
		DamageMultiplier: 1.0,
		CritMultiplier:   destruction.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: conflagrateCoeff,
		Charges:          2,
		RechargeTime:     time.Second * 12,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if destruction.FABAura.IsActive() {
				destruction.FABAura.Deactivate(sim)
			}

			// keep charges in sync
			destruction.FABConflagrate.ConsumeCharge(sim)
			baseDamage := destruction.CalcAndRollDamageRange(sim, conflagrateScale, conflagrateVariance)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			var emberGain int32 = 1

			// ember lottery
			if sim.Proc(0.15, "Ember Lottery") {
				emberGain *= 2
			}

			if result.DidCrit() {
				emberGain += 1
			}

			destruction.BurningEmbers.Gain(sim, emberGain, spell.ActionID)
		},
	})
}
