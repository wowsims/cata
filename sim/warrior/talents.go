package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
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
	warrior.applyWarMachine()
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

func (warrior *Warrior) applyWarMachine() {
	if warrior.Talents.WarAcademy == 0 {
		return
	}

	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskMortalStrike |
			SpellMaskRagingBlow |
			SpellMaskDevastate |
			SpellMaskVictoryRush |
			SpellMaskSlam,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 1.0 + (0.05 * float64(warrior.Talents.WarAcademy)),
	})
}

func (warrior *Warrior) applyBattleTrance() {
	if warrior.Talents.BattleTrance == 0 {
		return
	}

	var affectedSpellMask int64 = 0
	for _, spell := range warrior.Spellbook {
		if spell.DefaultCast.Cost > 5 && (spell.ClassSpellMask&SpellMaskSpecialAttack) != 0 {
			affectedSpellMask |= spell.ClassSpellMask
		}
	}

	// mask off the special attack bit so we don't accidentally match against everything
	affectedSpellMask &= ^SpellMaskSpecialAttack

	btMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask:  affectedSpellMask,
		Kind:       core.SpellMod_PowerCost_Flat,
		FloatValue: -5,
	})

	var lastTriggertime int64 = 0
	btAura := warrior.RegisterAura(core.Aura{
		Label:    "Battle Trance",
		ActionID: core.ActionID{SpellID: 12964},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			btMod.Activate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Battle Trance affects the spells that proc it, so make sure we don't eat the proc with the same attack
			// that just proced it
			if (spell.ClassSpellMask&affectedSpellMask) != 0 && lastTriggertime != int64(sim.CurrentTime) {
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			btMod.Deactivate()
		},
	})

	procChance := 0.05 * float64(warrior.Talents.BattleTrance)
	triggerSpellMask := SpellMaskBloodthirst | SpellMaskMortalStrike | SpellMaskShieldSlam
	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Battle Trance Trigger",
		Icd: &core.Cooldown{
			Timer:    warrior.NewTimer(),
			Duration: time.Second * 5,
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask & triggerSpellMask) == 0 {
				return
			}

			if !aura.Icd.IsReady(sim) {
				return
			}

			if sim.Proc(procChance, "Battle Trance Trigger") {
				aura.Icd.Use(sim)
				lastTriggertime = int64(sim.CurrentTime)
				btAura.Activate(sim)
			}
		},
	}))
}

func (warrior *Warrior) applyCruelty() {
	if warrior.Talents.Cruelty == 0 {
		return
	}
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskBloodthirst | SpellMaskMortalStrike | SpellMaskShieldSlam,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5 * float64(warrior.Talents.Cruelty) * core.CritRatingPerCritChance,
	})
}

func (warrior *Warrior) applyExecutioner() {
	if warrior.Talents.Executioner == 0 {
		return
	}
	executionerBuff := warrior.RegisterAura(core.Aura{
		Label:     "Executioner",
		ActionID:  core.ActionID{SpellID: 90806},
		Duration:  time.Second * 9,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			oldSpeed := 0.05 * float64(oldStacks)
			newSpeed := 0.05 * float64(newStacks)
			aura.Unit.MultiplyMeleeSpeed(sim, (1.0+newSpeed)/(1.0+oldSpeed))
		},
	})

	procChance := 0.5 * float64(warrior.Talents.Executioner)
	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Executioner Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if (spell.ClassSpellMask & SpellMaskExecute) == 0 {
				return
			}

			if sim.Proc(procChance, "Executioner Trigger") {
				executionerBuff.Activate(sim)
				executionerBuff.AddStack(sim)
			}
		},
	}))

}

