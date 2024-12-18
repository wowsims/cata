package core

import (
	"fmt"
	"slices"

	"github.com/wowsims/cata/sim/core/proto"
)

type ApplySetItemEffect func(agent Agent, setName string)
type ItemSet struct {
	ID              int32
	Name            string
	AlternativeName string

	// Maps set piece requirement to an ApplyEffect function that will be called
	// before the Sim starts.
	//
	// The function should apply any benefits provided by the set bonus.
	Bonuses map[int32]ApplySetItemEffect
}

var ItemSetSlots = []proto.ItemSlot{
	proto.ItemSlot_ItemSlotHead,
	proto.ItemSlot_ItemSlotShoulder,
	proto.ItemSlot_ItemSlotChest,
	proto.ItemSlot_ItemSlotHands,
	proto.ItemSlot_ItemSlotLegs,
}

func (set ItemSet) Items() []Item {
	var items []Item
	for _, item := range ItemsByID {
		if item.SetName == "" {
			continue
		}
		if item.SetName == set.Name || item.SetName == set.AlternativeName {
			items = append(items, item)
		}
	}
	// Sort so the order of IDs is always consistent, for tests.
	slices.SortFunc(items, func(a, b Item) int {
		return int(a.ID - b.ID)
	})
	return items
}

var sets []*ItemSet

// Registers a new ItemSet with item IDs populated.
func NewItemSet(set ItemSet) *ItemSet {
	foundID := set.ID == 0
	foundName := false
	foundAlternativeName := set.AlternativeName == ""
	for _, item := range ItemsByID {
		if item.SetName == "" {
			continue
		}
		foundID = foundID || (item.SetID > 0 && item.SetID == set.ID)
		foundName = foundName || item.SetName == set.Name
		foundAlternativeName = foundAlternativeName || item.SetName == set.AlternativeName
		if foundID && foundName && foundAlternativeName {
			break
		}
	}

	if WITH_DB {
		if !foundID {
			panic(fmt.Sprintf("No items found for set id %d", set.ID))
		}
		if !foundName {
			panic("No items found for set " + set.Name)
		}
		if len(set.AlternativeName) > 0 && !foundAlternativeName {
			panic("No items found for set alternative " + set.AlternativeName)
		}
	}

	sets = append(sets, &set)
	return &set
}

func (character *Character) HasSetBonus(set *ItemSet, numItems int32) bool {
	if character.Env != nil && character.Env.IsFinalized() {
		panic("HasSetBonus is very slow and should never be called after finalization. Try caching the value during construction instead!")
	}

	if _, ok := set.Bonuses[numItems]; !ok {
		panic(fmt.Sprintf("Item set %s does not have a bonus with %d pieces.", set.Name, numItems))
	}

	var count int32
	for _, item := range character.Equipment {
		if item.SetName == "" {
			continue
		}
		if item.SetName == set.Name || item.SetName == set.AlternativeName || (item.SetID > 0 && item.SetID == set.ID) {
			count++
			if count >= numItems {
				return true
			}
		}
	}

	return false
}

type SetBonus struct {
	// Name of the set.
	Name string

	// Number of pieces required for this bonus.
	NumPieces int32

	// Function for applying the effects of this set bonus.
	BonusEffect ApplySetItemEffect
}

// Returns a list describing all active set bonuses.
func (character *Character) GetActiveSetBonuses() []SetBonus {
	return character.GetSetBonuses(character.Equipment)
}

func (character *Character) GetSetBonuses(equipment Equipment) []SetBonus {
	var activeBonuses []SetBonus

	setItemCount := make(map[*ItemSet]int32)
	for _, item := range equipment {
		if item.SetName == "" {
			continue
		}

		var foundSet *ItemSet = nil

		if item.SetID > 0 {
			// Try finding by ID first to make sure sets with different names but share id all point to the same count.
			for _, set := range sets {
				if set.ID == item.SetID {
					foundSet = set
					break
				}
			}
		}

		if foundSet == nil {
			for _, set := range sets {
				if set.Name == item.SetName || set.AlternativeName == item.SetName {
					foundSet = set
					break
				}
			}
		}

		if foundSet != nil {
			setItemCount[foundSet]++
			if bonusEffect, ok := foundSet.Bonuses[setItemCount[foundSet]]; ok {
				activeBonuses = append(activeBonuses, SetBonus{
					Name:        foundSet.Name,
					NumPieces:   setItemCount[foundSet],
					BonusEffect: bonusEffect,
				})
			}
		}
	}

	return activeBonuses
}

func (character *Character) HasActiveSetBonus(setName string, count int32) bool {
	activeSetBonuses := character.GetActiveSetBonuses()

	for _, activeSetBonus := range activeSetBonuses {
		if activeSetBonus.Name == setName && activeSetBonus.NumPieces >= count {
			return true
		}
	}

	return false
}

// Apply effects from item set bonuses.
func (character *Character) applyItemSetBonusEffects(agent Agent) {
	activeSetBonuses := character.GetActiveSetBonuses()

	for _, activeSetBonus := range activeSetBonuses {
		activeSetBonus.BonusEffect(agent, activeSetBonus.Name)
	}

	if character.ItemSwap.IsEnabled() {
		unequippedSetBonuses := FilterSlice(character.GetSetBonuses(character.ItemSwap.unEquippedItems), func(unequippeBonus SetBonus) bool {
			return !character.HasActiveSetBonus(unequippeBonus.Name, unequippeBonus.NumPieces)
		})

		for _, unequippedSetBonus := range unequippedSetBonuses {
			unequippedSetBonus.BonusEffect(agent, unequippedSetBonus.Name)
		}
	}
}

