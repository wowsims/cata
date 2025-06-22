package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var ObliterateActionID = core.ActionID{SpellID: 49020}

/*
A brutal instant attack that deals 250% weapon damage.
Total damage is increased by 12.5% for each of your diseases on the target.
*/
func (fdk *FrostDeathKnight) registerObliterate() {
	ohSpell := fdk.RegisterSpell(core.SpellConfig{
		ActionID:       ObliterateActionID.WithTag(2), // Actually 66198
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: death_knight.DeathKnightSpellObliterate,

		DamageMultiplier: 2.5,
		CritMultiplier:   fdk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= fdk.GetDiseaseMulti(target, 1.0, 0.125)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	fdk.RegisterSpell(core.SpellConfig{
		ActionID:       ObliterateActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellObliterate,

		MaxRange: core.MaxMeleeRange,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 20,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 2.5,
		CritMultiplier:   fdk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= fdk.GetDiseaseMulti(target, 1.0, 0.125)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)

			if result.Landed() && fdk.ThreatOfThassarianAura.IsActive() {
				ohSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
