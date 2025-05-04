package shadow

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/priest"
)

const mbScale = 2.638
const mbCoeff = 1.909
const mbVariance = 0.055

func (shadow *ShadowPriest) registerMindBlastSpell() {
	shadow.MindBlast = shadow.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8092},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: priest.PriestSpellMindBlast,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           shadow.DefaultCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
			PercentModifier: 100,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    shadow.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ThreatMultiplier: 1,
		BonusCoefficient: mbCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shadow.CalcAndRollDamageRange(sim, mbScale, mbVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				shadow.ShadowOrbs.Gain(1, spell.ActionID, sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
