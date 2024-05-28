package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
	"time"
)

func (paladin *Paladin) ApplyHolyTalents() {
	paladin.ApplyArbiterOfTheLight()
	paladin.ApplyProtectorOfTheInnocent()
	paladin.ApplyJudgementsOfThePure()
	paladin.ApplyBlazingLight()
	paladin.ApplyDenounce()
}

func (paladin *Paladin) ApplyArbiterOfTheLight() {
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

func (paladin *Paladin) ApplyBlazingLight() {
	if paladin.Talents.BlazingLight == 0 {
		return
	}
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskExorcism | SpellMaskHolyShock,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(paladin.Talents.BlazingLight),
	})
}

func (paladin *Paladin) ApplyDenounce() {
	if paladin.Talents.Denounce == 0 {
		return
	}
	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskExorcism,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -core.TernaryFloat64(paladin.Talents.Denounce == 1, 0.38, 0.75),
	})
}
