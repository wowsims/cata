package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO: Cleanup death strike the same way we did for plague strike
var DeathStrikeActionID = core.ActionID{SpellID: 49998}

func (dk *DeathKnight) registerDeathStrikeSpell() {
	damageTakenInFive := 0.0

	core.MakePermanent(dk.GetOrRegisterAura(core.Aura{
		Label: "Death Strike Damage Taken",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				damageTaken := result.Damage
				damageTakenInFive += damageTaken

				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: sim.CurrentTime + time.Second*5,
					OnAction: func(s *core.Simulation) {
						damageTakenInFive -= damageTaken
					},
				})
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageTakenInFive = 0.0
		},
	}))

	healingSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       DeathStrikeActionID.WithTag(3),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskSpellHealing,
		ClassSpellMask: DeathKnightSpellDeathStrikeHeal,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
	})

	ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       DeathStrikeActionID.WithTag(2),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: DeathKnightSpellDeathStrike,

		DamageMultiplier: 1.5,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassBaseScaling*0.14699999988 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)

			if result.Landed() {
				healingSpell.Cast(sim, &dk.Unit)
				healingSpell.CalcAndDealHealing(sim, &dk.Unit, max(damageTakenInFive*0.05, dk.MaxHealth()*0.07), healingSpell.OutcomeHealing)
			}
		},
	})

	dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       DeathStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellDeathStrike,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 20,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.5,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassBaseScaling*0.29399999976 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			spell.SpendRefundableCost(sim, result)
			dk.ThreatOfThassarianProc(sim, result, ohSpell)

			if result.Landed() {
				healingSpell.Cast(sim, &dk.Unit)
				healingSpell.CalcAndDealHealing(sim, &dk.Unit, max(damageTakenInFive*0.2, dk.MaxHealth()*0.07), healingSpell.OutcomeHealing)
			}

			spell.DealDamage(sim, result)
		},
	})
}

// func (dk *DeathKnight) registerDrwDeathStrikeSpell() {
// 	bonusBaseDamage := dk.sigilOfAwarenessBonus()
// 	hasGlyph := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDeathStrike)

// 	dk.RuneWeapon.DeathStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
// 		ActionID:    DeathStrikeActionID.WithTag(1),
// 		SpellSchool: core.SpellSchoolPhysical,
// 		ProcMask:    core.ProcMaskMeleeMHSpecial,
// 		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

// 		BonusCritRating:  (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
// 		DamageMultiplier: .75 * dk.improvedDeathStrikeDamageBonus(),
// 		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := 297 + bonusBaseDamage + dk.DrwWeaponDamage(sim, spell)

// 			if hasGlyph {
// 				baseDamage *= 1 + 0.01*min(dk.CurrentRunicPower(), 25)
// 			}
// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
// 		},
// 	})
// }
