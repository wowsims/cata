import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Component } from '../components/component.js';
import { CopyButton } from '../components/copy_button.js';
import { Input, InputConfig } from '../components/input.js';
import { Player } from '../player.js';
import { PlayerSpecs } from '../player_specs';
import { Class, Spec } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { TypedEvent } from '../typed_event.js';
import { isRightClick } from '../utils.js';
import { classGlyphsConfig } from './factory';
import { GlyphsPicker } from './glyphs_picker';

export interface TalentsPickerConfig<ModObject, TalentsProto> extends InputConfig<ModObject, string> {
	playerClass: Class;
	playerSpec: Spec;
	tree: TalentsConfig<TalentsProto>;
}

export class TalentsPicker<ModObject extends Player<any>, TalentsProto> extends Input<ModObject, string> {
	readonly modObject: ModObject;

	private readonly config: TalentsPickerConfig<ModObject, TalentsProto>;

	readonly tree: TalentTreePicker<TalentsProto>;

	constructor(parent: HTMLElement, modObject: ModObject, config: TalentsPickerConfig<ModObject, TalentsProto>) {
		super(parent, 'talents-picker-root', modObject, { ...config });
		this.modObject = modObject;
		this.config = config;

		const containerElemRef = ref<HTMLDivElement>();
		const actionsContainerRef = ref<HTMLDivElement>();

		const talentsListRef = ref<HTMLDivElement>();

		this.rootElem.appendChild(
			<div className="talents-picker-inner" ref={containerElemRef}>
				<div className="talents-picker-header">
					<div className="talents-picker-actions" ref={actionsContainerRef} />
				</div>
				<div id="talents" className="talents-picker-list" ref={talentsListRef} />
			</div>,
		);

		const talentsListContainer = talentsListRef.value!;

		new CopyButton(actionsContainerRef.value!, {
			extraCssClasses: ['btn-sm', 'btn-outline-primary', 'copy-talents'],
			getContent: () => modObject.getTalentsString(),
			text: 'Copy',
			tooltip: 'Copy talent string',
		});

		this.tree = new TalentTreePicker(talentsListContainer, this.config.tree, this, config.playerSpec);
		this.tree.rows.forEach(row => row.forEach(talent => talent.setSelected(false)));

		if (this.isPlayer()) {
			new GlyphsPicker(this.rootElem, this.modObject, classGlyphsConfig[this.modObject.getClass()]);
		}

		this.init();
	}

	getInputElem(): HTMLElement {
		return this.rootElem;
	}

	getInputValue(): string {
		return this.tree.getTalentsString();
	}

	setInputValue(newValue: string) {
		this.tree.setTalentsString(newValue);
	}

	isPlayer(): this is TalentsPicker<Player<any>, TalentsProto> {
		return !!(this.modObject as unknown as Player<any>)?.playerClass;
	}
}

class TalentTreePicker<TalentsProto> extends Component {
	readonly rows: Array<Array<TalentPicker<TalentsProto>>>;
	readonly picker: TalentsPicker<any, TalentsProto>;

	constructor(parent: HTMLElement, config: TalentTreeConfig<TalentsProto>, picker: TalentsPicker<any, TalentsProto>, playerSpec: Spec) {
		super(parent, 'talent-tree-picker-root');
		this.picker = picker;

		const resetButton = ref<HTMLButtonElement>();
		this.rows = config.talents.reduce<Array<Array<TalentPicker<TalentsProto>>>>((acc, talent) => {
			if (!acc[talent.location.rowIdx]) acc[talent.location.rowIdx] = [];

			acc[talent.location.rowIdx][talent.location.colIdx] = new TalentPicker(null, talent, this);

			return acc;
		}, []);

		this.rootElem.replaceChildren(
			<>
				<div className="talent-tree-header">
					<img src={this.getTreeIcon(playerSpec)} className="talent-tree-icon" />
					<span className="talent-tree-title">{PlayerSpecs.fromProto(playerSpec).friendlyName}</span>
					<button ref={resetButton} className="talent-tree-reset btn link-danger">
						<i className="fa fa-times"></i>
					</button>
				</div>
				<div className="talent-tree-background" style={{ backgroundImage: `url('${config.backgroundUrl}')` }} />
				<div className="talent-tree-main">
					{this.rows.map((row, rowIdx) => (
						<div className="talent-tree-row">
							<div className="talent-tree-level">{(rowIdx + 1) * 15}</div>
							{row.map(talent => talent.rootElem)}
						</div>
					))}
				</div>
			</>,
		);

		const resetBtn = resetButton.value!;
		tippy(resetBtn, { content: 'Reset talent points' });
		resetBtn.addEventListener('click', _event => this.resetPoints());
	}

	getTalent(location: TalentLocation): TalentPicker<TalentsProto> {
		const talent = this.rows[location.rowIdx].find(talent => talent.getCol() == location.colIdx);
		if (!talent) throw new Error('No talent found with location: ' + location);
		return talent;
	}

	getTalentsString(): string {
		const selectedTalents = Array.from(Array(6), (_, rowIdx) => {
			const talent = this.rows[rowIdx].find(talent => talent.getRow() == rowIdx && talent.isSelected());
			return talent ? talent.getCol() + 1 : 0;
		});
		return selectedTalents.join('');
	}

