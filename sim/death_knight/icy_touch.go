package death_knight

import (
	"github.com/wowsims/mop/sim/core"
)

var IcyTouchActionID = core.ActionID{SpellID: 45477}

func (dk *DeathKnight) registerIcyTouchSpell() {
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       IcyTouchActionID,
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellIcyTouch,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.46799999475 + spell.MeleeAttackPower()*0.2

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.SpendRefundableCost(sim, result)

			if result.Landed() {
				dk.FrostFeverSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}

func (dk *DeathKnight) registerDrwIcyTouchSpell() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    IcyTouchActionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.46799999475 + spell.MeleeAttackPower()*0.2

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				dk.RuneWeapon.FrostFeverSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
