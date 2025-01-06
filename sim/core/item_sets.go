package core

import (
	"fmt"
	"slices"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
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
	BonusEffect ApplySetBonus
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
		setBonusAura := character.makeSetBonusStatusAura(activeSetBonus.Name, activeSetBonus.NumPieces, true)
		activeSetBonus.BonusEffect(agent, setBonusAura)
	}

	if character.ItemSwap.IsEnabled() {
		unequippedSetBonuses := FilterSlice(character.GetSetBonuses(character.ItemSwap.unEquippedItems), func(unequippedBonus SetBonus) bool {
			return !character.HasActiveSetBonus(unequippedBonus.Name, unequippedBonus.NumPieces)
		})

		for _, unequippedSetBonus := range unequippedSetBonuses {
			setBonusAura := character.makeSetBonusStatusAura(unequippedSetBonus.Name, unequippedSetBonus.NumPieces, false)
			unequippedSetBonus.BonusEffect(agent, setBonusAura)
		}
	}
}

func (character *Character) makeSetBonusStatusAura(setName string, numPieces int32, activeAtStart bool) *Aura {
	statusAura := character.GetOrRegisterAura(Aura{
		Label:      fmt.Sprintf("%s %dP", setName, numPieces),
		BuildPhase: Ternary(activeAtStart, CharacterBuildPhaseGear, CharacterBuildPhaseNone),
		Duration:   NeverExpires,
	})

	if activeAtStart {
		statusAura = MakePermanent(statusAura)
	}

	if character.ItemSwap.IsEnabled() {
		character.RegisterItemSwapCallback(ItemSetSlots, func(sim *Simulation, _ proto.ItemSlot) {
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

// Creates a new ProcTriggerAura that is dependent on the set bonus being active
// This should only be used if the dependent Aura is:
// 1. On the a different Unit than the setBonus is registered to (usually the Character)
// 2. You need to register multiple dependent Aura's for the same Unit
func (setBonusTracker *Aura) MakeDependentProcTriggerAura(unit *Unit, config ProcTrigger) *Aura {
	oldExtraCondition := config.ExtraCondition
	config.ExtraCondition = func(sim *Simulation, spell *Spell, result *SpellResult) bool {
		return setBonusTracker.IsActive() && ((oldExtraCondition == nil) || oldExtraCondition(sim, spell, result))
	}

	aura := MakeProcTriggerAura(unit, config)

	return aura
}

// Attaches a ProcTrigger to the set bonus
// Preffered use-case.
// For non standard use-cases see: MakeDependentProcTriggerAura
func (setBonusTracker *Aura) AttachProcTrigger(config ProcTrigger) {
	ApplyProcTriggerCallback(setBonusTracker.Unit, setBonusTracker, config)
}

// Attaches a SpellMod to the set bonus
func (setBonusTracker *Aura) AttachSpellMod(spellModConfig SpellModConfig) {
	setBonusDep := setBonusTracker.Unit.AddDynamicMod(spellModConfig)

	setBonusTracker.ApplyOnGain(func(_ *Aura, _ *Simulation) {
		setBonusDep.Activate()
	})

	setBonusTracker.ApplyOnExpire(func(_ *Aura, _ *Simulation) {
		setBonusDep.Deactivate()
	})
}

// Adds Stats to the set bonus
func (setBonusTracker *Aura) AttachStatsBuff(stats stats.Stats) {
	setBonusTracker.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		aura.Unit.AddStatsDynamic(sim, stats)
	})

	setBonusTracker.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		aura.Unit.AddStatsDynamic(sim, stats.Invert())
	})
}

// Adds a Stat to the set bonus
func (setBonusTracker *Aura) AttachStatBuff(stat stats.Stat, value float64) {
	statsToAdd := stats.Stats{}
	statsToAdd[stat] = value
	setBonusTracker.AttachStatsBuff(statsToAdd)
}
