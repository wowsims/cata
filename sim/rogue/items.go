package rogue

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var Tier11 = core.NewItemSet(core.ItemSet{
	Name: "Wind Dancer's Regalia",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			// +5% Crit to Backstab, Mutilate, and Sinister Strike
			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_BonusCrit_Percent,
				FloatValue: 5,
				ClassMask:  RogueSpellBackstab | RogueSpellMutilate | RogueSpellSinisterStrike,
			})
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// 1% Chance on Auto Attack to increase crit of next Evis or Envenom by +100% for 15 seconds
			rogue := agent.(RogueAgent).GetRogue()

			t11Proc := rogue.RegisterAura(core.Aura{
				Label:    "Deadly Scheme Proc",
				ActionID: core.ActionID{SpellID: 90472},
				Duration: time.Second * 15,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.Envenom.BonusCritPercent += 100
					rogue.Eviscerate.BonusCritPercent += 100
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.Envenom.BonusCritPercent -= 100
					rogue.Eviscerate.BonusCritPercent -= 100
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == rogue.Envenom || spell == rogue.Eviscerate {
						aura.Deactivate(sim)
					}
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Deadly Scheme Aura",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeWhiteHit,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.01,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					t11Proc.Activate(sim)
				},
			})
		},
	},
})

func MakeT12StatAura(action core.ActionID, stat stats.Stat, name string) core.Aura {
	var lastRatingGain float64
	return core.Aura{
		Label:    name,
		ActionID: action,
		Duration: 30 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			lastRatingGain = aura.Unit.GetStat(stat) * 0.25
			aura.Unit.AddStatDynamic(sim, stat, lastRatingGain)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.AddStatDynamic(sim, stat, -lastRatingGain)
		},
	}
}

var Tier12 = core.NewItemSet(core.ItemSet{
	Name: "Vestments of the Dark Phoenix",
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your melee critical strikes deal 6% additional damage as Fire over 4 sec.
			// Rolls like ignite
			// Tentatively, this is just Ignite. Testing required to validate behavior.
			rogue := agent.GetCharacter()

			cata.RegisterIgniteEffect(&rogue.Unit, cata.IgniteConfig{
				ActionID:         core.ActionID{SpellID: 99173},
				DotAuraLabel:     "Burning Wounds",
				IncludeAuraDelay: true,
				SetBonusAura:     setBonusAura,

				ProcTrigger: core.ProcTrigger{
					Name:     "Rogue T12 2P Bonus",
					Callback: core.CallbackOnSpellHitDealt,
					ProcMask: core.ProcMaskMelee,
					Outcome:  core.OutcomeCrit,
				},

				DamageCalculator: func(result *core.SpellResult) float64 {
					return result.Damage * .06
				},
			})

		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			// Your Tricks of the Trade ability also causes you to gain a 25% increase to Haste, Mastery, or Critical Strike chosen at random for 30 sec.
			// Cannot pick the same stat twice in a row. No other logic appears to exist
			// Not a dynamic 1.25% mod; snapshots stats and applies that much as bonus rating for duration
			// Links to all buffs: https://www.wowhead.com/spell=99175/item-rogue-t12-4p-bonus#comments:id=1507073
			rogue := agent.(RogueAgent).GetRogue()

			// Aura for adding 25% of current rating as extra rating
			hasteAura := rogue.GetOrRegisterAura(MakeT12StatAura(core.ActionID{SpellID: 99186}, stats.HasteRating, "Future on Fire"))
			critAura := rogue.GetOrRegisterAura(MakeT12StatAura(core.ActionID{SpellID: 99187}, stats.CritRating, "Fiery Devastation"))
			mastAura := rogue.GetOrRegisterAura(MakeT12StatAura(core.ActionID{SpellID: 99188}, stats.MasteryRating, "Master of Flames"))
			auraArray := [3]*core.Aura{hasteAura, critAura, mastAura}

			// Proc aura watching for ToT threat transfer start
			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Rogue T12 4P Bonus",
				Callback:       core.CallbackOnApplyEffects,
				ClassSpellMask: RogueSpellTricksOfTheTrade,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if rogue.T12ToTLastBuff == 3 { // any of 3
						randomStat := int(math.Mod(sim.RandomFloat("Rogue T12 4P Bonus Initial")*10, 3))
						rogue.T12ToTLastBuff = randomStat
						auraArray[rogue.T12ToTLastBuff].Activate(sim)
					} else { // cannot re-roll same
						randomStat := int(math.Mod(sim.RandomFloat("Rogue T12 4P Bonus")*10, 1)) + 1
						rogue.T12ToTLastBuff = (rogue.T12ToTLastBuff + randomStat) % 3
						auraArray[rogue.T12ToTLastBuff].Activate(sim)
					}
				},
			})
		},
	},
})

