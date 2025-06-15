package death_knight

import (
	"github.com/wowsims/mop/sim/core"
)

var obliterateActionID = core.ActionID{SpellID: 49020}

/*
A brutal instant attack that deals 250% weapon damage.
Total damage is increased by 12.5% for each of your diseases on the target.
*/
func (dk *DeathKnight) registerObliterateSpell() {
	ohSpell := fdk.RegisterSpell(core.SpellConfig{
		ActionID:       obliterateActionID.WithTag(2), // Actually 66198
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellObliterate,

		DamageMultiplier: 2.5,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.28900000453 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.125)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

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

		DamageMultiplier: 2.5,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.57800000906 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.125)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)

			if result.Landed() && dk.ThreatOfThassarianAura.IsActive() {
				ohSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
