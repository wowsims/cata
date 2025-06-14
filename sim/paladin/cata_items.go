package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// PvP set
var ItemSetCataclysmGladiatorsVindication = core.NewItemSet(core.ItemSet{
	ID:   917,
	Name: "Gladiator's Vindication",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases Strength by 70.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Strength, 70)
		},
		// Increases the range of your Judgment by 10 yards.
		// Increases Strength by 90.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			setBonusAura.AttachStatBuff(stats.Strength, 90)
			setBonusAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
				paladin.Judgment.MaxRange += 10
			}).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
				paladin.Judgment.MaxRange -= 10
			})
		},
	},
})

// Increases the damage dealt by your Crusader Strike ability by 5%.
func (paladin *Paladin) addCataclysmPvpGloves() {
	paladin.RegisterPvPGloveMod(
		[]int32{64844, 70649, 60414, 65591, 72379, 70250, 70488, 73707, 73570},
		core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  SpellMaskCrusaderStrike,
			FloatValue: 0.05,
		})
}

// Tier 11 Ret
var ItemSetReinforcedSapphiriumBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battleplate",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the damage done by your Templar's Verdict ability by 10%.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskTemplarsVerdict,
				FloatValue: 0.1,
			})
		},
		// Your Inquisition ability's duration is calculated as if you had one additional Holy Power.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			// Handled in retribution/inquisition.go
			setBonusAura.ExposeToAPL(90299)
			paladin.T11Ret4pc = setBonusAura
		},
	},
})

// Tier 11 Prot
var ItemSetReinforcedSapphiriumBattlearmor = core.NewItemSet(core.ItemSet{
	Name: "Reinforced Sapphirium Battlearmor",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the damage done by your Crusader Strike ability by 10%.
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskCrusaderStrike,
				FloatValue: 0.1,
			})
		},
		// Increases the duration of your Guardian of Ancient Kings ability by 50%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Custom,
				ClassMask: SpellMaskGuardianOfAncientKings,
				ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
					if paladin.AncientPowerAura != nil {
						paladin.AncientPowerAura.Duration = core.DurationFromSeconds(paladin.GoakAura.Duration.Seconds() * 1.5)
					}
					paladin.GoakAura.Duration = core.DurationFromSeconds(paladin.GoakAura.Duration.Seconds() * 1.5)
				},
				RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
					if paladin.AncientPowerAura != nil {
						paladin.AncientPowerAura.Duration = paladin.goakBaseDuration()
					}
					paladin.GoakAura.Duration = paladin.goakBaseDuration()
				},
			})
		},
	},
})

// Tier 12 Ret
var ItemSetBattleplateOfImmolation = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Immolation",
	Bonuses: map[int32]core.ApplySetBonus{
		// Your Crusader Strike deals 15% additional damage as Fire damage over 4 sec.
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
				},

				DamageCalculator: func(result *core.SpellResult) float64 {
					return result.Damage * 0.15
				},
			})
		},
		// Increases damage done by your Judgment by 25%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskJudgment,
				FloatValue: 0.25,
			})
		},
	},
})

// Tier 12 Prot
var ItemSetBattlearmorOfImmolation = core.NewItemSet(core.ItemSet{
	Name: "Battlearmor of Immolation",
	Bonuses: map[int32]core.ApplySetBonus{
		// Your Shield of the Righteous deals 20% additional damage as Fire damage.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			paladin := agent.(PaladinAgent).GetPaladin()
			if paladin.Spec != proto.Spec_SpecProtectionPaladin {
				return
			}

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
		// When your Divine Protection expires, you gain an additional 12% parry chance for 10 sec.
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

// Tier 13 Ret
var ItemSetBattleplateOfRadiantGlory = core.NewItemSet(core.ItemSet{
	Name: "Battleplate of Radiant Glory",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the damage done by Crusader Strike by 15%.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskCrusaderStrike,
				FloatValue: 0.15,
			})
		},
		// Increases damage done by your Templar's Verdict by 20%.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			if agent.GetCharacter().Spec != proto.Spec_SpecRetributionPaladin {
				return
			}

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskTemplarsVerdict,
				FloatValue: 0.2,
			})
		},
	},
})

// Tier 13 Prot
var ItemSetArmorOfRadiantGlory = core.NewItemSet(core.ItemSet{
	Name: "Armor of Radiant Glory",
	Bonuses: map[int32]core.ApplySetBonus{
		// Your Judgment ability now also grants a physical absorption shield equal to 25% of the damage it dealt.
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
				ClassSpellMask: SpellMaskJudgment,
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
		// Reduces the cooldown of Devotion Aura by 30 sec and increases the radius of its effect by 60 yards.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				ClassMask: SpellMaskDevotionAura,
				TimeValue: time.Second * -30,
			})
		},
	},
})
