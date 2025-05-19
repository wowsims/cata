package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Tooltip:

You Jab the target, dealing ${1.5*$<low>} to ${1.5*$<high>} damage and generating

$stnc=$?a103985[${1.1}][${1.0}]
$dw1=$?a108561[${1}][${0.898882275}]
$dw=$?a115697[${1}][${$<dw1>}]
$bm=$?s120267[${0.85}][${1}]
$off1=$?a108561[${0}][${1}]
$off=$?a115697[${0}][${$<off1>}]
$offl=$?!s124146[${$mwb/2/$mws}][${$owb/2/$ows}]
$offh=$?!s124146[${$MWB/2/$mws}][${$OWB/2/$ows}]
$mist=$?a121278[${0.5}][${1}]
$low=${$<bm>*$<stnc>*($<dw>*(($mwb)/($MWS)*$<mist>+$<off>*$<offl>)+($AP/14)-1)}
$high=${$<bm>*$<stnc>*($<dw>*(($MWB)/($MWS)*$<mist>+$<off>*$<offh>)+($AP/14)+1)}

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
