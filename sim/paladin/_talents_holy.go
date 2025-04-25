package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
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
		ClassMask:  SpellMaskJudgement | SpellMaskTemplarsVerdict,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 6 * float64(paladin.Talents.ArbiterOfTheLight),
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

	hasteMultiplier := 1 + 0.01*3*float64(paladin.Talents.JudgementsOfThePure)
	spiritRegenAmount := 0.1 * float64(paladin.Talents.JudgementsOfThePure)

	paladin.JudgementsOfThePureAura = paladin.GetOrRegisterAura(core.Aura{
		Label:    "Judgements of the Pure" + paladin.Label,
		ActionID: actionId,
		Duration: 60 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			paladin.MultiplyCastSpeed(hasteMultiplier)
			paladin.MultiplyMeleeSpeed(sim, hasteMultiplier)
			paladin.PseudoStats.SpiritRegenRateCombat += spiritRegenAmount
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.MultiplyCastSpeed(1 / hasteMultiplier)
			paladin.MultiplyMeleeSpeed(sim, 1/hasteMultiplier)
			paladin.PseudoStats.SpiritRegenRateCombat -= spiritRegenAmount
		},
	})

	core.MakeProcTriggerAura(&paladin.Unit, core.ProcTrigger{
		Name:           "Judgements of the Pure (Proc)" + paladin.Label,
		ActionID:       actionId,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: SpellMaskJudgement,

		ProcChance: 1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.JudgementsOfThePureAura.Activate(sim)
		},
	})
}

func (paladin *Paladin) applyBlazingLight() {
	if paladin.Talents.BlazingLight == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskExorcism | SpellMaskGlyphOfExorcism | SpellMaskHolyShock,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.1 * float64(paladin.Talents.BlazingLight),
	})
}

func (paladin *Paladin) applyDenounce() {
	if paladin.Talents.Denounce == 0 {
		return
	}

	paladin.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskExorcism,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -([]int32{0, 38, 75}[paladin.Talents.Denounce]),
	})
}
