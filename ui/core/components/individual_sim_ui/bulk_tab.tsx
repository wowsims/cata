import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../../css_utils';
import { IndividualSimUI } from '../../individual_sim_ui';
import { BulkComboResult, BulkSettings, ItemSpecWithSlot, ProgressMetrics, TalentLoadout } from '../../proto/api';
import { EquipmentSpec, GemColor, ItemSlot, ItemSpec, SimDatabase, SimEnchant, SimGem, SimItem } from '../../proto/common';
import { SavedTalents, UIEnchant, UIGem, UIItem, UIItem_FactionRestriction } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { Database } from '../../proto_utils/database';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { getEmptyGemSocketIconUrl } from '../../proto_utils/gems';
import { Stats } from '../../proto_utils/stats';
import { canEquipItem, getEligibleItemSlots } from '../../proto_utils/utils';
import { TypedEvent } from '../../typed_event';
import { EventID } from '../../typed_event.js';
import { cloneChildren, noop } from '../../utils';
import { WorkerProgressCallback } from '../../worker_pool';
import { BaseModal } from '../base_modal';
import { BooleanPicker } from '../boolean_picker';
import { Component } from '../component';
import { ContentBlock } from '../content_block';
import { GearData, ItemData, ItemList, ItemRenderer, SelectorModal, SelectorModalTabs } from '../gear_picker/gear_picker';
import { Importer } from '../importers';
import { ResultsViewer } from '../results_viewer';
import { SimTab } from '../sim_tab';

export class BulkGearJsonImporter extends Importer {
	private readonly simUI: IndividualSimUI<any>;
	private readonly bulkUI: BulkTab;
	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, bulkUI: BulkTab) {
		super(parent, simUI, 'Bag Item Import', true);
		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.descriptionElem.appendChild(
			<>
				<p>Import bag items from a JSON file, which can be created by the WowSimsExporter in-game AddOn.</p>
				<p>To import, upload the file or paste the text below, then click, 'Import'.</p>
			</>,
		);
	}

	async onImport(data: string) {
		try {
			const equipment = EquipmentSpec.fromJsonString(data, { ignoreUnknownFields: true });
			if (equipment?.items?.length > 0) {
				const db = await Database.loadLeftoversIfNecessary(equipment);
				const items = equipment.items.filter(spec => spec.id > 0 && db.lookupItemSpec(spec));
				if (items.length > 0) {
					this.bulkUI.addItems(items);
				}
			}
			this.close();
		} catch (e: any) {
			console.warn(e);
			alert(e.toString());
		}
	}
}

class BulkSimResultRenderer {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, result: BulkComboResult, baseResult: BulkComboResult) {
		const dpsDelta = result.unitMetrics!.dps!.avg! - baseResult.unitMetrics!.dps!.avg;

		const equipButtonRef = ref<HTMLButtonElement>();
		const dpsDeltaRef = ref<HTMLDivElement>();
		const itemsContainerRef = ref<HTMLDivElement>();
		parent.appendChild(
			<>
				<div className="results-sim">
					<div className="bulk-result-body-dps bulk-items-text-line results-sim-dps damage-metrics">
						<span className="topline-result-avg">{this.formatDps(result.unitMetrics!.dps!.avg)}</span>

						<span ref={dpsDeltaRef} className={clsx(dpsDelta >= 0 ? 'bulk-result-header-positive' : 'bulk-result-header-negative')}>
							{this.formatDpsDelta(dpsDelta)}
						</span>

						<p className="talent-loadout-text">
							{result.talentLoadout && typeof result.talentLoadout === 'object' ? (
								typeof result.talentLoadout.name === 'string' && <>Talent loadout used: {result.talentLoadout.name}</>
							) : (
								<>Current talents</>
							)}
						</p>
					</div>
				</div>
				<div ref={itemsContainerRef} className="bulk-gear-combo"></div>
				{!!result.itemsAdded?.length && (
					<button ref={equipButtonRef} className="btn btn-primary bulk-equipit">
						Equip
					</button>
				)}
			</>,
		);

