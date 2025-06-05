package warrior

import (
	"github.com/wowsims/mop/sim/core"
)

func (war *Warrior) registerExecuteSpell() {
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5308},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskExecute,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		CritMultiplier:   war.DefaultCritMultiplier(),
		DamageMultiplier: 1.0,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20() || war.T16Dps4P.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := war.CalcScalingSpellDmg(5.52500009537) + spell.MeleeAttackPower()*2.54999995232
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
