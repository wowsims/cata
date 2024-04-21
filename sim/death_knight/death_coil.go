package death_knight

import (
	"github.com/wowsims/cata/sim/core"
)

var DeathCoilActionID = core.ActionID{SpellID: 47541}

func (dk *DeathKnight) registerDeathCoilSpell() {
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathCoilActionID,
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DeathKnightSpellDeathCoil,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassBaseScaling*0.87599998713 + spell.MeleeAttackPower()*0.23
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

// func (dk *DeathKnight) registerDrwDeathCoilSpell() {
// 	bonusFlatDamage := 443 + dk.sigilOfTheWildBuckBonus() + dk.sigilOfTheVengefulHeartDeathCoil()

// 	dk.RuneWeapon.DeathCoil = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
// 		ActionID:    DeathCoilActionID,
// 		SpellSchool: core.SpellSchoolShadow,
// 		ProcMask:    core.ProcMaskSpellDamage,

// 		BonusCritRating: dk.darkrunedBattlegearCritBonus() * core.CritRatingPerCritChance,
// 		DamageMultiplier: (1.0 + float64(dk.Talents.Morbidity)*0.05) *
// 			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDarkDeath), 1.15, 1.0),
// 		CritMultiplier:   dk.RuneWeapon.DefaultMeleeCritMultiplier(),
// 		ThreatMultiplier: 1.0,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := bonusFlatDamage + 0.15*dk.RuneWeapon.getImpurityBonus(spell)
// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
// 		},
// 	})
// }
