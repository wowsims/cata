package destruction

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warlock"
)

func (destro *DestructionWarlock) registerChaosBolt() {
	if !destro.Talents.ChaosBolt {
		return
	}

	destro.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 50796},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellChaosBolt,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2500 * time.Millisecond,
			},
			CD: core.Cooldown{
				Timer:    destro.NewTimer(),
				Duration: 12 * time.Second,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           destro.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.628,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destro.CalcAndRollDamageRange(sim, warlock.Coefficient_ChaosBolt, warlock.Variance_ChaosBolt)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
