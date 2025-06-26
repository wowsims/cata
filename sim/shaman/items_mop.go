package shaman

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
)

// T14 elem
// (2) Set: Increases the damage done by your Lightning Bolt spell by 5%.
// (4) Set: Your Rolling Thunder ability now grants 2 Lightning Shield charges each time it triggers.
var ItemSetRegaliaOfTheFirebird = core.NewItemSet(core.ItemSet{
	Name:                    "Regalia of the Firebird",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.05,
				ClassMask:  SpellMaskLightningBolt | SpellMaskLightningBoltOverload,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.T14Ele4pc = setBonusAura
			//talents_elemental.go -> Rolling Thunder
		},
	},
})

// T15 elem
// (2) Set: Your Lightning Bolt, Chain Lighting, and Lava Beam hits have a 10% chance to cause a Lightning Strike at the target's location, dealing 32375 to 37625 Nature damage divided among all non-crowd controlled targets within 10 yards.
// (4) Set: The cooldown of your Ascendance is reduced by 1 sec each time you cast Lava Burst.
var ItemSetRegaliaOfTheWitchDoctor = core.NewItemSet(core.ItemSet{
	Name:                    "Regalia of the Witch Doctor",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()

			lightningStrike := shaman.RegisterSpell(core.SpellConfig{
				ActionID:       core.ActionID{SpellID: 138146},
				SpellSchool:    core.SpellSchoolNature,
				ProcMask:       core.ProcMaskSpellProc,
				CritMultiplier: shaman.DefaultCritMultiplier(),
				MissileSpeed:   20,
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := sim.RollWithLabel(32375, 37625, "Lighting Strike 2pT14")
					nTargets := shaman.Env.GetNumTargets()
					results := make([]*core.SpellResult, nTargets)
					for i, aoeTarget := range sim.Encounter.TargetUnits {
						results[i] = spell.CalcDamage(sim, aoeTarget, baseDamage/float64(nTargets), spell.OutcomeMagicHitAndCrit)
					}
					spell.WaitTravelTime(sim, func(sim *core.Simulation) {
						for i, _ := range sim.Encounter.TargetUnits {
							spell.DealDamage(sim, results[i])
						}
					})
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Regalia of the Witch Doctor 2P",
				Callback:       core.CallbackOnSpellHitDealt,
				Outcome:        core.OutcomeLanded,
				ProcChance:     0.1,
				ClassSpellMask: SpellMaskLightningBolt | SpellMaskChainLightningOverload | SpellMaskLavaBeam | SpellMaskLavaBeamOverload | SpellMaskChainLightning | SpellMaskChainLightningOverload,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					lightningStrike.Cast(sim, result.Target)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Regalia of the Witch Doctor 4P",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskLavaBurst,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shaman.Ascendance.CD.Reduce(time.Millisecond * 1500) //simc says 1.5s in a hotfix
					shaman.UpdateMajorCooldowns()
				},
			})
		},
	},
})

// T16 elem
// (2) Set: Fulmination increases all Fire and Nature damage dealt to that target from the Shaman by 4% for 2 sec per Lightning Shield charge consumed.
// (4) Set: Your Lightning Bolt and Chain Lightning spells have a chance to summon a Lightning Elemental to fight by your side for 10 sec.
var ItemSetCelestialHarmonyRegalia = core.NewItemSet(core.ItemSet{
	Name:                    "Celestial Harmony Regalia",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			debuffAuras := shaman.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
				return target.GetOrRegisterAura(core.Aura{
					Label:     "Elemental Discharge - " + shaman.Label,
					ActionID:  core.ActionID{SpellID: 144999},
					Duration:  time.Second * 2,
					MaxStacks: 6,
					OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
						aura.Duration = time.Second * 2 * time.Duration(newStacks)
						aura.Refresh(sim)
						core.EnableDamageDoneByCaster(DDBC_2PT16, DDBC_Total, shaman.AttackTables[aura.Unit.UnitIndex], func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
							if spell.SpellSchool.Matches(core.SpellSchoolNature | core.SpellSchoolFire) {
								//TODO Does the damage taken also increases with LS stacks ?
								return 1.0 + float64(newStacks)*0.04
							}
							return 1.0
						})
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						aura.Duration = time.Second * 2
						core.DisableDamageDoneByCaster(DDBC_2PT16, shaman.AttackTables[aura.Unit.UnitIndex])
					},
				})
			})
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Celestial Harmony Regalia 2P",
				Callback:       core.CallbackOnSpellHitDealt,
				Outcome:        core.OutcomeLanded,
				ClassSpellMask: SpellMaskFulmination,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					debuff := debuffAuras.Get(result.Target)
					debuff.Activate(sim)
					debuff.SetStacks(sim, shaman.LightningShieldAura.GetStacks()-1)
				},
			})

		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Celestial Harmony Regalia 4P",
				Callback:       core.CallbackOnSpellHitDealt,
				Outcome:        core.OutcomeLanded,
				ClassSpellMask: SpellMaskLightningBolt | SpellMaskChainLightning,
				ICD:            time.Second * 60,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					//TODO
				},
			})
		},
	},
})

