package warlock

import (
	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerFelFlameSpell() {
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 77799},
		SpellSchool:    core.SpellSchoolShadow | core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellFelFlame,
		MissileSpeed:   38,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.30199998617,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, 0.24799999595, 0.15000000596)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					if warlock.ImmolateDot != nil {
						immoDot := warlock.ImmolateDot.Dot(target)
						if immoDot.IsActive() {
							immoDot.TickCount = max(0, immoDot.TickCount-2)
							immoDot.TakeSnapshot(sim, false)
						}
					}

					if warlock.UnstableAffliction != nil {
						unstableAff := warlock.UnstableAffliction.Dot(target)
						if unstableAff != nil && unstableAff.IsActive() {
							unstableAff.TickCount = max(0, unstableAff.TickCount-2)
							unstableAff.TakeSnapshot(sim, false)
						}
					}
				}
			})
		},
	})
}
