package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerEviscerate() {
	coefficient := 0.32600000501
	resourceCoefficient := 0.47600001097
	apScalingPerComboPoint := 0.091

	avgBaseDamage := coefficient * rogue.ClassSpellScaling
	damagePerComboPoint := resourceCoefficient * rogue.ClassSpellScaling
	baseMinDamage := avgBaseDamage * 0.5

	rogue.Eviscerate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2098},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagFinisher | SpellFlagColdBlooded | core.SpellFlagAPL,
		MetricSplits:   6,
		ClassSpellMask: RogueSpellEviscerate,

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

		BonusCritPercent: core.TernaryFloat64(
			rogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfEviscerate), 10, 0),
		DamageMultiplier: 1,
		DamageMultiplierAdditive: 1 +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.CoupDeGrace] +
			[]float64{0.0, 0.07, 0.14, 0.2}[rogue.Talents.Aggression],
		CritMultiplier:   rogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)

			comboPoints := float64(rogue.ComboPoints())
			baseDamage := baseMinDamage +
				sim.RandomFloat("Eviscerate")*avgBaseDamage +
				damagePerComboPoint*comboPoints +
				apScalingPerComboPoint*comboPoints*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.ApplyFinisher(sim, spell)
				rogue.ApplyCutToTheChase(sim)
			} else {
				spell.IssueRefund(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}
