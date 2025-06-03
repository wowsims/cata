package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (fireElemental *FireElemental) registerFireBlast() {
	fireElemental.FireBlast = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57984},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 276,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 5,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.CritMultiplier(1.0, 0), // Spell 85801
		ThreatMultiplier: 1,
		BonusCoefficient: 0.429,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO these are approximation, from base SP
			spell.CalcAndDealDamage(sim, target, sim.Roll(220, 268), spell.OutcomeMagicHitAndCrit) //Estimated from beta testing
		},
	})
}

func (fireElemental *FireElemental) registerFireNova() {
	fireElemental.FireNova = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 424340},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAoE,

		ManaCost: core.ManaCostOptions{
			FlatCost: 207,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 5,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.CritMultiplier(1.0, 0), // Spell 85801
		ThreatMultiplier: 1,
		BonusCoefficient: 1.00,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(453, 537)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (fireElemental *FireElemental) registerFireShieldAura() {
	actionID := core.ActionID{SpellID: 13376}

	//dummy spell
	spell := fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAoE,

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.CritMultiplier(1.0, 0), // Spell 85801
		ThreatMultiplier: 1,
		BonusCoefficient: 0.032,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "FireShield",
			},
			NumberOfTicks: 40,
			TickLength:    time.Second * 3,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// TODO is this the right affect should it be Capped?
				// TODO these are approximation, from base SP
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, 102, dot.Spell.OutcomeMagicHitAndCrit) //Estimated from beta testing
				}
			},
		},
	})

	fireElemental.FireShieldAura = fireElemental.RegisterAura(core.Aura{
		Label:    "Fire Shield",
		ActionID: actionID,
		Duration: time.Minute * 2,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			spell.AOEDot().Apply(sim)
		},
	})
}
