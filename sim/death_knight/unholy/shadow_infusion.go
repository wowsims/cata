package unholy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Your successful Death Coils empower your active Ghoul, increasing its damage dealt by 10% for 30 sec.
Stacks up to 5 times.
*/
func (uhdk *UnholyDeathKnight) registerShadowInfusion() {
	damageMod := uhdk.Ghoul.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.1,
	})

	actionID := core.ActionID{SpellID: 91342}
	uhdk.Ghoul.ShadowInfusionAura = uhdk.Ghoul.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Infusion" + uhdk.Ghoul.Label,
		ActionID:  actionID,
		Duration:  time.Second * 30,
		MaxStacks: 5,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damageMod.UpdateFloatValue(float64(newStacks) * 0.1)
		},
	}).AttachDependentAura(uhdk.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Infusion" + uhdk.Label,
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 5,
	}))

	core.MakeProcTriggerAura(&uhdk.Unit, core.ProcTrigger{
		Name:           "Shadow Infusion Trigger" + uhdk.Label,
		ActionID:       core.ActionID{SpellID: 49572},
		Callback:       core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
		ClassSpellMask: death_knight.DeathKnightSpellDeathCoil | death_knight.DeathKnightSpellDeathCoilHeal,
		Outcome:        core.OutcomeLanded,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return !uhdk.Ghoul.DarkTransformationAura.IsActive()
		},

		Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			uhdk.Ghoul.ShadowInfusionAura.Activate(sim)
			uhdk.Ghoul.ShadowInfusionAura.AddStack(sim)
		},
	})
}
