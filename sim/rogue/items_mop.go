package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

var Tier14 = core.NewItemSet(core.ItemSet{
	Name:                    "Battlegear of the Thousandfold Blades",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the damage done by your Venomous Wounds ability by 20%,
			// increases the damage done by your Sinister Strike ability by 15%,
			// and increases the damage done by your Backstab ability by 10%.
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  RogueSpellVenomousWounds,
				FloatValue: 0.2,
			})
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  RogueSpellSinisterStrike,
				FloatValue: 0.15,
			})
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_DamageDone_Pct,
				ClassMask:  RogueSpellBackstab,
				FloatValue: 0.1,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the duration of your Shadow Blades ability by Combat: 6 /  Assassination,  Subtlety: 12 sec.
			rogue := agent.(RogueAgent).GetRogue()
			addTime := time.Second * time.Duration(core.Ternary(rogue.Spec == proto.Spec_SpecCombatRogue, 6, 12))
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_BuffDuration_Flat,
				ClassMask: RogueSpellShadowBlades,
				TimeValue: addTime,
			})
		},
	},
})

var Tier15 = core.NewItemSet(core.ItemSet{
	Name:                    "Nine-Tail Battlegear",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Increases the duration of your finishing moves as if you had used an additional combo point, up to a maximum of 6 combo points.
			rogue := agent.(RogueAgent).GetRogue()

			rogue.Has2PT15 = true
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Shadow Blades also reduces the cost of all your abilities by 15%.
			// Additionally, reduces the GCD of all rogue abilities by 300ms
			rogue := agent.(RogueAgent).GetRogue()
			energyMod := rogue.AddDynamicMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Pct,
				ClassMask:  RogueSpellsAll,
				FloatValue: -0.15,
			})
			gcdMod := rogue.AddDynamicMod(core.SpellModConfig{
				Kind:      core.SpellMod_GlobalCooldown_Flat,
				ClassMask: RogueSpellActives,
				TimeValue: time.Millisecond * -300,
			})
			aura := rogue.RegisterAura(core.Aura{
				Label:    "Shadow Blades Energy Cost Reduction",
				ActionID: core.ActionID{SpellID: 138151},
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					energyMod.Activate()
					gcdMod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					energyMod.Deactivate()
					gcdMod.Deactivate()
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Rogue T15 4P Bonus",
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: RogueSpellShadowBlades,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
				},
			})
		},
	},
})

var Tier16 = core.NewItemSet(core.ItemSet{
	Name:                    "Barbed Assassin Battlegear",
	DisabledInChallengeMode: true,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// When you generate a combo point from Revealing Strike's effect, Honor Among Thieves, or Seal Fate
			// your next combo point generating ability has its energy cost reduced by {Subtlety: 2, Assassination: 6, Combat: 15].
			// Stacks up to 5 times.
			rogue := agent.(RogueAgent).GetRogue()

			energyReduction := 0
			switch rogue.Spec {
			case proto.Spec_SpecSubtletyRogue:
				energyReduction = -2
			case proto.Spec_SpecAssassinationRogue:
				energyReduction = -6
			default:
				energyReduction = -15
			}

			energyMod := rogue.AddDynamicMod(core.SpellModConfig{
				Kind:      core.SpellMod_PowerCost_Flat,
				ClassMask: RogueSpellGenerator,
				IntValue:  0, // Set dynamically
			})

			// This aura gets activated by the applicable spell scripts
			rogue.T16EnergyAura = rogue.RegisterAura(core.Aura{
				Label:     "Silent Blades",
				ActionID:  core.ActionID{SpellID: 145193},
				Duration:  time.Second * 30,
				MaxStacks: 5,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					energyMod.UpdateIntValue(aura.GetStacks() * int32(energyReduction))
					energyMod.Activate()
				},
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					energyMod.UpdateIntValue(aura.GetStacks() * int32(energyReduction))
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					energyMod.Deactivate()
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.Flags.Matches(SpellFlagBuilder) && spell.DefaultCast.Cost > 0 {
						// Free action casts (such as Dispatch w/ Blindside) will not consume the aura
						aura.Deactivate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Killing Spree deals 10% more damage every time it strikes a target.
			// Abilities against a target with Vendetta on it increase your mastery by 250 for 5 sec, stacking up to 20 times.
			// Every time you Backstab, you have a 4% chance to replace your Backstab with Ambush that can be used regardless of Stealth.
			rogue := agent.(RogueAgent).GetRogue()

			if rogue.Spec == proto.Spec_SpecSubtletyRogue {
				aura := rogue.RegisterAura(core.Aura{
					Label:    "Sleight of Hand",
					ActionID: core.ActionID{SpellID: 145211},
					Duration: time.Second * 10,
					OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if spell == rogue.Ambush {
							aura.Deactivate(sim)
						}
					},
				})

				setBonusAura.AttachProcTrigger(core.ProcTrigger{
					Name:           "Rogue T16 4P Bonus",
					Callback:       core.CallbackOnApplyEffects,
					ClassSpellMask: RogueSpellBackstab,
					Outcome:        core.OutcomeLanded,
					ProcChance:     0.04,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						aura.Activate(sim)
					},
				})
			} else if rogue.Spec == proto.Spec_SpecCombatRogue {
				rogue.T16SpecMod = rogue.AddDynamicMod(core.SpellModConfig{
					Kind:       core.SpellMod_DamageDone_Pct,
					ClassMask:  RogueSpellKillingSpreeHit,
					FloatValue: 0.1, // Set dynamically in Killing Spree
				})
			} else if rogue.Spec == proto.Spec_SpecAssassinationRogue {
				aura := rogue.RegisterAura(core.Aura{
					Label:     "Toxicologist",
					ActionID:  core.ActionID{SpellID: 145249},
					Duration:  time.Second * 5,
					MaxStacks: 20,
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						aura.AddStack(sim)
						aura.Activate(sim)
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						aura.SetStacks(sim, 0)
					},
					OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
						change := newStacks - oldStacks
						aura.Unit.AddStatDynamic(sim, stats.MasteryRating, float64(250*change))
					},
				})

				setBonusAura.AttachProcTrigger(core.ProcTrigger{
					Name:           "Rogue T16 4P Bonus",
					Callback:       core.CallbackOnApplyEffects,
					ClassSpellMask: RogueSpellVendetta,
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						aura.Activate(sim)
					},
				})

				setBonusAura.AttachProcTrigger(core.ProcTrigger{
					Name:           "Toxicologist Trigger",
					Callback:       core.CallbackOnSpellHitDealt,
					ClassSpellMask: RogueSpellActives,
					ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
						return rogue.Vendetta.RelatedAuraArrays.AnyActive(result.Target)
					},
					Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						aura.AddStack(sim)
						aura.Refresh(sim)
					},
				})
			}
		},
	},
})
