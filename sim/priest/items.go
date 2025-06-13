package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

// var ItemSetVestmentsOfAbsolution = core.NewItemSet(core.ItemSet{
// 	Name: "Vestments of Absolution",
// 	Bonuses: map[int32]core.ApplySetBonus{
// 		2: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_PowerCost_Pct,
// 				IntValue: -100,
// 				ClassMask:  PriestSpellPrayerOfHealing,
// 			})
// 		},
// 		4: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_DamageDone_Flat,
// 				FloatValue: 0.05,
// 				ClassMask:  PriestSpellGreaterHeal,
// 			})
// 		},
// 	},
// })

// var ItemSetValorous = core.NewItemSet(core.ItemSet{
// 	Name: "Garb of Faith",
// 	Bonuses: map[int32]core.ApplySetBonus{
// 		2: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_PowerCost_Pct,
// 				IntValue: -10,
// 				ClassMask:  PriestSpellMindBlast,
// 			})
// 		},
// 		4: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_BonusCrit_Rating,
// 				FloatValue: 10 * core.CritRatingPerCritChance,
// 				ClassMask:  PriestSpellShadowWordDeath,
// 			})
// 		},
// 	},
// })

// var ItemSetRegaliaOfFaith = core.NewItemSet(core.ItemSet{
// 	Name: "Regalia of Faith",
// 	Bonuses: map[int32]core.ApplySetBonus{
// 		2: func(agent core.Agent, setBonusAura *core.Aura) {
// 			// Not implemented
// 		},
// 		4: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_PowerCost_Pct,
// 				IntValue: -5,
// 				ClassMask:  PriestSpellGreaterHeal,
// 			})
// 		},
// 	},
// })

// var ItemSetConquerorSanct = core.NewItemSet(core.ItemSet{
// 	Name: "Sanctification Garb",
// 	Bonuses: map[int32]core.ApplySetBonus{
// 		2: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_DamageDone_Flat,
// 				FloatValue: 0.15,
// 				ClassMask:  PriestSpellDevouringPlague,
// 			})
// 		},
// 		4: func(agent core.Agent, setBonusAura *core.Aura) {
// 			priest := agent.(PriestAgent).GetPriest()
// 			procAura := priest.NewTemporaryStatsAura("Devious Mind", core.ActionID{SpellID: 64907}, stats.Stats{stats.SpellHaste: 240}, time.Second*4)

// 			priest.RegisterAura(core.Aura{
// 				Label:    "Devious Mind Proc",
// 				Duration: core.NeverExpires,
// 				OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 					aura.Activate(sim)
// 				},
// 				// TODO: Does this affect the spell that procs it?
// 				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 					if spell.ClassSpellMask == PriestSpellMindBlast {
// 						procAura.Activate(sim)
// 					}
// 				},
// 			})
// 		},
// 	},
// })

// var ItemSetSanctificationRegalia = core.NewItemSet(core.ItemSet{
// 	Name: "Sanctification Regalia",
// 	Bonuses: map[int32]core.ApplySetBonus{
// 		2: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_BonusCrit_Rating,
// 				FloatValue: 10 * core.CritRatingPerCritChance,
// 				ClassMask:  PriestSpellPrayerOfHealing,
// 			})
// 		},
// 		4: func(agent core.Agent, setBonusAura *core.Aura) {
// 			priest := agent.(PriestAgent).GetPriest()
// 			procAura := priest.NewTemporaryStatsAura("Sanctification Reglia 4pc", core.ActionID{SpellID: 64912}, stats.Stats{stats.SpellPower: 250}, time.Second*5)

// 			priest.RegisterAura(core.Aura{
// 				Label:    "Sancitifcation Reglia 4pc",
// 				Duration: core.NeverExpires,
// 				OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 					aura.Activate(sim)
// 				},
// 				// TODO: Does this affect the spell that procs it?
// 				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 					if spell == priest.PowerWordShield {
// 						procAura.Activate(sim)
// 					}
// 				},
// 			})
// 		},
// 	},
// })

// var ItemSetZabras = core.NewItemSet(core.ItemSet{
// 	Name:            "Zabra's Regalia",
// 	AlternativeName: "Velen's Regalia",
// 	Bonuses: map[int32]core.ApplySetBonus{
// 		2: func(agent core.Agent, setBonusAura *core.Aura) {
// 			// Modifies dot length, need to implement later again
// 			// Requieres tests and proper modification of SpellMods
// 		},
// 		4: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_BonusCrit_Rating,
// 				FloatValue: 5 * core.CritRatingPerCritChance,
// 				ClassMask:  PriestSpellMindFlay,
// 			})
// 		},
// 	},
// })

