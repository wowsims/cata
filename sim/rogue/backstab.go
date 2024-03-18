package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (rogue *Rogue) registerBackstabSpell() {
	hasGlyph := rogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfBackstab)
	baseDamage := RogueBaseDamageScalar*.307 + 10

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.GetGeneratorCostModifier(60 - []float64{0, 7, 14, 20}[rogue.Talents.SlaughterFromTheShadows]),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand)
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(Tier9, 4), 5*core.CritRatingPerCritChance, 0) +
			10*core.CritRatingPerCritChance*float64(rogue.Talents.PuncturingWounds),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively (up to 142%).
		DamageMultiplier: 2.07*(1+
			0.1*float64(rogue.Talents.Opportunity)+
			[]float64{0.0, .07, .14, .20}[rogue.Talents.Aggression]+
			core.TernaryFloat64(rogue.HasSetBonus(Tier6, 4), 0.06, 0)) +
			core.TernaryFloat64(rogue.Spec == proto.Spec_SpecSubtletyRogue, .4, 0),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				if hasGlyph && result.DidCrit() {
					rogue.AddEnergy(sim, 5, rogue.EnergyRefundMetrics)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
