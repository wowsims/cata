package beast_mastery

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (bmHunter *BeastMasteryHunter) registerKillCommandSpell() {
	if bmHunter.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 34026}

	bmHunter.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMelee,
		ClassSpellMask: hunter.HunterSpellKillCommand,
		Flags:          core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    bmHunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           bmHunter.DefaultCritMultiplier(),
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if bmHunter.Pet != nil && bmHunter.Pet.KillCommand != nil {
				bmHunter.Pet.KillCommand.Cast(sim, target)
			}
		},
	})
}
