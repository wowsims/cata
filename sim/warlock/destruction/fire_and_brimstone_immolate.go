package destruction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (destruction *DestructionWarlock) registerFireAndBrimstoneImmolate() {
	fabImmolate := destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 108686},
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
			ModifyCast: func(sim *core.Simulation, _ *core.Spell, _ *core.Cast) {
				if !destruction.FABAura.IsActive() {
					destruction.FABAura.Activate(sim)
				}
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   destruction.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: immolateCoeff,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return destruction.BurningEmbers.CanSpend(10)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			reduction := destruction.getFABReduction()
			spell.DamageMultiplier *= reduction
			spell.RelatedDotSpell.DamageMultiplier *= reduction

			destruction.BurningEmbers.Spend(sim, 10, spell.ActionID)
			for _, enemy := range sim.Environment.Encounter.TargetUnits {
				result := spell.CalcDamage(sim, enemy, destruction.CalcScalingSpellDmg(immolateScale), spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.RelatedDotSpell.Cast(sim, enemy)
				}

				if result.DidCrit() {
					destruction.BurningEmbers.Gain(sim, 1, spell.ActionID)
				}

				spell.DealDamage(sim, result)
			}

			spell.DamageMultiplier /= reduction
			spell.RelatedDotSpell.DamageMultiplier /= reduction
		},
	})

	fabImmolate.RelatedDotSpell = destruction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 108686}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: warlock.WarlockSpellImmolateDot,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   destruction.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FAB - Immolate (DoT)",
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ClassSpellMask == warlock.WarlockSpellImmolate && spell != fabImmolate {
						if fabImmolate.RelatedDotSpell.Dot(result.Target).IsActive() {
							fabImmolate.RelatedDotSpell.Dot(result.Target).Deactivate(sim)
						}
					}
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
			// both immolate versions are mutually exlusive and the og one will always be stronger
			// baf does not overwrite default dot
			if destruction.Immolate.Dot(target).IsActive() {
				return
			}

			destruction.ApplyDotWithPandemic(spell.Dot(target), sim)
		},
	})
}
