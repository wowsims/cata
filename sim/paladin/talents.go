package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (paladin *Paladin) ApplyTalents() {
	paladin.ApplyCrusade()
	paladin.ApplyRuleOfLaw()
	paladin.ApplySealsOfThePure()
	paladin.ApplyArbitorOfTheLight()
	paladin.ApplyProtectorOfTheInnocent()
	paladin.ApplyJudgementsOfThePure()
}

// Retribution Talents first two rows
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

// Protection Talents first two rows
func (paladin *Paladin) ApplySealsOfThePure() {
	if paladin.Talents.SealsOfThePure == 0 {
		return
	}
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskSealOfRighteousness | SpellMaskSealOfTruth | SpellmaskSealofJustice,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.06 * float64(paladin.Talents.SealsOfThePure),
	})
}

// Holy Talents first two rows
func (paladin *Paladin) ApplyArbitorOfTheLight() {
	if paladin.Talents.ArbiterOfTheLight == 0 {
		return
	}
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskJudgement,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 6 * float64(paladin.Talents.ArbiterOfTheLight) * core.CritRatingPerCritChance,
	})
}

func (paladin *Paladin) ApplyProtectorOfTheInnocent() {
	if paladin.Talents.ProtectorOfTheInnocent == 0 {
		return
	}
	// TODO: Implement as a aura
}

// Might need Rework
func (paladin *Paladin) ApplyJudgementsOfThePure() {
	if paladin.Talents.JudgementsOfThePure == 0 {
		return
	}
	actionId := core.ActionID{SpellID: 53657}

	hasteAmount := 3 * float64(paladin.Talents.JudgementsOfThePure) * core.HasteRatingPerHastePercent

	jotpAura := paladin.GetOrRegisterAura(core.Aura{
		Label:    "Judgements of the Pure",
		ActionID: actionId,
		Duration: 60 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {

			paladin.AddStatDynamic(sim, stats.SpellHaste, hasteAmount)
			paladin.AddStatDynamic(sim, stats.MeleeHaste, hasteAmount)

			// TODO: Add Spirit Mod
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.SpellHaste, -hasteAmount)
			paladin.AddStatDynamic(sim, stats.MeleeHaste, -hasteAmount)

			// Todo: Remove Spirit Mod
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Judgements of the Pure",
		ActionID:       actionId,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskJudgement,

		ProcChance: 1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			jotpAura.Activate(sim)
		},
	})
}