var Tier13 = core.NewItemSet(core.ItemSet{
	Name: "Blackfang Battleweave",
	Bonuses: map[int32]core.ApplySetBonus{
		// After triggering Tricks of the Trade, your abilities cost 20% less energy for 6 sec.
		// This is implemented as it is because the 20% reduction is applied -before- talents/glyphs/passives, which is not how SpellMod_PowerCost_Pct operates
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			rogue := agent.(RogueAgent).GetRogue()

			aura := rogue.GetOrRegisterAura(core.Aura{
				Label:    "Tricks of Time",
				ActionID: core.ActionID{SpellID: 105864},
				Duration: time.Second * 6,
			})

			// Bonus 60 Enmergy
			aura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -12,
				ClassMask:  RogueSpellAmbush | RogueSpellBackstab | RogueSpellMutilate,
			})

			// Bonus 45 Enmergy
			aura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -9,
				ClassMask:  RogueSpellSinisterStrike | RogueSpellGouge | RogueSpellGarrote,
			})

			// Bonus 40 Enmergy
			aura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -8,
				ClassMask:  RogueSpellRevealingStrike,
			})

			// Bonus 35 Enmergy
			aura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -7,
				ClassMask:  RogueSpellEviscerate | RogueSpellEnvenom | RogueSpellHemorrhage,
			})

			// Bonus 30 Enmergy
			aura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -6,
				ClassMask:  RogueSpellRecuperate,
			})

			// Bonus 25 Enmergy
			aura.AttachSpellMod(core.SpellModConfig{
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -5,
				ClassMask:  RogueSpellSliceAndDice | RogueSpellRupture,
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:           "Rogue T13 2P Bonus",
				Callback:       core.CallbackOnApplyEffects,
				ClassSpellMask: RogueSpellTricksOfTheTrade,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					aura.Activate(sim)
				},
			})
		},
		// Increases the duration of Shadow Dance by 2 sec, Adrenaline Rush by 3 sec, and Vendetta by 9 sec.
		// Implemented in respective spells
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			rogue := agent.(RogueAgent).GetRogue()

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Custom,
				ClassMask: RogueSpellVendetta,
				ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
					for _, aura := range spell.RelatedAuras[0] {
						if aura != nil {
							aura.Duration = rogue.getVendettaDuration(9)
						}
					}
				},
				RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
					for _, aura := range spell.RelatedAuras[0] {
						if aura != nil {
							aura.Duration = rogue.getVendettaDuration(0)
						}
					}
				},
			})

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Custom,
				ClassMask: RogueSpellAdrenalineRush,
				ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
					rogue.AdrenalineRushAura.Duration += time.Second * 3
				},
				RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
					rogue.AdrenalineRushAura.Duration -= time.Second * 3
				},
			})

			setBonusAura.AttachSpellMod(core.SpellModConfig{
				Kind:      core.SpellMod_Custom,
				ClassMask: RogueSpellShadowDance,
				ApplyCustom: func(mod *core.SpellMod, spell *core.Spell) {
					rogue.ShadowDanceAura.Duration += time.Second * 2
				},
				RemoveCustom: func(mod *core.SpellMod, spell *core.Spell) {
					rogue.ShadowDanceAura.Duration -= time.Second * 2
				},
			})
		},
	},
})

// Pulled from old Shadowcraft/SimC logic.
// There exists Blizzard sourced numbers, but those were from MoP beta. TBD which is valid.
// The final difference between the Blizzard numbers and old TC numbers is exceedingly small either way.
func getFangsProcRate(character *core.Character) float64 {
	switch character.Spec {
	case proto.Spec_SpecSubtletyRogue:
		return 0.275
	case proto.Spec_SpecAssassinationRogue:
		return 0.235
	default:
		return 0.095
	}
}

// Fear + Vengeance
var JawsOfRetribution = core.NewItemSet(core.ItemSet{
	Name:  "Jaws of Retribution",
	Slots: []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand},
	Bonuses: map[int32]core.ApplySetBonus{
		// Your melee attacks have a chance to grant Suffering, increasing your Agility by 2, stacking up to 50 times.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()

			agiAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Suffering",
					ActionID:  core.ActionID{SpellID: 109959},
					Duration:  time.Second * 30,
					MaxStacks: 50,
				},
				BonusPerStack: stats.Stats{stats.Agility: 2},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Rogue Legendary Daggers Stage 1",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				Outcome:    core.OutcomeLanded,
				ProcChance: getFangsProcRate(character),
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					agiAura.Activate(sim)
					agiAura.AddStack(sim)
				},
			})
		},
	},
})

