package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var BloodStrikeActionID = core.ActionID{SpellID: 45902}

func (dk *DeathKnight) registerBloodStrikeSpell() {
	ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       BloodStrikeActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellBloodStrike,

		DamageMultiplier:         0.8,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.37799999118 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.025)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

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
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.75599998236 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.125)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hasReaping {
				spell.SpendRefundableCostAndConvertBloodRune(sim, result, 1)
			} else {
				spell.SpendRefundableCost(sim, result)
			}
			dk.ThreatOfThassarianProc(sim, result, ohSpell)

			spell.DealDamage(sim, result)
		},
	})
}

func (dk *DeathKnight) registerDrwBloodStrikeSpell() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    BloodStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.75599998236 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
