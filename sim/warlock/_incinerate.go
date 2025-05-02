package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) registerIncinerate() {
	shadowAndFlameProcChance := []float64{0.0, 0.33, 0.66, 1.0}[warlock.Talents.ShadowAndFlame]

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 29722},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   24,
		ClassSpellMask: WarlockSpellIncinerate,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 14},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2500 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.53899997473,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, 0.57300001383, 0.15000000596)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if warlock.Immolate.Dot(target).IsActive() {
				result.Damage *= 7.0 / 6.0
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() && sim.Proc(shadowAndFlameProcChance, "S&F Proc") {
					core.ShadowAndFlameAura(target).Activate(sim)
				}
			})
		},
	})
}
