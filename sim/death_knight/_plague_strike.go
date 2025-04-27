package death_knight

import (
	"github.com/wowsims/mop/sim/core"
)

var PlagueStrikeActionID = core.ActionID{SpellID: 45462}

func (dk *DeathKnight) registerPlagueStrikeSpell() {
	ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       PlagueStrikeActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellPlagueStrike,

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.18700000644 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       PlagueStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellPlagueStrike,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.37400001287 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)
			dk.ThreatOfThassarianProc(sim, result, ohSpell)
			if result.Landed() {
				dk.BloodPlagueSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (dk *DeathKnight) registerDrwPlagueStrikeSpell() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    PlagueStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.37400001287 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if result.Landed() {
				dk.RuneWeapon.BloodPlagueSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
