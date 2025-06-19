package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Increases the damage dealt by your Blood Boil by 10%, and when you land a melee attack on a target that is infected with your Blood Plague, there is a 10% chance that your next Blood Boil or Death and Decay will consume no runes.
(Proc chance: 10%)
*/
func (bdk *BloodDeathKnight) registerCrimsonScourge() {
	var crimsonScourgeAura *core.Aura
	crimsonScourgeAura = bdk.RegisterAura(core.Aura{
		Label:    "Crimson Scourge" + bdk.Label,
		ActionID: core.ActionID{SpellID: 81141},
		Duration: time.Second * 15,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: death_knight.DeathKnightSpellBloodBoil | death_knight.DeathKnightSpellDeathAndDecay,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.CurCast.Cost <= 0
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			crimsonScourgeAura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  death_knight.DeathKnightSpellBloodBoil | death_knight.DeathKnightSpellDeathAndDecay,
		FloatValue: -2.0,
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:       "Crimson Scourge Trigger" + bdk.Label,
		ActionID:   core.ActionID{SpellID: 81136},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,

		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return bdk.BloodPlagueSpell.Dot(result.Target).IsActive()
		},

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			crimsonScourgeAura.Activate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.DeathKnightSpellBloodBoil,
		FloatValue: 0.1,
	})
}
