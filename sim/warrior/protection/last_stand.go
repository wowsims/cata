package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerLastStand() {
	actionID := core.ActionID{SpellID: 12975}
	healthMetrics := war.NewHealthMetrics(actionID)

	var bonusHealth float64
	war.LastStandAura = war.RegisterAura(core.Aura{
		Label:    "Last Stand",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusHealth = war.MaxHealth() * 0.3
			war.AddStatsDynamic(sim, stats.Stats{stats.Health: bonusHealth})
			war.GainHealth(sim, bonusHealth, healthMetrics)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.AddStatsDynamic(sim, stats.Stats{stats.Health: -bonusHealth})
		},
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: warrior.SpellMaskLastStand,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !war.RallyingCryAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			war.LastStandAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}
