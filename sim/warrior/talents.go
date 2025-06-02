package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warrior *Warrior) ApplyTalents() {
	// Level 15
	warrior.registerJuggernaut()

	// Level 30
	warrior.registerImpendingVictory()

	// Level 45

	// Level 60
	warrior.registerBladestorm()
	warrior.registerShockwave()

	// Level 75

	// Level 90
}

func (war *Warrior) registerJuggernaut() {
	if !war.Talents.Juggernaut {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskCharge,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -8 * time.Second,
	})
}

func (war *Warrior) registerImpendingVictory() {
	if !war.Talents.ImpendingVictory {
		return
	}

	actionID := core.ActionID{SpellID: 103840}
	healthMetrics := war.NewHealthMetrics(actionID)

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskImpendingVictory,

		RageCost: core.RageCostOptions{
			Cost:   10,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 56 + spell.MeleeAttackPower()*0.56
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				war.GainHealth(sim, war.MaxHealth()*0.2, healthMetrics)
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (war *Warrior) registerBladestorm() {
	if !war.Talents.Bladestorm {
		return
	}

	actionID := core.ActionID{SpellID: 46924}

	damageMultiplier := 1.2
	if war.Spec == proto.Spec_SpecArmsWarrior {
		damageMultiplier = 1.8
	} else if war.Spec == proto.Spec_SpecProtectionWarrior {
		damageMultiplier = 1.6
	}

	mhSpell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(1), // Real Spell ID: 50622
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskBladestormMH,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, enemyTarget := range sim.Encounter.ActiveTargets {
				baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, &enemyTarget.Unit, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	})

	ohSpell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(2), // Real Spell ID: 95738,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskBladestormOH,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: damageMultiplier,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, enemyTarget := range sim.Encounter.ActiveTargets {
				baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, &enemyTarget.Unit, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(0),
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskBladestorm,
		Flags:          core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 60,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultCritMultiplier(),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Bladestorm",
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mhSpell.Cast(sim, target)
				if war.OffHand() != nil && war.OffHand().WeaponType != proto.WeaponType_WeaponTypeUnknown {
					ohSpell.Cast(sim, target)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerShockwave() {
	if !war.Talents.Shockwave {
		return
	}

	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 46968},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: SpellMaskShockwave,
		Flags:          core.SpellFlagAoE | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: 40 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numLandedHits := 0
			baseDamage := 0.75 * spell.MeleeAttackPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

				if result.Landed() {
					numLandedHits++
				}
			}
			if numLandedHits >= 3 {
				spell.CD.Reduce(time.Second * 20)
			}
		},
	})
}

// func (warrior *Warrior) applyToughness() {
// 	if warrior.Talents.Toughness == 0 {
// 		return
// 	}
// 	warrior.ApplyEquipScaling(stats.Armor, []float64{1.0, 1.03, 1.06, 1.1}[warrior.Talents.Toughness])
// }

// func (warrior *Warrior) applyShieldSpecialization() {
// 	if warrior.Talents.ShieldSpecialization == 0 {
// 		return
// 	}
// 	extraBlockRage := 5 * float64(warrior.Talents.ShieldSpecialization)

// 	metrics := warrior.NewRageMetrics(core.ActionID{SpellID: 12725})
// 	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
// 		Name:     "Shield Specialization Rage Trigger",
// 		Callback: core.CallbackOnSpellHitTaken,
// 		Outcome:  core.OutcomeBlock,
// 		ICD:      1500 * time.Millisecond,
// 		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			warrior.AddRage(sim, extraBlockRage, metrics)
// 		},
// 	})
// }

// func (warrior *Warrior) applyShieldMastery() {
// 	if warrior.Talents.ShieldMastery == 0 {
// 		return
// 	}
// 	warrior.AddStaticMod(core.SpellModConfig{
// 		ClassMask: SpellMaskShieldBlock,
// 		Kind:      core.SpellMod_Cooldown_Flat,
// 		TimeValue: time.Duration(-10*warrior.Talents.ShieldMastery) * time.Second,
// 	})

// 	warrior.AddStaticMod(core.SpellModConfig{
// 		ClassMask: SpellMaskShieldWall,
// 		Kind:      core.SpellMod_Cooldown_Flat,
// 		TimeValue: time.Duration(-30*warrior.Talents.ShieldMastery) * time.Second,
// 	})

// 	actionID := core.ActionID{SpellID: 84608}
// 	magicDamageReduction := 1.0 - []float64{0.0, 0.07, 0.14, 0.2}[warrior.Talents.ShieldMastery]
// 	sbMagicDamageReductionAura := warrior.RegisterAura(core.Aura{
// 		Label:    "Shield Mastery",
// 		ActionID: actionID,
// 		Duration: 6 * time.Second,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= magicDamageReduction
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= magicDamageReduction
// 			warrior.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= magicDamageReduction
// 		},
// 	})

// 	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
// 		Name:           "Shield Mastery Trigger",
// 		ActionID:       actionID,
// 		Callback:       core.CallbackOnCastComplete,
// 		ClassSpellMask: SpellMaskShieldBlock,
// 		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			sbMagicDamageReductionAura.Activate(sim)
// 		},
// 	})

// }

// func (warrior *Warrior) applyHoldTheLine() {
// 	if warrior.Talents.HoldTheLine == 0 {
// 		return
// 	}
// 	buff := warrior.RegisterAura(core.Aura{
// 		Label:    "Hold the Line",
// 		ActionID: core.ActionID{SpellID: 84621},
// 		Duration: 5 * time.Second * time.Duration(warrior.Talents.HoldTheLine),
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			warrior.CriticalBlockChance[1] += 0.1
// 			warrior.AddStatDynamic(sim, stats.BlockPercent, 10)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			warrior.CriticalBlockChance[1] -= 0.1
// 			warrior.AddStatDynamic(sim, stats.BlockPercent, -10)
// 		},
// 	})

// 	core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
// 		Name:     "Hold the Line Trigger",
// 		Callback: core.CallbackOnSpellHitTaken,
// 		Outcome:  core.OutcomeParry,
// 		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			buff.Activate(sim)
// 		},
// 	})
// }

// func (warrior *Warrior) applyGagOrder() {
// 	if warrior.Talents.GagOrder == 0 {
// 		return
// 	}

// 	warrior.AddStaticMod(core.SpellModConfig{
// 		ClassMask: SpellMaskHeroicThrow,
// 		Kind:      core.SpellMod_Cooldown_Flat,
// 		TimeValue: time.Duration(-15*warrior.Talents.GagOrder) * time.Second,
// 	})

// }
