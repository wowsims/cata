import { Button, convertToLegacy, Icon, Link } from '@wowsims/ui';
import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../individual_sim_ui';
import { BulkComboResult, BulkSettings, ItemSpecWithSlot, ProgressMetrics, TalentLoadout } from '../../proto/api';
import { EquipmentSpec, GemColor, ItemSlot, ItemSpec, SimDatabase, SimEnchant, SimGem, SimItem } from '../../proto/common';
import { SavedTalents, UIEnchant, UIGem, UIItem, UIItem_FactionRestriction } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { Database } from '../../proto_utils/database';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { getEmptyGemSocketIconUrl } from '../../proto_utils/gems';
import { canEquipItem, getEligibleItemSlots } from '../../proto_utils/utils';
import { TypedEvent } from '../../typed_event';
import { EventID } from '../../typed_event.js';
import { cloneChildren } from '../../utils';
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
		const dpsDivParent = document.createElement('div');
		dpsDivParent.classList.add('results-sim');

		const dpsDiv = document.createElement('div');
		dpsDiv.classList.add('bulk-result-body-dps', 'bulk-items-text-line', 'results-sim-dps', 'damage-metrics');
		dpsDivParent.appendChild(dpsDiv);

		const dpsNumber = document.createElement('span');
		dpsNumber.textContent = this.formatDps(result.unitMetrics!.dps!.avg);
		dpsNumber.classList.add('topline-result-avg');
		dpsDiv.appendChild(dpsNumber);

		const dpsDelta = result.unitMetrics!.dps!.avg! - baseResult.unitMetrics!.dps!.avg;
		const dpsDeltaSpan = document.createElement('span');
		dpsDeltaSpan.textContent = `${this.formatDpsDelta(dpsDelta)}`;
		dpsDeltaSpan.classList.add(dpsDelta >= 0 ? 'bulk-result-header-positive' : 'bulk-result-header-negative');
		dpsDiv.appendChild(dpsDeltaSpan);

		const itemsContainer = document.createElement('div');
		itemsContainer.classList.add('bulk-gear-combo');
		parent.appendChild(itemsContainer);
		parent.appendChild(dpsDivParent);

		const talentText = document.createElement('p');
		talentText.classList.add('talent-loadout-text');
		if (result.talentLoadout && typeof result.talentLoadout === 'object') {
			if (typeof result.talentLoadout.name === 'string') {
				talentText.textContent = 'Talent loadout used: ' + result.talentLoadout.name;
			}
		} else {
			talentText.textContent = 'Current talents';
		}

		dpsDiv.appendChild(talentText);
		if (result.itemsAdded && result.itemsAdded.length > 0) {
			const equipBtn = document.createElement('button');
			equipBtn.textContent = 'Equip';
			equipBtn.classList.add('btn', 'btn-primary', 'bulk-equipit');
			equipBtn.onclick = () => {
				result.itemsAdded.forEach(itemAdded => {
					const item = simUI.sim.db.lookupItemSpec(itemAdded.item!);
					simUI.player.equipItem(TypedEvent.nextEventID(), itemAdded.slot, item);
					simUI.simHeader.activateTab('gear-tab');
				});
			};

			parent.appendChild(equipBtn);

			for (const is of result.itemsAdded) {
				const item = simUI.sim.db.lookupItemSpec(is.item!);
				const renderer = new ItemRenderer(parent, itemsContainer, simUI.player);
				renderer.update(item!);

				const p = document.createElement('a');
				p.classList.add('bulk-result-item-slot');
				p.textContent = this.itemSlotName(is);
				renderer.nameElem.appendChild(p);
			}
		} else if (!result.talentLoadout || typeof result.talentLoadout !== 'object') {
			const p = document.createElement('p');
			p.textContent = 'No changes - this is your currently equipped gear!';
			parent.appendChild(p);
			dpsDeltaSpan.textContent = '';
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
	private readonly selectorModal: SelectorModal;

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
		this.selectorModal = new SelectorModal(this.simUI.rootElem, this.simUI, this.simUI.player);

		this.simUI.sim.waitForInit().then(() => {
			this.setItem(item);
			const slot = getEligibleItemSlots(this.item.item)[0];
			const eligibleEnchants = this.simUI.sim.db.getEnchants(slot);
			const openEnchantGemSelector = (event: Event) => {
				event.preventDefault();

				if (!!eligibleEnchants.length) {
					this.selectorModal.openTab(slot, SelectorModalTabs.Enchants, this.createGearData());
				} else if (!!this.item._gems.length) {
					this.selectorModal.openTab(slot, SelectorModalTabs.Gem1, this.createGearData());
				}

				const destroyItemButton = document.createElement('button');
				destroyItemButton.textContent = 'Remove from Batch';
				destroyItemButton.classList.add('btn', 'btn-danger');
				destroyItemButton.onclick = () => {
					bulkUI.setItems(
						bulkUI.getItems().filter((_, idx) => {
							return idx !== this.index;
						}),
					);
					this.selectorModal.close();
				};
				const closeX = this.selectorModal.header?.querySelector('.close-button');
				if (!!closeX) {
					this.selectorModal.header?.insertBefore(destroyItemButton, closeX);
				}
			};

			this.itemElem.iconElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.nameElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.enchantElem.addEventListener('click', openEnchantGemSelector);
		});
	}

	setItem(newItem: EquippedItem | null) {
		this.itemElem.clear();
		if (newItem != null) {
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
			equipItem: (_, equippedItem) => {
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

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<any>) {
		super(parentElem, simUI, { identifier: 'bulk-tab', title: 'Batch' });
		this.simUI = simUI;

		this.leftPanel = (<div className="bulk-tab-left tab-panel-left">{this.column1}</div>) as HTMLDivElement;
		this.rightPanel = (<div className="bulk-tab-right tab-panel-right" />) as HTMLDivElement;

		this.pendingDiv = (<div className="results-pending-overlay d-flex hide" />) as HTMLDivElement;
		this.pendingResults = new ResultsViewer(this.pendingDiv);
		this.pendingResults.hideAll();

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);
		this.contentContainer.appendChild(this.pendingDiv);

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
						this.gemIconElements[idx].src = filledId.iconUrl;
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

	addItems(items: Array<ItemSpec>) {
		if (this.items.length == 0) {
			this.items = items;
		} else {
			this.items = this.items.concat(items);
		}
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	setItems(items: Array<ItemSpec>) {
		this.items = items;
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	clearItems() {
		this.items = new Array<ItemSpec>();
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
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

		const noticeWorkInProgress = (
			<div className="bulk-items-text-line">
				<i>
					Notice: This is under very early but active development and experimental. You may also need to update your WoW AddOn if you want to import
					your bags.
				</i>
			</div>
		);
		itemsBlock.bodyElement.appendChild(noticeWorkInProgress);

		const itemTextIntro = (
			<div className="bulk-items-text-line">
				<i>
					Notice: This is under very early but active development and experimental. You may also need to update your WoW AddOn if you want to import
					your bags.
				</i>
			</div>
		);
		itemsBlock.bodyElement.appendChild(itemTextIntro);

		const itemList = (<div className="tab-panel-colbulk-gear-combo" />) as HTMLElement;
		itemsBlock.bodyElement.appendChild(itemList);

		this.itemsChangedEmitter.on(() => {
			itemList.innerHTML = '';
			if (this.items.length > 0) {
				itemTextIntro.textContent = 'The following items will be simmed together with your equipped gear.';
				for (let i = 0; i < this.items.length; ++i) {
					const spec = this.items[i];
					const item = this.simUI.sim.db.lookupItemSpec(spec);
					new BulkItemPicker(itemList, this.simUI, this, item!, i);
				}
			}
		});

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
			resultsBlock.bodyElement.innerHTML = '';

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

		const bulkSimButton = (
			<Button variant="primary" className={clsx('bulk-settings-button', 'w-100')}>
				Simulate Batch
			</Button>
		) as HTMLButtonElement;
		settingsBlock.bodyElement.appendChild(bulkSimButton);
		bulkSimButton.addEventListener('click', () => {
			this.pendingDiv.classList.remove('hide');
			this.leftPanel.classList.add('blurred');
			this.rightPanel.classList.add('blurred');
			const defaultState = cloneChildren(bulkSimButton);

			bulkSimButton.disabled = true;
			bulkSimButton.classList.add('disabled');
			bulkSimButton.replaceChildren(
				<>
					<Icon icon="spinner" className="fa-spin" /> Running
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

				if (combinations == 0) {
					combinations = progressMetrics.totalSims;
				}
				if (this.fastMode) {
					if (rounds == 0 && progressMetrics.totalSims > 0) {
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

				if (progressMetrics.finalBulkResult != null) {
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
			<Button variant="secondary" className={clsx('bulk-settings-button', 'w-100')} iconLeft="download">
				Import From Bags
			</Button>
		);
		settingsBlock.bodyElement.appendChild(importButton);
		importButton.addEventListener('click', () => new BulkGearJsonImporter(this.simUI.rootElem, this.simUI, this));

		const importFavsButton = (
			<Button variant="secondary" className={clsx('bulk-settings-button', 'w-100')}>
				Import Favorites
			</Button>
		);
		settingsBlock.bodyElement.appendChild(importFavsButton);
		importFavsButton.addEventListener('click', () => {
			const filters = this.simUI.player.sim.getFilters();
			const items = filters.favoriteItems.map(itemID => {
				return ItemSpec.create({ id: itemID });
			});
			this.addItems(items);
		});

		const searchButton = (
			<Button variant="secondary" className={clsx('bulk-settings-button', 'w-100')} iconLeft="search">
				Add Item
			</Button>
		) as HTMLButtonElement;
		settingsBlock.bodyElement.appendChild(searchButton);

		const searchText = (<input type="text" placeholder="search..." className="hide" />) as HTMLInputElement;
		settingsBlock.bodyElement.appendChild(searchText);

		const searchResults = (<ul className="batch-search-results hide" />) as HTMLUListElement;
		settingsBlock.bodyElement.appendChild(searchResults);
		let allItems = Array<UIItem>();

		searchText.addEventListener('keyup', event => {
			if (event.key == 'Enter') {
				const toAdd = Array<ItemSpec>();
				searchResults.childNodes.forEach(node => {
					const strID = (node as HTMLElement).getAttribute('data-item-id');
					if (strID != null) {
						toAdd.push(ItemSpec.create({ id: Number.parseInt(strID) }));
					}
				});
				this.addItems(toAdd);
			}
		});

		searchText.addEventListener('input', _event => {
			const searchString = searchText.value;
			searchResults.innerHTML = '';
			if (searchString.length == 0) {
				return;
			}
			const pieces = searchString.split(' ');

			let displayCount = 0;
			allItems.every(item => {
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
					const itemElement = document.createElement('li');
					itemElement.innerHTML = `<span>${item.name}</span>`;
					itemElement.setAttribute('data-item-id', item.id.toString());
					itemElement.addEventListener('click', _ev => {
						this.addItems(Array<ItemSpec>(ItemSpec.create({ id: item.id })));
					});
					if (item.heroic) {
						const htxt = document.createElement('span');
						htxt.style.color = 'green';
						htxt.innerText = '[H]';
						itemElement.appendChild(htxt);
					}
					if (item.factionRestriction == UIItem_FactionRestriction.HORDE_ONLY) {
						const ftxt = document.createElement('span');
						ftxt.style.color = 'red';
						ftxt.innerText = '(H)';
						itemElement.appendChild(ftxt);
					}
					if (item.factionRestriction == UIItem_FactionRestriction.ALLIANCE_ONLY) {
						const ftxt = document.createElement('span');
						ftxt.style.color = 'blue';
						ftxt.innerText = '(A)';
						itemElement.appendChild(ftxt);
					}
					searchResults.append(itemElement);
					displayCount++;
				}

				return displayCount < 10;
			});
		});

		const baseSearchHTML = `${convertToLegacy(<Icon icon="search" />)} Add Item`;
		searchButton.innerHTML = baseSearchHTML;
		searchButton.addEventListener('click', () => {
			if (searchText.classList.contains('hide')) {
				searchButton.innerHTML = 'Close Search Results';
				allItems = this.simUI.sim.db.getAllItems().filter(item => canEquipItem(item, this.simUI.player.getPlayerSpec(), undefined));
				searchText.classList.remove('hide');
				searchResults.classList.remove('hide');
				searchText.focus();
			} else {
				searchButton.innerHTML = baseSearchHTML;
				searchText.classList.add('hide');
				searchResults.innerHTML = '';
				searchResults.classList.add('hide');
			}
		});

		const clearButton = (
			<Button variant="secondary" className={clsx('bulk-settings-button', 'w-100')}>
				Clear All
			</Button>
		) as HTMLButtonElement;
		settingsBlock.bodyElement.appendChild(clearButton);
		clearButton.addEventListener('click', () => {
			this.clearItems();
			resultsBlock.rootElem.hidden = true;
			resultsBlock.bodyElement.innerHTML = '';
		});

		const talentsContainerRef = ref<HTMLDivElement>();
		// Talents to sim
		const talentsToSimDiv = (
			<div className={clsx('talents-picker-container', 'd-flex', !this.simTalents && 'hide')}>
				<label>Pick talents to sim (will increase time to sim)</label>
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
		const handleToggle = (frag: HTMLElement, load: TalentLoadout) => {
			const chipDiv = frag.querySelector('.saved-data-set-chip');
			const exists = this.savedTalents.some(talent => talent.name === load.name); // Replace 'id' with your unique identifier

			console.log('Exists:', exists);
			console.log('Load Object:', load);
			console.log('Saved Talents Before Update:', this.savedTalents);

			if (exists) {
				// If the object exists, find its index and remove it
				const indexToRemove = this.savedTalents.findIndex(talent => talent.name === load.name);
				this.savedTalents.splice(indexToRemove, 1);
				chipDiv?.classList.remove('active');
			} else {
				// If the object does not exist, add it
				this.savedTalents.push(load);
				chipDiv?.classList.add('active');
			}

			console.log('Updated savedTalents:', this.savedTalents);
		};
		for (const name in jsonData) {
			try {
				console.log(name, jsonData[name]);
				const savedTalentLoadout = SavedTalents.fromJson(jsonData[name]);
				const loadout = {
					talentsString: savedTalentLoadout.talentsString,
					glyphs: savedTalentLoadout.glyphs,
					name: name,
				};

				const index = this.savedTalents.findIndex(talent => JSON.stringify(talent) === JSON.stringify(loadout));
				const talentFragment = document.createElement('fragment');
				talentFragment.appendChild(
					<div className={clsx('saved-data-set-chip badge rounded-pill', index !== -1 && 'active')}>
						<Link as="button" className="saved-data-set-name">
							{name}
						</Link>
					</div>,
				);

				console.log('Adding event for loadout', loadout);
				// Wrap the event listener addition in an IIFE
				((talentFragment, loadout) => {
					talentFragment.addEventListener('click', () => handleToggle(talentFragment, loadout));
				})(talentFragment, loadout);

				talentsContainerRef.value?.appendChild(talentFragment);
			} catch (e) {
				console.log(e);
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}

		//////////////////////
		////////////////////////////////////

		// Default Gem Options
		const defaultGemDiv = (
			<div className={clsx('default-gem-container', 'd-flex', !this.autoGem && 'hide')}>
				<label>Defaults for Auto Gem</label>
				<div className="sockets-container">
					{[GemColor.GemColorRed, GemColor.GemColorYellow, GemColor.GemColorBlue, GemColor.GemColorMeta].map((socketColor, socketIndex) => {
						const iconRef = ref<HTMLImageElement>();
						const socketRef = ref<HTMLImageElement>();
						const gemContainer = (
							<div className="gem-socket-container">
								<img ref={iconRef} className="gem-icon" />
								<img ref={socketRef} className="socket-icon" />
							</div>
						);
						this.gemIconElements.push(iconRef.value!);
						socketRef.value!.src = getEmptyGemSocketIconUrl(socketColor);

						let selector: GemSelectorModal;
						const handleChoose = (itemData: ItemData<UIGem>) => {
							this.defaultGems[socketIndex] = itemData.item;
							this.storeSettings();
							ActionId.fromItemId(itemData.id)
								.fill()
								.then(filledId => {
									this.gemIconElements[socketIndex].src = filledId.iconUrl;
								});
							selector.close();
						};

						const openGemSelector = (_color: GemColor, _socketIndex: number) => {
							if (!selector) selector = new GemSelectorModal(this.simUI.rootElem, this.simUI, socketColor, handleChoose);
							selector.show();
						};

						this.gemIconElements[socketIndex].addEventListener('click', () => openGemSelector(socketColor, socketIndex));
						gemContainer.addEventListener('click', () => openGemSelector(socketColor, socketIndex));

						return gemContainer;
					})}
				</div>
			</div>
		);

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-fast-mode',
			label: 'Fast Mode',
			labelTooltip: 'Fast mode reduces accuracy but will run faster.',
			changedEvent: () => this.itemsChangedEmitter,
			getValue: () => this.fastMode,
			setValue: (_, obj, value) => {
				obj.fastMode = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-combinations',
			label: 'Combinations',
			labelTooltip:
				'When checked bulk simulator will create all possible combinations of the items. When disabled trinkets and rings will still run all combinations becausee they have two slots to fill each.',
			changedEvent: () => this.itemsChangedEmitter,
			getValue: () => this.doCombos,
			setValue: (_, obj, value) => {
				obj.doCombos = value;
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-auto-enchant',
			label: 'Auto Enchant',
			labelTooltip: 'When checked bulk simulator apply the current enchant for a slot to each replacement item it can.',
			changedEvent: () => this.itemsChangedEmitter,
			getValue: () => this.autoEnchant,
			setValue: (_, obj, value) => {
				obj.autoEnchant = value;
				defaultGemDiv.classList[!value ? 'add' : 'remove']('hide');
			},
		});
		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-auto-gem',
			label: 'Auto Gem',
			labelTooltip: 'When checked bulk simulator will fill any un-filled gem sockets with default gems.',
			changedEvent: () => this.itemsChangedEmitter,
			getValue: () => this.autoGem,
			setValue: (_, obj, value) => {
				obj.autoGem = value;
				defaultGemDiv.classList[!value ? 'add' : 'remove']('hide');
			},
		});

		new BooleanPicker<BulkTab>(settingsBlock.bodyElement, this, {
			id: 'bulk-sim-talents',
			label: 'Sim Talents',
			labelTooltip: 'When checked bulk simulator will sim chosen talent setups. Warning, it might cause the bulk sim to run for a lot longer',
			changedEvent: () => this.itemsChangedEmitter,
			getValue: () => this.simTalents,
			setValue: (_, obj, value) => {
				obj.simTalents = value;
				defaultGemDiv.classList[!value ? 'add' : 'remove']('hide');
			},
		});

		settingsBlock.bodyElement.appendChild(defaultGemDiv);
		settingsBlock.bodyElement.appendChild(talentsToSimDiv);
	}

	private setSimProgress(progress: ProgressMetrics, iterPerSecond: number, currentRound: number, rounds: number, combinations: number) {
		const secondsRemain = ((progress.totalIterations - progress.completedIterations) / iterPerSecond).toFixed();

		this.pendingResults.setContent(
			<div className="results-sim">
				<div>{combinations} total combinations.</div>
				<div>{rounds > 0 ? `${currentRound} / ${rounds} refining rounds` : ''}</div>
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

	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, socketColor: GemColor, onSelect: (itemData: ItemData<UIGem>) => void) {
		super(parent, 'selector-modal', { disposeOnClose: false });

		this.simUI = simUI;
		this.onSelect = onSelect;
		this.socketColor = socketColor;
		this.ilist = null;

		window.scrollTo({ top: 0 });

		this.header!.insertAdjacentHTML('afterbegin', `<span>Choose Default Gem</span>`);
		this.body.innerHTML = `<div class="tab-content selector-modal-tab-content"></div>`;
		this.contentElem = this.rootElem.querySelector('.selector-modal-tab-content') as HTMLElement;
	}

	show() {
		// construct item list the first time its opened.
		// This makes startup faster and also means we are sure to have item database loaded.
		if (this.ilist == null) {
			this.ilist = new ItemList<UIGem>(
				this.body,
				this.simUI,
				ItemSlot.ItemSlotHead,
				SelectorModalTabs.Gem1,
				this.simUI.player,
				'Gem1',
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
						baseEP: 0,
						ignoreEPFilter: true,
						onEquip: (_eventID, _gem: UIGem) => {
							return;
						},
					};
				}),
				this.socketColor,
				gem => {
					return this.simUI.player.computeGemEP(gem);
				},
				() => {
					return null;
				},
				() => {
					return;
				},
				this.onSelect,
			);

			// let invokeUpdate = () => {this.ilist?.updateSelected()}
			const applyFilter = () => {
				this.ilist?.applyFilters();
			};
			// Add event handlers
			// this.itemsChangedEmitter.on(invokeUpdate);

			this.addOnDisposeCallback(() => this.ilist?.dispose());

			this.simUI.sim.phaseChangeEmitter.on(applyFilter);
			this.simUI.sim.filtersChangeEmitter.on(applyFilter);
			// gearData.changeEvent.on(applyFilter);
		}

		this.open();
	}
}
