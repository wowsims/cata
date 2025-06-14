package hunter

import (
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

				sharedDmg := spell.RangedAttackPower() * 0.2
				sharedDmg += hunter.CalcAndRollDamageRange(sim, 0.69999998808, 1)
				// Here we assume the Glaive Toss hits every single target in the encounter.
				// This might be unrealistic, but until we have more spatial parameters, this is what we should do.
				results := make([]*core.SpellResult, numHits)
				for i := int32(0); i < numHits; i++ {
					unit := hunter.Env.GetTargetUnit(i)
					// Primary always hits, secondaries only on a successful roll
					if unit != target {
						successChance := hunter.Options.GlaiveTossSuccess / 100.0
						roll := sim.RollWithLabel(0, 1, "GlaiveTossSuccess")
						if roll > successChance {
							continue
						}
					}
					dmg := sharedDmg
					if unit == target {
						dmg *= 4 // primary target takes 4× damage
					}
					results[i] = spell.CalcDamage(sim, unit, dmg, spell.OutcomeRangedHitAndCrit)
				}

				for _, res := range results {
					spell.DealDamage(sim, res)
				}
			},
		})
	}

	firstGlaive := registerGlaive(120755)
	secondGlaive := registerGlaive(120756)

	hunter.GlaiveToss = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 117050},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskProc,
		Flags:        core.SpellFlagAPL,
		MaxRange:     40,
		MissileSpeed: 15,
		MinRange:     0,
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
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				firstGlaive.Cast(sim, target)
				secondGlaive.Cast(sim, target)
			})
		},
	})
}
