package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) applyRetributionTalents() {
	paladin.applyCrusade()
	paladin.applyRuleOfLaw()
	paladin.applyPursuitOfJustice()
	paladin.applySanctityOfBattle()
	paladin.applySealsOfCommand()
	paladin.applySanctifiedWrath()
	paladin.applyCommunion()
	paladin.applyArtOfWar()
	paladin.applyDivineStorm()
	paladin.applyDivinePurpose()
	paladin.applyInquiryOfFaith()
	paladin.applyZealotry()
}

func (paladin *Paladin) applyCrusade() {
	if paladin.Talents.Crusade == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskCrusaderStrike | SpellMaskHammerOfTheRighteous | SpellMaskTemplarsVerdict | SpellMaskHolyShock,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(paladin.Talents.Crusade),
	})

	// TODO: Add Healing Mod for Holy Shock if healing sim gets implemented
}

func (paladin *Paladin) applyRuleOfLaw() {
	if paladin.Talents.RuleOfLaw == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskCrusaderStrike | SpellMaskWordOfGlory | SpellMaskHammerOfTheRighteous,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5 * float64(paladin.Talents.RuleOfLaw),
	})
}

func (paladin *Paladin) applyPursuitOfJustice() {
	if paladin.Talents.PursuitOfJustice == 0 {
		return
	}

	spellID := []int32{0, 26022, 26023}[paladin.Talents.PursuitOfJustice]
	multiplier := []float64{0, 0.08, 0.15}[paladin.Talents.PursuitOfJustice]
	paladin.NewMovementSpeedAura("Pursuit of Justice", core.ActionID{SpellID: spellID}, multiplier)
}

func (paladin *Paladin) applySanctityOfBattle() {
	if !paladin.Talents.SanctityOfBattle {
		return
	}

	baseSpenderCooldown := float64(paladin.sharedBuilderBaseCD.Milliseconds())

	spenderCooldownMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: SpellMaskBuilder,
	})

	updateTimeValue := func(castSpeed float64) {
		spenderCooldownMod.UpdateTimeValue(-(time.Millisecond * time.Duration(baseSpenderCooldown-baseSpenderCooldown*castSpeed)))
	}

	paladin.AddOnCastSpeedChanged(func(_ float64, castSpeed float64) {
		updateTimeValue(castSpeed)
	})

	core.MakePermanent(paladin.GetOrRegisterAura(core.Aura{
		Label:    "Sanctity of Battle",
		ActionID: core.ActionID{SpellID: 25956},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			updateTimeValue(paladin.CastSpeed)
			spenderCooldownMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			spenderCooldownMod.Deactivate()
		},
	}))
}

func (paladin *Paladin) applySealsOfCommand() {
	if !paladin.Talents.SealsOfCommand {
		return
	}

	// Seals of Command
	sealsOfCommandProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 20424},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskSealsOfCommand,

		DamageMultiplier: 0.07,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Seals of Command",
		ActionID:       core.ActionID{SpellID: 85126},
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskSealOfTruth | SpellMaskSealOfRighteousness | SpellMaskSealOfJustice,
		ProcChance:     1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			sealsOfCommandProc.Cast(sim, result.Target)
		},
	})
}

func (paladin *Paladin) applySanctifiedWrath() {
	if paladin.Talents.SanctifiedWrath == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskHammerOfWrath,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 2 * float64(paladin.Talents.SanctifiedWrath),
	})
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskAvengingWrath,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * time.Duration(20*paladin.Talents.SanctifiedWrath)),
	})

	// Hammer of Wrath execute restriction removal is handled in hammer_of_wrath.go
}

func (paladin *Paladin) applyCommunion() {
	if !paladin.Talents.Communion {
		return
	}

	paladin.PseudoStats.DamageDealtMultiplier *= 1.02
}

