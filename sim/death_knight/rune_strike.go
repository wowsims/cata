package death_knight

import (
	"github.com/wowsims/cata/sim/core"
)

var RuneStrikeActionID = core.ActionID{SpellID: 56815}

func (dk *DeathKnight) registerRuneStrikeSpell() {
	ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       RuneStrikeActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellRuneStrike,

		DamageMultiplier: 1.8,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.75,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.MeleeAttackPower()*0.05

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       RuneStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellRuneStrike,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 30,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return dk.HasActiveAura("Blood Presence")
		},

		DamageMultiplier: 1.8,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.75,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + spell.MeleeAttackPower()*0.1

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)

			spell.SpendRefundableCost(sim, result)
			dk.ThreatOfThassarianProc(sim, result, ohSpell)

			spell.DealDamage(sim, result)
		},
	})
}

// func (dk *DeathKnight) registerDrwRuneStrikeSpell() {
// 	runeStrikeGlyphCritBonus := core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfRuneStrike), 10.0, 0.0)

// 	dk.RuneWeapon.RuneStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
// 		ActionID:    RuneStrikeActionID.WithTag(1),
// 		SpellSchool: core.SpellSchoolPhysical,
// 		ProcMask:    core.ProcMaskMeleeMHSpecial,
// 		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

// 		BonusCritRating: (dk.annihilationCritBonus() + runeStrikeGlyphCritBonus) * core.CritRatingPerCritChance,
// 		DamageMultiplier: 1.5 *
// 			dk.darkrunedPlateRuneStrikeDamageBonus(),
// 		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
// 		ThreatMultiplier: 1.75,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := 0.15*spell.MeleeAttackPower() + dk.DrwWeaponDamage(sim, spell)
// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
// 		},
// 	})
// }
