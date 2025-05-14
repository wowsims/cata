package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (paladin *Paladin) applyProtectionTalents() {
	paladin.applySealsOfThePure()
	paladin.applyToughness()
	paladin.applyHallowedGround()
	paladin.applySanctuary()
	paladin.applyWrathOfTheLightbringer()
	paladin.applyHammerOfTheRighteous()
	paladin.applyReckoning()
	paladin.applyShieldOfTheRighteous()
	paladin.applyGrandCrusader()
	paladin.applyHolyShield()
	paladin.applySacredDuty()
	paladin.applyShieldOfTheTemplar()
	paladin.applyArdentDefender()
}

func (paladin *Paladin) applySealsOfThePure() {
	if paladin.Talents.SealsOfThePure == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskSealOfRighteousness | SpellMaskSealOfTruth | SpellMaskSealOfJustice | SpellMaskCensure,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.06 * float64(paladin.Talents.SealsOfThePure),
	})
}

func (paladin *Paladin) applyToughness() {
	if paladin.Talents.Toughness == 0 {
		return
	}

	paladin.ApplyEquipScaling(stats.Armor, []float64{1.0, 1.03, 1.06, 1.1}[paladin.Talents.Toughness])
}

func (paladin *Paladin) applyHallowedGround() {
	if paladin.Talents.HallowedGround == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskConsecration,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.2 * float64(paladin.Talents.HallowedGround),
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskConsecration,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -40 * paladin.Talents.HallowedGround,
	})
}

func (paladin *Paladin) applySanctuary() {
	if paladin.Talents.Sanctuary == 0 {
		return
	}

	paladin.PseudoStats.ReducedCritTakenChance += 0.02 * float64(paladin.Talents.Sanctuary)
	paladin.PseudoStats.DamageTakenMultiplier *= 1.0 - []float64{0, 0.03, 0.07, 0.1}[paladin.Talents.Sanctuary]

	manaReturnActionID := core.ActionID{SpellID: []int32{0, 57319, 84626, 84627}[paladin.Talents.Sanctuary]}
	manaMetrics := paladin.NewManaMetrics(manaReturnActionID)
	manaReturnPct := 0.01 * float64(paladin.Talents.Sanctuary)

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:     "Sanctuary" + paladin.Label,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeBlock | core.OutcomeDodge,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AddMana(sim, manaReturnPct*paladin.MaxMana(), manaMetrics)
		},
	})
}

func (paladin *Paladin) applyHammerOfTheRighteous() {
	if !paladin.Talents.HammerOfTheRighteous {
		return
	}

	aoeMinDamage, aoeMaxDamage :=
		core.CalcScalingSpellEffectVarianceMinMax(proto.Class_ClassPaladin, 0.70800000429, 0.40000000596)

	numTargets := paladin.Env.GetNumTargets()
	actionId := core.ActionID{SpellID: 53595}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)

	hammerOfTheRighteousAoe := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88263},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskHammerOfTheRighteousAoe,

		MaxRange: 8,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.RollWithLabel(aoeMinDamage, aoeMaxDamage, "Hammer of the Righteous"+paladin.Label) +
				0.18000000715*spell.MeleeAttackPower()
			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMagicCrit)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	paladin.HammerOfTheRighteous = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHammerOfTheRighteousMelee,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.BuilderCooldown(),
				Duration: paladin.sharedBuilderBaseCD,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.MainHand().HandType != proto.HandType_HandTypeTwoHand
		},

		DamageMultiplier: 0.3,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.GainHolyPower(sim, 1, hpMetrics)
				hammerOfTheRighteousAoe.Cast(sim, target)
			}

			spell.DealOutcome(sim, result)
		},
	})
}

func (paladin *Paladin) applyWrathOfTheLightbringer() {
	if paladin.Talents.WrathOfTheLightbringer == 0 {
		return
	}

	dmgIncrease := 0.5 * float64(paladin.Talents.WrathOfTheLightbringer)

	// For some reason, only Crusader Strike and JoT are additive, while the rest are multiplicative.
	// Dunno if this is actually correct but that's how simc does it.
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskCrusaderStrike | SpellMaskJudgementOfTruth,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: dmgIncrease,
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskJudgementOfJustice | SpellMaskJudgementOfInsight | SpellMaskJudgementOfRighteousness,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: dmgIncrease,
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskHammerOfWrath | SpellMaskHolyWrath,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 15 * float64(paladin.Talents.WrathOfTheLightbringer),
	})
}

func (paladin *Paladin) applyReckoning() {
	if paladin.Talents.Reckoning == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 20178}
	procChance := 0.1 * float64(paladin.Talents.Reckoning)

	var reckoningSpell *core.Spell

	procAura := paladin.RegisterAura(core.Aura{
		Label:     "Reckoning Proc" + paladin.Label,
		ActionID:  actionID,
		Duration:  time.Second * 8,
		MaxStacks: 4,

		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			config := *paladin.AutoAttacks.MHConfig()
			config.ActionID = actionID
			reckoningSpell = paladin.GetOrRegisterSpell(config)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == paladin.AutoAttacks.MHAuto() {
				reckoningSpell.Cast(sim, result.Target)
				aura.RemoveStack(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       "Reckoning" + paladin.Label,
		ProcMask:   core.ProcMaskMelee,
		ProcChance: procChance,
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeBlock,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			procAura.Activate(sim)
			procAura.SetStacks(sim, 4)
		},
	})
}

