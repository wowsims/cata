package mop

import (
	"time"

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
	})

	for _, version := range []core.ItemVersion{core.ItemVersionLFR, core.ItemVersionNormal, core.ItemVersionHeroic, core.ItemVersionThunderforged, core.ItemVersionHeroicThunderforged} {
		labelSuffix := []string{" (Celestial)", "", " (Heroic)", " (Thunderforged)", " (Heroic Thunderforged)"}[version]

		// Renataki's Soul Charm
		renatakiItemID := []int32{95625, 94512, 96369, 95997, 96741}[version]
		if renatakiItemID != 0 {
			core.NewItemEffect(renatakiItemID, func(agent core.Agent, state proto.ItemLevelState) {
				character := agent.GetCharacter()
				label := "Blades of Renataki"

				statValue := core.GetItemEffectScaling(renatakiItemID, 0.44999998808, state)

				_, aura := character.NewTemporaryStatBuffWithStacks(core.TemporaryStatBuffWithStacksConfig{
					AuraLabel:            label + labelSuffix,
					ActionID:             core.ActionID{SpellID: 138756},
					Duration:             time.Second * 20,
					MaxStacks:            10,
					TimePerStack:         time.Second * 1,
					BonusPerStack:        stats.Stats{stats.Agility: statValue},
					StackingAuraActionID: core.ActionID{SpellID: 138737},
					StackingAuraLabel:    "Item - Proc Stacking Agility" + labelSuffix,
				})

				core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
					Name:    label,
					Harmful: true,
					ICD:     time.Second * 10,
					DPM: character.NewRPPMProcManager(renatakiItemID, false, core.ProcMaskDirect|core.ProcMaskProc, core.RPPMConfig{
						PPM: 1.21000003815,
					}),
					Outcome:  core.OutcomeLanded,
					Callback: core.CallbackOnSpellHitDealt,
					Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
						aura.Activate(sim)
					},
				})
			})
		}

	}

}
