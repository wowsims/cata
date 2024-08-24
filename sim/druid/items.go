package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// T11 Feral
var ItemSetStormridersBattlegarb = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Battlegarb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in rake.go and lacerate.go
		},
		4: func(agent core.Agent) {
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
		},
	},
})

// T11 Balance
var ItemSetStormridersRegalia = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases the critical strike chance of your Insect Swarm and Moonfire spells by 5%
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 5,
				ClassMask:  DruidSpellDoT | DruidSpellMoonfire | DruidSpellSunfire,
			})
		},
		// Whenever Eclipse triggers, your critical strike chance with spells is increased by 15% for 8 sec.  Each critical strike you achieve reduces that bonus by 5%
		4: func(agent core.Agent) {
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
				if gained {
					tierSet4pAura.Activate(sim)
				} else {
					tierSet4pAura.Deactivate(sim)
				}
			})
		},
	},
})

// T12 Balance
var ItemSetObsidianArborweaveRegalia = core.NewItemSet(core.ItemSet{
	Name: "Obsidian Arborweave Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		// You have a chance to summon a Burning Treant to assist you in battle for 15 sec when you cast Wrath or Starfire. (Proc chance: 20%, 45s cooldown)
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
				ActionID:       core.ActionID{SpellID: 99019},
				Name:           "Item - Druid T12 2P Bonus",
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
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			druid.OnSpellRegistered(func(spell *core.Spell) {
				if spell.ClassSpellMask == DruidSpellWrath {
					druid.SetSpellEclipseEnergy(spell.SpellID, WrathBaseEnergyGain, WrathBaseEnergyGain+3)
				}

				if spell.ClassSpellMask == DruidSpellStarfire {
					druid.SetSpellEclipseEnergy(spell.SpellID, StarfireBaseEnergyGain, StarfireBaseEnergyGain+5)
				}
			})
		},
	},
})

// PvP Feral
var ItemSetGladiatorsSanctuary = core.NewItemSet(core.ItemSet{
	ID:   922,
	Name: "Gladiator's Sanctuary",

	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.AddStats(stats.Stats{
				stats.Agility: 70,
			})
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.AddStats(stats.Stats{
				stats.Agility: 90,
			})
		},
	},
})

func init() {
}
