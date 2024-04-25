package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (druid *Druid) registerStarsurgeSpell() {
	spellCoeff := 1.228

	druid.Starsurge = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 98995},
		SpellSchool: core.SpellSchoolArcane | core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagOmenTrigger | core.SpellFlagAPL,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           druid.DefaultSpellCritMultiplier(),
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.17,
			Multiplier: 1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ThreatMultiplier: 1,

		BonusCoefficient: 1.228,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1363, 1752) + (spell.SpellPower() * spellCoeff)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
