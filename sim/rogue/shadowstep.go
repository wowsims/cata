package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerShadowstepCD() {
	actionID := core.ActionID{SpellID: 36554}

	rogue.ShadowstepAura = rogue.RegisterAura(core.Aura{
		Label:    "Shadowstep",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// TODO: Movement Speed?
		},
	})

	rogue.Shadowstep = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: RogueSpellShadowstep,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 24,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// TODO: Teleport?
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: rogue.ShadowstepAura,
	})
}