// Sleeper + Dreamer
var MawOfOblivion = core.NewItemSet(core.ItemSet{
	Name:  "Maw of Oblivion",
	Slots: []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand},
	Bonuses: map[int32]core.ApplySetBonus{
		// Your melee attacks have a chance to grant Nightmare, increasing your Agility by 5, stacking up to 50 times.
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()

			agiAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Nightmare",
					ActionID:  core.ActionID{SpellID: 109955},
					Duration:  time.Second * 30,
					MaxStacks: 50,
				},
				BonusPerStack: stats.Stats{stats.Agility: 5},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Rogue Legendary Daggers Stage 2",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				Outcome:    core.OutcomeLanded,
				ProcChance: getFangsProcRate(agent.GetCharacter()),
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					agiAura.Activate(sim)
					agiAura.AddStack(sim)
				},
			})
		},
	},
})

// Golad + Tiriosh
var FangsOfTheFather = core.NewItemSet(core.ItemSet{
	Name:  "Fangs of the Father",
	Slots: []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand},
	Bonuses: map[int32]core.ApplySetBonus{
		// Your melee attacks have a chance to grant Shadows of the Destroyer, increasing your Agility by 17, stacking up to 50 times.
		// Each application past 30 grants an increasing chance to trigger Fury of the Destroyer.
		// When triggered, this consumes all applications of Shadows of the Destroyer, immediately granting 5 combo points and cause your finishing moves to generate 5 combo points.
		// Lasts 6 sec.

		// Tooltip is deceptive. The stacks of Shadows of the Destroyer only clear when the 5 Combo Point effect ends
		2: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()
			cpMetrics := character.NewComboPointMetrics(core.ActionID{SpellID: 109950})

			agiAura := core.MakeStackingAura(character, core.StackingStatAura{
				Aura: core.Aura{
					Label:     "Shadows of the Destroyer",
					ActionID:  core.ActionID{SpellID: 109941},
					Duration:  time.Second * 30,
					MaxStacks: 50,
				},
				BonusPerStack: stats.Stats{stats.Agility: 17},
			})

			wingsProc := character.GetOrRegisterAura(core.Aura{
				Label:    "Fury of the Destroyer",
				ActionID: core.ActionID{SpellID: 109949},
				Duration: time.Second * 6,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddComboPoints(sim, 5, cpMetrics)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.Flags.Matches(SpellFlagFinisher) {
						aura.Unit.AddComboPoints(sim, 5, cpMetrics)
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					agiAura.SetStacks(sim, 0)
					agiAura.Deactivate(sim)
				},
			})

			setBonusAura.AttachProcTrigger(core.ProcTrigger{
				Name:       "Rogue Legendary Daggers Stage 3",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				Outcome:    core.OutcomeLanded,
				ProcChance: getFangsProcRate(character),
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					// Adding a stack and activating the combo point effect is mutually exclusive.
					// Agility bonus is lost when combo point effect ends
					stacks := float64(agiAura.GetStacks())
					if stacks > 30 && !wingsProc.IsActive() {
						if stacks == 50 || sim.Proc(1.0/(50-stacks), "Fangs of the Father") {
							wingsProc.Activate(sim)
						} else {
							agiAura.Activate(sim)
							agiAura.AddStack(sim)
						}
					} else {
						agiAura.Activate(sim)
						agiAura.AddStack(sim)
					}
				},
			})
		},
	},
})

var CataPVPSet = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Vestments",
	ID:   914,
	Bonuses: map[int32]core.ApplySetBonus{
		2: func(_ core.Agent, setBonusAura *core.Aura) {
			setBonusAura.AttachStatBuff(stats.Agility, 70)
		},
		4: func(agent core.Agent, setBonusAura *core.Aura) {
			character := agent.GetCharacter()
			setBonusAura.AttachStatBuff(stats.Agility, 90)

			actionID := core.ActionID{SpellID: 21975}
			energyMetrics := character.NewEnergyMetrics(actionID)

			setBonusAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
				character.UpdateMaxEnergy(sim, 10, energyMetrics)
			})
			setBonusAura.ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
				character.UpdateMaxEnergy(sim, -10, energyMetrics)
			})
			setBonusAura.ExposeToAPL(actionID.SpellID)
		},
	},
})

// 45% SS/RvS Modifier for Legendary MH Dagger
func makeWeightedBladesModifier(itemID int32) {
	core.NewItemEffect(itemID, func(agent core.Agent) {
		agent.GetCharacter().AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			FloatValue: 0.45,
			ClassMask:  RogueSpellWeightedBlades,
		})
	})
}

func init() {
	makeWeightedBladesModifier(77945)
	makeWeightedBladesModifier(77947)
	makeWeightedBladesModifier(77949)
}
