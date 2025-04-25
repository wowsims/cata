package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerPreparationCD() {
	if !subRogue.Talents.Preparation {
		return
	}

	subRogue.Preparation = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 14185},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellPreparation,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    subRogue.NewTimer(),
				Duration: time.Minute * 5,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			// Spells affected by Preparation are: Shadowstep, Vanish (Overkill/Master of Subtlety), Sprint
			// If Glyph of Preparation is applied, Smoke Bomb, Dismantle, and Kick are also affected
			subRogue.Shadowstep.CD.Reset()
			subRogue.Vanish.CD.Reset()
		},
	})

	subRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    subRogue.Preparation,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !subRogue.Vanish.CD.IsReady(sim)
		},
	})
}
