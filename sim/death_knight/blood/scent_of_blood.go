package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Your successful main-hand autoattacks, dodges and parries have a chance to increase the healing and minimum healing done by your next Death Strike within 20 sec by 20%, and to generate 10 Runic Power.
This effect stacks up to 5 times.
(1s cooldown)
*/
func (bdk *BloodDeathKnight) registerScentOfBlood() {
	actionID := core.ActionID{SpellID: 50421}
	rpMetrics := bdk.NewRunicPowerMetrics(actionID)

	dsMod := bdk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: death_knight.DeathKnightSpellDeathStrikeHeal,
	})

	var scentOfBloodAura *core.Aura
	scentOfBloodAura = bdk.RegisterAura(core.Aura{
		Label:     "Scent of Blood" + bdk.Label,
		ActionID:  actionID,
		Duration:  time.Second * 20,
		MaxStacks: 5,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dsMod.Activate()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			dsMod.UpdateFloatValue(0.2 * float64(newStacks))
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dsMod.Deactivate()
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: death_knight.DeathKnightSpellDeathStrike,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			scentOfBloodAura.Deactivate(sim)
		},
	})

	icd := core.Cooldown{
		Timer:    bdk.NewTimer(),
		Duration: time.Second * 1,
	}

	scentOfBloodHandler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !icd.IsReady(sim) {
			return
		}

		if !sim.Proc(bdk.AutoAttacks.MH().SwingSpeed/3.6, "Scent of Blood Proc") {
			return
		}

		icd.Use(sim)

		if !bdk.Talents.Conversion || !bdk.ConversionAura.IsActive() {
			bdk.AddRunicPower(sim, 10.0, rpMetrics)
		}

		scentOfBloodAura.Activate(sim)
		scentOfBloodAura.AddStack(sim)
	}

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:     "Scent Of Blood Auto Trigger" + bdk.Label,
		ActionID: core.ActionID{SpellID: 148211},
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMHAuto,
		Outcome:  core.OutcomeLanded,

		Handler: scentOfBloodHandler,
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:     "Scent Of Blood Avoidance Trigger" + bdk.Label,
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeDodge | core.OutcomeParry,

		Handler: scentOfBloodHandler,
	})
}
