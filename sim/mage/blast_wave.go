package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerBlastWaveSpell() {
	if !mage.Talents.BlastWave {
		return
	}

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11113},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellBlastWave,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:    time.Millisecond * 500,
				GCDMin: time.Millisecond * 500,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 15,
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		BonusCoefficient:         0.193,
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var targetCount int32
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				targetCount++
				baseDamage := sim.Roll(1047, 1233)
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			// if targetCount > 1 {
			// 	mage.FlamestrikeBW.Cast(sim, target)
			// }
		},
	})

	// mage.FlamestrikeBW = mage.RegisterSpell(core.SpellConfig{
	// 	ActionID:       core.ActionID{SpellID: 88148}.WithTag(1),
	// 	SpellSchool:    core.SpellSchoolFire,
	// 	ProcMask:       core.ProcMaskSpellDamage,
	// 	ClassSpellMask: MageSpellFlamestrike,

	// 	ManaCost: core.ManaCostOptions{
	// 		BaseCost: 0.30,
	// 	},
	// 	Cast: core.CastConfig{
	// 		DefaultCast: core.Cast{
	// 			NonEmpty: true,
	// 		},
	// 	},
	// 	DamageMultiplierAdditive: 1,
	// 	CritMultiplier:           mage.DefaultSpellCritMultiplier(),
	// 	BonusCoefficient:         0.146,
	// 	ThreatMultiplier:         1,
	// 	Dot: core.DotConfig{
	// 		IsAOE: true,
	// 		Aura: core.Aura{
	// 			Label: "Flamestrike (Blast Wave)",
	// 		},
	// 		NumberOfTicks: 4,
	// 		TickLength:    time.Second * 2,
	// 		OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
	// 			target := mage.CurrentTarget
	// 			baseDamage := 0.103 * mage.ClassSpellScaling
	// 			dot.Snapshot(target, baseDamage)
	// 		},
	// 		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
	// 			for _, aoeTarget := range sim.Encounter.TargetUnits {
	// 				dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeSnapshotCrit)
	// 			}
	// 		},
	// 		BonusCoefficient: 0.061,
	// 	},

	// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 		for _, aoeTarget := range sim.Encounter.TargetUnits {
	// 			baseDamage := 0.662 * mage.ClassSpellScaling
	// 			baseDamage *= sim.Encounter.AOECapMultiplier()
	// 			spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 		}
	// 		spell.AOEDot().Apply(sim)
	// 	},
	// })
}
