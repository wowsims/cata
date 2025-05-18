package destruction

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const felFlameVariance = 0.1
const felFlameScale = 0.85
const felFlameCoeff = 0.85

func (destruction DestructionWarlock) registerFelflame() {
	destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 77799},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellFelFlame,
		MissileSpeed:   38,
		ManaCost:       core.ManaCostOptions{BaseCostPercent: 3},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		DamageMultiplier: 1.0,
		CritMultiplier:   destruction.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: felFlameCoeff,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destruction.CalcAndRollDamageRange(sim, felFlameScale, felFlameVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			destruction.BurningEmbers.Gain(core.TernaryInt32(result.DidCrit(), 2, 1), spell.ActionID, sim)
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
