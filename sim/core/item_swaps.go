package core

import (
	"slices"
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
	equipmentStats  ItemSwapStats

	initialized bool
}

type ItemSwapStats struct {
	allSlots    stats.Stats
	weaponSlots stats.Stats
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

	var prepullBonusStats stats.Stats
	if itemSwap.PrepullBonusStats != nil {
		prepullBonusStats = stats.FromUnitStatsProto(itemSwap.PrepullBonusStats)
	}

	equipmentStats := calcItemSwapStatsOffset(character.Equipment, swapItems, prepullBonusStats, slots)

	character.ItemSwap = ItemSwap{
		isFuryWarrior:        character.Spec == proto.Spec_SpecFuryWarrior,
		mhCritMultiplier:     mhCritMultiplier,
		ohCritMultiplier:     ohCritMultiplier,
		rangedCritMultiplier: rangedCritMultiplier,
		slots:                slots,
		originalEquip:        character.Equipment,
		swapEquip:            swapItems,
		unEquippedItems:      swapItems,
		equipmentStats:       equipmentStats,
		swapSet:              proto.APLActionItemSwap_Main,
		initialized:          false,
	}
}

func (swap *ItemSwap) initialize(character *Character) {
	swap.character = character
}

func (character *Character) RegisterItemSwapCallback(slots []proto.ItemSlot, callback OnItemSwap) {
	if character == nil || !character.ItemSwap.IsEnabled() || len(slots) == 0 {
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
// This will also get the eligible slots for the item
func (swap *ItemSwap) RegisterProc(itemID int32, aura *Aura) {
	slots := swap.EligibleSlotsForItem(itemID)
	swap.RegisterProcWithSlots(itemID, aura, slots)
}

// Helper for handling Item Effects that use the itemID to toggle the aura on and off
func (swap *ItemSwap) RegisterProcWithSlots(itemID int32, aura *Aura, slots []proto.ItemSlot) {
	swap.registerProcInternal(ItemSwapProcConfig{
		ItemID: itemID,
		Aura:   aura,
		Slots:  slots,
	})
}

// Helper for handling Enchant Effects that use the effectID to toggle the aura on and off
func (swap *ItemSwap) RegisterEnchantProc(effectID int32, aura *Aura) {
	slots := swap.EligibleSlotsForEffect(effectID)
	swap.RegisterEnchantProcWithSlots(effectID, aura, slots)
}
func (swap *ItemSwap) RegisterEnchantProcWithSlots(effectID int32, aura *Aura, slots []proto.ItemSlot) {
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
	character := swap.character

	character.RegisterItemSwapCallback(config.Slots, func(sim *Simulation, _ proto.ItemSlot) {
		isItemSlotMatch := false
		if isItemProc {
			isItemSlotMatch = character.hasItemEquipped(config.ItemID, config.Slots)
		} else if isEnchantEffectProc {
			isItemSlotMatch = character.hasEnchantEquipped(config.EnchantId, config.Slots)
		}

		if isItemSlotMatch {
			if !config.Aura.IsActive() {
				config.Aura.Activate(sim)
			}
		} else {
			config.Aura.Deactivate(sim)
		}
		// Enchant effects such as Weapon/Back do not trigger an ICD
		if isItemProc && config.Aura.Icd != nil {
			config.Aura.Icd.Use(sim)
		}
	})
}

// Helper for handling Item On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) RegisterActive(itemID int32) {
	slots := swap.EligibleSlotsForItem(itemID)
	if !swap.CanRegisterItemCallback(itemID, slots) {
		return
	}

	swap.character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		hasItemEquipped := swap.character.hasItemEquipped(itemID, slots)

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

func (swap *ItemSwap) IsValidSwap(swapSet proto.APLActionItemSwap_SwapSet) bool {
	return swap.swapSet != swapSet
}

func (swap *ItemSwap) IsSwapped() bool {
	return swap.swapSet == proto.APLActionItemSwap_Swap1
}

func (character *Character) hasItemEquipped(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return character.Equipment.containsItemInSlots(itemID, possibleSlots)
}

func (character *Character) hasEnchantEquipped(effectID int32, possibleSlots []proto.ItemSlot) bool {
	return character.Equipment.containsEnchantInSlots(effectID, possibleSlots)
}

func (swap *ItemSwap) GetEquippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.character.Equipment[slot]
}

func (swap *ItemSwap) GetUnequippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.unEquippedItems[slot]
}

func (swap *ItemSwap) CanRegisterItemCallback(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return swap.ItemExistsInMainEquip(itemID, possibleSlots) || swap.ItemExistsInSwapEquip(itemID, possibleSlots)
}

func (swap *ItemSwap) ItemExistsInMainEquip(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return swap.originalEquip.containsItemInSlots(itemID, possibleSlots)
}

func (swap *ItemSwap) ItemExistsInSwapEquip(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return swap.swapEquip.containsItemInSlots(itemID, possibleSlots)
}

