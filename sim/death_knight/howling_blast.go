package death_knight

import (
	"github.com/wowsims/cata/sim/core"
)

var HowlingBlastActionID = core.ActionID{SpellID: 49184}

func (dk *DeathKnight) registerHowlingBlastSpell() {
	if !dk.Talents.HowlingBlast {
		return
	}

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       HowlingBlastActionID,
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellHowlingBlast,

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
		ThreatMultiplier: 1,

		CritMultiplier: dk.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, sim.GetNumTargets())
			for idx, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := dk.ClassBaseScaling*1.17499995232 + 0.44*spell.MeleeAttackPower()

				if aoeTarget != target {
					baseDamage *= 0.5
				}
				results[idx] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				if aoeTarget == target {
					spell.SpendRefundableCost(sim, results[idx])
				}
			}

			for _, result := range results {
				spell.DealDamage(sim, result)
			}
		},
	})
}
