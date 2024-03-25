package warrior

import (
	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterExecuteSpell() {
	var rageMetrics *core.ResourceMetrics
	warrior.Execute = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47471},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   10,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.IsExecutePhase20() && warrior.StanceMatches(BattleStance|BerserkerStance)
		},

		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			availableRage := spell.Unit.CurrentRage()
			extraRage := core.TernaryFloat64(availableRage >= 20, 20, availableRage)
			warrior.SpendRage(sim, extraRage, rageMetrics)
			rageMetrics.Events--

			ap := spell.MeleeAttackPower()
			baseDamage := 10.0 + (ap * 0.437)
			extraDamageScale := extraRage / 20.0
			extraDamage := ((ap * 0.874) - 1.0) * extraDamageScale

			result := spell.CalcAndDealDamage(sim, target, baseDamage+extraDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
	rageMetrics = warrior.Execute.Cost.(*core.RageCost).ResourceMetrics
}
