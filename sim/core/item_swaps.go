package core

import (
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type OnItemSwap func(*Simulation, proto.ItemSlot)

type ItemSwap struct {
	character       *Character
	onSwapCallbacks [NumItemSlots][]OnItemSwap

	mhCritMultiplier     float64
	ohCritMultiplier     float64
	rangedCritMultiplier float64

	// Which slots to actually swap.
	slots []proto.ItemSlot

	// Holds the original equip
	originalEquip Equipment
	// Holds the items that are selected for swapping
	swapEquip Equipment
	// Holds items that are currently not equipped
	unEquippedItems Equipment
	swapSet         proto.APLActionItemSwap_SwapSet
	initialized     bool
}

/**
 * TODO All the extra parameters here and the code in multiple places for handling the Weapon struct is really messy,
 * we'll need to figure out something cleaner as this will be quite error-prone
**/
func (character *Character) enableItemSwap(itemSwap *proto.ItemSwap, mhCritMultiplier float64, ohCritMultiplier float64, rangedCritMultiplier float64) {
	var slots []proto.ItemSlot
	var swapItems Equipment
	hasItemSwap := make(map[proto.ItemSlot]bool)

	for slot, item := range itemSwap.Items {
		itemSlot := proto.ItemSlot(slot)
		hasItemSwap[itemSlot] = item != nil && item.Id != 0
		swapItems[itemSlot] = toItem(item)
	}

	has2H := swapItems[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand
	hasMh := character.HasMHWeapon()
	hasOh := character.HasOHWeapon()

	// Handle MH and OH together, because present MH + empty OH --> swap MH and unequip OH
	hasItemSwap = FilterMap(hasItemSwap, func(itemSlot proto.ItemSlot, v bool) bool {
		if itemSlot == proto.ItemSlot_ItemSlotMainHand || itemSlot == proto.ItemSlot_ItemSlotOffHand {
			if itemSlot == proto.ItemSlot_ItemSlotMainHand || (itemSlot == proto.ItemSlot_ItemSlotOffHand && hasMh) {
				return true
			} else if itemSlot == proto.ItemSlot_ItemSlotOffHand || (has2H && hasOh) {
				return true
			} else {
				return false
			}
		} else {
			return true
		}
	})

	for slot, hasSlotItemSwap := range hasItemSwap {
		if hasSlotItemSwap {
			slots = append(slots, slot)
		}
	}

	if len(slots) == 0 {
		return
	}

	character.ItemSwap = ItemSwap{
		mhCritMultiplier:     mhCritMultiplier,
		ohCritMultiplier:     ohCritMultiplier,
		rangedCritMultiplier: rangedCritMultiplier,
		slots:                slots,
		originalEquip:        character.Equipment,
		swapEquip:            swapItems,
		unEquippedItems:      swapItems,
		swapSet:              proto.APLActionItemSwap_Unknown,
		initialized:          false,
	}
}

func (swap *ItemSwap) initialize(character *Character) {
	swap.character = character
}

func (character *Character) RegisterItemSwapCallback(slots []proto.ItemSlot, callback OnItemSwap) {
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}

	var slot proto.ItemSlot
	for _, slot = range slots {
		character.ItemSwap.onSwapCallbacks[slot] = append(character.ItemSwap.onSwapCallbacks[slot], callback)
	}
}

// Helper for handling Effects that use PPMManager to toggle the aura on/off
func (swap *ItemSwap) RegisterPPMEffect(effectID int32, ppm float64, ppmm *PPMManager, aura *Aura) {
	character := swap.character
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}
	character.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}, func(sim *Simulation, slot proto.ItemSlot) {
		procMask := character.GetProcMaskForEnchant(effectID)
		*ppmm = character.AutoAttacks.NewPPMManager(ppm, procMask)

		if ppmm.Chance(procMask) == 0 {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

// Helper for handling procs that use PPMManager to toggle the aura on/off
func (swap *ItemSwap) RegisterPPMEffectWithCustomProcMask(procMask ProcMask, ppm float64, ppmm *PPMManager) {
	character := swap.character
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}
	character.RegisterItemSwapCallback([]proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand, proto.ItemSlot_ItemSlotOffHand}, func(sim *Simulation, slot proto.ItemSlot) {
		*ppmm = character.AutoAttacks.NewPPMManager(ppm, procMask)
	})
}

