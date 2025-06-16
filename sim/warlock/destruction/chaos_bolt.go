package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warlock"
)

var chaosBoltVariance = 0.2
var chaosBoltScale = 2.5875
var chaosBoltCoeff = 2.5875
var chaosBoltDotCoeff = 0.1294
var chaosBoltDotScale = 0.1294

func (destro *DestructionWarlock) registerChaosBolt() {
	destro.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 116858},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellChaosBolt,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 3000 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           destro.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         chaosBoltCoeff,
		BonusCritPercent:         100,
		MissileSpeed:             16,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Chaosbolt (DoT)",
			},
			NumberOfTicks:    3,
			TickLength:       time.Second,
			BonusCoefficient: chaosBoltDotCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, destro.CalcScalingSpellDmg(chaosBoltDotScale))
				dot.SnapshotAttackerMultiplier *= (1 + destro.GetStat(stats.SpellCritPercent)/100)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickMagicCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := destro.CalcAndRollDamageRange(sim, chaosBoltScale, chaosBoltVariance)
			spell.DamageMultiplier *= (1 + destro.GetStat(stats.SpellCritPercent)/100)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DamageMultiplier /= (1 + destro.GetStat(stats.SpellCritPercent)/100)

			// check again we can actually spend as Dark Soul might have run out before the cast finishes
			if result.Landed() && destro.BurningEmbers.CanSpend(core.TernaryInt32(destro.T15_2pc.IsActive(), 8, 10)) {
				destro.BurningEmbers.Spend(sim, core.TernaryInt32(destro.T15_2pc.IsActive(), 8, 10), spell.ActionID)
			} else {
				return
			}

			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destro.BurningEmbers.CanSpend(core.TernaryInt32(destro.T15_2pc.IsActive(), 8, 10))
		},
	})
}
