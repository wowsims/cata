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

	isFuryWarrior        bool
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
	if hasItemSwap[proto.ItemSlot_ItemSlotOffHand] && hasMh {
		hasItemSwap[proto.ItemSlot_ItemSlotMainHand] = true
	}

	if has2H && hasOh {
		hasItemSwap[proto.ItemSlot_ItemSlotOffHand] = true
	}

	slots := SetToSortedSlice(hasItemSwap)

	if len(slots) == 0 {
		return
	}

	character.ItemSwap = ItemSwap{
		isFuryWarrior:        character.Spec == proto.Spec_SpecFuryWarrior,
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

	if (character.Env != nil) && character.Env.IsFinalized() {
		panic("Tried to add a new item swap callback in a finalized environment!")
	}

	for _, slot := range slots {
		character.ItemSwap.onSwapCallbacks[slot] = append(character.ItemSwap.onSwapCallbacks[slot], callback)
	}
}

// Helper for handling Item Effects that use the itemID to toggle the aura on and off
func (swap *ItemSwap) RegisterProc(itemID int32, aura *Aura, slots []proto.ItemSlot) {
	swap.registerProcInternal(ItemSwapProcConfig{
		ItemID: itemID,
		Aura:   aura,
		Slots:  slots,
	})
}

// Helper for handling Enchant Effects that use the effectID to toggle the aura on and off
func (swap *ItemSwap) RegisterEnchantProc(effectID int32, aura *Aura, slots []proto.ItemSlot) {
	swap.registerProcInternal(ItemSwapProcConfig{
		EnchantId: effectID,
		Aura:      aura,
		Slots:     slots,
	})
}

type ItemSwapProcConfig struct {
	ItemID    int32
	EnchantId int32
	Aura      *Aura
	Slots     []proto.ItemSlot
}

func (swap *ItemSwap) registerProcInternal(config ItemSwapProcConfig) {
	isItemProc := config.ItemID != 0
	isEnchantEffectProc := config.EnchantId != 0

	swap.character.RegisterItemSwapCallback(config.Slots, func(sim *Simulation, slot proto.ItemSlot) {
		var isItemSlotMatch = false

		if isItemProc {
			isItemSlotMatch = swap.HasEquippedItem(config.ItemID, config.Slots)
		} else if isEnchantEffectProc {
			isItemSlotMatch = swap.HasEquippedEnchant(config.EnchantId, config.Slots)
		}

		if isItemSlotMatch {
			if !config.Aura.IsActive() {
				config.Aura.Activate(sim)
				// Enchant effects such as Weapon/Back do not trigger an ICD
				if isItemProc && config.Aura.Icd != nil {
					config.Aura.Icd.Use(sim)
				}
			}
		} else {
			config.Aura.Deactivate(sim)
		}
	})
}

// Helper for handling Item On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) RegisterActive(itemID int32, slots []proto.ItemSlot) {
	swap.character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		isSwapItem := swap.ItemExistsInSwapSet(itemID)
		if !isSwapItem {
			return
		}
		hasItemEquipped := swap.character.HasItemEquipped(itemID)

		spell := swap.character.GetSpell(ActionID{ItemID: itemID})
		if spell != nil {
			aura := swap.character.GetAuraByID(spell.ActionID)
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

			spell.CD.Set(sim.CurrentTime + max(spell.CD.TimeToReady(sim), time.Second*30))
		}
	})
}

// Helper for handling Enchant On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) ProcessTinker(spell *Spell, slots []proto.ItemSlot) {
	swap.character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		if spell == nil || !swap.initialized {
			return
		}
		spell.CD.Set(sim.CurrentTime + max(spell.CD.TimeToReady(sim), time.Second*30))
	})
}

func (swap *ItemSwap) IsEnabled() bool {
	return swap.character != nil && len(swap.slots) > 0
}

func (swap *ItemSwap) IsSwapped() bool {
	return swap.swapSet == proto.APLActionItemSwap_Swap1
}

func (character *Character) HasItemEquipped(itemID int32) bool {
	for _, item := range character.Equipment {
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

func (swap *ItemSwap) HasEquippedItem(itemID int32, possibleSlots []proto.ItemSlot) bool {
	for _, slot := range possibleSlots {
		if swap.character.Equipment[slot].ID == itemID {
			return true
		}
	}
	return false
}

func (swap *ItemSwap) HasEquippedEnchant(effectID int32, possibleSlots []proto.ItemSlot) bool {
	for _, slot := range possibleSlots {
		if swap.character.Equipment[slot].Enchant.EffectID == effectID {
			return true
		}
	}
	return false
}

func (swap *ItemSwap) ItemExistsInSwapSet(itemID int32) bool {
	for _, item := range swap.swapEquip {
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

func (swap *ItemSwap) SwapItems(sim *Simulation, swapSet proto.APLActionItemSwap_SwapSet, isReset bool) {
	if !swap.IsEnabled() || swap.swapSet == swapSet && !isReset {
		return
	}

	character := swap.character

	weaponSlotSwapped := false
	newStats := stats.Stats{}
	has2H := swap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotMainHand).HandType == proto.HandType_HandTypeTwoHand
	isPrepull := sim.CurrentTime < 0

	for _, slot := range swap.slots {
		if !isReset && !isPrepull && (slot < proto.ItemSlot_ItemSlotMainHand || slot > proto.ItemSlot_ItemSlotRanged) {
			continue
		}

		//will swap both on the MainHand Slot for 2H.
		if has2H && slot == proto.ItemSlot_ItemSlotOffHand && !swap.isFuryWarrior {
			continue
		}

		if ok, swapStats := swap.swapItem(slot, has2H, isReset); ok {
			newStats = newStats.Add(swapStats)
			weaponSlotSwapped = slot == proto.ItemSlot_ItemSlotMainHand || slot == proto.ItemSlot_ItemSlotOffHand || slot == proto.ItemSlot_ItemSlotRanged || weaponSlotSwapped

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
		if weaponSlotSwapped {
			if character.AutoAttacks.AutoSwingMelee {
				character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
			}
			if character.AutoAttacks.AutoSwingRanged {
				character.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime)
			}
		}

		character.ExtendGCDUntil(sim, max(character.NextGCDAt(), sim.CurrentTime+GCDDefault))
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
	if has2H && slot == proto.ItemSlot_ItemSlotMainHand && !swap.isFuryWarrior {
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
	swap.initialized = false
	if !swap.IsEnabled() {
		return
	}

	swap.SwapItems(sim, proto.APLActionItemSwap_Main, true)

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
