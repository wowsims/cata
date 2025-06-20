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
			FlatCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.42899999022,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 13.8 //Magic number from beta testing
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (fireElemental *FireElemental) registerFireNova() {
	levelScalingMultiplier := 91.517600 / 12.102900
	fireElemental.FireNova = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 117588},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAoE,

		ManaCost: core.ManaCostOptions{
			FlatCost: 30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1.00,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(49*levelScalingMultiplier, 58*levelScalingMultiplier) //Estimated from beta testing 49 58
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (fireElemental *FireElemental) registerImmolate() {
	actionID := core.ActionID{SpellID: 118297}

	fireElemental.Immolate = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: 1,
		CritMultiplier:   fireElemental.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1.0,

		ManaCost: core.ManaCostOptions{
			FlatCost: 95,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !fireElemental.IsGuardian()
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate",
			},
			NumberOfTicks:       7,
			TickLength:          time.Second * 3,
			BonusCoefficient:    0.34999999404,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, fireElemental.shamanOwner.CalcScalingSpellDmg(0.62800002098))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := fireElemental.shamanOwner.CalcScalingSpellDmg(1.79499995708)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.Dot(target).Apply(sim)
		},
	})
}

func (fireElemental *FireElemental) registerEmpower() {
	actionID := core.ActionID{SpellID: 118350}
	buffAura := fireElemental.shamanOwner.RegisterAura(core.Aura{
		Label:    "Empower",
		ActionID: actionID,
		Duration: core.NeverExpires,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.05,
	})

	fireElemental.Empower = fireElemental.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       core.SpellFlagChanneled,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    fireElemental.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Empower",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					buffAura.Activate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					buffAura.Deactivate(sim)
				},
			},
			NumberOfTicks:        1,
			TickLength:           time.Second * 60,
			AffectedByCastSpeed:  false,
			HasteReducesDuration: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !fireElemental.IsGuardian()
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Hot(target).Apply(sim)
		},
	})
}