		if (!!result.itemsAdded?.length) {
			equipButtonRef.value?.addEventListener('click', () => {
				result.itemsAdded.forEach(itemAdded => {
					const item = simUI.sim.db.lookupItemSpec(itemAdded.item!);
					simUI.player.equipItem(TypedEvent.nextEventID(), itemAdded.slot, item);
					simUI.simHeader.activateTab('gear-tab');
				});
			});

			const items = (<></>) as HTMLElement;
			for (const is of result.itemsAdded) {
				const itemContainer = (<div className="bulk-result-item" />) as HTMLElement;
				const item = simUI.sim.db.lookupItemSpec(is.item!);
				const renderer = new ItemRenderer(items, itemContainer, simUI.player);
				renderer.update(item!);
				renderer.nameElem.appendChild(<a className="bulk-result-item-slot">{this.itemSlotName(is)}</a>);
				items.appendChild(itemContainer);
			}
			itemsContainerRef.value?.appendChild(items);
		} else if (!result.talentLoadout || typeof result.talentLoadout !== 'object') {
			dpsDeltaRef.value?.classList.add('hide');
			parent.appendChild(<p>No changes - this is your currently equipped gear!</p>);
		}
	}

	private formatDps(dps: number): string {
		return (Math.round(dps * 100) / 100).toFixed(2);
	}

	private formatDpsDelta(delta: number): string {
		return (delta >= 0 ? '+' : '') + this.formatDps(delta);
	}

	private itemSlotName(is: ItemSpecWithSlot): string {
		return JSON.parse(ItemSpecWithSlot.toJsonString(is, { emitDefaultValues: true }))['slot'].replace('ItemSlot', '');
	}
}

export class BulkItemPicker extends Component {
	private readonly itemElem: ItemRenderer;
	readonly simUI: IndividualSimUI<any>;
	readonly bulkUI: BulkTab;
	readonly index: number;

	protected item: EquippedItem;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, bulkUI: BulkTab, item: EquippedItem, index: number) {
		super(parent, 'bulk-item-picker');
		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.index = index;
		this.item = item;
		this.itemElem = new ItemRenderer(parent, this.rootElem, simUI.player);

		this.simUI.sim.waitForInit().then(() => {
			this.setItem(item);
			const slot = getEligibleItemSlots(this.item.item)[0];
			const eligibleEnchants = this.simUI.sim.db.getEnchants(slot);
			const eligibleReforges = this.item?.item ? this.simUI.player.getAvailableReforgings(this.item.getWithRandomSuffixStats()) : [];
			const eligibleRandomSuffixes = this.item.item.randomSuffixOptions;

			const openEnchantGemSelector = (event: Event) => {
				event.preventDefault();

				if (!!eligibleEnchants.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.Enchants, this.createGearData());
				} else if (!!eligibleRandomSuffixes.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.RandomSuffixes, this.createGearData());
				} else if (!!eligibleReforges.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.Reforging, this.createGearData());
				} else if (!!this.item._gems.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.Gem1, this.createGearData());
				}

				const destroyItemButton = <button className="btn btn-danger">Remove from Batch</button>;
				destroyItemButton.addEventListener('click', () => {
					bulkUI.setItems(
						bulkUI.getItems().filter((_, idx) => {
							return idx != this.index;
						}),
					);
					this.bulkUI.selectorModal.close();
				});
				const closeX = this.bulkUI.selectorModal.header?.querySelector('.close-button');
				if (!!closeX) {
					this.bulkUI.selectorModal.header?.insertBefore(destroyItemButton, closeX);
				}
			};

			this.itemElem.iconElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.nameElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.enchantElem.addEventListener('click', openEnchantGemSelector);
		});
	}

	setItem(newItem: EquippedItem | null) {
		this.itemElem.clear();
		if (!!newItem) {
			this.itemElem.update(newItem);
			this.item = newItem;
		} else {
			this.itemElem.rootElem.style.opacity = '30%';
			this.itemElem.iconElem.style.backgroundImage = `url('/cata/assets/item_slots/empty.jpg')`;
			this.itemElem.nameElem.textContent = 'Add new item (not implemented)';
			this.itemElem.rootElem.style.alignItems = 'center';
		}
	}

	private createGearData(): GearData {
		const changeEvent = new TypedEvent<void>();
		return {
			equipItem: (_, equippedItem: EquippedItem | null) => {
				if (equippedItem) {
					const allItems = this.bulkUI.getItems();
					allItems[this.index] = equippedItem.asSpec();
					this.item = equippedItem;
					this.bulkUI.setItems(allItems);
					changeEvent.emit(TypedEvent.nextEventID());
				}
			},
			getEquippedItem: () => this.item,
			changeEvent: changeEvent,
		};
	}
}

