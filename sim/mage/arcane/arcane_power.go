package arcane

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/mage"
)

func (arcane *ArcaneMage) registerArcanePowerCD() {

	hasGlyph := !arcane.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcanePower)

	arcanePowerDamageMod := arcane.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellsAllDamaging,
		FloatValue: 0.15,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	arcanePowerCostMod := arcane.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellsAllDamaging,
		FloatValue: 0.1,
		Kind:       core.SpellMod_PowerCost_Pct,
	})

	arcanePowerDurationMod := arcane.AddDynamicMod(core.SpellModConfig{
		ClassMask: mage.MageSpellsAllDamaging,
		IntValue:  10,
		Kind:      core.SpellMod_BuffDuration_Flat,
	})

	arcanePowerCooldownMod := arcane.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellArcanePower,
		FloatValue: 1.0,
		Kind:       core.SpellMod_Cooldown_Multiplier,
	})

	actionID := core.ActionID{SpellID: 12472}
	arcane.arcanePowerAura = arcane.RegisterAura(core.Aura{
		Label:    "Icy Veins",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if hasGlyph {
				arcanePowerCooldownMod.Activate()
				arcanePowerDurationMod.Activate()
			}
			arcanePowerDamageMod.Activate()
			arcanePowerCostMod.Activate()
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			if hasGlyph {
				arcanePowerCooldownMod.Deactivate()
				arcanePowerDurationMod.Deactivate()
			}
			arcanePowerDamageMod.Deactivate()
			arcanePowerCostMod.Deactivate()
		},
	})

	arcane.arcanePower = arcane.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: mage.MageSpellArcanePower,
		Flags:          core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    arcane.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			arcane.arcanePowerAura.Activate(sim)
		},
	})

	arcane.AddMajorCooldown(core.MajorCooldown{
		Spell: arcane.arcanePower,
		Type:  core.CooldownTypeDPS,
	})
}
