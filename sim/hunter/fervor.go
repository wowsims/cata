package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerFervorSpell() {
	if !hunter.Talents.Fervor {
		return
	}
	actionID := core.ActionID{SpellID: 82726}

	focusMetrics := hunter.NewFocusMetrics(actionID)
	fervorSpell := hunter.RegisterSpell(core.SpellConfig{
		ClassSpellMask: HunterSpellFervor,
		ActionID:       actionID,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
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
			if hunter.Pet != nil {
				hunter.Pet.AddFocus(sim, 50, focusMetrics)
			}
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 10,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					hunter.AddFocus(sim, 5, focusMetrics)
					hunter.Pet.AddFocus(sim, 5, focusMetrics)
				},
			})
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: fervorSpell,
		Type:  core.CooldownTypeDPS,
	})
}
