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

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Mending",
			Duration: core.NeverExpires,
			Icd: &core.Cooldown{
				Duration: time.Second * 15,
				Timer:    character.NewTimer(),
			},

			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			// Same proc mask an description as Avalanche
			// https://www.wowhead.com/spell=74197/avalanche#comments:id=1374594
			// Might have independent Melee proc without ICD but PPM based.
			// TODO: Verify
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.ProcMask.Matches(core.ProcMaskRanged) || !aura.Icd.IsReady(sim) {
					return
				}

				if sim.Proc(0.15, "Mending") {
					aura.Icd.Use(sim)
					procSpell.Cast(sim, aura.Unit)
				}
			},

			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || spell.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}

				if aura.Icd.IsReady(sim) && sim.Proc(0.15, "Mending") {
					aura.Icd.Use(sim)
					procSpell.Cast(sim, aura.Unit)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(3241, aura)
	})

	// Enchant: 4067, Spell: 74197 - Enchant Weapon - Avalanche
	// http://elitistjerks.com/f79/t110302-enhsim_cataclysm/p4/#post1832162
	// Research indicates that the proc itself does not behave as game tables suggest <.<
	core.NewEnchantEffect(4067, func(agent core.Agent) {
		character := agent.GetCharacter()
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 74196},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellOrProc,

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

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Avalanche",
			Duration: core.NeverExpires,
			Icd: &core.Cooldown{
				Duration: time.Second * 10,
				Timer:    character.NewTimer(),
			},

			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
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
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4067, 5.0, &ppmm, aura)
	})

	// Enchant: 4074, Spell: 74211 - Enchant Weapon - Elemental Slayer
	core.NewEnchantEffect(4074, func(agent core.Agent) {
		character := agent.GetCharacter()
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 74208},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellOrProc,

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
		ppmm := character.AutoAttacks.NewPPMManager(2.0, procMask)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Elemental Slayer",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			// https://www.wowhead.com/spell=74197/avalanche#comments:id=1374594
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Elemental Slayer") {
					procSpell.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4074, 2.0, &ppmm, aura)
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
		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Hurricane",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

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
		})

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

		handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskSpellDamage|core.ProcMaskSpellHealing) {
				return
			}

			if aura.Icd.IsReady(sim) && sim.Proc(0.25, "Heartsong") {
				aura.Icd.Use(sim)
				statAura.Activate(sim)
			}
		}

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Heartsong",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			Icd: &core.Cooldown{
				Duration: time.Second * 20,
				Timer:    character.NewTimer(),
			},

			OnSpellHitDealt:       handler,
			OnPeriodicDamageDealt: handler,
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

		handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskSpellDamage|core.ProcMaskSpellHealing) {
				return
			}

			if aura.Icd.IsReady(sim) && sim.Proc(1.0/3.0, "Power Torrent") {
				aura.Icd.Use(sim)
				statAura.Activate(sim)
			}
		}

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Power Torrent",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			Icd: &core.Cooldown{
				Duration: time.Second * 45,
				Timer:    character.NewTimer(),
			},

			OnSpellHitDealt:       handler,
			OnPeriodicDamageDealt: handler,
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
		ppmm := character.AutoAttacks.NewPPMManager(2.5, procMask)
		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Windwalk",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Windwalk") {
					statAura.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4098, 2.5, &ppmm, aura)
	})

	// Enchant: 4099, Spell: 74246 - Enchant Weapon - Landslide
	core.NewEnchantEffect(4099, func(agent core.Agent) {
		character := agent.GetCharacter()

		mainHand := character.NewTemporaryStatsAura(
			"Landslide Proc",
			core.ActionID{SpellID: 74245},
			stats.Stats{stats.AttackPower: 1000},
			time.Second*12,
		)

		offHand := character.NewTemporaryStatsAura(
			"Landslide Proc",
			core.ActionID{SpellID: 74245},
			stats.Stats{stats.AttackPower: 1000},
			time.Second*12,
		)

		procMask := character.GetProcMaskForEnchant(4099)
		ppmm := character.AutoAttacks.NewPPMManager(1, procMask)
		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Landslide",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Landslide") {
					if spell.IsMH() {
						mainHand.Activate(sim)
					} else {
						offHand.Activate(sim)
					}
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(4099, 1, &ppmm, aura)
	})

	// Enchant: 4115, Spell: 75172 - Lightweave Embroidery
	core.NewEnchantEffect(4115, func(agent core.Agent) {
		character := agent.GetCharacter()

		statAura := character.NewTemporaryStatsAura(
			"Lightweave Embroidery Proc",
			core.ActionID{SpellID: 75171},
			stats.Stats{stats.Intellect: 480},
			time.Second*15,
		)

		handler := func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && !spell.ProcMask.Matches(core.ProcMaskSpellDamage|core.ProcMaskSpellHealing) {
				return
			}

			if aura.Icd.IsReady(sim) && sim.Proc(0.35, "Lightweave Eombridery Cata") {
				aura.Icd.Use(sim)
				statAura.Activate(sim)
			}
		}

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Lightweave Embroidery Cata",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			Icd: &core.Cooldown{
				Duration: time.Second * 60,
				Timer:    character.NewTimer(),
			},

			OnSpellHitDealt:       handler,
			OnPeriodicDamageDealt: handler,
			OnHealDealt:           handler,
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4115, aura)
	})

	// Enchant: 4179, Spell: 82175 - Synapse Springs
	core.NewEnchantEffect(4179, func(agent core.Agent) {
		character := agent.GetCharacter()

		agiAura := character.NewTemporaryStatsAura(
			"Hyperspeed Acceleration - Agi",
			core.ActionID{SpellID: 96228},
			stats.Stats{stats.Agility: 480},
			time.Second*10,
		)

		strAura := character.NewTemporaryStatsAura(
			"Hyperspeed Acceleration - Str",
			core.ActionID{SpellID: 96229},
			stats.Stats{stats.Strength: 480},
			time.Second*10,
		)

		intAura := character.NewTemporaryStatsAura(
			"Hyperspeed Acceleration - Int",
			core.ActionID{SpellID: 96230},
			stats.Stats{stats.Intellect: 480},
			time.Second*10,
		)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 82175},
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolMagic,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *core.Simulation, unit *core.Unit, _ *core.Spell) {
				intStat := unit.GetStat(stats.Intellect)
				strStat := unit.GetStat(stats.Strength)
				agiStat := unit.GetStat(stats.Agility)
				if intStat > strStat && intStat > agiStat {
					intAura.Activate(sim)
				} else if agiStat > intStat && agiStat > strStat {
					agiAura.Activate(sim)
				} else {
					strAura.Activate(sim)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Enchant: 4180, Spell: 82177 - Quickflip Deflection Plates
	core.NewEnchantEffect(4180, func(agent core.Agent) {
		character := agent.GetCharacter()
		statAura := character.NewTemporaryStatsAura(
			"Quickflip Deflection Plates Buff",
			core.ActionID{SpellID: 82176},
			stats.Stats{stats.Armor: 1500},
			time.Second*12,
		)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 82177},
			SpellSchool: core.SpellSchoolPhysical | core.SpellSchoolMagic,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				statAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeSurvival,
		})
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

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Gnomish X-Ray Scope",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			Icd: &core.Cooldown{
				Duration: time.Second * 40,
				Timer:    character.NewTimer(),
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}

				if aura.Icd.IsReady(sim) && sim.Proc(0.1, "Gnomish X-Ray Scope") {
					aura.Icd.Use(sim)
					statAura.Activate(sim)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4175, aura)
	})

	// Enchant: 4181, Spell: 82180 - Tazik Shocker
	core.NewEnchantEffect(4179, func(agent core.Agent) {
		character := agent.GetCharacter()
		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 82175},
			SpellSchool: core.SpellSchoolNature,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 120,
				},
			},

			ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
				// Benerfits from enhancement mastery
				// Ele crit dmg multi
				// Moonkin eclipse, so basically everything
				spell.CalcAndDealDamage(sim, unit, sim.Roll(4320, 961), spell.OutcomeMagicHitAndCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Enchant: 4182, Spell: 82200 - Spinal Healing Injector
	core.NewEnchantEffect(4182, func(agent core.Agent) {
		character := agent.GetCharacter()
		healthMetric := character.NewHealthMetrics(core.ActionID{SpellID: 82184})
		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 82200},
			SpellSchool: core.SpellSchoolNone,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagCombatPotion,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
				result := sim.Roll(27000, 33000)
				if character.HasAlchStone() {
					result *= 1.4
				}

				character.GainHealth(sim, result, healthMetric)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeSurvival,
		})
	})

	// Enchant: 4183, Spell: 82201 - Z50 Mana Gulper
	core.NewEnchantEffect(4183, func(agent core.Agent) {
		character := agent.GetCharacter()
		manaMetric := character.NewManaMetrics(core.ActionID{SpellID: 82186})
		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 82201},
			SpellSchool: core.SpellSchoolNone,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPotion,

			// TODO: In theory those ingi on-use enchants share a CD with potions
			// The potion CD timer is not available right now
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
				mana := sim.Roll(10730, 12470)
				if character.HasAlchStone() {
					mana *= 1.4
				}

				character.AddMana(sim, mana, manaMetric)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return c.HasManaBar() && (c.MaxMana()-c.CurrentMana()) > 10730
			},
			Spell:    spell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeMana,
		})
	})

	// Enchant: 4187, Spell: 84424 - Invisibility Field -- Non combat enchant won't implement
	// Enchant: 4214, Spell: 84425 - Cardboard Assassin -- Non combat enchant won't implement
	// Enchant: 4188, Spell: 84427 - Grounded Plasma Shield -- Shield Absorbs not supported?

	// Enchant: 4215, Spell: 92433, Item: 55055 - Elementium Shield Spike
	core.NewEnchantEffect(4215, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 55055}

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

		aura := character.RegisterAura(core.Aura{
			Label:    "Elementium Shield Spike",
			ActionID: actionID,
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Outcome.Matches(core.OutcomeBlock) && spell.ProcMask.Matches(core.ProcMaskMelee) {
					procSpell.Cast(sim, spell.Unit)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4215, aura)
	})

	// Enchant: 4216, Spell: 92437, Item: 55056  - Pyrium Shield Spike
	core.NewEnchantEffect(4216, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: 55055}

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

		aura := character.RegisterAura(core.Aura{
			Label:    "Pyrium Shield Spike",
			ActionID: actionID,
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Outcome.Matches(core.OutcomeBlock) && spell.ProcMask.Matches(core.ProcMaskMelee) {
					procSpell.Cast(sim, spell.Unit)
				}
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

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Flintlocke's Woodchucker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},

			Icd: &core.Cooldown{
				Duration: time.Second * 40,
				Timer:    character.NewTimer(),
			},

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskRanged) {
					return
				}

				if aura.Icd.IsReady(sim) && sim.Proc(0.1, "Flintlocke's Woodchucker") {
					aura.Icd.Use(sim)
					statAura.Activate(sim)
					dmgProc.Cast(sim, result.Target)
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(4267, aura)
	})
}
