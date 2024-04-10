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
			BaseCost:   0.07,
			Multiplier: 1 - 0.05*float64(shaman.Talents.Convection) - shaman.GetMentalQuicknessBonus(),
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

func (shaman *Shaman) registerUnleashFlame(unleashElementsTimer *core.Timer) {
	spellCoeff := 0.429

	var affectedSpells []*core.Spell

	unleashFlameAura := shaman.RegisterAura(core.Aura{
		Label:    "Unleash Flame",
		ActionID: core.ActionID{SpellID: 73683},
		Duration: time.Second * 8,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: need to confirm all spells that benefit from this
			affectedSpells = core.FilterSlice([]*core.Spell{
				shaman.LavaBurst,
				shaman.FlameShock,
				shaman.FireNova,
				shaman.LavaLash,
				shaman.MagmaTotem,
				shaman.SearingTotem,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.2
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier -= 0.2
			}
		},
	})

	config := shaman.newUnleashElementsSpellConfig(unleashElementsTimer)
	config.SpellSchool = core.SpellSchoolFire
	config.ProcMask = core.ProcMaskSpellDamage

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 1118 + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		spell.DealDamage(sim, result)
		unleashFlameAura.Activate(sim)
	}

	shaman.UnleashFlame = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashFrost(unleashElementsTimer *core.Timer) {
	spellCoeff := 0.386

	config := shaman.newUnleashElementsSpellConfig(unleashElementsTimer)
	config.SpellSchool = core.SpellSchoolFrost
	config.ProcMask = core.ProcMaskSpellDamage

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 873 + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		spell.DealDamage(sim, result)
	}

	shaman.UnleashFrost = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashWind(unleashElementsTimer *core.Timer) {

	unleashWindAura := shaman.RegisterAura(core.Aura{
		Label:     "Unleash Wind",
		ActionID:  core.ActionID{SpellID: 73681},
		Duration:  time.Second * 12,
		MaxStacks: 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1.4)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1/1.4)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				aura.RemoveStack(sim)
			}
		},
	})

	config := shaman.newUnleashElementsSpellConfig(unleashElementsTimer)
	config.SpellSchool = core.SpellSchoolPhysical
	config.ProcMask = core.ProcMaskMeleeSpecial

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		// TODO: 175% weapon damage
		unleashWindAura.Activate(sim)
	}

	shaman.UnleashWind = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashLife(unleashElementsTimer *core.Timer) {
	spellCoeff := 0.201
	//TODO: does this benefit from shaman.Talents.TidalWaves?
	//bonusCoeff := 0.02 * float64(shaman.Talents.TidalWaves)

	var affectedSpells []*core.Spell

	unleashLifeAura := shaman.RegisterAura(core.Aura{
		Label:    "Unleash Life",
		ActionID: core.ActionID{SpellID: 73685},
		Duration: time.Second * 8,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: need to confirm all spells that benefit from this
			affectedSpells = core.FilterSlice([]*core.Spell{
				shaman.ChainHeal,
				shaman.HealingWave,
				shaman.HealingSurge,
				shaman.GreaterHealingWave,
				shaman.Riptide,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.2
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.CostMultiplier -= 0.2
			}
		},
	})

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
		unleashLifeAura.Activate(sim)
	}

	shaman.UnleashLife = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerUnleashElements() {
	unleashElementsTimer := shaman.NewTimer()
	shaman.registerUnleashFlame(unleashElementsTimer)
	shaman.registerUnleashFrost(unleashElementsTimer)
	shaman.registerUnleashWind(unleashElementsTimer)
	shaman.registerUnleashLife(unleashElementsTimer)
}
