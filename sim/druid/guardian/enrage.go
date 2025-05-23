package guardian

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (bear *GuardianDruid) registerEnrageSpell() {
	actionID := core.ActionID{SpellID: 5229}
	rageMetrics := bear.NewRageMetrics(actionID)

	bear.EnrageAura = bear.RegisterAura(core.Aura{
		Label:    "Enrage",
		ActionID: actionID,
		Duration: 10 * time.Second + 1, // add 1 ns duration offset in order to guarantee that the final tick fires

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second,
				NumTicks: 10,
				Priority: core.ActionPriorityRegen,

				OnAction: func(sim *core.Simulation) {
					if aura.IsActive() {
						bear.AddRage(sim, 1, rageMetrics)
					}
				},
			})
		},
	})

	bear.BearFormAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		if !bear.Env.MeasuringStats {
			bear.EnrageAura.Deactivate(sim)
		}
	})

	bear.Enrage = bear.RegisterSpell(druid.Bear, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},

			CD: core.Cooldown{
				Timer:    bear.NewTimer(),
				Duration: time.Minute,
			},

			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bear.AddRage(sim, 20, rageMetrics)
			bear.EnrageAura.Activate(sim)
		},
	})

	bear.AddMajorCooldown(core.MajorCooldown{
		Spell: bear.Enrage.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