func (paladin *Paladin) applyArtOfWar() {
	if paladin.Talents.TheArtOfWar == 0 {
		return
	}

	exorcismCastTimeMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  SpellMaskExorcism,
		FloatValue: -1.0,
	})

	exorcismCostMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: SpellMaskExorcism,
		IntValue:  -100,
	})

	exorcismDamageMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  SpellMaskExorcism | SpellMaskGlyphOfExorcism,
		FloatValue: 1.0,
	})

	artOfWarInstantCast := paladin.RegisterAura(core.Aura{
		Label:    "The Art Of War" + paladin.Label,
		ActionID: core.ActionID{SpellID: 59578},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			exorcismCastTimeMod.Activate()
			exorcismCostMod.Activate()
			exorcismDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			exorcismCastTimeMod.Deactivate()
			exorcismCostMod.Deactivate()
			exorcismDamageMod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&SpellMaskExorcism != 0 {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       "Art of War" + paladin.Label,
		ActionID:   core.ActionID{SpellID: 87138},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
		Outcome:    core.OutcomeLanded,
		ProcChance: []float64{0, 0.07, 0.14, 0.20}[paladin.Talents.TheArtOfWar],
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			artOfWarInstantCast.Activate(sim)
		},
	})
}

// Divine Storm is a non-ap normalised instant attack that has a weapon damage % modifier with a 1.0 coefficient.
// It does this damage to all targets in range.
// DS also heals up to 3 party or raid members for 25% of the total damage caused.
// The heal has threat implications, but given prot paladin cannot get enough talent
// points to take DS, we'll ignore it for now.
func (paladin *Paladin) applyDivineStorm() {
	if !paladin.Talents.DivineStorm {
		return
	}

	numTargets := paladin.Env.GetNumTargets()
	actionId := core.ActionID{SpellID: 53385}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)

	paladin.DivineStorm = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskDivineStorm,

		MaxRange: 8,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5,
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

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			numHits := 0
			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				result := spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				if result.Landed() {
					numHits += 1
				}
				results[idx] = result
			}

			if numHits >= 4 {
				paladin.GainHolyPower(sim, 1, hpMetrics)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})
}

func (paladin *Paladin) applyDivinePurpose() {
	if paladin.Talents.DivinePurpose == 0 {
		return
	}

	duration := time.Second * 8
	paladin.DivinePurposeAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Purpose" + paladin.Label,
		ActionID: core.ActionID{SpellID: 90174},
		Duration: duration,

		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&SpellMaskCanConsumeDivinePurpose != 0 && aura.RemainingDuration(sim) < duration {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := []float64{0, 0.07, 0.15}[paladin.Talents.DivinePurpose]
	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Divine Purpose (proc)" + paladin.Label,
		ActionID:       core.ActionID{SpellID: 90174},
		Callback:       core.CallbackOnSpellHitDealt | core.CallbackOnCastComplete,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskCanTriggerDivinePurpose,
		ProcChance:     1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result == nil && spell.ClassSpellMask&SpellMaskInquisition == 0 {
				return
			}

			if sim.Proc(procChance, "Divine Purpose"+paladin.Label) {
				paladin.DivinePurposeAura.Activate(sim)
			}
		},
	})
}

func (paladin *Paladin) applyInquiryOfFaith() {
	if paladin.Talents.InquiryOfFaith == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskCensure,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(paladin.Talents.InquiryOfFaith),
	})

	// Inquisition duration is handled in inquisition.go
}

func (paladin *Paladin) applyZealotry() {
	if !paladin.Talents.Zealotry {
		return
	}

	actionId := core.ActionID{SpellID: 85696}
	duration := time.Second * 20

	paladin.ZealotryAura = paladin.RegisterAura(core.Aura{
		Label:    "Zealotry" + paladin.Label,
		ActionID: actionId,
		Duration: duration,
		// Holy Power logic is handled for each ability
	})

	paladin.Zealotry = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		Flags:          core.SpellFlagAPL,
		ProcMask:       core.ProcMaskEmpty,
		SpellSchool:    core.SpellSchoolHoly,
		ClassSpellMask: SpellMaskZealotry,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 2 * time.Minute,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.GetHolyPowerValue() >= 3
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: paladin.ZealotryAura,
	})
}
