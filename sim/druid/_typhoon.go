package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerTyphoonSpell() {
	if !druid.Talents.Typhoon {
		return
	}

	hasTyphoonGlyph := druid.HasMinorGlyph(proto.DruidMinorGlyph_GlyphOfTyphoon)
	hasMonsoonGlyph := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMonsoon)

	druid.Typhoon = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 61391},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellTyphoon,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL | SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 16,
			PercentModifier: 100 - (8 * core.TernaryInt32(hasTyphoonGlyph, 1, 0)),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * (20 - core.TernaryDuration(hasMonsoonGlyph, 3, 0)),
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   1,
		BonusCoefficient: 0.126,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				baseDamage := core.CalcScalingSpellAverageEffect(proto.Class_ClassDruid, 1.316)
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
			})
		},
	})
}
