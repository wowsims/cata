package hunter

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerGlaiveTossSpell() {
	if !hunter.Talents.GlaiveToss {
		return
	}

	registerGlaive := func(spellID int32) *core.Spell {
		return hunter.RegisterSpell(core.SpellConfig{
			ActionID:                 core.ActionID{SpellID: spellID},
			SpellSchool:              core.SpellSchoolPhysical,
			ProcMask:                 core.ProcMaskRangedSpecial,
			Flags:                    core.SpellFlagMeleeMetrics,
			MissileSpeed:             18,
			BonusCritPercent:         0,
			DamageMultiplierAdditive: 1,
			DamageMultiplier:         1,
			CritMultiplier:           hunter.DefaultCritMultiplier(),
			ThreatMultiplier:         1,
			BonusCoefficient:         1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				numHits := hunter.Env.GetNumTargets()
				sharedDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower())
				sharedDmg += spell.RangedAttackPower() * 0.2
				sharedDmg += (435.8 + sim.RandomFloat(fmt.Sprintf("Glaive Toss-%v", spellID))*872)
				// Here we assume the Glaive Toss hits every single target in the encounter.
				// This might be unrealistic, but until we have more spatial parameters, this is what we should do.
				results := make([]*core.SpellResult, numHits)
				for i := int32(0); i < numHits; i++ {
					unit := hunter.Env.GetTargetUnit(i)
					dmg := sharedDmg
					if unit == target {
						dmg *= 4 // primary target takes 4× damage
					}
					results[i] = spell.CalcDamage(sim, unit, dmg, spell.OutcomeRangedHitAndCrit)
				}

				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					for _, res := range results {
						spell.DealDamage(sim, res)
					}
				})
			},
		})
	}

	firstGlaive := registerGlaive(120755)
	secondGlaive := registerGlaive(120756)

	hunter.GlaiveToss = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 117050},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskProc,
		Flags:       core.SpellFlagAPL,
		MaxRange:    40,
		MinRange:    0,
		FocusCost: core.FocusCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: 15 * time.Second,
			},
		},
		DamageMultiplierAdditive: 1,

		CritMultiplier:   hunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			firstGlaive.Cast(sim, target)
			secondGlaive.Cast(sim, target)
		},
	})
}
