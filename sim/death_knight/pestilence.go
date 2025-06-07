package death_knight

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var PestilenceActionID = core.ActionID{SpellID: 50842}

func (dk *DeathKnight) registerPestilenceSpell() {
	hasReaping := dk.Inputs.Spec == proto.Spec_SpecUnholyDeathKnight

	dk.PestilenceSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 50842},
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellPestilence,

		MaxRange: core.MaxMeleeRange,

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
						spell.SpendRefundableCostAndConvertBloodRune(sim, result, 1)
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
