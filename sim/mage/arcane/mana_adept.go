package arcane

import "github.com/wowsims/mop/sim/core"

func (arcane *ArcaneMage) ApplyMastery() {
	arcaneMastery := arcane.AddDynamicMod(core.SpellModConfig{
		School: core.SpellSchoolArcane | core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolHoly | core.SpellSchoolNature | core.SpellSchoolShadow,
		Kind:   core.SpellMod_DamageDone_Pct,
	})

	arcane.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		arcaneMastery.UpdateFloatValue(arcane.ArcaneMasteryValue())
	})

	core.MakePermanent(arcane.GetOrRegisterAura(core.Aura{
		Label:    "Mana Adept",
		ActionID: core.ActionID{SpellID: 76547},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMastery.UpdateFloatValue(arcane.ArcaneMasteryValue())
			arcaneMastery.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMastery.Deactivate()
		},
	}))

	core.MakeProcTriggerAura(&arcane.Unit, core.ProcTrigger{
		Name:     "Arcane Mastery Mana Updater",
		Callback: core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			arcaneMastery.UpdateFloatValue(arcane.ArcaneMasteryValue())
		},
	})
}
