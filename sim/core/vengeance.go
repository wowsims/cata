package core

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core/stats"
)

const VengeanceScaling = 0.018 // Might be reverted to 0.015 in a later patch

func (character *Character) RegisterVengeance(spellID int32, requiredAura *Aura) {
	// First register the exposed Vengeance buff Aura, which we will model
	// as discrete stacks with 1 AP granted per stack for ease of tracking
	// in the timeline and APLs.
	buffAura := MakeStackingAura(character, StackingStatAura{
		Aura: Aura{
			Label:     "Vengeance",
			ActionID:  ActionID{SpellID: spellID},
			Duration:  time.Second * 20,
			MaxStacks: math.MaxInt32,
		},

		BonusPerStack: stats.Stats{stats.AttackPower: 1},
	})

	// Then set up the proc trigger.
	vengeanceTrigger := ProcTrigger{
		Name:     "Vengeance Trigger",
		Callback: CallbackOnSpellHitTaken,

		Handler: func(sim *Simulation, spell *Spell, result *SpellResult) {
			// Check that the caster is an NPC.
			if spell.Unit.Type != EnemyUnit {
				return
			}

			// Vengeance uses pre-outcome, pre-mitigation damage.
			rawDamage := result.PreOutcomeDamage / result.ArmorMultiplier

			// The Weakened Blows debuff does not reduce Vengeance gains.
			// TODO: The game similarly hardcodes a correction for Demoralizing Banner, add that in once we implement the debuff in the sim.
			if (spell.SpellSchool == SpellSchoolPhysical) && spell.Unit.GetAura("Weakened Blows").IsActive() {
				rawDamage /= 0.9
			}

			// Note that result.PreOutcomeDamage does not include the impact of the tank's various DamageTakenMultiplier PseudoStats.
			// By default this is the desired behavior, since it means that tank DRs are automatically divided out in the calculation.
			// However, *detrimental* contributions to the relevant DamageTakenMultiplier PseudoStats *do* increase Vengeance gains in-game.
			// This can be relevant on certain bosses, such as Ignite Armor stacks increasing Vengeance gains on Iron Juggernaut in SoO.
			// TODO: Find a simple way to keep track of only detrimental contributions to DamageTakenMultiplier (and school-specific variants) with minimal overhead.

			// Apply baseline scaling to the raw damage value.
			rawVengeance := VengeanceScaling * rawDamage

			// Spells that are not mitigated by armor generate 2.5x more Vengeance.
			if (spell.SpellSchool != SpellSchoolPhysical) || spell.Flags.Matches(SpellFlagIgnoreArmor) {
				rawVengeance *= 2.5
			}

			// TODO: Is the 0.5x Vengeance multiplier for non-periodic AoE spells still a thing for the new version of Vengeance in Classic?

			// TODO: Weapon-based specials may be normalizing out spell.DamageMultiplier as well?

			// If the buff Aura is currently active, then perform decaying average with previous Vengeance.
			newVengeance := rawVengeance

			if buffAura.IsActive() {
				newVengeance += float64(buffAura.GetStacks()) * buffAura.RemainingDuration(sim).Seconds() / buffAura.Duration.Seconds()
			}

			// Compare to minimum ramp-up Vengeance value based on equilibrium estimate.
			var inferredAttackInterval time.Duration

			if spell.IsMH() {
				// TODO: Is this supposed to be the base speed prior to attack speed multipliers?
				inferredAttackInterval = spell.Unit.AutoAttacks.MainhandSwingSpeed()
			} else if spell.IsOH() {
				inferredAttackInterval = spell.Unit.AutoAttacks.OffhandSwingSpeed()
			} else {
				inferredAttackInterval = time.Minute
			}

			// TODO: Does this also need the 2.5x multiplier for spells and the 0.5x AoE multiplier in it?
			inferredEquilibriumVengeance := VengeanceScaling * rawDamage * buffAura.Duration.Seconds() / inferredAttackInterval.Seconds()

			if newVengeance < 0.5*inferredEquilibriumVengeance {
				if sim.Log != nil {
					result.Target.Log(sim, "Triggered Vengeance ramp-up mechanism because newVengeance = %.1f and inferredEquilibriumVengeance = %.1f .", newVengeance, inferredEquilibriumVengeance)
				}

				newVengeance = 0.5 * inferredEquilibriumVengeance
			}

			// Apply HP cap.
			newVengeance = min(newVengeance, result.Target.MaxHealth())

			if sim.Log != nil {
				result.Target.Log(sim, "Updated Vengeance for %s due to %s from %s. Raw damage value = %.1f, raw Vengeance contribution = %.1f, new Vengeance value = %.1f .", result.Target.Label, spell.ActionID, spell.Unit.Label, rawDamage, rawVengeance, newVengeance)
			}

			// Activate or refresh the buff Aura and set stacks.
			buffAura.Activate(sim)
			buffAura.SetStacks(sim, int32(math.Round(newVengeance)))
		},
	}

	// Finally, either create a new hidden Aura for the Vengeance trigger,
	// or attach it to the supplied parent Aura (Bear Form for Druids,
	// Defensive Stance for Warriors).
	if requiredAura == nil {
		MakeProcTriggerAura(&character.Unit, vengeanceTrigger)
	} else {
		requiredAura.AttachProcTrigger(vengeanceTrigger)
	}
}
