// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment, ref } from 'tsx-vanilla';

import { BaseModal } from '../components/base_modal.js';
import { Component } from '../components/component.js';
import { ContentBlock } from '../components/content_block.js';
import { Input } from '../components/input.js';
import { getLanguageCode } from '../constants/lang';
import { setItemQualityCssClass } from '../css_utils.js';
import { Player } from '../player.js';
import { Glyphs, ItemQuality } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { stringComparator } from '../utils.js';

export type GlyphConfig = {
	name: string;
	description: string;
	iconUrl: string;
};

export type GlyphsConfig = {
	primeGlyphs: Record<number, GlyphConfig>;
	majorGlyphs: Record<number, GlyphConfig>;
	minorGlyphs: Record<number, GlyphConfig>;
};

interface GlyphData {
	id: number;
	name: string;
	description: string;
	iconUrl: string;
	quality: ItemQuality | null;
}

const emptyGlyphData: GlyphData = {
	id: 0,
	name: 'Empty',
	description: '',
	iconUrl: 'https://wow.zamimg.com/images/wow/icons/medium/inventoryslot_empty.jpg',
	quality: null,
};

export class GlyphsPicker extends Component {
	private readonly glyphsConfig: GlyphsConfig;

	primeGlyphPickers: Array<GlyphPicker> = [];
	majorGlyphPickers: Array<GlyphPicker> = [];
	minorGlyphPickers: Array<GlyphPicker> = [];

	constructor(parent: HTMLElement, player: Player<any>, glyphsConfig: GlyphsConfig) {
		super(parent, 'glyphs-picker-root');
		this.glyphsConfig = glyphsConfig;

		const primeGlyphs = Object.keys(glyphsConfig.primeGlyphs).map(idStr => Number(idStr));
		const majorGlyphs = Object.keys(glyphsConfig.majorGlyphs).map(idStr => Number(idStr));
		const minorGlyphs = Object.keys(glyphsConfig.minorGlyphs).map(idStr => Number(idStr));

		const primeGlyphsData = primeGlyphs.map(glyph => this.getGlyphData(glyph));
		const majorGlyphsData = majorGlyphs.map(glyph => this.getGlyphData(glyph));
		const minorGlyphsData = minorGlyphs.map(glyph => this.getGlyphData(glyph));

		primeGlyphsData.sort((a, b) => stringComparator(a.name, b.name));
		majorGlyphsData.sort((a, b) => stringComparator(a.name, b.name));
		minorGlyphsData.sort((a, b) => stringComparator(a.name, b.name));

		const primeGlyphsBlock = new ContentBlock(this.rootElem, 'prime-glyphs', {
			header: { title: 'Prime Glyphs', extraCssClasses: ['border-0', 'mb-1'] },
		});

		const majorGlyphsBlock = new ContentBlock(this.rootElem, 'major-glyphs', {
			header: { title: 'Major Glyphs', extraCssClasses: ['border-0', 'mb-1'] },
		});

		const minorGlyphsBlock = new ContentBlock(this.rootElem, 'minor-glyphs', {
			header: { title: 'Minor Glyphs', extraCssClasses: ['border-0', 'mb-1'] },
		});

		this.primeGlyphPickers = (['prime1', 'prime2', 'prime3'] as Array<keyof Glyphs>).map(glyphField => {
			return new GlyphPicker(primeGlyphsBlock.bodyElement, player, primeGlyphsData, glyphField);
		});

		this.majorGlyphPickers = (['major1', 'major2', 'major3'] as Array<keyof Glyphs>).map(glyphField => {
			return new GlyphPicker(majorGlyphsBlock.bodyElement, player, majorGlyphsData, glyphField);
		});

		this.minorGlyphPickers = (['minor1', 'minor2', 'minor3'] as Array<keyof Glyphs>).map(glyphField => {
			return new GlyphPicker(minorGlyphsBlock.bodyElement, player, minorGlyphsData, glyphField);
		});
	}

	// In case we ever want to parse description from tooltip HTML.
	//static descriptionRegex = /<a href=\\"\/wotlk.*>(.*)<\/a>/g;
	getGlyphData(glyph: number): GlyphData {
		const glyphConfig = this.glyphsConfig.primeGlyphs[glyph] || this.glyphsConfig.majorGlyphs[glyph] || this.glyphsConfig.minorGlyphs[glyph];

		return {
			id: glyph,
			name: glyphConfig.name,
			description: glyphConfig.description,
			iconUrl: glyphConfig.iconUrl,
			quality: ItemQuality.ItemQualityCommon,
		};
	}
}

class GlyphPicker extends Input<Player<any>, number> {
	readonly player: Player<any>;

	selectedGlyph: GlyphData | undefined;

	private readonly glyphOptions: Array<GlyphData>;

	private readonly anchorElem: HTMLAnchorElement;
	private readonly iconElem: HTMLImageElement;
	private readonly nameElem: HTMLSpanElement;