export class BulkTab extends SimTab {
	readonly simUI: IndividualSimUI<any>;

	readonly itemsChangedEmitter = new TypedEvent<void>();

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	readonly column1: HTMLElement = this.buildColumn(1, 'raid-settings-col');

	protected items: Array<ItemSpec> = new Array<ItemSpec>();

	private pendingResults: ResultsViewer;
	private pendingDiv: HTMLDivElement;

	// TODO: Make a real options probably
	private doCombos: boolean;
	private fastMode: boolean;
	private autoGem: boolean;
	private simTalents: boolean;
	private autoEnchant: boolean;
	private defaultGems: SimGem[];
	private savedTalents: TalentLoadout[];
	private gemIconElements: HTMLImageElement[];
	readonly selectorModal: SelectorModal;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<any>) {
		super(parentElem, simUI, { identifier: 'bulk-tab', title: 'Batch' });
		this.simUI = simUI;

		this.leftPanel = (<div className="bulk-tab-left tab-panel-left">{this.column1}</div>) as HTMLDivElement;
		this.rightPanel = (<div className="bulk-tab-right tab-panel-right" />) as HTMLDivElement;

		this.pendingDiv = (<div className="results-pending-overlay d-flex hide" />) as HTMLDivElement;
		this.pendingResults = new ResultsViewer(this.pendingDiv);
		this.pendingResults.hideAll();
		this.selectorModal = new SelectorModal(this.simUI.rootElem, this.simUI, this.simUI.player, undefined, {
			id: 'bulk-selector-modal',
			disabledTabs: [SelectorModalTabs.Items],
		});

		this.contentContainer.appendChild(
			<>
				{this.leftPanel}
				{this.rightPanel}
				{this.pendingDiv}
			</>,
		);

		this.doCombos = true;
		this.fastMode = true;
		this.autoGem = true;
		this.autoEnchant = true;
		this.savedTalents = [];
		this.simTalents = false;
		this.defaultGems = [UIGem.create(), UIGem.create(), UIGem.create(), UIGem.create()];
		this.gemIconElements = [];
		this.buildTabContent();

		this.simUI.sim.waitForInit().then(() => {
			this.loadSettings();
		});
	}

	private getSettingsKey(): string {
		return this.simUI.getStorageKey('bulk-settings.v1');
	}

	private loadSettings() {
		const storedSettings = window.localStorage.getItem(this.getSettingsKey());
		if (storedSettings != null) {
			const settings = BulkSettings.fromJsonString(storedSettings, {
				ignoreUnknownFields: true,
			});

			this.doCombos = settings.combinations;
			this.fastMode = settings.fastMode;
			this.autoEnchant = settings.autoEnchant;
			this.savedTalents = settings.talentsToSim;
			this.autoGem = settings.autoGem;
			this.simTalents = settings.simTalents;
			this.defaultGems = new Array<SimGem>(
				SimGem.create({ id: settings.defaultRedGem }),
				SimGem.create({ id: settings.defaultYellowGem }),
				SimGem.create({ id: settings.defaultBlueGem }),
				SimGem.create({ id: settings.defaultMetaGem }),
			);

			this.defaultGems.forEach((gem, idx) => {
				ActionId.fromItemId(gem.id)
					.fill()
					.then(filledId => {
						if (gem.id) {
							this.gemIconElements[idx].src = filledId.iconUrl;
							this.gemIconElements[idx].classList.remove('hide');
						}
					});
			});
		}
	}

	private storeSettings() {
		const settings = this.createBulkSettings();
		const setStr = BulkSettings.toJsonString(settings, { enumAsInteger: true });
		window.localStorage.setItem(this.getSettingsKey(), setStr);
	}

	protected createBulkSettings(): BulkSettings {
		return BulkSettings.create({
			items: this.items,
			// TODO(Riotdog-GehennasEU): Make all of these configurable.
			// For now, it's always constant iteration combinations mode for "sim my bags".
			combinations: this.doCombos,
			fastMode: this.fastMode,
			autoEnchant: this.autoEnchant,
			autoGem: this.autoGem,
			simTalents: this.simTalents,
			talentsToSim: this.savedTalents,
			defaultRedGem: this.defaultGems[0].id,
			defaultYellowGem: this.defaultGems[1].id,
			defaultBlueGem: this.defaultGems[2].id,
			defaultMetaGem: this.defaultGems[3].id,
			iterationsPerCombo: this.simUI.sim.getIterations(), // TODO(Riotdog-GehennasEU): Define a new UI element for the iteration setting.
		});
	}

	protected createBulkItemsDatabase(): SimDatabase {
		const itemsDb = SimDatabase.create();
		for (const is of this.items) {
			const item = this.simUI.sim.db.lookupItemSpec(is);
			if (!item) {
				throw new Error(`item with ID ${is.id} not found in database`);
			}
			itemsDb.items.push(SimItem.fromJson(UIItem.toJson(item.item), { ignoreUnknownFields: true }));
			if (item.enchant) {
				itemsDb.enchants.push(
					SimEnchant.fromJson(UIEnchant.toJson(item.enchant), {
						ignoreUnknownFields: true,
					}),
				);
			}
			for (const gem of item.gems) {
				if (gem) {
					itemsDb.gems.push(SimGem.fromJson(UIGem.toJson(gem), { ignoreUnknownFields: true }));
				}
			}
		}
		for (const gem of this.defaultGems) {
			if (gem.id > 0) {
				itemsDb.gems.push(gem);
			}
		}
		return itemsDb;
	}

	addItem(item: ItemSpec) {
		this.addItems([item]);
	}
	addItems(items: ItemSpec[]) {
		this.items = [...(this.items || []), ...items];
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	setItems(items: ItemSpec[]) {
		this.items = items;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	removeItem(item: ItemSpec) {
		const indexToRemove = this.items.findIndex(i => ItemSpec.equals(i, item));
		if (indexToRemove === -1) return;
		this.items.splice(indexToRemove, 1);
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}
	clearItems() {
		this.items = new Array<ItemSpec>();
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	hasItem(item: ItemSpec) {
		return this.items.some(i => ItemSpec.equals(i, item));
	}

	getItems(): Array<ItemSpec> {
		const result = new Array<ItemSpec>();
		this.items.forEach(spec => {
			result.push(ItemSpec.clone(spec));
		});
		return result;
	}

	setCombinations(doCombos: boolean) {
		this.doCombos = doCombos;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	setFastMode(fastMode: boolean) {
		this.fastMode = fastMode;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	protected async runBulkSim(onProgress: WorkerProgressCallback) {
		this.pendingResults.setPending();

		try {
			await this.simUI.sim.runBulkSim(this.createBulkSettings(), this.createBulkItemsDatabase(), onProgress);
		} catch (e) {
			this.simUI.handleCrash(e);
		}
	}

	protected buildTabContent() {
		const itemsBlock = new ContentBlock(this.column1, 'bulk-items', {
			header: { title: 'Items' },
		});
		itemsBlock.bodyElement.classList.add('gear-picker-root', 'gear-picker-root-bulk');

		const itemTextIntro = (
			<div className="bulk-items-text-line">
				<i>
					Notice: This is under very early but active development and experimental. You may also need to update your WoW AddOn if you want to import
					your bags.
				</i>
			</div>
		);

		const itemList = (<div className="tab-panel-col bulk-gear-combo" />) as HTMLElement;

		this.itemsChangedEmitter.on(() => {
			const items = (<></>) as HTMLElement;
			if (!!this.items.length) {
				itemTextIntro.textContent = 'The following items will be simmed together with your equipped gear.';
				for (let i = 0; i < this.items.length; ++i) {
					const spec = this.items[i];
					const item = this.simUI.sim.db.lookupItemSpec(spec);
					new BulkItemPicker(items, this.simUI, this, item!, i);
				}
			}
			itemList.replaceChildren(items);
		});

		itemsBlock.bodyElement.appendChild(
			<>
				{itemTextIntro}
				{itemList}
			</>,
		);

		this.clearItems();

		const resultsBlock = new ContentBlock(this.column1, 'bulk-results', {
			header: {
				title: 'Results',
				extraCssClasses: ['bulk-results-header'],
			},
		});

		resultsBlock.rootElem.hidden = true;
		resultsBlock.bodyElement.classList.add('gear-picker-root', 'gear-picker-root-bulk', 'tab-panel-col');

		this.simUI.sim.bulkSimStartEmitter.on(() => {
			resultsBlock.rootElem.hidden = true;
		});

		this.simUI.sim.bulkSimResultEmitter.on((_, bulkSimResult) => {
			resultsBlock.rootElem.hidden = bulkSimResult.results.length == 0;
			resultsBlock.bodyElement.replaceChildren();

			for (const r of bulkSimResult.results) {
				const resultBlock = new ContentBlock(resultsBlock.bodyElement, 'bulk-result', {
					header: { title: '' },
					bodyClasses: ['bulk-results-body'],
				});
				new BulkSimResultRenderer(resultBlock.bodyElement, this.simUI, r, bulkSimResult.equippedGearResult!);
			}
		});

		const settingsBlock = new ContentBlock(this.rightPanel, 'bulk-settings', {
			header: { title: 'Setup' },
		});

		const bulkSimButton = (<button className="btn btn-primary w-100 bulk-settings-button">Simulate Batch</button>) as HTMLButtonElement;
		bulkSimButton.addEventListener('click', () => {
			this.pendingDiv.classList.remove('hide');
			this.leftPanel.classList.add('blurred');
			this.rightPanel.classList.add('blurred');

			const defaultState = cloneChildren(bulkSimButton);
			bulkSimButton.disabled = true;
			bulkSimButton.classList.add('disabled');
			bulkSimButton.replaceChildren(
				<>
					<i className="fa fa-spinner fa-spin" /> Running
				</>,
			);

			let simStart = new Date().getTime();
			let lastTotal = 0;
			let rounds = 0;
			let currentRound = 0;
			let combinations = 0;

			this.runBulkSim((progressMetrics: ProgressMetrics) => {
				const msSinceStart = new Date().getTime() - simStart;
				const iterPerSecond = progressMetrics.completedIterations / (msSinceStart / 1000);

				if (combinations === 0) {
					combinations = progressMetrics.totalSims;
				}
				if (this.fastMode) {
					if (rounds === 0 && progressMetrics.totalSims > 0) {
						rounds = Math.ceil(Math.log(progressMetrics.totalSims / 20) / Math.log(2)) + 1;
						currentRound = 1;
					}
					if (progressMetrics.totalSims < lastTotal) {
						currentRound += 1;
						simStart = new Date().getTime();
					}
				}

				this.setSimProgress(progressMetrics, iterPerSecond, currentRound, rounds, combinations);
				lastTotal = progressMetrics.totalSims;

				if (!!progressMetrics.finalBulkResult) {
					// reset state
					this.pendingDiv.classList.add('hide');
					this.leftPanel.classList.remove('blurred');
					this.rightPanel.classList.remove('blurred');

					this.pendingResults.hideAll();
					bulkSimButton.disabled = false;
					bulkSimButton.classList.remove('disabled');
					bulkSimButton.replaceChildren(...defaultState);
				}
			});
		});

		const importButton = (
			<button className="btn btn-secondary w-100 bulk-settings-button">
				<i className="fa fa-download" /> Import From Bags
			</button>
		) as HTMLButtonElement;
		importButton.addEventListener('click', () => new BulkGearJsonImporter(this.simUI.rootElem, this.simUI, this));

		const importFavsButton = (
			<button className="btn btn-secondary w-100 bulk-settings-button">
				<i className="fa fa-download" /> Import Favorites
			</button>
		);
		importFavsButton.addEventListener('click', () => {
			const filters = this.simUI.player.sim.getFilters();
			const items = filters.favoriteItems.map(itemID => {
				return ItemSpec.create({ id: itemID });
			});
			this.addItems(items);
		});

		const searchInputRef = ref<HTMLInputElement>();
		const searchResultsRef = ref<HTMLUListElement>();
		const searchWrapper = (
			<div className="search-wrapper hide">
				<input ref={searchInputRef} type="text" placeholder="Search..." className="batch-search-input form-control hide" />
				<ul ref={searchResultsRef} className="batch-search-results hide"></ul>
			</div>
		);

		let allItems = Array<UIItem>();

		searchInputRef.value?.addEventListener('keyup', event => {
			if (event.key == 'Enter') {
				const toAdd = Array<ItemSpec>();
				searchResultsRef.value?.childNodes.forEach(node => {
					const strID = (node as HTMLElement).getAttribute('data-item-id');
					if (strID != null) {
						toAdd.push(ItemSpec.create({ id: Number.parseInt(strID) }));
					}
				});
				this.addItems(toAdd);
			}
		});

		searchInputRef.value?.addEventListener('input', _event => {
			const searchString = searchInputRef.value?.value || '';

			if (!searchString.length) {
				searchResultsRef.value?.replaceChildren();
				searchResultsRef.value?.classList.add('hide');
				return;
			}

			const pieces = searchString.split(' ');
			const items = <></>;

			allItems.forEach(item => {
				let matched = true;
				const lcName = item.name.toLowerCase();
				const lcSetName = item.setName.toLowerCase();

				pieces.forEach(piece => {
					const lcPiece = piece.toLowerCase();
					if (!lcName.includes(lcPiece) && !lcSetName.includes(lcPiece)) {
						matched = false;
						return false;
					}
					return true;
				});

				if (matched) {
					const itemRef = ref<HTMLLIElement>();
					const itemNameRef = ref<HTMLSpanElement>();
					items.appendChild(
						<li ref={itemRef} dataset={{ itemId: item.id.toString() }}>
							<span ref={itemNameRef}>{item.name}</span>
							{item.heroic && <span className="item-quality-uncommon">[H]</span>}
							{item.factionRestriction === UIItem_FactionRestriction.HORDE_ONLY && <span className="faction-horde">(H)</span>}
							{item.factionRestriction === UIItem_FactionRestriction.ALLIANCE_ONLY && <span className="faction-alliance">(A)</span>}
						</li>,
					);
					setItemQualityCssClass(itemNameRef.value!, item.quality);
					itemRef.value?.addEventListener('click', () => {
						this.addItems(Array<ItemSpec>(ItemSpec.create({ id: item.id })));
					});
				}
			});
			searchResultsRef.value?.replaceChildren(items);
			searchResultsRef.value?.classList.remove('hide');
		});

		const searchButtonContents = (
			<>
				<i className="fa fa-search" /> Add Item
			</>
		);

		const searchButton = <button className="btn btn-secondary w-100 bulk-settings-button">{searchButtonContents.cloneNode(true)}</button>;
		searchButton.addEventListener('click', () => {
			if (searchInputRef.value?.classList.contains('hide')) {
				searchWrapper?.classList.remove('hide');
				searchButton.replaceChildren(<>Close Search Results</>);
				allItems = this.simUI.sim.db.getAllItems().filter(item => canEquipItem(item, this.simUI.player.getPlayerSpec(), undefined));
				searchInputRef.value?.classList.remove('hide');
				if (searchInputRef.value?.value) searchResultsRef.value?.classList.remove('hide');
				searchInputRef.value?.focus();
			} else {
				searchButton.replaceChildren(searchButtonContents.cloneNode(true));
				searchWrapper?.classList.add('hide');
				searchInputRef.value?.classList.add('hide');
				searchResultsRef.value?.replaceChildren();
				searchResultsRef.value?.classList.add('hide');
			}
		});

		const clearButton = <button className="btn btn-secondary w-100 bulk-settings-button">Clear all</button>;
		clearButton.addEventListener('click', () => {
			this.clearItems();
			resultsBlock.rootElem.hidden = true;
			resultsBlock.bodyElement.replaceChildren();
		});

		// Talents to sim
		const talentsContainerRef = ref<HTMLDivElement>();
		const talentsToSimDiv = (
			<div className={clsx('talents-picker-container', !this.simTalents && 'hide')}>
				<label className="mb-2">Pick talents to sim (will increase time to sim)</label>
				<div ref={talentsContainerRef} className="talents-container"></div>
			</div>
		);

		const dataStr = window.localStorage.getItem(this.simUI.getSavedTalentsStorageKey());

		let jsonData;
		try {
			if (dataStr !== null) {
				jsonData = JSON.parse(dataStr);
			}
		} catch (e) {
			console.warn('Invalid json for local storage value: ' + dataStr);
		}

		const handleToggle = (element: HTMLElement, load: TalentLoadout) => {
			const exists = this.savedTalents.some(talent => talent.name === load.name); // Replace 'id' with your unique identifier
			// console.log('Exists:', exists);
			// console.log('Load Object:', load);
			// console.log('Saved Talents Before Update:', this.savedTalents);

			if (exists) {
				// If the object exists, find its index and remove it
				const indexToRemove = this.savedTalents.findIndex(talent => talent.name === load.name);
				this.savedTalents.splice(indexToRemove, 1);
				element?.classList.remove('active');
			} else {
				// If the object does not exist, add it
				this.savedTalents.push(load);
				element?.classList.add('active');
			}

			// console.log('Updated savedTalents:', this.savedTalents);
		};

		for (const name in jsonData) {
			try {
				const savedTalentLoadout = SavedTalents.fromJson(jsonData[name]);
				const loadout = {
					talentsString: savedTalentLoadout.talentsString,
					glyphs: savedTalentLoadout.glyphs,
					name: name,
				};

				const index = this.savedTalents.findIndex(talent => JSON.stringify(talent) === JSON.stringify(loadout));
				const talentChipRef = ref<HTMLDivElement>();
				const talentButtonRef = ref<HTMLButtonElement>();

				// console.log('Adding event for loadout', loadout);
				talentsContainerRef.value!.appendChild(
					<div ref={talentChipRef} className={clsx('saved-data-set-chip badge rounded-pill', index !== -1 && 'active')}>
						<button ref={talentButtonRef} className="saved-data-set-name">
							{name}
						</button>
					</div>,
				);
				talentButtonRef.value!.addEventListener('click', () => handleToggle(talentChipRef.value!, loadout));
			} catch (e) {
				console.log(e);
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}

		//////////////////////
		////////////////////////////////////

		// Default Gem Options
		const socketsContainerRef = ref<HTMLDivElement>();
		const defaultGemDiv = (
			<div className={clsx('default-gem-container', !this.autoGem && 'hide')}>
				<label className="mb-2">Defaults for Auto Gem</label>
				<div ref={socketsContainerRef} className="sockets-container"></div>
			</div>
		);

		Array<GemColor>(GemColor.GemColorRed, GemColor.GemColorYellow, GemColor.GemColorBlue, GemColor.GemColorMeta).forEach((socketColor, socketIndex) => {
			const gemContainerRef = ref<HTMLDivElement>();
			const gemIconRef = ref<HTMLImageElement>();
			const socketIconRef = ref<HTMLImageElement>();

			socketsContainerRef.value!.appendChild(
				<div ref={gemContainerRef} className="gem-socket-container">
					<img ref={gemIconRef} className="gem-icon hide" />
					<img ref={socketIconRef} className="socket-icon" />
				</div>,
			);

			this.gemIconElements.push(gemIconRef.value!);
			socketIconRef.value!.src = getEmptyGemSocketIconUrl(socketColor);

			let selector: GemSelectorModal;

			const onSelectHandler = (itemData: ItemData<UIGem>) => {
				this.defaultGems[socketIndex] = itemData.item;
				this.storeSettings();
				ActionId.fromItemId(itemData.id)
					.fill()
					.then(filledId => {
						if (itemData.id) {
							this.gemIconElements[socketIndex].src = filledId.iconUrl;
							this.gemIconElements[socketIndex].classList.remove('hide');
						}
					});
				selector.close();
			};

			const onRemoveHandler = () => {
				this.defaultGems[socketIndex] = UIGem.create();
				this.storeSettings();
				this.gemIconElements[socketIndex].classList.add('hide');
				this.gemIconElements[socketIndex].src = '';
				selector.close();
			};

			const openGemSelector = () => {
				if (!selector) selector = new GemSelectorModal(this.simUI.rootElem, this.simUI, socketColor, onSelectHandler, onRemoveHandler);
				selector.show();
			};

			this.gemIconElements[socketIndex].addEventListener('click', openGemSelector);
			gemContainerRef.value?.addEventListener('click', openGemSelector);
		});

		settingsBlock.bodyElement.appendChild(
			<>
				{bulkSimButton}
				{importButton}
				{importFavsButton}
				{searchButton}
				{searchWrapper}
				{clearButton}
				{defaultGemDiv}
				{talentsToSimDiv}
			</>,
		);

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-fast-mode',
			label: 'Fast Mode',
			labelTooltip: 'Fast mode reduces accuracy but will run faster.',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.fastMode,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.fastMode = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-combinations',
			label: 'Combinations',
			labelTooltip:
				'When checked bulk simulator will create all possible combinations of the items. When disabled trinkets and rings will still run all combinations becausee they have two slots to fill each.',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.doCombos,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.doCombos = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-auto-enchant',
			label: 'Auto Enchant',
			labelTooltip: 'When checked bulk simulator apply the current enchant for a slot to each replacement item it can.',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.autoEnchant,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.autoEnchant = value;
				defaultGemDiv.classList[value ? 'remove' : 'add']('hide');
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-auto-gem',
			label: 'Auto Gem',
			labelTooltip: 'When checked bulk simulator will fill any un-filled gem sockets with default gems.',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.autoGem,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.autoGem = value;
				defaultGemDiv.classList[value ? 'remove' : 'add']('hide');
			},
		});

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-sim-talents',
			label: 'Sim Talents',
			labelTooltip: 'When checked bulk simulator will sim chosen talent setups. Warning, it might cause the bulk sim to run for a lot longer',
			changedEvent: (_obj: BulkTab) => this.itemsChangedEmitter,
			getValue: _obj => this.simTalents,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.simTalents = value;
				talentsToSimDiv.classList[value ? 'remove' : 'add']('hide');
			},
		});
	}

	private setSimProgress(progress: ProgressMetrics, iterPerSecond: number, currentRound: number, rounds: number, combinations: number) {
		const secondsRemain = ((progress.totalIterations - progress.completedIterations) / iterPerSecond).toFixed();

		this.pendingResults.setContent(
			<div className="results-sim">
				<div>{combinations} total combinations.</div>
				<div>
					{rounds > 0 && (
						<>
							{currentRound} / {rounds} refining rounds
						</>
					)}
				</div>
				<div>
					{progress.completedSims} / {progress.totalSims}
					<br />
					simulations complete
				</div>
				<div>
					{progress.completedIterations} / {progress.totalIterations}
					<br />
					iterations complete
				</div>
				<div>{secondsRemain} seconds remaining.</div>
			</div>,
		);
	}
}

