package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (bdk *BloodDeathKnight) registerCrimsonScourge() {
	var crimsonScourgeAura *core.Aura
	crimsonScourgeAura = bdk.RegisterAura(core.Aura{
		Label:    "Crimson Scourge" + bdk.Label,
		ActionID: core.ActionID{SpellID: 81141},
		Duration: time.Second * 15,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: death_knight.DeathKnightSpellBloodBoil | death_knight.DeathKnightSpellDeathAndDecay,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.CurCast.Cost > 0 {
				return
			}

			crimsonScourgeAura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  death_knight.DeathKnightSpellBloodBoil | death_knight.DeathKnightSpellDeathAndDecay,
		FloatValue: -100,
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:     "Crimson Scourge Trigger" + bdk.Label,
		ActionID: core.ActionID{SpellID: 81136},
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !bdk.BloodPlagueSpell.Dot(result.Target).IsActive() {
				return
			}

			if !sim.Proc(0.1, "Crimson Scourge Proc") {
				return
			}

			crimsonScourgeAura.Activate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.DeathKnightSpellBloodBoil,
		FloatValue: 0.1,
	})
}
