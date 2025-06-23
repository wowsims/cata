package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerVanishSpell() {
	rogue.Vanish = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1856},
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: RogueSpellVanish,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Pause auto attacks
			rogue.AutoAttacks.CancelAutoSwing(sim)
			// Apply stealth
			rogue.StealthAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Vanish,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDrums,
	})
}
