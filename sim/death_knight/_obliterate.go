package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var obliterateActionID = core.ActionID{SpellID: 49020}

func (dk *DeathKnight) registerObliterateSpell() {
	ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       obliterateActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellObliterate,

		DamageMultiplier:         1.5,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.28900000453 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.125)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	hasBloodRites := dk.Inputs.Spec == proto.Spec_SpecBloodDeathKnight

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       obliterateActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellObliterate,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 20,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.5,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.57800000906 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.125)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if hasBloodRites {
				spell.SpendRefundableCostAndConvertFrostOrUnholyRune(sim, result, 1)
			} else {
				spell.SpendRefundableCost(sim, result)
			}
			dk.ThreatOfThassarianProc(sim, result, ohSpell)

			spell.DealDamage(sim, result)
		},
	})
}
