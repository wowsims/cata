package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// T11
var ItemSetMaleficRaiment = core.NewItemSet(core.ItemSet{
	Name: "Shadowflame Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_CastTime_Pct,
				ClassMask:  WarlockSpellChaosBolt | WarlockSpellHandOfGuldan | WarlockSpellHaunt,
				FloatValue: -0.1,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()

			dmgMod := warlock.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  WarlockSpellFelFlame,
				FloatValue: 3.0,
			})

			aura := warlock.RegisterAura(core.Aura{
				Label:     "Fel Spark",
				ActionID:  core.ActionID{SpellID: 89937},
				Duration:  15 * time.Second,
				MaxStacks: 2,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Deactivate()
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Matches(WarlockSpellFelFlame) && result.Landed() {
						aura.RemoveStack(sim)
					}
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Item - Warlock T11 4P Bonus",
				ActionID:       core.ActionID{SpellID: 89935},
				ClassSpellMask: WarlockSpellImmolateDot | WarlockSpellUnstableAffliction,
				Callback:       core.CallbackOnPeriodicDamageDealt,
				ProcChance:     0.02,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
					aura.SetStacks(sim, 2)
				},
			})
		},
	},
})

// T12
type FieryImpPet struct {
	core.Pet

	FlameBlast *core.Spell
}

func (warlock *Warlock) NewFieryImp() *FieryImpPet {
	baseStats := stats.Stats{stats.SpellCritPercent: 0} // rough guess, seems to only get crit from debuffs?

	statInheritance := func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{}
	}

	imp := &FieryImpPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                     "Fiery Imp",
			Owner:                    &warlock.Character,
			BaseStats:                baseStats,
			NonHitExpStatInheritance: statInheritance,
			EnabledOnStart:           false,
			IsGuardian:               true,
		}),
	}

	warlock.AddPet(imp)
	imp.registerFlameBlast(warlock)

	return imp
}

func (imp *FieryImpPet) GetPet() *core.Pet { return &imp.Pet }

func (imp *FieryImpPet) Initialize() {}

func (imp *FieryImpPet) Reset(_ *core.Simulation) {}

func (imp *FieryImpPet) ExecuteCustomRotation(sim *core.Simulation) {
	if imp.FlameBlast.CanCast(sim, imp.CurrentTarget) {
		imp.FlameBlast.Cast(sim, imp.CurrentTarget)
		delay := time.Duration(sim.RollWithLabel(150.0, 750.0, "Imp Cast Delay")) * time.Millisecond
		imp.WaitUntil(sim, imp.NextGCDAt()+delay)
		return
	}
}

func (pet *FieryImpPet) registerFlameBlast(warlock *Warlock) {
	pet.FlameBlast = pet.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 99226},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 16,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           pet.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := sim.Roll(2750+1, 2750+1514)
			result := spell.CalcDamage(sim, target, baseDmg, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

var ItemSetBalespidersBurningVestments = core.NewItemSet(core.ItemSet{
	Name: "Balespider's Burning Vestments",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Item - Warlock T12 2P Bonus",
				ActionID:   core.ActionID{SpellID: 99220},
				ProcChance: 0.05,
				ICD:        45 * time.Second,
				Callback:   core.CallbackOnPeriodicDamageDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					warlock.FieryImp.EnableWithTimeout(sim, warlock.FieryImp, 15*time.Second)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()

			dmgMod := warlock.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				School:     core.SpellSchoolShadow | core.SpellSchoolFire,
				FloatValue: 0.20,
			})

			aura := warlock.RegisterAura(core.Aura{
				Label:    "Apocalypse",
				ActionID: core.ActionID{SpellID: 99232},
				Duration: 8 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					dmgMod.Deactivate()
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Item - Warlock T12 4P Bonus",
				ActionID:       core.ActionID{SpellID: 99229},
				ClassSpellMask: WarlockSpellShadowBolt | WarlockSpellIncinerate | WarlockSpellSoulFire | WarlockSpellDrainSoul,
				ProcChance:     0.05,
				Callback:       core.CallbackOnCastComplete,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetGladiatorsFelshroud = core.NewItemSet(core.ItemSet{
	ID:   910,
	Name: "Gladiator's Felshroud",

	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Intellect, 70)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Intellect, 90)
			// TODO: enable if we ever implement death coil
			// lock.AddStaticMod(core.SpellModConfig{
			// 	Kind:       core.SpellMod_Cooldown_Flat,
			// 	ClassMask:  WarlockSpellDeathCoil,
			// 	FloatValue: -30 * time.Second,
			// })
		},
	},
})

