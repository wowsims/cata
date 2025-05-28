package cata

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type ItemVersion int32

const (
	ItemVersionLFR ItemVersion = iota
	ItemVersionNormal
	ItemVersionHeroic
)

func init() {
	core.NewItemEffect(63839, func(agent core.Agent) {
		character := agent.GetCharacter()
		spreadDot := character.RegisterSpell(core.SpellConfig{
			ActionID:                 core.ActionID{SpellID: 91076},
			SpellSchool:              core.SpellSchoolNature,
			DamageMultiplier:         1,
			DamageMultiplierAdditive: 1,
			ThreatMultiplier:         1,
			CritMultiplier:           character.DefaultCritMultiplier(),
			ProcMask:                 core.ProcMaskEmpty,
			Flags:                    core.SpellFlagNoOnCastComplete,
			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Vengful Wisp - 2",
				},
				NumberOfTicks:       5,
				TickLength:          3 * time.Second,
				AffectedByCastSpeed: false,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 1176)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
					if sim.Proc(0.1, "Vengeful Wisp") {
						// select random proc target
						spreadTarget := sim.Encounter.TargetUnits[int(sim.Roll(0, float64(len(sim.Encounter.TargetUnits))))]

						// refresh dot on next step - refreshing potentially on aura expire
						// which will cause nasty things to happen
						core.StartDelayedAction(sim, core.DelayedActionOptions{
							DoAt:     sim.CurrentTime + 1,
							Priority: core.ActionPriorityDOT,
							OnAction: func(s *core.Simulation) {
								dot.Spell.Dot(spreadTarget).Apply(s)
							},
						})

					}
				},
			},
		})

		trinketDot := character.RegisterSpell(core.SpellConfig{
			ActionID:                 core.ActionID{SpellID: 91075},
			SpellSchool:              core.SpellSchoolNature,
			DamageMultiplier:         1,
			DamageMultiplierAdditive: 1,
			ThreatMultiplier:         1,
			ProcMask:                 core.ProcMaskEmpty,
			Flags:                    core.SpellFlagNoOnCastComplete,
			CritMultiplier:           character.DefaultCritMultiplier(),
			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Vengful Wisp - 1",
				},
				NumberOfTicks:       5,
				TickLength:          3 * time.Second,
				AffectedByCastSpeed: false,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 1176)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

					if sim.Proc(0.1, "Vengeful Wisp") {
						// select random proc target
						spreadTarget := sim.Encounter.TargetUnits[int(sim.Roll(0, float64(len(sim.Encounter.TargetUnits))))]
						spreadDot.Dot(spreadTarget).Apply(sim) // refresh self on
					}
				},
			},
		})

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Vengful Wisp",
			ActionID:   core.ActionID{ItemID: 63839},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.1,
			ICD:        time.Second * 100,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				trinketDot.Dot(result.Target).Apply(sim)
			},
		})

		character.ItemSwap.RegisterProc(63839, triggerAura)
	})

	core.NewItemEffect(64645, func(agent core.Agent) {
		character := agent.GetCharacter()
		storedMana := 0.0

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Tyrande's Favorite Doll (Mana)",
			ActionID: core.ActionID{ItemID: 64645},
			Callback: core.CallbackOnCastComplete,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				// Only Mana is converted
				if !character.HasManaBar() || spell.DefaultCast.Cost == 0 {
					return
				}
				storedMana = min(4200, storedMana+spell.DefaultCast.Cost*0.2)
			},
		})
		oldOnReset := triggerAura.OnReset
		triggerAura.OnReset = func(aura *core.Aura, sim *core.Simulation) {
			storedMana = 0
			oldOnReset(aura, sim)
		}

		character.ItemSwap.RegisterProc(64645, triggerAura)

		sharedTimer := character.GetOffensiveTrinketCD()
		manaMetric := character.NewManaMetrics(core.ActionID{SpellID: 92601})
		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:                 core.ActionID{ItemID: 64645},
			SpellSchool:              core.SpellSchoolArcane,
			Flags:                    core.SpellFlagNoOnCastComplete,
			ProcMask:                 core.ProcMaskEmpty,
			DamageMultiplier:         1,
			DamageMultiplierAdditive: 1,
			ThreatMultiplier:         1,
			CritMultiplier:           character.DefaultCritMultiplier(),
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 1,
					Timer:    sharedTimer,
				},
			},
			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return character.HasManaBar()
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, storedMana, spell.OutcomeMagicHitAndCrit)
				}
				character.AddMana(sim, storedMana, manaMetric)
				storedMana = 0
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    spell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS | core.CooldownTypeMana,
		})
	})

	core.NewItemEffect(59461, func(agent core.Agent) {
		character := agent.GetCharacter()

		dummyAura := character.RegisterAura(core.Aura{
			Label:     "Raw Fury",
			ActionID:  core.ActionID{SpellID: 91832},
			Duration:  time.Second * 15,
			MaxStacks: 5,
		})

		triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Raw Fury Aura",
			ActionID:   core.ActionID{ItemID: 59461},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.5,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 5,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				dummyAura.Activate(sim)
				dummyAura.AddStack(sim)
			},
		}))

		character.ItemSwap.RegisterProc(59461, triggerAura)

		buffAura := character.NewTemporaryStatsAura("Forged Fury", core.ActionID{SpellID: 91836}, stats.Stats{stats.Strength: 1926}, time.Second*20)
		sharedCD := character.GetOffensiveTrinketCD()
		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 59461},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,
			Cast: core.CastConfig{
				SharedCD: core.Cooldown{
					Timer:    sharedCD,
					Duration: time.Second * 20,
				},
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dummyAura.Deactivate(sim)
				buffAura.Activate(sim)
			},
			ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
				return dummyAura.GetStacks() == 5
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    trinketSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
			BuffAura: buffAura,

			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return dummyAura.GetStacks() == 5
			},
		})
	})

	core.NewItemEffect(68996, func(agent core.Agent) {
		character := agent.GetCharacter()

		totalAbsorbed := 0.0
		stayWithdrawnAura := character.RegisterAura(core.Aura{
			Label:    "Stay Withdrawn",
			Duration: time.Second * 10,
			ActionID: core.ActionID{SpellID: 96993},
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				tickAmount := totalAbsorbed * 0.08

				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period:   time.Second * 2,
					NumTicks: 5,
					OnAction: func(sim *core.Simulation) {
						character.RemoveHealth(sim, tickAmount)

						// Hacky way to force a fake stacks change log for the timeline to mimic ticks
						if sim.Log != nil {
							stacks := aura.GetStacks()
							aura.Unit.Log(sim, "%s stacks: %d --> %d", aura.ActionID, stacks, stacks)
						}
					},
					CleanUp: func(sim *core.Simulation) {
						totalAbsorbed = 0
					},
				})
			},
		})

		maxShieldStrength := 56980.0
		absorbAura := character.RegisterAura(core.Aura{
			Label:     "Stay of Execution",
			ActionID:  core.ActionID{ItemID: 68996, SpellID: 96988},
			Duration:  time.Second * 30,
			MaxStacks: int32(maxShieldStrength),
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if totalAbsorbed > 0 {
					stayWithdrawnAura.Activate(sim)
					stacks := int32(totalAbsorbed * 0.08)
					if stacks > 0 {
						stayWithdrawnAura.MaxStacks = stacks
						stayWithdrawnAura.SetStacks(sim, stacks)
					}
				}
			},
		})

		character.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if absorbAura.IsActive() && result.Damage > 0 && totalAbsorbed < maxShieldStrength {
				remainingAbsorb := maxShieldStrength - totalAbsorbed
				absorbedDamage := min(result.Damage*0.2, remainingAbsorb)
				result.Damage -= absorbedDamage
				totalAbsorbed = min(maxShieldStrength, totalAbsorbed+absorbedDamage)
				absorbAura.SetStacks(sim, int32(totalAbsorbed))

				if sim.Log != nil {
					character.Log(sim, "Stay of Execution absorbed %.1f damage", absorbedDamage)
				}
			}
		})

		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 68996},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,
			Cast: core.CastConfig{
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 20,
				},
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				absorbAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    trinketSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeSurvival,
		})
	})

	for _, version := range []ItemVersion{ItemVersionNormal, ItemVersionHeroic} {
		heroic := version == ItemVersionHeroic
		labelSuffix := core.Ternary(heroic, " (Heroic)", "")

		leadenItemID := core.TernaryInt32(heroic, 56347, 55816)
		core.NewItemEffect(leadenItemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: core.TernaryInt32(heroic, 92184, 92179)}
			armorBonus := core.TernaryFloat64(heroic, 3420, 2580)

			procAura := character.NewTemporaryStatsAura("Leaden Despair Proc"+labelSuffix, actionID, stats.Stats{stats.Armor: armorBonus}, time.Second*10)

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 30,
			}
			procAura.Icd = &icd

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     "Leaden Despair Trigger" + labelSuffix,
				Callback: core.CallbackOnSpellHitTaken,
				ProcMask: core.ProcMaskDirect,
				ActionID: core.ActionID{ItemID: leadenItemID},
				Harmful:  true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.35 {
						icd.Use(sim)
						procAura.Activate(sim)
					}
				},
			})
		})

		heartItemID := core.TernaryInt32(heroic, 65110, 59514)
		core.NewItemEffect(heartItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			procAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Heart's Revelation" + labelSuffix,
					ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 92325, 91027)},
					Duration:  time.Second * 15,
					MaxStacks: 5,
				},
				BonusPerStack: stats.Stats{stats.SpellPower: core.TernaryFloat64(heroic, 87, 77)},
			})

			triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Heart of Ignacious Aura" + labelSuffix,
				ActionID:   core.ActionID{ItemID: heartItemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskSpellDamage,
				ProcChance: 1,
				Outcome:    core.OutcomeLanded,
				ICD:        time.Second * 2,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
					procAura.AddStack(sim)
				},
			}))

			character.ItemSwap.RegisterProc(heartItemID, triggerAura)

			hastePerStack := core.TernaryFloat64(heroic, 363, 321)
			buffAura := character.RegisterAura(core.Aura{
				Label:     "Heart's Judgement" + labelSuffix,
				ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 92328, 91041)},
				Duration:  time.Second * 20,
				MaxStacks: 5,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					deltaHasteRating := hastePerStack * float64(newStacks-oldStacks)
					character.AddStatDynamic(sim, stats.HasteRating, deltaHasteRating)
				},
			})

			sharedCD := character.GetOffensiveTrinketCD()
			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: heartItemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    sharedCD,
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					buffAura.Activate(sim)
					buffAura.SetStacks(sim, procAura.GetStacks())
					procAura.Deactivate(sim)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeDPS,
				ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
					return procAura.GetStacks() == 5
				},
			})

			character.AddStatProcBuff(heartItemID, procAura, false, core.TrinketSlots())
		})

		jarItemID := core.TernaryInt32(heroic, 65029, 59354)
		core.NewItemEffect(jarItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			manaReturn := core.TernaryFloat64(heroic, 7260, 6420)

			procAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Inner Eye" + labelSuffix,
					ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 91320, 92329)},
					Duration:  time.Second * 15,
					MaxStacks: 5,
					Icd: &core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Second * 30,
					},
				},
				BonusPerStack: stats.Stats{stats.Spirit: core.TernaryFloat64(heroic, 116, 103)},
			})

			triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Jar of Ancient Remedies Aura" + labelSuffix,
				ActionID:   core.ActionID{ItemID: jarItemID},
				Callback:   core.CallbackOnHealDealt,
				ProcMask:   core.ProcMaskSpellHealing,
				ProcChance: 1,
				Outcome:    core.OutcomeLanded,
				ICD:        time.Second * 2,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if procAura.Icd.IsReady(sim) {
						procAura.Activate(sim)
						procAura.AddStack(sim)
					}
				},
			}))

			character.ItemSwap.RegisterProc(jarItemID, triggerAura)

			manaMetric := character.NewManaMetrics(core.ActionID{SpellID: core.TernaryInt32(heroic, 92331, 91322)})
			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: jarItemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
					return character.HasManaBar()
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					procAura.Deactivate(sim)
					// can not regain stacks for 30 seconds
					procAura.Icd.Use(sim)
					character.AddMana(sim, manaReturn, manaMetric)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeMana,
			})
		})

		matrixItemID := core.TernaryInt32(heroic, 69150, 68994)
		core.NewItemEffect(matrixItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			bonusStats := core.TernaryFloat64(heroic, 1834, 1624)

			procAuraCrit := character.NewTemporaryStatsAura(
				"Matrix Restabilizer Crit Proc"+labelSuffix,
				core.ActionID{SpellID: core.TernaryInt32(heroic, 97140, 96978)},
				stats.Stats{stats.CritRating: bonusStats},
				time.Second*30)
			procAuraHaste := character.NewTemporaryStatsAura(
				"Matrix Restabilizer Haste Proc"+labelSuffix,
				core.ActionID{SpellID: core.TernaryInt32(heroic, 97139, 96977)},
				stats.Stats{stats.HasteRating: bonusStats},
				time.Second*30)
			procAuraMastery := character.NewTemporaryStatsAura(
				"Matrix Restabilizer Mastery Proc"+labelSuffix,
				core.ActionID{SpellID: core.TernaryInt32(heroic, 97141, 96979)},
				stats.Stats{stats.MasteryRating: bonusStats},
				time.Second*30)

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Second * 105,
			}

			procAuraCrit.Icd = &icd
			procAuraHaste.Icd = &icd
			procAuraMastery.Icd = &icd

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Matrix Restabilizer Trigger" + labelSuffix,
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged,
				ProcChance: 0.2,
				ActionID:   core.ActionID{ItemID: matrixItemID},
				Harmful:    true,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					if icd.IsReady(sim) {
						statType := character.GetHighestStatType([]stats.Stat{stats.CritRating, stats.HasteRating, stats.MasteryRating})
						switch statType {
						case stats.CritRating:
							procAuraCrit.Activate(sim)
						case stats.HasteRating:
							procAuraHaste.Activate(sim)
						case stats.MasteryRating:
							procAuraMastery.Activate(sim)
						default:
							panic("unexpected statType")
						}
						icd.Use(sim)
					}
				},
			})

			character.ItemSwap.RegisterProc(matrixItemID, triggerAura)
		})

		apparatusItemID := core.TernaryInt32(heroic, 69113, 68972)
		core.NewItemEffect(apparatusItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			bonusPerStack := core.TernaryFloat64(heroic, 575, 508)
			buffDuration := time.Second * 15

			buffAuraCrit := character.NewTemporaryStatBuffWithStacks(
				"Blessing of the Shaper Crit"+labelSuffix,
				core.ActionID{SpellID: 96928},
				stats.Stats{stats.CritRating: bonusPerStack},
				5,
				buffDuration)

			buffAuraHaste := character.NewTemporaryStatBuffWithStacks(
				"Blessing of the Shaper Haste"+labelSuffix,
				core.ActionID{SpellID: 96927},
				stats.Stats{stats.HasteRating: bonusPerStack},
				5,
				buffDuration)

			buffAuraMastery := character.NewTemporaryStatBuffWithStacks(
				"Blessing of the Shaper Mastery"+labelSuffix,
				core.ActionID{SpellID: 96929},
				stats.Stats{stats.MasteryRating: bonusPerStack},
				5,
				buffDuration)

			buffAura := character.RegisterAura(core.Aura{
				Label:    "Apparatus of Khaz'goroth" + labelSuffix,
				ActionID: core.ActionID{ItemID: apparatusItemID},
				Duration: buffDuration,
			})

			titanicPower := character.RegisterAura(core.Aura{
				Label:     "Titanic Power" + labelSuffix,
				ActionID:  core.ActionID{SpellID: 96923},
				Duration:  time.Second * 30,
				MaxStacks: 5,
			})

			offHandBlacklist := map[int32]bool{
				// DK: Threat of Thassarian crits doesn't proc
				49143: true, // Frost Strike
				49998: true, // Death Strike
				85948: true, // Festering Strike
				49020: true, // Obliterate
				45462: true, // Plague Strike
				56815: true, // Rune Strike

				// Warrior: Whirlwind off-hand crits doesn't trigger procs
				1680: true, // Whirlwind
			}

			triggerAura := core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Titanic Power Trigger" + labelSuffix,
				ActionID:   core.ActionID{SpellID: 96924},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 1,
				Outcome:    core.OutcomeCrit,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if buffAuraCrit.IsActive() || buffAuraHaste.IsActive() || buffAuraMastery.IsActive() {
						return
					}

					// Warrior: Raging Blow crits doesn't trigger procs (both MH and OH)
					if spell.ActionID.SpellID == 85288 {
						return
					}

					// Off-hand blacklist
					if _, blacklisted := offHandBlacklist[spell.ActionID.SpellID]; spell.ProcMask.Matches(core.ProcMaskMeleeOHSpecial) && blacklisted {
						return
					}

					titanicPower.Activate(sim)
					titanicPower.AddStack(sim)
				},
			}))

			character.ItemSwap.RegisterProc(apparatusItemID, triggerAura)

			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: apparatusItemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    character.GetOffensiveTrinketCD(),
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					statType := character.GetHighestStatType([]stats.Stat{stats.CritRating, stats.HasteRating, stats.MasteryRating})

					switch statType {
					case stats.CritRating:
						buffAuraCrit.Activate(sim)
						buffAuraCrit.SetStacks(sim, titanicPower.GetStacks())
					case stats.HasteRating:
						buffAuraHaste.Activate(sim)
						buffAuraHaste.SetStacks(sim, titanicPower.GetStacks())
					case stats.MasteryRating:
						buffAuraMastery.Activate(sim)
						buffAuraMastery.SetStacks(sim, titanicPower.GetStacks())
					default:
						panic("unexpected statType")
					}

					buffAura.Activate(sim)
					titanicPower.Deactivate(sim)
				},
				ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
					return titanicPower.IsActive()
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeDPS,
				ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
					return titanicPower.IsActive()
				},
			})
		})

		vesselItemID := core.TernaryInt32(heroic, 69167, 68995)
		core.NewItemEffect(vesselItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			procAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Accelerated" + labelSuffix,
					ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 97142, 96980)},
					Duration:  time.Second * 20,
					MaxStacks: 5,
				},
				BonusPerStack: stats.Stats{stats.CritRating: core.TernaryFloat64(heroic, 92, 82)},
			})

			offHandBlacklist := map[int32]bool{
				// Warrior: Slam and Whirlwind off-hand crits doesn't trigger procs
				1464: true, // Slam
				1680: true, // Whirlwind
			}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{ItemID: vesselItemID},
				Name:       "Vessel of Acceleration" + labelSuffix,
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrMeleeProc,
				Outcome:    core.OutcomeCrit,
				ProcChance: 1,
				Harmful:    false,
				Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					// Off-hand blacklist
					if _, blacklisted := offHandBlacklist[spell.ActionID.SpellID]; spell.ProcMask.Matches(core.ProcMaskMeleeOHSpecial) && blacklisted {
						return
					}

					procAura.Activate(sim)
					procAura.AddStack(sim)
				},
			})

			character.ItemSwap.RegisterProc(vesselItemID, triggerAura)
		})

		jawsItemID := core.TernaryInt32(heroic, 69111, 68926)
		core.NewItemEffect(jawsItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			if !character.HasManaBar() {
				return
			}

			manaMod := character.AddDynamicMod(core.SpellModConfig{
				School:   core.SpellSchoolHoly | core.SpellSchoolNature,
				IntValue: 0,
				Kind:     core.SpellMod_PowerCost_Flat,
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Millisecond * 900,
			}

			manaReturn := core.TernaryInt32(heroic, -115, -110)

			victoriousAura := character.GetOrRegisterAura(core.Aura{
				Label:     "Victorious" + labelSuffix,
				ActionID:  core.ActionID{SpellID: core.TernaryInt32(heroic, 97120, 96907)},
				Duration:  time.Second * 20,
				MaxStacks: 10,
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.ProcMask&(core.ProcMaskSpellDamage|core.ProcMaskSpellHealing) == 0 {
						return
					}

					if !icd.IsReady(sim) {
						return
					}

					icd.Use(sim)

					aura.AddStack(sim)
				},
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Deactivate(sim)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					manaMod.Activate()
				},
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					manaMod.UpdateIntValue(manaReturn * newStacks)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					manaMod.Deactivate()
					manaMod.UpdateIntValue(0)
				},
			})

			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: jawsItemID},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    character.GetOffensiveTrinketCD(),
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					victoriousAura.Activate(sim)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeMana,
			})
		})

		spindleItemID := core.TernaryInt32(heroic, 69138, 68981)
		core.NewItemEffect(spindleItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			shieldStrength := core.TernaryFloat64(heroic, 19283, 17095)
			actionID := core.ActionID{ItemID: spindleItemID, SpellID: core.TernaryInt32(heroic, 97129, 96945)}
			duration := time.Second * 30

			shield := character.NewDamageAbsorptionAura("Loom of Fate"+labelSuffix, actionID, duration, func(unit *core.Unit) float64 {
				return shieldStrength
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute,
			}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Spidersilk Spindle Trigger" + labelSuffix,
				Callback:   core.CallbackOnSpellHitTaken,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					preHitHp := character.CurrentHealth() + result.Damage
					if icd.IsReady(sim) && spell.SpellSchool == core.SpellSchoolPhysical &&
						character.CurrentHealthPercent() < 0.35 && preHitHp >= 0.35 {
						icd.Use(sim)
						shield.Activate(sim)
					}
				},
			})

			character.ItemSwap.RegisterProc(spindleItemID, triggerAura)
		})

		scalesOfLifeItemID := core.TernaryInt32(heroic, 69109, 68915)
		core.NewItemEffect(scalesOfLifeItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			// Assuming full stack since sim doesn't track overhealing
			maxHeal := core.TernaryFloat64(heroic, 19283, 17095)
			trinketSpell := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: scalesOfLifeItemID},
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskSpellHealing,
				Flags:       core.SpellFlagNoOnCastComplete,

				Cast: core.CastConfig{
					SharedCD: core.Cooldown{
						Timer:    character.GetDefensiveTrinketCD(),
						Duration: time.Second * 20,
					},
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 1,
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealHealing(sim, spell.Unit, maxHeal, spell.OutcomeHealingCrit)
				},
			})

			character.AddMajorCooldown(core.MajorCooldown{
				Spell:    trinketSpell,
				Priority: core.CooldownPriorityDefault,
				Type:     core.CooldownTypeSurvival,
			})
		})
	}

	for _, version := range []ItemVersion{ItemVersionLFR, ItemVersionNormal, ItemVersionHeroic} {
		labelSuffix := []string{" (LFR)", "", " (Heroic)"}[version]

		vialItemID := []int32{77979, 77207, 77999}[version]
		core.NewItemEffect(vialItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			actionID := core.ActionID{SpellID: []int32{109721, 107994, 109724}[version]}
			minDmg := []float64{3568, 4028, 4546}[version]
			maxDmg := []float64{5353, 6042, 6819}[version]

			// 01/26/25: Testing during the first Dragon Soul PTR is consistent with these scaling
			// coefficients from simC, refer to Discord discussion from here onwards:
			// https://discord.com/channels/891730968493305867/1034670402150092841/1332728469427322982
			apMod := []float64{0.266, 0.3, 0.339}[version]

			// The AP scaling calculation can use either Melee or
			// Ranged AP depending on how the proc was triggered, so
			// keep track of it separately.
			var apSnapshot float64

			lightningStrike := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty, // ProcMask is set in the trigger aura
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(), // even ranged procs use the melee Crit multiplier as of first PTR (1.5x for Hunters)
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := sim.Roll(minDmg, maxDmg) + apMod*apSnapshot
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Vial of Shadows Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: vialItemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.45,
				ICD:        time.Second * 9,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskMelee | core.ProcMaskMeleeProc) {
						apSnapshot = spell.MeleeAttackPower()
					} else {
						apSnapshot = spell.RangedAttackPower()
					}
					lightningStrike.ProcMask = core.Ternary(spell.ProcMask.Matches(core.ProcMaskRanged), core.ProcMaskRangedProc, core.ProcMaskMeleeProc)
					lightningStrike.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(vialItemID, triggerAura)

		})

		fetishItemID := []int32{77982, 77210, 78002}[version]
		core.NewItemEffect(fetishItemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			numTargets := character.Env.GetNumTargets()

			actionID := core.ActionID{SpellID: []int32{109753, 107998, 109755}[version]}
			minDmg := []float64{8029, 9063, 10230}[version]
			maxDmg := []float64{12044, 13594, 15345}[version]
			apMod := []float64{0.598, 0.675, 0.762}[version]

			whirlingMaw := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskMeleeProc,
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					results := make([]*core.SpellResult, numTargets)

					for idx := int32(0); idx < numTargets; idx++ {
						baseDamage := sim.Roll(minDmg, maxDmg) +
							apMod*spell.MeleeAttackPower()
						results[idx] = spell.CalcDamage(sim, sim.Environment.GetTargetUnit(idx), baseDamage, spell.OutcomeMeleeSpecialCritOnly)
					}

					for idx := int32(0); idx < numTargets; idx++ {
						spell.DealDamage(sim, results[idx])
					}
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Bone-Link Fetish Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: fetishItemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.15,
				ICD:        time.Second * 27,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					whirlingMaw.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(fetishItemID, triggerAura)
		})

		cunningItemID := []int32{77980, 77208, 78000}[version]
		core.NewItemEffect(cunningItemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			numTargets := character.Env.GetNumTargets()

			actionID := core.ActionID{SpellID: []int32{109798, 108005, 109800}[version]}
			minDmg := []float64{2498, 2820, 3183}[version]
			maxDmg := []float64{3747, 4230, 4774}[version]
			spMod := []float64{0.277, 0.313, 0.353}[version]

			shadowboltVolley := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolShadow,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				MissileSpeed: 20,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				BonusCoefficient: spMod,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					results := make([]*core.SpellResult, numTargets)

					for idx := int32(0); idx < numTargets; idx++ {
						results[idx] = spell.CalcDamage(sim, sim.Environment.GetTargetUnit(idx), sim.Roll(minDmg, maxDmg), spell.OutcomeMagicCrit)
					}

					spell.WaitTravelTime(sim, func(sim *core.Simulation) {
						for idx := int32(0); idx < numTargets; idx++ {
							spell.DealDamage(sim, results[idx])
						}
					})
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Cunning of the Cruel Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: cunningItemID},
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				ProcMask:   core.ProcMaskSpellOrSpellProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.45,
				ICD:        time.Second * 9,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shadowboltVolley.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(cunningItemID, triggerAura)

		})

		prideItemID := []int32{77983, 77211, 78003}[version]
		core.NewItemEffect(prideItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			absorbModifier := []float64{0.43, 0.50, 0.56}[version]
			actionID := core.ActionID{ItemID: prideItemID, SpellID: 108008}
			duration := time.Second * 6

			shieldStrength := 0.0
			shield := character.NewDamageAbsorptionAura("Indomitable"+labelSuffix, actionID, duration, func(unit *core.Unit) float64 {
				return shieldStrength
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Minute,
			}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Indomitable Pride Trigger" + labelSuffix,
				Callback:   core.CallbackOnSpellHitTaken,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					preHitHp := character.CurrentHealth() + result.Damage
					if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.50 && preHitHp >= 0.50 {
						shieldStrength = result.Damage * absorbModifier
						if shieldStrength > 1 {
							icd.Use(sim)
							shield.Activate(sim)
						}
					}
				},
			})

			character.ItemSwap.RegisterProc(prideItemID, triggerAura)
		})

		// Kiril, Fury of Beasts
		// Equip: Your melee and ranged attacks have a chance to trigger Fury of the Beast, granting 107 Agility and 10% increased size every 1 sec.
		// This effect stacks a maximum of 10 times and lasts 20 sec.
		// (Proc chance: 15%, 55s cooldown)
		// TODO: Verify if the aura is cancelled when swapping druid forms
		// Video from 4.3.0 showing that it doesn't: https://www.youtube.com/watch?v=A6PYbDRaH6E
		// Comment from 4.3.3 stating that it does: https://www.wowhead.com/mop-classic/item=77194/kiril-fury-of-beasts#comments:id=1639024
		kirilItemID := []int32{78482, 77194, 78473}[version]
		core.NewItemEffect(kirilItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			beastFuryAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Beast Fury" + labelSuffix,
					ActionID:  core.ActionID{SpellID: []int32{109860, 108016, 109863}[version]},
					Duration:  time.Second * 20,
					MaxStacks: 10,
				},
				BonusPerStack: stats.Stats{stats.Agility: []float64{95, 107, 120}[version]},
			})

			furyOfTheBeastAura := character.RegisterAura(core.Aura{
				Label:    "Fury of the Beast" + labelSuffix,
				ActionID: core.ActionID{SpellID: []int32{109861, 108011, 109864}[version]},
				Duration: time.Second * 20,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					beastFuryAura.Activate(sim)
					core.StartPeriodicAction(sim, core.PeriodicActionOptions{
						Period:   time.Second,
						NumTicks: 10,
						OnAction: func(sim *core.Simulation) {
							beastFuryAura.AddStack(sim)
						},
					})
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Fury of the Beast Trigger" + labelSuffix,
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskMeleeProc,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.15,
				ICD:        time.Second * 55,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					furyOfTheBeastAura.Activate(sim)
				},
			})

			character.ItemSwap.RegisterProc(kirilItemID, triggerAura)
		})

		// These spells ignore the slot the weapon is in.
		// Any other ability should only trigger the proc if the weapon is in the right slot.
		ignoresSlot := map[int32]bool{
			23881: true, // Bloodthirst
			6544:  true, // Heroic Leap
		}

		// Souldrinker
		// Equip: Your melee attacks have a chance to drain your target's health, damaging the target for an amount equal to 1.3%/1.5%/1.7% of your maximum health and healing you for twice that amount.
		// (Proc chance: 15%)
		souldrinkerItemID := []int32{78488, 77193, 78479}[version]
		core.NewItemEffect(souldrinkerItemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: []int32{109828, 108022, 109831}[version]}
			label := fmt.Sprintf("Drain Life Trigger %s", labelSuffix)
			hpModifier := []float64{0.013, 0.015, 0.017}[version]
			meleeWeaponSlots := core.MeleeWeaponSlots()

			var damageDealt float64
			drainLifeHeal := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID.WithTag(2),
				SpellSchool: core.SpellSchoolShadow,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealHealing(sim, target, damageDealt*2, spell.OutcomeAlwaysHit)
				},
			})

			drainLife := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID.WithTag(1),
				SpellSchool: core.SpellSchoolShadow,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := character.MaxHealth() * hpModifier
					damageDealt = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit).Damage

					drainLifeHeal.Cast(sim, &character.Unit)
				},
			})

			makeProcTrigger := func(character *core.Character, isMH bool) {
				itemSlot := core.Ternary(isMH, meleeWeaponSlots[:1], meleeWeaponSlots[1:])
				procMask := core.Ternary(isMH, core.ProcMaskMeleeMH|core.ProcMaskMeleeProc, core.ProcMaskMeleeOH)

				aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
					Name:     fmt.Sprintf("%s %s", label, core.Ternary(isMH, "MH", "OH")),
					ActionID: core.ActionID{ItemID: souldrinkerItemID},
					ProcMask: core.ProcMaskMelee,
					Outcome:  core.OutcomeLanded,
					Callback: core.CallbackOnSpellHitDealt,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if _, ignore := ignoresSlot[spell.ActionID.SpellID]; !spell.ProcMask.Matches(procMask) && !ignore {
							return
						}

						if sim.Proc(0.15, label) {
							drainLife.Cast(sim, result.Target)
						}
					},
				})

				character.ItemSwap.RegisterProcWithSlots(souldrinkerItemID, aura, itemSlot)
			}

			if character.ItemSwap.CouldHaveItemEquippedInSlot(souldrinkerItemID, proto.ItemSlot_ItemSlotMainHand) {
				makeProcTrigger(character, true)
			}
			if character.ItemSwap.CouldHaveItemEquippedInSlot(souldrinkerItemID, proto.ItemSlot_ItemSlotOffHand) {
				makeProcTrigger(character, false)
			}
		})

		// No'Kaled, the Elements of Death
		// Equip: Your melee attacks have a chance to blast your enemy with Fire, Shadow, or Frost, dealing 6781/7654/8640 to 10171/11481/12960 damage.
		// (Proc chance: 7%)
		nokaledItemID := []int32{78481, 77188, 78472}[version]
		core.NewItemEffect(nokaledItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			defaultProcMask := core.ProcMaskMeleeProc
			procMask := character.GetDynamicProcMaskForWeaponEffect(nokaledItemID)

			minDamage := []float64{6781, 7654, 8640}[version]
			maxDamage := []float64{10171, 11481, 12960}[version]

			registerSpell := func(actionID core.ActionID, spellSchool core.SpellSchool) *core.Spell {
				return character.RegisterSpell(core.SpellConfig{
					ActionID:    actionID,
					SpellSchool: spellSchool,
					ProcMask:    core.ProcMaskEmpty,
					Flags:       core.SpellFlagPassiveSpell,
					MaxRange:    45,

					DamageMultiplier: 1,
					CritMultiplier:   character.DefaultCritMultiplier(),
					ThreatMultiplier: 1,

					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						baseDamage := sim.RollWithLabel(minDamage, maxDamage, "No'Kaled, the Elements of Death")
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
					},
				})
			}

			flameblast := registerSpell(
				core.ActionID{SpellID: []int32{109871, 107785, 109872}[version]},
				core.SpellSchoolFire)

			iceblast := registerSpell(
				core.ActionID{SpellID: []int32{109869, 107789, 109870}[version]},
				core.SpellSchoolFrost)

			shadowblast := registerSpell(
				core.ActionID{SpellID: []int32{109867, 107787, 109868}[version]},
				core.SpellSchoolShadow)

			spells := []*core.Spell{flameblast, iceblast, shadowblast}

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:     "No'Kaled Trigger" + labelSuffix,
				ActionID: core.ActionID{ItemID: nokaledItemID},
				Callback: core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}

					if _, ignore := ignoresSlot[spell.ActionID.SpellID]; !spell.ProcMask.Matches(*procMask|defaultProcMask) && !ignore {
						return
					}

					if sim.Proc(0.07, "No'Kaled, the Elements of Death") {
						spell := spells[int(sim.RollWithLabel(0, float64(len(spells)), "No'Kaled spell to cast"))]
						spell.Cast(sim, result.Target)
					}
				},
			})

			character.ItemSwap.RegisterProc(nokaledItemID, triggerAura)
		})

		// Rathrak, the Poisonous Mind
		// Equip: Your harmful spellcasts have a chance to poison all enemies near your target for 7715/8710/9830 nature damage over 10 sec.
		// (Proc chance: 15%, 17s cooldown)
		rathrakItemID := []int32{78484, 77195, 78475}[version]
		core.NewItemEffect(rathrakItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			tickDamage := []float64{7715, 8710, 9830}[version] / 5

			blastOfCorruption := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: []int32{109851, 107831, 109854}[version]},
				SpellSchool: core.SpellSchoolNature,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				Dot: core.DotConfig{
					IsAOE: true,
					Aura: core.Aura{
						Label: "Blast of Corruption" + labelSuffix,
					},
					NumberOfTicks:       5,
					TickLength:          time.Second * 2,
					AffectedByCastSpeed: false,

					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						for _, aoeTarget := range sim.Encounter.TargetUnits {
							result := dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, tickDamage, dot.Spell.OutcomeMagicCritNoHitCounter)

							if result.DidCrit() {
								dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
							} else {
								dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
							}
						}
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.AOEDot().Apply(sim)
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Rathrak Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: rathrakItemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskSpellOrSpellProc,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.15,
				ICD:        time.Second * 17,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					blastOfCorruption.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(rathrakItemID, triggerAura)

		})

		// Vishanka, Jaws of the Earth
		// Equip: Your ranged attacks have a chance to deal 7040/7950/8970 damage over 2 sec.
		// (Proc chance: 15%, 17s cooldown)
		// Time between ticks: 200ms
		vishankaItemID := []int32{78480, 78359, 78471}[version]
		core.NewItemEffect(vishankaItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			tickDamage := []float64{7040, 7950, 8970}[version] / 10

			speakingOfRage := character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: []int32{109856, 107821, 109858}[version]},
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Speaking of Rage" + labelSuffix,
					},
					NumberOfTicks:       10,
					TickLength:          time.Millisecond * 200,
					AffectedByCastSpeed: false,

					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						result := dot.Spell.CalcAndDealPeriodicDamage(sim, target, tickDamage, dot.Spell.OutcomeRangedCritOnlyNoHitCounter)

						if result.DidCrit() {
							dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
						} else {
							dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
						}
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.CritMultiplier(1, 0),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).Apply(sim)
				},
			})

			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Vishanka Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: vishankaItemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskRanged | core.ProcMaskRangedProc,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.15,
				ICD:        time.Second * 17,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					speakingOfRage.Cast(sim, result.Target)
				},
			})

			character.ItemSwap.RegisterProc(vishankaItemID, triggerAura)
		})

		// Ti'tahk, the Steps of Time
		// Equip: Your spells have a chance to grant you 1708/1928/2176 haste rating for 10 sec and 342/386/435 haste rating to up to 3 allies within 20 yards.
		// (Proc chance: 15%, 50s cooldown)
		// The buff has two effects, one for the caster and one shared.
		// * The first effect is 1366/1542/1741 haste rating on the caster.
		// * The second effect is 342/386/435 haste rating on the caster and up to 3 allies within 20 yards.
		// E.g. for the LFR version it's 1366 + 342 = 1708 haste rating for the caster with a shared ID so we just combine them.
		// TODO: Add the shared buff as an optional misc raid buff?
		//       Could be annoying with 3 different versions, uptime etc.
		titahkItemID := []int32{78486, 77190, 78477}[version]
		titahkAuraID := []int32{109842, 107804, 109844}[version]
		titahkBonus := []float64{1366 + 342, 1542 + 386, 1741 + 435}[version]
		titahkLabel := fmt.Sprintf("Ti'tahk, the Steps of Time %s", labelSuffix)
		core.NewItemEffect(titahkItemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura(titahkLabel, core.ActionID{SpellID: titahkAuraID}, stats.Stats{stats.HasteRating: titahkBonus}, time.Second*10)

			handler := func(triggerAura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOrMeleeProc | core.ProcMaskRangedOrRangedProc) {
					return
				}
				if triggerAura.Icd.IsReady(sim) && sim.Proc(0.15, fmt.Sprintf("%s Trigger", titahkLabel)) {
					procAura.Activate(sim)
					triggerAura.Icd.Use(sim)
				}
			}

			triggerAura := character.RegisterAura(core.Aura{
				Label:    fmt.Sprintf("%s Trigger", titahkLabel),
				ActionID: core.ActionID{ItemID: titahkItemID},
				Duration: core.NeverExpires,
				Icd: &core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 50,
				},
				OnSpellHitDealt:       handler,
				OnPeriodicDamageDealt: handler,
				OnHealDealt:           handler,
				OnPeriodicHealDealt:   handler,
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					handler(aura, sim, spell, nil)
				},
			})

			character.ItemSwap.RegisterProc(titahkItemID, triggerAura)
		})
	}
}

