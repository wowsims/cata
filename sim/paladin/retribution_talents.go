package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"time"
)

func (paladin *Paladin) ApplyRetributionTalents() {
	paladin.ApplyCrusade()
	paladin.ApplyRuleOfLaw()
	paladin.ApplySealsOfCommand()
	paladin.ApplySanctifiedWrath()
	paladin.ApplyCommunion()
	paladin.ApplyArtOfWar()
	paladin.ApplyDivinePurpose()
	paladin.ApplyInquiryOfFaith()
}

func (paladin *Paladin) ApplyCrusade() {
	if paladin.Talents.Crusade == 0 {
		return
	}
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskCrusaderStrike | SpellMaskDivineStorm | SpellMaskTemplarsVerdict | SpellMaskHolyShock,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(paladin.Talents.Crusade),
	})

	// TODO: Add Healing Mod for Holy Shock
}

func (paladin *Paladin) ApplyRuleOfLaw() {
	if paladin.Talents.RuleOfLaw == 0 {
		return
	}
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskCrusaderStrike | SpellMaskWordOfGlory | SpellMaskHammerOfTheRighteous,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5 * float64(paladin.Talents.RuleOfLaw) * core.CritRatingPerCritChance,
	})
}

func (paladin *Paladin) ApplySealsOfCommand() {
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
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
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

func (paladin *Paladin) ApplySanctifiedWrath() {
	if paladin.Talents.SanctifiedWrath == 0 {
		return
	}
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskHammerOfWrath,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 0.02 * float64(paladin.Talents.SanctifiedWrath) * core.CritRatingPerCritChance,
	})
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskAvengingWrath,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Second * 20 * time.Duration(paladin.Talents.SanctifiedWrath),
	})
}

func (paladin *Paladin) ApplyCommunion() {
	if !paladin.Talents.Communion {
		return
	}

	paladin.PseudoStats.DamageDealtMultiplier *= 1.02
}

func (paladin *Paladin) ApplyArtOfWar() {
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
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskExorcism,
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

func (paladin *Paladin) ApplyDivinePurpose() {
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
		ProcChance:     core.TernaryFloat64(paladin.Talents.DivinePurpose == 1, 0.07, 0.15),

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.DivinePurposeAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) ApplyInquiryOfFaith() {
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
