package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var DeathStrikeActionID = core.ActionID{SpellID: 49998}

/*
Focuses dark power into a strike that deals 185% weapon damage plus 499 to an enemy and heals you for 20% of the damage you have sustained from non-player sources during the preceding 5 sec (minimum of at least 7% of your maximum health).
This attack cannot be parried.
*/
func (dk *DeathKnight) registerDeathStrike() {
	damageTakenInFive := 0.0

	hasBloodRites := dk.Inputs.Spec == proto.Spec_SpecBloodDeathKnight

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

	healingSpell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 45470},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: DeathKnightSpellDeathStrikeHeal,
		Flags:          core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healValue := damageTakenInFive * dk.deathStrikeHealingMultiplier
			healValueModed := spell.CalcHealing(sim, target, healValue, spell.OutcomeHealingNoHitCounter).Damage

			minHeal := spell.Unit.MaxHealth() * 0.07

			flags := spell.Flags
			healing := healValue
			if healValueModed < minHeal {
				// Remove caster modifiers for spell when doing min heal
				spell.Flags |= core.SpellFlagIgnoreAttackerModifiers
				healing = minHeal

				// Scent of Blood healing modifier is applied to the min heal
				// This **should** also be the only thing modifying the DamageMultiplier of this spell
				healing *= spell.DamageMultiplier
			}

			spell.CalcAndDealHealing(sim, target, healing, spell.OutcomeHealing)

			// Add back caster modifiers
			spell.Flags = flags
		},
	})

	var ohSpell *core.Spell
	if dk.Spec == proto.Spec_SpecFrostDeathKnight {
		ohSpell = dk.registerOffHandDeathStrike()
	}

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathStrikeActionID.WithTag(1),
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellDeathStrike,

		MaxRange: core.MaxMeleeRange,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 20,
			// Not actually refundable, but setting this to `true` if specced into blood
			// makes the default SpendCost function skip handling the rune cost and
			// lets us manually spend it with death rune conversion in ApplyEffects.
			Refundable: hasBloodRites,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		DamageMultiplier: 1.85,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.40000000596) +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialNoParry)

			if hasBloodRites {
				spell.SpendCostAndConvertFrostOrUnholyRune(sim, result.Landed())
			}

			if result.Landed() && dk.ThreatOfThassarianAura.IsActive() {
				ohSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)

			healingSpell.Cast(sim, &dk.Unit)
		},
	})
}

func (dk *DeathKnight) registerOffHandDeathStrike() *core.Spell {
	return dk.RegisterSpell(core.SpellConfig{
		ActionID:       DeathStrikeActionID.WithTag(2), // Actually 66188
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: DeathKnightSpellDeathStrike,

		DamageMultiplier: 1.85,
		CritMultiplier:   dk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.20000000298) +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}

func (dk *DeathKnight) registerDrwDeathStrike() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.CalcScalingSpellDmg(0.40000000596) +
				dk.RuneWeapon.StrikeWeapon.CalculateWeaponDamage(sim, spell.MeleeAttackPower()) +
				dk.RuneWeapon.StrikeWeaponDamage

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialNoParry)
		},
	})
}
