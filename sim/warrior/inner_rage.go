package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warrior *Warrior) RegisterInnerRage() {
	actionID := core.ActionID{SpellID: 1134}

	costMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskHeroicStrike | SpellMaskCleave,
		Kind:       core.SpellMod_Cooldown_Multiplier,
		FloatValue: -0.5,
	})

	warrior.InnerRageAura = warrior.RegisterAura(core.Aura{
		Label:    "Inner Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Deactivate()
		},
	})

	ir := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagMCD | core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskInnerRage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ThreatMultiplier: 0.0,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warrior.InnerRageAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ir,
		Type:  core.CooldownTypeDPS,
	})
}
