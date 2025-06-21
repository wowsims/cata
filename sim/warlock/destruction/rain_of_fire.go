package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

// Rain of Fire does not need to be channeled anymore for destruction warlocks
// RoF applies a hidden dot to the target which exhibits all the usual dot mechanics
// This also causes rof to not stack on the same target
//
// Measured proc rate for rof on 2 targets is 70 procs on 550 ticks ~12.5% = 1/8th

var rofScale = 0.15
var rofCoeff = 0.15

func (destruction DestructionWarlock) registerRainOfFire() {
	baseDamage := destruction.CalcScalingSpellDmg(rofScale)
	destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 104232},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellRainOfFire,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 6.25,
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           destruction.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         rofCoeff,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Rain of Fire (DoT)",
				ActionID: core.ActionID{SpellID: 104232}.WithTag(1),
			},

			TickLength:           time.Second,
			NumberOfTicks:        8,
			HasteReducesDuration: true,
			IsAOE:                true,
			BonusCoefficient:     rofCoeff,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					result := dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, baseDamage, dot.OutcomeTickMagicCrit)
					if result.Landed() && sim.Proc(0.125, "RoF - Ember Proc") {
						destruction.BurningEmbers.Gain(sim, 2, dot.ActionID)
					}
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
