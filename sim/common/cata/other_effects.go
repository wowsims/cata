package cata

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {
	core.NewItemEffect(55816, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 92179}

		procAura := character.NewTemporaryStatsAura("Leaden Despair Proc", actionID, stats.Stats{stats.Armor: 2580}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 30,
		}
		procAura.Icd = &icd

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Leaden Despair Trigger",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskDirect,
			ActionID: core.ActionID{ItemID: 55816},
			Harmful:  true,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.35 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

	core.NewItemEffect(56347, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 92184}

		procAura := character.NewTemporaryStatsAura("Leaden Despair (Heroic) Proc", actionID, stats.Stats{stats.Armor: 3420}, time.Second*10)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 30,
		}
		procAura.Icd = &icd

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Leaden Despair (Heroic) Trigger",
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskDirect,
			ActionID: core.ActionID{ItemID: 56347},
			Harmful:  true,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				if icd.IsReady(sim) && character.CurrentHealthPercent() < 0.35 {
					icd.Use(sim)
					procAura.Activate(sim)
				}
			},
		})
	})

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

	core.NewItemEffect(59514, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Heart's Revelation",
				ActionID:  core.ActionID{SpellID: 91027},
				Duration:  time.Second * 15,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.SpellPower: 77},
		})

		core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Heart of Ignacious Aura",
			ActionID:   core.ActionID{ItemID: 59514},
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

		buffAura := character.RegisterAura(core.Aura{
			Label:     "Heart's Judgement",
			ActionID:  core.ActionID{SpellID: 91041},
			Duration:  time.Second * 20,
			MaxStacks: 5,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				deltaHasteRating := float64(321) * float64(newStacks-oldStacks)
				character.AddStatDynamic(sim, stats.HasteRating, deltaHasteRating)
			},
		})

		sharedCD := character.GetOffensiveTrinketCD()
		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 59514},
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

	core.NewItemEffect(65110, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Heart's Revelation (Heroic)",
				ActionID:  core.ActionID{SpellID: 92325},
				Duration:  time.Second * 15,
				MaxStacks: 5,
			},
			BonusPerStack: stats.Stats{stats.SpellPower: 87},
		})

		core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Heart of Ignacious (Heroic) Aura",
			ActionID:   core.ActionID{ItemID: 65110},
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

		buffAura := character.RegisterAura(core.Aura{
			Label:     "Heart's Judgement (Heroic)",
			ActionID:  core.ActionID{SpellID: 92328},
			Duration:  time.Second * 20,
			MaxStacks: 5,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				deltaHasteRating := float64(363) * float64(newStacks-oldStacks)
				character.AddStatDynamic(sim, stats.HasteRating, deltaHasteRating)
			},
		})

		sharedCD := character.GetOffensiveTrinketCD()
		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 65110},
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

	core.NewItemEffect(59354, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Inner Eye",
				ActionID:  core.ActionID{SpellID: 91320},
				Duration:  time.Second * 15,
				MaxStacks: 5,
				Icd: &core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 30,
				},
			},
			BonusPerStack: stats.Stats{stats.Spirit: 103},
		})

		core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Jar of Ancient Remedies Aura",
			ActionID:   core.ActionID{ItemID: 59354},
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

		manaMetric := character.NewManaMetrics(core.ActionID{SpellID: 91322})
		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 59461},
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
				character.AddMana(sim, 6420, manaMetric)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    trinketSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeMana,
		})
	})

	core.NewItemEffect(65029, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Inner Eye (Heroic)",
				ActionID:  core.ActionID{SpellID: 92329},
				Duration:  time.Second * 15,
				MaxStacks: 5,
				Icd: &core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 30,
				},
			},
			BonusPerStack: stats.Stats{stats.Spirit: 116},
		})

		core.MakePermanent(core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Jar of Ancient Remedies (Heroic) Aura",
			ActionID:   core.ActionID{ItemID: 65029},
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

		manaMetric := character.NewManaMetrics(core.ActionID{SpellID: 92331})
		trinketSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 65029},
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
				character.AddMana(sim, 7260, manaMetric)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    trinketSpell,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeMana,
		})
	})

	// Normal
	registerApparatusOfKhazGoroth(apparatusConfig{
		ItemID:        68972,
		BonusPerStack: 508,
	})

	// Heroic
	registerApparatusOfKhazGoroth(apparatusConfig{
		ItemID:        69113,
		BonusPerStack: 575,
		Heroic:        true,
	})

	core.NewItemEffect(68994, func(agent core.Agent) {
		character := agent.GetCharacter()

		bonusStats := 1624.0
		procAuraCrit := character.NewTemporaryStatsAura("Matrix Restabilizer Crit Proc", core.ActionID{SpellID: 96978}, stats.Stats{stats.CritRating: bonusStats}, time.Second*30)
		procAuraHaste := character.NewTemporaryStatsAura("Matrix Restabilizer Haste Proc", core.ActionID{SpellID: 96977}, stats.Stats{stats.HasteRating: bonusStats}, time.Second*30)
		procAuraMastery := character.NewTemporaryStatsAura("Matrix Restabilizer Mastery Proc", core.ActionID{SpellID: 96979}, stats.Stats{stats.MasteryRating: bonusStats}, time.Second*30)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 105,
		}

		procAuraCrit.Icd = &icd
		procAuraHaste.Icd = &icd
		procAuraMastery.Icd = &icd

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Matrix Restabilizer Trigger",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			ProcChance: 0.2,
			ActionID:   core.ActionID{ItemID: 68994},
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

	core.NewItemEffect(69150, func(agent core.Agent) {
		character := agent.GetCharacter()

		bonusStats := 1834.0
		procAuraCrit := character.NewTemporaryStatsAura("Matrix Restabilizer Crit Proc (Heroic)", core.ActionID{SpellID: 97140}, stats.Stats{stats.CritRating: bonusStats}, time.Second*30)
		procAuraHaste := character.NewTemporaryStatsAura("Matrix Restabilizer Haste Proc (Heroic)", core.ActionID{SpellID: 97139}, stats.Stats{stats.HasteRating: bonusStats}, time.Second*30)
		procAuraMastery := character.NewTemporaryStatsAura("Matrix Restabilizer Mastery Proc (Heroic)", core.ActionID{SpellID: 97141}, stats.Stats{stats.MasteryRating: bonusStats}, time.Second*30)

		icd := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 105,
		}

		procAuraCrit.Icd = &icd
		procAuraHaste.Icd = &icd
		procAuraMastery.Icd = &icd

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Matrix Restabilizer Trigger (Heroic)",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			ProcChance: 0.2,
			ActionID:   core.ActionID{ItemID: 69150},
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
}

type apparatusConfig struct {
	ItemID        int32
	BonusPerStack float64
	Heroic        bool
}

func registerApparatusOfKhazGoroth(config apparatusConfig) {
	core.NewItemEffect(config.ItemID, func(agent core.Agent) {
		character := agent.GetCharacter()

		labelSuffix := core.Ternary(config.Heroic, " (Heroic)", "")
		buffDuration := time.Second * 15

		buffAuraCrit := character.NewTemporaryStatBuffWithStacks(
			"Blessing of the Shaper Crit"+labelSuffix,
			core.ActionID{SpellID: 96928},
			stats.Stats{stats.CritRating: config.BonusPerStack},
			5,
			buffDuration)

		buffAuraHaste := character.NewTemporaryStatBuffWithStacks(
			"Blessing of the Shaper Haste"+labelSuffix,
			core.ActionID{SpellID: 96927},
			stats.Stats{stats.HasteRating: config.BonusPerStack},
			5,
			buffDuration)

		buffAuraMastery := character.NewTemporaryStatBuffWithStacks(
			"Blessing of the Shaper Mastery"+labelSuffix,
			core.ActionID{SpellID: 96929},
			stats.Stats{stats.MasteryRating: config.BonusPerStack},
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
			ActionID:    core.ActionID{ItemID: config.ItemID},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,
			Cast: core.CastConfig{
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 15,
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

	procTrigger := config.ProcTrigger
	procTrigger.Handler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		target := result.Target
		dot := igniteSpell.Dot(target)
		outstandingDamage := dot.OutstandingDmg()
		newDamage := config.DamageCalculator(result)
		totalDamage := outstandingDamage + newDamage

		if config.IncludeAuraDelay {
			waitTime := time.Millisecond * time.Duration(sim.Roll(375, 625))
			applyDotAt := sim.CurrentTime + waitTime

			if sim.Log != nil {
				unit.Log(sim, "Schedule travel (%0.2f s) for %s", waitTime.Seconds(), config.DotAuraLabel)

				if dot.IsActive() && (dot.NextTickAt() < applyDotAt) {
					unit.Log(sim, "%s rolled with %0.3f damage both ticking and rolled into next", config.DotAuraLabel, outstandingDamage)
				}
			}

			core.StartDelayedAction(sim, core.DelayedActionOptions{
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
