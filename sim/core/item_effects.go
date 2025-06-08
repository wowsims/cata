package core

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// Function for applying permanent effects to an Agent.
//
// Passing Character instead of Agent would work for almost all cases,
// but there are occasionally class-specific item effects.
type ApplyEffect func(Agent, proto.ItemLevelState)

var itemEffects = map[int32]ApplyEffect{}
var enchantEffects = map[int32]ApplyEffect{}

// IDs of item effects which should be used for tests.
var itemEffectsForTest []int32
var enchantEffectsForTest []int32

// This value can be set before adding item effects, to control whether they are included in tests.
var AddEffectsToTest = true

func HasItemEffect(id int32) bool {
	_, ok := itemEffects[id]
	return ok
}
func HasItemEffectForTest(id int32) bool {
	return slices.Contains(itemEffectsForTest, id)
}

func HasEnchantEffect(id int32) bool {
	_, ok := enchantEffects[id]
	return ok
}

// Registers an ApplyEffect function which will be called before the Sim
// starts, for any Agent that is wearing the item.
func NewItemEffect(id int32, itemEffect ApplyEffect) {
	if WITH_DB {
		if _, hasItem := ItemsByID[id]; !hasItem {
			if _, hasGem := GemsByID[id]; !hasGem {
				panic(fmt.Sprintf("No item with ID: %d", id))
			}
		}
	}

	if HasItemEffect(id) {
		panic(fmt.Sprintf("Cannot add multiple effects for one item: %d, %#v", id, itemEffect))
	}

	itemEffects[id] = itemEffect
	if AddEffectsToTest {
		itemEffectsForTest = append(itemEffectsForTest, id)
	}
}

func NewEnchantEffect(id int32, enchantEffect ApplyEffect) {
	if WITH_DB {
		if _, ok := EnchantsByEffectID[id]; !ok {
			panic(fmt.Sprintf("No enchant with ID: %d", id))
		}
	}

	if HasEnchantEffect(id) {
		panic(fmt.Sprintf("Cannot add multiple effects for one enchant: %d, %#v", id, enchantEffect))
	}

	enchantEffects[id] = enchantEffect
	if AddEffectsToTest {
		enchantEffectsForTest = append(enchantEffectsForTest, id)
	}
}

func (equipment *Equipment) applyItemEffects(agent Agent, registeredItemEffects map[int32]bool, registeredItemEnchantEffects map[int32]bool, includeGemEffects bool) {
	for _, eq := range equipment {
		if applyItemEffect, ok := itemEffects[eq.ID]; ok && !registeredItemEffects[eq.ID] {
			applyItemEffect(agent, eq.UpgradeStep)
			registeredItemEffects[eq.ID] = true
		}

		if includeGemEffects {
			for _, g := range eq.Gems {
				if applyGemEffect, ok := itemEffects[g.ID]; ok {
					applyGemEffect(agent, proto.ItemLevelState_Base)
				}
			}
		}

		if applyEnchantEffect, ok := enchantEffects[eq.Enchant.EffectID]; ok && !registeredItemEnchantEffects[eq.Enchant.EffectID] {
			applyEnchantEffect(agent, proto.ItemLevelState_Base)
			registeredItemEnchantEffects[eq.Enchant.EffectID] = true
		}
	}
}

// Helpers for making common types of active item effects.

func NewSimpleStatItemActiveEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration, sharedCDFunc func(*Character) Cooldown, otherEffects ApplyEffect) {
	NewItemEffect(itemID, func(agent Agent, _ proto.ItemLevelState) {
		registerCD := MakeTemporaryStatsOnUseCDRegistration(
			"ItemActive-"+strconv.Itoa(int(itemID)),
			bonus,
			duration,
			SpellConfig{
				ActionID: ActionID{ItemID: itemID},
			},
			func(character *Character) Cooldown {
				return Cooldown{
					Timer:    character.NewTimer(),
					Duration: cooldown,
				}
			},
			sharedCDFunc,
		)

		registerCD(agent, proto.ItemLevelState_Base)
		if otherEffects != nil {
			otherEffects(agent, proto.ItemLevelState_Base)
		}
	})
}

// No shared CD
func NewSimpleStatItemEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, func(character *Character) Cooldown {
		return Cooldown{}
	}, nil)
}

func NewSimpleStatOffensiveTrinketEffectWithOtherEffects(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration, otherEffects ApplyEffect) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, func(character *Character) Cooldown {
		return Cooldown{
			Timer:    character.GetOffensiveTrinketCD(),
			Duration: duration,
		}
	}, otherEffects)
}
func NewSimpleStatOffensiveTrinketEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatOffensiveTrinketEffectWithOtherEffects(itemID, bonus, duration, cooldown, nil)
}

func NewSimpleStatDefensiveTrinketEffect(itemID int32, bonus stats.Stats, duration time.Duration, cooldown time.Duration) {
	NewSimpleStatItemActiveEffect(itemID, bonus, duration, cooldown, func(character *Character) Cooldown {
		return Cooldown{
			Timer:    character.GetDefensiveTrinketCD(),
			Duration: duration,
		}
	}, nil)
}
