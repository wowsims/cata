package demonology

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const tocScale = 0.767
const tocVariance = 0.1
const tocCoeff = 0.767

func (demonology *DemonologyWarlock) registerTouchOfChaos() {
	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 103964},
		SpellSchool:    core.SpellSchoolChaos,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellTouchOfChaos,
		MissileSpeed:   120,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           demonology.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         tocCoeff,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return demonology.IsInMeta() && demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 28, 40))
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, demonology.CalcAndRollDamageRange(sim, tocScale, tocVariance), spell.OutcomeMagicHitAndCrit)
			demonology.DemonicFury.Spend(sim, core.TernaryInt32(demonology.T15_2pc.IsActive(), 28, 40), spell.ActionID)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)

				corruption := demonology.Corruption.Dot(target)
				if corruption.IsActive() {
					corruption.TakeSnapshot(sim, false)

					// most sane way I can think off to keep tick count but update haste tick rate and roll over properly
					// duration is actually extended on refresh for lower haste
					state := corruption.SaveState(sim)
					corruption.ApplyRollover(sim)
					state.ExtraTicks = 0
					state.TickPeriod = corruption.TickPeriod()
					state.RemainingDuration = state.TickPeriod*time.Duration(state.TicksRemaining) + state.NextTickIn
					corruption.RestoreState(state, sim)

					// add up to the max duration or up to 5 seconds
					maxLength := math.Min(float64(corruption.BaseDuration()+corruption.BaseDuration()/2), float64(corruption.RemainingDuration(sim)+time.Second*5))

					for idx := 0; float64(corruption.RemainingDuration(sim)+corruption.TickPeriod()) < maxLength; idx++ {
						corruption.AddTick()
					}
				}
			})
		},
	})
}
