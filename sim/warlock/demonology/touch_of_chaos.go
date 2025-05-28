package demonology

import (
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
			return demonology.Metamorphosis.RelatedSelfBuff.IsActive() && demonology.DemonicFury.CanSpend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 28, 40))
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, demonology.CalcAndRollDamageRange(sim, tocScale, tocVariance), spell.OutcomeMagicHitAndCrit)
			demonology.DemonicFury.Spend(core.TernaryInt32(demonology.T15_2pc.IsActive(), 28, 40), spell.ActionID, sim)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)

				corruption := demonology.Corruption.Dot(target)
				if corruption.IsActive() {
					corruption.TakeSnapshot(sim, false)

					// add two ticks up to a limit of pandemic
					for idx := 0; idx < 2; idx++ {
						extended := float64(corruption.RemainingDuration(sim) + corruption.TickPeriod())
						maxLength := float64(corruption.BaseDuration() + corruption.BaseDuration()/2)
						if extended < maxLength {
							corruption.AddTick()
						}
					}
				}
			})
		},
	})
}
