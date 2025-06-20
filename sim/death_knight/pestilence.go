package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var PestilenceActionID = core.ActionID{SpellID: 50842}

// Spreads existing Blood Plague and Frost Fever infections from your target to all other enemies within 10 yards.
func (dk *DeathKnight) registerPestilence() {
	hasReaping := dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight
	maxRange := core.MaxMeleeRange + core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfPestilence), 5, 0)

	dk.PestilenceSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       PestilenceActionID,
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellPestilence,

		MaxRange: maxRange,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 0,
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			frostFeverActive := dk.FrostFeverSpell.Dot(target).IsActive()
			bloodPlagueActive := dk.BloodPlagueSpell.Dot(target).IsActive()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealOutcome(sim, aoeTarget, spell.OutcomeMagicHit)

				if aoeTarget == target {
					if hasReaping {
						spell.SpendRefundableCostAndConvertBloodRune(sim, result.Landed())
					} else {
						spell.SpendRefundableCost(sim, result)
					}
				}

				if result.Landed() {
					if aoeTarget != target {
						if frostFeverActive {
							dk.FrostFeverSpell.Cast(sim, aoeTarget)
						}
						if bloodPlagueActive {
							dk.BloodPlagueSpell.Cast(sim, aoeTarget)
						}
					}
				}
			}
		},
	})
}

func (dk *DeathKnight) registerDrwPestilence() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50842},
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		MaxRange: core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			frostFeverActive := dk.RuneWeapon.FrostFeverSpell.Dot(target).IsActive()
			bloodPlagueActive := dk.RuneWeapon.BloodPlagueSpell.Dot(target).IsActive()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealOutcome(sim, aoeTarget, spell.OutcomeMagicHit)

				if result.Landed() {
					if aoeTarget != target {
						if frostFeverActive {
							dk.RuneWeapon.FrostFeverSpell.Cast(sim, aoeTarget)
						}
						if bloodPlagueActive {
							dk.RuneWeapon.BloodPlagueSpell.Cast(sim, aoeTarget)
						}
					}
				}
			}
		},
	})
}
