import { ref } from 'tsx-vanilla';

import { BaseModal } from '../components/base_modal.js';
import { Component } from '../components/component.js';
import { ContentBlock } from '../components/content_block.js';
import { Input } from '../components/input.js';
import { setItemQualityCssClass } from '../css_utils.js';
import { Player } from '../player.js';
import { Glyphs, ItemQuality } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { Database } from '../proto_utils/database.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { stringComparator } from '../utils.js';

export type GlyphConfig = {
	name: string;
	description: string;
	iconUrl: string;
};

export type GlyphsConfig = {
	majorGlyphs: Record<number, GlyphConfig>;
	minorGlyphs: Record<number, GlyphConfig>;
};

interface GlyphData {
	id: number;
	name: string;
	description: string;
	iconUrl: string;
	quality: ItemQuality | null;
	spellId: number;
}

const emptyGlyphData: GlyphData = {
	id: 0,
	name: 'Empty',
	description: '',
	iconUrl: 'https://wow.zamimg.com/images/wow/icons/medium/inventoryslot_empty.jpg',
	quality: null,
	spellId: 0,
};

export class GlyphsPicker extends Component {
	private readonly glyphsConfig: GlyphsConfig;
	readonly selectorModal: GlyphSelectorModal;
	readonly player: Player<any>;
	majorGlyphPickers: Array<GlyphPicker> = [];
	minorGlyphPickers: Array<GlyphPicker> = [];

	constructor(parent: HTMLElement, player: Player<any>, glyphsConfig: GlyphsConfig) {
		super(parent, 'glyphs-picker-root');
		this.glyphsConfig = glyphsConfig;
		this.player = player;

		const majorGlyphs = Object.keys(glyphsConfig.majorGlyphs).map(idStr => Number(idStr));
		const minorGlyphs = Object.keys(glyphsConfig.minorGlyphs).map(idStr => Number(idStr));

		const majorGlyphsBlock = new ContentBlock(this.rootElem, 'major-glyphs', {
			header: { title: 'Major Glyphs', extraCssClasses: ['border-0'] },
		});

		const minorGlyphsBlock = new ContentBlock(this.rootElem, 'minor-glyphs', {
			header: { title: 'Minor Glyphs', extraCssClasses: ['border-0'] },
		});
		this.selectorModal = new GlyphSelectorModal(this.rootElem.closest('.individual-sim-ui')!);

		Database.get().then(db => {
			const majorGlyphsData = majorGlyphs.map(glyph => this.getGlyphData(glyph, db));
			const minorGlyphsData = minorGlyphs.map(glyph => this.getGlyphData(glyph, db));

			majorGlyphsData.sort((a, b) => stringComparator(a.name, b.name));
			minorGlyphsData.sort((a, b) => stringComparator(a.name, b.name));

			this.majorGlyphPickers = (['major1', 'major2', 'major3'] as Array<keyof Glyphs>).map(
				glyphField =>
					new GlyphPicker(majorGlyphsBlock.bodyElement, {
						label: 'Major',
						player,
						selectorModal: this.selectorModal,
						glyphOptions: majorGlyphsData,
						glyphField,
					}),
			);

			this.minorGlyphPickers = (['minor1', 'minor2', 'minor3'] as Array<keyof Glyphs>).map(
				glyphField =>
					new GlyphPicker(minorGlyphsBlock.bodyElement, {
						label: 'Minor',
						player,
						selectorModal: this.selectorModal,
						glyphOptions: minorGlyphsData,
						glyphField,
					}),
			);
		});
	}

	// In case we ever want to parse description from tooltip HTML.
	//static descriptionRegex = /<a href=\\"\/wotlk.*>(.*)<\/a>/g;
	getGlyphData(glyph: number, db: Database): GlyphData {
		const glyphConfig = this.glyphsConfig.majorGlyphs[glyph] || this.glyphsConfig.minorGlyphs[glyph];

		return {
			id: glyph,
			name: glyphConfig.name,
			description: glyphConfig.description,
			iconUrl: glyphConfig.iconUrl,
			quality: ItemQuality.ItemQualityCommon,
			spellId: db.glyphItemToSpellId(glyph),
		};
	}
}

