package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
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
	critMod := frost.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageSpellFrostfireBolt,
		FloatValue: frost.GetStat(stats.SpellCritPercent)*2 + 50,
		Kind:       core.SpellMod_BonusCrit_Percent,
	})

	buff.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
		critMod.Activate()
	}

	buff.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
		critMod.Deactivate()
	})

	core.MakeProcTriggerAura(&frost.Unit, core.ProcTrigger{
		Name:           "Brain Freeze - Trigger",
		ClassSpellMask: mage.MageSpellFrostfireBolt,
		Callback:       core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buff.Activate(sim)
			buff.AddStack(sim)
		},
	})

}
