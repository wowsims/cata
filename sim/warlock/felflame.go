package warlock

import "github.com/wowsims/mop/sim/core"

const felFlameVariance = 0.1
const felFlameScale = 0.85
const felFlameCoeff = 0.85

func (warlock *Warlock) RegisterFelflame(callback WarlockSpellCastedCallback) *core.Spell {
	return warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 77799},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellFelFlame,
		MissileSpeed:   38,
		ManaCost:       core.ManaCostOptions{BaseCostPercent: 6},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		DamageMultiplier: 1.0,
		CritMultiplier:   warlock.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: felFlameCoeff,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, felFlameScale, felFlameVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if callback != nil {
				callback([]core.SpellResult{*result}, spell, sim)
			}

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
