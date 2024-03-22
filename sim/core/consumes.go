package core

import (
	"time"

	"github.com/wowsims/cata/sim/core"
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
	alchemyFlaskBonus := core.TernaryInt(character.HasProfession(proto.Profession_Alchemy), 80, 0)
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
				stats.Stamina: 450,
			})
		case proto.Flask_FlaskOfFlowingWater:
			character.AddStats(stats.Stats{
				stats.Spirit: 300 + alchemyFlaskBonus,
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
				stats.Resilience: 50,
			})
			if character.HasProfession(proto.Profession_Alchemy) {
				character.AddStats(stats.Stats{
					stats.Resilience: 82,
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
				stats.Mastery: 225,
			})
		case proto.BattleElixir_ElixirOfMightySpeed:
			character.AddStats(stats.Stats{
				stats.MeleeHaste: 225,
				stats.SpellHaste: 225,
			})
		case proto.BattleElixir_ElixirOfImpossibleAccuracy:
			character.AddStats(stats.Stats{
				stats.MeleeHit: 225,
				stats.SpellHit: 225,
			})
		case proto.BattleElixir_ElixirOfTheCobra:
			character.AddStats(stats.Stats{
				stats.MeleeCrit: 225,
				stats.SpellCrit: 225,
			})
		case proto.BattleElixir_ElixirOfTheNaga:
			character.AddStats(stats.Stats{
				stats.Expertise: 225,
			})
		case proto.BattleElixir_GhostElixir:
			character.AddStats(stats.Stats{
				stats.Spirit: 225,
			})
		case proto.BattleElixir_ElixirOfAccuracy:
			character.AddStats(stats.Stats{
				stats.MeleeHit: 45,
				stats.SpellHit: 45,
			})
		case proto.BattleElixir_ElixirOfArmorPiercing:
			character.AddStats(stats.Stats{
				stats.ArmorPenetration: 45,
			})
		case proto.BattleElixir_ElixirOfDeadlyStrikes:
			character.AddStats(stats.Stats{
				stats.MeleeCrit: 45,
				stats.SpellCrit: 45,
			})
		case proto.BattleElixir_ElixirOfExpertise:
			character.AddStats(stats.Stats{
				stats.Expertise: 45,
			})
		case proto.BattleElixir_ElixirOfLightningSpeed:
			character.AddStats(stats.Stats{
				stats.MeleeHaste: 45,
				stats.SpellHaste: 45,
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
				stats.Defense: 45,
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
		case proto.GuardianElixir_GiftOfArthas:
			character.AddStats(stats.Stats{
				stats.ShadowResistance: 10,
			})

			debuffAuras := (&character.Unit).NewEnemyAuraArray(GiftOfArthasAura)

			actionID := ActionID{SpellID: 11374}
			goaProc := character.RegisterSpell(SpellConfig{
				ActionID:    actionID,
				SpellSchool: SpellSchoolNature,
				ProcMask:    ProcMaskEmpty,

				ThreatMultiplier: 1,
				FlatThreatBonus:  90,

				ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
					debuffAuras.Get(target).Activate(sim)
					spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
				},
			})

			character.RegisterAura(Aura{
				Label:    "Gift of Arthas",
				Duration: NeverExpires,
				OnReset: func(aura *Aura, sim *Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
					if result.Landed() &&
						spell.SpellSchool == SpellSchoolPhysical &&
						sim.RandomFloat("Gift of Arthas") < 0.3 {
						goaProc.Cast(sim, spell.Unit)
					}
				},
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
			stats.ArmorPenetration: 40,
			stats.Stamina:          40,
		})
	case proto.Food_FoodMegaMammothMeal:
		character.AddStats(stats.Stats{
			stats.AttackPower:       80,
			stats.RangedAttackPower: 80,
			stats.Stamina:           40,
		})
	case proto.Food_FoodSpicedWormBurger:
		character.AddStats(stats.Stats{
			stats.MeleeCrit: 40,
			stats.SpellCrit: 40,
			stats.Stamina:   40,
		})
	case proto.Food_FoodRhinoliciousWormsteak:
		character.AddStats(stats.Stats{
			stats.Expertise: 40,
			stats.Stamina:   40,
		})
	case proto.Food_FoodImperialMantaSteak:
		character.AddStats(stats.Stats{
			stats.MeleeHaste: 40,
			stats.SpellHaste: 40,
			stats.Stamina:    40,
		})
	case proto.Food_FoodSnapperExtreme:
		character.AddStats(stats.Stats{
			stats.MeleeHit: 40,
			stats.SpellHit: 40,
			stats.Stamina:  40,
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
			stats.SpellCrit: 20,
			stats.Spirit:    20,
		})
	case proto.Food_FoodSpicyHotTalbuk:
		character.AddStats(stats.Stats{
			stats.MeleeHit: 20,
			stats.Spirit:   20,
		})
	case proto.Food_FoodFishermansFeast:
		character.AddStats(stats.Stats{
			stats.Stamina: 30,
			stats.Spirit:  20,
		})
	}

	registerPotionCD(agent, consumes)
	registerConjuredCD(agent, consumes)
	registerExplosivesCD(agent, consumes)
}

