package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (monk *Monk) registerExpelHarm() {
	actionID := core.ActionID{SpellID: 115072}
	chiMetrics := monk.NewChiMetrics(actionID)
	healingDone := 0.0

	expelHarmDamageSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 115129},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagPassiveSpell,
		MissileSpeed: 20,
		MaxRange:     10,

		DamageMultiplier: 0.5,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, healingDone, spell.OutcomeMagicHitAndCrit)
		},
	})

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagHelpful | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellExpelHarm,

		EnergyCost: core.EnergyCostOptions{
			Cost:   core.TernaryInt32(monk.StanceMatches(WiseSerpent), 0, 40),
			Refund: 0.8,
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
		ThreatMultiplier: 1.0,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			hpBefore := target.GetStat(stats.Health)
			// Can only target ourselves for now
			result := spell.CalcHealing(sim, &monk.Unit, baseDamage, spell.OutcomeHealing)
			hpAfter := target.GetStat(stats.Health)
			healingDone = hpAfter - hpBefore

			if result.Landed() && healingDone > 0 {
				// Should be the closest target
				expelHarmDamageSpell.Cast(sim, monk.CurrentTarget)
			}

			chiGain := core.TernaryInt32(monk.StanceMatches(FierceTiger), 2, 1)
			monk.AddChi(sim, spell, chiGain, chiMetrics)
		},
	})
}
