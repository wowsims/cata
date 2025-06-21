package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var IcyTouchActionID = core.ActionID{SpellID: 45477}

/*
Chills the target for (<560-607> + 0.319 * <AP>) Frost damage

-- Glyph of Icy Touch --

, dispels 1 beneficial Magic effect

-- /Glyph of Icy Touch --

and infects them with Frost Fever, a disease that deals periodic frost damage for 30 sec.
*/
func (dk *DeathKnight) registerIcyTouch() {
	hasReaping := dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       IcyTouchActionID,
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellIcyTouch,

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
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcAndRollDamageRange(sim, 0.46799999475, 0.08299999684) + spell.MeleeAttackPower()*0.319

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if hasReaping {
				spell.SpendRefundableCostAndConvertFrostRune(sim, result.Landed())
			} else {
				spell.SpendRefundableCost(sim, result)
			}

			if result.Landed() {
				dk.FrostFeverSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (dk *DeathKnight) registerDrwIcyTouch() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    IcyTouchActionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcAndRollDamageRange(sim, 0.46799999475, 0.08299999684) + spell.MeleeAttackPower()*0.319

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				dk.RuneWeapon.FrostFeverSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
