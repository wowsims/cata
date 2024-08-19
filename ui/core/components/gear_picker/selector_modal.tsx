import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Player, ReforgeData } from '../../player';
import { GemColor, ItemQuality, ItemRandomSuffix, ItemSlot } from '../../proto/common';
import { UIEnchant as Enchant, UIGem as Gem, UIItem as Item } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { gemMatchesSocket, getEmptyGemSocketIconUrl } from '../../proto_utils/gems';
import { shortSecondaryStatNames, slotNames } from '../../proto_utils/names';
import { Stats } from '../../proto_utils/stats';
import { SimUI } from '../../sim_ui';
import { EventID, TypedEvent } from '../../typed_event';
import { mod, randomUUID, sanitizeId } from '../../utils';
import { BaseModal } from '../base_modal';
import GearPicker from './gear_picker';
import ItemList, { GearData, ItemData, ItemListType } from './item_list';
import { createGemContainer, getEmptySlotIconUrl } from './utils';

export enum SelectorModalTabs {
	Items = 'Items',
	RandomSuffixes = 'Random Suffix',
	Enchants = 'Enchants',
	Reforging = 'Reforging',
	Gem1 = 'Gem1',
	Gem2 = 'Gem2',
	Gem3 = 'Gem3',
}

type SelectorModalOptions = {
	// This will add a unique ID to the modal, allowing multiple of the same modals to exist
	id: string;
	// Prevents rendering of certail tabs
	disabledTabs?: SelectorModalTabs[];
};
export default class SelectorModal extends BaseModal {
	private readonly simUI: SimUI;
	private player: Player<any>;
	private gearPicker: GearPicker | undefined;
	private ilists: ItemList<ItemListType>[] = [];

	private readonly itemSlotTabElems: HTMLElement[] = [];
	private readonly titleElem: HTMLElement;
	private readonly tabsElem: HTMLElement;
	private readonly contentElem: HTMLElement;

