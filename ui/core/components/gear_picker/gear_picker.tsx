import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../../css_utils';
import { Player } from '../../player';
import { GemColor, ItemSlot, ItemType } from '../../proto/common';
import { UIEnchant as Enchant, UIGem as Gem } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { getEnchantDescription } from '../../proto_utils/enchants';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { getEmptyGemSocketIconUrl } from '../../proto_utils/gems';
import { shortSecondaryStatNames, slotNames } from '../../proto_utils/names';
import { SimUI } from '../../sim_ui';
import { EventID } from '../../typed_event';
import { Component } from '../component';
import QuickSwapList from '../quick_swap';
import { GearData } from './item_list';
import { addQuickEnchantPopover } from './quick_enchant_popover';
import { addQuickGemPopover } from './quick_gem_popover';
import SelectorModal, { SelectorModalTabs } from './selector_modal';

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
		<a ref={gemContainerElem} className="gem-socket-container" href="javascript:void(0)" attributes={{ role: 'button' }} dataset={{ socketIdx: index }}>
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

export class GearPicker extends Component {
	// ItemSlot is used as the index
	readonly itemPickers: Array<ItemPicker>;
	readonly selectorModal: SelectorModal;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'gear-picker-root');

		const leftSide = document.createElement('div');
		leftSide.classList.add('gear-picker-left', 'tab-panel-col');
		this.rootElem.appendChild(leftSide);

		const rightSide = document.createElement('div');
		rightSide.classList.add('gear-picker-right', 'tab-panel-col');
		this.rootElem.appendChild(rightSide);

		const leftItemPickers = [
			ItemSlot.ItemSlotHead,
			ItemSlot.ItemSlotNeck,
			ItemSlot.ItemSlotShoulder,
			ItemSlot.ItemSlotBack,
			ItemSlot.ItemSlotChest,
			ItemSlot.ItemSlotWrist,
			ItemSlot.ItemSlotMainHand,
			ItemSlot.ItemSlotOffHand,
			ItemSlot.ItemSlotRanged,
		].map(slot => new ItemPicker(leftSide, this, simUI, player, slot));

		const rightItemPickers = [
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
			ItemSlot.ItemSlotFinger1,
			ItemSlot.ItemSlotFinger2,
			ItemSlot.ItemSlotTrinket1,
			ItemSlot.ItemSlotTrinket2,
		].map(slot => new ItemPicker(rightSide, this, simUI, player, slot));

		this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);

		this.selectorModal = new SelectorModal(simUI.rootElem, simUI, player, this, { id: 'gear-picker-selector-modal' });
	}
}

export class ItemRenderer extends Component {
	private readonly player: Player<any>;

	readonly iconElem: HTMLAnchorElement;
	readonly nameElem: HTMLAnchorElement;
	readonly enchantElem: HTMLAnchorElement;
	readonly reforgeElem: HTMLAnchorElement;
	readonly socketsContainerElem: HTMLElement;
	socketsElem: HTMLAnchorElement[] = [];

	constructor(parent: HTMLElement, root: HTMLElement, player: Player<any>) {
		super(parent, 'item-picker-root', root);
		this.player = player;

		const iconElem = ref<HTMLAnchorElement>();
		const nameElem = ref<HTMLAnchorElement>();
		const enchantElem = ref<HTMLAnchorElement>();
		const reforgeElem = ref<HTMLAnchorElement>();
		const sce = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<>
				<div className="item-picker-icon-wrapper">
					<a ref={iconElem} className="item-picker-icon" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<div ref={sce} className="item-picker-sockets-container"></div>
				</div>
				<div className="item-picker-labels-container">
					<a ref={nameElem} className="item-picker-name" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<a ref={enchantElem} className="item-picker-enchant hide" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<a ref={reforgeElem} className="item-picker-reforge hide" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
				</div>
			</>,
		);

