package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const agonyScale = 0.0255 * 1.18
const agonyCoeff = 0.0255 * 1.18

func (affliction *AfflictionWarlock) registerAgony() {
	affliction.Agony = affliction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 980},
		Flags:          core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellDamage,
		SpellSchool:    core.SpellSchoolShadow,
		ClassSpellMask: warlock.WarlockSpellAgony,

		ThreatMultiplier: 1,
		DamageMultiplier: 1,
		BonusCoefficient: agonyCoeff,
		CritMultiplier:   affliction.DefaultCritMultiplier(),

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Agony",
				MaxStacks: 10,
			},

			TickLength:          2 * time.Second,
			NumberOfTicks:       12,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, affliction.CalcScalingSpellDmg(agonyScale))
			},

			BonusCoefficient: agonyCoeff,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				var stacks int32 = 10

				// on the last tick the aura seems to be deactivated first
				if dot.Aura.IsActive() {
					dot.Aura.AddStack(sim)
					stacks = dot.Aura.GetStacks()
				}

				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeMagicHitAndSnapshotCrit)
				result.Damage *= float64(stacks)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit).Landed() {
				affliction.ApplyDotWithPandemic(spell.Dot(target), sim)
				spell.Dot(target).AddStack(sim)
			}
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			dot := spell.Dot(target)

			// Always compare fully stacked agony damage
			if useSnapshot {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
				result.Damage *= 10
				result.Damage /= dot.TickPeriod().Seconds()
				return result
			} else {
				result := spell.CalcPeriodicDamage(sim, target, affliction.CalcScalingSpellDmg(agonyScale), spell.OutcomeExpectedMagicCrit)
				result.Damage *= 10
				result.Damage /= dot.CalcTickPeriod().Round(time.Millisecond).Seconds()
				return result
			}
		},
	})
}
