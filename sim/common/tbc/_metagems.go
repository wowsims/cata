package tbc

// func init() {
// 	core.AddEffectsToTest = false
// 	// Keep these in order by item ID.

// 	core.NewItemEffect(25893, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()
// 		procAura := character.NewTemporaryStatsAura("Mystic Focus Proc", core.ActionID{ItemID: 25893}, stats.Stats{stats.HasteRating: 320}, time.Second*4)

// 		icd := core.Cooldown{
// 			Timer:    character.NewTimer(),
// 			Duration: time.Second * 35,
// 		}

// 		character.RegisterAura(core.Aura{
// 			Label:    "Mystical Skyfire Diamond",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 				if !icd.IsReady(sim) || sim.RandomFloat("Mystical Skyfire Diamond") > 0.15 {
// 					return
// 				}
// 				icd.Use(sim)
// 				procAura.Activate(sim)
// 			},
// 		})
// 	})

// 	core.NewItemEffect(25899, func(agent core.Agent, _ proto.ItemLevelState) {
// 		agent.GetCharacter().PseudoStats.BonusDamage += 3
// 	})

// 	core.NewItemEffect(25901, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()
// 		icd := core.Cooldown{
// 			Timer:    character.NewTimer(),
// 			Duration: time.Second * 15,
// 		}
// 		manaMetrics := character.NewManaMetrics(core.ActionID{ItemID: 25901})

// 		character.RegisterAura(core.Aura{
// 			Label:    "Insightful Earthstorm Diamond",
// 			Duration: core.NeverExpires,
// 			OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 				aura.Activate(sim)
// 			},
// 			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 				if !icd.IsReady(sim) || sim.RandomFloat("Insightful Earthstorm Diamond") > 0.04 {
// 					return
// 				}
// 				icd.Use(sim)
// 				character.AddMana(sim, 300, manaMetrics)
// 			},
// 		})
// 	})

// 	core.NewItemEffect(32410, func(agent core.Agent, _ proto.ItemLevelState) {
// 		character := agent.GetCharacter()
// 		procAura := character.NewTemporaryStatsAura("Thundering Skyfire Diamond Proc", core.ActionID{ItemID: 32410}, stats.Stats{stats.HasteRating: 240}, time.Second*6)

// 		icd := core.Cooldown{
// 			Timer:    character.NewTimer(),
// 			Duration: time.Second * 40,
// 		}
// 		dpm := character.AutoAttacks.NewPPMManager(1.5, core.ProcMaskWhiteHit) // Mask 68, melee or ranged auto attacks.

// 		character.RegisterAura(core.Aura{
// 			Label:    "Thundering Skyfire Diamond",
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

// 				if dpm.Proc(sim, spell.ProcMask, "Thundering Skyfire Diamond") {
// 					icd.Use(sim)
// 					procAura.Activate(sim)
// 				}
// 			},
// 		})
// 	})

// 	// Eternal Earthstorm
// 	core.NewItemEffect(35501, func(agent core.Agent, _ proto.ItemLevelState) {
// 		agent.GetCharacter().PseudoStats.BlockDamageReduction += 0.01
// 	})

// 	core.NewItemEffect(35503, func(agent core.Agent, _ proto.ItemLevelState) {
// 		agent.GetCharacter().MultiplyStat(stats.Intellect, 1.02)
// 	})

// 	// These are handled in character.go, but create empty effects, so they are included in tests.
// 	core.NewItemEffect(34220, func(_ core.Agent, _ proto.ItemLevelState) {}) // Chaotic Skyfire Diamond
// 	core.NewItemEffect(32409, func(_ core.Agent, _ proto.ItemLevelState) {}) // Relentless Earthstorm Diamond

// 	core.AddEffectsToTest = true
// }
