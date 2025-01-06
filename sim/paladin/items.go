package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// Tier 11 ret
var ItemSetReinforcedSapphiriumBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battleplate",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  SpellMaskTemplarsVerdict,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.(PaladinAgent).GetPaladin()

			setBonusAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				character.HasT11Ret4pc = true
			})
			setBonusAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				character.HasT11Ret4pc = false
			})

			// Handled in inquisition.go
			setBonusAura.ExposeToAPL(90299)
		},
	},
})

// Tier 12 ret
var ItemSetBattleplateOfImmolation = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Immolation",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			cata.RegisterIgniteEffect(&paladin.Unit, cata.IgniteConfig{
				ActionID:           core.ActionID{SpellID: 35395}.WithTag(3), // actual 99092
				DisableCastMetrics: true,
				DotAuraLabel:       "Flames of the Faithful" + paladin.Label,
				IncludeAuraDelay:   true,
				SetBonusAura:       setBonusAura,

				ProcTrigger: core.ProcTrigger{
					Name:           "Flames of the Faithful" + paladin.Label,
					Callback:       core.CallbackOnSpellHitDealt,
					ClassSpellMask: SpellMaskCrusaderStrike,
					Outcome:        core.OutcomeLanded,
					ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
						return setBonusAura.IsActive()
					},
				},

				DamageCalculator: func(result *core.SpellResult) float64 {
					return result.Damage * 0.15
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Handled in talents_retribution.go
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(SpellMaskZealotry) {
					return
				}

				onEquip := func() {
					paladin.ZealotryAura.Duration += time.Second * 15
				}

				setBonusAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					onEquip()
				})
				setBonusAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					paladin.ZealotryAura.Duration -= time.Second * 15
				})

				if setBonusAura.IsActive() {
					onEquip()
				}
			})

			setBonusAura.ExposeToAPL(99116)
		},
	},
})

// Tier 13 ret
var ItemSetBattleplateOfRadiantGlory = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Radiant Glory",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			// Actual buff credited with the Holy Power gain is Virtuous Empowerment
			hpMetrics := paladin.NewHolyPowerMetrics(core.ActionID{SpellID: 105767})

			// Used for checking "Is Aura Known" in the APL
			paladin.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 105767},
				Label:    "Virtuous Empowerment" + paladin.Label,
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "T13 2pc trigger" + paladin.Label,
				ActionID:       core.ActionID{SpellID: 105765},
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskJudgement,
				ProcChance:     1,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					// TODO: Measure the aura update delay distribution on PTR.
					var delaySeconds float64
					if sim.Proc(0.75, "T13 2pc") {
						delaySeconds = 0.010 * sim.RandomFloat("T13 2pc"+paladin.Label)
					} else {
						delaySeconds = 0.090 + 0.020*sim.RandomFloat("T13 2pc"+paladin.Label)
					}

					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt:     sim.CurrentTime + core.DurationFromSeconds(delaySeconds),
						Priority: core.ActionPriorityRegen,

						OnAction: func(_ *core.Simulation) {
							paladin.GainHolyPower(sim, 1, hpMetrics)
						},
					})
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			damageMod := paladin.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskModifiedByZealOfTheCrusader,
				FloatValue: 0.18,
			})

			zealOfTheCrusader := paladin.RegisterAura(core.Aura{
				Label:    "Zeal of the Crusader" + paladin.Label,
				ActionID: core.ActionID{SpellID: 105819},
				Duration: time.Second * 20,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					damageMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					damageMod.Deactivate()
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "T13 4pc trigger" + paladin.Label,
				ActionID:       core.ActionID{SpellID: 105820},
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskZealotry,
				ProcChance:     1,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					zealOfTheCrusader.Activate(sim)
				},
			})
		},
	},
})

// PvP set
var ItemSetGladiatorsVindication = core.NewItemSet(core.ItemSet{
	ID:   917,
	Name: "Gladiator's Vindication",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Strength, 70)
		},
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Strength, 90)
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: SpellMaskJudgementBase,
				TimeValue: -1 * time.Second,
			})
		},
	},
})

