package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var frostStrikeActionID = core.ActionID{SpellID: 49143}

// Instantly strike the enemy, causing 115% weapon damage as Frost damage.
func (fdk *FrostDeathKnight) registerFrostStrikeSpell() {
	ohSpell := fdk.RegisterSpell(core.SpellConfig{
		ActionID:       frostStrikeActionID.WithTag(2), // Actually 66196
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: death_knight.DeathKnightSpellFrostStrike,

		DamageMultiplier: 1.3,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.12399999797 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	fdk.RegisterSpell(core.SpellConfig{
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

		DamageMultiplier: 1.3,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.24699999392 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)

			if result.Landed() && dk.ThreatOfThassarianAura.IsActive() {
				ohSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
