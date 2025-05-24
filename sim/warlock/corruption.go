package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const corruptionScale = 0.165
const corruptionCoeff = 0.165

func (warlock *Warlock) RegisterCorruption(callback WarlockSpellCastedCallback) *core.Spell {
	warlock.Corruption = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 172},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellCorruption,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 1.25},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Corruption",
			},
			NumberOfTicks:       9,
			TickLength:          2 * time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    corruptionCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, warlock.CalcScalingSpellDmg(corruptionScale))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				callback([]core.SpellResult{*result}, dot.Spell, sim)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
			} else {
				return spell.CalcPeriodicDamage(sim, target, 1033/6, spell.OutcomeExpectedMagicCrit)
			}
		},
	})

	return warlock.Corruption
}
