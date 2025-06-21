package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (fdk *FrostDeathKnight) registerKillingMachine() {
	mask := death_knight.DeathKnightSpellObliterate | death_knight.DeathKnightSpellFrostStrike
	if fdk.CouldHaveSetBonus(death_knight.ItemSetBattleplateOfTheAllConsumingMaw, 4) {
		mask |= death_knight.DeathKnightSpellSoulReaper
	}

	var killingMachineAura *core.Aura
	killingMachineAura = fdk.RegisterAura(core.Aura{
		Label:    "Killing Machine" + fdk.Label,
		ActionID: core.ActionID{SpellID: 51124},
		Duration: time.Second * 10,
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: mask,
		ProcMask:       core.ProcMaskMeleeMH,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			killingMachineAura.Deactivate(sim)
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_BonusCrit_Percent,
		ClassMask:  mask,
		FloatValue: 100,
	})

	// Dummy spell to react with triggers
	kmProcSpell := fdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51124},
		Flags:          core.SpellFlagNoLogs | core.SpellFlagNoMetrics,
		ClassSpellMask: death_knight.DeathKnightSpellKillingMachine,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: killingMachineAura,
	})

	core.MakeProcTriggerAura(&fdk.Unit, core.ProcTrigger{
		Name:     "Killing Machine Trigger" + fdk.Label,
		ActionID: core.ActionID{SpellID: 51128},
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeWhiteHit,
		Outcome:  core.OutcomeLanded,
		DPM:      fdk.NewStaticLegacyPPMManager(6.0, core.ProcMaskMeleeWhiteHit),

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			kmProcSpell.Cast(sim, &fdk.Unit)
		},
	})
}
