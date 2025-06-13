package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerNaturesSwiftness() {
	actionID := core.ActionID{SpellID: 132158}
	cdTimer := druid.NewTimer()
	cd := time.Minute * 1

	htCastTimeMod := druid.AddDynamicMod(core.SpellModConfig{
		ClassMask:  DruidSpellHealingTouch,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
	})

	nsAura := druid.RegisterAura(core.Aura{
		Label:    "Nature's Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			htCastTimeMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			htCastTimeMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(DruidSpellHealingTouch) {
				return
			}

			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			druid.UpdateMajorCooldowns()
		},
	})

	druid.NaturesSwiftness = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			nsAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.NaturesSwiftness.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
