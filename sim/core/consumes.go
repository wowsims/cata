package core

import (
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
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
	if consumes.Flask != proto.Flask_FlaskUnknown {
		switch consumes.Flask {
		case proto.Flask_FlaskOfTitanicStrength:
			character.AddStats(stats.Stats{
				stats.Strength: 300 + alchemyFlaskBonus,
			})
		case proto.Flask_FlaskOfTheWinds:
			character.AddStats(stats.Stats{
				stats.Agility: 300 + alchemyFlaskBonus,
			})
		case proto.Flask_FlaskOfSteelskin:
			character.AddStats(stats.Stats{
				stats.Stamina: 450 + alchemyFlaskBonus*1.5,
			})
		case proto.Flask_FlaskOfFlowingWater:
			character.AddStats(stats.Stats{
				stats.Spirit: 300 + alchemyFlaskBonus,
			})
		case proto.Flask_FlaskOfTheDraconicMind:
			character.AddStats(stats.Stats{
				stats.Intellect: 300 + alchemyFlaskBonus,
			})
		case proto.Flask_FlaskOfTheFrostWyrm:
			character.AddStats(stats.Stats{
				stats.SpellPower: 125,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.SpellPower: 47,
				})
			}
		case proto.Flask_FlaskOfEndlessRage:
			character.AddStats(stats.Stats{
				stats.AttackPower:       180,
				stats.RangedAttackPower: 180,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.AttackPower:       80,
					stats.RangedAttackPower: 80,
				})
			}
		case proto.Flask_FlaskOfPureMojo:
			character.AddStats(stats.Stats{
				stats.MP5: 45,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.MP5: 20,
				})
			}
		case proto.Flask_FlaskOfStoneblood:
			character.AddStats(stats.Stats{
				stats.Health: 1300,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.Health: 650,
				})
			}
		case proto.Flask_LesserFlaskOfToughness:
			character.AddStats(stats.Stats{
				stats.ResilienceRating: 50,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.ResilienceRating: 82,
				})
			}
		case proto.Flask_LesserFlaskOfResistance:
			character.AddStats(stats.Stats{
				stats.ArcaneResistance: 50,
				stats.FireResistance:   50,
				stats.FrostResistance:  50,
				stats.NatureResistance: 50,
				stats.ShadowResistance: 50,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.ArcaneResistance: 40,
					stats.FireResistance:   40,
					stats.FrostResistance:  40,
					stats.NatureResistance: 40,
					stats.ShadowResistance: 40,
				})
			}
		}
	} else {
		switch consumes.BattleElixir {
		case proto.BattleElixir_ElixirOfTheMaster:
			character.AddStats(stats.Stats{
				stats.MasteryRating: 225 + alchemyBattleElixirBonus,
			})
		case proto.BattleElixir_ElixirOfMightySpeed:
			character.AddStats(stats.Stats{
				stats.HasteRating: 225 + alchemyBattleElixirBonus,
			})
		case proto.BattleElixir_ElixirOfImpossibleAccuracy:
			character.AddStats(stats.Stats{
				stats.HitRating: 225 + alchemyBattleElixirBonus,
			})
		case proto.BattleElixir_ElixirOfTheCobra:
			character.AddStats(stats.Stats{
				stats.CritRating: 225 + alchemyBattleElixirBonus,
			})
		case proto.BattleElixir_ElixirOfTheNaga:
			character.AddStats(stats.Stats{
				stats.ExpertiseRating: 225 + alchemyBattleElixirBonus,
			})
		case proto.BattleElixir_GhostElixir:
			character.AddStats(stats.Stats{
				stats.Spirit: 225 + alchemyBattleElixirBonus,
			})
		case proto.BattleElixir_ElixirOfAccuracy:
			character.AddStats(stats.Stats{
				stats.HitRating: 45,
			})
		case proto.BattleElixir_ElixirOfArmorPiercing:
			character.AddStats(stats.Stats{
				stats.Agility:    25,
				stats.CritRating: 25,
			})
		case proto.BattleElixir_ElixirOfDeadlyStrikes:
			character.AddStats(stats.Stats{
				stats.CritRating: 45,
			})
		case proto.BattleElixir_ElixirOfExpertise:
			character.AddStats(stats.Stats{
				stats.ExpertiseRating: 45,
			})
		case proto.BattleElixir_ElixirOfLightningSpeed:
			character.AddStats(stats.Stats{
				stats.HasteRating: 45,
			})
		case proto.BattleElixir_ElixirOfMightyAgility:
			character.AddStats(stats.Stats{
				stats.Agility: 45,
			})
		case proto.BattleElixir_ElixirOfMightyStrength:
			character.AddStats(stats.Stats{
				stats.Strength: 45,
			})
		case proto.BattleElixir_GurusElixir:
			character.AddStats(stats.Stats{
				stats.Agility:   20,
				stats.Strength:  20,
				stats.Stamina:   20,
				stats.Intellect: 20,
				stats.Spirit:    20,
			})
		case proto.BattleElixir_SpellpowerElixir:
			character.AddStats(stats.Stats{
				stats.SpellPower: 58,
			})
		case proto.BattleElixir_WrathElixir:
			character.AddStats(stats.Stats{
				stats.AttackPower:       90,
				stats.RangedAttackPower: 90,
			})
		case proto.BattleElixir_ElixirOfDemonslaying:
			if character.CurrentTarget.MobType == proto.MobType_MobTypeDemon {
				character.PseudoStats.MobTypeAttackPower += 265
			}
		}

		switch consumes.GuardianElixir {
		case proto.GuardianElixir_ElixirOfDeepEarth:
			character.AddStats(stats.Stats{
				stats.Armor: 900,
			})
		case proto.GuardianElixir_PrismaticElixir:
			character.AddStats(stats.Stats{
				stats.ArcaneResistance: 90,
				stats.FireResistance:   90,
				stats.FrostResistance:  90,
				stats.NatureResistance: 90,
				stats.ShadowResistance: 90,
			})
		case proto.GuardianElixir_ElixirOfMightyDefense:
			character.AddStats(stats.Stats{
				stats.Armor: 180,
			})
		case proto.GuardianElixir_ElixirOfMightyFortitude:
			character.AddStats(stats.Stats{
				stats.Health: 350,
			})
		case proto.GuardianElixir_ElixirOfMightyMageblood:
			character.AddStats(stats.Stats{
				stats.MP5: 30,
			})
		case proto.GuardianElixir_ElixirOfMightyThoughts:
			character.AddStats(stats.Stats{
				stats.Intellect: 45,
			})
		case proto.GuardianElixir_ElixirOfProtection:
			character.AddStats(stats.Stats{
				stats.Armor: 800,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.Armor: 280,
				})
			}
		case proto.GuardianElixir_ElixirOfSpirit:
			character.AddStats(stats.Stats{
				stats.Spirit: 50,
			})
		}
	}

	switch consumes.Food {
	case proto.Food_FoodFishFeast:
		character.AddStats(stats.Stats{
			stats.AttackPower:       80,
			stats.RangedAttackPower: 80,
			stats.SpellPower:        46,
			stats.Stamina:           40,
		})
	case proto.Food_FoodGreatFeast:
		character.AddStats(stats.Stats{
			stats.AttackPower:       60,
			stats.RangedAttackPower: 60,
			stats.SpellPower:        35,
			stats.Stamina:           30,
		})
	case proto.Food_FoodBlackenedDragonfin:
		character.AddStats(stats.Stats{
			stats.Agility: 40,
			stats.Stamina: 40,
		})
	case proto.Food_FoodHeartyRhino:
		character.AddStats(stats.Stats{
			stats.CritRating: 40,
			stats.Stamina:    40,
		})
	case proto.Food_FoodMegaMammothMeal:
		character.AddStats(stats.Stats{
			stats.AttackPower:       80,
			stats.RangedAttackPower: 80,
			stats.Stamina:           40,
		})
	case proto.Food_FoodSpicedWormBurger:
		character.AddStats(stats.Stats{
			stats.CritRating: 40,
			stats.Stamina:    40,
		})
	case proto.Food_FoodRhinoliciousWormsteak:
		character.AddStats(stats.Stats{
			stats.ExpertiseRating: 40,
			stats.Stamina:         40,
		})
	case proto.Food_FoodImperialMantaSteak:
		character.AddStats(stats.Stats{
			stats.HasteRating: 40,
			stats.Stamina:     40,
		})
	case proto.Food_FoodSnapperExtreme:
		character.AddStats(stats.Stats{
			stats.HitRating: 40,
			stats.Stamina:   40,
		})
	case proto.Food_FoodMightyRhinoDogs:
		character.AddStats(stats.Stats{
			stats.MP5:     16,
			stats.Stamina: 40,
		})
	case proto.Food_FoodFirecrackerSalmon:
		character.AddStats(stats.Stats{
			stats.SpellPower: 46,
			stats.Stamina:    40,
		})
	case proto.Food_FoodCuttlesteak:
		character.AddStats(stats.Stats{
			stats.Spirit:  40,
			stats.Stamina: 40,
		})
	case proto.Food_FoodDragonfinFilet:
		character.AddStats(stats.Stats{
			stats.Strength: 40,
			stats.Stamina:  40,
		})
	case proto.Food_FoodBlackenedBasilisk:
		character.AddStats(stats.Stats{
			stats.SpellPower: 23,
			stats.Spirit:     20,
		})
	case proto.Food_FoodGrilledMudfish:
		character.AddStats(stats.Stats{
			stats.Agility: 20,
			stats.Spirit:  20,
		})
	case proto.Food_FoodRavagerDog:
		character.AddStats(stats.Stats{
			stats.AttackPower:       40,
			stats.RangedAttackPower: 40,
			stats.Spirit:            20,
		})
	case proto.Food_FoodRoastedClefthoof:
		character.AddStats(stats.Stats{
			stats.Strength: 20,
			stats.Spirit:   20,
		})
	case proto.Food_FoodSkullfishSoup:
		character.AddStats(stats.Stats{
			stats.CritRating: 20,
			stats.Spirit:     20,
		})
	case proto.Food_FoodSpicyHotTalbuk:
		character.AddStats(stats.Stats{
			stats.HitRating: 20,
			stats.Spirit:    20,
		})
	case proto.Food_FoodFishermansFeast:
		character.AddStats(stats.Stats{
			stats.Stamina: 30,
			stats.Spirit:  20,
		})
	case proto.Food_FoodSeafoodFeast:
		character.AddStat(stats.Stamina, 90)
		character.AddStat(character.GetHighestStatType([]stats.Stat{stats.Strength, stats.Agility, stats.Intellect}), 90)
	case proto.Food_FoodFortuneCookie:
		character.AddStat(stats.Stamina, 90)
		character.AddStat(character.GetHighestStatType([]stats.Stat{stats.Strength, stats.Agility, stats.Intellect}), 90)
	case proto.Food_FoodSeveredSagefish:
		character.AddStats(stats.Stats{
			stats.Stamina:   90,
			stats.Intellect: 90,
		})
	case proto.Food_FoodBeerBasedCrocolisk:
		character.AddStats(stats.Stats{
			stats.Stamina:  90,
			stats.Strength: 90,
		})
	case proto.Food_FoodSkeweredEel:
		character.AddStats(stats.Stats{
			stats.Stamina: 90,
			stats.Agility: 90,
		})
	case proto.Food_FoodDeliciousSagefishTail:
		character.AddStats(stats.Stats{
			stats.Stamina: 90,
			stats.Spirit:  90,
		})
	case proto.Food_FoodBasiliskLiverdog:
		character.AddStats(stats.Stats{
			stats.Stamina:     90,
			stats.HasteRating: 90,
		})
	case proto.Food_FoodBakedRockfish:
		character.AddStats(stats.Stats{
			stats.Stamina:    90,
			stats.CritRating: 90,
		})
	case proto.Food_FoodCrocoliskAuGratin:
		character.AddStats(stats.Stats{
			stats.Stamina:         90,
			stats.ExpertiseRating: 90,
		})
	case proto.Food_FoodGrilledDragon:
		character.AddStats(stats.Stats{
			stats.Stamina:   90,
			stats.HitRating: 90,
		})
	case proto.Food_FoodLavascaleMinestrone:
		character.AddStats(stats.Stats{
			stats.Stamina:       90,
			stats.MasteryRating: 90,
		})
	case proto.Food_FoodBlackbellySushi:
		character.AddStats(stats.Stats{
			stats.Stamina:     90,
			stats.ParryRating: 90,
		})
	case proto.Food_FoodMushroomSauceMudfish:
		character.AddStats(stats.Stats{
			stats.Stamina:     90,
			stats.DodgeRating: 90,
		})
	}

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
	registerTinkerHandsCD(agent, consumes)
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion
	startingPotion := consumes.PrepopPotion

	potionCD := character.GetPotionCD()
	// if character.Spec == proto.Spec_SpecBalanceDruid {
	// 	// Create both pots spells so they will be selectable in APL UI regardless of settings.
	// 	speedMCD := makePotionActivation(proto.Potions_PotionOfSpeed, character, potionCD)
	// 	wildMagicMCD := makePotionActivation(proto.Potions_PotionOfWildMagic, character, potionCD)
	// 	speedMCD.Spell.Flags |= SpellFlagAPL | SpellFlagMCD
	// 	wildMagicMCD.Spell.Flags |= SpellFlagAPL | SpellFlagMCD
	// }

	if defaultPotion == proto.Potions_UnknownPotion && startingPotion == proto.Potions_UnknownPotion {
		return
	}

	startingMCD := makePotionActivation(startingPotion, character, potionCD)
	if startingMCD.Spell != nil {
		startingMCD.Spell.Flags |= SpellFlagPrepullPotion
	}

	var defaultMCD MajorCooldown
	if defaultPotion == startingPotion {
		defaultMCD = startingMCD
	} else {
		defaultMCD = makePotionActivation(defaultPotion, character, potionCD)
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

func makePotionActivation(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
	mcd := makePotionActivationInternal(potionType, character, potionCD)
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

func makePotionActivationInternal(potionType proto.Potions, character *Character, potionCD *Timer) MajorCooldown {
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

	if potionType == proto.Potions_MythicalHealingPotion {
		actionID := ActionID{ItemID: 57191}
		healthMetrics := character.NewHealthMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					healthGain := sim.RollWithLabel(22500, 27500, "MythicalHealingPotion")

					if alchStoneEquipped && potionType == proto.Potions_MythicalHealingPotion {
						healthGain *= 1.40
					}
					character.GainHealth(sim, healthGain*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_MythicalManaPotion {
		actionID := ActionID{ItemID: 57192}
		manaMetrics := character.NewManaMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCombat() * 5
				manaGain := 10750.0
				if alchStoneEquipped && potionType == proto.Potions_MythicalManaPotion {
					manaGain *= 1.4
				}
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					manaGain := sim.RollWithLabel(9250, 10750, "MythicalManaPotion")
					if alchStoneEquipped && potionType == proto.Potions_MythicalManaPotion {
						manaGain *= 1.4
					}
					character.AddMana(sim, manaGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_GolembloodPotion {
		actionID := ActionID{ItemID: 58146}
		aura := character.NewTemporaryStatsAura("Golemblood Potion", actionID, stats.Stats{stats.Strength: 1200}, time.Second*25)
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
	} else if potionType == proto.Potions_PotionOfTheTolvir {
		actionID := ActionID{ItemID: 58145}
		aura := character.NewTemporaryStatsAura("Potion of the Tol'vir", actionID, stats.Stats{stats.Agility: 1200}, time.Second*25)
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
	} else if potionType == proto.Potions_PotionOfConcentration {
		//actionID := ActionID{ItemID: 57194}
		// Todo: Implement. Has a cast time of 10seconds and you regain mana while casting it
		// Not sure about exact functionality
	} else if potionType == proto.Potions_VolcanicPotion {
		actionID := ActionID{ItemID: 58091}
		aura := character.NewTemporaryStatsAura("Volcanic Potion", actionID, stats.Stats{stats.Intellect: 1200}, time.Second*25)
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
	} else if potionType == proto.Potions_EarthenPotion {
		actionID := ActionID{ItemID: 58090}
		aura := character.NewTemporaryStatsAura("Earthen Potion", actionID, stats.Stats{stats.Armor: 4800}, time.Second*25) // Adjust stats as necessary
		return MajorCooldown{
			Type: CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
				},
			}),
		}
	} else if potionType == proto.Potions_MightyRejuvenationPotion {
		actionID := ActionID{ItemID: 57193}
		// No specific aura stats provided; adjust as needed
		manaMetrics := character.NewManaMetrics(actionID)
		healthMetrics := character.NewHealthMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeSurvival,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only pop if we have less than the max mana provided by the potion minus 1mp5 tick.
				totalRegen := character.ManaRegenPerSecondWhileCombat() * 5
				manaGain := 11000.0
				if alchStoneEquipped && potionType == proto.Potions_MythicalManaPotion {
					manaGain *= 1.4
				}
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,

				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					resourceGain := sim.RollWithLabel(9000, 11000, "MightyRejuvPotion") // Todo: Does it roll once or twice?

					if alchStoneEquipped && potionType == proto.Potions_MightyRejuvenationPotion {
						resourceGain *= 1.40
					}
					character.GainHealth(sim, resourceGain*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
					character.AddMana(sim, resourceGain, manaMetrics)
				},
			}),
		}
	} else if potionType == proto.Potions_PotionOfSpeed {
		actionID := ActionID{ItemID: 40211}
		aura := character.NewTemporaryStatsAura("Potion of Speed", actionID, stats.Stats{stats.HasteRating: 500}, time.Second*15)
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
	} else if potionType == proto.Potions_HastePotion {
		actionID := ActionID{ItemID: 22838}
		aura := character.NewTemporaryStatsAura("Haste Potion", actionID, stats.Stats{stats.HasteRating: 400}, time.Second*15)
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
	} else if potionType == proto.Potions_MightyRagePotion {
		actionID := ActionID{ItemID: 13442}
		aura := character.NewTemporaryStatsAura("Mighty Rage Potion", actionID, stats.Stats{stats.Strength: 60}, time.Second*15)
		rageMetrics := character.NewRageMetrics(actionID)
		return MajorCooldown{
			Type: CooldownTypeDPS,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				if character.Class == proto.Class_ClassWarrior {
					return character.CurrentRage() < 25
				}
				return true
			},
			Spell: character.GetOrRegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					aura.Activate(sim)
					if character.Class == proto.Class_ClassWarrior {
						bonusRage := sim.RollWithLabel(45, 75, "Mighty Rage Potion")
						character.AddRage(sim, bonusRage, rageMetrics)
					}
				},
			}),
		}
	} else if potionType == proto.Potions_FlameCap {
		actionID := ActionID{ItemID: 22788}

		flameCapProc := character.RegisterSpell(SpellConfig{
			ActionID:    actionID,
			ProcMask:    ProcMaskEmpty,
			SpellSchool: SpellSchoolFire,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHitAndCrit)
			},
		})

		const procChance = 0.185
		var fireSpells []*Spell
		character.OnSpellRegistered(func(spell *Spell) {
			if spell.SpellSchool.Matches(SpellSchoolFire) {
				fireSpells = append(fireSpells, spell)
			}
		})

		flameCapAura := character.RegisterAura(Aura{
			Label:    "Flame Cap",
			ActionID: actionID,
			Duration: time.Minute,
			OnGain: func(aura *Aura, sim *Simulation) {
				for _, spell := range fireSpells {
					spell.BonusSpellPower += 80
				}
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				for _, spell := range fireSpells {
					spell.BonusSpellPower -= 80
				}
			},
			OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(ProcMaskMeleeOrRanged) {
					return
				}
				if sim.RandomFloat("Flame Cap Melee") > procChance {
					return
				}

				flameCapProc.Cast(sim, result.Target)
			},
		})

		return MajorCooldown{
			Type: CooldownTypeDPS,
			Spell: character.RegisterSpell(SpellConfig{
				ActionID: actionID,
				Flags:    SpellFlagNoOnCastComplete,
				Cast:     potionCast,
				ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
					flameCapAura.Activate(sim)
				},
			}),
		}
	} else {
		return MajorCooldown{}
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
	}
}
