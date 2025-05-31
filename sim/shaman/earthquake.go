package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) registerEarthquakeSpell() {

	earthquakePulse := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 77478},
		Flags:            core.SpellFlagAoE | SpellFlagFocusable | core.SpellFlagIgnoreArmor,
		SpellSchool:      core.SpellSchoolPhysical,
		ClassSpellMask:   SpellMaskEarthquake,
		ProcMask:         core.ProcMaskSpellProc,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Casts-- // Do not count pulses as casts
			// Coefficient damage calculated manually because it's a Nature spell but deals Physical damage
			baseDamage := shaman.ClassSpellScaling*0.32400000095 + 0.11*spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	shaman.Earthquake = shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 77478},
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: 2500 * time.Millisecond,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Earthquake",
			},
			NumberOfTicks:        10,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				earthquakePulse.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			dot.Apply(sim)
		},
	})
}
