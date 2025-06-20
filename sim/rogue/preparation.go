package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerPreparationCD() {

	rogue.Preparation = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 14185},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: RogueSpellPreparation,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 5,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Spells affected by Preparation are:, Vanish (Overkill/Master of Subtlety), Sprint, Evasion, Dismantle
			rogue.Vanish.CD.Reset()
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Preparation,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
	})
}