// Takes in the SpellResult for the triggering spell, and returns the damage per
// tick of a *fresh* Ignite triggered by that spell. Roll-over damage
// calculations for existing Ignites are handled internally.
type IgniteDamageCalculator func(result *core.SpellResult) float64

type IgniteConfig struct {
	ActionID           core.ActionID
	ClassSpellMask     int64
	DisableCastMetrics bool
	DotAuraLabel       string
	DotAuraTag         string
	ProcTrigger        core.ProcTrigger // Ignores the Handler field and creates a custom one, but uses all others.
	DamageCalculator   IgniteDamageCalculator
	IncludeAuraDelay   bool // "munching" and "free roll-over" interactions
	SpellSchool        core.SpellSchool
	NumberOfTicks      int32
	TickLength         time.Duration
	SetBonusAura       *core.Aura
}

func RegisterIgniteEffect(unit *core.Unit, config IgniteConfig) *core.Spell {
	spellFlags := core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete

	if config.DisableCastMetrics {
		spellFlags |= core.SpellFlagPassiveSpell
	}

	if config.SpellSchool == 0 {
		config.SpellSchool = core.SpellSchoolFire
	}

	if config.NumberOfTicks == 0 {
		config.NumberOfTicks = 2
	}

	if config.TickLength == 0 {
		config.TickLength = time.Second * 2
	}

	igniteSpell := unit.RegisterSpell(core.SpellConfig{
		ActionID:         config.ActionID,
		SpellSchool:      config.SpellSchool,
		ProcMask:         core.ProcMaskSpellProc,
		ClassSpellMask:   config.ClassSpellMask,
		Flags:            spellFlags,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     config.DotAuraLabel,
				Tag:       config.DotAuraTag,
				MaxStacks: math.MaxInt32,
			},

			NumberOfTicks:       config.NumberOfTicks,
			TickLength:          config.TickLength,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.Spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	refreshIgnite := func(sim *core.Simulation, target *core.Unit, damagePerTick float64) {
		// Cata Ignite
		// 1st ignite application = 4s, split into 2 ticks (2s, 0s)
		// Ignite refreshes: Duration = 4s + MODULO(remaining duration, 2), max 6s. Split damage over 3 ticks at 4s, 2s, 0s.
		dot := igniteSpell.Dot(target)
		dot.SnapshotBaseDamage = damagePerTick
		igniteSpell.Cast(sim, target)
		dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
	}

	var scheduledRefresh *core.PendingAction
	procTrigger := config.ProcTrigger
	procTrigger.Handler = func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
		target := result.Target
		dot := igniteSpell.Dot(target)
		outstandingDamage := dot.OutstandingDmg()
		newDamage := config.DamageCalculator(result)
		totalDamage := outstandingDamage + newDamage
		newTickCount := dot.BaseTickCount + core.TernaryInt32(dot.IsActive(), 1, 0)
		damagePerTick := totalDamage / float64(newTickCount)

		if config.IncludeAuraDelay {
			// Rough 2-bucket model for the aura update delay distribution based
			// on PTR measurements. Most updates occur on either the same or very
			// next spell batch after the proc, and can therefore be modeled by a
			// 0-10 ms random draw. But a reasonable minority fraction take ~10x
			// longer than this to fire. The origin of these longer delays is
			// likely not actually random in reality, but can be treated that way
			// in practice since the player cannot play around them.
			var delaySeconds float64

			if sim.Proc(0.75, "Aura Delay") {
				delaySeconds = 0.010 * sim.RandomFloat("Aura Delay")
			} else {
				delaySeconds = 0.090 + 0.020*sim.RandomFloat("Aura Delay")
			}

			applyDotAt := sim.CurrentTime + core.DurationFromSeconds(delaySeconds)

			// Cancel any prior aura updates already in the queue
			if (scheduledRefresh != nil) && (scheduledRefresh.NextActionAt > sim.CurrentTime) {
				scheduledRefresh.Cancel(sim)

				if sim.Log != nil {
					unit.Log(sim, "Previous %s proc was munched due to server aura delay", config.DotAuraLabel)
				}
			}

			// Schedule a delayed refresh of the DoT with cached damagePerTick value (allowing for "free roll-overs")
			if sim.Log != nil {
				unit.Log(sim, "Schedule travel (%0.1f ms) for %s", delaySeconds*1000, config.DotAuraLabel)

				if dot.IsActive() && (dot.NextTickAt() < applyDotAt) {
					unit.Log(sim, "%s rolled with %0.3f damage both ticking and rolled into next", config.DotAuraLabel, outstandingDamage)
				}
			}

			scheduledRefresh = core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt:     applyDotAt,
				Priority: core.ActionPriorityDOT,

				OnAction: func(_ *core.Simulation) {
					refreshIgnite(sim, target, damagePerTick)
				},
			})
		} else {
			refreshIgnite(sim, target, damagePerTick)
		}
	}

	if config.SetBonusAura != nil {
		config.SetBonusAura.AttachProcTrigger(procTrigger)
	} else {
		core.MakeProcTriggerAura(unit, procTrigger)
	}

	return igniteSpell
}
