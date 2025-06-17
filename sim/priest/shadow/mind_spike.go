package shadow

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

const mindSpikeCoeff = 1.304
const mindSpikeScale = 1.277
const mindSpikeVariance = 0.054

func (shadow *ShadowPriest) registerMindSpike() {
	shadow.MindSpike = shadow.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 73510},
		SpellSchool:              core.SpellSchoolShadow,
		ProcMask:                 core.ProcMaskSpellDamage,
		Flags:                    core.SpellFlagAPL,
		ClassSpellMask:           priest.PriestSpellMindSpike,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           shadow.DefaultCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		ThreatMultiplier: 1,
		BonusCoefficient: mindSpikeCoeff,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(
				sim,
				target,
				shadow.CalcAndRollDamageRange(sim, mindSpikeScale, mindSpikeVariance),
				spell.OutcomeMagicHitAndCrit,
			)
			if result.Landed() {
				if shadow.SurgeOfDarkness == nil || !shadow.SurgeOfDarkness.IsActive() {
					shadow.ShadowWordPain.Dot(target).Deactivate(sim)

					// only access those if they're actually registered and talented
					if shadow.VampiricTouch != nil {
						shadow.VampiricTouch.Dot(target).Deactivate(sim)
					}
					if shadow.DevouringPlague != nil {
						shadow.DevouringPlague.Dot(target).Deactivate(sim)
					}
				}
			}

			// delay hit for dummy effect of SurgeDarkness so aura is active
			spell.DealDamage(sim, result)
		},
	})
}