		this.iconElem = iconElem.value!;
		this.nameElem = nameElem.value!;
		this.reforgeElem = reforgeElem.value!;
		this.enchantElem = enchantElem.value!;
		this.socketsContainerElem = sce.value!;
	}

	clear() {
		this.nameElem.removeAttribute('data-wowhead');
		this.nameElem.removeAttribute('href');
		this.iconElem.removeAttribute('data-wowhead');
		this.iconElem.removeAttribute('href');
		this.enchantElem.removeAttribute('data-wowhead');
		this.enchantElem.removeAttribute('href');
		this.enchantElem.classList.add('hide');
		this.reforgeElem.classList.add('hide');

		this.iconElem.style.backgroundImage = '';
		this.enchantElem.innerText = '';
		this.reforgeElem.innerText = '';
		this.socketsContainerElem.innerText = '';
		this.socketsElem = [];
		this.nameElem.textContent = '';
	}

	update(newItem: EquippedItem) {
		this.nameElem.textContent = newItem.item.name;

		if (newItem.randomSuffix) {
			this.nameElem.textContent += ' ' + newItem.randomSuffix.name;
		}

		if (newItem.item.heroic) {
			this.nameElem.insertAdjacentElement('beforeend', createHeroicLabel());
		} else {
			this.nameElem.querySelector('.heroic-label')?.remove();
		}

		if (newItem.reforge) {
			const reforgeData = this.player.getReforgeData(newItem, newItem.reforge);
			const fromText = shortSecondaryStatNames.get(newItem.reforge?.fromStat[0]);
			const toText = shortSecondaryStatNames.get(newItem.reforge?.toStat[0]);
			this.reforgeElem.innerText = `Reforged ${Math.abs(reforgeData.fromAmount)} ${fromText} → ${reforgeData.toAmount} ${toText}`;
			this.reforgeElem.classList.remove('hide');
		} else {
			this.reforgeElem.innerText = '';
			this.reforgeElem.classList.add('hide');
		}

		setItemQualityCssClass(this.nameElem, newItem.item.quality);

		this.player.setWowheadData(newItem, this.iconElem);
		this.player.setWowheadData(newItem, this.nameElem);

		newItem
			.asActionId()
			.fill()
			.then(filledId => {
				filledId.setBackgroundAndHref(this.iconElem);
				filledId.setWowheadHref(this.nameElem);
			});

		if (newItem.enchant) {
			getEnchantDescription(newItem.enchant).then(description => {
				this.enchantElem.textContent = description;
			});
			// Make enchant text hover have a tooltip.
			if (newItem.enchant.spellId) {
				this.enchantElem.href = ActionId.makeSpellUrl(newItem.enchant.spellId);
				ActionId.makeSpellTooltipData(newItem.enchant.spellId).then(url => {
					this.enchantElem.dataset.wowhead = url;
				});
			} else {
				this.enchantElem.href = ActionId.makeItemUrl(newItem.enchant.itemId);
				ActionId.makeItemTooltipData(newItem.enchant.itemId).then(url => {
					this.enchantElem.dataset.wowhead = url;
				});
			}
			this.enchantElem.dataset.whtticon = 'false';
			this.enchantElem.classList.remove('hide');
		} else {
			this.enchantElem.classList.add('hide');
		}

		newItem.allSocketColors().forEach((socketColor, gemIdx) => {
			const gemContainer = createGemContainer(socketColor, newItem.gems[gemIdx], gemIdx);
			if (gemIdx === newItem.numPossibleSockets - 1 && [ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(newItem.item.type)) {
				const updateProfession = () => {
					if (this.player.isBlacksmithing()) {
						gemContainer.classList.remove('hide');
					} else {
						gemContainer.classList.add('hide');
					}
				};
				this.player.professionChangeEmitter.on(updateProfession);
				updateProfession();
			}
			this.socketsElem.push(gemContainer);
			this.socketsContainerElem.appendChild(gemContainer);
		});
	}
}

