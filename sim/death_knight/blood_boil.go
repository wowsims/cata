package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var BloodBoilActionID = core.ActionID{SpellID: 48721}

/*
Boils the blood of all enemies within 10 yards, dealing (<3472-4245> + <AP> * 0.11) Shadow damage.
Deals 50% additional damage to targets infected with Blood Plague or Frost Fever.
*/
func (dk *DeathKnight) registerBloodBoil() {
	rpMetric := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 65658})
	hasGlyphOfFesteringBlood := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfFesteringBlood)
	hasReaping := dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight
	results := make([]*core.SpellResult, dk.Env.GetNumTargets())

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       BloodBoilActionID,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellBloodBoil,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
			// Not actually refundable, but setting this to `true` if specced into blood
			// makes the default SpendCost function skip handling the rune cost and
			// lets us manually spend it with death rune conversion in ApplyEffects.
			Refundable: hasReaping,
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
			anyHit := false
			for idx, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := dk.CalcAndRollDamageRange(sim, 3.09599995613, 0.20000000298) +
					0.1099999994*spell.MeleeAttackPower()
				baseDamage *= core.TernaryFloat64(hasGlyphOfFesteringBlood || dk.DiseasesAreActive(aoeTarget), 1.5, 1.0)

				results[idx] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				anyHit = anyHit || results[idx].Landed()
			}

			if hasReaping {
				spell.SpendRefundableCostAndConvertBloodRune(sim, true)
			}

			if anyHit {
				dk.AddRunicPower(sim, 10, rpMetric)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}

func (dk *DeathKnight) registerDrwBloodBoil() *core.Spell {
	results := make([]*core.SpellResult, dk.Env.GetNumTargets())
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    BloodBoilActionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAoE,
		ProcMask:    core.ProcMaskSpellDamage,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for idx, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := dk.CalcAndRollDamageRange(sim, 3.09599995613, 0.20000000298) +
					0.1099999994*spell.MeleeAttackPower()

				// TODO: Is DRW damage affected by Glyph of Festering Blood?
				// TODO: Verify if owner's diseases count, simc says so
				anyDiseasesActive := dk.DiseasesAreActive(aoeTarget) || dk.RuneWeapon.DiseasesAreActive(aoeTarget)
				baseDamage *= core.TernaryFloat64(anyDiseasesActive, 1.5, 1.0)

				results[idx] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