// T14 enh
// (2) Set: Increases the damage done by your Lava Lash ability by 15%.
// (4) Set: Increases the critical strike chance bonus from your Stormstrike ability by an additional 15%.
var ItemSetBattlegearOfTheFirebird = core.NewItemSet(core.ItemSet{
	Name:                    "Battlegear of the Firebird",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  SpellMaskLavaLash,
				FloatValue: 0.15,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.T14Enh4pc = setBonusAura
			//in shaman.go
		},
	},
})

// T15 enh
// (2) Set: Your Stormstrike also grants you 2 additional charges of Maelstrom Weapon.
// (4) Set: The cooldown of your Feral Spirits is reduced by 8 sec each time Windfury Weapon is triggered.
var ItemSetBattlegearOfTheWitchDoctor = core.NewItemSet(core.ItemSet{
	Name:                    "Battlegear of the Witch Doctor",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Battlegear of the Witch Doctor 2P",
				Callback:       core.CallbackOnSpellHitDealt,
				Outcome:        core.OutcomeLanded,
				ClassSpellMask: SpellMaskStormstrikeCast,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shaman.MaelstromWeaponAura.Activate(sim)
					shaman.MaelstromWeaponAura.SetStacks(sim, shaman.MaelstromWeaponAura.GetStacks()+2)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Battlegear of the Witch Doctor 4P",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskWindfuryWeapon,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shaman.FeralSpirit.CD.Reduce(time.Second * 8)
					shaman.UpdateMajorCooldowns()
				},
			})
		},
	},
})

// T16 enh
// 2 pieces: For 10 sec after using Unleash Elements, your attacks have a chance to unleash a random weapon imbue.
// 4 pieces: When Flame Shock deals periodic damage, you have a 5% chance to gain 5 stacks of Searing Flames and reset the cooldown of Lava Lash.
var ItemSetCelesialHarmonyBattlegear = core.NewItemSet(core.ItemSet{
	Name:                    "Celestial Harmony Battlegear",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			var imbueSpells []*core.Spell
			shaman.OnSpellRegistered(func(spell *core.Spell) {
				if spell.Matches(SpellMaskWindfuryWeapon | SpellMaskFrostbrandWeapon | SpellMaskFlametongueWeapon) {
					imbueSpells = append(imbueSpells, spell)
				}
			})
			procAura := core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
				Name:       "Celestial Harmony Battlegear 2P Proc",
				Callback:   core.CallbackOnSpellHitDealt,
				Outcome:    core.OutcomeLanded,
				ProcMask:   core.ProcMaskMeleeOrMeleeProc,
				ICD:        time.Millisecond * 100,
				ProcChance: 0.1,
				Duration:   time.Second * 10,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if len(imbueSpells) == 0 {
						return
					}
					rand := int(math.Floor(sim.RollWithLabel(0, float64(len(imbueSpells)), "Enh 4PT16 Proc")))
					imbueSpells[rand].Cast(sim, result.Target)
				},
			})
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Celestial Harmony Battlegear 2P",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskUnleashElements,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Celestial Harmony Battlegear 4P",
				Callback:       core.CallbackOnPeriodicDamageDealt,
				ClassSpellMask: SpellMaskFlameShockDot,
				ProcChance:     0.05,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					sfAura := shaman.GetAura("Searing Flames")
					if sfAura != nil {
						sfAura.Activate(sim)
						sfAura.SetStacks(sim, sfAura.GetStacks()+5)
						shaman.LavaLash.CD.Reset()
					}
				},
			})
		},
	},
})

// S12 enh
// (2) Set: Increases the chance to trigger your Maelstrom Weapon talent by 20%.
// (4) Set: While your weapon is imbued with Flametongue Weapon, your attacks also slow the target's movement speed by 50% for 3 sec.
var ItemSetGladiatorsEarthshaker = core.NewItemSet(core.ItemSet{
	Name:                    "Gladiator's Earthshaker",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			shaman := agent.(ShamanAgent).GetShaman()
			shaman.S12Enh2pc = setBonusAura
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {},
	},
})

func init() {
}
