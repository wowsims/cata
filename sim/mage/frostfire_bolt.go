package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var frostfireBoltCoefficient = 1.5 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "EffetBonusCoefficient"
var frostfireBoltScaling = 1.5     // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "Coefficient"
var frostfireBoltVariance = 0.24   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "Variance"

func (mage *Mage) registerFrostfireBoltSpell() {

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
			if mage.BrainFreezeAura.IsActive() {
				mage.BrainFreezeAura.Deactivate(sim)
			}
			if mage.IcyVeinsAura.IsActive() && hasGlyph {
				for _ = range 3 {
					baseDamage := mage.CalcAndRollDamageRange(sim, frostfireBoltScaling, frostfireBoltVariance)
					result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
					result.Damage = result.Damage * .4
					spell.WaitTravelTime(sim, func(sim *core.Simulation) {
						spell.DealDamage(sim, result)
					})
					if result.Landed() {
						mage.HandleIcicleGeneration(sim, target, result.Damage)
					}
				}
			} else {
				baseDamage := mage.CalcAndRollDamageRange(sim, frostfireBoltScaling, frostfireBoltVariance)
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
					if mage.BrainFreezeAura.IsActive() {
						mage.BrainFreezeAura.Deactivate(sim)
					}
				})
				if result.Landed() {
					mage.HandleIcicleGeneration(sim, target, result.Damage)
				}
			}

		},
	})
}
