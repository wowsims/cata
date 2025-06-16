package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const carrionSwarmScale = 0.5
const carrionSwarmVariance = 0.1
const carrionSwarmCoeff = 0.5

func (demonology *DemonologyWarlock) registerCarrionSwarm() {
	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 103967},
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellDamage,
		SpellSchool:    core.SpellSchoolShadow,
		ClassSpellMask: warlock.WarlockSpellCarrionSwarm,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCDMin: time.Millisecond * 500,
				GCD:    time.Millisecond * 1000,
			},
			CD: core.Cooldown{
				Timer:    demonology.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: carrionSwarmCoeff,
		CritMultiplier:   demonology.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return demonology.IsInMeta() && demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 35, 50))
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			demonology.DemonicFury.Spend(sim, core.TernaryInt32(demonology.T15_2pc.IsActive(), 35, 50), spell.ActionID)
			for _, enemy := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(
					sim,
					enemy,
					demonology.CalcAndRollDamageRange(sim, carrionSwarmScale, carrionSwarmVariance),
					spell.OutcomeMagicHitAndCrit,
				)
			}
		},
	})
}