func (paladin *Paladin) addBloodthirstyGloves() {
	spellMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  SpellMaskCrusaderStrike,
		FloatValue: 0.05,
	})

	checkGloves := func() {
		switch paladin.Hands().ID {
		case 64844, 70649, 60414, 65591, 72379, 70250, 70488, 73707, 73570:
			spellMod.Activate()
		default:
			spellMod.Deactivate()
		}
	}

	if paladin.ItemSwap.IsEnabled() {
		paladin.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotHands}, func(_ *core.Simulation, _ proto.ItemSlot) {
			checkGloves()
		})
	} else {
		checkGloves()
	}
}

// Tier 11 prot
var ItemSetReinforcedSapphiriumBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battlearmor",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				ClassMask:  SpellMaskCrusaderStrike,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			paladin.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(SpellMaskGuardianOfAncientKings) {
					return
				}

				goakBaseDuration := paladin.GoakAura.Duration
				acientPowerBaseDuration := paladin.GoakAura.Duration

				applyT11Prot4pcBonus := func(duration time.Duration) time.Duration {
					return time.Millisecond * time.Duration(float64(duration.Milliseconds())*1.5)
				}

				onEquip := func() {
					if paladin.AncientPowerAura != nil {
						paladin.AncientPowerAura.Duration = applyT11Prot4pcBonus(acientPowerBaseDuration)
					}
					paladin.GoakAura.Duration = applyT11Prot4pcBonus(goakBaseDuration)
				}

				setBonusAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
					onEquip()
				})

				setBonusAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
					if paladin.AncientPowerAura != nil {
						paladin.AncientPowerAura.Duration = acientPowerBaseDuration
					}
					paladin.GoakAura.Duration = goakBaseDuration
				})

				if setBonusAura.IsActive() {
					onEquip()
				}
			})
		},
	},
})

// Tier 12 prot
var ItemSetBattlearmorOfImmolation = core.NewItemSet(core.ItemSet{
	Name: "Battlearmor of Immolation",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			procDamage := 0.0

			righteousFlames := paladin.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 53600}.WithTag(3), // actual 99075
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,
				Flags: core.SpellFlagIgnoreModifiers |
					core.SpellFlagBinary |
					core.SpellFlagNoOnCastComplete |
					core.SpellFlagNoOnDamageDealt,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.CalcAndDealDamage(sim, target, procDamage, spell.OutcomeAlwaysHit)
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Righteous Flames" + paladin.Label,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskShieldOfTheRighteous,
				Outcome:        core.OutcomeLanded,
				ProcChance:     1,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procDamage = result.Damage * 0.2
					righteousFlames.Cast(sim, result.Target)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			flamingAegis := paladin.GetOrRegisterAura(core.Aura{
				Label:    "Flaming Aegis" + paladin.Label,
				ActionID: core.ActionID{SpellID: 99090},
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					paladin.PseudoStats.BaseParryChance += 0.12
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					paladin.PseudoStats.BaseParryChance -= 0.12
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "T12 4pc trigger" + paladin.Label,
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskDivineProtection,
				ProcChance:     1,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					core.StartDelayedAction(sim, core.DelayedActionOptions{
						DoAt:     sim.CurrentTime + paladin.DivineProtectionAura.Duration,
						Priority: core.ActionPriorityLow,

						OnAction: func(_ *core.Simulation) {
							flamingAegis.Activate(sim)
						},
					})
				},
			})
		},
	},
})

// Tier 13 prot
var ItemSetArmorOfRadiantGlory = core.NewItemSet(core.ItemSet{
	Name: "Armor of Radiant Glory",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			actionID := core.ActionID{SpellID: 105801}
			duration := time.Second * 6

			shieldStrength := 0.0
			shield := paladin.NewDamageAbsorptionAura("Delayed Judgement"+paladin.Label, actionID, duration, func(unit *core.Unit) float64 {
				return shieldStrength
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Delayed Judgement Proc" + paladin.Label,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskJudgement,
				Outcome:        core.OutcomeLanded,

				ProcChance: 1,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shieldStrength = result.Damage * 0.25
					if shieldStrength > 1 {
						shield.Activate(sim)
					}
				},
			})

			setBonusAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				shieldStrength = 0
			})
		},
		4: func(_ core.Agent, _ *core.Aura) {
			// Divine Guardian not implemented since it's a raid cooldown and doesn't affect the Paladin
		},
	},
})
