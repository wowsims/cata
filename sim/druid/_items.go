package druid

import (
	"time"

	cata "github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// T11 Feral
var ItemSetStormridersBattlegarb = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Battlegarb",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Implemented in rake.go and lacerate.go
			druid := agent.(DruidAgent).GetDruid()
			druid.T11Feral2pBonus = setBonusAura
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			druid := agent.(DruidAgent).GetDruid()
			var apDepByStackCount = map[int32]*stats.StatDependency{}

			for i := 1; i <= 3; i++ {
				apDepByStackCount[int32(i)] = druid.NewDynamicMultiplyStat(stats.AttackPower, 1.0+0.01*float64(i))
			}

			druid.StrengthOfThePantherAura = druid.RegisterAura(core.Aura{
				Label:     "Strength of the Panther",
				ActionID:  core.ActionID{SpellID: 90166},
				Duration:  time.Second * 30,
				MaxStacks: 3,

				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
					if oldStacks > 0 {
						druid.DisableDynamicStatDep(sim, apDepByStackCount[oldStacks])
					}

					if newStacks > 0 {
						druid.EnableDynamicStatDep(sim, apDepByStackCount[newStacks])
					}
				},
			})

			druid.T11Feral4pBonus = setBonusAura
		},
	},
})

// T11 Balance
var ItemSetStormridersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		// Increases the critical strike chance of your Insect Swarm and Moonfire spells by 5%
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 5,
				ClassMask:  DruidSpellDoT | DruidSpellMoonfire | DruidSpellSunfire,
			})
		},
		// Whenever Eclipse triggers, your critical strike chance with spells is increased by 15% for 8 sec.  Each critical strike you achieve reduces that bonus by 5%
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			druid := agent.(DruidAgent).GetDruid()

			tierSet4pMod := druid.AddDynamicMod(core.SpellModConfig{
				School: core.SpellSchoolArcane | core.SpellSchoolNature,
				Kind:   core.SpellMod_BonusCrit_Percent,
			})

			tierSet4pAura := druid.RegisterAura(core.Aura{
				ActionID:  core.ActionID{SpellID: 90163},
				Label:     "Druid T11 Balance 4P Bonus",
				Duration:  time.Second * 8,
				MaxStacks: 3,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.SetStacks(sim, aura.MaxStacks)

					tierSet4pMod.UpdateFloatValue(float64(aura.GetStacks()) * 5)
					tierSet4pMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					tierSet4pMod.Deactivate()
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.DidCrit() && aura.GetStacks() > 0 {
						aura.RemoveStack(sim)
						tierSet4pMod.UpdateFloatValue(float64(aura.GetStacks()) * 5)
					}
				},
			})

			druid.AddEclipseCallback(func(_ Eclipse, gained bool, sim *core.Simulation) {
				if setBonusAura.IsActive() {
					if gained {
						tierSet4pAura.Activate(sim)
					} else {
						tierSet4pAura.Deactivate(sim)
					}
				}
			})
		},
	},
})

// T12 Feral
var ItemSetObsidianArborweaveBattlegarb = core.NewItemSet(core.ItemSet{
	Name: "Obsidian Arborweave Battlegarb",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// TODO: Verify behavior after PTR testing
			druid := agent.(DruidAgent).GetDruid()
			cata.RegisterIgniteEffect(&druid.Unit, cata.IgniteConfig{
				ActionID:         core.ActionID{SpellID: 99002},
				DotAuraLabel:     "Fiery Claws",
				IncludeAuraDelay: true,
				ParentAura:       setBonusAura,

				ProcTrigger: core.ProcTrigger{
					Name:           "Fiery Claws Trigger",
					Callback:       core.CallbackOnSpellHitDealt,
					ClassSpellMask: DruidSpellMangle | DruidSpellMaul | DruidSpellShred,
					Outcome:        core.OutcomeLanded,
				},

				DamageCalculator: func(result *core.SpellResult) float64 {
					return result.Damage * 0.1
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Full implementation in berserk.go and barkskin.go
			druid := agent.(DruidAgent).GetDruid()
			druid.T12Feral4pBonus = setBonusAura

			if !druid.InForm(Bear) {
				return
			}

			druid.SmokescreenAura = druid.RegisterAura(core.Aura{
				Label:    "Smokescreen",
				ActionID: core.ActionID{SpellID: 99011},
				Duration: time.Second * 12,

				OnGain: func(_ *core.Aura, _ *core.Simulation) {
					druid.PseudoStats.BaseDodgeChance += 0.1
				},

				OnExpire: func(_ *core.Aura, _ *core.Simulation) {
					druid.PseudoStats.BaseDodgeChance -= 0.1
				},
			})
		},
	},
})

// T12 Balance
var ItemSetObsidianArborweaveRegalia = core.NewItemSet(core.ItemSet{
	Name: "Obsidian Arborweave Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		// You have a chance to summon a Burning Treant to assist you in battle for 15 sec when you cast Wrath or Starfire. (Proc chance: 20%, 45s cooldown)
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			druid := agent.(DruidAgent).GetDruid()

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				ActionID:       core.ActionID{SpellID: 99019},
				Name:           "Item - Druid T12 Balance 2P Bonus",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: DruidSpellWrath | DruidSpellStarfire,
				ProcChance:     0.20,
				ICD:            time.Second * 45,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					druid.BurningTreant.EnableWithTimeout(sim, druid.BurningTreant, time.Second*15)
				},
			})
		},
		// While not in an Eclipse state, your Wrath generates 3 additional Lunar Energy and your Starfire generates 5 additional Solar Energy.
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			druid := agent.(DruidAgent).GetDruid()

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Custom,
				ClassMask: DruidSpellWrath,
				ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
					druid.SetSpellEclipseEnergy(DruidSpellWrath, WrathBaseEnergyGain, Wrath4PT12EnergyGain)
				},
				RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
					druid.SetSpellEclipseEnergy(DruidSpellWrath, WrathBaseEnergyGain, WrathBaseEnergyGain)
				},
			})

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Custom,
				ClassMask: DruidSpellStarfire,
				ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
					druid.SetSpellEclipseEnergy(DruidSpellStarfire, StarfireBaseEnergyGain, Starfire4PT12EnergyGain)
				},
				RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
					druid.SetSpellEclipseEnergy(DruidSpellStarfire, StarfireBaseEnergyGain, StarfireBaseEnergyGain)
				},
			})

			setBonusAura.ExposeToAPL(99049)
		},
	},
})

