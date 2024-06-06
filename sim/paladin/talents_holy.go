package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
	"time"
)

func (paladin *Paladin) applyHolyTalents() {
	paladin.applyArbiterOfTheLight()
	paladin.applyProtectorOfTheInnocent()
	paladin.applyJudgementsOfThePure()
	paladin.applyBlazingLight()
	paladin.applyDenounce()
}

func (paladin *Paladin) applyArbiterOfTheLight() {
	if paladin.Talents.ArbiterOfTheLight == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskJudgement,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 6 * float64(paladin.Talents.ArbiterOfTheLight) * core.CritRatingPerCritChance,
	})
}

func (paladin *Paladin) applyProtectorOfTheInnocent() {
	if paladin.Talents.ProtectorOfTheInnocent == 0 {
		return
	}

	// TODO: Implement as a aura
}

// Might need Rework
func (paladin *Paladin) applyJudgementsOfThePure() {
	if paladin.Talents.JudgementsOfThePure == 0 {
		return
	}

	actionId := core.ActionID{SpellID: 53657}

	hasteAmount := 3 * float64(paladin.Talents.JudgementsOfThePure) * core.HasteRatingPerHastePercent
	spiritRegenAmount := 0.1 * float64(paladin.Talents.JudgementsOfThePure)

	jotpAura := paladin.GetOrRegisterAura(core.Aura{
		Label:    "Judgements of the Pure",
		ActionID: actionId,
		Duration: 60 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.SpellHaste, hasteAmount)
			paladin.AddStatDynamic(sim, stats.MeleeHaste, hasteAmount)
			paladin.PseudoStats.SpiritRegenRateCombat += spiritRegenAmount
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.AddStatDynamic(sim, stats.SpellHaste, -hasteAmount)
			paladin.AddStatDynamic(sim, stats.MeleeHaste, -hasteAmount)
			paladin.PseudoStats.SpiritRegenRateCombat -= spiritRegenAmount
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

func (paladin *Paladin) applyBlazingLight() {
	if paladin.Talents.BlazingLight == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskExorcism | SpellMaskHolyShock,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(paladin.Talents.BlazingLight),
	})
}

func (paladin *Paladin) applyDenounce() {
	if paladin.Talents.Denounce == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskExorcism,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -([]float64{0, 0.38, 0.75}[paladin.Talents.Denounce]),
	})
}
