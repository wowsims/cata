package unholy

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var BloodStrikeActionID = core.ActionID{SpellID: 45902}

/*
Instantly strike the enemy, causing 40% weapon damage plus 942.
Damage is increased by 12.5% for each of your diseases on the target.
*/
func (uhdk *UnholyDeathKnight) registerBloodStrike() {
	uhdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       BloodStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellBloodStrike,

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

		DamageMultiplier: 0.4,
		CritMultiplier:   uhdk.DefaultCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := uhdk.CalcScalingSpellDmg(0.75599998236) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= uhdk.GetDiseaseMulti(target, 1.0, 0.125)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCostAndConvertBloodRune(sim, result.Landed())

			spell.DealDamage(sim, result)
		},
	})
}
