import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { SortDirection } from '../../constants/other';
import { EP_TOOLTIP } from '../../constants/tooltips';
import { setItemQualityCssClass } from '../../css_utils';
import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { Class, GemColor, ItemLevelState, ItemQuality, ItemRandomSuffix, ItemSlot, ItemSpec } from '../../proto/common';
import { DatabaseFilters, RepFaction, UIEnchant as Enchant, UIGem as Gem, UIItem as Item, UIItem_FactionRestriction } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { getUniqueEnchantString } from '../../proto_utils/enchants';
import { EquippedItem, ReforgeData } from '../../proto_utils/equipped_item';
import { difficultyNames, professionNames, REP_FACTION_NAMES, REP_FACTION_QUARTERMASTERS, REP_LEVEL_NAMES } from '../../proto_utils/names';
import { getPVPSeasonFromItem, isPVPItem } from '../../proto_utils/utils';
import { Sim } from '../../sim';
import { SimUI } from '../../sim_ui';
import { EventID, TypedEvent } from '../../typed_event';
import { formatDeltaTextElem } from '../../utils';
import {
	makePhaseSelector,
	makeShow1hWeaponsSelector,
	makeShow2hWeaponsSelector,
	makeShowEPValuesSelector,
	makeShowMatchingGemsSelector,
} from '../inputs/other_inputs';
import { ItemNotice } from '../item_notice/item_notice';
import Toast from '../toast';
import { Clusterize } from '../virtual_scroll/clusterize';
import { FiltersMenu } from './filters_menu';
import { SelectorModalTabs } from './selector_modal';
import { createHeroicLabel } from './utils';

export interface ItemData<T extends ItemListType> {
	item: T;
	name: string | HTMLElement;
	id: number;
	actionId: ActionId;
	quality: ItemQuality;
	phase: number;
	ilvl?: number;
	ignoreEPFilter: boolean;
	heroic: boolean;
	onEquip: (eventID: EventID, item: T) => void;
}

interface ItemDataWithIdx<T extends ItemListType> {
	idx: number;
	data: ItemData<T>;
}

export interface GearData {
	equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => void;
	getEquippedItem: () => EquippedItem | null;
	changeEvent: TypedEvent<any>;
}

export type ItemListType = Item | Enchant | Gem | ReforgeData | ItemRandomSuffix | ItemLevelState;
enum ItemListSortBy {
	EP,
	ILVL,
}

export default class ItemList<T extends ItemListType> {
	private listElem: HTMLElement;
	private readonly simUI: SimUI;
	private readonly player: Player<any>;
	public id: string;
	public label: string;
	private slot: ItemSlot;
	private itemData: Array<ItemData<T>>;
	private itemsToDisplay: Array<number>;
	private currentFilters: DatabaseFilters;
	private searchInput: HTMLInputElement;
	private socketColor: GemColor;
	private computeEP: (item: T) => number;
	private equippedToItemFn: (equippedItem: EquippedItem | null) => T | null | undefined;
	private gearData: GearData;
	private tabContent: Element;
	private onItemClick: (itemData: ItemData<T>) => void;
	private scroller: Clusterize;

	private sortBy = ItemListSortBy.ILVL;
	private sortDirection = SortDirection.DESC;

