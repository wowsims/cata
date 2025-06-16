package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

var ItemSetMistsGladiatorsVindication = core.NewItemSet(core.ItemSet{
	ID:   1111,
	Name: "Gladiator's Vindication",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
		},
		/*
			You gain a charge of Holy Power whenever you take direct damage.
			This effect cannot occur more than once every 8 seconds.
		*/
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			actionID := core.ActionID{SpellID: 131649}
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback: core.CallbackOnSpellHitTaken,
				Harmful:  true,
				ICD:      time.Second * 8,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					paladin.HolyPower.Gain(sim, 1, actionID)
				},
			})
		},
	},
})

// Increases the range of your Judgment by 10 yards.
func (paladin *Paladin) addMistsPvpGloves() {
	paladin.RegisterPvPGloveMod(
		[]int32{84419, 84834, 85027, 91269, 91270, 91622, 93528, 94343, 98844, 99871, 100013, 100365, 100573, 102630, 102827, 103243, 103440},
		core.SpellModConfig{
			Kind:      core.SpellMod_Custom,
			ClassMask: SpellMaskJudgment,
			ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
				spell.MaxRange += 10
			},
			RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
				spell.MaxRange -= 10
			},
		})
}

// Tier 14 Ret
var ItemSetWhiteTigerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "White Tiger Battlegear",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the damage done by your Templar's Verdict ability by 15%.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			if agent.GetCharacter().Spec != proto.Spec_SpecRetributionPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskTemplarsVerdict,
				FloatValue: 0.15,
			}).ExposeToAPL(123108)
		},
		// Your Seals and Judgments deal 10% additional damage.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskJudgment | SpellMaskSeals, // Censure?!
				FloatValue: 0.10,
			}).ExposeToAPL(70762)
		},
	},
})

// Tier 14 Prot
var ItemSetWhiteTigerPlate = core.NewItemSet(core.ItemSet{
	Name: "White Tiger Plate",
	Bonuses: map[int32]core.ApplySetBonus{
		// Reduces the cooldown of your Ardent Defender ability by 60 sec.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecProtectionPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: SpellMaskArdentDefender,
				TimeValue: time.Second * -60,
			}).ExposeToAPL(123104)
		},
		// Increases the healing done by your Word of Glory spell by 10% and increases the damage reduction of your Shield of the Righteous ability by 10%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskWordOfGlory,
				FloatValue: 0.1,
			}).AttachAdditivePseudoStatBuff(&paladin.ShieldOfTheRighteousMultiplicativeMultiplier, 0.1)

			setBonusAura.ExposeToAPL(123107)
		},
	},
})

// Tier 14 Holy
var ItemSetWhiteTigerVestments = core.NewItemSet(core.ItemSet{
	Name: "White Tiger Vestments",
	Bonuses: map[int32]core.ApplySetBonus{
		// Reduces the mana cost of your Holy Radiance spell by 10%.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecHolyPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				ClassMask:  SpellMaskHolyRadiance,
				FloatValue: -0.1,
			}).ExposeToAPL(123102)
		},
		// Reduces the cooldown of your Holy Shock spell by 1 sec.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecHolyPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: SpellMaskHolyShock,
				TimeValue: time.Second * -1,
			}).ExposeToAPL(123103)
		},
	},
})

func (paladin *Paladin) registerHolyDamageTemplarsVerdict() *core.Spell {
	actionID := core.ActionID{SpellID: 85256}

	return paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID.WithTag(2), // Actual 138165
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: SpellMaskTemplarsVerdict,

		DamageMultiplier: 2.75,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + paladin.CalcScalingSpellDmg(0.55000001192)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.HolyPower.Spend(sim, 3, actionID)
				paladin.T15Ret4pc.Deactivate(sim)
			}

			spell.DealDamage(sim, result)
		},
	})
}

