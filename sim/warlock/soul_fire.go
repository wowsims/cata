package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerSoulFireSpell() {
	warlock.SoulFire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6353},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSoulFire,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 4,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.726,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Deactivate(sim)
			result := spell.CalcDamage(sim, target, 2862, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