	constructor(parent: HTMLElement, player: Player<any>, glyphOptions: Array<GlyphData>, glyphField: keyof Glyphs) {
		super(parent, 'glyph-picker-root', player, {
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
			<a ref={anchorElemRef} attributes={{ role: 'button' }} className="d-flex w-100">
				<img ref={iconElemRef} className="item-picker-icon" />
				<div className="item-picker-labels-container">
					<span ref={nameElemRef} className="item-picker-name" />
				</div>
			</a>,
		);

		this.anchorElem = anchorElemRef.value!;
		this.iconElem = iconElemRef.value!;
		this.nameElem = nameElemRef.value!;

		const selectorModal = new GlyphSelectorModal(this.rootElem.closest('.individual-sim-ui')!, this, this.glyphOptions);
		const openGlyphSelectorModal = (event: Event) => {
			event.preventDefault();
			selectorModal.open();
		};

		this.anchorElem.addEventListener('click', openGlyphSelectorModal);
		this.addOnDisposeCallback(() => this.anchorElem.removeEventListener('click', openGlyphSelectorModal));

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
			const lang = getLanguageCode();
			const langPrefix = lang ? `${lang}.` : '';

			this.anchorElem.href = ActionId.makeItemUrl(this.selectedGlyph.id);
			this.anchorElem.dataset.wowhead = `domain=${langPrefix}cata&dataEnv=11`;
			this.anchorElem.dataset.whtticon = 'false';

			this.iconElem.src = this.selectedGlyph.iconUrl;

			this.nameElem.textContent = this.selectedGlyph.name;
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
	constructor(parent: HTMLElement, glyphPicker: GlyphPicker, glyphOptions: Array<GlyphData>) {
		super(parent, 'glyph-modal', { title: 'Glyphs' });

		this.body.innerHTML = `
			<div class="input-root">
				<input class="selector-modal-search form-control" type="text" placeholder="Search...">
			</div>
			<ul class="selector-modal-list"></ul>
		`;

		const listElem = this.rootElem.getElementsByClassName('selector-modal-list')[0] as HTMLElement;

		glyphOptions = [emptyGlyphData].concat(glyphOptions);
		const listItemElems = glyphOptions.map((glyphData, glyphIdx) => {
			const listItemElem = document.createElement('li');
			listItemElem.classList.add('selector-modal-list-item');
			listElem.appendChild(listItemElem);

			listItemElem.dataset.idx = String(glyphIdx);

			listItemElem.innerHTML = `
				<a class="selector-modal-list-item-link">
					<img class="selector-modal-list-item-icon" />
					<label class="selector-modal-list-item-name">${glyphData.name}</label>
					<span class="selector-modal-list-item-description">${glyphData.description}</span>
				</a>
      `;

			const anchorElem = listItemElem.children[0] as HTMLAnchorElement;
			const iconElem = listItemElem.querySelector('.selector-modal-list-item-icon') as HTMLImageElement;
			const nameElem = listItemElem.querySelector('.selector-modal-list-item-name') as HTMLElement;

			anchorElem.href = glyphData.id == 0 ? '' : ActionId.makeItemUrl(glyphData.id);
			anchorElem.addEventListener('click', (event: Event) => {
				event.preventDefault();
				glyphPicker.setValue(TypedEvent.nextEventID(), glyphData.id);
			});
			iconElem.src = glyphData.iconUrl;
			setItemQualityCssClass(nameElem, glyphData.quality);

			return listItemElem;
		});

		const updateSelected = () => {
			const selectedGlyphId = glyphPicker.selectedGlyph?.id ?? 0;

			listItemElems.forEach(elem => {
				const listItemIdx = parseInt(elem.dataset.idx!);
				const listItemData = glyphOptions[listItemIdx];

				elem.classList.remove('active');
				if (listItemData.id == selectedGlyphId) {
					elem.classList.add('active');
				}
			});
		};
		updateSelected();

		const applyFilters = () => {
			let validItemElems = listItemElems;

			validItemElems = validItemElems.filter(elem => {
				const listItemIdx = parseInt(elem.dataset.idx!);
				const listItemData = glyphOptions[listItemIdx];

				if (searchInput.value.length > 0) {
					const searchQuery = searchInput.value.toLowerCase().split(' ');
					const name = listItemData.name.toLowerCase();

					let include = true;
					searchQuery.forEach(v => {
						if (!name.includes(v)) include = false;
					});
					if (!include) {
						return false;
					}
				}

				return true;
			});

			listElem.innerHTML = ``;
			listElem.append(...validItemElems);
		};

		const searchInput = this.rootElem.getElementsByClassName('selector-modal-search')[0] as HTMLInputElement;
		searchInput.addEventListener('input', applyFilters);
		searchInput.addEventListener('keyup', ev => {
			if (ev.key == 'Enter') {
				listItemElems.find(ele => {
					if (ele.classList.contains('hidden')) {
						return false;
					}
					const nameElem = ele.getElementsByClassName('selector-modal-list-item-name')[0] as HTMLElement;
					nameElem.click();
					return true;
				});
			}
		});

		glyphPicker.player.glyphsChangeEmitter.on(() => {
			applyFilters();
			updateSelected();
		});
	}
}