func (paladin *Paladin) applyShieldOfTheRighteous() {
	if !paladin.Talents.ShieldOfTheRighteous {
		return
	}

	actionId := core.ActionID{SpellID: 53600}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)

	shieldDmg := core.CalcScalingSpellAverageEffect(proto.Class_ClassPaladin, 0.59299999475)

	paladin.ShieldOfTheRighteous = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskShieldOfTheRighteous,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.GetHolyPowerValue() > 0
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := []float64{0, 1, 3, 6}[paladin.GetHolyPowerValue()] *
				(shieldDmg + 0.10000000149*spell.MeleeAttackPower())

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.SpendHolyPower(sim, hpMetrics)
			}

			spell.DealOutcome(sim, result)
		},
	})
}

func (paladin *Paladin) applyShieldOfTheTemplar() {
	if paladin.Talents.ShieldOfTheTemplar == 0 {
		return
	}

	actionId := core.ActionID{SpellID: 84854}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskGuardianOfAncientKings,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * time.Duration(40*paladin.Talents.ShieldOfTheTemplar)),
	})

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskAvengingWrath,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * time.Duration(20*paladin.Talents.ShieldOfTheTemplar)),
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Divine Plea Templar Effect" + paladin.Label,
		ActionID:       actionId,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskDivinePlea,
		ProcChance:     1,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.GainHolyPower(sim, 3, hpMetrics)
		},
	})

}

func (paladin *Paladin) applyGrandCrusader() {
	if paladin.Talents.GrandCrusader == 0 {
		return
	}

	paladin.GrandCrusaderAura = paladin.RegisterAura(core.Aura{
		Label:    "Grand Crusader (Proc)" + paladin.Label,
		ActionID: core.ActionID{SpellID: 85043},
		Duration: time.Second * 6,

		// Dummy effect. Implemented in avengers_shield.go

		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&SpellMaskAvengersShield != 0 {
				paladin.GrandCrusaderAura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Grand Crusader" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 85416},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskBuilder,
		ProcChance:     []float64{0, 0.05, 0.10}[paladin.Talents.GrandCrusader],
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.AvengersShield.CD.Reset()
			paladin.GrandCrusaderAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) applyHolyShield() {
	if !paladin.Talents.HolyShield {
		return
	}

	holyShieldAura := paladin.RegisterAura(core.Aura{
		Label:    "Holy Shield" + paladin.Label,
		ActionID: core.ActionID{SpellID: 20925},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.BlockDamageReduction += 0.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.BlockDamageReduction -= 0.2
		},
	})

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20925},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskHolyShield,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			holyShieldAura.Activate(sim)
		},
	})
}

// 25/50% chance on Judgement/AS to apply 100% crit to next SotR
func (paladin *Paladin) applySacredDuty() {
	if paladin.Talents.SacredDuty == 0 {
		return
	}

	critMod := paladin.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskShieldOfTheRighteous,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 100,
	})

	paladin.SacredDutyAura = paladin.RegisterAura(core.Aura{
		Label:    "Sacred Duty (Proc)" + paladin.Label,
		ActionID: core.ActionID{SpellID: 85433},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			critMod.Activate()
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			critMod.Deactivate()
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&SpellMaskShieldOfTheRighteous != 0 && result.DidCrit() {
				paladin.SacredDutyAura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Sacred Duty" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 53710},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskAvengersShield | SpellMaskJudgement,
		ProcChance:     []float64{0, 0.25, 0.50}[paladin.Talents.SacredDuty],
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.SacredDutyAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) applyArdentDefender() {
	if !paladin.Talents.ArdentDefender {
		return
	}

	actionID := core.ActionID{SpellID: 31850}

	adAura := paladin.RegisterAura(core.Aura{
		Label:    "Ardent Defender" + paladin.Label,
		ActionID: actionID,
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	ardentDefender := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		SpellSchool:    core.SpellSchoolHoly,
		ClassSpellMask: SpellMaskArdentDefender,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			adAura.Activate(sim)
		},
	})

	adHealAmount := 0.0

	// Spell to heal you when AD has procced; fire this before fatal damage so that a Death is not detected
	adHeal := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 66235},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful,

		CritMultiplier:   1,
		ThreatMultiplier: 0,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, &paladin.Unit, adHealAmount, spell.OutcomeHealing)
		},
	})

	// >= 15% hp, hit gets reduced so we end up at 15% without heal
	// < 15% hp, hit gets reduced to 0 and we heal the remaining health up to 15%
	paladin.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
		if adAura.IsActive() && result.Damage >= paladin.CurrentHealth() {
			maxHealth := paladin.MaxHealth()
			currentHealth := paladin.CurrentHealth()
			incomingDamage := result.Damage

			if currentHealth/maxHealth >= 0.15 {
				// Incoming attack gets reduced so we end up at 15% hp
				// TODO: Overkill counted as absorb but not as healing in logs
				result.Damage = currentHealth - maxHealth*0.15
				if sim.Log != nil {
					paladin.Log(sim, "Ardent Defender absorbed %.1f damage", incomingDamage-result.Damage)
				}
			} else {
				// Incoming attack gets reduced to 0
				// Heal up to 15% hp
				// TODO: Overkill counted as absorb but not as healing in logs
				result.Damage = 0
				adHealAmount = maxHealth*0.15 - currentHealth
				adHeal.Cast(sim, &paladin.Unit)
				if sim.Log != nil {
					paladin.Log(sim, "Ardent Defender absorbed %.1f damage and healed for %.1f", incomingDamage, adHealAmount)
				}
			}

			adAura.Deactivate(sim)
		}
	})

	if paladin.Spec == proto.Spec_SpecProtectionPaladin {
		paladin.AddMajorCooldown(core.MajorCooldown{
			Spell: ardentDefender,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
