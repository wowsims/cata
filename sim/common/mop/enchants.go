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

		for _, aura := range auras {
			character.AddStatProcBuff(4441, aura, true, core.AllWeaponSlots())
		}
	})

	// Permanently enchants a melee weapon to sometimes increase your Intellect by 0 when healing or dealing
	// damage with spells. If less than 25% of your mana remains when the effect is triggered, your Spirit will
	// also increase by 0.
	newJadeSpiritEnchant := func(name string, effectId int32, procEffectId int32, buffEffectId int32, icd time.Duration) {
		core.NewEnchantEffect(effectId, func(agent core.Agent, _ proto.ItemLevelState) {
			character := agent.GetCharacter()
			duration := time.Second * 12

			intellect := character.NewTemporaryStatsAura(
				name+" - Intellect",
				core.ActionID{SpellID: buffEffectId}.WithTag(1),
				stats.Stats{stats.Intellect: 1650},
				duration,
			)
			spirit := character.NewTemporaryStatsAura(
				name+" - Spirit",
				core.ActionID{SpellID: buffEffectId}.WithTag(2),
				stats.Stats{stats.Spirit: 750},
				duration,
			)

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     "Enchant Weapon - " + name,
				Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
				Harmful:  true,
				ActionID: core.ActionID{SpellID: procEffectId},
				ICD:      icd,
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

			character.AddStatProcBuff(4442, intellect, true, core.AllWeaponSlots())
		})
	}

	newJadeSpiritEnchant("Jade Spirit", 4442, 120033, 104993, 3*time.Second)
	// TODO: Currently the PVP variant has no ICD, TBD if this is intended.
	newJadeSpiritEnchant("Spirit of Conquest", 5124, 142536, 142535, 0)

	// Permanently enchants a melee weapon to sometimes increase your Strength or Agility by 0 when dealing melee
	// damage. Your highest stat is always chosen.
	newDancingSteelEnchant := func(name string, effectId int32, procEffectId int32, agiEffectId int32, strEffectId int32) {
		core.NewEnchantEffect(effectId, func(agent core.Agent, _ proto.ItemLevelState) {
			character := agent.GetCharacter()
			duration := time.Second * 12

			createDancingSteelAuras := func(tag int32) map[stats.Stat]*core.StatBuffAura {
				labelSuffix := core.Ternary(tag == 1, " Main Hand", " (Off Hand)")
				slot := core.Ternary(tag == 1, proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand)
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
				for _, aura := range auras {
					character.AddStatProcBuff(effectId, aura, true, []proto.ItemSlot{slot})
					character.AddStatProcBuff(effectId, aura, true, []proto.ItemSlot{slot})
				}
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

	// Permanently enchants a melee weapon to make your damaging melee strikes sometimes activate a Mogu protection
	// spell, absorbing up to 8000 damage.
	core.NewEnchantEffect(4445, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		shield := character.NewDamageAbsorptionAura(core.AbsorptionAuraConfig{
			Aura: core.Aura{
				Label:    "Colossus" + character.Label,
				ActionID: core.ActionID{SpellID: 116631},
				Duration: time.Second * 10,
			},
			ShieldStrengthCalculator: func(_ *core.Unit) float64 {
				return 8000
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Enchant Weapon - Colossus",
			Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			Harmful:  true,
			ActionID: core.ActionID{SpellID: 118314},
			DPM: character.NewRPPMProcManager(
				4445,
				true,
				core.ProcMaskDirect|core.ProcMaskProc,
				core.RPPMConfig{
					PPM: 5.5,
				}.WithHasteMod(),
			),
			Outcome: core.OutcomeLanded,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				shield.Activate(sim)
			},
		})
	})

	// Permanently enchants a melee weapon to sometimes increase your dodge by 1650 for 7s when dealing melee
	// damage.
	core.NewEnchantEffect(4446, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		duration := time.Second * 7

		aura := character.NewTemporaryStatsAura(
			"River's Song",
			core.ActionID{SpellID: 116660}.WithTag(1),
			stats.Stats{stats.DodgeRating: 1650},
			duration,
		)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Enchant Weapon - River's Song",
			Callback: core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			Harmful:  true,
			ActionID: core.ActionID{SpellID: 104441},
			ICD:      time.Millisecond * 100,
			DPM: character.NewRPPMProcManager(
				4446,
				true,
				core.ProcMaskDirect|core.ProcMaskProc,
				core.RPPMConfig{
					PPM:         3.67,
					Coefficient: 1.0,
				}.WithHasteMod(),
			),
			Outcome: core.OutcomeLanded,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				aura.Activate(sim)
			},
		})
	})

	// Synapse Springs
	core.NewEnchantEffect(4898, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		if !character.HasProfession(proto.Profession_Engineering) {
			return
		}

		bonus := stats.Stats{}
		bonus[character.GetHighestStatType([]stats.Stat{
			stats.Strength, stats.Agility, stats.Intellect,
		})] = 1920

		core.RegisterTemporaryStatsOnUseCD(character,
			"Synapse Springs",
			bonus,
			10*time.Second,
			core.SpellConfig{
				ActionID: core.ActionID{SpellID: 126734},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute,
					},
					SharedCD: core.Cooldown{
						Timer:    character.GetOffensiveTrinketCD(),
						Duration: 10 * time.Second,
					},
				},
			})
	})

	// Phase Fingers
	core.NewEnchantEffect(4697, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()
		if !character.HasProfession(proto.Profession_Engineering) {
			return
		}

		core.RegisterTemporaryStatsOnUseCD(character,
			"Phase Fingers",
			stats.Stats{stats.DodgeRating: 2880},
			10*time.Second,
			core.SpellConfig{
				ActionID: core.ActionID{SpellID: 108788},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute,
					},
					SharedCD: core.Cooldown{
						Timer:    character.GetDefensiveTrinketCD(),
						Duration: 10 * time.Second,
					},
				},
			})
	})
}