// T13
var ItemSetVestmentsOfTheFacelessShroud = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Faceless Shroud",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				TimeValue: -time.Minute * 4,
				ClassMask: WarlockSpellSummonDoomguard | WarlockSpellSummonInfernal,
			})

			summonDuration := core.TernaryDuration(warlock.Spec == proto.Spec_SpecDemonologyWarlock, 20*time.Second, 30*time.Second)

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_BuffDuration_Flat,
				ClassMask: WarlockSpellSummonDoomguard | WarlockSpellSummonInfernal,
				TimeValue: summonDuration,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()

			spDep := warlock.NewDynamicMultiplyStat(stats.SpellPower, 1.1)
			aura := warlock.RegisterAura(core.Aura{
				Label:    "Temporal Ruin",
				ActionID: core.ActionID{SpellID: 105786},
				Duration: 10 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.EnableDynamicStatDep(sim, spDep)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.DisableDynamicStatDep(sim, spDep)
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Item - Warlock T13 4P Bonus (Soulburn)",
				ActionID:       core.ActionID{SpellID: 105787},
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: WarlockSpellSoulBurn,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
				},
			})

			setBonusAura.ExposeToAPL(105787)

			warlock.T13_4pc = setBonusAura
		},
	},
})

// T14
var ItemSetShaSkinRegalia = core.NewItemSet(core.ItemSet{
	Name: "Sha-Skin Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.1,
				ClassMask:  WarlockSpellCorruption,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.05,
				ClassMask:  WarlockSpellIncinerate | WarlockSpellFaBIncinerate,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.02,
				ClassMask:  WarlockSpellShadowBolt | WarlockSpellDemonicSlash | WarlockSpellTouchOfChaos,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			buff := agent.GetCharacter().RegisterAura(core.Aura{
				Label:    "Sha-Skin Regalia - 4P Buff",
				ActionID: core.ActionID{SpellID: 148463},
				Duration: time.Second * 20,
			}).AttachMultiplicativePseudoStatBuff(&agent.GetCharacter().PseudoStats.DamageDealtMultiplier, 1.10)

			agent.GetCharacter().OnSpellRegistered(func(spell *core.Spell) {
				if spell.Matches(WarlockDarkSoulSpell) {
					spell.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
						buff.Activate(sim)
					})
				}
			})
		},
	},
})

// T15
var ItemSetRegaliaOfTheThousandfeldHells = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Thousandfold Hells",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.T15_2pc = agent.GetCharacter().RegisterAura(core.Aura{
				Label:    "Regalia of the Thousandfold Hells - 2P Buff",
				Duration: time.Second * 20,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotNumberOfTicks_Flat,
				IntValue:  2,
				ClassMask: WarlockSpellHaunt,
			})

			agent.GetCharacter().OnSpellRegistered(func(spell *core.Spell) {
				if spell.Matches(WarlockDarkSoulSpell) {
					spell.RelatedSelfBuff.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
						warlock.T15_2pc.Activate(sim)
					})
				}
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.05,
				ClassMask:  WarlockSpellMaleficGrasp | WarlockSpellDrainSoul,
			})

			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.T15_4pc = setBonusAura
		},
	},
})

