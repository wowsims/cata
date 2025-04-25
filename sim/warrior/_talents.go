package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// Applies the effects of "common" talents: talents in the first two rows of each tree that any spec could theoretically take
// Because cata restricts you to 10 points in a different tree, anything more is inaccessible. The rest of the trees are handled in each
// spec's implementation
func (warrior *Warrior) ApplyCommonTalents() {
	warrior.applyArmsCommonTalents()
	warrior.applyFuryCommonTalents()
	warrior.applyProtectionCommonTalents()
}

func (warrior *Warrior) applyArmsCommonTalents() {
	warrior.applyWarAcademy()
	warrior.RegisterDeepWounds()
}

func (warrior *Warrior) applyFuryCommonTalents() {
	warrior.applyBattleTrance()
	warrior.applyCruelty()
	warrior.applyExecutioner()
}

func (warrior *Warrior) applyProtectionCommonTalents() {
	warrior.applyIncite()
	warrior.applyToughness()
	warrior.applyBloodAndThunder()
	warrior.applyShieldSpecialization()
	warrior.applyShieldMastery()
	warrior.applyHoldTheLine()
	warrior.applyGagOrder()
}

func (warrior *Warrior) applyWarAcademy() {
	if warrior.Talents.WarAcademy == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskMortalStrike |
			SpellMaskRagingBlow |
			SpellMaskDevastate |
			SpellMaskVictoryRush |
			SpellMaskSlam,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: (0.05 * float64(warrior.Talents.WarAcademy)),
	})
}

const battleTranceAffectedSpellsMask = SpellMaskCleave |
	SpellMaskColossusSmash |
	SpellMaskExecute |
	SpellMaskHeroicStrike |
	SpellMaskRend |
	SpellMaskShatteringThrow |
	SpellMaskSlam |
	SpellMaskSunderArmor |
	SpellMaskThunderClap |
	SpellMaskWhirlwind |
	SpellMaskShieldSlam |
	SpellMaskConcussionBlow |
	SpellMaskDevastate |
	SpellMaskShockwave |
	SpellMaskBloodthirst |
	SpellMaskRagingBlow |
	SpellMaskMortalStrike |
	SpellMaskBladestorm

func (warrior *Warrior) applyBattleTrance() {
	if warrior.Talents.BattleTrance == 0 {
		return
	}

	btMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask: battleTranceAffectedSpellsMask,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	})

	triggerSpellMask := SpellMaskBloodthirst | SpellMaskMortalStrike | SpellMaskShieldSlam

	actionID := core.ActionID{SpellID: 12964}
	btAura := warrior.RegisterAura(core.Aura{
		Label:    "Battle Trance",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			btMod.Activate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Battle Trance affects the spells that proc it, so make sure we don't eat the proc with the same attack
			// that just proced it
			if (spell.ClassSpellMask&triggerSpellMask) != 0 && aura.TimeActive(sim) == 0 {
				return
			}

			if (spell.ClassSpellMask & battleTranceAffectedSpellsMask) != 0 {
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			btMod.Deactivate()
		},
	})

	procChance := 0.05 * float64(warrior.Talents.BattleTrance)

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Battle Trance Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ProcChance:     procChance,
		ICD:            5 * time.Second,
		ClassSpellMask: triggerSpellMask,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			btAura.Activate(sim)
		},
	})
}

func (warrior *Warrior) applyCruelty() {
	if warrior.Talents.Cruelty == 0 {
		return
	}
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskBloodthirst | SpellMaskMortalStrike | SpellMaskShieldSlam,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5 * float64(warrior.Talents.Cruelty),
	})
}

func (warrior *Warrior) applyExecutioner() {
	if warrior.Talents.Executioner == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 90806}
	executionerBuff := warrior.RegisterAura(core.Aura{
		Label:     "Executioner",
		ActionID:  actionID,
		Duration:  time.Second * 9,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			oldSpeed := 0.05 * float64(oldStacks)
			newSpeed := 0.05 * float64(newStacks)
			aura.Unit.MultiplyMeleeSpeed(sim, (1.0+newSpeed)/(1.0+oldSpeed))
		},
	})

	procChance := 0.5 * float64(warrior.Talents.Executioner)
	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Executioner Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskExecute,
		ProcChance:     procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			executionerBuff.Activate(sim)
			executionerBuff.AddStack(sim)
		},
	})
}

