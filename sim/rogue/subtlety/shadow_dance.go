package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerShadowDanceCD() {
	actionID := core.ActionID{SpellID: 51713}

	ambushReduction := subRogue.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Flat,
		ClassMask: rogue.RogueSpellAmbush,
		IntValue:  -20,
	})

	subRogue.ShadowDanceAura = subRogue.RegisterAura(core.Aura{
		Label:    "Shadow Dance",
		ActionID: actionID,
		Duration: time.Second * 8,
		// Can now cast opening abilities outside of stealth
		// Covered in rogue.go by IsStealthed()
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ambushReduction.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ambushReduction.Deactivate()
		},
	})

	subRogue.ShadowDance = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: rogue.RogueSpellShadowDance,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    subRogue.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			subRogue.BreakStealth(sim)
			subRogue.ShadowDanceAura.Activate(sim)
		},
		RelatedSelfBuff: subRogue.ShadowDanceAura,
	})

	subRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    subRogue.ShadowDance,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
	})
}