// T16
var ItemSetRegaliaOfTheHornedNightmare = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Horned Nightmare",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()
			var buff *core.Aura
			switch warlock.Spec {
			case proto.Spec_SpecAfflictionWarlock:
				buff = warlock.RegisterAura(core.Aura{
					ActionID: core.ActionID{SpellID: 145082},
					Label:    "Regalia of the Horned Nightmare - Affli - 2pc",
					Duration: time.Second * 10,
				}).AttachSpellMod(core.SpellModConfig{
					Kind:       core.SpellMod_DamageDone_Pct,
					FloatValue: 0.15,
					ClassMask:  WarlockSpellDrainSoul | WarlockSpellMaleficGrasp,
				})

				warlock.T16_2pc_buff = buff
				setBonusAura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Matches(WarlockSpellUnstableAffliction) && result.DidCrit() && sim.Proc(0.5, "T16 - 2pc") {
						buff.Activate(sim)
						return
					}
				}
			case proto.Spec_SpecDemonologyWarlock:
				// TODO: Research if all pets or just the primary pet is affected
				buffAction := core.ActionID{SpellID: 145075}
				applyBuffAura := func(unit *core.Unit) {
					unit.RegisterAura(core.Aura{
						ActionID: buffAction,
						Label:    "Regalia of the Horned Nightmare - Demo - 2pc",
						Duration: time.Second * 10,
					}).AttachMultiplicativePseudoStatBuff(&unit.PseudoStats.DamageDealtMultiplier, 1.2)
				}

				applyBuffAura(&warlock.Unit)
				for _, pet := range warlock.Pets {
					if pet.IsGuardian() {
						continue
					}

					applyBuffAura(&pet.Unit)
				}

				setBonusAura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Matches(WarlockSpellSoulFire) && sim.Proc(0.2, "T16 - 2pc") {
						warlock.GetAuraByID(buffAction).Activate(sim)
						for _, pet := range warlock.Pets {
							if pet.IsGuardian() {
								continue
							}

							if !pet.IsActive() {
								continue
							}

							pet.GetAuraByID(buffAction).Activate(sim)
						}
					}
				}
			case proto.Spec_SpecDestructionWarlock:
				buff = warlock.RegisterAura(core.Aura{
					ActionID: core.ActionID{SpellID: 145075},
					Label:    "Regalia of the Horned Nightmare - Destro - 2pc",
					Duration: time.Second * 10,
				}).AttachSpellMod(core.SpellModConfig{
					Kind:       core.SpellMod_BonusCrit_Percent,
					FloatValue: 0.1,
					ClassMask:  WarlockSpellImmolate | WarlockSpellImmolateDot | WarlockSpellIncinerate | WarlockSpellFaBIncinerate,
				})

				setBonusAura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.Matches(WarlockSpellConflagrate|WarlockSpellFaBConflagrate) && result.DidCrit() && sim.Proc(0.2, "T16 - 2pc") {
						buff.Activate(sim)
						return
					}
				}
			default:
				return
			}
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			warlock := agent.(WarlockAgent).GetWarlock()
			switch agent.GetCharacter().Spec {
			case proto.Spec_SpecAfflictionWarlock:
				warlock.OnSpellRegistered(func(spell *core.Spell) {
					if !spell.Matches(WarlockSpellHaunt) {
						return
					}

					for _, enemy := range warlock.Env.Encounter.TargetUnits {
						for _, array := range spell.RelatedAuraArrays {
							array.Get(enemy).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
								if sim.Proc(0.1, "T16 4p") {
									warlock.GetSecondaryResourceBar().Gain(sim, 1, spell.ActionID)
								}
							})
						}
					}
				})
			case proto.Spec_SpecDemonologyWarlock:
				setBonusAura.AttachProcTrigger(core.ProcTrigger{
					Callback:       core.CallbackOnCastComplete,
					ClassSpellMask: WarlockSpellShadowBolt | WarlockSpellTouchOfChaos,
					ProcChance:     0.08,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						// TODO: NOT IMPLEMENTED - Need to verify how this interacts with existing Shadow Flame DoT
					},
				})
			case proto.Spec_SpecDestructionWarlock:
				buff := warlock.RegisterAura(core.Aura{
					ActionID: core.ActionID{SpellID: 145164},
					Label:    "Regalia of the Horned Nightmare - Demo - 4pc",
					Duration: time.Second * 5,
					Icd: &core.Cooldown{
						Timer:    warlock.NewTimer(),
						Duration: time.Second * 10,
					},
				}).AttachStatBuff(stats.CritRating, core.CritRatingPerCritPercent*15)

				warlock.GetSecondaryResourceBar().RegisterOnGain(func(
					sim *core.Simulation,
					gain, realGain int32,
					actionID core.ActionID,
				) {
					if realGain == 0 || buff.Icd.IsReady(sim) {
						return
					}

					old := warlock.GetSecondaryResourceBar().Value() - realGain
					if int(old/10) == int(warlock.GetSecondaryResourceBar().Value()/10) {
						return
					}

					buff.Activate(sim)
				})
			}
		},
	},
})
