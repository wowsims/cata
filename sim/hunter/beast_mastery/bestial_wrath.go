package beast_mastery

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (bmHunter *BeastMasteryHunter) registerBestialWrathCD() {
	if bmHunter.Pet == nil {
		return
	}

	bwCostMod := bmHunter.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  hunter.HunterSpellsAll,
		FloatValue: -0.5,
	})

	actionID := core.ActionID{SpellID: 19574}

	bestialWrathPetAura := bmHunter.Pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})

	bestialWrathAura := bmHunter.RegisterAura(core.Aura{
		Label:    "Bestial Wrath",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
			bwCostMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
			bwCostMod.Deactivate()
		},
	})
	core.RegisterPercentDamageModifierEffect(bestialWrathAura, 1.1)

	bwSpell := bmHunter.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: hunter.HunterSpellBestialWrath,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 1,
			},
			CD: core.Cooldown{
				Timer:    bmHunter.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bestialWrathPetAura.Activate(sim)
			bestialWrathAura.Activate(sim)
		},
	})

	bmHunter.AddMajorCooldown(core.MajorCooldown{
		Spell: bwSpell,
		Type:  core.CooldownTypeDPS,
	})
}
