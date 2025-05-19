package windwalker

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (ww *WindwalkerMonk) registerPassives() {
	ww.registerComboBreaker()
}

func (ww *WindwalkerMonk) registerComboBreaker() {
	registerComboBreakerAuraAndTrigger := func(labelSuffix string, spellID int32, triggerSpellMask int64) {
		aura := ww.RegisterAura(core.Aura{
			Label:    fmt.Sprintf("Combo Breaker: %s %s", labelSuffix, ww.Label),
			ActionID: core.ActionID{SpellID: spellID},
			Duration: time.Second * 20,

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.Matches(triggerSpellMask) || !result.Landed() {
					return
				}
				aura.Deactivate(sim)
			},
		})

		core.MakeProcTriggerAura(&ww.Unit, core.ProcTrigger{
			Name:           fmt.Sprintf("Combo Breaker: %s Trigger %s", labelSuffix, ww.Label),
			Callback:       core.CallbackOnSpellHitDealt,
			ClassSpellMask: monk.MonkSpellJab,
			Outcome:        core.OutcomeLanded,
			ProcChance:     0.12,

			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				aura.Activate(sim)
			},
		})
	}

	registerComboBreakerAuraAndTrigger(
		"Blackout Kick",
		116768,
		monk.MonkSpellBlackoutKick,
	)

	registerComboBreakerAuraAndTrigger(
		"Tiger Palm",
		118864,
		monk.MonkSpellTigerPalm,
	)
}
