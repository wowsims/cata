package fury

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *FuryWarrior) ApplyTalents() {
	war.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86526)
	war.Warrior.ApplyCommonTalents()

	war.RegisterDeathWish()
	war.RegisterRagingBlow()

	war.applyFlurry()
	war.applyEnrage()
	war.applyRampage()
	war.applyMeatCleaver()
	war.applyBloodsurge()
	war.applyIntensityRage()
	war.applySingleMindedFury()

	war.ApplyGlyphs()
}

func (war *FuryWarrior) applyFlurry() {
	if war.Talents.Flurry == 0 {
		return
	}

	atkSpeedBonus := 1.0 + []float64{0.0, 0.08, 0.16, 0.25}[war.Talents.Flurry]
	actionID := core.ActionID{SpellID: 12968}
	flurryAura := war.RegisterAura(core.Aura{
		Label:     "Flurry",
		ActionID:  actionID,
		Duration:  15 * time.Second,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			war.MultiplyMeleeSpeed(sim, atkSpeedBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.MultiplyMeleeSpeed(sim, 1.0/atkSpeedBonus)
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Flurry Trigger",
		ActionID: actionID,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrMeleeProc,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
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

func (war *FuryWarrior) applyEnrage() {
	if war.Talents.Enrage == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 14202}
	procChance := 0.03 * float64(war.Talents.Enrage)
	baseDamageBonus := []float64{0.0, 0.03, 0.07, 0.1}[war.Talents.Enrage]
	var bonusSnapshot float64
	enrageAura := war.RegisterAura(core.Aura{
		Label:    "Enrage",
		Tag:      warrior.EnrageTag,
		ActionID: actionID,
		Duration: 9 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bonusSnapshot = 1.0 + (baseDamageBonus * war.EnrageEffectMultiplier)
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= bonusSnapshot
			//core.RegisterPercentDamageModifierEffect(aura, bonusSnapshot)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			war.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= bonusSnapshot
		},
	})

	core.RegisterPercentDamageModifierEffect(enrageAura, 1+baseDamageBonus)

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:       "Enrage Trigger",
		ActionID:   actionID,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			enrageAura.Activate(sim)
		},
	})
}

func (war *FuryWarrior) applyRampage() {
	// Raid buff is handled in warrior.ApplyRaidBuffs

	war.AddStat(stats.PhysicalCritPercent, 2)
}

func (war *FuryWarrior) applyMeatCleaver() {
	if war.Talents.MeatCleaver == 0 {
		return
	}

	buffMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskCleave | warrior.SpellMaskWhirlwind | warrior.SpellMaskWhirlwindOh,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.0,
	})

	actionID := core.ActionID{SpellID: 85739}
	bonusPerStack := 0.05 * float64(war.Talents.MeatCleaver)
	buff := war.RegisterAura(core.Aura{
		Label:     "Meat Cleaver",
		ActionID:  actionID,
		Duration:  10 * time.Second,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks != 0 {
				bonus := bonusPerStack * float64(newStacks)
				buffMod.UpdateFloatValue(bonus)
				buffMod.Activate()
			} else {
				buffMod.Deactivate()
			}
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Meat Cleaver Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnCastComplete,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: warrior.SpellMaskCleave | warrior.SpellMaskWhirlwind,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buff.Activate(sim)
			buff.AddStack(sim)
		},
	})
}

func (war *FuryWarrior) applyBloodsurge() {
	if war.Talents.Bloodsurge == 0 {
		return
	}

	castTimeMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskSlam,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1.0,
	})

	costMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskSlam,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	})

	damageMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskSlam,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.2,
	})

	actionID := core.ActionID{SpellID: 46916}
	buff := war.RegisterAura(core.Aura{
		Label:    "Bloodsurge",
		ActionID: actionID,
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
			costMod.Activate()
			damageMod.Activate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask & warrior.SpellMaskSlam) != 0 {
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
			costMod.Deactivate()
			damageMod.Deactivate()
		},
	})

	procChance := 0.1 * float64(war.Talents.Bloodsurge)
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Bloodsurge Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: warrior.SpellMaskBloodthirst,
		ProcChance:     procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buff.Activate(sim)
		},
	})
}

func (war *FuryWarrior) applyIntensityRage() {
	if war.Talents.IntensifyRage == 0 {
		return
	}

	cdr := 1.0 - 0.1*float64(war.Talents.IntensifyRage)
	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskBerserkerRage | warrior.SpellMaskRecklessness | warrior.SpellMaskDeathWish,
		Kind:       core.SpellMod_Cooldown_Multiplier,
		FloatValue: cdr,
	})
}

func (war *FuryWarrior) applySingleMindedFury() {
	if !war.Talents.SingleMindedFury {
		return
	}

	war.PseudoStats.DamageDealtMultiplier *= 1.2
	// Slam's extra hit is handled in its implementation
}
