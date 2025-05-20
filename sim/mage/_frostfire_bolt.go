package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (mage *Mage) registerFrostfireBoltSpell() {

	hasGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostfireBolt)

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

		castSpeedMod := mage.AddDynamicMod(core.SpellModConfig{
			ClassMask:  mage.MageSpellFrostfireBolt,
			FloatValue: -500 * time.Millisecond,
			Kind:       core.SpellMod_CastTime_Flat,
		})

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2750,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: 1.5, // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "EffetBonusCoefficient"
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.5 * mage.ClassSpellScaling // Per https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=44614 Field "Coefficient"
			if hasGlyph{
				castSpeedMod.Activate()
			}
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					spell.DealDamage(sim, result)
				}
			})
		},
	})
}
