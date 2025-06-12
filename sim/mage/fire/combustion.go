package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerCombustionSpell() {
	actionID := core.ActionID{SpellID: 11129}

	combustionVariance := 0.17   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Variance"
	combustionScaling := 1.0     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Coefficient"
	combustionCoefficient := 1.0 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "BonusCoefficient"

	fire.combustion = fire.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage, // need to check proc mask for impact damage
		ClassSpellMask: mage.MageSpellCombustionApplication,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    fire.NewTimer(),
				Duration: time.Second * 45,
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           fire.DefaultCritMultiplier(),
		BonusCoefficient:         combustionCoefficient,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := fire.CalcAndRollDamageRange(sim, combustionScaling, combustionVariance)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.DealDamage(sim, result)
				spell.RelatedDotSpell.Cast(sim, target)
			}
		},
	})

	calculatedDotTick := func(target *core.Unit) float64 {
		tickDamage := 0.0
		spell := fire.Ignite
		dot := spell.Dot(target)
		if dot.IsActive() {
			tickDamage += dot.SnapshotBaseDamage / 2
		}
		return tickDamage
	}

	fire.combustion.RelatedDotSpell = fire.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: mage.MageSpellCombustion,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   fire.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Combustion Dot",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				tickBase := calculatedDotTick(target)
				dot.Snapshot(target, tickBase)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			tickBase := calculatedDotTick(target)
			result := spell.CalcPeriodicDamage(sim, target, tickBase, spell.OutcomeExpectedMagicAlwaysHit)

			critChance := spell.SpellCritChance(target)
			critMod := (critChance * (spell.CritMultiplier - 1))
			result.Damage *= 1 + critMod

			return result
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
