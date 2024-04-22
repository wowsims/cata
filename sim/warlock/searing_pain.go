package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerSearingPainSpell() {
	warlock.SearingPain = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 5676},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
		},

		BonusCoefficient:         0.37799999118,
		DamageMultiplierAdditive: 1 + warlock.GrandFirestoneBonus(),
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, sim.Roll(347, 410), spell.OutcomeMagicHitAndCrit)
		},
	})
}
