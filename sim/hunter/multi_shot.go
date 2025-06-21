package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hunter *Hunter) registerMultiShotSpell() {

	hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2643},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskRangedSpecial,
		ClassSpellMask: HunterSpellMultiShot,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagRanged,
		MissileSpeed:   40,
		MinRange:       0,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		BonusCritPercent:         0,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         0.6,
		CritMultiplier:           hunter.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := hunter.Env.GetNumTargets() // Multi is uncapped in Cata

			sharedDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower())

			baseDamageArray := make([]*core.SpellResult, numHits)
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				currentTarget := hunter.Env.GetTargetUnit(hitIndex)
				baseDamage := sharedDmg
				baseDamageArray[hitIndex] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeRangedHitAndCrit)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.DealDamage(sim, baseDamageArray[hitIndex])

					//Serpent Spread
					if hunter.Spec == proto.Spec_SpecSurvivalHunter {

						ss := hunter.SerpentSting.Dot(curTarget)
						hunter.ImprovedSerpentSting.Cast(sim, curTarget)
						ss.BaseTickCount = 5
						ss.Apply(sim)
					}

					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			})

		},
	})
}
