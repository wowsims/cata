package mop

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {

	core.NewItemEffect(75274, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()
		duration := time.Second * 15

		statValue := core.GetItemEffectScaling(75274, 2.66700005531, state)

		auras := make(map[stats.Stat]*core.StatBuffAura, 2)
		auras[stats.Strength] = character.NewTemporaryStatsAura(
			"Strength",
			core.ActionID{SpellID: 60229},
			stats.Stats{stats.Strength: statValue},
			duration,
		)
		auras[stats.Agility] = character.NewTemporaryStatsAura(
			"Agility",
			core.ActionID{SpellID: 60233},
			stats.Stats{stats.Agility: statValue},
			duration,
		)
		auras[stats.Intellect] = character.NewTemporaryStatsAura(
			"Intellect",
			core.ActionID{SpellID: 60234},
			stats.Stats{stats.Intellect: statValue},
			duration,
		)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Zen Alchemist Stone",
			ActionID:   core.ActionID{SpellID: 105574},
			ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
			Harmful:    true,
			ICD:        time.Second * 55,
			ProcChance: 0.25,
			Outcome:    core.OutcomeLanded,
			Callback:   core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				auras[character.GetHighestStatType([]stats.Stat{stats.Strength, stats.Agility, stats.Intellect})].Activate(sim)
			},
		})

		for _, aura := range auras {
			character.AddStatProcBuff(75274, aura, false, core.TrinketSlots())
		}
	})

	core.NewItemEffect(81266, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 126467}
		manaMetrics := character.NewManaMetrics(actionID)

		mana := core.GetItemEffectScaling(81266, 2.97199988365, state)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Price of Progress (Heroic)",
			ActionID:   actionID,
			ProcMask:   core.ProcMaskSpellHealing,
			Harmful:    true,
			ICD:        time.Second * 55,
			ProcChance: 0.10,
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicHealDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if character.HasManaBar() {
					character.AddMana(sim, mana, manaMetrics)
				}
			},
		})
	})

	// Renataki's Soul Charm
	// Your attacks  have a chance to grant Blades of Renataki, granting 1592 Agility every 1 sec for 10 sec.  (Approximately 1.21 procs per
	shared.ItemVersionMap{
		shared.ItemVersionLFR:                 95625,
		shared.ItemVersionNormal:              94512,
		shared.ItemVersionHeroic:              96369,
		shared.ItemVersionThunderforged:       95997,
		shared.ItemVersionHeroicThunderforged: 96741,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		label := "Blades of Renataki"

		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()

			statValue := core.GetItemEffectScaling(itemID, 0.44999998808, state)

			statBuffAura, aura := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				AuraLabel:            fmt.Sprintf("%s %s", label, versionLabel),
				ActionID:             core.ActionID{SpellID: 138756},
				Duration:             time.Second * 20,
				MaxStacks:            10,
				TimePerStack:         time.Second * 1,
				BonusPerStack:        stats.Stats{stats.Agility: statValue},
				StackingAuraActionID: core.ActionID{SpellID: 138737},
				StackingAuraLabel:    fmt.Sprintf("Item - Proc Stacking Agility %s", versionLabel),
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:    label,
				Harmful: true,
				ICD:     time.Second * 10,
				DPM: character.NewRPPMProcManager(itemID, false, core.ProcMaskDirect|core.ProcMaskProc, core.RPPMConfig{
					PPM: 1.21000003815,
				}),
				Outcome:  core.OutcomeLanded,
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					aura.Activate(sim)
				},
			})

			character.AddStatProcBuff(itemID, statBuffAura, false, core.TrinketSlots())
		})
	})

	// Delicate Vial of the Sanguinaire
	// When you dodge, you have a 4% chance to gain 963 mastery for 20s. This effect can stack up to 3 times.
	shared.ItemVersionMap{
		shared.ItemVersionLFR:                 95779,
		shared.ItemVersionNormal:              94518,
		shared.ItemVersionHeroic:              96523,
		shared.ItemVersionThunderforged:       96151,
		shared.ItemVersionHeroicThunderforged: 96895,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		label := "Delicate Vial of the Sanguinaire"

		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			statValue := core.GetItemEffectScaling(itemID, 2.97000002861, state)

			aura, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				Duration:             time.Second * 20,
				MaxStacks:            3,
				BonusPerStack:        stats.Stats{stats.MasteryRating: statValue},
				StackingAuraActionID: core.ActionID{SpellID: 138864},
				StackingAuraLabel:    fmt.Sprintf("Blood of Power %s", versionLabel),
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       label,
				ProcChance: 0.04,
				Outcome:    core.OutcomeDodge,
				Callback:   core.CallbackOnSpellHitTaken,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					aura.Activate(sim)
					aura.AddStack(sim)
				},
			})

			character.AddStatProcBuff(itemID, aura, false, core.TrinketSlots())
		})
	})

	// Primordius' Talisman of Rage
	// Your attacks have a chance to grant you 963 Strength for 10s. This effect can stack up to 5 times. (Approximately
	// 3.50 procs per minute)
	shared.ItemVersionMap{
		shared.ItemVersionLFR:                 95757,
		shared.ItemVersionNormal:              94519,
		shared.ItemVersionHeroic:              96501,
		shared.ItemVersionThunderforged:       96129,
		shared.ItemVersionHeroicThunderforged: 96873,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		label := "Primordius' Talisman of Rage"

		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			statValue := core.GetItemEffectScaling(itemID, 0.5189999938, state)

			aura, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				Duration:             time.Second * 10,
				MaxStacks:            5,
				BonusPerStack:        stats.Stats{stats.Strength: statValue},
				StackingAuraActionID: core.ActionID{SpellID: 138870},
				StackingAuraLabel:    fmt.Sprintf("Rampage %s", versionLabel),
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:    label,
				Harmful: true,
				DPM: character.NewRPPMProcManager(itemID, false, core.ProcMaskDirect|core.ProcMaskProc, core.RPPMConfig{
					PPM: 3.5,
				}),
				ICD:      time.Second * 5,
				Outcome:  core.OutcomeLanded,
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					aura.Activate(sim)
					aura.AddStack(sim)
				},
			})

			character.AddStatProcBuff(itemID, aura, false, core.TrinketSlots())
		})
	})

	// Talisman of Bloodlust
	// Your attacks have a chance to grant you 963 haste for 10s. This effect can stack up to 5 times. (Approximately
	// 3.50 procs per minute)
	shared.ItemVersionMap{
		shared.ItemVersionLFR:                 95748,
		shared.ItemVersionNormal:              94522,
		shared.ItemVersionHeroic:              96492,
		shared.ItemVersionThunderforged:       96120,
		shared.ItemVersionHeroicThunderforged: 96864,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		label := "Talisman of Bloodlust"

		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			statValue := core.GetItemEffectScaling(itemID, 0.5189999938, state)

			aura, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				Duration:             time.Second * 10,
				MaxStacks:            5,
				BonusPerStack:        stats.Stats{stats.HasteRating: statValue},
				StackingAuraActionID: core.ActionID{SpellID: 138895},
				StackingAuraLabel:    fmt.Sprintf("Frenzy %s", versionLabel),
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:    label,
				Harmful: true,
				DPM: character.NewRPPMProcManager(itemID, false, core.ProcMaskDirect|core.ProcMaskProc, core.RPPMConfig{
					PPM: 3.5,
				}),
				ICD:      time.Second * 5,
				Outcome:  core.OutcomeLanded,
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					aura.Activate(sim)
					aura.AddStack(sim)
				},
			})

			character.AddStatProcBuff(itemID, aura, false, core.TrinketSlots())
		})
	})

	// Gaze of the Twins
	// Your critical attacks have a chance to grant you 963 Critical Strike for 20s. This effect can stack up
	// to 3 times. (Approximately 0.72 procs per minute)
	shared.ItemVersionMap{
		shared.ItemVersionLFR:                 95799,
		shared.ItemVersionNormal:              94529,
		shared.ItemVersionHeroic:              96543,
		shared.ItemVersionThunderforged:       96171,
		shared.ItemVersionHeroicThunderforged: 96915,
	}.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
		label := "Gaze of the Twins"

		core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
			character := agent.GetCharacter()
			statValue := core.GetItemEffectScaling(itemID, 0.96799999475, state)

			aura, _ := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
				Duration:             time.Second * 20,
				MaxStacks:            3,
				BonusPerStack:        stats.Stats{stats.CritRating: statValue},
				StackingAuraActionID: core.ActionID{SpellID: 139170},
				StackingAuraLabel:    fmt.Sprintf("Eye of Brutality %s", versionLabel),
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:    label,
				Harmful: true,
				DPM: character.NewRPPMProcManager(itemID, false, core.ProcMaskDirect|core.ProcMaskProc, core.RPPMConfig{
					PPM: 0.72000002861,
				}.WithCritMod()),
				ICD:      time.Second * 10,
				Outcome:  core.OutcomeCrit,
				Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					aura.Activate(sim)
					aura.AddStack(sim)
				},
			})

			character.AddStatProcBuff(itemID, aura, false, core.TrinketSlots())
		})
	})

}
