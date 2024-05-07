package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) registerSoulFireSpell() {
	var improvedSoulFire *core.Aura = nil
	if warlock.Talents.ImprovedSoulFire > 0 {
		damageBonus := 1 + .04*float64(warlock.Talents.ImprovedSoulFire)

		improvedSoulFire = warlock.RegisterAura(core.Aura{
			Label:    "Improved Soul Fire",
			ActionID: core.ActionID{SpellID: 18120},
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				//TODO: Add or mult?
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= damageBonus
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= damageBonus
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				//TODO: Add or mult?
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= damageBonus
				warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= damageBonus
			},
		})
	}

	warlock.SoulFire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6353},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSoulFire,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 4,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.726,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 2862, spell.OutcomeMagicHitAndCrit)
			if result.Landed() && improvedSoulFire != nil {
				improvedSoulFire.Activate(sim)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
