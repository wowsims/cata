package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const immolateScale = 0.47 * 1.13 // Hotfix
const immolateCoeff = 0.47 * 1.13

func (destruction *DestructionWarlock) registerImmolate() {
	destruction.Immolate = destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 348},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellImmolate,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 3},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   destruction.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: immolateCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if destruction.FABAura.IsActive() {
				destruction.FABAura.Deactivate(sim)
			}

			result := spell.CalcDamage(sim, target, destruction.CalcScalingSpellDmg(immolateScale), spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.RelatedDotSpell.Cast(sim, target)
			}

			if result.DidCrit() {
				destruction.BurningEmbers.Gain(sim, 1, spell.ActionID)
			}

			spell.DealDamage(sim, result)
		},
	})

	destruction.Immolate.RelatedDotSpell = destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 348}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: warlock.WarlockSpellImmolateDot,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   destruction.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Immolate (DoT)",
			},
			NumberOfTicks:       5,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: true,
			BonusCoefficient:    immolateCoeff,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, destruction.CalcScalingSpellDmg(immolateScale))
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				if result.DidCrit() {
					destruction.BurningEmbers.Gain(sim, 1, dot.Spell.ActionID)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			destruction.ApplyDotWithPandemic(spell.Dot(target), sim)
		},
	})
}
