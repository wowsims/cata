package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) newUnleashElementsSpellConfig(unleashElementsTimer *core.Timer) core.SpellConfig {
	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: 73680},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    unleashElementsTimer,
				Duration: time.Second * 15,
			},
		},
	}
}

func (shaman *Shaman) registerUnleashElementsFlameTongue(unleashElementsTimer *core.Timer) {
	//TODO: confirm spell coefficient
	spellCoeff := 0.50

	config := shaman.newUnleashElementsSpellConfig(unleashElementsTimer)
	config.SpellSchool = core.SpellSchoolFire
	config.ProcMask = core.ProcMaskSpellDamage
	config.BonusHitRating = float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 1118 + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		spell.DealDamage(sim, result)
	}

	shaman.UnleashElementsFlameTongue = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashElementsFrostbrand(unleashElementsTimer *core.Timer) {
	//TODO: confirm spell coefficient
	spellCoeff := 0.40

	config := shaman.newUnleashElementsSpellConfig(unleashElementsTimer)
	config.SpellSchool = core.SpellSchoolFrost
	config.ProcMask = core.ProcMaskSpellDamage
	config.BonusHitRating = float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance

	//TODO: confirm spell coefficient
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 873 + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		spell.DealDamage(sim, result)
	}

	shaman.UnleashElementsFrostbrand = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashElementsWindfury(unleashElementsTimer *core.Timer) {
	config := shaman.newUnleashElementsSpellConfig(unleashElementsTimer)
	config.SpellSchool = core.SpellSchoolPhysical
	config.ProcMask = core.ProcMaskSpellDamage
	// TODO: Does this use spell hit?
	config.BonusHitRating = float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance

	// TODO: 175% weapon damage and apply buff here
	// config.ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
	// 	baseDamage := 2298 + spellCoeff*spell.SpellPower()
	// 	result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	// 	spell.DealDamage(sim, result)
	// },

	shaman.UnleashElementsWindfury = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashElementsEarthliving(unleashElementsTimer *core.Timer) {
	spellCoeff := 0.20
	//TODO: need?
	//bonusCoeff := 0.02 * float64(shaman.Talents.TidalWaves)

	config := shaman.newUnleashElementsSpellConfig(unleashElementsTimer)
	config.SpellSchool = core.SpellSchoolNature
	config.ProcMask = core.ProcMaskSpellHealing
	config.Flags = core.SpellFlagHelpful

	//TODO: apply buff for 30% on next direct heal
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		healPower := spell.HealingPower(target)
		baseHealing := 1996 + spellCoeff*healPower
		result := spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

		if result.Outcome.Matches(core.OutcomeCrit) {
			if shaman.Talents.AncestralAwakening > 0 {
				shaman.ancestralHealingAmount = result.Damage * 0.3

				// TODO: this should actually target the lowest health target in the raid.
				//  does it matter in a sim? We currently only simulate tanks taking damage (multiple tanks could be handled here though.)
				shaman.AncestralAwakening.Cast(sim, target)
			}
		}
	}

	shaman.UnleashElementsEarthliving = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashElements() {
	unleashElementsTimer := shaman.NewTimer()
	shaman.registerUnleashElementsFlameTongue(unleashElementsTimer)
	shaman.registerUnleashElementsFrostbrand(unleashElementsTimer)
	shaman.registerUnleashElementsWindfury(unleashElementsTimer)
	shaman.registerUnleashElementsEarthliving(unleashElementsTimer)
}
