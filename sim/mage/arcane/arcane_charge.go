package arcane

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (arcane *ArcaneMage) registerArcaneCharges() {
	abDamageMod := arcane.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellArcaneBlast | mage.MageSpellArcaneBarrage | mage.MageSpellArcaneMissilesTick,
		FloatValue: 0.5 * arcane.T15_4PC_ArcaneChargeEffect,
		Kind:       core.SpellMod_DamageDone_Flat,
	})
	abCostMod := arcane.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellArcaneBlast,
		FloatValue: 1.5 * arcane.T15_4PC_ArcaneChargeEffect,
		Kind:       core.SpellMod_PowerCost_Pct,
	})

	arcane.ArcaneChargesAura = arcane.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Charges Aura",
		ActionID:  core.ActionID{SpellID: 36032},
		Duration:  time.Second * 10,
		MaxStacks: 4,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			abDamageMod.Activate()
			abCostMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			abDamageMod.Deactivate()
			abCostMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			stacks := float64(newStacks)
			abDamageMod.UpdateFloatValue(.5 * arcane.T15_4PC_ArcaneChargeEffect * stacks)
			abCostMod.UpdateFloatValue(1.5 * arcane.T15_4PC_ArcaneChargeEffect * stacks)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(mage.MageSpellArcaneBarrage | mage.MageSpellEvocation) {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&arcane.Unit, core.ProcTrigger{
		Name:           "Arcane Charge Arcane Explosion - Trigger",
		ClassSpellMask: mage.MageSpellArcaneExplosion,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			arcane.ArcaneChargesAura.Activate(sim)
			if sim.Proc(.3, "ArcaneChargesProc") {
				arcane.ArcaneChargesAura.AddStack(sim)
			}
		},
	})

}
