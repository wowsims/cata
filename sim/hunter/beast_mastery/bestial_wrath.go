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

	duration := core.TernaryDuration(bmHunter.CouldHaveSetBonus(hunter.YaunGolSlayersBattlegear, 4), 16, 10) * time.Second

	bwCostMod := bmHunter.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  hunter.HunterSpellsAll | hunter.HunterSpellsTalents,
		FloatValue: -0.5,
	})
	bwDamageMod := bmHunter.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.1,
	})
	bwPetDamageMod := bmHunter.Pet.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.2,
	})

	actionID := core.ActionID{SpellID: 19574}

	bestialWrathPetAura := bmHunter.Pet.RegisterAura(core.Aura{
		Label:    "Bestial Wrath Pet",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bwPetDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			bwPetDamageMod.Deactivate()
		},
	})

	bestialWrathAura := bmHunter.RegisterAura(core.Aura{
		Label:    "Bestial Wrath",
		ActionID: actionID,
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			bwDamageMod.Activate()
			bwCostMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			bwDamageMod.Deactivate()
			bwCostMod.Deactivate()
		},
	})
	core.RegisterPercentDamageModifierEffect(bestialWrathAura, 1.1)

	bwSpell := bmHunter.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: hunter.HunterSpellBestialWrath,
		Flags:          core.SpellFlagReadinessTrinket,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
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
