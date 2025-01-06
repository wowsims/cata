package core

import (
	"fmt"
	"slices"

	"github.com/wowsims/cata/sim/core/proto"
)

type ApplySetBonus func(agent Agent, setBonusAura *Aura)

type ItemSet struct {
	ID              int32
	Name            string
	AlternativeName string

	// Maps set piece requirement to an ApplyEffect function that will be called
	// before the Sim starts.
	//
	// The function should apply any benefits provided by the set bonus.
	Bonuses map[int32]ApplySetBonus

	// Optional field to override the DefaultItemSetSlots
	// For Example: The set contains of 2 weapons
	Slots []proto.ItemSlot
}

var DefaultItemSetSlots = []proto.ItemSlot{
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

	if len(set.Slots) == 0 {
		set.Slots = DefaultItemSetSlots
	}

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
	BonusEffect ApplySetBonus

	// Optional field to override the DefaultItemSetSlots
	// For Example: The set contains of 2 weapons
	Slots []proto.ItemSlot
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
					Slots:       foundSet.Slots,
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
		setBonusAura := character.makeSetBonusStatusAura(activeSetBonus.Name, activeSetBonus.NumPieces, activeSetBonus.Slots, true)
		activeSetBonus.BonusEffect(agent, setBonusAura)
	}

	if character.ItemSwap.IsEnabled() {
		unequippedSetBonuses := FilterSlice(character.GetSetBonuses(character.ItemSwap.unEquippedItems), func(unequippedBonus SetBonus) bool {
			return !character.HasActiveSetBonus(unequippedBonus.Name, unequippedBonus.NumPieces)
		})

		for _, unequippedSetBonus := range unequippedSetBonuses {
			setBonusAura := character.makeSetBonusStatusAura(unequippedSetBonus.Name, unequippedSetBonus.NumPieces, unequippedSetBonus.Slots, false)
			unequippedSetBonus.BonusEffect(agent, setBonusAura)
		}
	}
}

func (character *Character) makeSetBonusStatusAura(setName string, numPieces int32, slots []proto.ItemSlot, activeAtStart bool) *Aura {
	statusAura := character.GetOrRegisterAura(Aura{
		Label:      fmt.Sprintf("%s %dP", setName, numPieces),
		BuildPhase: Ternary(activeAtStart, CharacterBuildPhaseGear, CharacterBuildPhaseNone),
		Duration:   NeverExpires,
	})

	if activeAtStart {
		statusAura = MakePermanent(statusAura)
	}

	if character.ItemSwap.IsEnabled() {
		character.RegisterItemSwapCallback(slots, func(sim *Simulation, _ proto.ItemSlot) {
			if character.HasActiveSetBonus(setName, numPieces) {
				statusAura.Activate(sim)
			} else {
				statusAura.Deactivate(sim)
			}
		})
	}

	return statusAura
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

// Adds a spellID to the set bonus so it can be exposed to the APL
func (setBonusTracker *Aura) ExposeToAPL(spellID int32) {
	setBonusTracker.ActionID = ActionID{SpellID: spellID}
}

// Adds a Spellmod to PVP GLoves
func (character *Character) RegisterPvPGloveMod(itemIDs []int32, config SpellModConfig) {
	spellMod := character.AddDynamicMod(config)

	checkGloves := func() {
		if slices.Contains(itemIDs, character.Hands().ID) {
			spellMod.Activate()
		} else {
			spellMod.Deactivate()
		}
	}

	if character.ItemSwap.IsEnabled() {
		character.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotHands}, func(_ *Simulation, _ proto.ItemSlot) {
			checkGloves()
		})
	} else {
		checkGloves()
	}
}
