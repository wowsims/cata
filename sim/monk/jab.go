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
func (monk *Monk) registerJab() {
	actionID := core.ActionID{SpellID: 100780}
	chiMetrics := monk.NewChiMetrics(actionID)

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellJab,
		MaxRange:       core.MaxMeleeRange,

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

		DamageMultiplier: 1.5,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				chiGain := core.TernaryInt32(monk.StanceMatches(FierceTiger), 2, 1)
				monk.AddChi(sim, spell, chiGain, chiMetrics)
			}
		},
	})
}