// Tier 15 Ret
var ItemSetBattlegearOfTheLightningEmperor = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of the Lightning Emperor",
	Bonuses: map[int32]core.ApplySetBonus{
		// Your Exorcism causes your target to take 6% increased Holy damage from your attacks for 6 sec.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecRetributionPaladin {
				return
			}

			exorcismAuras := paladin.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
				return unit.RegisterAura(core.Aura{
					Label:    "Exorcism" + unit.Label,
					ActionID: core.ActionID{SpellID: 138162},
					Duration: time.Second * 6,
				}).AttachMultiplicativePseudoStatBuff(
					&unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly], 1.06,
				)
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskExorcism,
				Outcome:        core.OutcomeLanded,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					exorcismAuras.Get(result.Target).Activate(sim)
				},
			}).ExposeToAPL(138159)
		},
		// Your Crusader Strike has a 40% chance to make your next Templar's Verdict deal all Holy damage.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecRetributionPaladin {
				return
			}

			templarsVerdictAura := paladin.RegisterAura(core.Aura{
				Label:    "Templar's Verdict" + paladin.Label,
				ActionID: core.ActionID{SpellID: 138169},
				Duration: core.NeverExpires,
			})
			paladin.T15Ret4pc = templarsVerdictAura
			paladin.T15Ret4pcTemplarsVerdict = paladin.registerHolyDamageTemplarsVerdict()

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskCrusaderStrike,
				Outcome:        core.OutcomeLanded,
				ProcChance:     0.4,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					templarsVerdictAura.Activate(sim)
				},
			}).ExposeToAPL(138164)
		},
	},
})

// Tier 15 Prot
var ItemSetPlateOfTheLightningEmperor = core.NewItemSet(core.ItemSet{
	Name: "Plate of the Lightning Emperor",
	Bonuses: map[int32]core.ApplySetBonus{
		// Casting Word of Glory or Eternal Flame also grants you 40% additional block chance for 5 sec per Holy Power.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			shieldOfGloryAura := paladin.RegisterAura(core.Aura{
				Label:    "Shield of Glory" + paladin.Label,
				ActionID: core.ActionID{SpellID: 138242},
				Duration: time.Second * 5,
			}).AttachAdditivePseudoStatBuff(&paladin.PseudoStats.BaseBlockChance, 0.4)

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskWordOfGlory,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shieldOfGloryAura.Duration = core.DurationFromSeconds(float64(paladin.DynamicHolyPowerSpent * 5.0))
					shieldOfGloryAura.Activate(sim)
				},
			}).ExposeToAPL(138238)
		},
		// You gain 1 Holy Power for each 20% of your health taken as damage while Divine Protection is active.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			hpGainMetrics := core.ActionID{SpellID: 138248}

			totalDamageTaken := 0.0

			paladin.OnSpellRegistered(func(spell *core.Spell) {
				if spell.Matches(SpellMaskDivineProtection) {
					paladin.DivineProtectionAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
						totalDamageTaken = 0.0
					})
				}
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback: core.CallbackOnSpellHitTaken,
				Outcome:  core.OutcomeLanded,
				Harmful:  true,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if paladin.DivineProtectionAura.IsActive() {
						totalDamageTaken += result.Damage
						if totalDamageTaken >= paladin.MaxHealth()*0.2 {
							paladin.HolyPower.Gain(sim, 1, hpGainMetrics)
							totalDamageTaken = 0
						}
					}
				},
			}).ExposeToAPL(138244)
		},
	},
})

// Tier 15 Holy
var ItemSetVestmentsOfTheLightningEmperor = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Lightning Emperor",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the healing done by your Daybreak ability by 50%.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecHolyPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskDaybreak,
				FloatValue: 0.5,
			}).ExposeToAPL(138291)
		},
		// Increases the healing transferred to your Beacon of Light target by 20%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecHolyPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskBeaconOfLight,
				FloatValue: 0.2,
			}).ExposeToAPL(138292)
		},
	},
})

// Tier 16 Ret
var ItemSetBattlegearOfWingedTriumph = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of Winged Triumph",
	Bonuses: map[int32]core.ApplySetBonus{
		// When Art of War activates, all damage is increased by 5% for 6 sec.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecRetributionPaladin {
				return
			}

			warriorOfTheLightAura := paladin.RegisterAura(core.Aura{
				Label:    "Warrior of the Light" + paladin.Label,
				ActionID: core.ActionID{SpellID: 144587},
				Duration: time.Second * 6,
			}).AttachMultiplicativePseudoStatBuff(&paladin.PseudoStats.DamageDealtMultiplier, 1.05)

			setBonusAura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
				paladin.TheArtOfWarAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					warriorOfTheLightAura.Activate(sim)
				})
			}).ExposeToAPL(144586)
		},
		// Holy Power consumers have a 25% chance to make your next Divine Storm free and deal 50% more damage.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecRetributionPaladin {
				return
			}

			paladin.DivineCrusaderAura = paladin.divinePurposeFactory("Divine Crusader", 144595, time.Second*12, func(aura *core.Aura, spell *core.Spell) bool {
				return spell.Matches(SpellMaskDivineStorm)
			}).AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskDivineStorm,
				FloatValue: 0.5,
			})

			setBonusAura.ExposeToAPL(144593)
		},
	},
})