class GemSelectorModal extends BaseModal {
	private readonly simUI: IndividualSimUI<any>;

	private readonly contentElem: HTMLElement;
	private ilist: ItemList<UIGem> | null;
	private socketColor: GemColor;
	private onSelect: (itemData: ItemData<UIGem>) => void;
	private onRemove: () => void;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, socketColor: GemColor, onSelect: (itemData: ItemData<UIGem>) => void, onRemove: () => void) {
		super(parent, 'selector-modal', { disposeOnClose: false });

		this.simUI = simUI;
		this.onSelect = onSelect;
		this.onRemove = onRemove;
		this.socketColor = socketColor;
		this.ilist = null;

		window.scrollTo({ top: 0 });

		this.header!.insertAdjacentElement('afterbegin', <h6 className="selector-modal-title mb-3">Choose Default Gem</h6>);
		const contentRef = ref<HTMLDivElement>();
		this.body.appendChild(<div ref={contentRef} className="tab-content selector-modal-tab-content"></div>);
		this.contentElem = contentRef.value!;
	}

	show() {
		// construct item list the first time its opened.
		// This makes startup faster and also means we are sure to have item database loaded.
		if (!this.ilist) {
			this.ilist = new ItemList<UIGem>(
				'bulk-tab-gem-selector',
				this.contentElem,
				this.simUI,
				ItemSlot.ItemSlotHead,
				SelectorModalTabs.Gem1,
				this.simUI.player,
				SelectorModalTabs.Gem1,
				{
					equipItem: (_eventID: EventID, _equippedItem: EquippedItem | null) => {
						return;
					},
					getEquippedItem: () => null,
					changeEvent: new TypedEvent(), // FIXME
				},
				this.simUI.player.getGems(this.socketColor).map((gem: UIGem) => {
					return {
						item: gem,
						id: gem.id,
						actionId: ActionId.fromItemId(gem.id),
						name: gem.name,
						quality: gem.quality,
						phase: gem.phase,
						heroic: false,
						baseEP: this.simUI.player.computeStatsEP(new Stats(gem.stats)),
						ignoreEPFilter: true,
						onEquip: noop,
					};
				}),
				this.socketColor,
				gem => {
					return this.simUI.player.computeGemEP(gem);
				},
				() => {
					return null;
				},
				this.onRemove,
				this.onSelect,
			);

			this.ilist.sizeRefresh();

			const applyFilter = () => this.ilist?.applyFilters();

			const phaseChangeEvent = this.simUI.sim.phaseChangeEmitter.on(applyFilter);
			const filtersChangeChangeEvent = this.simUI.sim.filtersChangeEmitter.on(applyFilter);

			this.addOnDisposeCallback(() => {
				phaseChangeEvent.dispose();
				filtersChangeChangeEvent.dispose();
				this.ilist?.dispose();
			});
		}

		this.open();
	}
}
