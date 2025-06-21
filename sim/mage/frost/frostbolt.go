package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

func (frostMage *FrostMage) registerFrostboltSpell() {
	frostboltVariance := 0.24   // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A116 Field: "Variance"
	frostboltScale := 1.5       // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A116 Field: "Coefficient"
	frostboltCoefficient := 1.5 // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=exact%253A116 Field: "BonusCoefficient"
	hasGlyph := frostMage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfIcyVeins)

	frostMage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 116},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellFrostbolt,
		MissileSpeed:   28,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   frostMage.DefaultCritMultiplier(),
		BonusCoefficient: frostboltCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			hasSplitBolts := frostMage.IcyVeinsAura.IsActive() && hasGlyph
			numberOfBolts := core.TernaryInt32(hasSplitBolts, 3, 1)
			damageMultiplier := core.TernaryFloat64(hasSplitBolts, 0.4, 1.0)

			spell.DamageMultiplier *= damageMultiplier
			for range numberOfBolts {
				baseDamage := frostMage.CalcAndRollDamageRange(sim, frostboltScale, frostboltVariance)
				result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					frostMage.ProcFingersOfFrost(sim, spell)
				}
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.DealDamage(sim, result)
					if result.Landed() {
						frostMage.GainIcicle(sim, target, result.Damage)
					}
				})
			}
			spell.DamageMultiplier /= damageMultiplier
		},
	})
}
