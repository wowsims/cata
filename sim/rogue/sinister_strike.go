package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (rogue *Rogue) registerSinisterStrikeSpell() {
	hasGlyphOfSinisterStrike := rogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfSinisterStrike)
	baseDamage := RogueBaseDamageScalar*.178 + 3

	rogue.SinisterStrike = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1752},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.GetGeneratorCostModifier(45 - 2*float64(rogue.Talents.ImprovedSinisterStrike)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(Tier9, 4), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: []float64{0.0, .07, .14, .20}[rogue.Talents.Aggression] +
			0.01*float64(rogue.Talents.ImprovedSinisterStrike) +
			core.TernaryFloat64(rogue.HasSetBonus(Tier6, 4), 0.06, 0),
		DamageMultiplier: 1.04,
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := baseDamage +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				points := int32(1)
				if hasGlyphOfSinisterStrike {
					if sim.RandomFloat("Glyph of Sinister Strike") < 0.2 {
						points += 1
					}
				}
				rogue.AddComboPoints(sim, points, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
