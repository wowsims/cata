import { ref } from 'tsx-vanilla';

import { Player } from '../../player';
import { ItemSlot, ItemType } from '../../proto/common';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { SimUI } from '../../sim_ui';
import { EventID } from '../../typed_event';
import { Component } from '../component';
import { GearData } from './item_list';
import SelectorModal, { SelectorModalTabs } from './selector_modal';
import { createGemContainer, getEmptySlotIconUrl } from './utils';

export default class IconItemSwapPicker extends Component {
	private readonly iconAnchor: HTMLAnchorElement;
	private readonly socketsContainerElem: HTMLElement;
	private readonly player: Player<any>;
	private readonly slot: ItemSlot;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'icon-picker-root');
		this.rootElem.classList.add('icon-picker');
		this.player = player;
		this.slot = slot;

		const iconAnchorRef = ref<HTMLAnchorElement>();
		const socketsContainerRef = ref<HTMLDivElement>();

		this.rootElem.prepend(
			<a ref={iconAnchorRef} className="icon-picker-button" href="#" attributes={{ role: 'button' }}>
				<div ref={socketsContainerRef} className="item-picker-sockets-container" />
			</a>,
		);

		this.iconAnchor = iconAnchorRef.value!;
		this.socketsContainerElem = socketsContainerRef.value!;

		const selectorModal = new SelectorModal(simUI.rootElem, simUI, this.player);

		player.sim.waitForInit().then(() => {
			this.iconAnchor.addEventListener('click', (event: Event) => {
				event.preventDefault();
				selectorModal.openTab(this.slot, SelectorModalTabs.Items, this.createGearData());
			});
		});

		player.itemSwapChangeEmitter.on(() => {
			this.update(player.getItemSwapGear().getEquippedItem(slot));
		});
	}

	update(newItem: EquippedItem | null) {
		this.iconAnchor.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconAnchor.removeAttribute('data-wowhead');
		this.iconAnchor.href = '#';

		if (newItem) {
			newItem.asActionId().fillAndSet(this.iconAnchor, true, true);
			this.player.setWowheadData(newItem, this.iconAnchor);

			this.socketsContainerElem.replaceChildren(
				<>
					{newItem.allSocketColors().map((socketColor, gemIdx) => {
						const gemContainer = createGemContainer(socketColor, newItem.gems[gemIdx], gemIdx);
						if (gemIdx === newItem.numPossibleSockets - 1 && [ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(newItem.item.type)) {
							const updateProfession = () => {
								gemContainer.classList[this.player.isBlacksmithing() ? 'remove' : 'add']('hide');
							};
							this.player.professionChangeEmitter.on(updateProfession);
							updateProfession();
						}
						return gemContainer;
					})}
				</>,
			);

			this.iconAnchor.classList.add('active');
		} else {
			this.socketsContainerElem.replaceChildren();
			this.iconAnchor.classList.remove('active');
		}
	}

	private createGearData(): GearData {
		return {
			equipItem: (eventID: EventID, newItem: EquippedItem | null) => {
				this.player.equipItemSwapitem(eventID, this.slot, newItem);
			},
			getEquippedItem: () => this.player.getItemSwapItem(this.slot),
			changeEvent: this.player.itemSwapChangeEmitter,
		};
	}
}
