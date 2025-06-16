package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var BloodStrikeActionID = core.ActionID{SpellID: 45902}

/*
Instantly strike the enemy, causing 40% weapon damage plus 942.
Damage is increased by 12.5% for each of your diseases on the target.
*/
func (dk *DeathKnight) registerBloodStrike() {
	hasReaping := dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       BloodStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellBloodStrike,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 0.8,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.75599998236 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.125)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCostAndConvertBloodRune(sim, result.Landed())

			spell.DealDamage(sim, result)
		},
	})
}
