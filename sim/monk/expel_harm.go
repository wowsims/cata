package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// $stnc=$?a103985[${1.0}][${1.0}]
// $dwm1=$?a108561[${1}][${0.898882275}]
// $dwm=$?a115697[${1}][${$<dwm1>}]
// $bm=$?s120267[${0.4}][${1}]
// $offm1=$?a108561[${0}][${1}]
// $offm=$?a115697[${0}][${$<offm1>}]
// $apc=$?s120267[${$AP/11}][${$AP/14}]
// $offlow=$?!s124146[${$mwb/2/$mws}][${$owb/2/$ows}]
// $offhigh=$?!s124146[${$MWB/2/$mws}][${$OWB/2/$ows}]
// $low=${$<stnc>*($<bm>*$<dwm>*(($mwb)/($MWS)+$<offm>*$<offlow>)+$<apc>-1)}
// $high=${$<stnc>*($<bm>*$<dwm>*(($MWB)/($MWS)+$<offm>*$<offhigh>)+$<apc>+1)}

func (monk *Monk) registerExpelHarm() {
	actionID := core.ActionID{SpellID: 115072}
	chiMetrics := monk.NewChiMetrics(actionID)
	healingDone := 0.0

	expelHarmDamageSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 115129},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagPassiveSpell | core.SpellFlagIgnoreAttackerModifiers,
		MissileSpeed: 20,
		MaxRange:     10,

		DamageMultiplier: 0.5,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(s *core.Simulation) {
				spell.CalcAndDealDamage(sim, target, healingDone, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)
			})
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellExpelHarm,

		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryInt32(monk.StanceMatches(WiseSerpent), 0, 40),
		},
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(monk.StanceMatches(WiseSerpent), 2.5, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    monk.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		DamageMultiplier: 7,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)
			hpBefore := spell.Unit.CurrentHealth()
			// Can only target ourselves for now
			spell.CalcAndDealHealing(sim, spell.Unit, baseDamage, spell.OutcomeHealing)
			hpAfter := spell.Unit.CurrentHealth()
			healingDone = hpAfter - hpBefore

			if healingDone > 0 {
				// Should be the closest target
				expelHarmDamageSpell.Cast(sim, monk.CurrentTarget)
			}

			chiGain := core.TernaryInt32(monk.StanceMatches(FierceTiger), 2, 1)
			monk.AddChi(sim, spell, chiGain, chiMetrics)
		},
	})
}
