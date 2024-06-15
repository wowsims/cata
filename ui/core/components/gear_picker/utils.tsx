import { ref } from 'tsx-vanilla';

import { GemColor, ItemSlot } from '../../proto/common';
import { UIGem as Gem } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { getEmptyGemSocketIconUrl } from '../../proto_utils/gems';

const emptySlotIcons: Record<ItemSlot, string> = {
	[ItemSlot.ItemSlotHead]: '/cata/assets/item_slots/head.jpg',
	[ItemSlot.ItemSlotNeck]: '/cata/assets/item_slots/neck.jpg',
	[ItemSlot.ItemSlotShoulder]: '/cata/assets/item_slots/shoulders.jpg',
	[ItemSlot.ItemSlotBack]: '/cata/assets/item_slots/shirt.jpg',
	[ItemSlot.ItemSlotChest]: '/cata/assets/item_slots/chest.jpg',
	[ItemSlot.ItemSlotWrist]: '/cata/assets/item_slots/wrists.jpg',
	[ItemSlot.ItemSlotHands]: '/cata/assets/item_slots/hands.jpg',
	[ItemSlot.ItemSlotWaist]: '/cata/assets/item_slots/waist.jpg',
	[ItemSlot.ItemSlotLegs]: '/cata/assets/item_slots/legs.jpg',
	[ItemSlot.ItemSlotFeet]: '/cata/assets/item_slots/feet.jpg',
	[ItemSlot.ItemSlotFinger1]: '/cata/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotFinger2]: '/cata/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotTrinket1]: '/cata/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotTrinket2]: '/cata/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotMainHand]: '/cata/assets/item_slots/mainhand.jpg',
	[ItemSlot.ItemSlotOffHand]: '/cata/assets/item_slots/offhand.jpg',
	[ItemSlot.ItemSlotRanged]: '/cata/assets/item_slots/ranged.jpg',
};
export function getEmptySlotIconUrl(slot: ItemSlot): string {
	return emptySlotIcons[slot];
}

export const createHeroicLabel = () => {
	return <span className="heroic-label">[H]</span>;
};

export const createGemContainer = (socketColor: GemColor, gem: Gem | null, index: number) => {
	const gemIconElem = ref<HTMLImageElement>();
	const gemContainerElem = ref<HTMLAnchorElement>();
	const gemContainer = (
		<a ref={gemContainerElem} className="gem-socket-container" href="javascript:void(0)" dataset={{ socketIdx: index }}>
			<img ref={gemIconElem} className={`gem-icon ${!gem ? 'hide' : ''}`} />
			<img className="socket-icon" src={getEmptyGemSocketIconUrl(socketColor)} />
		</a>
	);

	if (!!gem) {
		ActionId.fromItemId(gem.id)
			.fill()
			.then(filledId => {
				filledId.setWowheadHref(gemContainerElem.value!);
				gemIconElem.value!.src = filledId.iconUrl;
			});
	}
	return gemContainer as HTMLAnchorElement;
};
