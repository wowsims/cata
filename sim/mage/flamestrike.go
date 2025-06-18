package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerFlamestrikeSpell() {

	flameStrikeVariance := 0.2    // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2120 Field: "Variance"
	flameStrikeScaling := .46     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2120 Field: "Coefficient"
	flameStrikeCoefficient := .52 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2120 Field: "BonusCoefficient"
	flameStrikeDotScaling := .12
	flameStrikeDotCoefficient := .14

	mage.Flamestrike = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2120},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: MageSpellFlamestrike,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 6,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: flameStrikeCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := mage.CalcAndRollDamageRange(sim, flameStrikeScaling, flameStrikeVariance)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			spell.RelatedDotSpell.AOEDot().Apply(sim)
		},
	})

	mage.Flamestrike.RelatedDotSpell = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2120}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellFlamestrikeDot,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: flameStrikeDotCoefficient,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "FlameStrike DOT",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.Snapshot(target, mage.CalcScalingSpellDmg(flameStrikeDotScaling))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeSnapshotCrit)
				}
			},
		},
	})
}
