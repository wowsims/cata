package cata

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
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
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
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

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4066, aura)
	})

	// Enchant: 4067, Spell: 74197 - Enchant Weapon - Avalanche
	// http://elitistjerks.com/f79/t110302-enhsim_cataclysm/p4/#post1832162
	// Research indicates that the proc itself does not behave as game tables suggest <.<
	core.NewEnchantEffect(4067, func(agent core.Agent) {
		character := agent.GetCharacter()
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 74196},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskProc,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(463, 537), spell.OutcomeMagicHitAndCrit)
			},
		})

		// TODO: Verify PPM (currently based on elitist + simcraft)
		procMask := character.GetProcMaskForEnchant(4067)
		ppmm := character.AutoAttacks.NewPPMManager(5.0, procMask)
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
				if meleeIcd.IsReady(sim) && ppmm.Proc(sim, spell.ProcMask, "Avalanche") {
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

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4067, 5.0, &ppmm, aura)
	})

	// Enchant: 4074, Spell: 74211 - Enchant Weapon - Elemental Slayer
	core.NewEnchantEffect(4074, func(agent core.Agent) {
		character := agent.GetCharacter()
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 74208},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskProc,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(929, 1271), spell.OutcomeMagicHitAndCrit)
			},
		})

		// simpler spell than avalanche, only melee based proc chance
		// no verified PPM but silence effect + increased damage
		// will result in significant lower PPM
		// TODO: Verify PPM
		procMask := character.GetProcMaskForEnchant(4074)
		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Elemental Slayer",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      2.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4074, 2.0, aura.Ppmm, aura)
	})

	// Enchant: 4083, Spell: 74223 - Enchant Weapon - Hurricane
	core.NewEnchantEffect(4083, func(agent core.Agent) {
		character := agent.GetCharacter()

		procBuilder := func(name string, tag int32) *core.Aura {
			return character.NewTemporaryStatsAura(
				name,
				core.ActionID{SpellID: 74221, Tag: tag},
				stats.Stats{stats.MeleeHaste: 450, stats.SpellHaste: 450},
				time.Second*12,
			)
		}

		mhAura := procBuilder("Hurricane Enchant MH", 1)
		ohAura := procBuilder("Hurricane Enchant OH", 2)
		spAura := procBuilder("Hurricane Enchant Spell", 3)

		procMask := character.GetProcMaskForEnchant(4083)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

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

				if ppmm.Proc(sim, spell.ProcMask, "Hurricane") {
					if spell.IsMH() {
						mhAura.Activate(sim)
					} else {
						ohAura.Activate(sim)
					}

					return
				}

				if aura.Icd.IsReady(sim) && spell.ProcMask.Matches(core.ProcMaskSpellDamage) && sim.Proc(0.15, "Hurricane") {
					aura.Icd.Use(sim)
					hurricaneSpellProc(sim)
				}
			},

			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if aura.Icd.IsReady(sim) && sim.Proc(0.15, "Hurricane") {
					aura.Icd.Use(sim)
					hurricaneSpellProc(sim)
				}
			},
		}))

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4074, 1.0, &ppmm, aura)
	})

	// Enchant: 4084, Spell: 74225 - Enchant Weapon - Heartsong
	core.NewEnchantEffect(4084, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Heartsong Proc",
			core.ActionID{SpellID: 74224},
			stats.Stats{stats.Spirit: 200},
			time.Second*15,
		)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Heartsong",
			ActionID:   core.ActionID{SpellID: 74224},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 20,
			ProcChance: 0.25,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4084, aura)
	})

	// Enchant: 4097, Spell: 74242 - Enchant Weapon - Power Torrent
	core.NewEnchantEffect(4097, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Power Torrent Proc",
			core.ActionID{SpellID: 74241},
			stats.Stats{stats.Intellect: 500},
			time.Second*12,
		)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Power Torrent",
			ActionID:   core.ActionID{SpellID: 74224},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 45,
			ProcChance: 1.0 / 3.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4097, aura)
	})

	// Enchant: 4098, Spell: 74244 - Enchant Weapon - Windwalk
	core.NewEnchantEffect(4098, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Windwalk Proc",
			core.ActionID{SpellID: 74243},
			stats.Stats{stats.Dodge: 600},
			time.Second*10,
		)

		procMask := character.GetProcMaskForEnchant(4098)
		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Windwalk",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      2.5,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4098, 2.5, aura.Ppmm, aura)
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

		procMask := character.GetProcMaskForEnchant(4099)
		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Landslide",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.IsMH() {
					mainHand.Activate(sim)
				} else {
					offHand.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4099, 1, aura.Ppmm, aura)
	})

	// Enchant: 4115, Spell: 75172 - Lightweave Embroidery
	core.NewEnchantEffect(4115, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Lightweave Embroidery Proc",
			core.ActionID{SpellID: 75170},
			stats.Stats{stats.Intellect: 480},
			time.Second*15,
		)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Lightweave Embroidery Cata",
			ActionID:   core.ActionID{SpellID: 75171},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 60,
			ProcChance: 0.35,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4115, aura)
	})

	// Enchant: 4116, Spell: 75175 - Darkglow Embroidery
	core.NewEnchantEffect(4116, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Darkglow Embroidery Cata",
			core.ActionID{SpellID: 75175},
			stats.Stats{stats.Spirit: 480},
			time.Second*15,
		)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Darkglow Embroidery Cata",
			ActionID:   core.ActionID{SpellID: 75173},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 60,
			ProcChance: 0.30,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4116, aura)
	})

	// Enchant: 4118, Spell: 75178 - Swordguard Embroidery
	core.NewEnchantEffect(4118, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Swordguard Embroidery Cata",
			core.ActionID{SpellID: 75176},
			stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
			time.Second*15,
		)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Swordguard Embroidery Cataa",
			ActionID:   core.ActionID{SpellID: 75176},
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 55,
			ProcChance: 0.15,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
			},
		})
		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4118, aura)

	})

	// Enchant: 4175, Spell: 81932, Item: 59594 - Gnomish X-Ray Scope
	core.NewEnchantEffect(4175, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"X-Ray Targeting",
			core.ActionID{SpellID: 95712},
			stats.Stats{stats.RangedAttackPower: 800},
			time.Second*10,
		)

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Gnomish X-Ray Scope",
			ActionID:   core.ActionID{SpellID: 95712},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskRanged,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 40,
			ProcChance: 0.1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4175, aura)
	})

	// Enchant: 4215, Spell: 92433, Item: 55055 - Elementium Shield Spike
	core.NewEnchantEffect(4215, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 92432}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(90, 133)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Elementium Shield Spike",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeBlock,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, spell.Unit)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4215, aura)
	})

	// Enchant: 4216, Spell: 92437, Item: 55056  - Pyrium Shield Spike
	core.NewEnchantEffect(4216, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 92436}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(210, 350)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Pyrium Shield Spike",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeBlock,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, spell.Unit)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4216, aura)
	})

	// Enchant: 4267, Spell: 99623, Item: 70139 - Flintlocke's Woodchucker
	core.NewEnchantEffect(4267, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Flintlocke's Woodchucker Proc",
			core.ActionID{SpellID: 99621},
			stats.Stats{stats.Agility: 300},
			time.Second*10,
		)

		dmgProc := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 99621},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(550, 1650)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Flintlocke's Woodchucker",
			ActionID:   core.ActionID{SpellID: 99621},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskRanged,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 40,
			ProcChance: 0.1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				statAura.Activate(sim)
				dmgProc.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4267, aura)
	})
}
