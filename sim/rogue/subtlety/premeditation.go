package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (subRogue *SubtletyRogue) registerPremeditation() {
	comboMetrics := subRogue.NewComboPointMetrics(core.ActionID{SpellID: 14183})

	subRogue.Premeditation = subRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 14183},
		Flags:          core.SpellFlagAPL | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: rogue.RogueSpellPremeditation,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: 0,
				GCD:  0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    subRogue.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return subRogue.IsStealthed() || subRogue.HasActiveAura("Shadowmeld")
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			subRogue.AddComboPointsOrAnticipation(sim, 2, comboMetrics)
		},
	})

	subRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    subRogue.Premeditation,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityLow,
	})
}
