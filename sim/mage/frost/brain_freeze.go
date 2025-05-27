package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (frost *FrostMage) registerBrainFreeze() {
	/*
		https://www.wowhead.com/mop-classic/spell=44549/brain-freeze and https://www.wowhead.com/mop-classic/spell=44614/frostfire-bolt for more information.
	*/
	buff := frost.RegisterAura(core.Aura{
		Label:     "Brain Freeze",
		ActionID:  core.ActionID{SpellID: 44549},
		Duration:  time.Second * 15,
		MaxStacks: 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			frost.frostfireFrozenCritBuffMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			frost.frostfireFrozenCritBuffMod.Deactivate()
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -1,
		ClassMask:  mage.MageSpellFrostfireBolt,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1,
		ClassMask:  mage.MageSpellFrostfireBolt,
	})

	/*
		Shatter doubles the crit chance of spells against frozen targets and then adds an additional 50%, hence critChance * 2 + 50
		https://www.wowhead.com/mop-classic/spell=12982/shatter for more information.
	*/

	core.MakeProcTriggerAura(&frost.Unit, core.ProcTrigger{
		Name:           "Brain Freeze - Trigger",
		ClassSpellMask: mage.MageSpellLivingBombDot | mage.MageSpellLivingBombExplosion | mage.MageSpellFrostBomb,
		Callback:       core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// https://github.com/simulationcraft/simc/blob/e1190fed141feec2ec7a489e80caec5138c3a6ab/engine/class_modules/sc_mage.cpp#L4169
			if spell.Matches(mage.MageSpellLivingBombDot | mage.MageSpellLivingBombExplosion) {
				if sim.Proc(0.25, "BrainFreezeProc") {
					buff.Activate(sim)
				}
			} else if spell.Matches(mage.MageSpellFrostBomb) {
				if sim.Proc(1.0, "BrainFreezeProc") {
					buff.Activate(sim)
				}
			} else {
				if sim.Proc(0.09, "BrainFreezeProc") {
					buff.Activate(sim)
				}
			}
		},
	})

}