	constructor(
		id: string,
		parent: HTMLElement,
		simUI: SimUI,
		currentSlot: ItemSlot,
		currentTab: SelectorModalTabs,
		player: Player<any>,
		label: string,
		gearData: GearData,
		itemData: Array<ItemData<T>>,
		socketColor: GemColor,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => T | null | undefined,
		onRemove: (eventID: EventID) => void,
		onItemClick: (itemData: ItemData<T>) => void,
	) {
		this.id = id;
		this.label = label;
		this.simUI = simUI;
		this.player = player;
		this.itemData = itemData;
		this.socketColor = socketColor;
		this.computeEP = computeEP;
		this.equippedToItemFn = equippedToItemFn;
		this.onItemClick = onItemClick;

		this.slot = currentSlot;
		this.gearData = gearData;
		this.currentFilters = this.player.sim.getFilters();

		const selected = label === currentTab;
		const itemLabel = label === SelectorModalTabs.Reforging ? 'Reforge' : 'Item';

		const sortByIlvl = (event: MouseEvent) => {
			event.preventDefault();
			this.sort(ItemListSortBy.ILVL);
		};
		const sortByEP = (event: MouseEvent) => {
			event.preventDefault();
			this.sort(ItemListSortBy.EP);
		};

		const searchRef = ref<HTMLInputElement>();
		const epButtonRef = ref<HTMLButtonElement>();
		const filtersButtonRef = ref<HTMLButtonElement>();
		const showEpValuesRef = ref<HTMLDivElement>();
		const phaseSelectorRef = ref<HTMLDivElement>();
		const matchingGemsRef = ref<HTMLDivElement>();
		const show1hWeaponRef = ref<HTMLDivElement>();
		const show2hWeaponRef = ref<HTMLDivElement>();
		const modalListRef = ref<HTMLUListElement>();
		const removeButtonRef = ref<HTMLButtonElement>();
		const compareLabelRef = ref<HTMLElement>();

		const showEPOptions = ![ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2].includes(currentSlot);

		this.tabContent = (
			<div id={this.id} className={`selector-modal-tab-pane tab-pane fade ${selected ? 'active show' : ''}`}>
				<div className="selector-modal-filters">
					<input ref={searchRef} className="selector-modal-search form-control" type="text" placeholder="Search..." />
					{label === SelectorModalTabs.Items && (
						<button ref={filtersButtonRef} className="selector-modal-filters-button btn btn-primary">
							Filters
						</button>
					)}
					<div ref={phaseSelectorRef} className="selector-modal-phase-selector" />
					<div ref={show1hWeaponRef} className="sim-input selector-modal-boolean-option selector-modal-show-1h-weapons hide" />
					<div ref={show2hWeaponRef} className="sim-input selector-modal-boolean-option selector-modal-show-2h-weapons hide" />
					<div ref={matchingGemsRef} className="sim-input selector-modal-boolean-option selector-modal-show-matching-gems" />
					{showEPOptions && <div ref={showEpValuesRef} className="sim-input selector-modal-boolean-option selector-modal-show-ep-values" />}
					<button ref={removeButtonRef} className="selector-modal-remove-button btn btn-danger">
						Unequip Item
					</button>
				</div>
				<div className="selector-modal-list-labels">
					<span className="item-label">
						<small>{itemLabel}</small>
					</span>
					{label === SelectorModalTabs.Items && (
						<label className="source-label">
							<small>Source</small>
						</label>
					)}
					{(label === SelectorModalTabs.Items || label === SelectorModalTabs.Upgrades) && (
						<label className="ilvl-label interactive" onclick={sortByIlvl}>
							<small>ILvl</small>
						</label>
					)}
					{showEPOptions && (
						<span className="ep-label interactive" onclick={sortByEP}>
							<small>EP</small>
							<i className="fa-solid fa-plus-minus fa-2xs"></i>
							<button ref={epButtonRef} className="btn btn-link p-0 ms-1">
								<i className="far fa-question-circle fa-lg"></i>
							</button>
						</span>
					)}
					<span className="favorite-label"></span>
					<span ref={compareLabelRef} className="compare-label hide"></span>
				</div>
				<ul ref={modalListRef} className="selector-modal-list"></ul>
			</div>
		);

		parent.appendChild(this.tabContent);

		if (this.label === SelectorModalTabs.Items) {
			this.bindToggleCompare(compareLabelRef.value!);
		}

		if (
			label === SelectorModalTabs.Items &&
			player.getPlayerClass().weaponTypes.length > 0 &&
			(currentSlot === ItemSlot.ItemSlotMainHand || (currentSlot === ItemSlot.ItemSlotOffHand && player.getClass() === Class.ClassWarrior))
		) {
			if (show1hWeaponRef.value) makeShow1hWeaponsSelector(show1hWeaponRef.value, player.sim);
			if (show2hWeaponRef.value) makeShow2hWeaponsSelector(show2hWeaponRef.value, player.sim);
		}

		if (showEPOptions) {
			if (showEpValuesRef.value) makeShowEPValuesSelector(showEpValuesRef.value, player.sim);

			tippy(epButtonRef.value!, {
				content: EP_TOOLTIP,
			});
		}

		if (matchingGemsRef.value) {
			makeShowMatchingGemsSelector(matchingGemsRef.value, player.sim);
			if (!label.startsWith('Gem')) {
				matchingGemsRef.value?.classList.add('hide');
			}
		}

		// TODO: Turn this back on once we have proper phase data
		if (phaseSelectorRef.value) makePhaseSelector(phaseSelectorRef.value, player.sim);

		if (label === SelectorModalTabs.Items) {
			const filtersMenu = new FiltersMenu(parent, player, currentSlot);
			filtersButtonRef.value?.addEventListener('click', () => filtersMenu.open());
		}

		this.listElem = modalListRef.value!;
		this.itemsToDisplay = [];

		this.scroller = new Clusterize(
			{
				getNumberOfRows: () => {
					return this.itemsToDisplay.length;
				},
				generateRows: (startIdx, endIdx) => {
					const items = [];
					for (let i = startIdx; i < endIdx; ++i) {
						if (i >= this.itemsToDisplay.length) break;
						items.push(this.createItemElem({ idx: this.itemsToDisplay[i], data: this.itemData[this.itemsToDisplay[i]] }));
					}
					return items;
				},
			},
			{
				rows: [],
				scroll_elem: this.listElem,
				content_elem: this.listElem,
				item_height: 56,
				show_no_data_row: false,
				no_data_text: '',
				tag: 'li',
				rows_in_block: 16,
				blocks_in_cluster: 2,
			},
		);

		const removeButton = removeButtonRef.value;
		if (removeButton) {
			removeButton.addEventListener('click', _event => {
				onRemove(TypedEvent.nextEventID());
			});

			switch (label) {
				case SelectorModalTabs.Enchants:
					removeButton.textContent = 'Remove Enchant';
					break;
				case SelectorModalTabs.Tinkers:
					removeButton.textContent = 'Remove Tinkers';
					break;
				case SelectorModalTabs.Reforging:
					removeButton.textContent = 'Remove Reforge';
					break;
				case SelectorModalTabs.RandomSuffixes:
					removeButton.textContent = 'Remove Random Suffix';
					break;
				case SelectorModalTabs.Upgrades:
					removeButton.textContent = 'Remove Upgrade';
					break;
				case SelectorModalTabs.Gem1:
				case SelectorModalTabs.Gem2:
				case SelectorModalTabs.Gem3:
					removeButton.textContent = 'Remove Gem';
					break;
			}
		}

		this.updateSelected();

		this.searchInput = searchRef.value!;
		this.searchInput.addEventListener('input', () => this.applyFilters());
	}

