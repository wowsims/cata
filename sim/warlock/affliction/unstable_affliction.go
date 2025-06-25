package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const uaCoeff = 0.29
const uaScale = 0.29

func (affliction *AfflictionWarlock) registerUnstableAffliction() {
	affliction.UnstableAffliction = affliction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30108},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellUnstableAffliction,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 1.5},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           affliction.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura:                core.Aura{Label: "UnstableAffliction"},
			NumberOfTicks:       7,
			TickLength:          2 * time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    uaCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, affliction.CalcScalingSpellDmg(uaScale))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				affliction.ApplyDotWithPandemic(spell.Dot(target), sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			dot := spell.Dot(target)
			if useSnapshot {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
				result.Damage /= dot.TickPeriod().Seconds()
				return result
			} else {
				result := spell.CalcPeriodicDamage(sim, target, affliction.CalcScalingSpellDmg(uaScale), spell.OutcomeExpectedMagicCrit)
				result.Damage /= dot.CalcTickPeriod().Round(time.Millisecond).Seconds()
				return result
			}
		},
	})
}
