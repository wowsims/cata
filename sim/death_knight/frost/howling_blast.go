package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var HowlingBlastActionID = core.ActionID{SpellID: 49184}

// Blast the target with a frigid wind, dealing (<mastery> * (573 + 0.848 * <AP>)) Frost damage to that foe, and (0.5 * <mastery> * (573 + 0.848 * <AP>)) Frost damage to all other enemies within 10 yards, infecting all targets with Frost Fever.
func (fdk *FrostDeathKnight) registerHowlingBlast() {
	results := make([]*core.SpellResult, fdk.Env.GetNumTargets())

	fdk.RegisterSpell(core.SpellConfig{
		ActionID:       HowlingBlastActionID,
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellHowlingBlast,

		MaxRange: 30,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fdk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := fdk.CalcScalingSpellDmg(0.46000000834) + 0.848*spell.MeleeAttackPower()
				damageMultiplier := spell.DamageMultiplier

				if aoeTarget != target {
					// Beta changes 2025-06-16: https://www.wowhead.com/mop-classic/news/blood-death-knights-buffed-and-even-more-class-balance-adjustments-mists-of-377292
					// - Howling Blast’s damage to targets around the primary target has been increased to 65% of the damage dealt (was 50%). [5.2 Revert]
					spell.DamageMultiplier *= 0.65
				}

				results[idx] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				spell.DamageMultiplier = damageMultiplier

				if aoeTarget == target {
					spell.SpendRefundableCost(sim, results[idx])
				}
			}

			for _, result := range results {
				spell.DealDamage(sim, result)

				if result.Landed() {
					fdk.FrostFeverSpell.Cast(sim, result.Target)
				}
			}
		},
	})
}
