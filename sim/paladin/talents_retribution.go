package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"time"
)

func (paladin *Paladin) applyRetributionTalents() {
	paladin.applyCrusade()
	paladin.applyRuleOfLaw()
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
		ClassMask:  SpellMaskCrusaderStrike | SpellMaskDivineStorm | SpellMaskTemplarsVerdict | SpellMaskHolyShock,
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
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5 * float64(paladin.Talents.RuleOfLaw) * core.CritRatingPerCritChance,
	})
}

func (paladin *Paladin) applySanctityOfBattle() {
	if !paladin.Talents.SanctityOfBattle {
		return
	}

	spenderCooldownMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: SpellMaskCrusaderStrike | SpellMaskDivineStorm,
	})

	updateTimeValue := func(castSpeed float64) {
		spenderCooldownMod.UpdateTimeValue(-time.Duration(4500 - 4500*castSpeed))
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
		ActionID:    core.ActionID{SpellID: 20424},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 0.07,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
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
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 2 * float64(paladin.Talents.SanctifiedWrath) * core.CritRatingPerCritChance,
	})
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskAvengingWrath,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -(time.Second * 20 * time.Duration(paladin.Talents.SanctifiedWrath)),
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
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  SpellMaskExorcism,
		FloatValue: -1.0,
	})

	exorcismDamageMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  SpellMaskExorcism | SpellMaskGlyphOfExorcism,
		FloatValue: 1.0,
	})

	artOfWarInstantCast := paladin.RegisterAura(core.Aura{
		Label:    "Art Of War",
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
		Name:       "The Art of War",
		ActionID:   core.ActionID{SpellID: 87138},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeWhiteHit,
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
		ClassSpellMask: SpellMaskDivineStorm | SpellMaskSpecialAttack,

		MaxRange: 8,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.05,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 4500 * time.Millisecond,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),

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

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}

			if numHits >= 4 {
				paladin.GainHolyPower(sim, 1, hpMetrics)
			}
		},
	})
}

func (paladin *Paladin) applyDivinePurpose() {
	if paladin.Talents.DivinePurpose == 0 {
		return
	}

	paladin.DivinePurposeAura = paladin.RegisterAura(core.Aura{
		Label:    "Divine Purpose",
		ActionID: core.ActionID{SpellID: 90174},
		Duration: time.Second * 8,

		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&SpellMaskTemplarsVerdict != 0 ||
				spell.ClassSpellMask&SpellMaskInquisition != 0 ||
				spell.ClassSpellMask&SpellMaskZealotry != 0 {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Divine Purpose",
		ActionID:       core.ActionID{SpellID: 86172},
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskCanTriggerDivinePurpose,
		ProcChance:     []float64{0, 0.07, 0.15}[paladin.Talents.DivinePurpose],

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.DivinePurposeAura.Activate(sim)
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

	paladin.ZealotryAura = paladin.RegisterAura(core.Aura{
		Label:    "Zealotry",
		ActionID: actionId,
		Duration: 20 * time.Second,
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
			paladin.ZealotryAura.Activate(sim)
		},
	})
}
