package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var PlagueStrikeActionID = core.ActionID{SpellID: 45462}

/*
A vicious strike that deals 100% weapon damage plus 466 and infects the target with Blood Plague, a disease dealing Shadow damage over time

-- Ebon Plaguebringer --

and Frost Fever, a disease dealing Frost damage over time

-- /Ebon Plaguebringer --

.
*/
func (dk *DeathKnight) registerPlagueStrike() {
	var ohSpell *core.Spell
	if dk.Spec == proto.Spec_SpecFrostDeathKnight {
		ohSpell = dk.registerOffHandPlagueStrike()
	}

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       PlagueStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellPlagueStrike,

		MaxRange: core.MaxMeleeRange,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.37400001287) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)

			if result.Landed() {
				if dk.ThreatOfThassarianAura.IsActive() {
					ohSpell.Cast(sim, target)
				}

				dk.BloodPlagueSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (dk *DeathKnight) registerOffHandPlagueStrike() *core.Spell {
	return dk.RegisterSpell(core.SpellConfig{
		ActionID:       PlagueStrikeActionID.WithTag(2), // Actually 66216
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellPlagueStrike,

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.18700000644) +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}

func (dk *DeathKnight) registerDrwPlagueStrike() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    PlagueStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.37400001287) +
				dk.RuneWeapon.StrikeWeapon.CalculateWeaponDamage(sim, spell.MeleeAttackPower()) +
				dk.RuneWeapon.StrikeWeaponDamage

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if result.Landed() {
				dk.RuneWeapon.BloodPlagueSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
