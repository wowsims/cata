import { ItemSlot } from '../../../proto/common';
import { getEnumValues } from '../../../utils';

// Combines Fingers 1 and 2 and Trinket 1 and 2 into single groups
export enum BulkSimItemSlot {
	ItemSlotHead,
	ItemSlotNeck,
	ItemSlotShoulder,
	ItemSlotBack,
	ItemSlotChest,
	ItemSlotWrist,
	ItemSlotHands,
	ItemSlotWaist,
	ItemSlotLegs,
	ItemSlotFeet,
	ItemSlotFinger,
	ItemSlotTrinket,
	ItemSlotMainHand,
	ItemSlotOffHand,
	ItemSlotHandWeapon, // Weapon grouping slot for specs that can dual-wield
}

// Return all eligible bulk item slots.
// If the player can dual-wield, exclude main-hand/off-hand in favor of the grouped weapons slot
// Otherwise include main-hand/off-hand instead of the grouped weapons slot
export const getBulkItemSlots = (canDualWield: boolean) => {
	const allSlots = getEnumValues<BulkSimItemSlot>(BulkSimItemSlot);
	if (canDualWield) {
		return allSlots.filter(bulkSlot => ![BulkSimItemSlot.ItemSlotMainHand, BulkSimItemSlot.ItemSlotOffHand].includes(bulkSlot));
	} else {
		return allSlots.filter(bulkSlot => bulkSlot !== BulkSimItemSlot.ItemSlotHandWeapon);
	}
};

export const bulkSimSlotNames: Map<BulkSimItemSlot, string> = new Map([
	[BulkSimItemSlot.ItemSlotHead, 'Head'],
	[BulkSimItemSlot.ItemSlotNeck, 'Neck'],
	[BulkSimItemSlot.ItemSlotShoulder, 'Shoulders'],
	[BulkSimItemSlot.ItemSlotBack, 'Back'],
	[BulkSimItemSlot.ItemSlotChest, 'Chest'],
	[BulkSimItemSlot.ItemSlotWrist, 'Wrist'],
	[BulkSimItemSlot.ItemSlotHands, 'Hands'],
	[BulkSimItemSlot.ItemSlotWaist, 'Waist'],
	[BulkSimItemSlot.ItemSlotLegs, 'Legs'],
	[BulkSimItemSlot.ItemSlotFeet, 'Feet'],
	[BulkSimItemSlot.ItemSlotFinger, 'Rings'],
	[BulkSimItemSlot.ItemSlotTrinket, 'Trinkets'],
	[BulkSimItemSlot.ItemSlotMainHand, 'Main Hand'],
	[BulkSimItemSlot.ItemSlotOffHand, 'Off Hand'],
	[BulkSimItemSlot.ItemSlotHandWeapon, 'Weapons'],
]);

export const itemSlotToBulkSimItemSlot: Map<ItemSlot, BulkSimItemSlot> = new Map([
	[ItemSlot.ItemSlotHead, BulkSimItemSlot.ItemSlotHead],
	[ItemSlot.ItemSlotNeck, BulkSimItemSlot.ItemSlotNeck],
	[ItemSlot.ItemSlotShoulder, BulkSimItemSlot.ItemSlotShoulder],
	[ItemSlot.ItemSlotBack, BulkSimItemSlot.ItemSlotBack],
	[ItemSlot.ItemSlotChest, BulkSimItemSlot.ItemSlotChest],
	[ItemSlot.ItemSlotWrist, BulkSimItemSlot.ItemSlotWrist],
	[ItemSlot.ItemSlotHands, BulkSimItemSlot.ItemSlotHands],
	[ItemSlot.ItemSlotWaist, BulkSimItemSlot.ItemSlotWaist],
	[ItemSlot.ItemSlotLegs, BulkSimItemSlot.ItemSlotLegs],
	[ItemSlot.ItemSlotFeet, BulkSimItemSlot.ItemSlotFeet],
	[ItemSlot.ItemSlotFinger1, BulkSimItemSlot.ItemSlotFinger],
	[ItemSlot.ItemSlotFinger2, BulkSimItemSlot.ItemSlotFinger],
	[ItemSlot.ItemSlotTrinket1, BulkSimItemSlot.ItemSlotTrinket],
	[ItemSlot.ItemSlotTrinket2, BulkSimItemSlot.ItemSlotTrinket],
	[ItemSlot.ItemSlotMainHand, BulkSimItemSlot.ItemSlotMainHand],
	[ItemSlot.ItemSlotOffHand, BulkSimItemSlot.ItemSlotOffHand],
]);

export const getBulkItemSlotFromSlot = (slot: ItemSlot, canDualWield: boolean): BulkSimItemSlot => {
	if (canDualWield && [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand].includes(slot)) {
		return BulkSimItemSlot.ItemSlotHandWeapon;
	}
	return itemSlotToBulkSimItemSlot.get(slot)!;
};
