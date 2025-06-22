package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

// Rather than update a variable somewhere for one effect (Fury's Unshackled Fury) just take a callback
// to fetch its multiplier when needed
type RageMultiplierCB func() float64

func (war *Warrior) registerBerserkerRage() {

	actionID := core.ActionID{SpellID: 18499}
	duration := time.Second * 6
	// 2025-06-13 - Balance change
	// https://www.wowhead.com/blue-tracker/topic/eu/mists-of-pandaria-classic-development-notes-updated-6-june-571162
	if war.Spec == proto.Spec_SpecFuryWarrior {
		duration = time.Second * 8
	}

	war.BerserkerRageAura = war.RegisterAura(core.Aura{
		Label:    "Berserker Rage",
		ActionID: actionID,
		Duration: duration,
	})

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskBerserkerRage,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {

			war.EnrageAura.Deactivate(sim)
			war.EnrageAura.Activate(sim)
			war.BerserkerRageAura.Activate(sim)
		},
	})
}