func (swap *ItemSwap) EligibleSlotsForItem(itemID int32) []proto.ItemSlot {
	eligibleSlots := eligibleSlotsForItem(GetItemByID(itemID), swap.isFuryWarrior)

	if len(eligibleSlots) == 0 {
		return []proto.ItemSlot{}
	}

	if !swap.IsEnabled() {
		return eligibleSlots
	} else {
		return FilterSlice(eligibleSlots, func(slot proto.ItemSlot) bool {
			return (swap.originalEquip[slot].ID == itemID) || (swap.swapEquip[slot].ID == itemID)
		})
	}
}

func (swap *ItemSwap) EligibleSlotsForEffect(effectID int32) []proto.ItemSlot {
	var eligibleSlots []proto.ItemSlot

	if swap.IsEnabled() {
		for itemSlot := proto.ItemSlot(0); itemSlot < NumItemSlots; itemSlot++ {
			if swap.originalEquip.containsEnchantInSlot(effectID, itemSlot) || swap.swapEquip.containsEnchantInSlot(effectID, itemSlot) {
				eligibleSlots = append(eligibleSlots, itemSlot)
			}
		}
	}

	return eligibleSlots
}

func (swap *ItemSwap) SwapItems(sim *Simulation, swapSet proto.APLActionItemSwap_SwapSet, isReset bool) {
	if !swap.IsEnabled() || !swap.IsValidSwap(swapSet) && !isReset {
		return
	}

	character := swap.character

	weaponSlotSwapped := false
	isPrepull := sim.CurrentTime < 0

	for _, slot := range swap.slots {
		if !isReset && !isPrepull && (slot < proto.ItemSlot_ItemSlotMainHand || slot > proto.ItemSlot_ItemSlotRanged) {
			continue
		}

		swap.swapItem(sim, slot, isPrepull, isReset)
		weaponSlotSwapped = slot == proto.ItemSlot_ItemSlotMainHand || slot == proto.ItemSlot_ItemSlotOffHand || slot == proto.ItemSlot_ItemSlotRanged || weaponSlotSwapped

		for _, onSwap := range swap.onSwapCallbacks[slot] {
			onSwap(sim, slot)
		}
	}

	if !swap.IsValidSwap(swapSet) && isReset {
		return
	}

	statsToSwap := Ternary(isPrepull, swap.equipmentStats.allSlots, swap.equipmentStats.weaponSlots)
	if swap.IsSwapped() {
		statsToSwap = statsToSwap.Invert()
	}

	if sim.Log != nil {
		sim.Log("Item Swap - Stats Change: %v", statsToSwap.FlatString())
	}
	character.AddStatsDynamic(sim, statsToSwap)

	if !isPrepull && !isReset && weaponSlotSwapped {
		character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		character.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime)
		character.ExtendGCDUntil(sim, max(character.NextGCDAt(), sim.CurrentTime+GCDDefault))
	}

	swap.swapSet = swapSet
}

func (swap *ItemSwap) swapItem(sim *Simulation, slot proto.ItemSlot, isPrepull bool, isReset bool) {
	oldItem := *swap.GetEquippedItemBySlot(slot)

	if isReset {
		swap.character.Equipment[slot] = swap.originalEquip[slot]
	} else {
		swap.character.Equipment[slot] = swap.unEquippedItems[slot]
	}

	swap.unEquippedItems[slot] = oldItem

	if !isPrepull {
		switch slot {
		case proto.ItemSlot_ItemSlotMainHand:
			if swap.character.AutoAttacks.AutoSwingMelee {
				swap.character.AutoAttacks.SetMH(swap.character.WeaponFromMainHand(swap.mhCritMultiplier))
			}
		case proto.ItemSlot_ItemSlotOffHand:
			if swap.character.AutoAttacks.AutoSwingMelee {
				weapon := swap.character.WeaponFromOffHand(swap.ohCritMultiplier)
				swap.character.AutoAttacks.SetOH(weapon)
				swap.character.AutoAttacks.IsDualWielding = weapon.SwingSpeed != 0
				swap.character.AutoAttacks.EnableMeleeSwing(sim)
				swap.character.PseudoStats.CanBlock = swap.character.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
			}
		case proto.ItemSlot_ItemSlotRanged:
			if swap.character.AutoAttacks.AutoSwingRanged {
				swap.character.AutoAttacks.SetRanged(swap.character.WeaponFromRanged(swap.rangedCritMultiplier))
			}
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

func (swap *ItemSwap) doneIteration(sim *Simulation) {
	swap.reset(sim)
}

func calcItemSwapStatsOffset(originalEquipment Equipment, swapEquipment Equipment, prepullBonusStats stats.Stats, slots []proto.ItemSlot) ItemSwapStats {
	allSlotStats := stats.Stats{}
	weaponSlotStats := stats.Stats{}

	allWeaponSlots := AllWeaponSlots()
	allSlotStats = allSlotStats.Add(prepullBonusStats)

	for _, slot := range slots {
		slotStats := ItemEquipmentStats(swapEquipment[slot]).Subtract(ItemEquipmentStats(originalEquipment[slot]))
		allSlotStats = allSlotStats.Add(slotStats)

		if slices.Contains(allWeaponSlots, slot) {
			weaponSlotStats = weaponSlotStats.Add(slotStats)
		}
	}

	return ItemSwapStats{
		allSlots:    allSlotStats,
		weaponSlots: weaponSlotStats,
	}
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
