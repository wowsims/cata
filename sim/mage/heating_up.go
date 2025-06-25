package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// https://www.wowhead.com/mop-classic/spell=48107/heating-up#comments:id=1709419 For Information on heating up time specifics (.75s, .25s etc)

func (mage *Mage) registerHeatingUp() {

	mage.HeatingUp = mage.RegisterAura(core.Aura{
		Label:    "Heating Up",
		ActionID: core.ActionID{SpellID: 48107},
		Duration: time.Second * 10,
	})

	mage.InstantPyroblastAura = mage.RegisterAura(core.Aura{
		Label:    "Pyroblast!",
		ActionID: core.ActionID{SpellID: 48108},
		Duration: time.Second * 15,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -2.0,
		ClassMask:  MageSpellPyroblast,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
		ClassMask:  MageSpellPyroblast,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: .25,
		ClassMask:  MageSpellPyroblast | MageSpellPyroblastDot,
	})

}

func (mage *Mage) HeatingUpSpellHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, callback func()) {
	if spell.TravelTime() > time.Duration(FireSpellMaxTimeUntilResult) {
		core.StartDelayedAction(sim, core.DelayedActionOptions{
			DoAt: sim.CurrentTime + time.Duration(FireSpellMaxTimeUntilResult),
			OnAction: func(s *core.Simulation) {
				callback()
				mage.HandleHeatingUp(sim, spell, result)
			},
		})
	} else {
		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			callback()
			mage.HandleHeatingUp(sim, spell, result)
		})
	}
}

func (mage *Mage) HandleHeatingUp(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	if result.DidCrit() {
		if mage.HeatingUp.IsActive() {
			mage.InstantPyroblastAura.Activate(sim)
			mage.HeatingUp.Deactivate(sim)
		} else {
			mage.HeatingUp.Activate(sim)
		}
	} else {
		core.StartDelayedAction(sim, core.DelayedActionOptions{
			DoAt: sim.CurrentTime + time.Duration(HeatingUpDeactivateBuffer),
			OnAction: func(s *core.Simulation) {
				mage.HeatingUp.Deactivate(sim)
			},
		})
	}
}
