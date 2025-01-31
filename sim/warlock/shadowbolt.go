package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerShadowBolt() {
	shadowAndFlameProcChance := []float64{0.0, 0.33, 0.66, 1.0}[warlock.Talents.ShadowAndFlame]
	shadowAndFlameAuras := warlock.NewEnemyAuraArray(core.ShadowAndFlameAura)

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 686},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellShadowBolt,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{BaseCost: 0.10},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 3 * time.Second,
			},
		},

		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.75400000811,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, 0.62000000477, 0.1099999994)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() && sim.Proc(shadowAndFlameProcChance, "S&F Proc") {
					shadowAndFlameAuras.Get(result.Target).Activate(sim)
				}
			})
		},
	})
}
