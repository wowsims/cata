package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerMaulSpell() {
	numHits := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMaul) && druid.Env.GetNumTargets() > 1, 2, 1)

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

		DamageMultiplier: 1.1 * core.TernaryFloat64(druid.AssumeBleedActive, RendAndTearDamageMultiplier, 1),
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		FlatThreatBonus:  30,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			curTarget := target
			anyLanded := false

			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

				if hitIndex > 0 {
					baseDamage *= 0.5
				}

				result := spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

				if result.Landed() {
					anyLanded = true
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			if !anyLanded {
				spell.IssueRefund(sim)
			}
		},
	})
}