// Helper for handling Item Effects that use the itemID to toggle the aura on and off
func (swap *ItemSwap) RegisterProc(itemID int32, aura *Aura, slots []proto.ItemSlot) {
	character := swap.character
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}
	character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		hasItemEquipped := swap.HasItemEquipped(itemID)
		if hasItemEquipped {
			if !aura.IsActive() {
				aura.Activate(sim)
				if aura.Icd != nil {
					aura.Icd.Use(sim)
				}
			}
		} else {
			aura.Deactivate(sim)
		}
	})
}

// Helper for handling Enchant Effects that use the effectID to toggle the aura on and off
func (swap *ItemSwap) RegisterEnchantProc(effectID int32, aura *Aura, slots []proto.ItemSlot) {
	character := swap.character
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}
	character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		procMask := character.GetProcMaskForEnchant(effectID)

		if procMask == ProcMaskUnknown {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

// Helper for handling Item On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) RegisterActive(itemID int32, slots []proto.ItemSlot) {
	character := swap.character
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}
	character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		isSwapItem := swap.ItemExistsInSwapSet(itemID)
		if !isSwapItem {
			return
		}
		hasItemEquipped := swap.HasItemEquipped(itemID)
		itemSlot := swap.GetItemSwapItemSlot(itemID)
		if itemSlot == -1 {
			return
		}
		equippedItemID := swap.GetEquippedItemBySlot(itemSlot).ID
		spell := swap.character.GetSpell(ActionID{ItemID: itemID})
		if spell != nil {
			aura := character.GetAuraByID(spell.ActionID)
			if aura.IsActive() {
				aura.Deactivate(sim)
			}
			if !hasItemEquipped {
				spell.Flags |= SpellFlagSwapped
				return
			}
			spell.Flags &= ^SpellFlagSwapped
			if !swap.initialized {
				return
			}
			swappedItemID := swap.GetUnequippedItemBySlot(slot).ID

			if swappedItemID == equippedItemID && spell.CD.IsReady(sim) || swappedItemID != equippedItemID {
				spell.CD.Set(sim.CurrentTime + time.Second*30)
			}
		}
	})
}

// Helper for handling Enchant On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) ProcessTinker(spell *Spell, slots []proto.ItemSlot) {
	character := swap.character
	if character == nil || !character.ItemSwap.IsEnabled() {
		return
	}
	character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		if spell == nil || !swap.initialized {
			return
		}
		equippedItemID := swap.GetEquippedItemBySlot(slot).ID
		swappedItemID := swap.GetUnequippedItemBySlot(slot).ID
		if swappedItemID == equippedItemID && spell.CD.IsReady(sim) || swappedItemID != equippedItemID {
			spell.CD.Set(sim.CurrentTime + time.Second*30)
		}
	})
}

func (swap *ItemSwap) IsEnabled() bool {
	return swap.character != nil && len(swap.slots) > 0
}

func (swap *ItemSwap) IsSwapped() bool {
	return swap.swapSet == proto.APLActionItemSwap_Swap1
}

func (swap *ItemSwap) HasItemEquipped(itemID int32) bool {
	for _, item := range swap.character.Equipment {
		if item.ID == itemID {
			return true
		}
	}
	return false
}

func (swap *ItemSwap) GetEquippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.character.Equipment[slot]
}

func (swap *ItemSwap) GetUnequippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.unEquippedItems[slot]
}

func (swap *ItemSwap) GetItemSwapItemSlot(itemID int32) proto.ItemSlot {
	slotsToCheck := Ternary(swap.IsSwapped(), swap.swapEquip, swap.originalEquip)
	for slot, item := range slotsToCheck {
		if item.ID == itemID {
			return proto.ItemSlot(slot)
		}
	}
	return -1
}

func (swap *ItemSwap) ItemExistsInSwapSet(itemID int32) bool {
	for _, item := range swap.unEquippedItems {
		if item.ID == itemID {
			return true
		}
	}
	return false
}

func (swap *ItemSwap) CalcStatChanges(slots []proto.ItemSlot) stats.Stats {
	newStats := stats.Stats{}
	for _, slot := range slots {
		oldItemStats := swap.getItemStats(*swap.GetEquippedItemBySlot(slot))
		newItemStats := swap.getItemStats(*swap.GetUnequippedItemBySlot(slot))
		newStats = newStats.Add(newItemStats.Subtract(oldItemStats))
	}

	return newStats
}

