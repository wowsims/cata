package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerNaturesSwiftness() {
	actionID := core.ActionID{SpellID: 132158}
	cdTimer := druid.NewTimer()
	cd := time.Minute * 1

	castTimeModConfig := core.SpellModConfig{
		ClassMask:  DruidHealingNonInstantSpells,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
	}

	resourceCostModConfig := core.SpellModConfig{
		ClassMask:  DruidHealingNonInstantSpells,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -1,
	}

	healingPowerModConfig := core.SpellModConfig{
		ClassMask:  DruidHealingNonInstantSpells,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.5,
	}

	nsAura := druid.RegisterAura(core.Aura{
		Label:    "Nature's Swiftness",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(DruidSpellHealingTouch) {
				return
			}

			aura.Deactivate(sim)
			cdTimer.Set(sim.CurrentTime + cd)
			druid.UpdateMajorCooldowns()
		},
	}).AttachSpellMod(castTimeModConfig).AttachSpellMod(resourceCostModConfig).AttachSpellMod(healingPowerModConfig)

	druid.NaturesSwiftness = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:        actionID,
		Flags:           core.SpellFlagNoOnCastComplete,
		RelatedSelfBuff: nsAura,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.NaturesSwiftness.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
