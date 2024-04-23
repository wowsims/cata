package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerIncinerateSpell() {
	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 29722},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   24,
		ClassSpellMask: WarlockSpellIncinerate,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				//TODO: Check cast time
				CastTime: time.Millisecond * 2500,
			},
		},

		BonusCoefficient:         0.53899997473,
		DamageMultiplierAdditive: 1 + warlock.GrandFirestoneBonus(),
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//TODO: Check Damage
			var baseDamage float64
			if warlock.Immolate.Dot(target).IsActive() {
				baseDamage = sim.Roll(582+145, 676+169)
			} else {
				baseDamage = sim.Roll(582, 676)
			}

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