	private currentSlot: ItemSlot = ItemSlot.ItemSlotHead;
	private currentTab: SelectorModalTabs = SelectorModalTabs.Items;
	private disabledTabs: SelectorModalTabs[] = [];
	private options: SelectorModalOptions;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, gearPicker?: GearPicker, options?: Partial<SelectorModalOptions>) {
		super(parent, 'selector-modal', { disposeOnClose: false, size: 'xl' });

		this.simUI = simUI;
		this.player = player;
		this.gearPicker = gearPicker;
		this.options = { id: randomUUID(), ...options };
		this.disabledTabs = this.options.disabledTabs || [];

		this.addItemSlotTabs();

		this.header!.insertAdjacentElement(
			'afterbegin',
			<div>
				<h6 className="selector-modal-title" />
				<ul className="nav nav-tabs selector-modal-tabs"></ul>
			</div>,
		);

		this.body.appendChild(<div className="tab-content selector-modal-tab-content"></div>);

		this.titleElem = this.rootElem.querySelector<HTMLElement>('.selector-modal-title')!;
		this.tabsElem = this.rootElem.querySelector<HTMLElement>('.selector-modal-tabs')!;
		this.contentElem = this.rootElem.querySelector<HTMLElement>('.selector-modal-tab-content')!;

		this.body.appendChild(
			<div className="d-flex align-items-center form-text mt-3">
				<i className="fas fa-circle-exclamation fa-xl me-2"></i>
				<span>
					If gear is missing, check the selected phase and your gear filters.
					<br />
					If the problem persists, save any un-saved data, click the
					<i className="fas fa-cog mx-1"></i>
					to open your sim options, then click the "Restore Defaults".
				</span>
			</div>,
		);
	}

	openTab(selectedSlot: ItemSlot, selectedTab: SelectorModalTabs, gearData: GearData) {
		this.titleElem.textContent = slotNames.get(selectedSlot) ?? '';
		this.setData(selectedSlot, selectedTab, gearData);
		this.setActiveItemSlotTab(selectedSlot);
		this.open();
	}

	onShow() {
		if (this.gearPicker) {
			// Allow you to switch between gear picker slots with the up and down arrows
			const switchToPreviousItemSlotTab = this.switchToPreviousItemSlotTab.bind(this);
			const switchToNextItemSlotTab = this.switchToNextItemSlotTab.bind(this);

			document.addEventListener('keydown', switchToPreviousItemSlotTab);
			document.addEventListener('keydown', switchToNextItemSlotTab);

			this.addOnHideCallback(() => document.removeEventListener('keydown', switchToPreviousItemSlotTab));
			this.addOnHideCallback(() => document.removeEventListener('keydown', switchToNextItemSlotTab));
		}
	}

	private setData(selectedSlot: ItemSlot, selectedTab: SelectorModalTabs, gearData: GearData) {
		this.tabsElem.innerText = '';
		this.contentElem.innerText = '';
		this.ilists = [];

		const equippedItem = gearData.getEquippedItem();

		const eligibleItems = this.player.getItems(selectedSlot);
		const eligibleEnchants = this.player.getEnchants(selectedSlot);
		const eligibleReforges = equippedItem?.item ? this.player.getAvailableReforgings(equippedItem.getWithRandomSuffixStats()) : [];
		// If the enchant tab is selected but the item has no eligible enchants, default to items
		// If the reforge tab is selected but the item has no eligible reforges, default to items
		// If a gem tab is selected but the item has no eligible sockets, default to items
		if (
			(selectedTab === SelectorModalTabs.Enchants && !eligibleEnchants.length) ||
			(selectedTab === SelectorModalTabs.Reforging && !eligibleReforges.length) ||
			([SelectorModalTabs.Gem1, SelectorModalTabs.Gem2, SelectorModalTabs.Gem3].includes(selectedTab) &&
				equippedItem?.numSockets(this.player.isBlacksmithing()) === 0)
		) {
			selectedTab = SelectorModalTabs.Items;
		}

		this.currentTab = selectedTab;
		this.currentSlot = selectedSlot;

		const hasItemTab = !this.disabledTabs?.includes(SelectorModalTabs.Items);
		if (hasItemTab)
			this.addTab<Item>({
				id: sanitizeId(`${this.options.id}-${SelectorModalTabs.Items}`),
				label: SelectorModalTabs.Items,
				gearData,
				itemData: eligibleItems.map(item => {
					return {
						item: item,
						id: item.id,
						actionId: ActionId.fromItem(item),
						name: item.name,
						quality: item.quality,
						heroic: item.heroic,
						phase: item.phase,
						baseEP: this.player.computeItemEP(item, selectedSlot),
						ignoreEPFilter: false,
						onEquip: (eventID, item) => {
							const equippedItem = gearData.getEquippedItem();
							if (equippedItem) {
								gearData.equipItem(eventID, equippedItem.withItem(item));
							} else {
								gearData.equipItem(eventID, new EquippedItem(item));
							}
						},
					};
				}),
				computeEP: (item: Item) => this.player.computeItemEP(item, selectedSlot),
				equippedToItemFn: (equippedItem: EquippedItem | null) => equippedItem?.item,
				onRemove: (eventID: number) => {
					gearData.equipItem(eventID, null);
					this.removeTabs('Gem');
					this.removeTabs(SelectorModalTabs.RandomSuffixes);
				},
			});

		const hasEnchantTab = !this.disabledTabs?.includes(SelectorModalTabs.Enchants);
		if (hasEnchantTab)
			this.addTab<Enchant>({
				id: sanitizeId(`${this.options.id}-${SelectorModalTabs.Enchants}`),
				label: SelectorModalTabs.Enchants,
				gearData,
				itemData: eligibleEnchants.map(enchant => {
					return {
						item: enchant,
						id: enchant.effectId,
						actionId: enchant.itemId ? ActionId.fromItemId(enchant.itemId) : ActionId.fromSpellId(enchant.spellId),
						name: enchant.name,
						quality: enchant.quality,
						phase: enchant.phase || 1,
						baseEP: this.player.computeStatsEP(new Stats(enchant.stats)),
						ignoreEPFilter: true,
						heroic: false,
						onEquip: (eventID, enchant) => {
							const equippedItem = gearData.getEquippedItem();
							if (equippedItem) gearData.equipItem(eventID, equippedItem.withEnchant(enchant));
						},
					};
				}),
				computeEP: (enchant: Enchant) => this.player.computeEnchantEP(enchant),
				equippedToItemFn: (equippedItem: EquippedItem | null) => equippedItem?.enchant,
				onRemove: (eventID: number) => {
					const equippedItem = gearData.getEquippedItem();
					if (equippedItem) gearData.equipItem(eventID, equippedItem.withEnchant(null));
				},
			});

		const hasRandomSuffixTab = !this.disabledTabs?.includes(SelectorModalTabs.RandomSuffixes);
		if (hasRandomSuffixTab) this.addRandomSuffixTab(equippedItem, gearData);
		const hasReforgingTab = !this.disabledTabs?.includes(SelectorModalTabs.Reforging);
		if (hasReforgingTab) this.addReforgingTab(equippedItem, gearData);
		const hasGemsTab = ![SelectorModalTabs.Gem1, SelectorModalTabs.Gem2, SelectorModalTabs.Gem3].some(gem => this.disabledTabs?.includes(gem));
		if (hasGemsTab) this.addGemTabs(selectedSlot, equippedItem, gearData);

		this.ilists.find(list => selectedTab === list.label)?.sizeRefresh();
	}

	private addItemSlotTabs() {
		if (!this.gearPicker) {
			return;
		}

		this.dialog.prepend(
			<div className="gear-picker-modal-slots">
				{this.gearPicker.itemPickers.map(picker => {
					const anchorRef = ref<HTMLAnchorElement>();
					const wrapper = (
						<div className="item-picker-icon-wrapper" dataset={{ slot: picker.slot }}>
							<a
								ref={anchorRef}
								className="item-picker-icon"
								href="javascript:void(0)"
								onclick={(e: Event) => {
									e.preventDefault();
									if (picker.slot != this.currentSlot) {
										picker.openSelectorModal(this.currentTab);
									}
								}}
								dataset={{ whtticon: 'false' }}
							/>
						</div>
					) as HTMLElement;

					const setItemData = () => {
						if (picker.item) {
							this.player.setWowheadData(picker.item, anchorRef.value!);
							picker.item
								.asActionId()
								.fill()
								.then(filledId => {
									filledId.setBackgroundAndHref(anchorRef.value!);
								});
						} else {
							anchorRef.value!.style.backgroundImage = `url('${getEmptySlotIconUrl(picker.slot)}')`;
						}
					};
					setItemData();
					picker.onUpdate(() => setItemData());
					tippy(anchorRef.value!, {
						content: `Edit ${slotNames.get(picker.slot)}`,
						placement: 'left',
					});
					this.itemSlotTabElems.push(wrapper);
					return wrapper;
				})}
			</div>,
		);
	}

	private setActiveItemSlotTab(slot: ItemSlot) {
		this.itemSlotTabElems.forEach(elem => {
			if (elem.dataset.slot === slot.toString()) {
				elem.classList.add('active');
			} else if (elem.classList.contains('active')) {
				elem.classList.remove('active');
			}
		});
	}

	private switchToPreviousItemSlotTab(event: KeyboardEvent) {
		if (event.key === 'ArrowUp' && this.gearPicker) {
			event.preventDefault();
			const newSlot = mod(this.currentSlot - 1, Object.keys(ItemSlot).length / 2) as unknown as ItemSlot;
			this.gearPicker.itemPickers[newSlot].openSelectorModal(this.currentTab);
		}
	}

	private switchToNextItemSlotTab(event: KeyboardEvent) {
		if (event.key === 'ArrowDown' && this.gearPicker) {
			event.preventDefault();
			const newSlot = mod(this.currentSlot + 1, Object.keys(ItemSlot).length / 2) as unknown as ItemSlot;
			this.gearPicker.itemPickers[newSlot].openSelectorModal(this.currentTab);
		}
	}

	private addGemTabs(_slot: ItemSlot, equippedItem: EquippedItem | null, gearData: GearData) {
		if (!equippedItem) {
			return;
		}

		const socketBonusEP = this.player.computeStatsEP(new Stats(equippedItem.item.socketBonus)) / (equippedItem.item.gemSockets.length || 1);
		equippedItem.curSocketColors(this.player.isBlacksmithing()).forEach((socketColor, socketIdx) => {
			const label = SelectorModalTabs[`Gem${socketIdx + 1}` as keyof typeof SelectorModalTabs];
			this.addTab<Gem>({
				id: sanitizeId(`${this.options.id}-${label}`),
				label,
				gearData,
				itemData: this.player.getGems(socketColor).map((gem: Gem) => {
					return {
						item: gem,
						id: gem.id,
						actionId: ActionId.fromItemId(gem.id),
						name: gem.name,
						quality: gem.quality,
						phase: gem.phase,
						heroic: false,
						baseEP: this.player.computeStatsEP(new Stats(gem.stats)),
						ignoreEPFilter: true,
						onEquip: (eventID, gem) => {
							const equippedItem = gearData.getEquippedItem();
							if (equippedItem) gearData.equipItem(eventID, equippedItem.withGem(gem, socketIdx));
						},
					};
				}),
				computeEP: (gem: Gem) => {
					let gemEP = this.player.computeGemEP(gem);
					if (gemMatchesSocket(gem, socketColor)) {
						gemEP += socketBonusEP;
					}
					return gemEP;
				},
				equippedToItemFn: (equippedItem: EquippedItem | null) => equippedItem?.gems[socketIdx],
				onRemove: (eventID: number) => {
					const equippedItem = gearData.getEquippedItem();
					if (equippedItem) gearData.equipItem(eventID, equippedItem.withGem(null, socketIdx));
				},
				setTabContent: tabButton => {
					const gemContainer = createGemContainer(socketColor, null, socketIdx);
					tabButton.appendChild(gemContainer);
					tabButton.classList.add('selector-modal-tab-gem');

					const gemElem = tabButton.querySelector<HTMLElement>('.gem-icon')!;
					const emptySocketUrl = getEmptyGemSocketIconUrl(socketColor);

					const updateGemIcon = () => {
						const equippedItem = gearData.getEquippedItem();
						const gem = equippedItem?.gems[socketIdx];

						if (gem) {
							gemElem.classList.remove('hide');
							ActionId.fromItemId(gem.id)
								.fill()
								.then(filledId => {
									gemElem.setAttribute('src', filledId.iconUrl);
								});
						} else {
							gemElem.classList.add('hide');
							gemElem.setAttribute('src', emptySocketUrl);
						}
					};

					gearData.changeEvent.on(updateGemIcon);
					this.addOnDisposeCallback(() => gearData.changeEvent.off(updateGemIcon));
					updateGemIcon();
				},
				socketColor,
			});
		});
	}

	private addRandomSuffixTab(equippedItem: EquippedItem | null, gearData: GearData) {
		if (!equippedItem || !equippedItem.item.randomSuffixOptions.length) {
			return;
		}

		const itemProto = equippedItem.item;

		this.addTab<ItemRandomSuffix>({
			id: sanitizeId(`${this.options.id}-${SelectorModalTabs.RandomSuffixes}`),
			label: SelectorModalTabs.RandomSuffixes,
			gearData,
			itemData: this.player.getRandomSuffixes(itemProto).map((randomSuffix: ItemRandomSuffix) => {
				return {
					item: randomSuffix,
					id: randomSuffix.id,
					actionId: ActionId.fromRandomSuffix(itemProto, randomSuffix),
					name: randomSuffix.name,
					quality: itemProto.quality,
					phase: itemProto.phase,
					heroic: false,
					baseEP: this.player.computeRandomSuffixEP(randomSuffix),
					ignoreEPFilter: true,
					onEquip: (eventID, randomSuffix) => {
						const equippedItem = gearData.getEquippedItem();
						if (equippedItem) {
							gearData.equipItem(eventID, equippedItem.withItem(equippedItem.item).withRandomSuffix(randomSuffix));
						}
					},
				};
			}),
			computeEP: (randomSuffix: ItemRandomSuffix) => this.player.computeRandomSuffixEP(randomSuffix),
			equippedToItemFn: (equippedItem: EquippedItem | null) => equippedItem?.randomSuffix,
			onRemove: (eventID: number) => {
				const equippedItem = gearData.getEquippedItem();
				if (equippedItem) {
					gearData.equipItem(eventID, equippedItem.withItem(equippedItem.item).withRandomSuffix(null));
				}
				this.removeTabs(SelectorModalTabs.Reforging);
			},
		});
	}

	private addReforgingTab(equippedItem: EquippedItem | null, gearData: GearData) {
		if (!equippedItem || (equippedItem.hasRandomSuffixOptions() && !equippedItem.randomSuffix)) {
			return;
		}

		const itemProto = equippedItem.item;

		this.addTab<ReforgeData>({
			id: sanitizeId(`${this.options.id}-${SelectorModalTabs.Reforging}`),
			label: SelectorModalTabs.Reforging,
			gearData,
			itemData: this.player.getAvailableReforgings(equippedItem).map(reforgeData => {
				return {
					item: reforgeData,
					id: reforgeData.id,
					actionId: ActionId.fromReforge(itemProto, reforgeData.reforge),
					name: (
						<div>
							<span className="reforge-value negative">
								{reforgeData.fromAmount} {shortSecondaryStatNames.get(reforgeData.fromStat)}
							</span>
							<span className="reforge-value positive">
								+{reforgeData.toAmount} {shortSecondaryStatNames.get(reforgeData.toStat)}
							</span>
						</div>
					) as HTMLElement,
					quality: ItemQuality.ItemQualityCommon,
					phase: itemProto.phase,
					heroic: false,
					baseEP: this.player.computeReforgingEP(reforgeData),
					ignoreEPFilter: true,
					onEquip: (eventID, reforgeData) => {
						const equippedItem = gearData.getEquippedItem();
						if (equippedItem) {
							gearData.equipItem(eventID, equippedItem.withReforge(reforgeData.reforge));
						}
					},
				};
			}),
			computeEP: (reforge: ReforgeData) => this.player.computeReforgingEP(reforge),
			equippedToItemFn: (equippedItem: EquippedItem | null) =>
				equippedItem?.reforge ? this.player.getReforgeData(equippedItem, equippedItem.reforge) : null,
			onRemove: (eventID: number) => {
				const equippedItem = gearData.getEquippedItem();
				if (equippedItem) {
					gearData.equipItem(eventID, equippedItem.withItem(equippedItem.item).withRandomSuffix(equippedItem._randomSuffix));
				}
			},
		});
	}

	/**
	 * Adds one of the tabs for the item selector menu.
	 *
	 * T is expected to be Item, Enchant, or Gem. Tab menus for all 3 looks extremely
	 * similar so this function uses extra functions to do it generically.
	 */
	private addTab<T extends ItemListType>({
		id,
		label,
		gearData,
		itemData,
		computeEP,
		equippedToItemFn,
		onRemove,
		setTabContent,
		socketColor,
	}: {
		id: string;
		label: SelectorModalTabs;
		gearData: GearData;
		itemData: ItemData<T>[];
		computeEP: (item: T) => number;
		equippedToItemFn: (equippedItem: EquippedItem | null) => T | null | undefined;
		onRemove: (eventID: EventID) => void;
		setTabContent?: (tabElem: HTMLButtonElement) => void;
		socketColor?: GemColor;
	}) {
		if (!itemData.length) {
			return;
		}
		const selected = label === this.currentTab;
		const tabButton = ref<HTMLButtonElement>();
		this.tabsElem.appendChild(
			<li className="nav-item">
				<button
					ref={tabButton}
					className={clsx('nav-link selector-modal-item-tab', selected && 'active')}
					dataset={{
						label,
						bsToggle: 'tab',
						bsTarget: `#${id}`,
					}}
					attributes={{
						role: 'tab',
						'aria-selected': selected,
					}}
				/>
			</li>,
		);

		if (setTabContent) {
			setTabContent(tabButton.value!);
		} else {
			tabButton.value!.textContent = label;
		}

		const ilist = new ItemList(
			id,
			this.contentElem,
			this.simUI,
			this.currentSlot,
			this.currentTab,
			this.player,
			label,
			gearData,
			itemData,
			socketColor || GemColor.GemColorUnknown,
			computeEP,
			equippedToItemFn,
			onRemove,
			itemData => {
				const item = itemData;
				itemData.onEquip(TypedEvent.nextEventID(), item.item);

				const isItemChange = Item.is(item.item);
				const isRandomSuffixChange = !!item.actionId.randomSuffixId;
				// If the item changes, then gem slots and random suffix options will also change, so remove and recreate these tabs.
				if (isItemChange || isRandomSuffixChange) {
					if (!isRandomSuffixChange) {
						this.removeTabs(SelectorModalTabs.RandomSuffixes);
						this.addRandomSuffixTab(gearData.getEquippedItem(), gearData);
					}

					this.removeTabs(SelectorModalTabs.Reforging);
					this.addReforgingTab(gearData.getEquippedItem(), gearData);

					this.removeTabs('Gem');
					this.addGemTabs(this.currentSlot, gearData.getEquippedItem(), gearData);
				}
			},
		);

		const invokeUpdate = () => {
			ilist.updateSelected();
		};
		const applyFilter = () => {
			ilist.applyFilters();
		};
		const hideOrShowEPValues = () => {
			ilist.hideOrShowEPValues();
		};
		// Add event handlers
		gearData.changeEvent.on(invokeUpdate);

		this.player.sim.phaseChangeEmitter.on(applyFilter);
		this.player.sim.filtersChangeEmitter.on(applyFilter);
		this.player.sim.showEPValuesChangeEmitter.on(hideOrShowEPValues);

		this.addOnDisposeCallback(() => {
			gearData.changeEvent.off(invokeUpdate);
			this.player.sim.phaseChangeEmitter.off(applyFilter);
			this.player.sim.filtersChangeEmitter.off(applyFilter);
			this.player.sim.showEPValuesChangeEmitter.off(hideOrShowEPValues);
			ilist.dispose();
		});

		tabButton.value!.addEventListener('click', _event => {
			this.currentTab = label;
		});
		tabButton.value!.addEventListener('shown.bs.tab', _event => {
			ilist.sizeRefresh();
		});

		this.ilists.push(ilist as unknown as ItemList<ItemListType>);
	}

	private removeTabs(labelSubstring: string) {
		const tabElems = [...this.tabsElem.querySelectorAll<HTMLElement>('.selector-modal-item-tab')].filter(
			tab => tab.dataset?.label?.includes(labelSubstring),
		);

		const contentElems = tabElems.map(tabElem => document.querySelector(tabElem.dataset.bsTarget!)).filter(tabElem => Boolean(tabElem));
		tabElems.forEach(elem => elem.parentElement?.remove());
		contentElems.forEach(elem => elem!.remove());
	}
}
