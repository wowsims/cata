package mage

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerIceLanceSpell() {

	// Values found at https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=30455
	iceLanceScaling := 0.335
	iceLanceCoefficient := 0.335
	iceLanceVariance := 0.25
	hasGlyphIcyVeins := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIcyVeins)
	hasGlyphSplittingIce := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfSplittingIce)

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30455},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellIceLance,
		MissileSpeed:   38,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1 * 1.2, // 2013-09-23 Ice Lance's damage has been increased by 20%
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: iceLanceCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// The target does not entirely appear to be random, but I was unable to determine how to tell which to target. IE: sat in front of 3 dummies it will always hit 2 specific ones.
			randomTarget := mage.Env.NextTargetUnit(target)
			hasSplittingIce := hasGlyphSplittingIce && mage.Env.GetNumTargets() > 1
			hasSplitBolts := mage.IcyVeinsAura.IsActive() && hasGlyphIcyVeins
			numberOfBolts := core.TernaryInt32(hasSplitBolts, 3, 1)
			icyVeinsDamageMultiplier := core.TernaryFloat64(hasSplitBolts, 0.4, 1.0)
			// Testing it does not appear to be exactly half, so I believe that this does its own damage calc with variance, it can also crit.

			// Secondary Target hit
			spell.DamageMultiplier *= icyVeinsDamageMultiplier
			if hasSplittingIce {
				spell.DamageMultiplier /= 2
				for range numberOfBolts {
					baseDamage := mage.CalcAndRollDamageRange(sim, iceLanceScaling, iceLanceVariance)
					result := spell.CalcDamage(sim, randomTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
					spell.WaitTravelTime(sim, func(sim *core.Simulation) {
						spell.DealDamage(sim, result)
					})
				}
				spell.DamageMultiplier *= 2
			}

			// Main Target hit
			results := make([]*core.SpellResult, 0)
			for range numberOfBolts {
				baseDamage := mage.CalcAndRollDamageRange(sim, iceLanceScaling, iceLanceVariance)
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				results = append(results, result)
			}
			if mage.FingersOfFrostAura.IsActive() {
				mage.FingersOfFrostAura.RemoveStack(sim)
			}
			for _, result := range results {
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}
			spell.DamageMultiplier /= icyVeinsDamageMultiplier

			if mage.Spec == proto.Spec_SpecFrostMage {
				// Confirmed in game Icicles launch even if ice lance misses.
				for _, icicle := range mage.Icicles {
					if hasSplittingIce {
						mage.SpendIcicle(sim, randomTarget, icicle/2)
					}
					mage.SpendIcicle(sim, target, icicle)
				}
				mage.Icicles = make([]float64, 0)
			}

		},
	})
}