type GlyphPickerConfig = {
	label: string;
	player: Player<any>;
	glyphOptions: GlyphData[];
	glyphField: keyof Glyphs;
	selectorModal: GlyphSelectorModal;
};

class GlyphPicker extends Input<Player<any>, number> {
	readonly player: Player<any>;

	selectedGlyph: GlyphData | undefined;

	private readonly glyphOptions: GlyphData[];

	private readonly anchorElem: HTMLAnchorElement;
	private readonly iconElem: HTMLImageElement;
	private readonly nameElem: HTMLSpanElement;

	constructor(parent: HTMLElement, { player, selectorModal, glyphOptions, glyphField }: GlyphPickerConfig) {
		super(parent, 'glyph-picker-root', player, {
			id: `glyph-picker-glyph-${glyphField}`,
			inline: true,
			changedEvent: (player: Player<any>) => player.glyphsChangeEmitter,
			getValue: (player: Player<any>) => player.getGlyphs()[glyphField] as number,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const glyphs = player.getGlyphs();
				(glyphs[glyphField] as number) = newValue;
				player.setGlyphs(eventID, glyphs);
			},
		});
		this.rootElem.classList.add('item-picker-root');

		this.player = player;
		this.glyphOptions = glyphOptions;
		this.selectedGlyph = emptyGlyphData;

		const anchorElemRef = ref<HTMLAnchorElement>();
		const iconElemRef = ref<HTMLImageElement>();
		const nameElemRef = ref<HTMLSpanElement>();

		this.rootElem.appendChild(
			<a ref={anchorElemRef} attributes={{ role: 'button' }} className="glyph-link">
				<img ref={iconElemRef} className="item-picker-icon" />
				<div className="item-picker-labels-container">
					<span ref={nameElemRef} className="item-picker-name-container" />
				</div>
			</a>,
		);

		this.anchorElem = anchorElemRef.value!;
		this.iconElem = iconElemRef.value!;
		this.nameElem = nameElemRef.value!;

		this.anchorElem.addEventListener(
			'click',
			event => {
				event.preventDefault();
				selectorModal.openTab(this, glyphOptions);
			},
			{ signal: this.signal },
		);

		this.init();
	}

	getInputElem(): HTMLElement {
		return this.iconElem;
	}

	getInputValue(): number {
		return this.selectedGlyph?.id ?? 0;
	}

	setInputValue(newValue: number) {
		this.selectedGlyph = this.glyphOptions.find(glyphData => glyphData.id == newValue);

		if (this.selectedGlyph) {
			if (this.selectedGlyph.spellId) {
				this.anchorElem.href = ActionId.makeSpellUrl(this.selectedGlyph.spellId);
				ActionId.makeSpellTooltipData(this.selectedGlyph.spellId).then(url => {
					this.anchorElem.dataset.wowhead = url;
					this.anchorElem.dataset.whtticon = 'false';
				});
			} else {
				this.anchorElem.href = ActionId.makeItemUrl(this.selectedGlyph.id);
				ActionId.makeItemTooltipData(this.selectedGlyph.id).then(url => {
					this.anchorElem.dataset.wowhead = url;
					this.anchorElem.dataset.whtticon = 'false';
				});
			}

			this.iconElem.src = this.selectedGlyph.iconUrl;
			this.nameElem.textContent = this.selectedGlyph.name.replace(/Glyph of /, '');
		} else {
			this.clear();
		}
	}

	private clear() {
		this.anchorElem.removeAttribute('data-wowhead');
		this.anchorElem.removeAttribute('href');

		this.iconElem.src = emptyGlyphData.iconUrl;
		this.nameElem.textContent = emptyGlyphData.name;
	}
}

