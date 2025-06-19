package fury

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *FuryWarrior) registerCrazedBerserker() {
	// 2025-06-13 - Balance change
	// https://www.wowhead.com/blue-tracker/topic/eu/mists-of-pandaria-classic-development-notes-updated-6-june-571162
	patchedDamageMulti := 0.05
	war.AddStaticMod(core.SpellModConfig{
		Kind:     core.SpellMod_DamageDone_Pct,
		ProcMask: core.ProcMaskMeleeOH,
		// 2025-06-13 - Balance change
		// https://www.wowhead.com/blue-tracker/topic/eu/mists-of-pandaria-classic-development-notes-updated-6-june-571162
		FloatValue: 0.25 + patchedDamageMulti,
	})

	// 2025-06-13 - Balance change
	// https://www.wowhead.com/blue-tracker/topic/eu/mists-of-pandaria-classic-development-notes-updated-6-june-571162
	war.AutoAttacks.MHConfig().DamageMultiplier *= 1.1 + patchedDamageMulti
	war.AutoAttacks.OHConfig().DamageMultiplier *= 1.1 + patchedDamageMulti
}

func (war *FuryWarrior) registerFlurry() {

	flurryAura := war.RegisterAura(core.Aura{
		Label:     "Flurry",
		ActionID:  core.ActionID{SpellID: 12968},
		Duration:  15 * time.Second,
		MaxStacks: 3,
	}).AttachMultiplyMeleeSpeed(1.25)

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Flurry - Trigger",
		ActionID: core.ActionID{SpellID: 12972},
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.Proc(0.09, "Flurry") {
				flurryAura.Activate(sim)
				flurryAura.SetStacks(sim, flurryAura.MaxStacks)
				return
			}
			if flurryAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				flurryAura.RemoveStack(sim)
			}
		},
	})
}

func (war *FuryWarrior) registerBloodsurge() {
	actionID := core.ActionID{SpellID: 46916}

	war.BloodsurgeAura = war.RegisterAura(core.Aura{
		Label:     "Bloodsurge",
		ActionID:  actionID,
		Duration:  15 * time.Second,
		MaxStacks: 3,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskWildStrike,
		Kind:      core.SpellMod_PowerCost_Flat,
		IntValue:  -30,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskWildStrike,
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		TimeValue: time.Millisecond * -500,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Bloodsurge: Bloodthirst - Trigger",
		ClassSpellMask: warrior.SpellMaskBloodthirst,
		Outcome:        core.OutcomeLanded,
		Callback:       core.CallbackOnSpellHitDealt,
		ProcChance:     0.2,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.BloodsurgeAura.Activate(sim)
			war.BloodsurgeAura.SetStacks(sim, 3)
		},
	})
}

func (war *FuryWarrior) registerMeatCleaver() {
	actionID := core.ActionID{SpellID: 85739}

	war.MeatCleaverAura = war.RegisterAura(core.Aura{
		Label:     "Meat Cleaver",
		ActionID:  actionID,
		Duration:  10 * time.Second,
		MaxStacks: 3,
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Meat Cleaver: Whirlwind - Trigger",
		ClassSpellMask: warrior.SpellMaskWhirlwind,
		Outcome:        core.OutcomeLanded,
		Callback:       core.CallbackOnSpellHitDealt,
		ICD:            time.Millisecond * 500,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.MeatCleaverAura.Activate(sim)
			war.MeatCleaverAura.AddStack(sim)
		},
	})
}

func (war *FuryWarrior) registerSingleMindedFuryOrTitansGrip() {
	smf := war.RegisterAura(core.Aura{
		Label:    "Single-Minded Fury",
		ActionID: core.ActionID{SpellID: 81099},
		Duration: core.NeverExpires,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ProcMask:   core.ProcMaskMeleeOH,
		FloatValue: 0.35,
	}).AttachMultiplicativePseudoStatBuff(&war.Unit.PseudoStats.DamageDealtMultiplier, 1.35)

	tg := war.RegisterAura(core.Aura{
		Label:    "Titan's Grip",
		ActionID: core.ActionID{SpellID: 46917},
		Duration: core.NeverExpires,
	})

	if (war.MainHand().HandType == proto.HandType_HandTypeOneHand || war.MainHand().HandType == proto.HandType_HandTypeMainHand) &&
		(war.OffHand().HandType == proto.HandType_HandTypeOneHand || war.OffHand().HandType == proto.HandType_HandTypeOffHand) {
		core.MakePermanent(smf)
	} else {
		core.MakePermanent(tg)
	}
}

func (war *FuryWarrior) registerUnshackledFury() {
	masteryPoints := war.GetMasteryBonusMultiplier()
	masteryMod := war.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: masteryPoints,
		School:     core.SpellSchoolPhysical,
	})

	war.EnrageAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		masteryMod.Activate()
	})
	war.EnrageAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		masteryMod.Deactivate()
	})

	war.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryMod.UpdateFloatValue(war.GetMasteryBonusMultiplier())
	})
}
