import { Tab } from 'bootstrap';
import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { REPO_RELEASES_URL } from '../../constants/other';
import { IndividualSimUI } from '../../individual_sim_ui';
import { BulkSettings, ErrorOutcomeType, ProgressMetrics, TalentLoadout } from '../../proto/api';
import { GemColor, ItemRandomSuffix, ItemSlot, ItemSpec, ReforgeStat, Spec } from '../../proto/common';
import { ItemEffectRandPropPoints, SimDatabase, SimEnchant, SimGem, SimItem } from '../../proto/db';
import { SavedTalents, UIEnchant, UIGem, UIItem } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { getEmptyGemSocketIconUrl } from '../../proto_utils/gems';
import { canEquipItem, getEligibleItemSlots, isSecondaryItemSlot } from '../../proto_utils/utils';
import { RequestTypes } from '../../sim_signal_manager';
import { TypedEvent } from '../../typed_event';
import { getEnumValues, isExternal } from '../../utils';
import { WorkerProgressCallback } from '../../worker_pool';
import { ItemData } from '../gear_picker/item_list';
import SelectorModal from '../gear_picker/selector_modal';
import { BooleanPicker } from '../pickers/boolean_picker';
import { ResultsViewer } from '../results_viewer';
import { SimTab } from '../sim_tab';
import Toast from '../toast';
import BulkItemPickerGroup from './bulk/bulk_item_picker_group';
import BulkItemSearch from './bulk/bulk_item_search';
import BulkSimResultRenderer from './bulk/bulk_sim_results_renderer';
import GemSelectorModal from './bulk/gem_selector_modal';
import { BulkSimItemSlot, getBulkItemSlotFromSlot } from './bulk/utils';
import { BulkGearJsonImporter } from './importers';

const WEB_DEFAULT_ITERATIONS = 1000;
const WEB_ITERATIONS_LIMIT = 50_000;
const LOCAL_ITERATIONS_LIMIT = 1_000_000;

export class BulkTab extends SimTab {
	readonly simUI: IndividualSimUI<any>;
	readonly playerCanDualWield: boolean;
	readonly playerIsFuryWarrior: boolean;

	readonly itemsChangedEmitter = new TypedEvent<void>();
	readonly settingsChangedEmitter = new TypedEvent<void>();

	private readonly setupTabElem: HTMLElement;
	private readonly resultsTabElem: HTMLElement;
	private readonly combinationsElem: HTMLElement;
	private readonly bulkSimButton: HTMLButtonElement;
	private readonly settingsContainer: HTMLElement;
	private readonly booleanSettingsContainer: HTMLElement;

	private pendingDiv: HTMLDivElement;

	private setupTab: Tab;
	private resultsTab: Tab;
	private pendingResults: ResultsViewer;

	readonly selectorModal: SelectorModal;

	// The main array we will use to store items with indexes. Null values are the result of removed items to avoid having to shift pickers over and over.
	protected items: Array<ItemSpec | null> = new Array<ItemSpec | null>();
	protected pickerGroups: Map<BulkSimItemSlot, BulkItemPickerGroup> = new Map();

	protected combinations = 0;
	protected iterations = 0;

