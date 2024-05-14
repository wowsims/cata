package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

var ItemSetEarthenWarplate = core.NewItemSet(core.ItemSet{
	Name: "Earthen Warplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				ClassMask:  SpellMaskBloodthirst | SpellMaskMortalStrike,
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.05,
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 90294}

			apDep := make([]*stats.StatDependency, 3)
			for i := 1; i <= 3; i++ {
				apDep[i-1] = character.NewDynamicMultiplyStat(stats.AttackPower, 1.0+float64(i)*0.01)
			}

			buff := character.RegisterAura(core.Aura{
				Label:     "Rage of the Ages",
				ActionID:  actionID,
				Duration:  30 * time.Second,
				MaxStacks: 3,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					// Example from DK Death Eater
					if oldStacks > 0 {
						character.DisableDynamicStatDep(sim, apDep[oldStacks-1])
					}
					if newStacks > 0 {
						character.EnableDynamicStatDep(sim, apDep[newStacks-1])
					}
				},
			})

			core.MakeProcTriggerAura(&agent.GetCharacter().Unit, core.ProcTrigger{
				Name:           "Rage of the Ages Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskOverpower | SpellMaskRagingBlow,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					buff.Activate(sim)
					buff.AddStack(sim)
				},
			})
		},
	},
})

var ItemSetEarthenBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Earthen Battleplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				ClassMask:  SpellMaskShieldSlam,
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.05,
			})
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStaticMod(core.SpellModConfig{
				ClassMask:  SpellMaskShieldWall,
				Kind:       core.SpellMod_Cooldown_Multiplier,
				FloatValue: -0.5,
			})
		},
	},
})

var ItemSetMoltenGiantWarplate = core.NewItemSet(core.ItemSet{
	Name: "Molten Giant Warplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 99233}
			buff := character.RegisterAura(core.Aura{
				Label:    "Burning Rage",
				ActionID: actionID,
				Duration: 12 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.1
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:           "Burning Rage Trigger",
				ActionID:       actionID,
				ClassSpellMask: SpellMaskColossusSmash,
				Callback:       core.CallbackOnSpellHitDealt,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					buff.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 99237}

			fieryAttack := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty, // TODO (4.2) Test this
				Flags:       core.SpellFlagMeleeMetrics,

				CritMultiplier:   character.DefaultMeleeCritMultiplier(),
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicCrit) // TODO (4.1) Test hit table
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:           "Fiery Attack Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskMortalStrike | SpellMaskRagingBlow,
				ProcChance:     0.3,
				Outcome:        core.OutcomeLanded,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					fieryAttack.Cast(sim, result.Target)
				},
			})
		},
	},
})

var ItemSetMoltenGiantBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Molten Giant Battleplate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 99240}

			// TODO (4.2): Test if this rolls damage over like deep wounds or just resets it
			var shieldSlamDamage float64 = 0.0
			debuff := character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,

				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				Dot: core.DotConfig{
					Aura: core.Aura{
						Label: "Combust",
					},
					NumberOfTicks: 2,
					TickLength:    2 * time.Second,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
						dot.Snapshot(target, shieldSlamDamage/float64(dot.NumberOfTicks))
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
					},
				},

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					spell.Dot(target).Apply(sim)
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:           "Combust Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskShieldSlam,
				Outcome:        core.OutcomeLanded,

				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					shieldSlamDamage = result.Damage * 0.2
					debuff.Cast(sim, result.Target)
				},
			})
		},
		4: func(agent core.Agent) {
			character := agent.GetCharacter()

			character.RegisterAura(core.Aura{
				Label:    "T12 4P Bonus",
				ActionID: core.ActionID{SpellID: 99242},
				Duration: 10 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					character.AddStatDynamic(sim, stats.Parry, 6*core.ParryRatingPerParryChance)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					character.AddStatDynamic(sim, stats.Parry, -6*core.ParryRatingPerParryChance)
				},
			})

		},
	},
})

var ItemSetColossalDragonplateBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Colossal Dragonplate Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()

			mod := character.AddDynamicMod(core.SpellModConfig{
				ClassMask:  SpellMaskHeroicStrike,
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -10,
			})

			actionID := core.ActionID{SpellID: 105860}
			buffAura := character.RegisterAura(core.Aura{
				Label:    "Volatile Outrage",
				ActionID: actionID,
				Duration: 15 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mod.Deactivate()
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:           "Volatile Outrage Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskInnerRage,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					buffAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()

			actionID := core.ActionID{SpellID: 108126}
			procCS := warrior.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				Flags:       core.SpellFlagNoOnDamageDealt | core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					DefaultCast: core.Cast{
						GCD: 0,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					warrior.ColossusSmashAuras.Get(target).Activate(sim)
				},
			})

			// TODO (4.3): Check if this cares that the hit landed
			core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
				Name:           "Warrior T13 4P Bloodthirst Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskBloodthirst,
				ProcChance:     0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procCS.Cast(sim, result.Target)
				},
			})

			core.MakeProcTriggerAura(&warrior.Unit, core.ProcTrigger{
				Name:           "Warrior T13 4P Mortal Strike Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskMortalStrike,
				ProcChance:     0.13,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procCS.Cast(sim, result.Target)
				},
			})
		},
	},
})

var ItemSetColossalDragonplateArmor = core.NewItemSet(core.ItemSet{
	Name: "Colossal Dragonplate Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			actionID := core.ActionID{SpellID: 105909}

			//var shieldAmt float64 = 0.0
			shieldAura := character.RegisterAura(core.Aura{
				Label:    "Shield of Fury",
				ActionID: actionID,
				Duration: 6 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					// TODO: Shield mechanics NYI
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					// TODO: Shield mechanics NYI
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:           "Shield of Fury Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskRevenge,
				Outcome:        core.OutcomeLanded,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					// shieldAmt = result.Damage * 0.2
					shieldAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent) {
			// TODO: Implement this, turns Shield Wall into a raid buff
		},
	},
})

func init() {
}
