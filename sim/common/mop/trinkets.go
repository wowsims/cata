package mop

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type readinessTrinketConfig struct {
	itemVersionMap   shared.ItemVersionMap
	baseTrinketLabel string
	buffAuraLabel    string
	buffAuraID       int32
	buffedStat       stats.Stat
	buffDuration     time.Duration
	icd              time.Duration
	cdrAuraIDs       map[proto.Spec]int32
}

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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

		eligibleSlots := character.ItemSwap.EligibleSlotsForItem(75274)
		for _, aura := range auras {
			character.AddStatProcBuff(75274, aura, false, eligibleSlots)
		}
		character.ItemSwap.RegisterProcWithSlots(81266, triggerAura, eligibleSlots)
	})

	core.NewItemEffect(81266, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 126467}
		manaMetrics := character.NewManaMetrics(actionID)

		mana := core.GetItemEffectScaling(81266, 2.97199988365, state)

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

		eligibleSlots := character.ItemSwap.EligibleSlotsForItem(81266)
		character.ItemSwap.RegisterProcWithSlots(81266, triggerAura, eligibleSlots)
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

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			eligibleSlots := character.ItemSwap.EligibleSlotsForItem(itemID)
			character.AddStatProcBuff(itemID, statBuffAura, false, eligibleSlots)
			character.ItemSwap.RegisterProcWithSlots(itemID, triggerAura, eligibleSlots)
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

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       label,
				ProcChance: 0.04,
				Outcome:    core.OutcomeDodge,
				Callback:   core.CallbackOnSpellHitTaken,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					aura.Activate(sim)
					aura.AddStack(sim)
				},
			})

			eligibleSlots := character.ItemSwap.EligibleSlotsForItem(itemID)
			character.AddStatProcBuff(itemID, aura, false, eligibleSlots)
			character.ItemSwap.RegisterProcWithSlots(itemID, triggerAura, eligibleSlots)
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

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			eligibleSlots := character.ItemSwap.EligibleSlotsForItem(itemID)
			character.AddStatProcBuff(itemID, aura, false, eligibleSlots)
			character.ItemSwap.RegisterProcWithSlots(itemID, triggerAura, eligibleSlots)
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

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			eligibleSlots := character.ItemSwap.EligibleSlotsForItem(itemID)
			character.AddStatProcBuff(itemID, aura, false, eligibleSlots)
			character.ItemSwap.RegisterProcWithSlots(itemID, triggerAura, eligibleSlots)
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

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			eligibleSlots := character.ItemSwap.EligibleSlotsForItem(itemID)
			character.AddStatProcBuff(itemID, aura, false, eligibleSlots)
			character.ItemSwap.RegisterProcWithSlots(itemID, triggerAura, eligibleSlots)
		})
	})

	newReadinessTrinket := func(config *readinessTrinketConfig) {
		config.itemVersionMap.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
			core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
				character := agent.GetCharacter()

				auraID, exists := config.cdrAuraIDs[character.Spec]
				if exists {
					cdr := core.GetItemEffectScaling(itemID, 0.00989999995, state)
					core.MakePermanent(character.RegisterAura(core.Aura{
						Label:    fmt.Sprintf("Readiness %s", versionLabel),
						ActionID: core.ActionID{SpellID: auraID},
					}).AttachSpellMod(core.SpellModConfig{
						Kind:       core.SpellMod_Cooldown_Multiplier,
						SpellFlag:  core.SpellFlagReadinessTrinket,
						FloatValue: cdr,
					}))
				}

				stats := stats.Stats{}
				stats[config.buffedStat] = core.GetItemEffectScaling(itemID, 0.96799999475, state)

				aura := character.NewTemporaryStatsAura(
					fmt.Sprintf("%s %s", config.buffAuraLabel, versionLabel),
					core.ActionID{SpellID: config.buffAuraID},
					stats,
					config.buffDuration,
				)

				triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
					Name:       config.baseTrinketLabel,
					Harmful:    true,
					ProcChance: 0.15,
					ICD:        config.icd,
					ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
					Outcome:    core.OutcomeLanded,
					Callback:   core.CallbackOnSpellHitDealt,
					Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
						aura.Activate(sim)
					},
				})

				eligibleSlots := character.ItemSwap.EligibleSlotsForItem(itemID)
				character.AddStatProcBuff(itemID, aura, false, eligibleSlots)
				character.ItemSwap.RegisterProcWithSlots(itemID, triggerAura, eligibleSlots)

			})
		})
	}

	// Assurance of Consequence
	// Increases the cooldown recovery rate of six of your major abilities by 47%.
	// Effective for Agility-based damage roles only.
	//
	// Your attacks have a chance to grant you 14039 Agility for 20 sec.
	// (15% chance, 115 sec cooldown) (Proc chance: 15%, 1.917m cooldown)
	newReadinessTrinket(&readinessTrinketConfig{
		itemVersionMap: shared.ItemVersionMap{
			shared.ItemVersionLFR:             104974,
			shared.ItemVersionNormal:          102292,
			shared.ItemVersionHeroic:          104476,
			shared.ItemVersionWarforged:       105223,
			shared.ItemVersionHeroicWarforged: 105472,
			shared.ItemVersionFlexible:        104725,
		},
		baseTrinketLabel: "Assurance of Consequence",
		buffAuraLabel:    "Dextrous",
		buffAuraID:       146308,
		buffedStat:       stats.Agility,
		buffDuration:     time.Second * 20,
		icd:              time.Second * 115,
		cdrAuraIDs: map[proto.Spec]int32{
			// Druid
			// Missing: Bear Hug, Ironbark, Nature's Swiftness
			proto.Spec_SpecFeralDruid:       145961,
			proto.Spec_SpecGuardianDruid:    145962,
			proto.Spec_SpecRestorationDruid: 145963,
			// Hunter
			// Missing: Bestial Wrath
			proto.Spec_SpecBeastMasteryHunter: 145964,
			proto.Spec_SpecMarksmanshipHunter: 145965,
			proto.Spec_SpecSurvivalHunter:     145966,
			// Rogue
			// Missing: Cloak of Shadows, Evasion, JuJu Escape
			proto.Spec_SpecAssassinationRogue: 145983,
			proto.Spec_SpecCombatRogue:        145984,
			proto.Spec_SpecSubtletyRogue:      145985,
			// Priest - NOTE: Priests seem to have a Aura for this
			// Missing: Divine Hymn, Guardian Spirit, Hymn of Hope, Inner Focus, Pain Suppression, Power Word: Barrier, Void Shift
			// proto.Spec_SpecDisciplinePriest: 145981,
			// proto.Spec_SpecHolyPriest:       145982,
			// Shaman
			// Missing: Mana Tide Totem, Spirit Link Totem
			proto.Spec_SpecEnhancementShaman: 145986,
			proto.Spec_SpecRestorationShaman: 145988,
			// Monk
			// Missing: Zen Meditation, Life Cocoon, Revival, Thunder Focus Tea, Flying Serpent Kick
			proto.Spec_SpecBrewmasterMonk: 145967,
			proto.Spec_SpecMistweaverMonk: 145968,
			proto.Spec_SpecWindwalkerMonk: 145969,
		},
	})

	// Evil Eye of Galakras
	// Increases the cooldown recovery rate of six of your major abilities by 1%. Effective for Strength-based
	// damage roles only.
	//
	// Your attacks have a chance to grant you 11761 Strength for 10 sec.
	// (15% chance, 55 sec cooldown) (Proc chance: 15%, 55s cooldown)
	newReadinessTrinket(&readinessTrinketConfig{
		itemVersionMap: shared.ItemVersionMap{
			shared.ItemVersionLFR:             104993,
			shared.ItemVersionNormal:          102298,
			shared.ItemVersionHeroic:          104495,
			shared.ItemVersionWarforged:       105242,
			shared.ItemVersionHeroicWarforged: 105491,
			shared.ItemVersionFlexible:        104744,
		},
		baseTrinketLabel: "Evil Eye of Galakras",
		buffAuraLabel:    "Outrage",
		buffAuraID:       146245,
		buffedStat:       stats.Strength,
		buffDuration:     time.Second * 10,
		icd:              time.Second * 55,
		cdrAuraIDs: map[proto.Spec]int32{
			// Death Knight
			proto.Spec_SpecBloodDeathKnight:  145958,
			proto.Spec_SpecFrostDeathKnight:  145959,
			proto.Spec_SpecUnholyDeathKnight: 145960,
			// Paladin
			// Missing: Divine Plea, Hand Of Protection, Divine Shield, Hand Of Purity
			proto.Spec_SpecHolyPaladin:        145978,
			proto.Spec_SpecProtectionPaladin:  145976,
			proto.Spec_SpecRetributionPaladin: 145975,
			// Warrior
			// Missing: Die by the Sword, Mocking Banner
			proto.Spec_SpecArmsWarrior:       145990,
			proto.Spec_SpecFuryWarrior:       145991,
			proto.Spec_SpecProtectionWarrior: 145992,
		},
	})
}
