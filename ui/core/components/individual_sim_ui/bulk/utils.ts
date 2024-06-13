import { ItemSlot } from '../../../proto/common';

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
	ItemSlotRanged,
}

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
	[BulkSimItemSlot.ItemSlotRanged, 'Ranged'],
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
	[ItemSlot.ItemSlotRanged, BulkSimItemSlot.ItemSlotRanged],
]);