export class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly onUpdateCallbacks: (() => void)[] = [];

	private readonly itemElem: ItemRenderer;
	private readonly gearPicker: GearPicker;

	// All items and enchants that are eligible for this slot
	private _equippedItem: EquippedItem | null = null;

	private quickSwapEnchantPopover: QuickSwapList<Enchant> | null = null;
	private quickSwapGemPopover: QuickSwapList<Gem>[] = [];

	constructor(parent: HTMLElement, gearPicker: GearPicker, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');

		this.gearPicker = gearPicker;
		this.simUI = simUI;
		this.player = player;
		this.slot = slot;
		this.itemElem = new ItemRenderer(parent, this.rootElem, player);

		this.item = player.getEquippedItem(slot);

		player.sim.waitForInit().then(() => {
			const openGearSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Items);
			};
			const openReforgeSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Reforging);
			};

			this.itemElem.iconElem.addEventListener('click', openGearSelector);
			this.itemElem.nameElem.addEventListener('click', openGearSelector);
			this.itemElem.reforgeElem.addEventListener('click', openReforgeSelector);
			this.addQuickEnchantHelpers();
		});

		player.gearChangeEmitter.on(() => {
			this.item = this.player.getEquippedItem(this.slot);
			if (this._equippedItem) {
				if (this._equippedItem !== this.quickSwapEnchantPopover?.item) {
					this.quickSwapEnchantPopover?.update({ item: this._equippedItem });
				}
				this.addQuickGemHelpers();
			}
		});

		player.sim.filtersChangeEmitter.on(() => {
			if (this._equippedItem) {
				this.quickSwapEnchantPopover?.update({ item: this._equippedItem });
				this.quickSwapGemPopover.forEach(quickSwap => quickSwap.update({ item: this._equippedItem! }));
			}
		});

		player.sim.showQuickSwapChangeEmitter.on(() => {
			this.quickSwapEnchantPopover?.tooltip?.[this.player.sim.getShowQuickSwap() ? 'enable' : 'disable']();
			this.quickSwapGemPopover.forEach(quickSwap => quickSwap.tooltip?.[this.player.sim.getShowQuickSwap() ? 'enable' : 'disable']());
		});

		player.professionChangeEmitter.on(() => {
			if (!!this._equippedItem) {
				this.player.setWowheadData(this._equippedItem, this.itemElem.iconElem);
			}
		});
	}

	createGearData(): GearData {
		return {
			equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
				this.player.equipItem(eventID, this.slot, equippedItem);
			},
			getEquippedItem: () => this.player.getEquippedItem(this.slot),
			changeEvent: this.player.gearChangeEmitter,
		};
	}

	get item(): EquippedItem | null {
		return this._equippedItem;
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.itemElem.clear();
		// Clear quick swap gems array since gem sockets are rerendered every time
		this.quickSwapGemPopover = [];
		this.itemElem.nameElem.textContent = slotNames.get(this.slot) ?? '';
		setItemQualityCssClass(this.itemElem.nameElem, null);

		if (!!newItem) {
			this.itemElem.update(newItem);
		} else {
			this.itemElem.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		}

		this._equippedItem = newItem;
		this.onUpdateCallbacks.forEach(callback => callback());
	}

	onUpdate(callback: () => void) {
		this.onUpdateCallbacks.push(callback);
	}

	openSelectorModal(selectedTab: SelectorModalTabs) {
		this.gearPicker.selectorModal.openTab(this.slot, selectedTab, this.createGearData());
	}

	private addQuickGemHelpers() {
		if (!this._equippedItem) return;
		const openGemDetailTab = (socketIdx: number) => this.openSelectorModal(`Gem${socketIdx + 1}` as SelectorModalTabs);
		this.itemElem.socketsElem?.forEach(element => {
			const socketIdx = Number(element.dataset.socketIdx) || 0;
			element.addEventListener('click', event => {
				event.preventDefault();
				openGemDetailTab(0);
			});
			const popover = addQuickGemPopover(this.player, element, this._equippedItem!, this.slot, socketIdx, () => openGemDetailTab(socketIdx));
			if (!this.player.sim.getShowQuickSwap()) popover.tooltip?.disable();
			this.quickSwapGemPopover.push(popover);
		});
	}

	private addQuickEnchantHelpers() {
		if (!this._equippedItem) return;
		const openEnchantSelector = () => this.openSelectorModal(SelectorModalTabs.Enchants);
		this.itemElem.enchantElem.addEventListener('click', event => {
			event?.preventDefault();
			openEnchantSelector();
		});
		this.quickSwapEnchantPopover = addQuickEnchantPopover(this.player, this.itemElem.enchantElem, this._equippedItem, this.slot, openEnchantSelector);
		if (!this.player.sim.getShowQuickSwap()) this.quickSwapEnchantPopover.tooltip?.disable();
	}
}

export class IconItemSwapPicker extends Component {
	private readonly iconAnchor: HTMLAnchorElement;
	private readonly socketsContainerElem: HTMLElement;
	private readonly player: Player<any>;
	private readonly slot: ItemSlot;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'icon-picker-root');
		this.rootElem.classList.add('icon-picker');
		this.player = player;
		this.slot = slot;

		this.iconAnchor = document.createElement('a');
		this.iconAnchor.classList.add('icon-picker-button');
		this.iconAnchor.target = '_blank';
		this.rootElem.prepend(this.iconAnchor);

		this.socketsContainerElem = document.createElement('div');
		this.socketsContainerElem.classList.add('item-picker-sockets-container');
		this.iconAnchor.appendChild(this.socketsContainerElem);

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
		this.socketsContainerElem.innerText = '';

		if (newItem) {
			this.iconAnchor.classList.add('active');

			newItem.asActionId().fillAndSet(this.iconAnchor, true, true);
			this.player.setWowheadData(newItem, this.iconAnchor);

			newItem.allSocketColors().forEach((socketColor, gemIdx) => {
				this.socketsContainerElem.appendChild(createGemContainer(socketColor, newItem.gems[gemIdx], gemIdx));
			});
		} else {
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
