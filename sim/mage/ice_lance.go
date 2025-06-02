package mage

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

// Values found at https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=30455
var iceLanceScaling = 0.34
var iceLanceCoefficient = 0.34
var iceLanceVariance = 0.25

func (mage *Mage) registerIceLanceSpell() {

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

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultCritMultiplier(),
		BonusCoefficient:         iceLanceCoefficient,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if mage.FingersOfFrostAura.IsActive() {
				mage.FingersOfFrostAura.RemoveStack(sim)
			}

			// The target does not entirely appear to be random, but I was unable to determine how to tell which to target. IE: sat in front of 3 dummies it will always hit 2 specific ones.
			randomTarget := sim.Encounter.TargetUnits[int(sim.Roll(0, float64(len(sim.Encounter.TargetUnits))))]
			// Testing it does not appear to be exactly half, so I believe that this does its own damage calc with variance, it can also crit.
			if hasGlyphSplittingIce {
				baseDamage := mage.CalcAndRollDamageRange(sim, iceLanceScaling, iceLanceVariance) / 2
				if mage.IcyVeinsAura.IsActive() && hasGlyphIcyVeins {
					baseDamage = baseDamage * .4
					for _ = range 3 {
						result := spell.CalcDamage(sim, randomTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
						spell.WaitTravelTime(sim, func(sim *core.Simulation) {
							spell.DealDamage(sim, result)
						})
					}
				} else {
					result := spell.CalcDamage(sim, randomTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
					spell.WaitTravelTime(sim, func(sim *core.Simulation) {
						spell.DealDamage(sim, result)
					})
				}
			}
			baseDamage := mage.CalcAndRollDamageRange(sim, iceLanceScaling, iceLanceVariance)
			if mage.IcyVeinsAura.IsActive() && hasGlyphIcyVeins {
				baseDamage = baseDamage * .4
				for _ = range 3 {
					result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
					spell.WaitTravelTime(sim, func(sim *core.Simulation) {
						spell.DealDamage(sim, result)
					})
				}
			} else {
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
				})
			}

			if mage.Spec == proto.Spec_SpecFrostMage {
				//I've confirmed in game Icicles launch even if ice lance misses.
				for _, icicle := range mage.Icicles {
					if hasGlyphSplittingIce {
						mage.castIcicleWithDamage(sim, randomTarget, icicle/2)
					}
					mage.castIcicleWithDamage(sim, target, icicle)
				}
				mage.Icicles = make([]float64, 0)

			}

		},
	})
}
