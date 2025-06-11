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
	rpMetric := dk.NewRunicPowerMetrics(BloodBoilActionID)
	hasGlyphOfFesteringBlood := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfFesteringBlood)
	results := make([]*core.SpellResult, dk.Env.GetNumTargets())
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       BloodBoilActionID,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellBloodBoil,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost: 1,
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

			// TODO: Check if this still happens with Conversion talent active
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
				baseDamage *= core.TernaryFloat64(dk.RuneWeapon.DiseasesAreActive(aoeTarget), 1.5, 1.0)

				results[idx] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
