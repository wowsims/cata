package priest

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

var ItemSetVestmentsOfAbsolution = core.NewItemSet(core.ItemSet{
	Name: "Vestments of Absolution",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				FloatValue: -0.1,
				ClassMask:  PriestSpellPrayerOfHealing,
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.05,
				ClassMask:  PriestSpellGreaterHeal,
			})
		},
	},
})

var ItemSetValorous = core.NewItemSet(core.ItemSet{
	Name: "Garb of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				FloatValue: -0.1,
				ClassMask:  PriestSpellMindBlast,
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				FloatValue: 10 * core.CritRatingPerCritChance,
				ClassMask:  PriestSpellShadowWordDeath,
			})
		},
	},
})

var ItemSetRegaliaOfFaith = core.NewItemSet(core.ItemSet{
	Name: "Regalia of Faith",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Not implemented
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				FloatValue: -0.05,
				ClassMask:  PriestSpellGreaterHeal,
			})
		},
	},
})

var ItemSetConquerorSanct = core.NewItemSet(core.ItemSet{
	Name: "Sanctification Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.15,
				ClassMask:  PriestSpellDevouringPlague,
			})
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			procAura := priest.NewTemporaryStatsAura("Devious Mind", core.ActionID{SpellID: 64907}, stats.Stats{stats.SpellHaste: 240}, time.Second*4)

			priest.RegisterAura(core.Aura{
				Label:    "Devious Mind Proc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				// TODO: Does this affect the spell that procs it?
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ClassSpellMask == PriestSpellMindBlast {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetSanctificationRegalia = core.NewItemSet(core.ItemSet{
	Name: "Sanctification Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				FloatValue: 10 * core.CritRatingPerCritChance,
				ClassMask:  PriestSpellPrayerOfHealing,
			})
		},
		4: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()
			procAura := priest.NewTemporaryStatsAura("Sanctification Reglia 4pc", core.ActionID{SpellID: 64912}, stats.Stats{stats.SpellPower: 250}, time.Second*5)

			priest.RegisterAura(core.Aura{
				Label:    "Sancitifcation Reglia 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				// TODO: Does this affect the spell that procs it?
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == priest.PowerWordShield {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetZabras = core.NewItemSet(core.ItemSet{
	Name:            "Zabra's Regalia",
	AlternativeName: "Velen's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Modifies dot length, need to implement later again
			// Requieres tests and proper modification of SpellMods
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				FloatValue: 5 * core.CritRatingPerCritChance,
				ClassMask:  PriestSpellMindFlay,
			})
		},
	},
})

var ItemSetZabrasRaiment = core.NewItemSet(core.ItemSet{
	Name:            "Zabra's Raiment",
	AlternativeName: "Velen's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.15,
				ClassMask:  PriestSpellPrayerOfMending,
			})
		},
		4: func(agent core.Agent) {
			// changed in cata to flat 5% heal
			character := agent.GetCharacter()
			character.PseudoStats.DamageDealtMultiplier *= 1.05
		},
	},
})

var ItemSetCrimsonAcolyte = core.NewItemSet(core.ItemSet{
	Name: "Crimson Acolyte's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				FloatValue: 5 * core.CritRatingPerCritChance,
				ClassMask:  PriestSpellShadowWordPain | PriestSpellDevouringPlague | PriestSpellVampiricTouch | PriestSpellImprovedDevouringPlague,
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:      core.SpellMod_DotTickLength_Flat,
				TimeValue: -time.Millisecond * 170,
				ClassMask: PriestSpellMindFlay,
			})
		},
	},
})

var ItemSetCrimsonAcolytesRaiment = core.NewItemSet(core.ItemSet{
	Name: "Crimson Acolyte's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			priest := agent.(PriestAgent).GetPriest()

			var curAmount float64
			procSpell := priest.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 70770},
				SpellSchool: core.SpellSchoolHoly,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagHelpful,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Hot: core.DotConfig{
					Aura: core.Aura{
						Label: "CrimsonAcolyteRaiment2pc",
					},
					NumberOfTicks: 3,
					TickLength:    time.Second * 3,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
						dot.SnapshotBaseDamage = curAmount * 0.33
						dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
					},
				},
			})

			priest.RegisterAura(core.Aura{
				Label:    "Crimson Acolytes Raiment 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ClassSpellMask != PriestSpellFlashHeal || !sim.Proc(0.33, "Crimson Acolytes Raiment 2pc") {
						return
					}

					curAmount = result.Damage
					hot := procSpell.Hot(result.Target)
					hot.Apply(sim)
				},
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.05,
				ClassMask:  PriestSpellPowerWordShield,
			})

			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				FloatValue: 0.10,
				ClassMask:  PriestSpellCircleOfHealing,
			})
		},
	},
})

var ItemSetGladiatorsInvestiture = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Investiture",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Resilience, 400)
			agent.GetCharacter().AddStat(stats.Intellect, 70)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Intellect, 90)
		},
	},
})

var ItemSetGladiatorsRaiment = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Resilience, 400)
			agent.GetCharacter().AddStat(stats.Intellect, 70)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Intellect, 90)
		},
	},
})

// T11 - Shadow
var ItemSetMercurialRegalia = core.NewItemSet(core.ItemSet{
	Name: "Mercurial Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Rating,
				FloatValue: 5 * core.CritRatingPerCritChance,
				ClassMask:  PriestSpellMindFlay | PriestSpellMindSear,
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			character.AddStaticMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.3,
				ClassMask:  PriestSpellShadowyApparation,
			})
		},
	},
})

func init() {
}
