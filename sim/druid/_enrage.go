package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (druid *Druid) registerEnrageSpell() {
	actionID := core.ActionID{SpellID: 5229}
	rageMetrics := druid.NewRageMetrics(actionID)

	instantRage := 20.0
	primalMadnessRage := 6.0 * float64(druid.Talents.PrimalMadness)

	dmgBonus := 0.05 * float64(druid.Talents.KingOfTheJungle)

	druid.EnrageAura = druid.RegisterAura(core.Aura{
		Label:    "Enrage Aura",
		ActionID: actionID,
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.0 + dmgBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.0 + dmgBonus
		},
	})

	druid.Enrage = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.AddRage(sim, instantRage, rageMetrics)

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					if druid.EnrageAura.IsActive() {
						druid.AddRage(sim, 1, rageMetrics)
					}
				},
			})

			druid.EnrageAura.Activate(sim)

			if primalMadnessRage > 0 {
				druid.AddRage(sim, primalMadnessRage, druid.PrimalMadnessRageMetrics)
			}
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.Enrage.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
