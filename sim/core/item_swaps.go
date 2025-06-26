package core

import (
	"slices"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type OnItemSwap func(*Simulation, proto.ItemSlot)

type ItemSwap struct {
	character           *Character
	onItemSwapCallbacks [NumItemSlots][]OnItemSwap

	isFuryWarrior        bool
	isFeralDruid         bool
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

	for idx, itemSpec := range itemSwap.Items {
		itemSlot := proto.ItemSlot(idx)
		hasItemSwap[itemSlot] = itemSpec != nil && itemSpec.Id != 0
		swapItems[itemSlot] = toItem(itemSpec)
	}

	has2HSwap := swapItems[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand
	hasMhEquipped := character.HasMHWeapon()
	hasOhEquipped := character.HasOHWeapon()

	// Handle MH and OH together, because present MH + empty OH --> swap MH and unequip OH
	if hasItemSwap[proto.ItemSlot_ItemSlotOffHand] && hasMhEquipped {
		hasItemSwap[proto.ItemSlot_ItemSlotMainHand] = true
	}

	if has2HSwap && hasOhEquipped {
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

	equipmentStats := calcItemSwapStatsOffset(character.Equipment, swapItems, prepullBonusStats, slots, character.Spec)

	character.ItemSwap = ItemSwap{
		isFuryWarrior:        character.Spec == proto.Spec_SpecFuryWarrior,
		isFeralDruid:         character.Spec == proto.Spec_SpecFeralDruid || character.Spec == proto.Spec_SpecGuardianDruid,
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
		panic("Tried to add a new item swap callback for slots in a finalized environment!")
	}

	for _, slot := range slots {
		character.ItemSwap.onItemSwapCallbacks[slot] = append(character.ItemSwap.onItemSwapCallbacks[slot], callback)
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

	// Enchant effects such as Weapon/Back do not trigger an ICD
	shouldUpdateIcd := isItemProc && (config.Aura.Icd != nil)

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
			if swap.initialized && shouldUpdateIcd {
				config.Aura.Icd.Use(sim)
			}
		} else {
			config.Aura.Deactivate(sim)
			if swap.initialized && shouldUpdateIcd {
				// This is a hack to block ActivateAura APL
				// actions from executing for unequipped items.
				config.Aura.Icd.Set(NeverExpires)
			}
		}
	})
}

// Helper for handling Item On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) RegisterActive(itemID int32) {
	slots := swap.EligibleSlotsForItem(itemID)
	itemActionID := ActionID{ItemID: itemID}
	character := swap.character

	character.RegisterItemSwapCallback(slots, func(sim *Simulation, _ proto.ItemSlot) {
		spell := character.GetSpell(itemActionID)
		if spell == nil {
			return
		}

		aura := character.GetAuraByID(spell.ActionID)
		if aura.IsActive() {
			aura.Deactivate(sim)
		}

		hasItemEquipped := character.hasItemEquipped(itemID, slots)
		if !hasItemEquipped {
			spell.Flags |= SpellFlagSwapped
			return
		}

		spell.Flags &= ^SpellFlagSwapped

		if !swap.initialized {
			return
		}

		spell.CD.Set(sim.CurrentTime + max(spell.CD.TimeToReady(sim), time.Second*30))
	})
}

