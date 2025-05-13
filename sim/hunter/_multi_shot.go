package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerMultiShotSpell() {

	hunter.MultiShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2643},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskRangedSpecial,
		ClassSpellMask: HunterSpellMultiShot,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       5,
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
		DamageMultiplier:         1.21,
		CritMultiplier:           hunter.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := hunter.Env.GetNumTargets() // Multi is uncapped in Cata

			sharedDmg := hunter.AutoAttacks.Ranged().BaseDamage(sim)

			baseDamageArray := make([]*core.SpellResult, numHits)
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				currentTarget := hunter.Env.GetTargetUnit(hitIndex)
				baseDamage := sharedDmg + 0.2*spell.RangedAttackPower()
				baseDamageArray[hitIndex] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeRangedHitAndCrit)

			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.DealDamage(sim, baseDamageArray[hitIndex])
					if hunter.Talents.SerpentSpread > 0 {
						duration := time.Duration(3+(hunter.Talents.SerpentSpread*3)) * time.Second

						ss := hunter.SerpentSting.Dot(curTarget)
						if hunter.Talents.ImprovedSerpentSting > 0 && (!ss.IsActive() || ss.RemainingDuration(sim) <= duration) {
							hunter.ImprovedSerpentSting.Cast(sim, curTarget)
						}
						ss.BaseTickCount = (3 + (hunter.Talents.SerpentSpread * 3)) / 2
						ss.Apply(sim)
					}
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			})

		},
	})
}