// var ItemSetZabrasRaiment = core.NewItemSet(core.ItemSet{
// 	Name:            "Zabra's Raiment",
// 	AlternativeName: "Velen's Raiment",
// 	Bonuses: map[int32]core.ApplySetBonus{
// 		2: func(agent core.Agent, setBonusAura *core.Aura) {
// 			character := agent.GetCharacter()
// 			character.AddStaticMod(core.SpellModConfig{
// 				Kind:       core.SpellMod_DamageDone_Flat,
// 				FloatValue: 0.15,
// 				ClassMask:  PriestSpellPrayerOfMending,
// 			})
// 		},
// 		4: func(agent core.Agent, setBonusAura *core.Aura) {
// 			// changed in cata to flat 5% heal
// 			character := agent.GetCharacter()
// 			character.PseudoStats.DamageDealtMultiplier *= 1.05
// 		},
// 	},
// })

var ItemSetCrimsonAcolyte = core.NewItemSet(core.ItemSet{
	Name: "Crimson Acolyte's Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 5,
				ClassMask:  PriestSpellShadowWordPain | PriestSpellDevouringPlague | PriestSpellVampiricTouch | PriestSpellImprovedDevouringPlague,
			})
		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotTickLength_Flat,
				TimeValue: -time.Millisecond * 170,
				ClassMask: PriestSpellMindFlay,
			})
		},
	},
})

var ItemSetCrimsonAcolytesRaiment = core.NewItemSet(core.ItemSet{
	Name: "Crimson Acolyte's Raiment",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			priest := agent.(PriestAgent).GetPriest()

			var curAmount float64
			procSpell := priest.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 70770},
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Hot: core.DotConfig{
					Aura: core.Aura{
						Label: "Crimson Acolytes Raiment 2pc - Hot",
					},
					NumberOfTicks: 3,
					TickLength:    time.Second * 3,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
						dot.SnapshotBaseDamage = curAmount * 0.33
						dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
					},
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Crimson Acolytes Raiment 2pc",
				ClassSpellMask: PriestSpellFlashHeal,
				ProcChance:     1.0 / 3,
				Callback:       core.CallbackOnHealDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					curAmount = result.Damage
					hot := procSpell.Hot(result.Target)
					hot.Apply(sim)
				},
			})

			setBonusAura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
				curAmount = 0
			})

		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.05,
				ClassMask:  PriestSpellPowerWordShield,
			})

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.10,
				ClassMask:  PriestSpellCircleOfHealing,
			})
		},
	},
})

// T11 - Shadow
var ItemSetMercurialRegalia = core.NewItemSet(core.ItemSet{
	Name: "Mercurial Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 5,
				ClassMask:  PriestSpellMindFlay | PriestSpellMindSear,
			})
		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.3,
				ClassMask:  PriestSpellShadowyApparation,
			})
		},
	},
})

// T12 - Shadow
var ItemSetRegaliaOfTheCleansingFlame = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Cleansing Flame",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Fiend deals 20% extra damage as fire damage and cooldown reduced by 75 seconds
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				TimeValue: -time.Second * 75,
				ClassMask: PriestSpellShadowFiend,
			})

			priest := agent.(PriestAgent).GetPriest()
			shadowFlameProc := priest.ShadowfiendPet.RegisterSpell(core.SpellConfig{
				ActionID:         core.ActionID{SpellID: 99156},
				SpellSchool:      core.SpellSchoolFire,
				ProcMask:         core.ProcMaskEmpty,
				Flags:            core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods,
				BaseCost:         0,
				CritMultiplier:   1.0,
				DamageMultiplier: 1.0,
				ThreatMultiplier: 1.0,
				ManaCost: core.ManaCostOptions{
					BaseCostPercent: 0,
					PercentModifier: 100,
				},

				Cast: core.CastConfig{
					DefaultCast: core.Cast{
						GCD:      0,
						CastTime: 0,
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, spell.BonusSpellPower, spell.OutcomeAlwaysHit)
				},
			})

			setBonusAura.MakeDependentProcTriggerAura(&priest.ShadowfiendPet.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 99155},
				Name:       "Shadowflame (T12-2P)",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				Outcome:    core.OutcomeLanded,
				ProcChance: 1.0,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					// deal 20% of melee damage on hit
					shadowFlameProc.BonusSpellPower = result.Damage * 0.2
					shadowFlameProc.Cast(sim, result.Target)
				},
			})

		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()
			mbMod := character.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.15,
				ClassMask:  PriestSpellMindBlast,
			})

			mbAura := character.RegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 99158},
				Label:    "Dark Flames",
				Duration: core.NeverExpires,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mbMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mbMod.Deactivate()
				},
			})

			priest := agent.(PriestAgent).GetPriest()
			dotTracker := make([]int, len(priest.Env.AllUnits))
			setHandler := func(aura *core.Aura, sim *core.Simulation) {
				if setBonusAura.IsActive() && aura.IsActive() {
					dotTracker[aura.Unit.UnitIndex]++
				} else {
					dotTracker[aura.Unit.UnitIndex]--
				}

				// check if we have 3 dots anywhere
				for _, count := range dotTracker {
					if count == 3 {
						mbAura.Activate(sim)
						return
					}
				}

				mbAura.Deactivate(sim)
			}

			priest.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(PriestSpellShadowWordPain | PriestSpellVampiricTouch | PriestSpellDevouringPlague) {
					return
				}

				for _, target := range priest.Env.AllUnits {
					if target.Type == core.EnemyUnit {
						spell.Dot(target).ApplyOnGain(setHandler)
						spell.Dot(target).ApplyOnExpire(setHandler)
					}
				}
			})
		},
	},
})

