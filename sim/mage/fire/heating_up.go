package fire

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerHeatingUp() {

	fire.heatingUp = fire.RegisterAura(core.Aura{
		Label:    "Heating Up",
		ActionID: core.ActionID{SpellID: 48107},
		Duration: time.Second * 10,
	})

	fire.pyroblastAura = fire.RegisterAura(core.Aura{
		Label:    "Pyroblast!",
		ActionID: core.ActionID{SpellID: 48108},
		Duration: time.Second * 15,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -2.0,
		ClassMask:  mage.MageSpellPyroblast,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
		ClassMask:  mage.MageSpellPyroblast,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: .25,
		ClassMask:  mage.MageSpellPyroblast | mage.MageSpellPyroblastDot,
	})

	core.MakeProcTriggerAura(&fire.Unit, core.ProcTrigger{
		Name:           "Heating Up/Pyroblast! - Trigger",
		ClassSpellMask: mage.MageSpellFireball | mage.MageSpellScorch | mage.MageSpellInfernoBlast | mage.MageSpellPyroblast | mage.MageSpellCombustion,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				if fire.heatingUp.IsActive() {
					fire.pyroblastAura.Activate(sim)
					fire.heatingUp.Deactivate(sim)
				} else {
					fire.heatingUp.Activate(sim)
				}
			} else {
				fire.heatingUp.Deactivate(sim)
			}

		},
	})
}
