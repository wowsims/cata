package blood

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (bdk *BloodDeathKnight) registerRuneStrike() {
	bdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 56815},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMH, // Rune Strike triggers white hit procs as well so we give it both masks.
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellRuneStrike,

		MaxRange: core.MaxMeleeRange,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 30,
			Refundable:     true,
			RefundCost:     6,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 2,
		CritMultiplier:   bdk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.MeleeAttackPower()*0.1

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			spell.SpendRefundableCost(sim, result)

			spell.DealDamage(sim, result)
		},
	})
}

// func (bdk *BloodDeathKnight) registerDrwRuneStrikeSpell() *core.Spell {
// 	return bdk.RuneWeapon.RegisterSpell(core.SpellConfig{
// 		ActionID:    core.ActionID{SpellID: 62036},
// 		SpellSchool: core.SpellSchoolPhysical,
// 		ProcMask:    core.ProcMaskMeleeMH,
// 		Flags:       core.SpellFlagMeleeMetrics,

// 		MaxRange: core.MaxMeleeRange,

// 		DamageMultiplier: 2,
// 		CritMultiplier:   bdk.DefaultCritMultiplier(),
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.MeleeAttackPower()*0.1

// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
// 		},
// 	})
// }