	public sizeRefresh() {
		this.scroller.refresh(true);
		this.applyFilters();
	}

	public dispose() {
		this.scroller.dispose();
	}

	private getItemIdByItemType(item: ItemListType | null | undefined) {
		switch (this.label) {
			case SelectorModalTabs.Enchants:
				return (item as Enchant)?.effectId;
			case SelectorModalTabs.Tinkers:
				return (item as Enchant)?.effectId;
			case SelectorModalTabs.Reforging:
				return (item as ReforgeData)?.reforge!.id;
			case SelectorModalTabs.Items:
			case SelectorModalTabs.Gem1:
			case SelectorModalTabs.Gem2:
			case SelectorModalTabs.Gem3:
			case SelectorModalTabs.RandomSuffixes:
				return (item as Item | Gem | ItemRandomSuffix)?.id;
			case SelectorModalTabs.Upgrades:
				return item as ItemLevelState;
			default:
				return null;
		}
	}

	public updateSelected() {
		const newEquippedItem = this.gearData.getEquippedItem();
		const newItem = this.equippedToItemFn(newEquippedItem);
		const newItemId = this.getItemIdByItemType(newItem);
		const newEP = newItem !== undefined && newItem !== null ? this.computeEP(newItem) : 0;

		this.scroller.elementUpdate(item => {
			const idx = (item as HTMLElement).dataset.idx!;
			const itemData = this.itemData[parseFloat(idx)];

			if (itemData.id === newItemId) item.classList.add('active');
			else item.classList.remove('active');

			const epDeltaElem = item.querySelector<HTMLSpanElement>('.selector-modal-list-item-ep-delta');
			if (epDeltaElem) {
				epDeltaElem.textContent = '';
				if (itemData.item !== null) {
					const listItemEP = this.computeEP(itemData.item);
					if (newEP !== listItemEP) {
						formatDeltaTextElem(epDeltaElem, newEP, listItemEP, 0);
					}
				}
			}
		});
	}

