package death_knight

import (
	//"github.com/wowsims/mop/sim/core/proto"

	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (dk *DeathKnight) ApplyFrostTalents() {
	// Killing Machine
	dk.applyKillingMachine()
}

func (dk *DeathKnight) applyKillingMachine() {
	if dk.Talents.KillingMachine == 0 {
		return
	}

	kmMod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 100,
		ClassMask:  DeathKnightSpellObliterate | DeathKnightSpellFrostStrike,
	})

	kmAura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Killing Machine Proc",
		ActionID: core.ActionID{SpellID: 51124},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			kmMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			kmMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(DeathKnightSpellObliterate | DeathKnightSpellFrostStrike) {
				return
			}
			if !result.Landed() {
				return
			}
			if !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				return
			}
			aura.Deactivate(sim)
		},
	})

	// Dummy spell to react with triggers
	kmProcSpell := dk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51124},
		Flags:          core.SpellFlagNoLogs | core.SpellFlagNoMetrics,
		ClassSpellMask: DeathKnightSpellKillingMachine,
	})

	ppm := 2.0 * float64(dk.Talents.KillingMachine)
	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:     "Killing Machine",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeLanded,
		PPM:      ppm,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			kmAura.Activate(sim)
			kmProcSpell.Cast(sim, nil)
		},
	})
}
