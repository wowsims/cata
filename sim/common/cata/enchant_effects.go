package cata

import (
	"time"

	"github.com/wowsims/mop/sim/common/shared"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {

	// Enchant: 4066, Spell: 74195 - Enchant Weapon - Mending
	core.NewEnchantEffect(4066, func(agent core.Agent) {
		character := agent.GetCharacter()
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 74194})

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 74194},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagIgnoreModifiers | core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcHealing(sim, target, sim.Roll(744, 856), spell.OutcomeHealingCrit)
				target.GainHealth(sim, result.Damage, healthMetrics)
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Mending",
			ActionID:   core.ActionID{SpellID: 74194},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskMelee | core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 15,
			ProcChance: 0.15,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, spell.Unit)
			},
		})

		character.ItemSwap.RegisterEnchantProc(4066, aura)
	})

	// Enchant: 4067, Spell: 74197 - Enchant Weapon - Avalanche
	// http://elitistjerks.com/f79/t110302-enhsim_mopclysm/p4/#post1832162
	// Research indicates that the proc itself does not behave as game tables suggest <.<
	core.NewEnchantEffect(4067, func(agent core.Agent) {
		character := agent.GetCharacter()
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 74196},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellProc,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(463, 537), spell.OutcomeMagicHitAndCrit)
			},
		})

		ppm := 5.0
		dpm := character.NewDynamicLegacyProcForEnchant(4067, ppm, 0)
		meleeIcd := &core.Cooldown{
			Duration: time.Millisecond * 1,
			Timer:    character.NewTimer(),
		}

		aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Avalanche",
			Icd: &core.Cooldown{
				Duration: time.Second * 10,
				Timer:    character.NewTimer(),
			},

			// https://www.wowhead.com/spell=74197/avalanche#comments:id=1374594
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}

				// melee ICD - Melee and Spell Procs are independent
				// while melee attacks do not have a direct internal CD
				// only one proc per server batch can occur, so i.E. Storm Strike can not proc avalance twice
				if meleeIcd.IsReady(sim) && dpm.Proc(sim, spell.ProcMask, "Avalanche") {
					meleeIcd.Use(sim)
					procSpell.Cast(sim, result.Target)
					return
				}

				if aura.Icd.IsReady(sim) && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && sim.Proc(0.25, "Avalanche") {
					aura.Icd.Use(sim)
					procSpell.Cast(sim, result.Target)
				}
			},

			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}

				if aura.Icd.IsReady(sim) && sim.Proc(0.25, "Avalanche") {
					aura.Icd.Use(sim)
					procSpell.Cast(sim, result.Target)
				}
			},
		}))

		character.ItemSwap.RegisterEnchantProc(4067, aura)
	})

	// Enchant: 4074, Spell: 74211 - Enchant Weapon - Elemental Slayer
	core.NewEnchantEffect(4074, func(agent core.Agent) {
		character := agent.GetCharacter()
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 74208},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellProc,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(929, 1271), spell.OutcomeMagicHitAndCrit)
			},
		})

		// simpler spell than avalanche, only melee based proc chance
		// no verified PPM but silence effect + increased damage
		// will result in significant lower PPM
		// TODO: Verify PPM
		ppm := 2.0
		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Elemental Slayer",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			DPM:      character.NewDynamicLegacyProcForEnchant(4074, ppm, 0),
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterEnchantProc(4074, aura)
	})

	// Enchant: 4083, Spell: 74223 - Enchant Weapon - Hurricane
	core.NewEnchantEffect(4083, func(agent core.Agent) {
		character := agent.GetCharacter()

		procBuilder := func(name string, tag int32) *core.StatBuffAura {
			return character.NewTemporaryStatsAura(
				name,
				core.ActionID{SpellID: 74221, Tag: tag},
				stats.Stats{stats.HasteRating: 450},
				time.Second*12,
			)
		}

		mhAura := procBuilder("Hurricane Enchant MH", 1)
		ohAura := procBuilder("Hurricane Enchant OH", 2)
		spAura := procBuilder("Hurricane Enchant Spell", 3)

		ppm := 1.0
		dpm := character.NewDynamicLegacyProcForEnchant(4083, ppm, 0)

		hurricaneSpellProc := func(sim *core.Simulation) {
			if mhAura.IsActive() {
				mhAura.Refresh(sim)
				return
			}

			if ohAura.IsActive() {
				ohAura.Refresh(sim)
				return
			}

			spAura.Activate(sim)
		}

		// spell similar to avalanche
		// seprate proc chances for melee / spells
		// each weapon has a separate proc
		// if spell proc occurs first, both MH and OH weapon get a separate proc up to 3
		// if either MH or OH proc is running spell proc will refresh one of those
		tryProcAndCastSpell := func(sim *core.Simulation, aura *core.Aura) {
			if aura.Icd.IsReady(sim) && sim.Proc(0.15, "Hurricane") {
				aura.Icd.Use(sim)
				hurricaneSpellProc(sim)
			}
		}
		aura := core.MakePermanent(character.GetOrRegisterAura(core.Aura{
			Label: "Hurricane",

			Icd: &core.Cooldown{
				Duration: time.Second * 45,
				Timer:    character.NewTimer(),
			},

			// https://www.wowhead.com/spell=74197/avalanche#comments:id=1374594
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}

				if dpm.Proc(sim, spell.ProcMask, "Hurricane") {
					if spell.IsMH() {
						mhAura.Activate(sim)
					} else {
						ohAura.Activate(sim)
					}

					return
				}

				if spell.ProcMask.Matches(core.ProcMaskSpellDamage | core.ProcMaskSpellHealing) {
					tryProcAndCastSpell(sim, aura)
				}
			},
			OnPeriodicHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				tryProcAndCastSpell(sim, aura)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				tryProcAndCastSpell(sim, aura)
			},
		}))

		character.ItemSwap.RegisterEnchantProc(4083, aura)

		if character.AutoAttacks.AutoSwingMelee {
			character.AddStatProcBuff(4083, mhAura, true, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand})
			character.AddStatProcBuff(4083, ohAura, true, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand})
		} else {
			character.AddStatProcBuff(4083, spAura, true, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand})
		}
	})

	// Enchant: 4084, Spell: 74225 - Enchant Weapon - Heartsong
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Heartsong",
		EnchantID:  4084,
		AuraID:     74224,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 20,
		ProcChance: 0.25,
		Bonus:      stats.Stats{stats.Spirit: 200},
		Duration:   time.Second * 15,
	})

	// Enchant: 4097, Spell: 74242 - Enchant Weapon - Power Torrent
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Power Torrent",
		EnchantID:  4097,
		AuraID:     74241,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 45,
		ProcChance: 1.0 / 3.0,
		Bonus:      stats.Stats{stats.Intellect: 500},
		Duration:   time.Second * 12,
	})

	// Enchant: 4098, Spell: 74244 - Enchant Weapon - Windwalk
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:      "Windwalk",
		EnchantID: 4098,
		AuraID:    74243,
		Callback:  core.CallbackOnSpellHitDealt,
		Outcome:   core.OutcomeLanded,
		PPM:       1, // based on old Wowhead comments, TODO: measure in Classic
		Bonus:     stats.Stats{stats.DodgeRating: 600},
		Duration:  time.Second * 10,
	})

	// Enchant: 4099, Spell: 74246 - Enchant Weapon - Landslide
	core.NewEnchantEffect(4099, func(agent core.Agent) {
		character := agent.GetCharacter()

		mainHand := character.NewTemporaryStatsAura(
			"Landslide Proc MH",
			core.ActionID{SpellID: 74245, Tag: 1},
			stats.Stats{stats.AttackPower: 1000},
			time.Second*12,
		)

		offHand := character.NewTemporaryStatsAura(
			"Landslide Proc OH",
			core.ActionID{SpellID: 74245, Tag: 2},
			stats.Stats{stats.AttackPower: 1000},
			time.Second*12,
		)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Landslide",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			DPM:      character.NewDynamicLegacyProcForEnchant(4099, 1.0, 0),
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.IsMH() {
					mainHand.Activate(sim)
				} else {
					offHand.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterEnchantProc(4099, aura)
		if character.AutoAttacks.AutoSwingMelee {
			character.AddStatProcBuff(74245, mainHand, true, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand})
			character.AddStatProcBuff(74245, offHand, true, []proto.ItemSlot{proto.ItemSlot_ItemSlotOffHand})
		}
	})

	// Enchant: 4115, Spell: 75172 - Lightweave Embroidery
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Lightweave Embroidery Cata",
		EnchantID:  4115,
		AuraID:     75170,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 64,
		ProcChance: 0.25,
		Bonus:      stats.Stats{stats.Intellect: 580},
		Duration:   time.Second * 15,
	})

	// Enchant: 4116, Spell: 75175 - Darkglow Embroidery
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Darkglow Embroidery Cata",
		EnchantID:  4116,
		AuraID:     75173,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 57,
		ProcChance: 0.30,
		Bonus:      stats.Stats{stats.Spirit: 580},
		Duration:   time.Second * 15,
	})

	// Enchant: 4118, Spell: 75178 - Swordguard Embroidery
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Swordguard Embroidery Cata",
		EnchantID:  4118,
		AuraID:     75178,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 55,
		ProcChance: 0.15,
		Bonus:      stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
		Duration:   time.Second * 15,
	})

	// Enchant: 4175, Spell: 81932, Item: 59594 - Gnomish X-Ray Scope
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:       "Gnomish X-Ray Scope",
		EnchantID:  4175,
		ItemID:     59594,
		AuraID:     95712,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskRanged,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 40,
		ProcChance: 0.1,
		Bonus:      stats.Stats{stats.RangedAttackPower: 800},
		Duration:   time.Second * 10,
	})

	// Enchant: 4176, Item: 59595 - R19 Threatfinder
	core.NewEnchantEffect(4176, func(agent core.Agent) {
		agent.GetCharacter().AddBonusRangedHitPercent(88 / core.PhysicalHitRatingPerHitPercent)
	})

	// Enchant: 4177, Item: 59596 - Safety Catch Removal Kit
	core.NewEnchantEffect(4177, func(agent core.Agent) {
		character := agent.GetCharacter()
		// TODO: This should be ranged-only haste. For now just make it hunter-only.
		if character.Class == proto.Class_ClassHunter {
			character.AddStat(stats.HasteRating, 88)
		}
	})

	// Enchant: 4215, Spell: 92433, Item: 55055 - Elementium Shield Spike
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		EnchantID: 4215,
		SpellID:   92432,
		Trigger: core.ProcTrigger{
			Name:     "Elementium Shield Spike",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeBlock,
		},
		School:  core.SpellSchoolPhysical,
		MinDmg:  90,
		MaxDmg:  133,
		Outcome: shared.OutcomeMeleeCanCrit,
		IsMelee: true,
	})

	// Enchant: 4216, Spell: 92437, Item: 55056  - Pyrium Shield Spike
	shared.NewProcDamageEffect(shared.ProcDamageEffect{
		EnchantID: 4216,
		SpellID:   92436,
		Trigger: core.ProcTrigger{
			Name:     "Pyrium Shield Spike",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeBlock,
		},
		School:  core.SpellSchoolPhysical,
		Outcome: shared.OutcomeMeleeCanCrit,
		MinDmg:  210,
		MaxDmg:  350,
		IsMelee: true,
	})

	// Enchant: 4267, Spell: 99623, Item: 70139 - Flintlocke's Woodchucker
	shared.NewProcStatBonusEffectWithDamageProc(shared.ProcStatBonusEffect{
		Name:       "Flintlocke's Woodchucker",
		EnchantID:  4267,
		ItemID:     70139,
		AuraID:     99621,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskRanged,
		Outcome:    core.OutcomeLanded,
		ICD:        time.Second * 40,
		ProcChance: 0.1,
		Bonus:      stats.Stats{stats.Agility: 300},
		Duration:   time.Second * 10,
	},
		shared.DamageEffect{
			SpellID:  99621,
			School:   core.SpellSchoolPhysical,
			ProcMask: core.ProcMaskEmpty,
			Outcome:  shared.OutcomeRangedCanCrit,
			MinDmg:   550,
			MaxDmg:   1650,
			IsMelee:  true,
		})

	movementSpeedEnchants := []int32{
		3232, // Enchant Boots - Tuskarr's Vitality
		4104, // Enchant Boots - Lavawalker
		4105, // Enchant Boots - Assassin's Step
		4062, // Enchant Boots - Earthen Vitality
	}

	for _, enchantID := range movementSpeedEnchants {
		core.NewEnchantEffect(enchantID, func(agent core.Agent) {
			character := agent.GetCharacter()
			aura := character.NewMovementSpeedAura("Minor Run Speed", core.ActionID{SpellID: 13889}, 0.08)

			character.ItemSwap.RegisterEnchantProc(enchantID, aura)
		})
	}
}
