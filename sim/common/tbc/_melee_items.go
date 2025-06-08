package tbc

// func init() {
// 	core.AddEffectsToTest = false

// 	// Proc effects. Keep these in order by item ID.

// 	core.NewItemEffect(9449, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		// Assumes that the user will swap pummelers to have the buff for the whole fight.
// 		character.AddStat(stats.HasteRating, 500)
// 	})

// 	core.NewItemEffect(19019, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(19019, 6.0, 0)

// 		procActionID := core.ActionID{SpellID: 21992}

// 		singleTargetSpell := character.RegisterSpell(core.SpellConfig{
// 			ActionID:    procActionID.WithTag(1),
// 			SpellSchool: core.SpellSchoolNature,
// 			ProcMask:    core.ProcMaskEmpty,

// 			DamageMultiplier: 1,
// 			CritMultiplier:   character.DefaultSpellCritMultiplier(),
// 			ThreatMultiplier: 0.5,

// 			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 				spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)
// 			},
// 		})

// 		makeDebuffAura := func(target *core.Unit) *core.Aura {
// 			return target.GetOrRegisterAura(core.Aura{
// 				Label:    "Thunderfury",
// 				ActionID: procActionID,
// 				Duration: time.Second * 12,
// 				OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 					target.AddStatDynamic(sim, stats.NatureResistance, -25)
// 				},
// 				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 					target.AddStatDynamic(sim, stats.NatureResistance, 25)
// 				},
// 			})
// 		}

// 		numHits := min(5, character.Env.GetNumTargets())
// 		debuffAuras := make([]*core.Aura, len(character.Env.Encounter.TargetUnits))
// 		for i, target := range character.Env.Encounter.TargetUnits {
// 			debuffAuras[i] = makeDebuffAura(target)
// 		}

// 		bounceSpell := character.RegisterSpell(core.SpellConfig{
// 			ActionID:    procActionID.WithTag(2),
// 			SpellSchool: core.SpellSchoolNature,
// 			ProcMask:    core.ProcMaskEmpty,

// 			ThreatMultiplier: 1,
// 			FlatThreatBonus:  63,

// 			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 				curTarget := target
// 				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
// 					result := spell.CalcDamage(sim, curTarget, 0, spell.OutcomeMagicHit)
// 					if result.Landed() {
// 						debuffAuras[target.Index].Activate(sim)
// 					}
// 					spell.DealDamage(sim, result)
// 					curTarget = sim.Environment.NextTargetUnit(curTarget)
// 				}
// 			},
// 		})

// 		character.RegisterAura(core.Aura{
// 			Label:    "Thunderfury",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() {
// 					return
// 				}

// 				if dpm.Proc(sim, spell.ProcMask, "Thunderfury") {
// 					singleTargetSpell.Cast(sim, result.Target)
// 					bounceSpell.Cast(sim, result.Target)
// 				}
// 			},
// 		})
// 	})

// 	core.NewItemEffect(24114, func(agent core.Agent, _ proto.ItemLevelState) {
// 		agent.GetCharacter().PseudoStats.BonusDamage += 5
// 	})

// 	core.NewItemEffect(29297, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		const procChance = 0.03
// 		procAura := character.NewTemporaryStatsAura("Band of the Eternal Defender Proc", core.ActionID{ItemID: 29297}, stats.Stats{stats.Armor: 800}, time.Second*10)

// 		icd := core.Cooldown{
// 			Timer:    character.NewTimer(),
// 			Duration: time.Second * 60,
// 		}

// 		character.GetOrRegisterAura(core.Aura{
// 			Label:    "Band of the Eternal Defender",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() || !spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
// 					return
// 				}
// 				if !icd.IsReady(sim) {
// 					return
// 				}
// 				if sim.RandomFloat("Band of the Eternal Defender") < procChance {
// 					icd.Use(sim)
// 					procAura.Activate(sim)
// 				}
// 			},
// 		})
// 	})

// 	core.NewItemEffect(29301, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		procAura := character.NewTemporaryStatsAura("Band of the Eternal Champion Proc", core.ActionID{ItemID: 29301}, stats.Stats{stats.AttackPower: 160, stats.RangedAttackPower: 160}, time.Second*10)
// 		dpm := character.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMeleeOrRanged)

// 		icd := core.Cooldown{
// 			Timer:    character.NewTimer(),
// 			Duration: time.Second * 60,
// 		}

// 		character.GetOrRegisterAura(core.Aura{
// 			Label:    "Band of the Eternal Champion",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() {
// 					return
// 				}

// 				if !icd.IsReady(sim) {
// 					return
// 				}

// 				if dpm.Proc(sim, spell.ProcMask, "Band of the Eternal Champion") {
// 					icd.Use(sim)
// 					procAura.Activate(sim)
// 				}
// 			},
// 		})
// 	})

// 	core.NewItemEffect(29996, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(29996, 1.0, 0)

// 		actionID := core.ActionID{ItemID: 29996}

