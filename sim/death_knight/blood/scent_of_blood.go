package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (dk *BloodDeathKnight) registerScentOfBlood() {
	actionID := core.ActionID{SpellID: 50421}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	dsMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: death_knight.DeathKnightSpellDeathStrikeHeal,
	})

	var scentOfBloodAura *core.Aura
	scentOfBloodAura = dk.RegisterAura(core.Aura{
		Label:     "Scent of Blood" + dk.Label,
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
		Timer:    dk.NewTimer(),
		Duration: time.Second * 1,
	}

	scentOfBloodHandler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !icd.IsReady(sim) {
			return
		}

		if !sim.Proc(dk.AutoAttacks.MH().SwingSpeed/3.6, "Scent of Blood Proc") {
			return
		}

		icd.Use(sim)

		if !dk.Talents.Conversion || !dk.ConversionAura.IsActive() {
			dk.AddRunicPower(sim, 10.0, rpMetrics)
		}

		scentOfBloodAura.Activate(sim)
		scentOfBloodAura.AddStack(sim)
	}

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:     "Scent Of Blood Auto Trigger" + dk.Label,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeLanded,

		Handler: scentOfBloodHandler,
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:     "Scent Of Blood Avoidance Trigger" + dk.Label,
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeDodge | core.OutcomeParry,

		Handler: scentOfBloodHandler,
	})
}
