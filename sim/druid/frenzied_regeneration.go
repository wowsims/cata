package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (druid *Druid) registerFrenziedRegenerationCD() {
	actionID := core.ActionID{SpellID: 22842}
	healthMetrics := druid.NewHealthMetrics(actionID)
	rageMetrics := druid.NewRageMetrics(actionID)

	cdTimer := druid.NewTimer()
	cd := time.Minute * 3
	isGlyphed := druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFrenziedRegeneration)
	healingMulti := core.TernaryFloat64(isGlyphed, 1.3, 1.0)

	var bonusHealth float64
	druid.FrenziedRegenerationAura = druid.RegisterAura(core.Aura{
		Label:    "Frenzied Regeneration",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.HealingTakenMultiplier *= healingMulti
			bonusHealth = druid.MaxHealth() * 0.3
			druid.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})

			if druid.CurrentHealth() < bonusHealth {
				druid.GainHealth(sim, bonusHealth-druid.CurrentHealth(), healthMetrics)
			}
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.HealingTakenMultiplier /= healingMulti
			druid.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})

			if druid.CurrentHealth() > druid.MaxHealth() {
				druid.RemoveHealth(sim, druid.CurrentHealth()-druid.MaxHealth())
			}
		},
	})

	druid.FrenziedRegeneration = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			druid.FrenziedRegenerationAura.Activate(sim)

			if isGlyphed {
				return
			}

			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 20,
				Period:   time.Second * 1,
				Priority: core.ActionPriorityDOT,
				OnAction: func(sim *core.Simulation) {
					rageDumped := min(druid.CurrentRage(), 10.0)
					healthGained := rageDumped * 0.3 / 100 * druid.MaxHealth() * druid.PseudoStats.HealingTakenMultiplier

					if druid.FrenziedRegenerationAura.IsActive() {
						druid.SpendRage(sim, rageDumped, rageMetrics)
						druid.GainHealth(sim, healthGained, healthMetrics)
					}
				},
			})
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.FrenziedRegeneration.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