// 		var resourceMetricsRage *core.ResourceMetrics
// 		var resourceMetricsEnergy *core.ResourceMetrics
// 		if character.HasRageBar() {
// 			resourceMetricsRage = character.NewRageMetrics(actionID)
// 		}
// 		if character.HasEnergyBar() {
// 			resourceMetricsEnergy = character.NewEnergyMetrics(actionID)
// 		}

// 		character.GetOrRegisterAura(core.Aura{
// 			Label:    "Rod of the Sun King",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() {
// 					return
// 				}

// 				if dpm.Proc(sim, spell.ProcMask, "Rod of the Sun King") {
// 					switch spell.Unit.GetCurrentPowerBar() {
// 					case core.RageBar:
// 						spell.Unit.AddRage(sim, 5, resourceMetricsRage)
// 					case core.EnergyBar:
// 						spell.Unit.AddEnergy(sim, 10, resourceMetricsEnergy)
// 					}
// 				}
// 			},
// 		})
// 	})

// 	core.NewItemEffect(31193, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(31193, 0, 0.02)

// 		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
// 			ActionID:    core.ActionID{SpellID: 24585},
// 			SpellSchool: core.SpellSchoolShadow,
// 			ProcMask:    core.ProcMaskEmpty,

// 			DamageMultiplier: 1,
// 			CritMultiplier:   character.DefaultSpellCritMultiplier(),
// 			ThreatMultiplier: 1,

// 			BonusCoefficient: 1,

// 			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 				baseDamage := sim.Roll(48, 54)
// 				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
// 			},
// 		})

// 		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
// 			Name:     "Blade of Unquenched Thirst Trigger",
// 			ActionID: core.ActionID{ItemID: 31193},
// 			Callback: core.CallbackOnSpellHitDealt,
// 			Outcome:  core.OutcomeLanded,
// 			DPM:      dpm,
// 			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
// 				procSpell.Cast(sim, result.Target)
// 			},
// 		})
// 	})

// 	core.NewItemEffect(32262, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(12590, 1.0, 0)

// 		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
// 			ActionID:    core.ActionID{SpellID: 40291},
// 			SpellSchool: core.SpellSchoolShadow,
// 			ProcMask:    core.ProcMaskEmpty,

// 			DamageMultiplier: 1,
// 			CritMultiplier:   character.DefaultSpellCritMultiplier(),
// 			ThreatMultiplier: 1,

// 			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 				spell.CalcAndDealDamage(sim, target, 20, spell.OutcomeMagicHitAndCrit)
// 			},
// 		})

// 		procAura := character.GetOrRegisterAura(core.Aura{
// 			Label:    "Siphon Essence",
// 			ActionID: core.ActionID{SpellID: 40291},
// 			Duration: time.Second * 6,
// 			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
// 					return
// 				}

// 				procSpell.Cast(sim, result.Target)
// 			},
// 		})

// 		character.GetOrRegisterAura(core.Aura{
// 			Label:    "Syphon of the Nathrezim",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() {
// 					return
// 				}

// 				if dpm.Proc(sim, spell.ProcMask, "Syphon Of The Nathrezim") {
// 					procAura.Activate(sim)
// 				}
// 			},
// 		})
// 	})

// 	core.NewItemEffect(32375, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		const procChance = 0.02
// 		procAura := character.NewTemporaryStatsAura("Bulwark Of Azzinoth Proc", core.ActionID{ItemID: 32375}, stats.Stats{stats.Armor: 2000}, time.Second*10)

// 		character.GetOrRegisterAura(core.Aura{
// 			Label:    "Bulwark Of Azzinoth",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if result.Landed() && spell.SpellSchool.Matches(core.SpellSchoolPhysical) && sim.RandomFloat("Bulwark of Azzinoth") < procChance {
// 					procAura.Activate(sim)
// 				}
// 			},
// 		})
// 	})

// 	core.NewItemEffect(34473, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		procAura := character.NewTemporaryStatsAura("Commendation of Kael'Thas Proc", core.ActionID{ItemID: 34473}, stats.Stats{stats.DodgeRating: 152}, time.Second*10)

// 		icd := core.Cooldown{
// 			Timer:    character.NewTimer(),
// 			Duration: time.Second * 30,
// 		}

// 		character.GetOrRegisterAura(core.Aura{
// 			Label:    "Commendation of Kael'Thas",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
// 					return
// 				}

// 				if aura.Unit.CurrentHealthPercent() >= 0.35 {
// 					return
// 				}

// 				if !icd.IsReady(sim) {
// 					return
// 				}

// 				icd.Use(sim)
// 				procAura.Activate(sim)
// 			},
// 		})
// 	})

// 	core.NewItemEffect(12590, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()

// 		dpm := character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(12590, 1.0, 0)

// 		effectAura := character.NewTemporaryStatsAura("Felstriker Proc", core.ActionID{SpellID: 16551}, stats.Stats{stats.PhysicalCritPercent: 100}, time.Second*3)

// 		character.GetOrRegisterAura(core.Aura{
// 			Label:    "Felstriker",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 				if !result.Landed() {
// 					return
// 				}

// 				if dpm.Proc(sim, spell.ProcMask, "Felstriker") {
// 					effectAura.Activate(sim)
// 				}
// 			},
// 		})
// 	})

// 	core.AddEffectsToTest = true
// }
