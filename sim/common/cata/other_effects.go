package cata

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
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
			CritMultiplier:           character.DefaultSpellCritMultiplier(),
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
			CritMultiplier:           character.DefaultSpellCritMultiplier(),
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

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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
	})

	core.NewItemEffect(64645, func(agent core.Agent) {
		character := agent.GetCharacter()
		storedMana := 0.0

		core.MakePermanent(character.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: 64645},
			Label:    "Tyrande's Favorite Doll (Mana)",
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				storedMana = 0.0
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				// Only Mana is converted
				if !character.HasManaBar() {
					return
				}

				storedMana = math.Min(4200, storedMana+spell.DefaultCast.Cost*0.2)
			},
		}))

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
			CritMultiplier:           character.DefaultSpellCritMultiplier(),
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 1,
					Timer:    sharedTimer,
				},
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, storedMana, spell.OutcomeMagicHitAndCrit)
				}

				character.AddMana(sim, storedMana, manaMetric)
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

		core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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

			titanicPower := character.RegisterAura(core.Aura{
				Label:     "Titanic Power" + labelSuffix,
				ActionID:  core.ActionID{SpellID: 96923},
				Duration:  time.Second * 30,
				MaxStacks: 5,
			})

			core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Titanic Power Aura" + labelSuffix,
				ActionID:   core.ActionID{SpellID: 96924},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 1,
				Outcome:    core.OutcomeCrit,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if buffAuraCrit.IsActive() || buffAuraHaste.IsActive() || buffAuraMastery.IsActive() {
						return
					}

					titanicPower.Activate(sim)
					titanicPower.AddStack(sim)
				},
			}))

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

		jawsItemID := core.TernaryInt32(heroic, 69111, 68926)
		core.NewItemEffect(jawsItemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			if !character.HasManaBar() {
				return
			}

			manaMod := character.AddDynamicMod(core.SpellModConfig{
				School:     core.SpellSchoolHoly | core.SpellSchoolNature,
				FloatValue: 0,
				Kind:       core.SpellMod_PowerCost_Flat,
			})

			icd := core.Cooldown{
				Timer:    character.NewTimer(),
				Duration: time.Millisecond * 900,
			}

			manaReturn := core.TernaryFloat64(heroic, -115, -110)

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
					manaMod.UpdateFloatValue(manaReturn * float64(newStacks))
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					manaMod.Deactivate()
					manaMod.UpdateFloatValue(0)
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

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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
				CritMultiplier:   character.DefaultSpellCritMultiplier(),
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

			lightningStrike := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				Cast: core.CastConfig{
					DefaultCast: core.Cast{
						NonEmpty: true,
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultMeleeCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, sim.Roll(minDmg, maxDmg), spell.OutcomeMeleeSpecialCritOnly)
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Vial of Shadows Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: vialItemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.45,
				ICD:        time.Second * 9,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					lightningStrike.Cast(sim, result.Target)
				},
			})
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
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagPassiveSpell,

				Cast: core.CastConfig{
					DefaultCast: core.Cast{
						NonEmpty: true,
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultMeleeCritMultiplier(),
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

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Bone-Link Fetish Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: fetishItemID},
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged | core.ProcMaskProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.15,
				ICD:        time.Second * 27,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					whirlingMaw.Cast(sim, result.Target)
				},
			})
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

				Cast: core.CastConfig{
					DefaultCast: core.Cast{
						NonEmpty: true,
					},
				},

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultSpellCritMultiplier(),
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

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Cunning of the Cruel Trigger" + labelSuffix,
				ActionID:   core.ActionID{ItemID: cunningItemID},
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskProc,
				Outcome:    core.OutcomeLanded,
				Harmful:    true,
				ProcChance: 0.45,
				ICD:        time.Second * 9,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shadowboltVolley.Cast(sim, result.Target)
				},
			})
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

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
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
		})
	}
}

// Takes in the SpellResult for the triggering spell, and returns the damage per
// tick of a *fresh* Ignite triggered by that spell. Roll-over damage
// calculations for existing Ignites are handled internally.
type IgniteDamageCalculator func(result *core.SpellResult) float64

