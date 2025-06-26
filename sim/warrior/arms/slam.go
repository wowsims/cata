package arms

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) registerSlam() {

	actionID := core.ActionID{SpellID: 1464}

	var sweepingStrikesSlamDamage float64
	sweepingStrikesSlam := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1), // Real SpellID: 146361
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreAttackerModifiers | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: warrior.SpellMaskSweepingSlam,
		MinRange:       2,

		DamageMultiplier: 0.35,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sweepingStrikesSlamDamage, spell.OutcomeAlwaysHit)
		},
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: warrior.SpellMaskSlam,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   25,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 2.75,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			csDamageMultiplier := core.TernaryFloat64(war.ColossusSmashAuras.Get(target).IsActive(), 1.1, 1.0)
			baseDamage := war.CalcScalingSpellDmg(1) + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())*csDamageMultiplier
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if war.SweepingStrikesAura.IsActive() {
				sweepingStrikesSlamDamage = result.Damage
				for _, otherTarget := range sim.Encounter.TargetUnits {
					if otherTarget != target {
						sweepingStrikesSlam.Cast(sim, otherTarget)
					}
				}
			}

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
