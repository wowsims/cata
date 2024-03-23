package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerMultiShotSpell() {

	hunter.MultiShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 2643},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		BonusCritRating:          0,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1.2,
		CritMultiplier:           hunter.CritMultiplier(true, false, false),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := hunter.Env.GetNumTargets() // Multi is uncapped in Cata

			sharedDmg := hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				spell.BonusWeaponDamage() //

			baseDamageArray := make([]float64, numHits)
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				currentTarget := hunter.Env.GetTargetUnit(hitIndex)
				baseDamage := sharedDmg + 0.2*spell.RangedAttackPower(currentTarget)
				baseDamageArray[hitIndex] = baseDamage

			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.CalcAndDealDamage(sim, curTarget, baseDamageArray[hitIndex], spell.OutcomeRangedHitAndCrit)

					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			})

		},
	})
}
