package core

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// Registers all consume-related effects to the Agent.
func applyConsumeEffects(agent Agent) {
	character := agent.GetCharacter()
	consumables := character.Consumables
	if consumables == nil {
		return
	}
	alchemyFlaskBonus := TernaryFloat64(character.HasProfession(proto.Profession_Alchemy), 320, 0)
	alchemyBattleElixirBonus := TernaryFloat64(character.HasProfession(proto.Profession_Alchemy), 240, 0)
	if consumables.FlaskId != 0 {
		flask := ConsumablesByID[consumables.FlaskId]
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
		character.AddStats(flask.Stats)
	}

	if consumables.BattleElixirId != 0 {
		elixir := ConsumablesByID[consumables.BattleElixirId]
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
		character.AddStats(elixir.Stats)
	}

	if consumables.GuardianElixirId != 0 {
		elixir := ConsumablesByID[consumables.GuardianElixirId]
		if character.HasProfession(proto.Profession_Alchemy) && elixir.Stats[stats.Armor] > 0 {
			elixir.Stats[stats.Armor] += 280
		}
		character.AddStats(elixir.Stats)
	}
	if consumables.FoodId != 0 {
		food := ConsumablesByID[consumables.FoodId]
		isPanda := character.Race == proto.Race_RaceHordePandaren || character.Race == proto.Race_RaceAlliancePandaren
		var foodBuffStats stats.Stats
		if food.BuffsMainStat {
			buffAmount := TernaryFloat64(isPanda, food.Stats[stats.Stamina]*2, food.Stats[stats.Stamina])
			foodBuffStats[stats.Stamina] = buffAmount
			foodBuffStats[character.GetHighestStatType([]stats.Stat{stats.Strength, stats.Agility, stats.Intellect})] = buffAmount
		} else {
			if isPanda {
				for stat, amount := range food.Stats {
					food.Stats[stat] = amount * 2
				}
			}
			foodBuffStats = food.Stats
		}
		character.AddStats(foodBuffStats)
	}

	registerPotionCD(agent, consumables)
	registerConjuredCD(agent, consumables)
	registerExplosivesCD(agent, consumables)
}

var PotionAuraTag = "Potion"

func registerPotionCD(agent Agent, consumes *proto.ConsumesSpec) {
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

var AlchStoneItemIDs = []int32{136197, 80508, 96252, 96253, 96254, 44322, 44323, 44324}

func (character *Character) HasAlchStone() bool {
	alchStoneEquipped := false
	for _, itemID := range AlchStoneItemIDs {
		alchStoneEquipped = alchStoneEquipped || character.HasTrinketEquipped(itemID)
	}
	return character.HasProfession(proto.Profession_Alchemy) && alchStoneEquipped
}

func makePotionActivationSpell(potionId int32, character *Character, potionCD *Timer) MajorCooldown {
	potion := ConsumablesByID[potionId]
	mcd := makePotionActivationSpellInternal(potion, character, potionCD)
	cooldownDuration := TernaryDuration(potion.CooldownDuration > 0, potion.CooldownDuration, time.Minute*1)

	if mcd.Spell != nil {
		// Mark as 'Encounter Only' so that users are forced to select the generic Potion
		// placeholder action instead of specific potion spells, in APL prepull. This
		// prevents a mismatch between Consumes and Rotation settings.
		mcd.Spell.Flags |= SpellFlagEncounterOnly | SpellFlagPotion
		oldApplyEffects := mcd.Spell.ApplyEffects
		mcd.Spell.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			oldApplyEffects(sim, target, spell)
			if sim.CurrentTime < 0 {
				potionCD.Set(sim.CurrentTime + cooldownDuration)

				character.UpdateMajorCooldowns()
			}
		}
	}
	return mcd

}

type resourceGainConfig struct {
	resType proto.ResourceType
	min     float64
	spread  float64
}

