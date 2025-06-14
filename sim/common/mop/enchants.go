package mop

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {

	// Permanently enchants a melee weapon to sometimes increase your critical strike, haste, or mastery by 1500
	// for 12s when dealing damage or healing with spells and melee attacks.
	core.NewEnchantEffect(4441, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		duration := time.Second * 12

		haste := character.NewTemporaryStatsAura(
			"Windsong - Haste",
			core.ActionID{SpellID: 104423},
			stats.Stats{stats.HasteRating: 1500},
			duration,
		)
		crit := character.NewTemporaryStatsAura(
			"Windsong - Crit",
			core.ActionID{SpellID: 104509},
			stats.Stats{stats.CritRating: 1500},
			duration,
		)
		mastery := character.NewTemporaryStatsAura(
			"Windsong - Mastery",
			core.ActionID{SpellID: 104510},
			stats.Stats{stats.MasteryRating: 1500},
			duration,
		)

		auras := []*core.StatBuffAura{haste, crit, mastery}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Enchant Weapon - Windsong",
			Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			Harmful:  true,
			ActionID: core.ActionID{SpellID: 104561},
			DPM: character.NewRPPMProcManager(
				4441,
				true,
				core.ProcMaskDirect|core.ProcMaskProc,
				core.RPPMConfig{
					PPM: 2.2,
				},
			),
			Outcome: core.OutcomeLanded,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				aura := auras[int32(sim.RollWithLabel(0, 3, "Windsong Proc"))]
				aura.Activate(sim)
			},
		})
	})

	// Permanently enchants a melee weapon to sometimes increase your Intellect by 1650 when healing or dealing
	// damage with spells. If less than 25% of your mana remains when the effect is triggered, your Spirit will
	// also increase by 750.
	core.NewEnchantEffect(4442, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		duration := time.Second * 12

		intellect := character.NewTemporaryStatsAura(
			"Jade Spirit - Intellect",
			core.ActionID{SpellID: 104993}.WithTag(1),
			stats.Stats{stats.Intellect: 1650},
			duration,
		)
		spirit := character.NewTemporaryStatsAura(
			"Jade Spirit - Spirit",
			core.ActionID{SpellID: 104993}.WithTag(2),
			stats.Stats{stats.Spirit: 750},
			duration,
		)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Enchant Weapon - Jade Spirit",
			Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
			Harmful:  true,
			ActionID: core.ActionID{SpellID: 120033},
			ICD:      3 * time.Second,
			DPM: character.NewRPPMProcManager(
				4442,
				true,
				core.ProcMaskDirect|core.ProcMaskProc,
				core.RPPMConfig{
					PPM: 2.2,
				},
			),
			Outcome: core.OutcomeLanded,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				intellect.Activate(sim)
				if character.HasManaBar() && character.CurrentManaPercent() < 0.25 {
					spirit.Activate(sim)
				}
			},
		})
	})

	// Permanently enchants a melee weapon to sometimes increase your Strength or Agility by 0 when dealing melee
	// damage. Your highest stat is always chosen.
	newDancingSteelEnchant := func(name string, effectId int32, procEffectId int32, agiEffectId int32, strEffectId int32) {
		core.NewEnchantEffect(effectId, func(agent core.Agent, _ proto.ItemLevelState) {
			character := agent.GetCharacter()
			duration := time.Second * 12

			createDancingSteelAuras := func(tag int32) map[stats.Stat]*core.StatBuffAura {
				labelSuffix := core.Ternary(tag == 1, " Main Hand", " (Off Hand)")
				auras := make(map[stats.Stat]*core.StatBuffAura, 2)
				auras[stats.Agility] = character.NewTemporaryStatsAura(
					name+" - Agility"+labelSuffix,
					core.ActionID{SpellID: agiEffectId}.WithTag(tag),
					stats.Stats{stats.Agility: 1650},
					duration,
				)
				auras[stats.Strength] = character.NewTemporaryStatsAura(
					name+" - Strength"+labelSuffix,
					core.ActionID{SpellID: strEffectId}.WithTag(tag),
					stats.Stats{stats.Strength: 1650},
					duration,
				)
				return auras
			}

			mhAuras := createDancingSteelAuras(1)
			ohAuras := createDancingSteelAuras(2)

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     "Enchant Weapon - " + name,
				Callback: core.CallbackOnSpellHitDealt,
				Harmful:  true,
				ActionID: core.ActionID{SpellID: procEffectId},
				DPM: character.NewRPPMProcManager(
					effectId,
					true,
					core.ProcMaskMelee|core.ProcMaskMeleeProc,
					core.RPPMConfig{
						PPM: 2.53,
					},
				),
				Outcome: core.OutcomeLanded,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					core.Ternary(spell.IsOH(), ohAuras, mhAuras)[character.GetHighestStatType([]stats.Stat{stats.Strength, stats.Agility})].Activate(sim)
				},
			})
		})
	}

	newDancingSteelEnchant("Dancing Steel", 4444, 118333, 118334, 118335)
	newDancingSteelEnchant("Bloddy Dancing Steel", 5125, 142531, 142530, 142530)

}
