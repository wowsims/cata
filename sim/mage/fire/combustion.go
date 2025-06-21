package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerCombustionSpell() {

	hasGlyph := fire.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfCombustion)

	combustCD := core.Ternary(hasGlyph, time.Second*90, time.Second*45)
	combustDamageMultiplier := core.Ternary(hasGlyph, 2.0, 1.0)
	combustTickCount := core.Ternary(hasGlyph, 20, 10)

	actionID := core.ActionID{SpellID: 11129}

	combustionVariance := 0.17   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Variance"
	combustionScaling := 1.0     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "Coefficient"
	combustionCoefficient := 1.0 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.61217&filter%5BSpellID%5D=exact%253A2948 Field: "BonusCoefficient"

	fire.Combustion = fire.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage, // need to check proc mask for impact damage
		ClassSpellMask: mage.MageSpellCombustion,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    fire.NewTimer(),
				Duration: combustCD,
			},
		},
		DamageMultiplier: 1,
		CritMultiplier:   fire.DefaultCritMultiplier(),
		BonusCoefficient: combustionCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			fire.InfernoBlast.CD.Reset()
			spell.DamageMultiplier *= combustDamageMultiplier
			baseDamage := fire.CalcAndRollDamageRange(sim, combustionScaling, combustionVariance)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.RelatedDotSpell.Cast(sim, target)
			}
			spell.DamageMultiplier /= combustDamageMultiplier
		},
	})

	calculatedDotTick := func(sim *core.Simulation, target *core.Unit) float64 {
		spell := fire.Ignite
		dot := spell.Dot(target)
		if !dot.IsActive() {
			return 0.0
		}
		return dot.Spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick).Damage * .5
	}

	fire.Combustion.RelatedDotSpell = fire.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: mage.MageSpellCombustionDot,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   fire.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Combustion Dot",
			},
			NumberOfTicks:       int32(combustTickCount),
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				tickBase := calculatedDotTick(sim, target)
				dot.Snapshot(target, tickBase)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			tickBase := calculatedDotTick(sim, target)
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

	combustionTickCount := 0
	combustionTickDamage := 0.0
	updateCombustionTickCountEstimate := func() {
		combustionTickCount = int(fire.Combustion.RelatedDotSpell.Dot(fire.CurrentTarget).ExpectedTickCount())
	}
	updateCombustionTickDamageEstimate := func(sim *core.Simulation) {
		combustionTickDamage = fire.Combustion.RelatedDotSpell.ExpectedTickDamage(sim, fire.CurrentTarget)
	}

	updateCombustionTotalDamageEstimate := func() {
		combustionDotDamage := int32(float64(combustionTickCount) * combustionTickDamage)
		fire.combustionDotEstimate = combustionDotDamage
	}

	fire.AddOnCastSpeedChanged(func(old float64, new float64) {
		updateCombustionTickCountEstimate()
		updateCombustionTotalDamageEstimate()
	})

	fire.AddOnTemporaryStatsChange(func(sim *core.Simulation, _ *core.Aura, stats stats.Stats) {
		updateCombustionTickDamageEstimate(sim)
		updateCombustionTotalDamageEstimate()
	})

}
