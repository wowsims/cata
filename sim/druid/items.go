package druid

import (
	"time"

	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

// T11 Feral
var ItemSetStormridersBattlegarb = core.NewItemSet(core.ItemSet{
	Name: "Stormrider's Battlegarb",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Implemented in rake.go and lacerate.go
			druid := agent.(DruidAgent).GetDruid()
			setBonusAura.AttachBooleanToggle(druid.HasT11Feral2pBonus)
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

			setBonusAura.AttachBooleanToggle(druid.HasT11Feral4pBonus)
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
				SetBonusAura:     setBonusAura,

				ProcTrigger: core.ProcTrigger{
					Name:           "Fiery Claws Trigger",
					Callback:       core.CallbackOnSpellHitDealt,
					ClassSpellMask: DruidSpellMangle | DruidSpellMaul | DruidSpellShred,
					Outcome:        core.OutcomeLanded,
					ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
						return setBonusAura.IsActive()
					},
				},

				DamageCalculator: func(result *core.SpellResult) float64 {
					return result.Damage * 0.1
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Full implementation in berserk.go and barkskin.go
			druid := agent.(DruidAgent).GetDruid()

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
			setBonusAura.AttachBooleanToggle(druid.HasT12Feral4pBonus)
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

			druid.OnSpellRegistered(func(spell *core.Spell) {
				if spell.Matches(DruidSpellWrath) {
					onEquip := func() {
						druid.SetSpellEclipseEnergy(DruidSpellWrath, WrathBaseEnergyGain, Wrath4PT12EnergyGain)
					}

					setBonusAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
						onEquip()
					})

					setBonusAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
						druid.SetSpellEclipseEnergy(DruidSpellWrath, WrathBaseEnergyGain, WrathBaseEnergyGain)
					})

					if setBonusAura.IsActive() {
						onEquip()
					}

				}

				if spell.Matches(DruidSpellStarfire) {
					onEquip := func() {
						druid.SetSpellEclipseEnergy(DruidSpellStarfire, StarfireBaseEnergyGain, Starfire4PT12EnergyGain)
					}

					setBonusAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
						onEquip()
					})

					setBonusAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
						druid.SetSpellEclipseEnergy(DruidSpellStarfire, StarfireBaseEnergyGain, StarfireBaseEnergyGain)
					})

					if setBonusAura.IsActive() {
						onEquip()
					}
				}
			})

			setBonusAura.ExposeToAPL(99049)
		},
	},
})

// T13 Balance
var ItemSetDeepEarthRegalia = core.NewItemSet(core.ItemSet{
	Name: "Deep Earth Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		// Insect Swarm increases all damage done by your Starfire, Starsurge, and Wrath spells against that target by 3%
		2: func(_ core.Agent, _ *core.Aura) {
		},
		// Reduces the cooldown of Starsurge by 5 sec and increases its damage by 10%
		4: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.05,
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
