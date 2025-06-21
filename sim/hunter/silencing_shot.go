package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerSilencingShotSpell() {
	hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 34490},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		MinRange:    5,
		MaxRange:    40,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			focusMetics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34490})
			hunter.AddFocus(sim, 10, focusMetics)
		},
	})
}