	// TODO: Make a real options probably
	doCombos: boolean;
	fastMode: boolean;
	autoGem: boolean;
	simTalents: boolean;
	autoEnchant: boolean;
	defaultGems: SimGem[];
	savedTalents: TalentLoadout[];
	gemIconElements: HTMLImageElement[];

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<any>) {
		super(parentElem, simUI, { identifier: 'bulk-tab', title: 'Batch (<span class="text-success">New</span>)' });

		this.simUI = simUI;
		this.playerCanDualWield = this.simUI.player.getPlayerSpec().canDualWield;
		this.playerIsFuryWarrior = this.simUI.player.getSpec() === Spec.SpecFuryWarrior;

		const setupTabBtnRef = ref<HTMLButtonElement>();
		const setupTabRef = ref<HTMLDivElement>();
		const resultsTabBtnRef = ref<HTMLButtonElement>();
		const resultsTabRef = ref<HTMLDivElement>();
		const settingsContainerRef = ref<HTMLDivElement>();
		const combinationsElemRef = ref<HTMLHeadingElement>();
		const bulkSimBtnRef = ref<HTMLButtonElement>();
		const booleanSettingsContainerRef = ref<HTMLDivElement>();

		this.contentContainer.appendChild(
			<>
				<div className="bulk-tab-left tab-panel-left">
					<div className="bulk-tab-tabs">
						<ul className="nav nav-tabs" attributes={{ role: 'tablist' }}>
							<li className="nav-item" attributes={{ role: 'presentation' }}>
								<button
									className="nav-link active"
									type="button"
									attributes={{
										role: 'tab',
										// @ts-expect-error
										'aria-controls': 'bulkSetupTab',
										'aria-selected': true,
									}}
									dataset={{
										bsToggle: 'tab',
										bsTarget: `#bulkSetupTab`,
									}}
									ref={setupTabBtnRef}>
									Setup
								</button>
							</li>
							<li className="nav-item" attributes={{ role: 'presentation' }}>
								<button
									className="nav-link"
									type="button"
									attributes={{
										role: 'tab',
										// @ts-expect-error
										'aria-controls': 'bulkResultsTab',
										'aria-selected': false,
									}}
									dataset={{
										bsToggle: 'tab',
										bsTarget: `#bulkResultsTab`,
									}}
									ref={resultsTabBtnRef}>
									Results
								</button>
							</li>
						</ul>
						<div className="tab-content">
							<div id="bulkSetupTab" className="tab-pane fade active show" ref={setupTabRef} />
							<div id="bulkResultsTab" className="tab-pane fade show" ref={resultsTabRef}>
								<div className="d-flex align-items-center justify-content-center p-gap">Run a simulation to view results</div>
							</div>
						</div>
					</div>
				</div>
				<div className="bulk-tab-right tab-panel-right">
					<div className="bulk-settings-outer-container">
						<div className="bulk-settings-container" ref={settingsContainerRef}>
							<div className="bulk-combinations-count h4" ref={combinationsElemRef} />
							<button className="btn btn-primary bulk-settings-btn" ref={bulkSimBtnRef}>
								Simulate Batch
							</button>
							<div className="bulk-boolean-settings-container" ref={booleanSettingsContainerRef}></div>
						</div>
					</div>
					,
				</div>
			</>,
		);

		this.setupTabElem = setupTabRef.value!;
		this.resultsTabElem = resultsTabRef.value!;
		this.pendingDiv = (<div className="results-pending-overlay" />) as HTMLDivElement;

		this.combinationsElem = combinationsElemRef.value!;
		this.bulkSimButton = bulkSimBtnRef.value!;
		this.settingsContainer = settingsContainerRef.value!;
		this.booleanSettingsContainer = booleanSettingsContainerRef.value!;

		this.setupTab = new Tab(setupTabBtnRef.value!);
		this.resultsTab = new Tab(resultsTabBtnRef.value!);

		this.pendingResults = new ResultsViewer(this.pendingDiv);
		this.pendingResults.hideAll();
		this.selectorModal = new SelectorModal(this.simUI.rootElem, this.simUI, this.simUI.player, undefined, {
			id: 'bulk-selector-modal',
		});

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

			const loadEquippedItems = () => {
				// Clear all previously equipped items from the pickers
				for (const group of this.pickerGroups.values()) {
					if (group.has(-1)) {
						group.remove(-1);
					}
					if (group.has(-2)) {
						group.remove(-2);
					}
				}

				this.simUI.player.getEquippedItems().forEach((equippedItem, slot) => {
					const bulkSlot = getBulkItemSlotFromSlot(slot, this.playerCanDualWield);
					const group = this.pickerGroups.get(bulkSlot)!;
					const idx = this.isSecondaryItemSlot(slot) ? -2 : -1;
					if (equippedItem) {
						group.add(idx, equippedItem);
					}
				});
			};
			loadEquippedItems();

			TypedEvent.onAny([this.simUI.player.challengeModeChangeEmitter, this.simUI.player.gearChangeEmitter]).on(() => loadEquippedItems());
			this.itemsChangedEmitter.on(() => this.storeSettings());

			TypedEvent.onAny([this.itemsChangedEmitter, this.settingsChangedEmitter, this.simUI.sim.iterationsChangeEmitter]).on(() => {
				this.getCombinationsCount().then(result => this.combinationsElem.replaceChildren(result));
			});
			this.getCombinationsCount().then(result => this.combinationsElem.replaceChildren(result));
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

			this.addItems(settings.items);

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
			items: this.getItems(),
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
			iterationsPerCombo: this.getDefaultIterationsCount(),
		});
	}

	private getDefaultIterationsCount(): number {
		if (isExternal()) return WEB_DEFAULT_ITERATIONS;

		return this.simUI.sim.getIterations();
	}

	protected createBulkItemsDatabase(): SimDatabase {
		const itemsDb = SimDatabase.create();
		for (const is of this.items.values()) {
			if (!is) continue;

			const item = this.simUI.sim.db.lookupItemSpec(is);
			if (!item) {
				throw new Error(`item with ID ${is.id} not found in database`);
			}
			itemsDb.items.push(SimItem.fromJson(UIItem.toJson(item.item), { ignoreUnknownFields: true }));

			const ieRpp = this.simUI.sim.db.getItemEffectRandPropPoints(item.ilvl);
			if (ieRpp) {
				itemsDb.itemEffectRandPropPoints.push(ItemEffectRandPropPoints.create(this.simUI.sim.db.getItemEffectRandPropPoints(item.ilvl)));
			}

			if (item.enchant) {
				itemsDb.enchants.push(
					SimEnchant.fromJson(UIEnchant.toJson(item.enchant), {
						ignoreUnknownFields: true,
					}),
				);
			}
			if (item.randomSuffix) {
				itemsDb.randomSuffixes.push(
					ItemRandomSuffix.fromJson(ItemRandomSuffix.toJson(item.randomSuffix), {
						ignoreUnknownFields: true,
					}),
				);
			}
			if (item.reforge) {
				itemsDb.reforgeStats.push(
					ReforgeStat.fromJson(ReforgeStat.toJson(item.reforge), {
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

	// Add an item to its eligible bulk sim item slot(s). Mainly used for importing and search
	addItem(item: ItemSpec) {
		this.addItems([item]);
	}
	// Add items to their eligible bulk sim item slot(s). Mainly used for importing and search
	addItems(items: ItemSpec[]) {
		items.forEach(item => {
			const equippedItem = this.simUI.sim.db.lookupItemSpec(item)?.withChallengeMode(this.simUI.player.getChallengeModeEnabled()).withDynamicStats();
			if (!!equippedItem) {
				getEligibleItemSlots(equippedItem.item, this.playerIsFuryWarrior).forEach(slot => {
					// Avoid duplicating rings/trinkets/weapons
					if (this.isSecondaryItemSlot(slot)) return;

					const idx = this.items.push(item) - 1;
					const bulkSlot = getBulkItemSlotFromSlot(slot, this.playerCanDualWield);
					const group = this.pickerGroups.get(bulkSlot)!;
					group.add(idx, equippedItem);
				});
			}
		});

		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}
	// Add an item to a particular bulk sim item slot
	addItemToSlot(item: ItemSpec, bulkSlot: BulkSimItemSlot) {
		const equippedItem = this.simUI.sim.db.lookupItemSpec(item)?.withChallengeMode(this.simUI.player.getChallengeModeEnabled()).withDynamicStats();
		if (!!equippedItem) {
			const eligibleItemSlots = getEligibleItemSlots(equippedItem.item, this.playerIsFuryWarrior);
			if (!canEquipItem(equippedItem.item, this.simUI.player.getPlayerSpec(), eligibleItemSlots[0])) return;

			const idx = this.items.push(item) - 1;
			const group = this.pickerGroups.get(bulkSlot)!;
			group.add(idx, equippedItem);

			this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
		}
	}

	updateItem(idx: number, newItem: ItemSpec) {
		const equippedItem = this.simUI.sim.db.lookupItemSpec(newItem)?.withChallengeMode(this.simUI.player.getChallengeModeEnabled()).withDynamicStats();
		if (!!equippedItem) {
			this.items[idx] = newItem;

			getEligibleItemSlots(equippedItem.item, this.playerIsFuryWarrior).forEach(slot => {
				// Avoid duplicating rings/trinkets/weapons
				if (this.isSecondaryItemSlot(slot)) return;

				const bulkSlot = getBulkItemSlotFromSlot(slot, this.playerCanDualWield);
				const group = this.pickerGroups.get(bulkSlot)!;
				group.update(idx, equippedItem);
			});
		}

		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	removeItem(item: ItemSpec) {
		this.items.forEach((savedItem, idx) => {
			if (!!savedItem && savedItem.id === item.id) {
				this.removeItemByIndex(idx);
			}
		});
	}
	removeItemByIndex(idx: number) {
		if (idx < 0 || this.items.length < idx || !this.items[idx]) {
			new Toast({
				variant: 'error',
				body: 'Failed to remove item, please report this issue.',
			});
			return;
		}

		const item = this.items[idx]!;
		const equippedItem = this.simUI.sim.db.lookupItemSpec(item);
		if (!!equippedItem) {
			this.items[idx] = null;

			// Try to find the matching item within its eligible groups
			getEligibleItemSlots(equippedItem.item, this.playerIsFuryWarrior).forEach(slot => {
				const bulkSlot = getBulkItemSlotFromSlot(slot, this.playerCanDualWield);
				const group = this.pickerGroups.get(bulkSlot)!;
				group.remove(idx);
			});
			this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
		}
	}

	clearItems() {
		for (let idx = 0; idx < this.items.length; idx++) {
			this.removeItemByIndex(idx);
		}
		this.items = new Array<ItemSpec>();
		this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	hasItem(item: ItemSpec) {
		return this.items.some(i => !!i && ItemSpec.equals(i, item));
	}

	getItems(): Array<ItemSpec> {
		const result = new Array<ItemSpec>();
		this.items.forEach(spec => {
			if (!spec) return;

			result.push(ItemSpec.clone(spec));
		});
		return result;
	}

	private setCombinations(doCombos: boolean) {
		this.doCombos = doCombos;
		this.settingsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	private setFastMode(fastMode: boolean) {
		this.fastMode = fastMode;
		this.settingsChangedEmitter.emit(TypedEvent.nextEventID());
	}

	protected async runBulkSim(onProgress: WorkerProgressCallback) {
		this.pendingResults.setPending();

		try {
			const result = await this.simUI.sim.runBulkSim(this.createBulkSettings(), this.createBulkItemsDatabase(), onProgress);
			if (result.error?.type == ErrorOutcomeType.ErrorOutcomeAborted) {
				new Toast({
					variant: 'info',
					body: 'Bulk sim cancelled.',
				});
			}
		} catch (e) {
			this.isPending = false;
			this.simUI.handleCrash(e);
		}
	}

	protected async calculateBulkCombinations() {
		try {
			const result = await this.simUI.sim.calculateBulkCombinations(this.createBulkSettings(), this.createBulkItemsDatabase());
			this.combinations = result?.numCombinations ?? 0;
			this.iterations = result?.numIterations ?? 0;
		} catch (e) {
			this.simUI.handleCrash(e);
		}
	}

	protected buildTabContent() {
		this.buildSetupTabContent();
		this.buildResultsTabContent();
		this.buildBatchSettings();
	}

	private buildSetupTabContent() {
		const bagImportBtnRef = ref<HTMLButtonElement>();
		const favsImportBtnRef = ref<HTMLButtonElement>();
		const clearBtnRef = ref<HTMLButtonElement>();
		this.setupTabElem.appendChild(
			<>
				{/* // TODO: Remove once we're more comfortable with the state of Batch sim */}
				<p className="mb-0">
					<span className="bold">Batch Simming</span> is a new feature akin to the <span className="bold">Top Gear</span> sim on{' '}
					<a href="https://raidbots.com" target="_blank">
						Raidbots.com
					</a>{' '}
					that allows you to select multiple items and sim them find the best combinations.
					<br />
					This is an <span className="text-brand">Alpha</span> feature, so if you have feedback or find a bug, please report it!
				</p>
				{isExternal() && (
					<p className="mb-0">
						<a href={REPO_RELEASES_URL} target="_blank">
							<i className="fas fa-gauge-high me-1" />
							Download the local sim for faster results.
						</a>
					</p>
				)}
				<div className="bulk-gear-actions">
					<button className="btn btn-secondary" ref={bagImportBtnRef}>
						<i className="fa fa-download me-1" /> Import From Bags
					</button>
					<button className="btn btn-secondary" ref={favsImportBtnRef}>
						<i className="fa fa-download me-1" /> Import Favorites
					</button>
					<button className="btn btn-danger ms-auto" ref={clearBtnRef}>
						<i className="fas fa-times me-1" />
						Clear Items
					</button>
				</div>
			</>,
		);

		const bagImportButton = bagImportBtnRef.value!;
		const favsImportButton = favsImportBtnRef.value!;
		const clearButton = clearBtnRef.value!;

		bagImportButton.addEventListener('click', () => new BulkGearJsonImporter(this.simUI.rootElem, this.simUI, this).open());

		favsImportButton.addEventListener('click', () => {
			const filters = this.simUI.player.sim.getFilters();
			const items = filters.favoriteItems.map(itemID => ItemSpec.create({ id: itemID }));
			this.addItems(items);
		});

		clearButton.addEventListener('click', () => this.clearItems());

		new BulkItemSearch(this.setupTabElem, this.simUI, this);

		const itemList = (<div className="bulk-gear-combo" />) as HTMLElement;
		this.setupTabElem.appendChild(itemList);

		getEnumValues<BulkSimItemSlot>(BulkSimItemSlot).forEach(bulkSlot => {
			if (this.playerCanDualWield && [BulkSimItemSlot.ItemSlotMainHand, BulkSimItemSlot.ItemSlotOffHand].includes(bulkSlot)) return;
			if (!this.playerCanDualWield && bulkSlot === BulkSimItemSlot.ItemSlotHandWeapon) return;

			this.pickerGroups.set(bulkSlot, new BulkItemPickerGroup(itemList, this.simUI, this, bulkSlot));
		});
	}

	private buildResultsTabContent() {
		this.simUI.sim.bulkSimStartEmitter.on(() => this.resultsTabElem.replaceChildren());

		this.simUI.sim.bulkSimResultEmitter.on((_, bulkSimResult) => {
			for (const r of bulkSimResult.results) {
				new BulkSimResultRenderer(this.resultsTabElem, this.simUI, this, r, bulkSimResult.equippedGearResult!);
			}
			this.isPending = false;
			this.resultsTab.show();
		});
	}

	// Return whether or not the slot is considered secondary and the item should be grouped
	// This includes items in the Finger2 or Trinket2 slots, or OffHand for dual-wield specs
	private isSecondaryItemSlot(slot: ItemSlot) {
		return isSecondaryItemSlot(slot) || (this.playerCanDualWield && slot === ItemSlot.ItemSlotOffHand);
	}

	private set isPending(value: boolean) {
		if (value) {
			this.simUI.rootElem.classList.add('blurred');
			this.simUI.rootElem.insertAdjacentElement('afterend', this.pendingDiv);
		} else {
			this.simUI.rootElem.classList.remove('blurred');
			this.pendingDiv.remove();
			this.pendingResults.hideAll();
		}
	}

	protected buildBatchSettings() {
		let isRunning = false;
		this.bulkSimButton.addEventListener('click', async () => {
			if (isRunning) return;
			isRunning = true;
			this.bulkSimButton.disabled = true;
			this.isPending = true;
			let waitAbort = false;
			try {
				await this.simUI.sim.signalManager.abortType(RequestTypes.All);

				this.pendingResults.addAbortButton(async () => {
					if (waitAbort) return;
					try {
						waitAbort = true;
						await this.simUI.sim.signalManager.abortType(RequestTypes.BulkSim);
					} catch (error) {
						console.error('Error on bulk sim abort!');
						console.error(error);
					} finally {
						waitAbort = false;
						if (!isRunning) this.bulkSimButton.disabled = false;
					}
				});

				let simStart = new Date().getTime();
				let lastTotal = 0;
				let rounds = 0;
				let currentRound = 0;
				let combinations = 0;

				await this.runBulkSim((progressMetrics: ProgressMetrics) => {
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
				});
			} catch (error) {
				console.error(error);
			} finally {
				isRunning = false;
				if (!waitAbort) this.bulkSimButton.disabled = false;
				this.isPending = false;
			}
		});

		if (!isExternal()) {
			new BooleanPicker<BulkTab>(this.booleanSettingsContainer, this, {
				id: 'bulk-fast-mode',
				label: 'Fast Mode',
				labelTooltip: 'Fast mode reduces accuracy but will run faster.',
				changedEvent: _modObj => this.settingsChangedEmitter,
				getValue: _modObj => this.fastMode,
				setValue: (_, _modObj, newValue: boolean) => {
					this.setFastMode(newValue);
				},
			});
		}
		new BooleanPicker<BulkTab>(this.booleanSettingsContainer, this, {
			id: 'bulk-combinations',
			label: 'Combinations',
			labelTooltip:
				'When checked bulk simulator will create all possible combinations of the items. When disabled trinkets and rings will still run all combinations because they each have two slots to fill.',
			changedEvent: _modObj => this.settingsChangedEmitter,
			getValue: _modObj => this.doCombos,
			setValue: (_, _modObj, newValue: boolean) => {
				this.setCombinations(newValue);
			},
		});
		new BooleanPicker<BulkTab>(this.booleanSettingsContainer, this, {
			id: 'bulk-auto-enchant',
			label: 'Auto Enchant',
			labelTooltip: 'When checked bulk simulator apply the current enchant for a slot to each replacement item it can.',
			changedEvent: (_obj: BulkTab) => this.settingsChangedEmitter,
			getValue: _obj => this.autoEnchant,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.autoEnchant = value;
			},
		});

		const socketsContainerRef = ref<HTMLDivElement>();
		const defaultGemDiv = (
			<div className={clsx('default-gem-container', !this.autoGem && 'hide')}>
				<h6>Default Gems</h6>
				<div ref={socketsContainerRef} className="sockets-container"></div>
			</div>
		);

		const talentsContainerRef = ref<HTMLDivElement>();
		const talentsToSimDiv = (
			<div className={clsx('talents-picker-container', !this.simTalents && 'hide')}>
				<h6>Talents to Sim</h6>
				<div ref={talentsContainerRef} className="talents-container"></div>
			</div>
		);

		new BooleanPicker<BulkTab>(this.booleanSettingsContainer, this, {
			id: 'bulk-auto-gem',
			label: 'Auto Gem',
			labelTooltip: 'When checked bulk simulator will fill any un-filled gem sockets with default gems.',
			changedEvent: (_obj: BulkTab) => this.settingsChangedEmitter,
			getValue: _obj => this.autoGem,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.autoGem = value;
				defaultGemDiv.classList[value ? 'remove' : 'add']('hide');
			},
		});

		new BooleanPicker<BulkTab>(this.booleanSettingsContainer, this, {
			id: 'bulk-sim-talents',
			label: 'Sim Talents',
			labelTooltip: 'When checked bulk simulator will sim chosen talent setups. Warning, it might cause the bulk sim to run for a lot longer',
			changedEvent: (_obj: BulkTab) => this.settingsChangedEmitter,
			getValue: _obj => this.simTalents,
			setValue: (_, obj: BulkTab, value: boolean) => {
				obj.simTalents = value;
				talentsToSimDiv.classList[value ? 'remove' : 'add']('hide');
			},
		});

		this.settingsContainer.appendChild(
			<>
				{defaultGemDiv}
				{talentsToSimDiv}
			</>,
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

			this.settingsChangedEmitter.emit(TypedEvent.nextEventID());
		};

		const talentSelections = this.simUI.individualConfig.presets.talents.map(preset => {
			return {
				talentsString: preset.data.talentsString,
				glyphs: preset.data.glyphs,
				name: preset.name,
			};
		});

		for (const name in jsonData) {
			try {
				const savedTalentLoadout = SavedTalents.fromJson(jsonData[name]);
				talentSelections.push({
					talentsString: savedTalentLoadout.talentsString,
					glyphs: savedTalentLoadout.glyphs,
					name: name,
				});
			} catch (e) {
				console.log(e);
				console.warn('Failed parsing saved data: ' + jsonData[name]);
			}
		}

		talentSelections.forEach(selection => {
			const index = this.savedTalents.findIndex(talent => JSON.stringify(talent) === JSON.stringify(selection));
			const talentChipRef = ref<HTMLDivElement>();
			const talentButtonRef = ref<HTMLButtonElement>();

			talentsContainerRef.value!.appendChild(
				<div ref={talentChipRef} className={clsx('saved-data-set-chip badge rounded-pill', index !== -1 && 'active')}>
					<button ref={talentButtonRef} className="saved-data-set-name">
						{selection.name}
					</button>
				</div>,
			);
			talentButtonRef.value!.addEventListener('click', () => handleToggle(talentChipRef.value!, selection));
		});
	}

	private async getCombinationsCount(): Promise<Element> {
		await this.calculateBulkCombinations();

		const warningRef = ref<HTMLButtonElement>();
		const rtn = (
			<>
				<span className={clsx(this.showIterationsWarning() && 'text-danger')}>
					{this.combinations === 1 ? '1 Combination' : `${this.combinations} Combinations`}
					<br />
					<small>{this.iterations} Iterations</small>
				</span>
				{this.showIterationsWarning() && (
					<button className="warning link-warning" ref={warningRef}>
						<i className="fas fa-exclamation-triangle fa-2x" />
					</button>
				)}
			</>
		);

		if (!!warningRef.value) {
			tippy(warningRef.value, {
				content: `Simming over ${WEB_ITERATIONS_LIMIT} iterations in the web sim will take a long time or may not run at all.
					For larger batch sims we recommend using Fast Mode or downloading the executable and running on your machine.`,
				placement: 'left',
				popperOptions: {
					modifiers: [
						{
							name: 'flip',
							options: {
								fallbackPlacements: ['auto'],
							},
						},
					],
				},
			});
		}

		return rtn;
	}

	private showIterationsWarning(): boolean {
		return this.iterations > this.getIterationsLimit();
	}

	private getIterationsLimit(): number {
		return isExternal() ? WEB_ITERATIONS_LIMIT : LOCAL_ITERATIONS_LIMIT;
	}

	private setSimProgress(progress: ProgressMetrics, iterPerSecond: number, currentRound: number, rounds: number, combinations: number) {
		const secondsRemaining = ((progress.totalIterations - progress.completedIterations) / iterPerSecond).toFixed();
		if (isNaN(Number(secondsRemaining))) return;

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
				<div>{secondsRemaining} seconds remaining.</div>
			</div>,
		);
	}
}