func (warrior *Warrior) applyIncite() {
	if warrior.Talents.Incite == 0 {
		return
	}
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskHeroicStrike,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5 * float64(warrior.Talents.Incite),
	})

	inciteMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskHeroicStrike,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 200.0, // This is actually how Incite is implemented
	})

	actionID := core.ActionID{SpellID: 86627}
	var lastTriggerTime int64 = 0
	inciteAura := warrior.RegisterAura(core.Aura{
		Label:    "Incite",
		ActionID: actionID,
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			inciteMod.Activate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask&SpellMaskHeroicStrike) != 0 && result.DidCrit() && lastTriggerTime != int64(sim.CurrentTime) {
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			inciteMod.Deactivate()
		},
	})

	procChance := []float64{0.0, 0.33, 0.66, 1.0}[warrior.Talents.Incite]
	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Incite Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeCrit,
		ClassSpellMask: SpellMaskHeroicStrike,
		ProcChance:     procChance,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return !inciteAura.IsActive()
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			lastTriggerTime = int64(sim.CurrentTime)
			inciteAura.Activate(sim)
		},
	})

}

func (warrior *Warrior) applyToughness() {
	if warrior.Talents.Toughness == 0 {
		return
	}
	warrior.ApplyEquipScaling(stats.Armor, []float64{1.0, 1.03, 1.06, 1.1}[warrior.Talents.Toughness])
}

func (warrior *Warrior) applyBloodAndThunder() {
	if warrior.Talents.BloodAndThunder == 0 {
		return
	}
	procChance := 0.5 * float64(warrior.Talents.BloodAndThunder)
	var lastAppliedTime int64 = -1
	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Blood and Thunder Trigger",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskThunderClap,
		ProcChance:     procChance,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			// If the rend we're checking was applied this iteration, skip to avoid an explosion of B&T procs
			// (8 targets, Rend T1, TClap hits T1 + B&T applies Rend to 7 other targets, TClap hits T2 + applies Rend to 7 other targets, etc...)
			return result.Target.HasActiveAuraWithTag("Rend") && lastAppliedTime != int64(sim.CurrentTime)
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// B&T resnapshots all of the rends it applies and will overwrite "better" rends on any target the TC hits
			for _, target := range sim.Encounter.TargetUnits {
				rend := warrior.Rend.Dot(target)
				lastAppliedTime = int64(sim.CurrentTime)

				rend.Apply(sim)
				rend.TickOnce(sim)
			}
		},
	})
}

func (warrior *Warrior) applyShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}
	extraBlockRage := 5 * float64(warrior.Talents.ShieldSpecialization)

	metrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12725})
	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:     "Shield Specialization Rage Trigger",
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeBlock,
		ICD:      1500 * time.Millisecond,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, extraBlockRage, metrics)
		},
	})
}

func (warrior *Warrior) applyShieldMastery() {
	if warrior.Talents.ShieldMastery == 0 {
		return
	}
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskShieldBlock,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(-10*warrior.Talents.ShieldMastery) * time.Second,
	})

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskShieldWall,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(-30*warrior.Talents.ShieldMastery) * time.Second,
	})

	actionID := core.ActionID{SpellID: 84608}
	magicDamageReduction := 1.0 - []float64{0.0, 0.07, 0.14, 0.2}[warrior.Talents.ShieldMastery]
	sbMagicDamageReductionAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Mastery",
		ActionID: actionID,
		Duration: 6 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= magicDamageReduction
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= magicDamageReduction
			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= magicDamageReduction
		},
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:           "Shield Mastery Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskShieldBlock,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			sbMagicDamageReductionAura.Activate(sim)
		},
	})

}

func (warrior *Warrior) applyHoldTheLine() {
	if warrior.Talents.HoldTheLine == 0 {
		return
	}
	buff := warrior.RegisterAura(core.Aura{
		Label:    "Hold the Line",
		ActionID: core.ActionID{SpellID: 84621},
		Duration: 5 * time.Second * time.Duration(warrior.Talents.HoldTheLine),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.CriticalBlockChance[1] += 0.1
			warrior.AddStatDynamic(sim, stats.BlockPercent, 10)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.CriticalBlockChance[1] -= 0.1
			warrior.AddStatDynamic(sim, stats.BlockPercent, -10)
		},
	})

	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
		Name:     "Hold the Line Trigger",
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeParry,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buff.Activate(sim)
		},
	})
}

func (warrior *Warrior) applyGagOrder() {
	if warrior.Talents.GagOrder == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskHeroicThrow,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Duration(-15*warrior.Talents.GagOrder) * time.Second,
	})

}
