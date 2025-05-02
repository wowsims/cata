package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) applyFervorCD() {
	if !hunter.Talents.Fervor {
		return
	}
	actionID := core.ActionID{SpellID: 82726}

	focusMetrics := hunter.NewFocusMetrics(actionID)
	fervorSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 30,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AddFocus(sim, 50, focusMetrics)
			hunter.Pet.AddFocus(sim, 50, focusMetrics)
			//Todo: Recurring over 10 sec
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: fervorSpell,
		Type:  core.CooldownTypeDPS,
	})
}
