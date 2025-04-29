package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerStarfireSpell() {
	druid.SetSpellEclipseEnergy(DruidSpellStarfire, StarfireBaseEnergyGain, StarfireBaseEnergyGain)

	hasStarfireGlyph := druid.HasMajorGlyph(proto.DruidMajorGlyph(proto.DruidPrimeGlyph_GlyphOfStarfire))

	starfireGlyphSpell := druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 54845},
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagNoLogs,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			moonfireDot := druid.Moonfire.Dot(target)
			sunfireDot := druid.Sunfire.Dot(target)

			tryExtendDot(moonfireDot, &druid.ExtendingMoonfireStacks)
			tryExtendDot(sunfireDot, &druid.ExtendingMoonfireStacks)
		},
	})

	druid.Starfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2912},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellStarfire,
		Flags:          core.SpellFlagAPL | SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 11,
			PercentModifier: 100,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 3200,
			},
		},

		BonusCoefficient: 1.231,

		// TODO: Was the value of 1 here incorrect to begin with?
		BonusCritPercent: 1 / core.CritRatingPerCritPercent,

		DamageMultiplier: 1,

		CritMultiplier: druid.BalanceCritMultiplier(),

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			min, max := core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassDruid, 1.383, 0.22)
			baseDamage := sim.Roll(min, max)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				if hasStarfireGlyph {
					starfireGlyphSpell.Cast(sim, target)
				}

				spell.DealDamage(sim, result)
			}
		},
	})
}

func tryExtendDot(dot *core.Dot, extendingStacks *int) {
	if dot.IsActive() && *extendingStacks > 0 {
		*extendingStacks -= 1
		dot.UpdateExpires(dot.ExpiresAt() + time.Second*3)
	}
}
