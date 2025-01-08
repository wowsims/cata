package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// Dungeon Set 3 Tidefury Raiment
// (2) Set: Your Chain Lightning Spell now only loses 17% of its damage per jump.
// (4) Set: Your Water Shield ability grants an additional 56 mana each time it triggers and an additional 3 mana per 5 sec.
var ItemSetTidefury = core.NewItemSet(core.ItemSet{
	Name: "Tidefury Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Handled in chain_lightning.go
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			if shaman.SelfBuffs.Shield == proto.ShamanShield_WaterShield {
				shaman.AddStat(stats.MP5, 3)
			}
		},
	},
})

// var ItemSetSkyshatterRegalia = core.NewItemSet(core.ItemSet{
// 	Name: "Skyshatter Regalia",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			shaman := agent.(ShamanAgent).GetShaman()

// 			if shaman.Totems.Air == proto.AirTotem_NoAirTotem ||
// 				shaman.Totems.Water == proto.WaterTotem_NoWaterTotem ||
// 				shaman.Totems.Earth == proto.EarthTotem_NoEarthTotem ||
// 				shaman.Totems.Fire == proto.FireTotem_NoFireTotem {
// 				return
// 			}

// 			shaman.AddStat(stats.MP5, 19)
// 			shaman.AddStat(stats.SpellCrit, 35)
// 			shaman.AddStat(stats.SpellPower, 45)
// 		},
// 		4: func(agent core.Agent) {
// 			// Increases damage done by Lightning Bolt by 5%.
// 			// Implemented in lightning_bolt.go.
// 		},
// 	},
// })

// Skyshatter Harness
// 2 pieces: Your Earth Shock, Flame Shock, and Frost Shock abilities cost 10% less mana.
// 4 pieces: Whenever you use Stormstrike, you gain 70 attack power for 12 sec.

// var ItemSetSkyshatterHarness = core.NewItemSet(core.ItemSet{
// 	Name: "Skyshatter Harness",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		2: func(agent core.Agent) {
// 			// implemented in shocks.go
// 		},
// 		4: func(agent core.Agent) {
// 			// implemented in stormstrike.go
// 		},
// 	},
// })

// T11 elem
// (2) Set: Increases the critical strike chance of your Flame Shock spell by 10%.
// (4) Set: Reduces the cast time of your Lightning Bolt spell by 10%.
var ItemSetRagingElementsRegalia = core.NewItemSet(core.ItemSet{
	Name: "Regalia of the Raging Elements",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 10,
				ClassMask:  SpellMaskFlameShock,
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_CastTime_Pct,
				FloatValue: -0.1,
				ClassMask:  SpellMaskLightningBolt,
			})
		},
	},
})

// T12 elem
// (2) Set: Your Lightning Bolt has a 30% chance to reduce the remaining cooldown on your Fire Elemental Totem by 4 sec.
// (4) Set: Your Lava Surge talent also makes Lava Burst instant when it triggers.
var ItemSetVolcanicRegalia = core.NewItemSet(core.ItemSet{
	Name: "Volcanic Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			if shaman.useDragonSoul_2PT12 {
				shaman.RegisterAura(core.Aura{
					Label:    "Volcanic Regalia 2P",
					Duration: core.NeverExpires,
					OnReset: func(aura *core.Aura, sim *core.Simulation) {
						aura.Activate(sim)
					},
					OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if spell.ClassSpellMask != SpellMaskLightningBolt || !sim.Proc(0.3, "Volcanic Regalia 2P") || shaman.FireElementalTotem == nil {
							return
						}
						shaman.FireElementalTotem.CD.Reduce(4 * time.Second)
					},
				})
			} else {
				core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
					Name:       "Volcanic Regalia 2P",
					Callback:   core.CallbackOnSpellHitDealt,
					ProcMask:   core.ProcMaskSpellDamage,
					Outcome:    core.OutcomeLanded,
					ProcChance: 0.08,
					ICD:        time.Second * 105,
					ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
						return !spell.Matches(SpellMaskOverload)
					},
					Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
						shaman.FireElementalTotem.CD.Reset()
					},
				})
			}

		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			instantLavaSurgeMod := shaman.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_CastTime_Pct,
				FloatValue: -1,
				ClassMask:  SpellMaskLavaBurst,
			})
			shaman.VolcanicRegalia4PT12Aura = shaman.RegisterAura(core.Aura{
				Label:    "Volcano",
				ActionID: core.ActionID{SpellID: 99207},
				Duration: 10 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					instantLavaSurgeMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					instantLavaSurgeMod.Deactivate()
				},
			})
			//in talents.go under lava surge proc
		},
	},
})

// T13 elem
// (2) Set: Elemental Mastery also grants you 2000 mastery rating 15 sec.
// (4) Set: Each time Elemental Overload triggers, you gain 250 haste rating for 4 sec, stacking up to 3 times.
var ItemSetSpiritwalkersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Spiritwalker's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			//In talents.go under elemental mastery
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()
			procAura := shaman.RegisterAura(core.Aura{
				Label:     "Time Rupture",
				ActionID:  core.ActionID{SpellID: 105821},
				Duration:  4 * time.Second,
				MaxStacks: 3,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					changedHasteRating := (newStacks - oldStacks) * 250
					shaman.AddStatDynamic(sim, stats.HasteRating, float64(changedHasteRating))
				},
			})
			shaman.RegisterAura(core.Aura{
				Label:    "Spiritwalker's Regalia 4P",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				//TODO does dtr overloads proc this ?
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.ClassSpellMask&SpellMaskOverload > 0 {
						procAura.Activate(sim)
						procAura.AddStack(sim)
					}
				},
			})
		},
	},
})

