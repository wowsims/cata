package warrior

import (
	"math"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterExecuteSpell() {
	minRageAfterExecute := 5.0 * float64(warrior.Talents.SuddenDeath)

	var rageMetrics *core.ResourceMetrics
	warrior.Execute = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5308},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskExecute | SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

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

		CritMultiplier:   warrior.DefaultCritMultiplier(),
		DamageMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			availableRage := spell.Unit.CurrentRage()
			extraRage := math.Min(availableRage, 20)

			// Sudden Death: Keep X rage after using execute
			// Adjust extra rage spend to hit this floor

			spell.Unit.SpendRage(sim, extraRage, rageMetrics)
			rageMetrics.Events--

			rageAfterSpend := spell.Unit.CurrentRage()
			if rageAfterSpend < minRageAfterExecute {
				spell.Unit.AddRage(sim, minRageAfterExecute-rageAfterSpend, rageMetrics)
				rageMetrics.Events--
			}

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
	rageMetrics = warrior.Execute.RageMetrics()
}
