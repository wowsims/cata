package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var frostStrikeActionID = core.ActionID{SpellID: 49143}

func (dk *FrostDeathKnight) registerFrostStrikeSpell() {
	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       frostStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellFrostStrike,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
			Refundable:     true,
			RefundCost:     4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier:         1.3,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           dk.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.24699999392 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)
			dk.ThreatOfThassarianProc(sim, result, ohSpell)

			spell.DealDamage(sim, result)
		},
	})
}