// T11 enh
// (2) Set: Increases damage done by your Lava Lash and Stormstrike abilities by 10%.
// (4) Set: Increases the critical strike chance of your Lightning Bolt spell by 10%.
var ItemSetRagingElementsBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Battlegear of the Raging Elements",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: .10,
				ClassMask:  SpellMaskLavaLash | SpellMaskStormstrike,
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 10,
				ClassMask:  SpellMaskLightningBolt,
			})
		},
	},
})

func tier12StormstrikeBonus(_ *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
	if spell.ClassSpellMask&(SpellMaskFire|SpellMaskFlametongueWeapon) > 0 {
		return 1.06
	}
	return 1.0
}

// T12 enh
// (2) Set: Your Lava Lash gains an additional 5% damage increase per application of Searing Flames on the target.
// (4) Set: Your Stormstrike ability also causes the target to take 6% increased damage from your Fire Nova, Flame Shock, Flametongue Weapon, Lava Burst, Lava Lash, and Unleash Flame abilities.
var ItemSetVolcanicBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Volcanic Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			shaman.SearingFlamesMultiplier += 0.05
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			stormFireAuras := shaman.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
				return target.GetOrRegisterAura(core.Aura{
					Label:    "Stormfire-" + shaman.Label,
					ActionID: core.ActionID{SpellID: 99212},
					Duration: time.Second * 15,
					OnGain: func(aura *core.Aura, _ *core.Simulation) {
						core.EnableDamageDoneByCaster(DDBC_4pcT12, DDBC_Total, shaman.AttackTables[aura.Unit.UnitIndex], tier12StormstrikeBonus)
					},
					OnExpire: func(aura *core.Aura, _ *core.Simulation) {
						core.DisableDamageDoneByCaster(DDBC_4pcT12, shaman.AttackTables[aura.Unit.UnitIndex])
					},
				})
			})

			core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
				Name:           "Stormfire Trigger",
				ActionID:       core.ActionID{SpellID: 99213},
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskStormstrikeCast,
				Outcome:        core.OutcomeLanded,
				ProcChance:     1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					stormFire := stormFireAuras.Get(result.Target)
					stormFire.Activate(sim)
				},
			})
		},
	},
})

// T13 enh
// 2 pieces: While you have any stacks of Maelstrom Weapon, your Lightning Bolt, Chain Lightning, and healing spells deal 20% more healing or damage.
// 4 pieces: Your Feral Spirits have a 45% chance to grant you a charge of Maelstrom Weapon each time they deal damage.
var ItemSetSpiritwalkersBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Spiritwalker's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			if shaman.Talents.MaelstromWeapon == 0 {
				return
			}

			// Item sets are registered before talents, so MaelstromWeaponAura doesn't exist yet
			// Therefore we need to react on Feral Spirit registration to apply the logic
			shaman.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ClassSpellMask&SpellMaskFeralSpirit == 0 {
					return
				}

				dmgMod := shaman.AddDynamicMod(core.SpellModConfig{
					Kind:       core.SpellMod_DamageDone_Pct,
					FloatValue: 0.2,
					ClassMask:  SpellMaskLightningBolt | SpellMaskChainLightning,
				})

				temporalMaelstrom := shaman.RegisterAura(core.Aura{
					Label:    "Temporal Maelstrom" + shaman.Label,
					ActionID: core.ActionID{SpellID: 105869},
					Duration: time.Second * 30,
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						dmgMod.Activate()
						shaman.PseudoStats.HealingDealtMultiplier *= 1.2
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						dmgMod.Deactivate()
						shaman.PseudoStats.HealingDealtMultiplier /= 1.2
					},
				})

				shaman.MaelstromWeaponAura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
					temporalMaelstrom.Activate(sim)
				})

				shaman.MaelstromWeaponAura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
					temporalMaelstrom.Deactivate(sim)
				})
			})
		},
		4: func(agent core.Agent) {
			shaman := agent.(ShamanAgent).GetShaman()

			if !shaman.Talents.FeralSpirit || shaman.Talents.MaelstromWeapon == 0 {
				return
			}

			for _, wolf := range []*core.Unit{&shaman.SpiritWolves.SpiritWolf1.Unit, &shaman.SpiritWolves.SpiritWolf2.Unit} {
				core.MakeProcTriggerAura(wolf, core.ProcTrigger{
					Name:       "Spiritwalker's Battlegear 4pc" + shaman.Label,
					ActionID:   core.ActionID{SpellID: 105872},
					Callback:   core.CallbackOnSpellHitDealt,
					Outcome:    core.OutcomeLanded,
					ProcMask:   core.ProcMaskMelee,
					Harmful:    true,
					ProcChance: 0.45,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						shaman.MaelstromWeaponAura.Activate(sim)
						shaman.MaelstromWeaponAura.AddStack(sim)
					},
				})
			}
		},
	},
})

func init() {
}
