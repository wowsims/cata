package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const drainLifeScale = 0.334
const drainLifeCoeff = 0.334

func (warlock *Warlock) RegisterDrainLife(callback WarlockSpellCastedCallback) {
	manaMetric := warlock.NewManaMetrics(core.ActionID{SpellID: 689})
	healthMetric := warlock.NewHealthMetrics(core.ActionID{SpellID: 689})

	warlock.DrainLife = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 689},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDrainLife,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 1},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         drainLifeCoeff,

		Dot: core.DotConfig{
			Aura:                 core.Aura{Label: "Drain Life"},
			NumberOfTicks:        6,
			TickLength:           1 * time.Second,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			BonusCoefficient:     drainLifeCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, warlock.CalcScalingSpellDmg(drainLifeScale))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

				// Spend mana per tick
				warlock.SpendMana(sim, dot.Spell.Cost.GetCurrentCost(), manaMetric)
				warlock.GainHealth(sim, warlock.MaxHealth()*0.02, healthMetric)

				if callback != nil {
					callback([]core.SpellResult{*result}, dot.Spell, sim)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				spell.DealOutcome(sim, result)
			}
		},
	})
}
