import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../../../css_utils';
import { IndividualSimUI } from '../../../individual_sim_ui';
import { ItemLevelState, ItemSpec } from '../../../proto/common';
import { UIItem, UIItem_FactionRestriction } from '../../../proto/ui';
import { ActionId } from '../../../proto_utils/action_id';
import { canEquipItem, getEligibleItemSlots } from '../../../proto_utils/utils';
import { EventID, TypedEvent } from '../../../typed_event';
import { ContentBlock } from '../../content_block';
import { createNameDescriptionLabel } from '../../gear_picker/utils';
import { NumberPicker } from '../../pickers/number_picker';
import Toast from '../../toast';
import { BulkTab } from '../bulk_tab';
import { bulkSimSlotNames, itemSlotToBulkSimItemSlot } from './utils';

const MAX_SEARCH_RESULTS = 21;

export default class BulkItemSearch extends ContentBlock {
	readonly simUI: IndividualSimUI<any>;
	readonly bulkUI: BulkTab;

	// Can be used to remove any events in addEventListener
	// https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#add_an_abortable_listener
	private abortController: AbortController;
	private signal: AbortSignal;

	private readonly inputElem: HTMLInputElement;
	private readonly cancelSearchElem: HTMLButtonElement;
	private readonly searchResultElem: HTMLElement;

	readonly filtersChangeEmitter = new TypedEvent<void>();

	private allItems: Array<UIItem> = [];

