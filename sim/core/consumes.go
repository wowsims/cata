package core

import (
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/dbc"
)

// Registers all consume-related effects to the Agent.
func applyConsumeEffects(agent Agent) {
	character := agent.GetCharacter()
	consumes := character.Consumes
	if consumes == nil {
		return
	}
	alchemyFlaskBonus := TernaryFloat64(character.HasProfession(proto.Profession_Alchemy), 80, 0)
	alchemyBattleElixirBonus := TernaryFloat64(character.HasProfession(proto.Profession_Alchemy), 40, 0)
	if consumes.FlaskId != 0 {
		flask := ConsumableByID[consumes.FlaskId]
		if flask.Stats[stats.Strength] > 0 {
			flask.Stats[stats.Strength] += alchemyFlaskBonus
		} else if flask.Stats[stats.Agility] > 0 {
			flask.Stats[stats.Agility] += alchemyFlaskBonus
		} else if flask.Stats[stats.Intellect] > 0 {
			flask.Stats[stats.Intellect] += alchemyFlaskBonus
		} else if flask.Stats[stats.Spirit] > 0 {
			flask.Stats[stats.Spirit] += alchemyFlaskBonus
		} else if flask.Stats[stats.Stamina] > 0 {
			flask.Stats[stats.Stamina] += alchemyFlaskBonus * 1.5
		}
		character.AddStats(stats.FromProtoArray(flask.Stats))
	}

	if consumes.BattleElixirId != 0 {
		elixir := ConsumableByID[consumes.BattleElixirId]
		if elixir.Stats[stats.MasteryRating] > 0 {
			elixir.Stats[stats.MasteryRating] += alchemyBattleElixirBonus
		} else if elixir.Stats[stats.HasteRating] > 0 {
			elixir.Stats[stats.HasteRating] += alchemyBattleElixirBonus
		} else if elixir.Stats[stats.CritRating] > 0 {
			elixir.Stats[stats.CritRating] += alchemyBattleElixirBonus
		} else if elixir.Stats[stats.ExpertiseRating] > 0 {
			elixir.Stats[stats.ExpertiseRating] += alchemyBattleElixirBonus
		} else if elixir.Stats[stats.Spirit] > 0 {
			elixir.Stats[stats.Spirit] += alchemyBattleElixirBonus
		}
		character.AddStats(stats.FromProtoArray(elixir.Stats))
	}

	if consumes.GuardianElixirId != 0 {
		elixir := ConsumableByID[consumes.GuardianElixirId]
		if elixir.Stats[stats.Armor] > 0 {
			elixir.Stats[stats.Armor] += 280
		}
		character.AddStats(stats.FromProtoArray(elixir.Stats))
	}

	if consumes.FoodId != 0 {
		// This is an example of how we could do overrides
		if consumes.FoodId == 62290 || consumes.FoodId == 62649 {
			character.AddStat(stats.Stamina, 90)
			character.AddStat(character.GetHighestStatType([]stats.Stat{stats.Strength, stats.Agility, stats.Intellect}), 90)
		}
		food := ConsumableByID[consumes.FoodId]
		character.AddStats(stats.FromProtoArray(food.Stats))
	}
	// Todo: Implement this?
	// case proto.BattleElixir_ElixirOfDemonslaying:
	// 	if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
	// 		character.PseudoStats.MobTypeAttackPower += 265
	// 	}
	// }

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
	registerTinkerHandsCD(agent, consumes)
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	potion := consumes.PotId
	prepot := consumes.PrepotId

	potionCD := character.GetPotionCD()

	if potion == 0 && prepot == 0 {
		return
	}
	var mcd MajorCooldown
	if prepot != 0 {
		mcd = makePotionActivationSpell(prepot, character, potionCD)
		if mcd.Spell != nil {
			mcd.Spell.Flags |= SpellFlagPrepullPotion
		}
	}

	var defaultMCD MajorCooldown
	if potion == prepot {
		defaultMCD = mcd
	} else {
		if potion != 0 {
			defaultMCD = makePotionActivationSpell(potion, character, potionCD)
		}
	}
	if defaultMCD.Spell != nil {
		defaultMCD.Spell.Flags |= SpellFlagCombatPotion
		character.AddMajorCooldown(defaultMCD)
	}
}

var AlchStoneItemIDs = []int32{80508, 96252, 96253, 96254, 44322, 44323, 44324}

func (character *Character) HasAlchStone() bool {
	alchStoneEquipped := false
	for _, itemID := range AlchStoneItemIDs {
		alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
	}
	return character.HasProfession(proto.Profession_Alchemy) && alchStoneEquipped
}

