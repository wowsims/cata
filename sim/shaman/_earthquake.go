package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerEarthquakeSpell() {
	shaman.Earthquake = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77478},
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty, // D&D doesn't seem to proc things in game.

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second*30 - time.Second*5*time.Duration(dk.Talents.Morbidity),
			},
		},

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Earthquake",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 62 + 0.0475*dk.getImpurityBonus(dot.Spell)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					// DnD recalculates attack multipliers + crit dynamically on every tick so this is here on purpose
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[aoeTarget.UnitIndex]) * dk.RoRTSBonus(aoeTarget)
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(aoeTarget)
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
