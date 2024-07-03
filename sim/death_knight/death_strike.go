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

	doHealing := func(sim *core.Simulation, value float64) {
		healValue := damageTakenInFive * value
		healValueModed := healingSpell.CalcHealing(sim, healingSpell.Unit, healValue, healingSpell.OutcomeHealing).Damage

		minHeal := healingSpell.Unit.MaxHealth() * 0.07

		healing := healValue
		if healValueModed < minHeal {
			// Remove caster modifiers for spell when doing min heal
			healingSpell.Flags |= core.SpellFlagIgnoreAttackerModifiers
			healing = minHeal
		}
		healingSpell.Cast(sim, healingSpell.Unit)
		healingSpell.CalcAndDealHealing(sim, healingSpell.Unit, healing, healingSpell.OutcomeHealing)

		// Add back caster modifiers
		healingSpell.Flags ^= core.SpellFlagIgnoreAttackerModifiers
	}

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
			baseDamage := dk.ClassSpellScaling*0.14699999988 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
			doHealing(sim, 0.05)
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
			baseDamage := dk.ClassSpellScaling*0.29399999976 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			dk.ThreatOfThassarianProc(sim, result, ohSpell)

			spell.DealDamage(sim, result)
			doHealing(sim, 0.2)
		},
	})
}

func (dk *DeathKnight) registerDrwDeathStrikeSpell() *core.Spell {
	return dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dk.ClassSpellScaling*0.29399999976 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})
}
