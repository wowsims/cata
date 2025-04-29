package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
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

		ManaCost: core.ManaCostOptions{BaseCostPercent: 7},
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
		BonusCoefficient:         0.62800002098,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destro.CalcAndRollDamageRange(sim, 1.54700005054, 0.23800000548)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
