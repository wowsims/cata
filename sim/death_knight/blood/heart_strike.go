package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var HeartStrikeActionID = core.ActionID{SpellID: 55050}

/*
Instantly strike the target and up to two additional nearby enemies, causing 105% weapon damage plus 545 on the primary target, with each additional enemy struck taking 50% less damage than the previous target.
Damage dealt to each target is increased by an additional 15% for each of your diseases present.
*/
func (bdk *BloodDeathKnight) registerHeartStrike() {
	numHits := min(3, bdk.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)

	bdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       HeartStrikeActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellHeartStrike,

		MaxRange: core.MaxMeleeRange,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 1.05,
		CritMultiplier:   bdk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bdk.CalcScalingSpellDmg(0.43700000644) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			defaultMultiplier := spell.DamageMultiplier
			currentTarget := target
			for idx := int32(0); idx < numHits; idx++ {
				targetDamage := baseDamage * bdk.GetDiseaseMulti(currentTarget, 1.0, 0.15)

				results[idx] = spell.CalcDamage(sim, currentTarget, targetDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				if idx == 0 {
					spell.SpendRefundableCost(sim, results[idx])
				}

				spell.DamageMultiplier *= 0.5
				currentTarget = bdk.Env.NextTargetUnit(currentTarget)
			}

			spell.DamageMultiplier = defaultMultiplier

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}

func (bdk *BloodDeathKnight) registerDrwHeartStrike() *core.Spell {
	numHits := min(3, bdk.Env.GetNumTargets())
	results := make([]*core.SpellResult, numHits)
	return bdk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    HeartStrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bdk.CalcScalingSpellDmg(0.43700000644) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			defaultMultiplier := spell.DamageMultiplier
			currentTarget := target
			for idx := int32(0); idx < numHits; idx++ {
				targetDamage := baseDamage * bdk.RuneWeapon.GetDiseaseMulti(currentTarget, 1.0, 0.15)

				results[idx] = spell.CalcDamage(sim, currentTarget, targetDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				spell.DamageMultiplier *= 0.5
				currentTarget = bdk.Env.NextTargetUnit(currentTarget)
			}

			spell.DamageMultiplier = defaultMultiplier

			for _, result := range results {
				spell.DealDamage(sim, result)
				spell.DamageMultiplier /= 0.5
			}
		},
	})
}
