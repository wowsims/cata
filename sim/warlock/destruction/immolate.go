package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

const immolateScale = 0.47 * 1.13 // Hotfix
const immolateCoeff = 0.47 * 1.13

// Damage Done By Caster setup
const (
	DDBC_Immolate int = iota
	DDBC_Total
)

func (destruction *DestructionWarlock) registerImmolate() {
	actionID := core.ActionID{SpellID: 348}
	destruction.Immolate = destruction.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
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
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					core.EnableDamageDoneByCaster(DDBC_Immolate, DDBC_Total, destruction.AttackTables[aura.Unit.UnitIndex], immolateDamageDoneByCasterHandler)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					core.DisableDamageDoneByCaster(DDBC_Immolate, destruction.AttackTables[aura.Unit.UnitIndex])
				},
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

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			dot := spell.Dot(target)
			if useSnapshot {
				result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedSnapshotCrit)
				result.Damage /= dot.TickPeriod().Seconds()
				return result
			} else {
				result := spell.CalcPeriodicDamage(sim, target, destruction.CalcScalingSpellDmg(immolateCoeff), spell.OutcomeExpectedMagicCrit)
				result.Damage /= dot.CalcTickPeriod().Round(time.Millisecond).Seconds()
				return result
			}
		},
	})
}

func immolateDamageDoneByCasterHandler(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
	if spell.Matches(warlock.WarlockSpellRainOfFire) {
		return 1.5
	}

	return 1
}
