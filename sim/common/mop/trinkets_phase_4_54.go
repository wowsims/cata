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
	newReadinessTrinket := func(config *readinessTrinketConfig) {
		config.itemVersionMap.RegisterAll(func(version shared.ItemVersion, itemID int32, versionLabel string) {
			core.NewItemEffect(itemID, func(agent core.Agent, state proto.ItemLevelState) {
				character := agent.GetCharacter()

				auraID, exists := config.cdrAuraIDs[character.Spec]
				var cdrAura *core.Aura
				if exists {
					cdr := core.GetItemEffectScaling(itemID, 0.00989999995, state) / 100
					cdrAura = core.MakePermanent(character.RegisterAura(core.Aura{
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

				aura.Icd = triggerAura.Icd
				eligibleSlots := character.ItemSwap.EligibleSlotsForItem(itemID)
				character.AddStatProcBuff(itemID, aura, false, eligibleSlots)
				character.ItemSwap.RegisterProcWithSlots(itemID, triggerAura, eligibleSlots)
				if cdrAura != nil {
					character.ItemSwap.RegisterProcWithSlots(itemID, cdrAura, eligibleSlots)
				}
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