	public applyFilters() {
		this.currentFilters = this.player.sim.getFilters();
		let itemIdxs = new Array<number>(this.itemData.length);
		for (let i = 0; i < this.itemData.length; ++i) {
			itemIdxs[i] = i;
		}

		const currentEquippedItem = this.player.getEquippedItem(this.slot);

		if (this.label === SelectorModalTabs.Items) {
			itemIdxs = this.player.filterItemData(itemIdxs, i => this.itemData[i].item as unknown as Item, this.slot);
		} else if (this.label === SelectorModalTabs.Enchants || this.label === SelectorModalTabs.Tinkers) {
			itemIdxs = this.player.filterEnchantData(itemIdxs, i => this.itemData[i].item as unknown as Enchant, this.slot, currentEquippedItem);
		} else if (this.label === SelectorModalTabs.Gem1 || this.label === SelectorModalTabs.Gem2 || this.label === SelectorModalTabs.Gem3) {
			itemIdxs = this.player.filterGemData(itemIdxs, i => this.itemData[i].item as unknown as Gem, this.slot, this.socketColor);
		}

		itemIdxs = itemIdxs.filter(i => {
			const listItemData = this.itemData[i];

			if (listItemData.phase > this.player.sim.getPhase()) {
				return false;
			}

			if (!!this.searchInput.value.length) {
				const formatQuery = (value: string) => value.toLowerCase().replaceAll(/[^a-zA-Z0-9\s]/g, '');

				const searchQuery = formatQuery(this.searchInput.value).split(' ');
				const name = formatQuery(listItemData.name.toString());

				let include = true;
				searchQuery.some(v => {
					if (!name.includes(v)) include = false;
				});
				if (!include) {
					return false;
				}
			}

			return true;
		});

		if ([ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2].includes(this.slot)) {
			// Trinket EP is weird so just sort by ilvl instead.
			this.sortBy = ItemListSortBy.ILVL;
		} else {
			this.sortBy = ItemListSortBy.EP;
		}

		itemIdxs = this.sortIdxs(itemIdxs);

		this.itemsToDisplay = itemIdxs;
		this.scroller.update();

		this.hideOrShowEPValues();
	}

	public sort(sortBy: ItemListSortBy) {
		if (this.sortBy === sortBy) {
			this.sortDirection = 1 - this.sortDirection;
		} else {
			this.sortDirection = SortDirection.DESC;
		}
		this.sortBy = sortBy;
		this.itemsToDisplay = this.sortIdxs(this.itemsToDisplay);
		this.scroller.update();
	}