func makePotionActivationSpell(potionId int32, character *Character, potionCD *Timer) MajorCooldown {
	potion := ConsumableByID[potionId]
	mcd := makePotionActivationSpellInternal(potion, character, potionCD)
	if mcd.Spell != nil {
		// Mark as 'Encounter Only' so that users are forced to select the generic Potion
		// placeholder action instead of specific potion spells, in APL prepull. This
		// prevents a mismatch between Consumes and Rotation settings.
		mcd.Spell.Flags |= SpellFlagEncounterOnly | SpellFlagPotion
		oldApplyEffects := mcd.Spell.ApplyEffects
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			oldApplyEffects(sim, target, spell)
			if sim.CurrentTime < 0 {
				potionCD.Set(sim.CurrentTime + time.Minute)

				character.UpdateMajorCooldowns()
			}
		}
	}
	return mcd

}
func makePotionActivationSpellInternal(potion Consumable, character *Character, potionCD *Timer) MajorCooldown {
	alchStoneEquipped := character.HasAlchStone()
	potionCast := CastConfig{
		CD: Cooldown{
			Timer:    potionCD,
			Duration: time.Minute * 60, // Infinite CD
		},
		SharedCD: Cooldown{
			Timer:    character.GetPotionCD(),
			Duration: time.Minute * 60,
		},
	}

	//Simple stat potion

	actionID := ActionID{ItemID: potion.Id}
	if len(potion.EffectIds) == 0 {
		aura := character.NewTemporaryStatsAura(potion.Name, actionID, stats.FromProtoArray(potion.Stats), time.Second*time.Duration(potion.BuffDuration))
		return MajorCooldown{
			Type: CooldownTypeDPS,

			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if len(potion.EffectIds) >= 1 {
		// We assume its a hp or mana pot?
		// What other pots are there?
		mcd := MajorCooldown{
			Type: CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					for _, effectId := range potion.EffectIds {
						effect := EffectsById[effectId]
						if effect.EffectType == int32(dbc.E_HEAL) {
							healthMetrics := character.NewHealthMetrics(actionID)
							healthGain := sim.RollWithLabel(float64(effect.EffectBasePoints+1), float64(effect.EffectBasePoints+effect.EffectDieSides), potion.Name)
							if alchStoneEquipped {
								healthGain *= 1.40
							}
							character.GainHealth(sim, healthGain*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
						}
						if effect.EffectType == int32(dbc.E_ENERGIZE) {
							manaMetrics := character.NewManaMetrics(actionID)
							manaGain := sim.RollWithLabel(float64(effect.EffectBasePoints+1), float64(effect.EffectBasePoints+effect.EffectDieSides), potion.Name)
							if alchStoneEquipped {
								manaGain *= 1.4
							}
							character.AddMana(sim, manaGain, manaMetrics)
						}
					}
				},
			})}
		spell := character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast:     potionCast,
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				for _, effectId := range potion.EffectIds {
					effect := EffectsById[effectId]
					if effect.EffectType == int32(dbc.E_ENERGIZE) {
						manaMetrics := character.NewManaMetrics(actionID)
						manaGain := sim.RollWithLabel(float64(effect.EffectBasePoints+1), float64(effect.EffectBasePoints+effect.EffectDieSides), potion.Name)
						if alchStoneEquipped {
							manaGain *= 1.4
						}
						character.AddMana(sim, manaGain, manaMetrics)
						mcd.Type = CooldownTypeMana
						mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
							totalRegen := character.ManaRegenPerSecondWhileCombat() * 5
							manaGain := manaGain
							if alchStoneEquipped {
								manaGain *= 1.4
							}
							return character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
						}
					}
					if effect.EffectType == int32(dbc.E_HEAL) {
						healthMetrics := character.NewHealthMetrics(actionID)
						healthGain := sim.RollWithLabel(float64(effect.EffectBasePoints+1), float64(effect.EffectBasePoints+effect.EffectDieSides), potion.Name)
						if alchStoneEquipped {
							healthGain *= 1.40
						}
						character.GainHealth(sim, healthGain*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
						mcd.Type = CooldownTypeSurvival

					}

				}
			},
		})
		mcd.Spell = spell
		return mcd
	}

	return MajorCooldown{}
}

var ConjuredAuraTag = "Conjured"

func registerConjuredCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	conjuredType := consumes.DefaultConjured

	if conjuredType == proto.Conjured_ConjuredDarkRune {
		actionID := ActionID{ItemID: 20520}
		manaMetrics := character.NewManaMetrics(actionID)
		// damageTakenManaMetrics := character.NewManaMetrics(ActionID{SpellID: 33776})
		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 15,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				// Restores 900 to 1500 mana. (2 Min Cooldown)
				manaGain := sim.RollWithLabel(900, 1500, "dark rune")
				character.AddMana(sim, manaGain, manaMetrics)

				// if character.Class == proto.Class_ClassPaladin {
				// 	// Paladins gain extra mana from self-inflicted damage
				// 	// TO-DO: It is possible for damage to be resisted or to crit
				// 	// This would affect mana returns for Paladins
				// 	manaFromDamage := manaGain * 2.0 / 3.0 * 0.1
				// 	character.AddMana(sim, manaFromDamage, damageTakenManaMetrics, false)
				// }
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCombat() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 1500
			},
		})
	} else if conjuredType == proto.Conjured_ConjuredHealthstone {
		actionID := ActionID{ItemID: 5512}
		healthMetrics := character.NewHealthMetrics(actionID)

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				SharedCD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 2,
				},

				// Enforce only one HS per fight
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 60,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				character.GainHealth(sim, 0.45*character.baseStats[stats.Health], healthMetrics)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
		})
	}
}

var BigDaddyActionID = ActionID{SpellID: 89637}
var HighpoweredBoltGunActionID = ActionID{ItemID: 40771}

func registerExplosivesCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	if !character.HasProfession(proto.Profession_Engineering) {
		return
	}

	if consumes.ExplosiveBigDaddy {
		bomb := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    BigDaddyActionID,
			SpellSchool: SpellSchoolFire,
			ProcMask:    ProcMaskEmpty,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute,
				},

				DefaultCast: Cast{
					CastTime: time.Millisecond * 500,
				},

				ModifyCast: func(sim *Simulation, spell *Spell, cast *Cast) {
					spell.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
					spell.Unit.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime)
				},
			},

			// Explosives always have 1% resist chance, so just give them hit cap.
			BonusHitPercent:  100,
			DamageMultiplier: 1,
			CritMultiplier:   2,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				baseDamage := 5006 * sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    bomb,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
		})
	}

	if consumes.HighpoweredBoltGun {
		boltGun := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    ActionID{SpellID: 82207},
			SpellSchool: SpellSchoolFire,
			ProcMask:    ProcMaskEmpty,
			Flags:       SpellFlagNoOnCastComplete | SpellFlagCanCastWhileMoving,

			Cast: CastConfig{
				DefaultCast: Cast{
					GCD:      GCDDefault,
					CastTime: time.Second,
				},
				IgnoreHaste: true,
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
				SharedCD: Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 15,
				},
			},

			// Explosives always have 1% resist chance, so just give them hit cap.
			BonusHitPercent:  100,
			DamageMultiplier: 1,
			CritMultiplier:   2,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				spell.CalcAndDealDamage(sim, target, 8860, spell.OutcomeMagicHitAndCrit)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    boltGun,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
			ShouldActivate: func(s *Simulation, c *Character) bool {
				return false // Intentionally not automatically used
			},
		})
	}
}