func (swap *ItemSwap) SwapItems(sim *Simulation, swapSet proto.APLActionItemSwap_SwapSet, slots []proto.ItemSlot, isReset bool) {
	if !swap.IsEnabled() || swap.swapSet == swapSet {
		return
	}

	character := swap.character

	meleeWeaponSwapped := false
	newStats := stats.Stats{}
	has2H := swap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotMainHand).HandType == proto.HandType_HandTypeTwoHand
	isPrepull := sim.CurrentTime < 0
	for _, slot := range slots {
		if !isReset && !isPrepull && (slot < proto.ItemSlot_ItemSlotMainHand || slot > proto.ItemSlot_ItemSlotRanged) {
			continue
		}

		//will swap both on the MainHand Slot for 2H.
		if slot == proto.ItemSlot_ItemSlotOffHand && has2H {
			continue
		}

		if ok, swapStats := swap.swapItem(slot, has2H, isReset); ok {
			newStats = newStats.Add(swapStats)
			meleeWeaponSwapped = slot == proto.ItemSlot_ItemSlotMainHand || slot == proto.ItemSlot_ItemSlotOffHand || slot == proto.ItemSlot_ItemSlotRanged || meleeWeaponSwapped

			for _, onSwap := range swap.onSwapCallbacks[slot] {
				onSwap(sim, slot)
			}
		}
	}

	if sim.Log != nil {
		sim.Log("Item Swap - New Stats: %v", newStats.FlatString())
	}
	character.AddStatsDynamic(sim, newStats)

	if !isPrepull && !isReset {
		if character.AutoAttacks.AutoSwingMelee && meleeWeaponSwapped {
			character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		}

		// If GCD is ready then use the GCD, otherwise we assume it's being used along side a spell.
		if character.GCD.IsReady(sim) {
			character.ExtendGCDUntil(sim, max(character.NextGCDAt(), sim.CurrentTime+GCDDefault))
		}
	}

	swap.swapSet = swapSet
}

func (swap *ItemSwap) swapItem(slot proto.ItemSlot, has2H bool, isReset bool) (bool, stats.Stats) {
	oldItem := *swap.GetEquippedItemBySlot(slot)
	var newItem *Item
	if isReset {
		newItem = &swap.originalEquip[slot]
	} else {
		newItem = swap.GetUnequippedItemBySlot(slot)
	}

	swap.character.Equipment[slot] = *newItem
	oldItemStats := swap.getItemStats(oldItem)
	newItemStats := swap.getItemStats(*newItem)
	newStats := newItemStats.Subtract(oldItemStats)

	//2H will swap out the offhand also.
	if has2H && slot == proto.ItemSlot_ItemSlotMainHand {
		_, ohStats := swap.swapItem(proto.ItemSlot_ItemSlotOffHand, has2H, isReset)
		newStats = newStats.Add(ohStats)
	}

	swap.unEquippedItems[slot] = oldItem
	swap.swapWeapon(slot)

	return true, newStats
}

func (swap *ItemSwap) getItemStats(item Item) stats.Stats {
	return ItemEquipmentStats(item)
}

func (swap *ItemSwap) swapWeapon(slot proto.ItemSlot) {
	character := swap.character

	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:
		if character.AutoAttacks.AutoSwingMelee {
			character.AutoAttacks.SetMH(character.WeaponFromMainHand(swap.mhCritMultiplier))
		}
	case proto.ItemSlot_ItemSlotOffHand:
		if character.AutoAttacks.AutoSwingMelee {
			weapon := character.WeaponFromOffHand(swap.ohCritMultiplier)
			character.AutoAttacks.SetOH(weapon)

			character.AutoAttacks.IsDualWielding = weapon.SwingSpeed != 0
			character.PseudoStats.CanBlock = character.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
		}
	case proto.ItemSlot_ItemSlotRanged:
		if character.AutoAttacks.AutoSwingRanged {
			character.AutoAttacks.SetRanged(character.WeaponFromRanged(swap.rangedCritMultiplier))
		}
	}
}

func (swap *ItemSwap) reset(sim *Simulation) {
	if !swap.IsEnabled() {
		return
	}

	swap.SwapItems(sim, proto.APLActionItemSwap_Main, swap.slots, true)

	if !swap.initialized || swap.IsSwapped() {
		for _, slot := range swap.slots {
			for _, onSwap := range swap.onSwapCallbacks[slot] {
				onSwap(sim, slot)
			}
		}
	}

	swap.unEquippedItems = swap.swapEquip

	// This is used to set the initial spell flags for unequipped items.
	// Reset is called before the first iteration.
	swap.initialized = true
}

func toItem(itemSpec *proto.ItemSpec) Item {
	if itemSpec == nil || itemSpec.Id == 0 {
		return Item{}
	}

	return NewItem(ItemSpec{
		ID:           itemSpec.Id,
		Gems:         itemSpec.Gems,
		Enchant:      itemSpec.Enchant,
		RandomSuffix: itemSpec.RandomSuffix,
		Reforging:    itemSpec.Reforging,
	})
}
