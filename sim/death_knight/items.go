package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// T14 DPS
var ItemSetBattlegearOfTheLostCatacomb = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of the Lost Catacomb",
	ID:   1123,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Obliterate, Frost Strike, and Scourge Strike deal 4% increased damage.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec == proto.Spec_SpecBloodDeathKnight {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  DeathKnightSpellFrostStrike | DeathKnightSpellObliterate | DeathKnightSpellScourgeStrike,
				FloatValue: 0.04,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Pillar of Frost ability grants 5% additional Strength, and your Unholy Frenzy ability grants 10% additional haste.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec == proto.Spec_SpecBloodDeathKnight {
				return
			}

			// Handled in sim/core/buffs.go and sim/death_knight/frost/pillar_of_frost.go
			dk.T14Dps4pc = setBonusAura
		},
	},
})

// T14 Tank
var ItemSetPlateOfTheLostCatacomb = core.NewItemSet(core.ItemSet{
	Name: "Plate of the Lost Catacomb",
	ID:   1124,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Reduces the cooldown of your Vampiric Blood ability by 20 sec.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: DeathKnightSpellVampiricBlood,
				TimeValue: time.Second * -20,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the healing received from your Death Strike by 10%.
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			setBonusAura.AttachMultiplicativePseudoStatBuff(
				&dk.deathStrikeHealingMultiplier, 1.1,
			)
		},
	},
})

// T15 DPS
var ItemSetBattleplateOfTheAllConsumingMaw = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of the All-Consuming Maw",
	ID:   1152,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your attacks have a chance to raise the spirit of a fallen Zandalari as your Death Knight minion for 15 sec.
			// (Approximately 1.15 procs per minute)
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			risenZandalariSpell := dk.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 138342},
				SpellSchool: core.SpellSchoolPhysical,
				Flags:       core.SpellFlagPassiveSpell,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					for _, troll := range dk.FallenZandalari {
						if troll.IsActive() {
							continue
						}

						troll.EnableWithTimeout(sim, troll, time.Second*15)

						return
					}

					if sim.Log != nil {
						dk.Log(sim, "No Fallen Zandalari available for the T15 4pc to proc, this is unreasonable.")
					}
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback: core.CallbackOnSpellHitDealt,
				Outcome:  core.OutcomeLanded,
				DPM: dk.NewSetBonusRPPMProcManager(138343, setBonusAura, core.ProcMaskDirect, core.RPPMConfig{
					PPM: 1.14999997616,
				}),
				ICD: time.Millisecond * 250,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					risenZandalariSpell.Cast(sim, result.Target)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Soul Reaper ability now deals additional Shadow damage to targets below 45% instead of below 35%.
			// Additionally, Killing Machine now also increases the critical strike chance of Soul Reaper.
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			// KM effect handled in sim/death_knight/frost/killing_machine.go
			setBonusAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				dk.soulReaper45Percent = true
			}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				dk.soulReaper45Percent = false
			}).ExposeToAPL(138347)
		},
	},
})

// T15 Tank
var ItemSetPlateOfTheAllConsumingMaw = core.NewItemSet(core.ItemSet{
	Name: "Plate of the All-Consuming Maw",
	ID:   1151,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Reduces the cooldown of your Rune Tap ability by 10 sec and removes its Rune cost.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: DeathKnightSpellRuneTap,
				TimeValue: time.Second * -10,
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				ClassMask:  DeathKnightSpellRuneTap,
				FloatValue: -2.0,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Bone Shield ability grants you 15 Runic Power each time one of its charges is consumed.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			rpMetrics := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 138214})

			dk.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(DeathKnightSpellBoneShield) {
					return
				}

				dk.BoneShieldAura.ApplyOnStacksChange(func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if !setBonusAura.IsActive() {
						return
					}

					if newStacks < oldStacks {
						dk.AddRunicPower(sim, 15, rpMetrics)
					}
				})
			})
		},
	},
})