func makePotionActivationSpellInternal(potion Consumable, character *Character, potionCD *Timer) MajorCooldown {
	stoneMul := TernaryFloat64(character.HasAlchStone(), 1.4, 1.0)

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

	actionID := ActionID{ItemID: potion.Id}
	var aura *StatBuffAura
	mcd := MajorCooldown{
		Spell: character.GetOrRegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast:     potionCast,
		}),
	}
	if potion.BuffDuration > 0 {
		// Add stat buff aura if applicable
		aura = character.NewTemporaryStatsAura(potion.Name, actionID, potion.Stats, potion.BuffDuration)
		mcd.Type = aura.InferCDType()
	}
	var gains []resourceGainConfig
	resourceMetrics := make(map[proto.ResourceType]*ResourceMetrics)

	for _, effectID := range potion.EffectIds {
		e := SpellEffectsById[effectID]
		resourceType := e.GetResourceType()
		if e.Type == proto.EffectType_EffectTypeResourceGain && resourceType != 0 {
			if resourceType == proto.ResourceType_ResourceTypeMana && mcd.Type != CooldownTypeSurvival {
				mcd.Type = CooldownTypeMana
			} else if resourceType == proto.ResourceType_ResourceTypeHealth {
				mcd.Type = CooldownTypeSurvival
			} else {
				mcd.Type = CooldownTypeDPS
			}
			gains = append(gains, resourceGainConfig{
				resType: resourceType,
				min:     e.MinEffectSize,
				spread:  e.EffectSpread,
			})
			if _, exists := resourceMetrics[resourceType]; !exists {
				resourceMetrics[resourceType] = character.Metrics.NewResourceMetrics(actionID, resourceType)
			}
			// Preload resource types that are found on this item
			if resourceMetrics[resourceType] == nil {
				resourceMetrics[resourceType] = character.Metrics.NewResourceMetrics(actionID, resourceType)
			}
		}
	}

	mcd.Spell.ApplyEffects = func(sim *Simulation, _ *Unit, _ *Spell) {
		if aura != nil {
			aura.Activate(sim)
		}
		for _, config := range gains {
			gain := config.min + sim.RandomFloat(potion.Name)*config.spread
			gain *= stoneMul
			if config.resType == proto.ResourceType_ResourceTypeHealth {
				gain *= character.PseudoStats.HealingTakenMultiplier
			}
			character.ExecuteResourceGain(sim, config.resType, gain, resourceMetrics[config.resType])
		}
	}

	mcd.ShouldActivate = func(sim *Simulation, character *Character) bool {
		shouldActivate := true
		for _, config := range gains {
			switch config.resType {
			case proto.ResourceType_ResourceTypeMana:
				totalRegen := character.ManaRegenPerSecondWhileCombat() * 5
				manaGain := config.min + config.spread
				manaGain *= stoneMul
				shouldActivate = character.MaxMana()-(character.CurrentMana()+totalRegen) >= manaGain
			}
		}
		return shouldActivate
	}

	return mcd

}

var ConjuredAuraTag = "Conjured"

func registerConjuredCD(agent Agent, consumes *proto.ConsumesSpec) {
	character := agent.GetCharacter()

	//Todo: Implement dynamic handling like pots etc.
	switch consumes.ConjuredId {
	case 20520:
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
	case 5512:
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

func registerExplosivesCD(agent Agent, consumes *proto.ConsumesSpec) {
	//Todo: Get them dynamically from dbc data
	character := agent.GetCharacter()
	if !character.HasProfession(proto.Profession_Engineering) {
		return
	}
	switch consumes.ExplosiveId {
	case 89637:
		bomb := character.GetOrRegisterSpell(SpellConfig{
			ActionID:    BigDaddyActionID,
			SpellSchool: SpellSchoolFire,
			ProcMask:    ProcMaskEmpty,
			Flags:       SpellFlagAoE,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute,
				},

				DefaultCast: Cast{
					CastTime: time.Millisecond * 500,
				},

				ModifyCast: func(sim *Simulation, spell *Spell, cast *Cast) {
					spell.Unit.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime)
					spell.Unit.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime)
				},
			},

			// Explosives always have 1% resist chance, so just give them hit cap.
			BonusHitPercent:  100,
			DamageMultiplier: 1,
			CritMultiplier:   2,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, 5006, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    bomb,
			Type:     CooldownTypeDPS | CooldownTypeExplosive,
			Priority: CooldownPriorityLow + 10,
		})
	case 40771:
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
