package blood

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/death_knight"
)

var HeartStrikeActionID = core.ActionID{SpellID: 55050}

func (dk *BloodDeathKnight) registerHeartStrikeSpell() {
	maxHits := min(3, dk.Env.TotalTargetCount())
	results := make([]*core.SpellResult, maxHits)

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       HeartStrikeActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellHeartStrike,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.75,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.72799998522 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			numHits := min(maxHits, sim.Environment.ActiveTargetCount())
			currentTarget := target
			for idx := range numHits {
				targetDamage := baseDamage * dk.GetDiseaseMulti(currentTarget, 1.0, 0.15)

				results[idx] = spell.CalcDamage(sim, currentTarget, targetDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				if idx == 0 {
					spell.SpendRefundableCost(sim, results[idx])
				}

				spell.DamageMultiplier *= 0.75
				currentTarget = dk.Env.NextActiveTargetUnit(currentTarget)
			}

			for idx := range numHits {
				spell.DealDamage(sim, results[idx])
				spell.DamageMultiplier /= 0.75
			}
		},
	})
}

func (dk *BloodDeathKnight) registerDrwHeartStrikeSpell() *core.Spell {
	maxHits := min(3, dk.Env.TotalTargetCount())
	results := make([]*core.SpellResult, maxHits)
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    HeartStrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.72799998522 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			numHits := min(maxHits, dk.Env.ActiveTargetCount())
			currentTarget := target
			for idx := int32(0); idx < numHits; idx++ {
				targetDamage := baseDamage * dk.RuneWeapon.GetDiseaseMulti(currentTarget, 1.0, 0.15)

				results[idx] = spell.CalcDamage(sim, currentTarget, targetDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				spell.DamageMultiplier *= 0.75
				currentTarget = dk.Env.NextActiveTargetUnit(currentTarget)
			}

			for idx := range numHits {
				spell.DealDamage(sim, results[idx])
				spell.DamageMultiplier /= 0.75
			}
		},
	})
}