func ApplyPetConsumeEffects(pet *Character, ownerConsumes *proto.Consumes) {
	switch ownerConsumes.PetFood {
	case proto.PetFood_PetFoodSpicedMammothTreats:
		pet.AddStats(stats.Stats{
			stats.Strength: 30,
			stats.Stamina:  30,
		})
	case proto.PetFood_PetFoodKiblersBits:
		pet.AddStats(stats.Stats{
			stats.Strength: 20,
			stats.Stamina:  20,
		})
	}

	pet.AddStat(stats.Agility, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfAgility])
	pet.AddStat(stats.Strength, []float64{0, 5, 9, 13, 17, 20}[ownerConsumes.PetScrollOfStrength])
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	defaultPotion := consumes.DefaultPotion
	startingPotion := consumes.PrepopPotion

	potionCD := character.NewTimer()
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
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
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
		aura := character.NewTemporaryStatsAura("Golemblood Potion", actionID, stats.Stats{stats.Strength: 1200}, time.Second*15)
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
		aura := character.NewTemporaryStatsAura("Potion of the Tol'vir", actionID, stats.Stats{stats.Agility: 1200}, time.Second*15)
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
		aura := character.NewTemporaryStatsAura("Volcanic Potion", actionID, stats.Stats{stats.Intellect: 1200}, time.Second*15)
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
		aura := character.NewTemporaryStatsAura("Potion of Speed", actionID, stats.Stats{stats.MeleeHaste: 500, stats.SpellHaste: 500}, time.Second*15)
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
		aura := character.NewTemporaryStatsAura("Haste Potion", actionID, stats.Stats{stats.MeleeHaste: 400}, time.Second*15)
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
				totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
				return character.MaxMana()-(character.CurrentMana()+totalRegen) >= 1500
			},
		})
	} else if conjuredType == proto.Conjured_ConjuredFlameCap {
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

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				flameCapAura.Activate(sim)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	} else if conjuredType == proto.Conjured_ConjuredHealthstone {
		actionID := ActionID{ItemID: 36892}
		healthMetrics := character.NewHealthMetrics(actionID)

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.GetConjuredCD(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				character.GainHealth(sim, 4280*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeSurvival,
		})
	}
}

var ThermalSapperActionID = ActionID{ItemID: 42641}
var ExplosiveDecoyActionID = ActionID{ItemID: 40536}
var SaroniteBombActionID = ActionID{ItemID: 41119}
var CobaltFragBombActionID = ActionID{ItemID: 40771}

func registerExplosivesCD(agent Agent, consumes *proto.Consumes) {
	character := agent.GetCharacter()
	hasFiller := consumes.FillerExplosive != proto.Explosive_ExplosiveUnknown
	if !character.HasProfession(proto.Profession_Engineering) {
		return
	}
	if !consumes.ThermalSapper && !consumes.ExplosiveDecoy && !hasFiller {
		return
	}
	sharedTimer := character.NewTimer()

	if consumes.ThermalSapper {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    character.newThermalSapperSpell(sharedTimer),
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 30,
		})
	}

	if consumes.ExplosiveDecoy {
		character.AddMajorCooldown(MajorCooldown{
			Spell:    character.newExplosiveDecoySpell(sharedTimer),
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 20,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Decoy puts other explosives on 2m CD, so only use if there won't be enough
				// time to use another explosive OR there is no filler explosive.
				return sim.GetRemainingDuration() < time.Minute || !hasFiller
			},
		})
	}

	if hasFiller {
		var filler *Spell
		switch consumes.FillerExplosive {
		case proto.Explosive_ExplosiveSaroniteBomb:
			filler = character.newSaroniteBombSpell(sharedTimer)
		case proto.Explosive_ExplosiveCobaltFragBomb:
			filler = character.newCobaltFragBombSpell(sharedTimer)
		}

		character.AddMajorCooldown(MajorCooldown{
			Spell:    filler,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
		})
	}
}

// Creates a spell object for the common explosive case.
func (character *Character) newBasicExplosiveSpellConfig(sharedTimer *Timer, actionID ActionID, school SpellSchool, minDamage float64, maxDamage float64, cooldown Cooldown, _ float64, _ float64) SpellConfig {
	dealSelfDamage := actionID.SameAction(ThermalSapperActionID)

	return SpellConfig{
		ActionID:    actionID,
		SpellSchool: school,
		ProcMask:    ProcMaskEmpty,

		Cast: CastConfig{
			CD: cooldown,
			SharedCD: Cooldown{
				Timer:    sharedTimer,
				Duration: TernaryDuration(actionID.SameAction(ExplosiveDecoyActionID), time.Minute*2, time.Minute),
			},
		},

		// Explosives always have 1% resist chance, so just give them hit cap.
		BonusHitRating:   100 * SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(minDamage, maxDamage) * sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if dealSelfDamage {
				baseDamage := sim.Roll(minDamage, maxDamage)
				spell.CalcAndDealDamage(sim, &character.Unit, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	}
}
func (character *Character) newThermalSapperSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ThermalSapperActionID, SpellSchoolFire, 2188, 2812, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 5}, 2188, 2812))
}
func (character *Character) newExplosiveDecoySpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, ExplosiveDecoyActionID, SpellSchoolPhysical, 1440, 2160, Cooldown{Timer: character.NewTimer(), Duration: time.Minute * 2}, 0, 0))
}
func (character *Character) newSaroniteBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, SaroniteBombActionID, SpellSchoolFire, 1150, 1500, Cooldown{}, 0, 0))
}
func (character *Character) newCobaltFragBombSpell(sharedTimer *Timer) *Spell {
	return character.GetOrRegisterSpell(character.newBasicExplosiveSpellConfig(sharedTimer, CobaltFragBombActionID, SpellSchoolFire, 750, 1000, Cooldown{}, 0, 0))
}