// T13 Feral
var ItemSetDeepEarthBattlegarb = core.NewItemSet(core.ItemSet{
	Name: "Deep Earth Battlegarb",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			druid := agent.(DruidAgent).GetDruid()

			if druid.InForm(Bear) {
				setBonusAura.AttachProcTrigger(core.ProcTrigger{
					Name:           "T13 Savage Defense Trigger",
					Callback:       core.CallbackOnSpellHitDealt,
					ClassSpellMask: DruidSpellMangleBear,
					Outcome:        core.OutcomeCrit,

					Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
						if druid.PulverizeAura.IsActive() {
							druid.SavageDefenseAura.Activate(sim)
						}
					},
				})
			}

			if !druid.InForm(Cat) {
				return
			}

			// Rather than creating a whole extra Execute phase category just for this bonus, we will instead scale up ExecuteProportion_25 using linear interpolation. Note that we use ExecuteProportion_90 for Predatory Strikes (< 80%), which is why the math below looks funny.
			oldExecuteProportion_25 := druid.Env.Encounter.ExecuteProportion_25
			oldExecuteProportion_35 := druid.Env.Encounter.ExecuteProportion_35
			newExecuteProportion_25 := oldExecuteProportion_35*(1.0-(60.0-35.0)/(80.0-35.0)) + druid.Env.Encounter.ExecuteProportion_90*((60.0-35.0)/(80.0-35.0))
			newExecuteProportion_35 := 0.5 * (newExecuteProportion_25 + druid.Env.Encounter.ExecuteProportion_90) // We don't use this field anywhere, just need it to be any value above ExecuteProportion_25 but below ExecuteProportion_90 so that the transitions work properly.

			setBonusAura.ApplyOnGain(func(_ *core.Aura, _ *core.Simulation) {
				druid.Env.Encounter.ExecuteProportion_35 = newExecuteProportion_35
				druid.Env.Encounter.ExecuteProportion_25 = newExecuteProportion_25
			})

			setBonusAura.ApplyOnExpire(func(_ *core.Aura, _ *core.Simulation) {
				druid.Env.Encounter.ExecuteProportion_25 = oldExecuteProportion_25
				druid.Env.Encounter.ExecuteProportion_35 = oldExecuteProportion_35
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Implemented in tigers_fury.go
			druid := agent.(DruidAgent).GetDruid()
			druid.T13Feral4pBonus = setBonusAura
		},
	},
})

// T13 Balance
var ItemSetDeepEarthRegalia = core.NewItemSet(core.ItemSet{
	Name: "Deep Earth Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		// Insect Swarm increases all damage done by your Starfire, Starsurge, and Wrath spells against that target by 3%
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			druid := agent.(DruidAgent).GetDruid()

			t13InsectSwarmBonus := func(_ *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
				if spell.Matches(DruidSpellStarsurge | DruidSpellStarfire | DruidSpellWrath) {
					return 1.03
				}

				return 1.0
			}

			t13InsectSwarmBonusDummyAuras := druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
				return target.GetOrRegisterAura(core.Aura{
					ActionID: core.ActionID{SpellID: 105722},
					Label:    "Item - Druid T13 Balance 2P Bonus (Insect Swarm) - " + druid.Label,
					Duration: core.NeverExpires,
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						druid.AttackTables[aura.Unit.UnitIndex].DamageDoneByCasterMultiplier = t13InsectSwarmBonus
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						druid.AttackTables[aura.Unit.UnitIndex].DamageDoneByCasterMultiplier = nil
					},
				})
			})

			druid.OnSpellRegistered(func(spell *core.Spell) {
				if !spell.Matches(DruidSpellInsectSwarm) {
					return
				}

				for _, target := range druid.Env.Encounter.TargetUnits {
					spell.Dot(target).ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
						if setBonusAura.IsActive() {
							t13InsectSwarmBonusDummyAuras.Get(aura.Unit).Activate(sim)
						}
					})

					spell.Dot(target).ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
						t13InsectSwarmBonusDummyAuras.Get(aura.Unit).Deactivate(sim)
					})
				}
			})
		},
		// Reduces the cooldown of Starsurge by 5 sec and increases its damage by 10%
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.1,
				ClassMask:  DruidSpellStarsurge,
			})

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Cooldown_Flat,
				TimeValue: time.Second * -5,
				ClassMask: DruidSpellStarsurge,
			})
		},
	},
})

// PvP Feral
var ItemSetGladiatorsSanctuary = core.NewItemSet(core.ItemSet{
	ID:   922,
	Name: "Gladiator's Sanctuary",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Agility, 70)

		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Agility, 90)
		},
	},
})

func init() {
}
