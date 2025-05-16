package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Tooltip:

You Jab the target, dealing ${1.5*$<low>} to ${1.5*$<high>} damage and generating

-- Stance of the Fierce Tiger --

	2

-- else --

	1

--

	Chi.
*/
var jabActionID = core.ActionID{SpellID: 100780}

func jabSpellConfig(monk *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       jabActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellJab,
		MaxRange:       core.MaxMeleeRange,

		EnergyCost: overrides.EnergyCost,
		ManaCost:   overrides.ManaCost,
		Cast:       overrides.Cast,

		DamageMultiplier: 1.5,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: overrides.ApplyEffects,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
		config.Flags ^= core.SpellFlagAPL
	}

	return config
}
func (monk *Monk) registerJab() {
	chiMetrics := monk.NewChiMetrics(jabActionID)

	monk.RegisterSpell(jabSpellConfig(monk, false, core.SpellConfig{
		EnergyCost: core.EnergyCostOptions{
			Cost:   core.TernaryInt32(monk.StanceMatches(WiseSerpent), 0, 40),
			Refund: 0.8,
		},
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(monk.StanceMatches(WiseSerpent), 8, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				chiGain := core.TernaryInt32(monk.StanceMatches(FierceTiger), 2, 1)
				monk.AddChi(sim, spell, chiGain, chiMetrics)
			}
		},
	}))
}

func (pet *StormEarthAndFirePet) registerSEFJab() {
	pet.RegisterSpell(jabSpellConfig(pet.owner, true, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := pet.owner.CalculateMonkStrikeDamage(sim, spell)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	}))
}