	private sortIdxs(itemIdxs: Array<number>): number[] {
		let sortFn = (itemA: T, itemB: T) => {
			const first = (this.sortDirection === SortDirection.DESC ? itemB : itemA) as unknown as Item;
			const second = (this.sortDirection === SortDirection.DESC ? itemA : itemB) as unknown as Item;
			const diff = this.computeEP(first as T) - this.computeEP(second as T);
			// if EP is same, sort by ilvl
			if (Math.abs(diff) < 0.01)
				return (first.scalingOptions?.[ItemLevelState.Base].ilvl || first.ilvl) - (second.scalingOptions?.[ItemLevelState.Base].ilvl || second.ilvl);
			return diff;
		};
		switch (this.sortBy) {
			case ItemListSortBy.ILVL:
				sortFn = (itemA: T, itemB: T) => {
					const first = (this.sortDirection === SortDirection.DESC ? itemB : itemA) as unknown as Item;
					const second = (this.sortDirection === SortDirection.DESC ? itemA : itemB) as unknown as Item;
					return (
						(first.scalingOptions?.[ItemLevelState.Base].ilvl || first.ilvl) - (second.scalingOptions?.[ItemLevelState.Base].ilvl || second.ilvl)
					);
				};
				break;
		}

		return itemIdxs.sort((dataA, dataB) => {
			const itemA = this.itemData[dataA];
			const itemB = this.itemData[dataB];
			if (this.isItemFavorited(itemA) && !this.isItemFavorited(itemB)) return -1;
			if (this.isItemFavorited(itemB) && !this.isItemFavorited(itemA)) return 1;

			return sortFn(itemA.item, itemB.item);
		});
	}

	public hideOrShowEPValues() {
		const labels = this.tabContent.querySelectorAll('.ep-label');
		const container = this.tabContent.querySelectorAll('.selector-modal-list');
		const show = this.player.sim.getShowEPValues();
		const display = show ? '' : 'none';

		for (const label of labels) {
			(label as HTMLElement).style.display = display;
		}

		for (const c of container) {
			if (show) c.classList.remove('hide-ep');
			else c.classList.add('hide-ep');
		}
	}

