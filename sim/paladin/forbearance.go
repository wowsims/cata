package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerForbearance() {
	forbearanceAuras := paladin.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		return unit.RegisterAura(core.Aura{
			Label:    "Forbearance" + unit.Label,
			ActionID: core.ActionID{SpellID: 25771},
			Duration: time.Second * 60,
		})
	})

	paladin.OnSpellRegistered(func(spell *core.Spell) {
		if spell.Matches(SpellMaskCausesForbearance) {
			oldCastCondition := spell.ExtraCastCondition
			spell.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
				if target.IsOpponent(&paladin.Unit) {
					target = &paladin.Unit
				}

				aura := forbearanceAuras.Get(target)
				if aura.IsActive() {
					return false
				}

				return oldCastCondition == nil || oldCastCondition(sim, target)
			}
		}
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Forbearance On Heal Dealt Trigger" + paladin.Label,
		Callback:       core.CallbackOnHealDealt,
		ClassSpellMask: SpellMaskHandOfProtection | SpellMaskLayOnHands,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			target := result.Target
			if target.IsOpponent(&paladin.Unit) {
				target = &paladin.Unit
			}

			forbearanceAuras.Get(target).Activate(sim)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Forbearance On Cast Complete Trigger" + paladin.Label,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskDivineShield,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			forbearanceAuras.Get(&paladin.Unit).Activate(sim)
		},
	})
}
