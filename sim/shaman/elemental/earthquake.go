package elemental

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/shaman"
)

func (elemental *ElementalShaman) registerEarthquakeSpell() {

	earthquakePulse := elemental.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 77478},
		Flags:            shaman.SpellFlagShamanSpell | core.SpellFlagAoE | shaman.SpellFlagFocusable | core.SpellFlagIgnoreArmor,
		SpellSchool:      core.SpellSchoolPhysical,
		ClassSpellMask:   shaman.SpellMaskEarthquake,
		ProcMask:         core.ProcMaskSpellProc,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   elemental.DefaultCritMultiplier(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Casts-- // Do not count pulses as casts
			// Coefficient damage calculated manually because it's a Nature spell but deals Physical damage
			baseDamage := elemental.CalcScalingSpellDmg(0.32400000095) + 0.1099999994*spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	elemental.Earthquake = elemental.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 77478},
		Flags:    shaman.SpellFlagShamanSpell | core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 70.3,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: 2500 * time.Millisecond,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    elemental.NewTimer(),
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
