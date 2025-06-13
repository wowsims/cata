package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const shadowBoltScale = 1.38
const shadowBoltCoeff = 1.38

func (demonology *DemonologyWarlock) registerShadowBolt() {
	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 686},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellShadowBolt,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 5.5},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2500 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           demonology.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         shadowBoltCoeff,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !demonology.IsInMeta()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, demonology.CalcScalingSpellDmg(shadowBoltScale), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})

			if result.Landed() {
				demonology.DemonicFury.Gain(sim, 25, core.ActionID{SpellID: 686})
			}
		},
	})
}
