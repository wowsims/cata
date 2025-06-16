package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const maleficGraspScale = 0.132 * 1.5 // 2025.06.13 Changes to Beta - Malefic Damage increased by 50%
const maleficGraspCoeff = 0.132 * 1.5

func (affliction *AfflictionWarlock) registerMaleficGrasp() {
	affliction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 103103},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagChanneled,
		ClassSpellMask: warlock.WarlockSpellMaleficGrasp,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 1.5},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           affliction.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura:                 core.Aura{Label: "MaleficGrasp"},
			NumberOfTicks:        4,
			TickLength:           1 * time.Second,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			BonusCoefficient:     maleficGraspCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, affliction.CalcScalingSpellDmg(maleficGraspScale))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				if !result.Landed() {
					return
				}

				// 2025.06.13 Changes to Beta - Malefic DoT component increased to 50%
				affliction.ProcMaleficEffect(target, 0.5, sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