	setTalentsString(str: string) {
		const talentRows = str.split('').map(Number);
		this.rows.forEach((row, rowIdx) =>
			row.forEach(talent => {
				talent.setSelected(talent.getCol() + 1 === talentRows[rowIdx]);
			}),
		);
	}

	resetPoints() {
		this.rows.forEach(row => row.forEach(talent => talent.setSelected(false)));
		this.picker.inputChanged(TypedEvent.nextEventID());
	}

	private getTreeIcon(playerSpec: number): string {
		return PlayerSpecs.fromProto(playerSpec).getIcon('medium');
	}
}

class TalentPicker<TalentsProto> extends Component {
	readonly config: TalentConfig<TalentsProto>;
	private readonly tree: TalentTreePicker<TalentsProto>;

	private icon: HTMLElement;
	private longTouchTimer?: number;
	private zIdx: number;

	constructor(parent: HTMLElement | null, config: TalentConfig<TalentsProto>, tree: TalentTreePicker<TalentsProto>) {
		super(parent, 'talent-picker-root', document.createElement('a'));
		this.config = config;
		this.tree = tree;
		this.zIdx = 0;

		const iconRef = ref<HTMLDivElement>();

		this.rootElem.replaceChildren(
			<>
				<div ref={iconRef} className="talent-picker-icon"></div>
				<div className="talent-picker-label" dataset={{ whtticon: false }}>
					{config.fancyName}
				</div>
			</>,
		);

		this.icon = iconRef.value!;

		this.rootElem.addEventListener('click', event => {
			event.preventDefault();
		});
		this.rootElem.addEventListener('contextmenu', event => {
			event.preventDefault();
		});
		this.rootElem.addEventListener('touchmove', _event => {
			if (this.longTouchTimer != undefined) {
				clearTimeout(this.longTouchTimer);
				this.longTouchTimer = undefined;
			}
		});
		this.rootElem.addEventListener('touchstart', event => {
			event.preventDefault();
			this.longTouchTimer = window.setTimeout(() => {
				this.setSelected(false);
				this.tree.picker.inputChanged(TypedEvent.nextEventID());
				this.longTouchTimer = undefined;
			}, 750);
		});
		this.rootElem.addEventListener('touchend', event => {
			event.preventDefault();
			if (this.longTouchTimer != undefined) {
				clearTimeout(this.longTouchTimer);
				this.longTouchTimer = undefined;
			} else {
				return;
			}
			this.setSelected(true);
			this.tree.picker.inputChanged(TypedEvent.nextEventID());
		});
		this.rootElem.addEventListener('mousedown', event => {
			const shouldAdd = !isRightClick(event);
			this.setSelected(shouldAdd);
			this.tree.picker.inputChanged(TypedEvent.nextEventID());
		});
	}

	get zIndex() {
		return this.zIdx;
	}

	set zIndex(z: number) {
		this.zIdx = z;
		this.rootElem.style.zIndex = String(this.zIdx);
	}

	getRow(): number {
		return this.config.location.rowIdx;
	}

	getCol(): number {
		return this.config.location.colIdx;
	}

	isSelected(): boolean {
		return this.rootElem.dataset.selected === 'true';
	}

	setSelected(isSelected: boolean) {
		if (isSelected) {
			const currentlySet = this.tree.rows[this.getRow()].find(talent => talent.getCol() !== this.getCol() && talent.isSelected());
			if (currentlySet?.isSelected()) currentlySet.setSelected(false);
		}

		this.rootElem.dataset.selected = String(isSelected);

		ActionId.fromSpellId(this.config.spellId)
			.fill()
			.then(actionId => {
				actionId.setWowheadHref(this.rootElem as HTMLAnchorElement);
				this.icon.style.backgroundImage = `url('${actionId.iconUrl}')`;
			});
	}
}

export type TalentsConfig<TalentsProto> = TalentTreeConfig<TalentsProto>;

export type TalentTreeConfig<TalentsProto> = {
	backgroundUrl: string;
	talents: Array<TalentConfig<TalentsProto>>;
};

export type TalentLocation = {
	// 0-indexed row in the tree
	rowIdx: number;
	// 0-indexed column in the tree
	colIdx: number;
};

export type TalentConfig<TalentsProto> = {
	fieldName: keyof TalentsProto | string;
	fancyName: string;
	location: TalentLocation;
	spellId: number;
};

export function newTalentsConfig<TalentsProto>(talentConfig: TalentsConfig<TalentsProto>): TalentsConfig<TalentsProto> {
	talentConfig.talents.forEach((talent, i) => {
		// Validate that talents are given in the correct order (left-to-right top-to-bottom).
		if (i != 0) {
			const prevTalent = talentConfig.talents[i - 1];
			if (
				talent.location.rowIdx < prevTalent.location.rowIdx ||
				(talent.location.rowIdx == prevTalent.location.rowIdx && talent.location.colIdx <= prevTalent.location.colIdx)
			) {
				throw new Error(`Out-of-order talent: ${String(talent.fancyName)}`);
			}
		}
	});

	return talentConfig;
}
