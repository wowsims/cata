package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerEarthquakeSpell() {
	shaman.Earthquake = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77478},
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		//TODO: Not sure on the logic
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Earthquake",
			},
			NumberOfTicks:    10,
			TickLength:       time.Second * 1,
			BonusCoefficient: 0.119,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.SnapshotBaseDamage = 326
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(aoeTarget)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeMagicHitAndSnapshotCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)
		},
	})
}
