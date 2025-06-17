package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerLayOnHands() {
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 633},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellHealing,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful | core.SpellFlagIgnoreModifiers,
		ClassSpellMask: SpellMaskLayOnHands,

		MaxRange: 40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 10,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   0,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if target.IsOpponent(&paladin.Unit) {
				target = &paladin.Unit
			}

			spell.CalcAndDealHealing(sim, target, paladin.MaxHealth(), spell.OutcomeHealing)
		},
	})
}