// Helper for handling Enchant On Use effects to set a 30s cd on the related spell.
func (swap *ItemSwap) ProcessTinker(spell *Spell, slots []proto.ItemSlot) {
	swap.character.RegisterItemSwapCallback(slots, func(sim *Simulation, slot proto.ItemSlot) {
		if spell == nil || !swap.initialized {
			return
		}

		isUniqueItem := swap.GetEquippedItemBySlot(slot).ID != swap.GetUnequippedItemBySlot(slot).ID

		var newSpellCD time.Duration
		if isUniqueItem {
			// Unique items have a 30s CD regardless of the spell CD being > 30s or not
			newSpellCD = time.Second * 30
		} else {
			// Items with the same ItemID share the CD and does not get reset to 30s
			timeToReady := spell.CD.TimeToReady(sim)
			newSpellCD = TernaryDuration(timeToReady > 30, timeToReady, time.Second*30)
		}

		spell.CD.Set(sim.CurrentTime + newSpellCD)
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

func (swap *ItemSwap) CouldHaveItemEquippedInSlot(itemID int32, slot proto.ItemSlot) bool {
	return swap.character.Equipment.containsItemInSlots(itemID, []proto.ItemSlot{slot}) || swap.unEquippedItems.containsItemInSlots(itemID, []proto.ItemSlot{slot})
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

func (swap *ItemSwap) EligibleSlotsForItem(itemID int32) []proto.ItemSlot {
	eligibleSlots := eligibleSlotsForItem(GetItemByID(itemID), swap.isFuryWarrior)

	if len(eligibleSlots) == 0 {
		return []proto.ItemSlot{}
	}

	if !swap.IsEnabled() {
		return FilterSlice(eligibleSlots, func(slot proto.ItemSlot) bool {
			return (swap.character.Equipment[slot].ID == itemID)
		})
	}

	return FilterSlice(eligibleSlots, func(slot proto.ItemSlot) bool {
		return (swap.originalEquip[slot].ID == itemID) || (swap.swapEquip[slot].ID == itemID)
	})
}

func (swap *ItemSwap) EligibleSlotsForEffect(effectID int32) []proto.ItemSlot {
	var eligibleSlots []proto.ItemSlot

	for itemSlot := proto.ItemSlot(0); itemSlot < NumItemSlots; itemSlot++ {
		if !swap.IsEnabled() {
			if swap.character.Equipment.containsEnchantInSlot(effectID, itemSlot) {
				eligibleSlots = append(eligibleSlots, itemSlot)
			}
		} else {
			if swap.originalEquip.containsEnchantInSlot(effectID, itemSlot) || swap.swapEquip.containsEnchantInSlot(effectID, itemSlot) {
				eligibleSlots = append(eligibleSlots, itemSlot)
			}
		}
	}

	return eligibleSlots
}

func (swap *ItemSwap) SwapItems(sim *Simulation, swapSet proto.APLActionItemSwap_SwapSet, isReset bool) {
	if !swap.IsEnabled() || (!swap.IsValidSwap(swapSet) && !isReset) {
		return
	}

	character := swap.character
	weaponSlotSwapped := false
	isPrepull := sim.CurrentTime < 0

	for _, slot := range swap.slots {
		if slot == proto.ItemSlot_ItemSlotMainHand || slot == proto.ItemSlot_ItemSlotOffHand {
			weaponSlotSwapped = true
		} else if !isReset && !isPrepull {
			continue
		}

		swap.swapItem(sim, slot, isPrepull, isReset)

		for _, onSwapSlot := range swap.onItemSwapCallbacks[slot] {
			onSwapSlot(sim, slot)
		}
	}

	if !swap.IsValidSwap(swapSet) {
		return
	}

	statsToSwap := Ternary(isPrepull, swap.equipmentStats.allSlots, swap.equipmentStats.weaponSlots)
	if swap.IsSwapped() {
		statsToSwap = statsToSwap.Invert()
	}

	if sim.Log != nil {
		sim.Log("Item Swap - Stats Change: %v", statsToSwap.FlatString())
	}
	character.AddDynamicEquipStats(sim, statsToSwap)

	if !isPrepull && !isReset && weaponSlotSwapped {
		character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime)
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

	if isPrepull {
		return
	}

	character := swap.character

	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:

		// As of MoP Ranged Weapons are worn in the Main Hand
		// If we can get it, we equipeed a valid ranged weapon
		if character.Ranged() != nil {
			if character.AutoAttacks.AutoSwingRanged {
				character.AutoAttacks.SetRanged(character.WeaponFromRanged(swap.rangedCritMultiplier))
			}
		} else {
			// Feral's concept of Paws is handeled in the druid.go Initialize()
			// and doesn't need MH swap handling here.
			if character.AutoAttacks.AutoSwingMelee && !swap.isFeralDruid {
				character.AutoAttacks.SetMH(character.WeaponFromMainHand(swap.mhCritMultiplier))
			}
		}
	case proto.ItemSlot_ItemSlotOffHand:
		// OH slot handling is more involved because we need to dynamically toggle the OH weapon attack on/off
		// depending on the updated DW status after the swap.
		if character.AutoAttacks.AutoSwingMelee {
			weapon := character.WeaponFromOffHand(swap.ohCritMultiplier)
			isCurrentlyDualWielding := character.AutoAttacks.IsDualWielding
			character.AutoAttacks.SetOH(weapon)
			if !isPrepull && !isCurrentlyDualWielding {
				character.AutoAttacks.IsDualWielding = weapon.SwingSpeed != 0
				character.AutoAttacks.EnableMeleeSwing(sim)
			}
			character.PseudoStats.CanBlock = character.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
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

func calcItemSwapStatsOffset(originalEquipment Equipment, swapEquipment Equipment, prepullBonusStats stats.Stats, slots []proto.ItemSlot, spec proto.Spec) ItemSwapStats {
	allSlotStats := prepullBonusStats
	allWeaponSlots := AllWeaponSlots()
	swapStatEquipment := originalEquipment
	weaponStatEquipment := originalEquipment

	for _, slot := range slots {
		if swapEquipment.GetItemBySlot(slot) != nil {
			swapStatEquipment[slot] = swapEquipment[slot]
		}

		if slices.Contains(allWeaponSlots, slot) {
			weaponStatEquipment[slot] = swapEquipment[slot]
		}
	}

	allSlotStats = allSlotStats.Add(swapStatEquipment.Stats(spec).Subtract(originalEquipment.Stats(spec)))
	weaponSlotStats := weaponStatEquipment.Stats(spec).Subtract(originalEquipment.Stats(spec))
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
		ID:            itemSpec.Id,
		Gems:          itemSpec.Gems,
		Enchant:       itemSpec.Enchant,
		Tinker:        itemSpec.Tinker,
		RandomSuffix:  itemSpec.RandomSuffix,
		Reforging:     itemSpec.Reforging,
		UpgradeStep:   itemSpec.UpgradeStep,
		ChallengeMode: itemSpec.ChallengeMode,
	})
}
