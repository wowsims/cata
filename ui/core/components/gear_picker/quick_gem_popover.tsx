import { Player } from '../../player';
import { GemColor, ItemSlot } from '../../proto/common';
import { UIGem as Gem } from '../../proto/ui.js';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { TypedEvent } from '../../typed_event';
import QuickSwapList from '../quick_swap';

export const addQuickGemPopover = (
	player: Player<any>,
	tooltipElement: HTMLElement,
	item: EquippedItem,
	itemSlot: ItemSlot,
	socketIdx: number,
	openDetailTab: () => void,
) => {
	return new QuickSwapList({
		title: 'Favorite gems',
		emptyMessage: 'Add favorite gems.',
		tippyElement: tooltipElement,
		tippyConfig: {
			appendTo: document.querySelector('.sim-ui')!,
		},
		item,
		getItems: (currentItem: EquippedItem) => {
			const favoriteGems = player.sim
				.getFilters()
				.favoriteGems?.map(id => player.sim.db.lookupGem(id))
				.filter((gem): gem is Gem => {
					if (!gem) return false;
					return !(itemSlot !== ItemSlot.ItemSlotHead && gem?.color === GemColor.GemColorMeta);
				})
				.sort((a, b) => (a.color > b.color ? 1 : -1));

			return favoriteGems.map(gem => ({
				item: gem,
				active: currentItem.gems[socketIdx]?.id === gem.id,
			}));
		},
		onItemClick: clickedItem => {
			player.equipItem(TypedEvent.nextEventID(), itemSlot, item.withGem(clickedItem, socketIdx));
		},
		footerButton: {
			label: 'Open Gems',
			onClick: openDetailTab,
		},
	});
};
