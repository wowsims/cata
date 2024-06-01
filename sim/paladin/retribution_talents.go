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

	actionId := core.ActionID{SpellID: 20424}

	// Seals of Command
	paladin.SealsOfCommand = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionId,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1.0,
		CritMultiplier:   paladin.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := .07 * spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:       "Seals of Command",
		ActionID:   actionId,
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskMeleeSpecial | core.ProcMaskMeleeWhiteHit,
		ProcChance: 1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.IsMelee() {
				paladin.SealsOfCommand.Cast(sim, result.Target)
			}
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

	paladin.ArtOfWarInstantCast = paladin.RegisterAura(core.Aura{
		Label:    "Art Of War",
		ActionID: core.ActionID{SpellID: 53488},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.Exorcism.CastTimeMultiplier -= 1.0
			paladin.Exorcism.CostMultiplier -= 1.0
			paladin.Exorcism.DamageMultiplierAdditive += 1.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.Exorcism.CastTimeMultiplier += 1.0
			paladin.Exorcism.CostMultiplier += 1.0
			paladin.Exorcism.DamageMultiplierAdditive -= 1.0
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == paladin.Exorcism {
				aura.Deactivate(sim)
			}
		},
	})

	artOfWarChance := []float64{0, 0.07, 0.14, 0.20}[paladin.Talents.TheArtOfWar]

	paladin.RegisterAura(core.Aura{
		Label:    "The Art of War",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("Art of War Proc") < artOfWarChance {
				paladin.ArtOfWarInstantCast.Activate(sim)
			}
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
