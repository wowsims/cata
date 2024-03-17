package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (rogue *Rogue) registerAmbushSpell() {
	baseDamage := RogueBaseScalar * 0.327 + 28

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 8676},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.GetGeneratorCostModifier(60 - []float64{0, 7, 14, 20}[rogue.Talents.SlaughterFromTheShadows]),
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand) && rogue.IsStealthed()
		},

		BonusCritRating: 20*core.CritRatingPerCritChance*float64(rogue.Talents.ImprovedAmbush),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively.
		DamageMultiplier: core.TernaryFloat64(rogue.HasDagger(core.MainHand), 2.85, 1.97) * (1 +
			0.05*float64(rogue.Talents.ImprovedAmbush) +
			0.1*float64(rogue.Talents.Opportunity)),
		CritMultiplier:   rogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