// T16 DPS
var ItemSetBattleplateOfCyclopeanDread = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Cyclopean Dread",
	ID:   1200,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Killing Machine and Sudden Doom grant 500 Haste or Mastery, whichever is highest, for [Dark Transformation: 15 / (Hands * 2 + 4)] sec, stacking up to 10 times.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec == proto.Spec_SpecBloodDeathKnight {
				return
			}

			var duration time.Duration
			if dk.Spec == proto.Spec_SpecUnholyDeathKnight {
				duration = time.Second * 15
			} else if dk.MainHand().HandType == proto.HandType_HandTypeTwoHand {
				duration = time.Second * 8
			} else {
				duration = time.Second * 6
			}

			currentStat := stats.HasteRating
			deathShroudAura := dk.RegisterAura(core.Aura{
				Label:     "Death Shroud" + dk.Label,
				ActionID:  core.ActionID{SpellID: 144901},
				Duration:  duration,
				MaxStacks: 10,

				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					newStat := core.Ternary(
						dk.GetStat(stats.HasteRating) > dk.GetStat(stats.MasteryRating),
						stats.HasteRating,
						stats.MasteryRating)
					if currentStat == newStat {
						dk.AddStatDynamic(sim, currentStat, 500*float64(newStacks-oldStacks))
					} else {
						dk.AddStatDynamic(sim, currentStat, -500*float64(oldStacks))
						dk.AddStatDynamic(sim, newStat, 500*float64(newStacks))
						currentStat = newStat
					}
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: DeathKnightSpellKillingMachine | DeathKnightSpellSuddenDoom,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					deathShroudAura.Activate(sim)
					deathShroudAura.AddStack(sim)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Death Coil increases the duration of Dark Transformation by 2.0 sec per cast.
			// Special attacks while Pillar of Frost is active will impale your target with an icy spike.
			dk := agent.(DeathKnightAgent).GetDeathKnight()

			if dk.Spec == proto.Spec_SpecUnholyDeathKnight {
				setBonusAura.AttachProcTrigger(core.ProcTrigger{
					Callback:       core.CallbackOnSpellHitDealt | core.CallbackOnHealDealt,
					ClassSpellMask: DeathKnightSpellDeathCoil | DeathKnightSpellDeathCoilHeal,

					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if dk.Ghoul.DarkTransformationAura.IsActive() {
							dk.Ghoul.DarkTransformationAura.UpdateExpires(dk.Ghoul.DarkTransformationAura.ExpiresAt() + time.Second*2)
						}
					},
				})
			} else if dk.Spec == proto.Spec_SpecFrostDeathKnight {
				frozenPowerSpell := dk.RegisterSpell(core.SpellConfig{
					ActionID:    core.ActionID{SpellID: 147620},
					SpellSchool: core.SpellSchoolFrost,
					ProcMask:    core.ProcMaskEmpty,
					Flags:       core.SpellFlagPassiveSpell,

					DamageMultiplier: 1,
					CritMultiplier:   dk.DefaultCritMultiplier(),
					ThreatMultiplier: 1,

					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						baseDamage := 500.0 + 0.07999999821*spell.MeleeAttackPower()
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
					},
				})

				dk.OnSpellRegistered(func(spell *core.Spell) {
					if !spell.Matches(DeathKnightSpellPillarOfFrost) {
						return
					}

					dk.PillarOfFrostAura.AttachProcTrigger(core.ProcTrigger{
						Callback: core.CallbackOnSpellHitDealt,
						ProcMask: core.ProcMaskSpecial,
						Outcome:  core.OutcomeLanded,
						Harmful:  true,

						Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
							frozenPowerSpell.Cast(sim, result.Target)
						},
					})
				})
			}
		},
	},
})

// T16 Tank
var ItemSetPlateOfCyclopeanDread = core.NewItemSet(core.ItemSet{
	Name: "Plate of Cyclopean Dread",
	ID:   1201,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Every 10 Heart Strikes, Rune Strikes, Death Coils, Soul Reapers, or Blood Boils will add one charge to your next Bone Shield.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			dk.BoneWallAura = dk.RegisterAura(core.Aura{
				Label:     "Bone Wall" + dk.Label,
				ActionID:  core.ActionID{SpellID: 144948},
				Duration:  time.Minute * 2,
				MaxStacks: 6,
			}).AttachProcTrigger(core.ProcTrigger{
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: DeathKnightSpellBoneShield,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					dk.BoneWallAura.Deactivate(sim)
				},
			})

			boneWallDriver := dk.RegisterAura(core.Aura{
				Label:     "Bone Wall Driver" + dk.Label,
				ActionID:  core.ActionID{SpellID: 145719},
				Duration:  time.Minute * 10,
				MaxStacks: 10,

				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					if newStacks == 10 {
						dk.BoneWallAura.Activate(sim)
						dk.BoneWallAura.AddStack(sim)
						aura.SetStacks(sim, 1)
					}
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback: core.CallbackOnCastComplete,
				ClassSpellMask: DeathKnightSpellHeartStrike |
					DeathKnightSpellRuneStrike |
					DeathKnightSpellDeathCoil |
					DeathKnightSpellSoulReaper |
					DeathKnightSpellBloodBoil,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					boneWallDriver.Activate(sim)
					boneWallDriver.AddStack(sim)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Dancing Rune Weapon will reactivate all Frost and Unholy runes and convert them to Death runes.
			dk := agent.(DeathKnightAgent).GetDeathKnight()
			if dk.Spec != proto.Spec_SpecBloodDeathKnight {
				return
			}

			deathRuneMetrics := dk.NewDeathRuneMetrics(core.ActionID{SpellID: 144950})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: DeathKnightSpellDancingRuneWeapon,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					dk.RegenAllFrostAndUnholyRunesAsDeath(sim, deathRuneMetrics)
				},
			})
		},
	},
})
