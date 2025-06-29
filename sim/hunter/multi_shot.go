package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerMultiShotSpell() {

	hunter.MultiShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2643},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskRangedSpecial,
		ClassSpellMask: HunterSpellMultiShot,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
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
		CritMultiplier:           hunter.CritMultiplier(true, false, false),
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := hunter.Env.ActiveTargetCount() // Multi is uncapped in Cata

			sharedDmg := hunter.AutoAttacks.Ranged().BaseDamage(sim)

			baseDamageArray := make([]*core.SpellResult, numHits)
			for hitIndex, currentTarget := range sim.Encounter.ActiveTargetUnits {
				baseDamage := sharedDmg + 0.2*spell.RangedAttackPower(currentTarget)
				baseDamageArray[hitIndex] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeRangedHitAndCrit)

			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				for _, result := range baseDamageArray {
					spell.DealDamage(sim, result)
					if hunter.Talents.SerpentSpread > 0 {
						duration := time.Duration(3+(hunter.Talents.SerpentSpread*3)) * time.Second

						ss := hunter.SerpentSting.Dot(result.Target)
						if hunter.Talents.ImprovedSerpentSting > 0 && (!ss.IsActive() || ss.RemainingDuration(sim) <= duration) {
							hunter.ImprovedSerpentSting.Cast(sim, result.Target)
						}
						ss.BaseTickCount = (3 + (hunter.Talents.SerpentSpread * 3)) / 2
						ss.Apply(sim)
					}
				}
			})

		},
	})
}