class GlyphSelectorModal extends BaseModal {
	list: HTMLUListElement;
	listItems: HTMLLIElement[] = [];
	search: HTMLInputElement;
	glyphOptions: GlyphData[] = [];
	glyphPicker: GlyphPicker | null = null;
	constructor(parent: HTMLElement) {
		super(parent, 'glyph-modal', { title: 'Glyphs', disposeOnClose: false });

		const list = ref<HTMLUListElement>();
		const search = ref<HTMLInputElement>();

		this.body.appendChild(
			<>
				<div className="input-root">
					<input ref={search} className="selector-modal-search form-control" type="text" placeholder="Search..." />
				</div>
				<ul ref={list} className="selector-modal-list"></ul>
			</>,
		);

		this.list = list.value!;
		this.search = search.value!;

		this.search.addEventListener('input', () => this.applyFilters());
	}

	openTab(glyphPicker: GlyphPicker, glyphOptions: GlyphData[]) {
		this.setData(glyphPicker, glyphOptions);
		this.applyFilters();
		this.open();
	}

	private setData(glyphPicker: GlyphPicker, glyphOptions: GlyphData[]) {
		this.glyphPicker = glyphPicker;
		this.list.innerHTML = '';
		this.listItems = [];
		this.glyphOptions = [emptyGlyphData, ...glyphOptions];

		const listItemElems = this.glyphOptions.map((glyphData, glyphIdx) => {
			const anchorElem = ref<HTMLAnchorElement>();
			const iconElem = ref<HTMLImageElement>();
			const nameElem = ref<HTMLLabelElement>();

			const listItemElem = (
				<li
					className="selector-modal-list-item"
					dataset={{
						idx: String(glyphIdx),
					}}>
					<a ref={anchorElem} className="selector-modal-list-item-link">
						<img ref={iconElem} className="selector-modal-list-item-icon" />
						<label ref={nameElem} className="selector-modal-list-item-name">
							{glyphData.name}
						</label>
						<span className="selector-modal-list-item-description">{glyphData.description}</span>
					</a>
				</li>
			);

			if (anchorElem.value) {
				if (glyphData.spellId) {
					anchorElem.value.href = ActionId.makeSpellUrl(glyphData.spellId);
				} else {
					anchorElem.value.href = ActionId.makeItemUrl(glyphData.id);
				}
				anchorElem.value.addEventListener('click', event => {
					event.preventDefault();
					this.glyphPicker?.setValue(TypedEvent.nextEventID(), glyphData.id);
				});
			}
			if (iconElem.value) {
				iconElem.value.src = glyphData.iconUrl;
			}
			if (nameElem.value) setItemQualityCssClass(nameElem.value, glyphData.quality);

			return listItemElem as HTMLLIElement;
		});

		this.listItems = listItemElems;
		this.list.appendChild(<>{this.listItems}</>);

		this.glyphPicker.player.glyphsChangeEmitter.on(() => {
			this.applyFilters();
		});
	}

	applyFilters() {
		if (!this.glyphPicker) return;
		const selectedGlyphId = this.glyphPicker.selectedGlyph?.id ?? 0;

		this.listItems.forEach(elem => {
			const listItemIdx = parseInt(elem.dataset.idx!);
			const listItemData = this.glyphOptions[listItemIdx];
			elem.classList[listItemData.id == selectedGlyphId ? 'add' : 'remove']('active');
		});

		this.listItems.map(elem => {
			const listItemIdx = parseInt(elem.dataset.idx!);
			const listItemData = this.glyphOptions[listItemIdx];
			let action: 'add' | 'remove' = 'remove';

			if (this.search.value.length > 0) {
				const searchQuery = this.search.value.toLowerCase().split(' ');
				const name = listItemData.name.toLowerCase();

				searchQuery.forEach(v => {
					if (!name.includes(v) && action === 'remove') action = 'add';
				});
			}

			elem.classList[action]('d-none');
		});
	}
}
