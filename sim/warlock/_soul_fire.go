package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (warlock *Warlock) registerSoulFire() {
	var improvedSoulFire *core.Aura = nil
	if warlock.Talents.ImprovedSoulFire > 0 {
		damageBonus := 1 + .04*float64(warlock.Talents.ImprovedSoulFire)

		improvedSoulFire = warlock.RegisterAura(core.Aura{
			Label:    "Improved Soul Fire",
			ActionID: core.ActionID{SpellID: 18120},
			Duration: 20 * time.Second,
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

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6353},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSoulFire,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 9,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 4 * time.Second,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           warlock.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.72600001097,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := warlock.CalcAndRollDamageRange(sim, 2.54299998283, 0.22499999404)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if warlock.T13_4pc.IsActive() && warlock.SoulBurnAura.IsActive() && warlock.SoulShards.IsActive() {
				warlock.SoulShards.AddStack(sim)
			}

			warlock.SoulBurnAura.Deactivate(sim)
			if result.Landed() && improvedSoulFire != nil {
				improvedSoulFire.Activate(sim)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