type IgniteConfig struct {
	ActionID           core.ActionID
	DisableCastMetrics bool
	DotAuraLabel       string
	DotAuraTag         string
	ProcTrigger        core.ProcTrigger // Ignores the Handler field and creates a custom one, but uses all others.
	DamageCalculator   IgniteDamageCalculator
	IncludeAuraDelay   bool // "munching" and "free roll-over" interactions
}

func RegisterIgniteEffect(unit *core.Unit, config IgniteConfig) *core.Spell {
	spellFlags := core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete

	if config.DisableCastMetrics {
		spellFlags = spellFlags | core.SpellFlagPassiveSpell
	}

	igniteSpell := unit.RegisterSpell(core.SpellConfig{
		ActionID:         config.ActionID,
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskProc,
		Flags:            spellFlags,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     config.DotAuraLabel,
				Tag:       config.DotAuraTag,
				MaxStacks: math.MaxInt32,
			},

			NumberOfTicks:       2,
			TickLength:          time.Second * 2,
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

	refreshIgnite := func(sim *core.Simulation, target *core.Unit, totalDamage float64) {
		// Cata Ignite
		// 1st ignite application = 4s, split into 2 ticks (2s, 0s)
		// Ignite refreshes: Duration = 4s + MODULO(remaining duration, 2), max 6s. Split damage over 3 ticks at 4s, 2s, 0s.
		dot := igniteSpell.Dot(target)
		newTickCount := dot.BaseTickCount + core.TernaryInt32(dot.IsActive(), 1, 0)
		dot.SnapshotBaseDamage = totalDamage / float64(newTickCount)
		igniteSpell.Cast(sim, target)
		dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
	}

	var scheduledRefresh *core.PendingAction
	procTrigger := config.ProcTrigger
	procTrigger.Handler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		target := result.Target
		dot := igniteSpell.Dot(target)
		outstandingDamage := dot.OutstandingDmg()
		newDamage := config.DamageCalculator(result)
		totalDamage := outstandingDamage + newDamage

		if config.IncludeAuraDelay {
			// For now, assume that the mechanism driving random aura update
			// delays is the same as for random auto delays after rapidly
			// changing position. Therefore, use the fit delay parameters
			// from cat leap tests (see sim/druid/feral_charge.go) for
			// modeling consistency.
			// TODO: Measure the aura update delay distribution on PTR.
			waitTime := time.Millisecond * time.Duration(sim.Roll(150, 750))
			applyDotAt := sim.CurrentTime + waitTime

			// Check for max duration munching
			if dot.RemainingDuration(sim) > time.Second*4+waitTime {
				if sim.Log != nil {
					unit.Log(sim, "New %s proc was munched due to max %s duration", config.DotAuraLabel, config.DotAuraLabel)
				}

				return
			}

			// Cancel any prior aura updates already in the queue
			if (scheduledRefresh != nil) && (scheduledRefresh.NextActionAt > sim.CurrentTime) {
				scheduledRefresh.Cancel(sim)

				if sim.Log != nil {
					unit.Log(sim, "Previous %s proc was munched due to server aura delay", config.DotAuraLabel)
				}
			}

			// Schedule a delayed refresh of the DoT with cached outstandingDamage value (allowing for "free roll-overs")
			if sim.Log != nil {
				unit.Log(sim, "Schedule travel (%0.2f s) for %s", waitTime.Seconds(), config.DotAuraLabel)

				if dot.IsActive() && (dot.NextTickAt() < applyDotAt) {
					unit.Log(sim, "%s rolled with %0.3f damage both ticking and rolled into next", config.DotAuraLabel, outstandingDamage)
				}
			}

			scheduledRefresh = core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt:     applyDotAt,
				Priority: core.ActionPriorityDOT,

				OnAction: func(_ *core.Simulation) {
					refreshIgnite(sim, target, totalDamage)
				},
			})
		} else {
			refreshIgnite(sim, target, totalDamage)
		}
	}

	core.MakeProcTriggerAura(unit, procTrigger)
	return igniteSpell
}
