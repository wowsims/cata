package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (destruction *DestructionWarlock) registerFireAndBrimstoneConflagrate() {
	destruction.FABConflagrate = destruction.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 108685},
		Flags:            core.SpellFlagAPL,
		SpellSchool:      core.SpellSchoolFire,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ProcMask:         core.ProcMaskSpellDamage,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Charges:          2,
		RechargeTime:     time.Second * 12,
		ClassSpellMask:   warlock.WarlockSpellConflagrate,
		BonusCoefficient: conflagrateCoeff,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destruction.BurningEmbers.CanSpend(10)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !destruction.FABAura.IsActive() {
				destruction.FABAura.Activate(sim)
			}

			// reduce damage for this spell based on mastery
			reduction := destruction.getFABReduction()
			spell.DamageMultiplier *= reduction

			// keep charges in sync
			destruction.Conflagrate.ConsumeCharge(sim)
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(
					sim,
					aoeTarget,
					destruction.CalcAndRollDamageRange(sim, conflagrateScale, conflagrateVariance),
					spell.OutcomeMagicHitAndCrit)
			}
			spell.DamageMultiplier /= reduction
			destruction.BurningEmbers.Spend(10, spell.ActionID, sim)
		},
	})
}
