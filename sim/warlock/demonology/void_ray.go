package demonology

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const voidRayScale = 0.525
const voidRayVariance = 0.1
const voidRayCoeff = 0.234

func (demonology *DemonologyWarlock) registerVoidRay() {
	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 115422},
		SpellSchool:    core.SpellSchoolShadowFlame,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellVoidray,
		MissileSpeed:   38,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		DamageMultiplier: 1.0,
		CritMultiplier:   demonology.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: voidRayCoeff,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return demonology.IsInMeta() && demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 56, 80))
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			demonology.DemonicFury.Spend(sim, core.TernaryInt32(demonology.T15_2pc.IsActive(), 56, 80), spell.ActionID)
			for _, enemy := range sim.Encounter.TargetUnits {
				baseDamage := demonology.CalcAndRollDamageRange(sim, voidRayScale, voidRayVariance)
				spell.CalcAndDealDamage(sim, enemy, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