// T13 - Shadow
var ItemSetRegaliaOfDyingLight = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Dying Light",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  PriestSpellShadowWordDeath,
				FloatValue: 0.55,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			priest := agent.(PriestAgent).GetPriest()

			makeProcTriggerConfig := func(config core.ProcTrigger) core.ProcTrigger {
				return core.ProcTrigger{
					ActionID:       core.ActionID{SpellID: 105844},
					Name:           "Item - Priest T13 Shadow 4P Bonus (Shadowfiend and Shadowy Apparition)",
					Callback:       core.CallbackOnSpellHitDealt,
					Outcome:        core.OutcomeLanded,
					ProcChance:     1.0,
					ClassSpellMask: config.ClassSpellMask,
					ProcMask:       config.ProcMask,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if priest.ShadowOrbsAura != nil {
							priest.ShadowOrbsAura.Activate(sim)
							priest.ShadowOrbsAura.SetStacks(sim, 3)
						}
					},
				}
			}

			setBonusAura.MakeDependentProcTriggerAura(&priest.ShadowfiendPet.Unit, makeProcTriggerConfig(core.ProcTrigger{
				ProcMask: core.ProcMaskMelee,
			}))

			setBonusAura.AttachProcTrigger(makeProcTriggerConfig(core.ProcTrigger{
				ClassSpellMask: PriestSpellShadowyApparation,
			}))
		},
	},
})

// T14 - Shadow
var ItemSetRegaliaOfTheGuardianSperpent = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Guardian Serpent",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				ClassMask:  PriestSpellShadowWordPain,
				FloatValue: 10,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotNumberOfTicks_Flat,
				ClassMask: PriestSpellShadowWordPain | PriestSpellVampiricTouch,
				IntValue:  1,
			})
		},
	},
})

var ItemSetRegaliaOfTheExorcist = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Exorcist",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			priest := agent.(PriestAgent).GetPriest()
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Regalia of the Exorcist - 2P",
				SpellFlags:     core.SpellFlagPassiveSpell,
				ProcChance:     0.65,
				ClassSpellMask: PriestSpellShadowyApparation,
				Outcome:        core.OutcomeLanded,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if priest.ShadowWordPain != nil && priest.ShadowWordPain.Dot(result.Target).IsActive() {
						priest.ShadowWordPain.Dot(result.Target).AddTick()
					}

					if priest.VampiricTouch != nil && priest.VampiricTouch.Dot(result.Target).IsActive() {
						priest.VampiricTouch.Dot(result.Target).AddTick()
					}
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			priest := agent.(PriestAgent).GetPriest()
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Regalia of the Exorcist - 4P",
				ProcMask:       core.ProcMaskSpellDamage,
				ProcChance:     0.1,
				ClassSpellMask: PriestSpellVampiricTouch,
				Outcome:        core.OutcomeLanded,
				Callback:       core.CallbackOnPeriodicDamageDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					priest.ShadowyApparition.Cast(sim, result.Target)
				},
			})
		},
	},
})

var ItemSetRegaliaOfTheTernionGlory = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Ternion Glory",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_CritMultiplier_Flat,
				FloatValue: 0.4,
				ClassMask:  PriestSpellShadowyRecall,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			priest := agent.(PriestAgent).GetPriest()
			mod := priest.Unit.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.2,
				ClassMask:  PriestSpellShadowWordDeath | PriestSpellMindSpike | PriestSpellMindBlast,
			})

			var orbsSpend int32 = 0
			priest.Unit.GetSecondaryResourceBar().RegisterOnSpend(func(_ *core.Simulation, amount int32, _ core.ActionID) {
				orbsSpend = amount
			})

			aura := priest.Unit.RegisterAura(core.Aura{
				Label:    "Regalia of the Ternion Glory - 4P (Proc)",
				ActionID: core.ActionID{SpellID: 145180},
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mod.UpdateFloatValue(0.2 * float64(orbsSpend))
					mod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mod.Deactivate()
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.Matches(PriestSpellMindBlast | PriestSpellMindSpike | PriestSpellShadowWordDeath) {
						return
					}

					aura.Deactivate(sim)
				},
			})

			core.MakeProcTriggerAura(&priest.Unit, core.ProcTrigger{
				Name:           "Regalia of the Ternion Glory - 4P",
				Outcome:        core.OutcomeLanded,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: PriestSpellDevouringPlague,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
				},
			})
		},
	},
})

var shaWeaponIDs = []int32{86990, 86865, 86227}

func init() {
	for _, id := range shaWeaponIDs {
		core.NewItemEffect(id, func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			priest.AddStaticMod(core.SpellModConfig{
				Kind:      core.SpellMod_GlobalCooldown_Flat,
				TimeValue: -core.GCDDefault,
				ClassMask: PriestSpellShadowFiend | PriestSpellMindBender,
			})
		})
	}
}
