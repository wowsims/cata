package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerFrostfireBoltSpell() {

	frostfireBoltCoefficient := 1.5 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "EffetBonusCoefficient"
	frostfireBoltScaling := 1.5     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "Coefficient"
	frostfireBoltVariance := 0.24   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "Variance"

	hasGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIcyVeins)

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44614},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostfireBolt,
		MissileSpeed:   28,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2750,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: frostfireBoltCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hasSplitBolts := mage.IcyVeinsAura.IsActive() && hasGlyph
			numberOfBolts := core.TernaryInt32(hasSplitBolts, 3, 1)
			damageMultiplier := core.TernaryFloat64(hasSplitBolts, 0.4, 1.0)
			results := make([]*core.SpellResult, numberOfBolts)

			spell.DamageMultiplier *= damageMultiplier
			for idx := range numberOfBolts {
				baseDamage := mage.CalcAndRollDamageRange(sim, frostfireBoltScaling, frostfireBoltVariance)
				results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			spell.DamageMultiplier /= damageMultiplier
			mage.BrainFreezeAura.Deactivate(sim)

			for _, result := range results {
				if result.Landed() {
					mage.ProcFingersOfFrost(sim, spell)
				}
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
					if result.Landed() {
						mage.GainIcicle(sim, target, result.Damage)
					}
				})
			}
		},
	})
}