// Returns the names of all active set bonuses.
func (character *Character) GetActiveSetBonusNames() []string {
	activeSetBonuses := character.GetActiveSetBonuses()

	names := make([]string, len(activeSetBonuses))
	for i, activeSetBonus := range activeSetBonuses {
		names[i] = fmt.Sprintf("%s (%dpc)", activeSetBonus.Name, activeSetBonus.NumPieces)
	}
	return names
}

type CustomSetBonusCallbackConfig struct {
	// Override default behavior when the set bonus is gained.
	OnGain func(sim *Simulation, aura *Aura)
	// Override default behavior when the set bonus is lost.
	OnExpire func(sim *Simulation, aura *Aura)
}

// Adds a proc trigger aura that activates when the character has a set bonus.
func (character *Character) MakeProcTriggerAuraForSetBonus(setName string, numPieces int32, config ProcTrigger, customSetBonusCallbackConfig *CustomSetBonusCallbackConfig) *Aura {
	return character.factory_ProcTriggerAuraForSetBonus(setName, numPieces, &character.Unit, config, customSetBonusCallbackConfig)
}

func (character *Character) MakeProcTriggerAuraForSetBonusWithUnit(setName string, numPieces int32, unit *Unit, config ProcTrigger, customSetBonusCallbackConfig *CustomSetBonusCallbackConfig) *Aura {
	return character.factory_ProcTriggerAuraForSetBonus(setName, numPieces, unit, config, customSetBonusCallbackConfig)
}

func (character *Character) factory_ProcTriggerAuraForSetBonus(setName string, numPieces int32, unit *Unit, config ProcTrigger, customSetBonusCallbackConfig *CustomSetBonusCallbackConfig) *Aura {
	aura := MakeProcTriggerAura(unit, config)

	if customSetBonusCallbackConfig == nil {
		customSetBonusCallbackConfig = &CustomSetBonusCallbackConfig{}
	}

	character.registerSetBonusForItemSwap(setName, numPieces, ItemSwapRegistrationConfig{
		OnGain: func(sim *Simulation, _ proto.ItemSlot) {
			if customSetBonusCallbackConfig.OnGain != nil {
				customSetBonusCallbackConfig.OnGain(sim, aura)
			} else {
				aura.Activate(sim)
			}
		},
		OnExpire: func(sim *Simulation, _ proto.ItemSlot) {
			if customSetBonusCallbackConfig.OnExpire != nil {
				customSetBonusCallbackConfig.OnExpire(sim, aura)
			} else {
				aura.Deactivate(sim)
			}
		},
	})

	return aura
}

// Adds a static effect that activates when the character has a set bonus.
func (character *Character) MakeCallbackEffectForSetBonus(setName string, numPieces int32, callbackConfig CustomSetBonusCallbackConfig) {
	character.registerSetBonusForItemSwap(setName, numPieces, ItemSwapRegistrationConfig{
		OnGain: func(sim *Simulation, _ proto.ItemSlot) {
			if callbackConfig.OnGain != nil {
				callbackConfig.OnGain(sim, nil)
			}
		},
		OnExpire: func(sim *Simulation, _ proto.ItemSlot) {
			if callbackConfig.OnExpire != nil {
				callbackConfig.OnExpire(sim, nil)
			}
		},
		OnDisabled: func() {
			if callbackConfig.OnGain != nil {
				callbackConfig.OnGain(nil, nil)
			}
		},
	})
}

// Adds a static effect that activates when the character has a set bonus.
func (character *Character) MakeDynamicModForSetBonus(setName string, numPieces int32, spellModConfig SpellModConfig) {

	setBonusDep := character.AddDynamicMod(spellModConfig)
	character.registerSetBonusForItemSwap(setName, numPieces, ItemSwapRegistrationConfig{
		OnGain: func(_ *Simulation, _ proto.ItemSlot) {
			setBonusDep.Activate()
		},
		OnExpire: func(_ *Simulation, _ proto.ItemSlot) {
			setBonusDep.Deactivate()
		},
		OnDisabled: func() {
			setBonusDep.Activate()
		},
	})
}

type ItemSwapRegistrationConfig struct {
	// Will be called when your set bonus is gained.
	OnGain func(sim *Simulation, slot proto.ItemSlot)
	// Will be called when your set bonus is lost.
	OnExpire func(sim *Simulation, slot proto.ItemSlot)
	// Will be called when ItemSwap is disabled.
	OnDisabled func()
}

func (character *Character) registerSetBonusForItemSwap(setName string, numPieces int32, config ItemSwapRegistrationConfig) {
	if character.ItemSwap.IsEnabled() {
		character.RegisterItemSwapCallback(ItemSetSlots, func(sim *Simulation, slot proto.ItemSlot) {
			if character.HasActiveSetBonus(setName, numPieces) {
				if config.OnGain != nil {
					config.OnGain(sim, slot)
				}
			} else {
				if config.OnExpire != nil {
					config.OnExpire(sim, slot)
				}
			}
		})
	} else {
		if config.OnDisabled != nil {
			config.OnDisabled()
		}
	}
}