	private createItemElem(item: ItemDataWithIdx<T>): JSX.Element {
		const itemData = item.data;
		const itemEP = this.computeEP(itemData.item);
		const equippedItem = this.equippedToItemFn(this.gearData.getEquippedItem());
		const hasItem = equippedItem !== null && equippedItem !== undefined;
		const equippedItemID = this.getItemIdByItemType(equippedItem);
		const equippedItemEP = hasItem ? this.computeEP(equippedItem) : 0;

		const labelCellElem = ref<HTMLDivElement>();
		const nameElem = ref<HTMLLabelElement>();
		const anchorElem = ref<HTMLAnchorElement>();
		const iconElem = ref<HTMLImageElement>();
		const favoriteElem = ref<HTMLButtonElement>();
		const favoriteIconElem = ref<HTMLElement>();
		const compareContainer = ref<HTMLDivElement>();
		const compareButton = ref<HTMLButtonElement>();

		const listItemElem = (
			<li className={`selector-modal-list-item ${equippedItemID === itemData.id ? 'active' : ''}`} dataset={{ idx: item.idx.toString() }}>
				<div className="selector-modal-list-label-cell gap-1" ref={labelCellElem}>
					<a className="selector-modal-list-item-link" ref={anchorElem} dataset={{ whtticon: 'false' }}>
						<img className="selector-modal-list-item-icon" ref={iconElem}></img>
						<label className="selector-modal-list-item-name" ref={nameElem}>
							{typeof itemData.name === 'string' ? itemData.name : itemData.name.cloneNode(true)}
							{itemData.heroic && createHeroicLabel()}
						</label>
					</a>
				</div>
				{this.label === SelectorModalTabs.Items && (
					<div className="selector-modal-list-item-source-container">{this.getSourceInfo(itemData.item as unknown as Item, this.player.sim)}</div>
				)}
				{(this.label === SelectorModalTabs.Items || this.label === SelectorModalTabs.Upgrades) && (
					<div className="selector-modal-list-item-ilvl-container">{itemData.ilvl || (itemData.item as unknown as Item).ilvl}</div>
				)}
				{![ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2].includes(this.slot) && (
					<div className="selector-modal-list-item-ep">
						<span className="selector-modal-list-item-ep-value">
							{itemEP < 9.95 ? itemEP.toFixed(1).toString() : Math.round(itemEP).toString()}
						</span>
						<span
							className="selector-modal-list-item-ep-delta"
							ref={e => hasItem && equippedItemEP !== itemEP && formatDeltaTextElem(e, equippedItemEP, itemEP, 0)}
						/>
					</div>
				)}
				<div className="selector-modal-list-item-favorite-container">
					<button className="selector-modal-list-item-favorite btn btn-link p-0" ref={favoriteElem}>
						<i ref={favoriteIconElem} className="far fa-star fa-xl" />
					</button>
				</div>
				<div ref={compareContainer} className="selector-modal-list-item-compare-container hide">
					<button className="selector-modal-list-item-compare btn btn-link p-0" ref={compareButton}>
						<i className="fas fa-arrow-right-arrow-left fa-xl" />
					</button>
				</div>
			</li>
		);

		const toggleFavorite = (isFavorite: boolean) => {
			const filters = this.player.sim.getFilters();

			let favMethodName: keyof DatabaseFilters;
			let favId;
			switch (this.label) {
				case SelectorModalTabs.Items:
					favMethodName = 'favoriteItems';
					favId = itemData.id;
					break;
				case SelectorModalTabs.Enchants:
					favMethodName = 'favoriteEnchants';
					favId = getUniqueEnchantString(itemData.item as unknown as Enchant);
					break;
				case SelectorModalTabs.Tinkers:
					favMethodName = 'favoriteEnchants';
					favId = getUniqueEnchantString(itemData.item as unknown as Enchant);
					break;
				case SelectorModalTabs.Gem1:
				case SelectorModalTabs.Gem2:
				case SelectorModalTabs.Gem3:
					favMethodName = 'favoriteGems';
					favId = itemData.id;
					break;
				case SelectorModalTabs.RandomSuffixes:
					favMethodName = 'favoriteRandomSuffixes';
					favId = itemData.id;
					break;
				case SelectorModalTabs.Reforging:
					favMethodName = 'favoriteReforges';
					favId = itemData.id;
					break;
				default:
					return;
			}

			if (isFavorite) {
				filters[favMethodName].push(favId as never);
			} else {
				const favIdx = filters[favMethodName].indexOf(favId as never);
				if (favIdx !== -1) {
					filters[favMethodName].splice(favIdx, 1);
				}
			}

			favoriteElem.value!.classList.toggle('text-brand');
			favoriteIconElem.value!.classList.toggle('fas');
			favoriteIconElem.value!.classList.toggle('far');
			listItemElem.dataset.fav = isFavorite.toString();
			this.player.sim.setFilters(TypedEvent.nextEventID(), filters);
		};
		favoriteElem.value!.addEventListener('click', () => toggleFavorite(listItemElem.dataset.fav === 'false'));

		const isFavorite = this.isItemFavorited(itemData);
		if (isFavorite) {
			favoriteElem.value!.classList.add('text-brand');
			favoriteIconElem.value?.classList.add('fas');
			listItemElem.dataset.fav = 'true';
		} else {
			favoriteIconElem.value?.classList.add('far');
			listItemElem.dataset.fav = 'false';
		}

		const favoriteTooltip = tippy(favoriteElem.value!);
		const toggleFavoriteTooltipContent = (isFavorited: boolean) => favoriteTooltip.setContent(isFavorited ? 'Remove from favorites' : 'Add to favorites');
		toggleFavoriteTooltipContent(listItemElem.dataset.fav === 'true');

		if (this.label === SelectorModalTabs.Items) {
			const batchSimTooltip = tippy(compareButton.value!);

			this.bindToggleCompare(compareContainer.value!);
			const simUI = this.simUI instanceof IndividualSimUI ? this.simUI : null;
			if (simUI) {
				const checkHasItem = () => simUI.bt?.hasItem(ItemSpec.create({ id: itemData.id }));
				const toggleCompareButtonState = () => {
					const hasItem = checkHasItem();
					batchSimTooltip.setContent(hasItem ? 'Remove from Batch Sim' : 'Add to Batch Sim');
					compareButton.value!.classList[hasItem ? 'add' : 'remove']('text-brand');
				};

				toggleCompareButtonState();
				simUI.bt?.itemsChangedEmitter.on(() => {
					toggleCompareButtonState();
				});

				compareButton.value!.addEventListener('click', () => {
					const hasItem = checkHasItem();
					simUI.bt?.[hasItem ? 'removeItem' : 'addItem'](ItemSpec.create({ id: itemData.id }));

					new Toast({
						delay: 1000,
						variant: 'success',
						body: (
							<>
								<strong>{itemData.name}</strong> was {hasItem ? <>removed from the batch</> : <>added to the batch</>}.
							</>
						),
					});
					// TODO: should we open the bulk sim UI or should we run in the background showing progress, and then sort the items in the picker?
				});
			}
		}

		anchorElem.value!.addEventListener('click', (event: Event) => {
			event.preventDefault();
			if (event.target === favoriteElem.value) return false;
			this.onItemClick(itemData);
		});

		itemData.actionId.fill().then(filledId => {
			filledId.setWowheadHref(anchorElem.value!);
			iconElem.value!.src = filledId.iconUrl;
		});

		setItemQualityCssClass(nameElem.value!, itemData.quality);

		const notice = new ItemNotice(this.player, { itemId: itemData.id });
		if (notice.hasNotice) labelCellElem.value?.appendChild(notice.rootElem);

		return listItemElem;
	}

