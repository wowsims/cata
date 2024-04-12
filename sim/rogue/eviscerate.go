package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (rogue *Rogue) registerEviscerate() {
	rogue.Eviscerate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 2098},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagFinisher | SpellFlagColdBlooded | core.SpellFlagAPL,
		MetricSplits: 6,

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.8,
			RefundMetrics: rogue.EnergyRefundMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				spell.SetMetricsSplit(spell.Unit.ComboPoints())
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return rogue.ComboPoints() > 0
		},

		BonusCritRating: core.TernaryFloat64(
			rogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfEviscerate), 10*core.CritRatingPerCritChance, 0.0),
		DamageMultiplier: 1,
		DamageMultiplierAdditive: 1 +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.CoupDeGrace] +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.Aggression],
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			comboPoints := rogue.ComboPoints()
			flatBaseDamage := 354 + 517*float64(comboPoints)
			// tooltip implies 3..7% AP scaling, but testing shows it's fixed at 7% (3.4.0.46158)
			apRatio := 0.091 * float64(comboPoints)

			baseDamage := flatBaseDamage +
				// 254.0*sim.RandomFloat("Eviscerate") + TODO: Thebackstabi 3/18/2024 - Cataclysm has no spell variance ATM, unsure on damage range
				apRatio*spell.MeleeAttackPower()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
				rogue.ApplyCutToTheChase(sim)
			} else {
				spell.IssueRefund(sim)
			}

			//spell.DealDamage(sim, result)
		},
	})
}