func (warrior *Warrior) applyIncite() {
	if warrior.Talents.Incite == 0 {
		return
	}
	warrior.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskHeroicStrike,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5 * float64(warrior.Talents.Incite) * core.CritRatingPerCritChance,
	})

	inciteMod := warrior.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskHeroicStrike,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 200.0 * core.CritRatingPerCritChance, // This is actually how Incite is implemented
	})

	var lastTriggerTime int64 = 0
	inciteAura := warrior.RegisterAura(core.Aura{
		Label:    "Incite",
		ActionID: core.ActionID{SpellID: 86627},
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
	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Incite Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask & SpellMaskHeroicStrike) != 0 {
				if result.DidCrit() && !inciteAura.IsActive() {
					if sim.Proc(procChance, "Incite Trigger") {
						lastTriggerTime = int64(sim.CurrentTime)
						inciteAura.Activate(sim)
					}
				}
			}
		},
	}))

}

func (warrior *Warrior) applyToughness() {
	if warrior.Talents.Toughness == 0 {
		return
	}
	warrior.PseudoStats.ArmorMultiplier *= 1.0 + []float64{0.0, 0.03, 0.06, 0.1}[warrior.Talents.Toughness]
}

func (warrior *Warrior) applyBloodAndThunder() {
	if warrior.Talents.BloodAndThunder == 0 {
		return
	}
	procChance := 0.5 * float64(warrior.Talents.BloodAndThunder)
	var lastAppliedTime int64 = -1
	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Blood and Thunder Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask&SpellMaskThunderClap) != 0 && result.Target.HasActiveAuraWithTag("Rend") && sim.Proc(procChance, "Blood and Thunder") {

				// If the rend we're checking was applied this iteration, skip to avoid an explosion of B&T procs
				// (8 targets, Rend T1, TClap hits T1 + B&T applies Rend to 7 other targets, TClap hits T2 + applies Rend to 7 other targets, etc...)
				if lastAppliedTime == int64(sim.CurrentTime) {
					return
				}

				// B&T resnapshots all of the rends it applies and will overwrite "better" rends on any target the TC hits
				for _, target := range sim.Encounter.TargetUnits {
					rend := warrior.Rend.Dot(target)
					lastAppliedTime = int64(sim.CurrentTime)
					rend.Apply(sim)
					rend.TickOnce(sim)
				}
			}
		},
	}))

}

func (warrior *Warrior) applyShieldSpecialization() {
	if warrior.Talents.ShieldSpecialization == 0 {
		return
	}
	extraBlockRage := 5 * float64(warrior.Talents.ShieldSpecialization)

	metrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12725})
	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Shield Specialization Rage Trigger",
		Icd: &core.Cooldown{
			Timer:    warrior.NewTimer(),
			Duration: 1500 * time.Millisecond,
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !aura.Icd.IsReady(sim) {
				return
			}

			if result.Outcome.Matches(core.OutcomeBlock) {
				aura.Icd.Use(sim)
				warrior.AddRage(sim, extraBlockRage, metrics)
			}

			// TODO: Rage on spell reflect, if we ever decide to model that
		},
	}))

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

	magicDamageReduction := 1.0 - []float64{0.0, 0.07, 0.14, 0.2}[warrior.Talents.ShieldMastery]
	sbMagicDamageReductionAura := warrior.RegisterAura(core.Aura{
		Label:    "Shield Mastery",
		ActionID: core.ActionID{SpellID: 84608},
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

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Shield Mastery Trigger",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell.ClassSpellMask & SpellMaskShieldBlock) != 0 {
				sbMagicDamageReductionAura.Activate(sim)
			}
		},
	}))

}

func (warrior *Warrior) applyHoldTheLine() {
	if warrior.Talents.HoldTheLine == 0 {
		return
	}
	buff := warrior.RegisterAura(core.Aura{
		Label:    "Hold the Line",
		ActionID: core.ActionID{SpellID: 84621},
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.CriticalBlockChance += 0.1
			warrior.AddStatDynamic(sim, stats.Block, 10*core.BlockRatingPerBlockChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.CriticalBlockChance -= 0.1
			warrior.AddStatDynamic(sim, stats.Block, -10*core.BlockRatingPerBlockChance)
		},
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Hold the Line Trigger",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeParry) {
				buff.Activate(sim)
			}
		},
	}))

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