// Tier 16 Prot
var ItemSetPlateOfWingedTriumph = core.NewItemSet(core.ItemSet{
	Name: "Plate of Winged Triumph",
	Bonuses: map[int32]core.ApplySetBonus{
		// While Divine Protection is active, 75% of the damage taken is converted into a heal over time that activates when Divine Protection fades.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			totalDamageTaken := 0.0

			blessingOfTheGuardians := paladin.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 144581},
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskSpellHealing,
				Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

				Hot: core.DotConfig{
					Aura: core.Aura{
						Label: "Blessing of the Guardians" + paladin.Label,
					},
					TickLength:          time.Second,
					NumberOfTicks:       10,
					AffectedByCastSpeed: false,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
						dot.Snapshot(target, totalDamageTaken/10.0)
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeSnapshotCrit)
					},
				},

				DamageMultiplier: 0.75 * core.TernaryFloat64(paladin.Talents.UnbreakableSpirit, 0.5, 1.0),
				CritMultiplier:   paladin.DefaultCritMultiplier(),
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Hot(target).Apply(sim)
				},
			})

			paladin.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(SpellMaskDivineProtection) {
					return
				}

				paladin.DivineProtectionAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					totalDamageTaken = 0.0
				}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					if totalDamageTaken > 0 {
						blessingOfTheGuardians.Cast(sim, &paladin.Unit)
					}
				}).AttachProcTrigger(core.ProcTrigger{
					Callback: core.CallbackOnSpellHitTaken,
					Outcome:  core.OutcomeLanded,
					Harmful:  true,

					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						totalDamageTaken += result.Damage
					},
				})
			})

			setBonusAura.ExposeToAPL(144580)
		},
		// While at 3 or more stacks of Bastion of Glory, your next [Eternal Flame / Word of Glory] will consume no Holy Power and count as if 3 Holy Power were consumed.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecProtectionPaladin {
				return
			}

			paladin.BastionOfPowerAura = paladin.RegisterAura(core.Aura{
				Label:    "Bastion of Power" + paladin.Label,
				ActionID: core.ActionID{SpellID: 144569},
				Duration: time.Second * 20,
			}).AttachProcTrigger(core.ProcTrigger{
				Name:           "Bastion of Power Consume Trigger" + paladin.Label,
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskWordOfGlory,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					paladin.BastionOfPowerAura.Deactivate(sim)
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskShieldOfTheRighteous,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if paladin.BastionOfGloryAura.GetStacks() >= 3 {
						paladin.BastionOfPowerAura.Activate(sim)
					}
				},
			}).ExposeToAPL(144566)
		},
	},
})

// Tier 16 Holy
var ItemSetVestmentsOfWingedTriumph = core.NewItemSet(core.ItemSet{
	Name: "Vestments of Winged Triumph",
	Bonuses: map[int32]core.ApplySetBonus{
		// Infusion of Light also increases the healing done by Holy Light, Divine Light, and Holy Radiance by 25%.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecHolyPaladin {
				return
			}

			setBonusAura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
				paladin.InfusionOfLightAura.AttachDependentAura(paladin.RegisterAura(core.Aura{
					Label:    "Unyielding Faith" + paladin.Label,
					ActionID: core.ActionID{SpellID: 144624},
					Duration: time.Second * 15,
				}).AttachSpellMod(core.SpellModConfig{
					Kind:       core.SpellMod_DamageDone_Pct,
					ClassMask:  SpellMaskDivineLight | SpellMaskHolyLight | SpellMaskHolyRadiance,
					FloatValue: 0.25,
				}))
			}).ExposeToAPL(144625)
		},
		/*
			Reduces the cooldown of Divine Favor by 60 sec.
			While Divine Favor is active, Mastery is increased by 4500.
		*/
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecHolyPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: SpellMaskDivineFavor,
				TimeValue: time.Second * -60,
			})

			setBonusAura.ApplyOnInit(func(aura *core.Aura, sim *core.Simulation) {
				paladin.DivineFavorAura.AttachDependentAura(paladin.RegisterAura(core.Aura{
					Label:    "Favor of the Kings" + paladin.Label,
					ActionID: core.ActionID{SpellID: 144622},
					Duration: time.Second * 20,
				}).AttachStatBuff(stats.MasteryRating, 4500))
			}).ExposeToAPL(144613)
		},
	},
})
