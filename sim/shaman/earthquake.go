package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerEarthquakeSpell() {
	shaman.Earthquake = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 77478},
		Flags:            core.SpellFlagAPL | SpellFlagFocusable,
		SpellSchool:      core.SpellSchoolPhysical,
		ClassSpellMask:   SpellMaskEarthquake,
		ProcMask:         core.ProcMaskEmpty,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),

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
			IsAOE: true,
			Aura: core.Aura{
				Label: "Earthquake",
			},
			NumberOfTicks:        10,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteAffectsDuration: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Coefficient damage calculated manually because it's a Nature spell but deals Physical damage
				baseDamage := shaman.ClassSpellScaling*0.32400000095 + 0.11*dot.Spell.SpellPower()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
		},
	})
}
