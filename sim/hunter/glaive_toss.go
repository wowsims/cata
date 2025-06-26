package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerGlaiveTossSpell() {
	if !hunter.Talents.GlaiveToss {
		return
	}

	registerGlaive := func(spellID int32, tag int32) *core.Spell {
		return hunter.RegisterSpell(core.SpellConfig{
			ActionID:                 core.ActionID{SpellID: 117050}.WithTag(tag),
			SpellSchool:              core.SpellSchoolPhysical,
			ProcMask:                 core.ProcMaskRangedSpecial,
			ClassSpellMask:           HunterSpellGlaiveToss,
			Flags:                    core.SpellFlagMeleeMetrics | core.SpellFlagRanged | core.SpellFlagPassiveSpell,
			MissileSpeed:             18,
			BonusCritPercent:         0,
			DamageMultiplierAdditive: 1,
			DamageMultiplier:         1,
			CritMultiplier:           hunter.DefaultCritMultiplier(),
			ThreatMultiplier:         1,
			BonusCoefficient:         1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				numTargets := hunter.Env.GetNumTargets()
				sharedDmg := spell.RangedAttackPower()*0.2 + hunter.CalcAndRollDamageRange(sim, 0.7, 1)
				successChance := hunter.Options.GlaiveTossSuccess / 100.0

				runPass := func(skipPrimary bool) {
					for i := int32(0); i < numTargets; i++ {
						unit := hunter.Env.GetTargetUnit(i)

						// skip the main target on return
						if skipPrimary && unit == target {
							continue
						}

						if unit != target {
							if sim.RollWithLabel(0, 1, "GlaiveTossSuccess") > successChance {
								continue
							}
						}

						dmg := sharedDmg

						// If primary target we add multiplier
						if unit == target {
							dmg *= 4.0
						}

						res := spell.CalcDamage(sim, unit, dmg, spell.OutcomeRangedHitAndCrit)
						spell.DealDamage(sim, res)
					}
				}

				// Glaive Toss does two damage passes
				runPass(false)
				runPass(true) // Return pass, skips primary target
			},
		})
	}

	firstGlaive := registerGlaive(120755, 1)
	secondGlaive := registerGlaive(120756, 2)

	hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 117050},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskProc,
		ClassSpellMask: HunterSpellGlaiveToss,
		Flags:          core.SpellFlagAPL,
		MaxRange:       40,
		MissileSpeed:   15,
		MinRange:       0,
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
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHitNoHitCounter)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				firstGlaive.Cast(sim, target)
				secondGlaive.Cast(sim, target)
			})
		},
	})
}