	private _searchString = '';
	private minIlvl = 0;
	private maxIlvl = 0;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, bulkUI: BulkTab) {
		super(parent, 'bulk-item-search-root', { header: { title: 'Item Search' } });

		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.abortController = new AbortController();
		this.signal = this.abortController.signal;

		const searchContainerRef = ref<HTMLDivElement>();
		const searchInputRef = ref<HTMLInputElement>();
		const cancelSearchElemRef = ref<HTMLButtonElement>();
		const searchResultsRef = ref<HTMLUListElement>();
		const ilvlFiltersContainerRef = ref<HTMLDivElement>();
		this.bodyElement.appendChild(
			<div className="bulk-gear-search-container" ref={searchContainerRef}>
				<div className="d-flex flex-column">
					<label className="form-label" htmlFor="bulkGearSearch">
						Name
					</label>
					<div className="input-group">
						<input id="bulkGearSearch" className="form-control" type="text" placeholder="Search..." ref={searchInputRef} />
						<button className="btn btn-link cancel-bulk-gear-search-btn hide" ref={cancelSearchElemRef} type="button">
							<i className="fas fa-times" />
						</button>
					</div>
					<ul ref={searchResultsRef} className="bulk-gear-search-results dropdown-menu no-hover" />
				</div>
				<div className="bulk-gear-search-ilvl-filters" ref={ilvlFiltersContainerRef} />
			</div>,
		);

		this.inputElem = searchInputRef.value!;
		this.cancelSearchElem = cancelSearchElemRef.value!;
		this.searchResultElem = searchResultsRef.value!;

		tippy(this.cancelSearchElem, { content: 'Clear search' });
		this.cancelSearchElem.addEventListener('click', () => {
			this.searchString = '';
		});

		new NumberPicker(ilvlFiltersContainerRef.value!, this, {
			id: 'bulkGearSearchMinIlvl',
			label: 'Min ILvl',
			showZeroes: false,
			changedEvent: _ => this.filtersChangeEmitter,
			getValue: _ => this.minIlvl,
			setValue: (eventID: EventID, _, newValue: number) => {
				this.minIlvl = newValue;
				this.filtersChangeEmitter.emit(eventID);
			},
		});

		ilvlFiltersContainerRef.value!.appendChild(<span className="ilvl-filters-separator">-</span>);

		new NumberPicker(ilvlFiltersContainerRef.value!, this, {
			id: 'bulkGearSearchMaxIlvl',
			label: 'Max ILvl',
			showZeroes: false,
			changedEvent: _ => this.filtersChangeEmitter,
			getValue: _ => this.maxIlvl,
			setValue: (eventID: EventID, _, newValue: number) => {
				this.maxIlvl = newValue;
				this.filtersChangeEmitter.emit(eventID);
			},
		});

		this.simUI.sim.waitForInit().then(() => {
			this.allItems = this.simUI.sim.db
				.getAllItems()
				.filter(item => canEquipItem(item, this.simUI.player.getPlayerSpec(), undefined))
				.sort((a, b) => {
					const aIlvl = a.scalingOptions?.[ItemLevelState.Base].ilvl || a.ilvl;
					const bIlvl = b.scalingOptions?.[ItemLevelState.Base].ilvl || b.ilvl;
					if (aIlvl < bIlvl) return 1;
					else if (bIlvl < aIlvl) return -1;
					else return 0;
				});

			this.filtersChangeEmitter.on(() => {
				this.performSearch();
			});

			searchInputRef.value!.addEventListener('input', () => {
				this.searchString = searchInputRef.value!.value;
			});
		});
	}

	private get searchString(): string {
		return this._searchString;
	}

	private set searchString(newString: string) {
		this._searchString = newString;
		this.inputElem.value = newString;
		this.filtersChangeEmitter.emit(TypedEvent.nextEventID());
	}

	private performSearch() {
		if (!this.searchString.length) {
			this.cancelSearchElem.classList.add('hide');
			this.searchResultElem.classList.remove('show');
			return;
		}

		const pieces = this.searchString.split(' ');
		const items = <></>;

		let matchCount = 0;

		this.allItems.forEach(item => {
			const ilvl = item.scalingOptions?.[ItemLevelState.Base].ilvl || item.ilvl;
			if (this.maxIlvl != 0 && this.maxIlvl < ilvl) return false;
			if (this.minIlvl != 0 && this.minIlvl > ilvl) return false;

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

			if (matched) matchCount++;
			if (matched && matchCount <= MAX_SEARCH_RESULTS) {
				const itemRef = ref<HTMLAnchorElement>();
				const iconRef = ref<HTMLDivElement>();
				const itemNameRef = ref<HTMLSpanElement>();
				items.appendChild(
					<li>
						<a className="dropdown-item bulk-item-search-item" dataset={{ itemId: item.id.toString() }} ref={itemRef} target="_blank">
							<div className="bulk-item-search-item-icon-wrapper">
								<span className="item-picker-ilvl">{ilvl}</span>
								<div className="bulk-item-search-item-icon" ref={iconRef} />
							</div>
							<div className="d-flex flex-column ps-2">
								<div className="d-flex">
									<span ref={itemNameRef}>{item.name}</span>
									{item.nameDescription && createNameDescriptionLabel(item.nameDescription)}
									{item.factionRestriction === UIItem_FactionRestriction.HORDE_ONLY && <span className="faction-horde">(H)</span>}
									{item.factionRestriction === UIItem_FactionRestriction.ALLIANCE_ONLY && <span className="faction-alliance">(A)</span>}
								</div>
								<small>{bulkSimSlotNames.get(itemSlotToBulkSimItemSlot.get(getEligibleItemSlots(item)[0])!)}</small>
							</div>
						</a>
					</li>,
				);

				ActionId.fromItem(item)
					.fill()
					.then(id => {
						id.setBackground(iconRef.value!);
						id.setWowheadHref(itemRef.value!);
					});
				setItemQualityCssClass(itemNameRef.value!, item.quality);

				itemRef.value!.addEventListener(
					'click',
					event => {
						event.preventDefault();
						this.bulkUI.addItem(ItemSpec.create({ id: item.id }));

						new Toast({
							delay: 1000,
							variant: 'success',
							body: (
								<>
									<strong>{item.name}</strong> was added to the batch.
								</>
							),
						});
					},
					{ signal: this.signal },
				);
			}
		});

		this.searchResultElem.replaceChildren(
			<>
				{items}
				{matchCount > MAX_SEARCH_RESULTS && (
					<li className="bulk-item-search-item bulk-item-search-results-note">
						Showing {MAX_SEARCH_RESULTS} of {matchCount} total results.
					</li>
				)}
				{matchCount === 0 && <li className="bulk-item-search-item bulk-item-search-results-note">No results found.</li>}
			</>,
		);

		this.searchResultElem.classList.add('show');
		this.cancelSearchElem.classList.remove('hide');
	}
}