	private isItemFavorited(itemData: ItemData<T>): boolean {
		if (this.label === SelectorModalTabs.Items) {
			return this.currentFilters.favoriteItems.includes(itemData.id);
		} else if (this.label === SelectorModalTabs.Enchants) {
			return this.currentFilters.favoriteEnchants.includes(getUniqueEnchantString(itemData.item as unknown as Enchant));
		} else if (this.label === SelectorModalTabs.Tinkers) {
			return this.currentFilters.favoriteEnchants.includes(getUniqueEnchantString(itemData.item as unknown as Enchant));
		} else if (this.label.startsWith('Gem')) {
			return this.currentFilters.favoriteGems.includes(itemData.id);
		} else if (this.label === SelectorModalTabs.RandomSuffixes) {
			return this.currentFilters.favoriteRandomSuffixes.includes(itemData.id);
		} else if (this.label === SelectorModalTabs.Reforging) {
			return this.currentFilters.favoriteReforges.includes(itemData.id);
		}
		return false;
	}

	private getSourceInfo(item: Item, sim: Sim): JSX.Element {
		const makeAnchor = (href: string, inner: string | JSX.Element) => {
			return (
				<a href={href} target="_blank" dataset={{ whtticon: 'false' }}>
					<small>{inner}</small>
				</a>
			);
		};

		if (!item.sources?.length) {
			if (item.randomSuffixOptions.length) {
				return makeAnchor(`${ActionId.makeItemUrl(item.id)}#dropped-by`, 'World Drop');
			} else if (isPVPItem(item)) {
				const season = getPVPSeasonFromItem(item);
				if (!season) return <></>;

				return makeAnchor(
					ActionId.makeItemUrl(item.id),
					<span>
						{season}
						<br />
						PVP
					</span>,
				);
			}
			return <></>;
		}

		let source = item.sources[0];
		if (source.source.oneofKind === 'crafted') {
			const src = source.source.crafted;

			if (src.spellId) {
				return makeAnchor(ActionId.makeSpellUrl(src.spellId), professionNames.get(src.profession) ?? 'Unknown');
			}
			return makeAnchor(ActionId.makeItemUrl(item.id), professionNames.get(src.profession) ?? 'Unknown');
		} else if (source.source.oneofKind === 'drop') {
			const src = source.source.drop;
			const zone = sim.db.getZone(src.zoneId);
			const npc = sim.db.getNpc(src.npcId);
			if (!zone) {
				console.error('No zone found for item:', item);
				return <></>;
			}

			const category = src.category ? ` - ${src.category}` : '';
			if (npc) {
				return makeAnchor(
					ActionId.makeNpcUrl(npc.id),
					<span>
						{zone.name} ({difficultyNames.get(src.difficulty) ?? 'Unknown'})
						<br />
						{npc.name + category}
					</span>,
				);
			} else if (src.otherName) {
				return makeAnchor(
					ActionId.makeZoneUrl(zone.id),
					<span>
						{zone.name}
						<br />
						{src.otherName}
					</span>,
				);
			}
			return makeAnchor(ActionId.makeZoneUrl(zone.id), zone.name);
		} else if (source.source.oneofKind === 'quest' && source.source.quest.name) {
			const src = source.source.quest;
			return makeAnchor(
				ActionId.makeQuestUrl(src.id),
				<span>
					Quest
					{item.factionRestriction === UIItem_FactionRestriction.ALLIANCE_ONLY && (
						<img src="/mop/assets/img/alliance.png" className="ms-1" width="15" height="15" />
					)}
					{item.factionRestriction === UIItem_FactionRestriction.HORDE_ONLY && (
						<img src="/mop/assets/img/horde.png" className="ms-1" width="15" height="15" />
					)}
					<br />
					{src.name}
				</span>,
			);
		} else if ((source = item.sources.find(source => source.source.oneofKind === 'rep') ?? source).source.oneofKind === 'rep') {
			const factionNames = item.sources
				.filter(source => source.source.oneofKind === 'rep')
				.map(source =>
					source.source.oneofKind === 'rep' ? REP_FACTION_NAMES[source.source.rep.repFactionId] : REP_FACTION_NAMES[RepFaction.RepFactionUnknown],
				);
			const src = source.source.rep;
			const npcId = REP_FACTION_QUARTERMASTERS[src.repFactionId];
			return makeAnchor(
				ActionId.makeNpcUrl(npcId),
				<>
					{factionNames.map(name => (
						<span>
							{name}
							{item.factionRestriction === UIItem_FactionRestriction.ALLIANCE_ONLY && (
								<img src="/mop/assets/img/alliance.png" className="ms-1" width="15" height="15" />
							)}
							{item.factionRestriction === UIItem_FactionRestriction.HORDE_ONLY && (
								<img src="/mop/assets/img/horde.png" className="ms-1" width="15" height="15" />
							)}
							<br />
						</span>
					))}
					<span>{REP_LEVEL_NAMES[src.repLevel]}</span>
				</>,
			);
		} else if (isPVPItem(item)) {
			const season = getPVPSeasonFromItem(item);
			if (!season) return <></>;

			return makeAnchor(
				ActionId.makeItemUrl(item.id),
				<span>
					{season}
					<br />
					PVP
				</span>,
			);
		} else if (source.source.oneofKind === 'soldBy') {
			const src = source.source.soldBy;
			return makeAnchor(
				ActionId.makeNpcUrl(src.npcId),
				<span>
					Sold by
					<br />
					{src.npcName}
				</span>,
			);
		}
		return <></>;
	}

	private bindToggleCompare(element: Element) {
		const toggleCompare = () => element.classList[!this.player.sim.getShowExperimental() ? 'add' : 'remove']('hide');
		toggleCompare();
		this.player.sim.showExperimentalChangeEmitter.on(() => {
			toggleCompare();
		});
	}
}
