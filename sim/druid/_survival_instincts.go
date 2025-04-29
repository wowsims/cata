package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerSurvivalInstinctsCD() {
	if !druid.Talents.SurvivalInstincts {
		return
	}

	actionID := core.ActionID{SpellID: 61336}

	cdTimer := druid.NewTimer()
	cd := time.Minute * 3
	getDuration := func() time.Duration {
		return core.TernaryDuration(druid.T11Feral4pBonus.IsActive(), time.Second*18, time.Second*12)
	}

	druid.SurvivalInstinctsAura = druid.RegisterAura(core.Aura{
		Label:    "Survival Instincts",
		ActionID: actionID,
		Duration: getDuration(),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageTakenMultiplier *= 0.5
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.DamageTakenMultiplier /= 0.5
		},
	})

	druid.SurvivalInstincts = druid.RegisterSpell(Cat|Bear, core.SpellConfig{
		ActionID: actionID,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			druid.SurvivalInstinctsAura.Duration = getDuration()
			druid.SurvivalInstinctsAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.SurvivalInstincts.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
