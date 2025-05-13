package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerMaulSpell() {
	flatBaseDamage := 34.0
	numHits := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMaul) && druid.Env.GetNumTargets() > 1, 2, 1)
	rendAndTearMod := []float64{1.0, 1.07, 1.13, 1.2}[druid.Talents.RendAndTear]

	druid.Maul = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6807},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellMaul,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 3,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		FlatThreatBonus:  30,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatBaseDamage + 0.19*spell.MeleeAttackPower()

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				modifier := 1.0
				if druid.BleedCategories.Get(curTarget).AnyActive() {
					modifier += .3
				}
				if druid.AssumeBleedActive || druid.Rip.Dot(curTarget).IsActive() || druid.Rake.Dot(curTarget).IsActive() || druid.Lacerate.Dot(curTarget).IsActive() {
					modifier *= rendAndTearMod
				}
				if hitIndex > 0 {
					modifier *= 0.5
				}

				result := spell.CalcAndDealDamage(sim, curTarget, baseDamage*modifier, spell.OutcomeMeleeSpecialHitAndCrit)

				if !result.Landed() {
					spell.IssueRefund(sim)
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