func registerTinkerHandsCD(agent Agent, consumes *proto.Consumes) {
	if consumes.TinkerHands == proto.TinkerHands_TinkerHandsNone {
		return
	}
	character := agent.GetCharacter()
	if !character.HasProfession(proto.Profession_Engineering) {
		return
	}

	switch consumes.TinkerHands {
	case proto.TinkerHands_TinkerHandsSynapseSprings:
		// Enchant: 4179, Spell: 82174 - Synapse Springs
		statType := character.GetHighestStatType([]stats.Stat{stats.Intellect, stats.Strength, stats.Agility})

		var actionID ActionID
		var highestStat stats.Stats
		var label string
		switch statType {
		case stats.Intellect:
			actionID = ActionID{SpellID: 96230}
			highestStat = stats.Stats{stats.Intellect: 480}
			label = "Synapse Springs - Int"
		case stats.Agility:
			actionID = ActionID{SpellID: 96228}
			highestStat = stats.Stats{stats.Agility: 480}
			label = "Synapse Springs - Agi"
		case stats.Strength:
			actionID = ActionID{SpellID: 96229}
			highestStat = stats.Stats{stats.Strength: 480}
			label = "Synapse Springs - Str"
		default:
			panic("Stat type doesn't match any defined case")
		}

		aura := character.NewTemporaryStatsAura(
			label,
			actionID,
			highestStat,
			time.Second*10,
		)

		spell := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    ActionID{SpellID: 82174},
			SpellSchool: SpellSchoolPhysical,
			Flags:       SpellFlagNoOnCastComplete,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
				SharedCD: Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 10,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    spell,
			Priority: CooldownPriorityLow,
			Type:     CooldownTypeDPS,
		})
		character.ItemSwap.ProcessTinker(spell, []proto.ItemSlot{proto.ItemSlot_ItemSlotHands})
	case proto.TinkerHands_TinkerHandsQuickflipDeflectionPlates:
		// Enchant: 4180, Spell: 82176 - Quickflip Deflection Plates
		actionID := ActionID{SpellID: 82176}
		statAura := character.NewTemporaryStatsAura(
			"Quickflip Deflection Plates Buff",
			actionID,
			stats.Stats{stats.Armor: 1500},
			time.Second*12,
		)

		spell := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    actionID,
			SpellSchool: SpellSchoolPhysical,
			Flags:       SpellFlagNoOnCastComplete,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				statAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    spell,
			Priority: CooldownPriorityLow,
			Type:     CooldownTypeSurvival,
		})
		character.ItemSwap.ProcessTinker(spell, []proto.ItemSlot{proto.ItemSlot_ItemSlotHands})
	case proto.TinkerHands_TinkerHandsTazikShocker:
		// Enchant: 4181, Spell: 82180 - Tazik Shocker
		actionID := ActionID{SpellID: 82179}
		spell := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    actionID,
			SpellSchool: SpellSchoolNature,
			Flags:       SpellFlagNoOnCastComplete,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 120,
				},
			},

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *Simulation, unit *Unit, spell *Spell) {
				// Benerfits from enhancement mastery
				// Ele crit dmg multi
				// Moonkin eclipse, so basically everything
				spell.CalcAndDealDamage(sim, unit, sim.Roll(4320, 961), spell.OutcomeMagicHitAndCrit)
			},
		})
		character.ItemSwap.ProcessTinker(spell, []proto.ItemSlot{proto.ItemSlot_ItemSlotHands})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    spell,
			Priority: CooldownPriorityLow,
			Type:     CooldownTypeDPS,
		})
	case proto.TinkerHands_TinkerHandsSpinalHealingInjector:
		// Enchant: 4182, Spell: 82184 - Spinal Healing Injector
		actionID := ActionID{SpellID: 82184}
		healthMetric := character.NewHealthMetrics(actionID)
		spell := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    actionID,
			SpellSchool: SpellSchoolPhysical,
			Flags:       SpellFlagNoOnCastComplete | SpellFlagCombatPotion,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
				SharedCD: Cooldown{
					Timer:    character.GetPotionCD(),
					Duration: time.Minute * 60,
				},
			},

			ApplyEffects: func(sim *Simulation, unit *Unit, spell *Spell) {
				result := sim.Roll(27000, 33000)
				if character.HasAlchStone() {
					result *= 1.4
				}

				character.GainHealth(sim, result, healthMetric)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    spell,
			Priority: CooldownPriorityLow,
			Type:     CooldownTypeSurvival,
		})
		character.ItemSwap.ProcessTinker(spell, []proto.ItemSlot{proto.ItemSlot_ItemSlotHands})
	case proto.TinkerHands_TinkerHandsZ50ManaGulper:
		// Enchant: 4183, Spell: 82186 - Z50 Mana Gulper
		actionId := ActionID{SpellID: 82186}
		manaMetric := character.NewManaMetrics(actionId)
		spell := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    actionId,
			SpellSchool: SpellSchoolPhysical,
			Flags:       SpellFlagNoOnCastComplete | SpellFlagPotion,

			// TODO: In theory those ingi on-use enchants share a CD with potions
			// The potion CD timer is not available right now
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Second * 60,
				},
				SharedCD: Cooldown{
					Timer:    character.GetPotionCD(),
					Duration: time.Minute * 60,
				},
			},

			ApplyEffects: func(sim *Simulation, unit *Unit, spell *Spell) {
				mana := sim.Roll(10730, 12470)
				if character.HasAlchStone() {
					mana *= 1.4
				}

				character.AddMana(sim, mana, manaMetric)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			ShouldActivate: func(s *Simulation, c *Character) bool {
				return c.HasManaBar() && (c.MaxMana()-c.CurrentMana()) > 10730
			},
			Spell:    spell,
			Priority: CooldownPriorityLow,
			Type:     CooldownTypeMana,
		})
		character.ItemSwap.ProcessTinker(spell, []proto.ItemSlot{proto.ItemSlot_ItemSlotHands})
	}
}
