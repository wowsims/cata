package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerDeathAndDecaySpell() {
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 43265},
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty, // D&D doesn't seem to proc things in game.
		ClassSpellMask: DeathKnightSpellDeathAndDecay,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1.9,
		CritMultiplier:   dk.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Death and Decay" + dk.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// DnD recalculates everything on each tick
				baseDamage := 26 + dot.Spell.MeleeAttackPower()*0.06400000304
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.SpellMetrics[aoeTarget.UnitIndex].Casts++
					dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
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
