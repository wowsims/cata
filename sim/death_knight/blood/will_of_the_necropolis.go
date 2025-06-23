package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
When a damaging attack brings you below 30% of your maximum health, the cooldown on your Rune Tap ability is refreshed and your next Rune Tap has no cost, and all damage taken is reduced by 25% for 8 sec.
This effect cannot occur more than once every 45 seconds.
(45s cooldown)
*/
func (bdk *BloodDeathKnight) registerWillOfTheNecropolis() {
	wotnDmgReductionAura := bdk.RegisterAura(core.Aura{
		Label:    "Will of The Necropolis Damage Reduction" + bdk.Label,
		ActionID: core.ActionID{SpellID: 81162},
		Duration: time.Second * 8,
	}).AttachMultiplicativePseudoStatBuff(
		&bdk.Unit.PseudoStats.DamageTakenMultiplier, 0.75,
	)

	var wotnRuneTapCostAura *core.Aura
	wotnRuneTapCostAura = bdk.RegisterAura(core.Aura{
		Label:    "Will of The Necropolis Rune Tap Cost" + bdk.Label,
		ActionID: core.ActionID{SpellID: 96171},
		Duration: time.Second * 8,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: death_knight.DeathKnightSpellRuneTap,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.CurCast.Cost <= 0
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			wotnRuneTapCostAura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  death_knight.DeathKnightSpellRuneTap,
		FloatValue: -2.0,
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:     "Will of The Necropolis Trigger" + bdk.Label,
		ActionID: core.ActionID{SpellID: 81164},
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeLanded,
		Harmful:  true,
		ICD:      time.Second * 45,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			maxHealth := bdk.MaxHealth()
			threshold := maxHealth * 0.3
			currentHealth := bdk.CurrentHealth()

			return currentHealth < threshold && currentHealth+result.Damage >= threshold
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			wotnDmgReductionAura.Activate(sim)
			wotnRuneTapCostAura.Activate(sim)
			bdk.RuneTapSpell.CD.Reset()
		},
	})
}
