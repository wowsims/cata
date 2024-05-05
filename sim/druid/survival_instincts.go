package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (druid *Druid) registerSurvivalInstinctsCD() {
	if !druid.Talents.SurvivalInstincts {
		return
	}

	actionID := core.ActionID{SpellID: 61336}

	cdTimer := druid.NewTimer()
	cd := time.Minute * 3

	druid.SurvivalInstinctsAura = druid.RegisterAura(core.Aura{
		Label:    "Survival Instincts",
		ActionID: actionID,
		Duration: core.TernaryDuration(druid.HasSetBonus(ItemSetStormridersBattlegarb, 4), time.Second * 18, time.Second * 12),
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
			druid.SurvivalInstinctsAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.SurvivalInstincts.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}
