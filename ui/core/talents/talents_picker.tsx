import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Component } from '../components/component.js';
import { CopyButton } from '../components/copy_button.js';
import { Input, InputConfig } from '../components/input.js';
import { Player } from '../player.js';
import { PlayerClasses } from '../player_classes';
import { Class, Spec } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { TypedEvent } from '../typed_event.js';
import { isRightClick, sum } from '../utils.js';
import { classGlyphsConfig } from './factory';
import { GlyphsPicker } from './glyphs_picker';
import { HunterPet } from './hunter_pet';

export const MAX_POINTS_PLAYER_MOP = 6;
export const MAX_POINTS_PLAYER = 41;
const MAX_POINTS_HUNTER_PET = 17;
const MAX_POINTS_HUNTER_PET_BM = 21;

const PLAYER_TREES_UNLOCK_POINT_THRESHOLD = 31;

export interface TalentsPickerConfig<ModObject, TalentsProto> extends InputConfig<ModObject, string> {
	playerClass: Class;
	pointsPerRow: number;
	trees: TalentsConfig<TalentsProto>;
}

export class TalentsPicker<ModObject extends Player<any> | HunterPet<any>, TalentsProto> extends Input<ModObject, string> {
	readonly modObject: ModObject;

	readonly numRows: number;
	readonly numCols: number;
	readonly pointsPerRow: number;
	readonly mopTalents: boolean;

	maxPoints: number;

	private readonly config: TalentsPickerConfig<ModObject, TalentsProto>;

	readonly trees: Array<TalentTreePicker<TalentsProto>>;

	constructor(parent: HTMLElement, modObject: ModObject, config: TalentsPickerConfig<ModObject, TalentsProto>) {
		super(parent, 'talents-picker-root', modObject, { ...config });
		this.modObject = modObject;
		this.config = config;
		this.mopTalents = config.playerClass === Class.ClassMonk;
		this.pointsPerRow = this.mopTalents ? 1 : config.pointsPerRow;
		this.numRows = Math.max(...config.trees.map(treeConfig => treeConfig.talents.map(talentConfig => talentConfig.location.rowIdx).flat()).flat()) + 1;
		this.numCols = Math.max(...config.trees.map(treeConfig => treeConfig.talents.map(talentConfig => talentConfig.location.colIdx).flat()).flat()) + 1;

		if (this.isHunterPet()) {
			if (this.modObject.player.isSpec(Spec.SpecBeastMasteryHunter)) {
				this.maxPoints = MAX_POINTS_HUNTER_PET_BM;
			} else {
				this.maxPoints = MAX_POINTS_HUNTER_PET;
			}
		} else {
			this.maxPoints = this.mopTalents ? MAX_POINTS_PLAYER_MOP : MAX_POINTS_PLAYER;
		}

		const getPointsRemaining = (): number => this.maxPoints - modObject.getTalentTreePoints().reduce((sum, points) => sum + points, 0);

		const containerElemRef = ref<HTMLDivElement>();
		const pointsRemainingElemRef = ref<HTMLSpanElement>();
		const actionsContainerRef = ref<HTMLDivElement>();

		const carouselContainerRef = ref<HTMLDivElement>();
		const carouselPrevBtnRef = ref<HTMLButtonElement>();
		const carouselNextBtnRef = ref<HTMLButtonElement>();

		this.rootElem.appendChild(
			<div className="talents-picker-inner" ref={containerElemRef}>
				<div className="talents-picker-header">
					<div className="d-flex">
						<label>
							Points Remaining:
							<span ref={pointsRemainingElemRef}>{getPointsRemaining()}</span>
						</label>
					</div>
					<div className="talents-picker-actions" ref={actionsContainerRef} />
				</div>
				<div id="talents-carousel" className="carousel slide">
					<div className="carousel-inner" ref={carouselContainerRef}></div>
					<button className="carousel-control-prev" type="button" ref={carouselPrevBtnRef}>
						<span className="carousel-control-prev-icon" attributes={{ 'aria-hidden': true }} />
						<span className="visually-hidden">Previous</span>
					</button>
					<button className="carousel-control-next" type="button" ref={carouselNextBtnRef}>
						<span className="carousel-control-next-icon" attributes={{ 'aria-hidden': true }} />
						<span className="visually-hidden">Next</span>
					</button>
				</div>
			</div>,
		);

		const containerElem = containerElemRef.value!;
		const carouselContainer = carouselContainerRef.value!;
		const carouselPrevBtn = carouselPrevBtnRef.value!;
		const carouselNextBtn = carouselNextBtnRef.value!;

		modObject.talentsChangeEmitter.on(() => (pointsRemainingElemRef.value!.textContent = `${getPointsRemaining()}`));

		new CopyButton(actionsContainerRef.value!, {
			extraCssClasses: ['btn-sm', 'btn-outline-primary', 'copy-talents'],
			getContent: () => modObject.getTalentsString(),
			text: 'Copy',
			tooltip: 'Copy talent string',
		});

		this.trees = this.config.trees.map((treeConfig, i) => {
			const carouselItem = document.createElement('div');
			carouselContainer.appendChild(carouselItem);

			carouselItem.classList.add('carousel-item');
			// Set middle talents active by default for mobile slider
			if (i === 1) carouselItem.classList.add('active');

			return new TalentTreePicker(carouselItem, treeConfig, this, config.playerClass, i);
		});
		this.trees.forEach(tree => tree.talents.forEach(talent => talent.setPoints(0, false)));

		if (!this.isHunterPet()) {
			let carouselitemIdx = 0;
			const slidePrev = () => {
				if (carouselitemIdx >= 1) return;
				carouselitemIdx += 1;
				carouselContainer.style.transform = `translateX(${33.3 * carouselitemIdx}%)`;
				carouselContainer.children[Math.abs(carouselitemIdx - 2) % 3]!.classList.remove('active');
				carouselContainer.children[Math.abs(carouselitemIdx - 1) % 3]!.classList.add('active');
			};
			const slideNext = () => {
				if (carouselitemIdx <= -1) return;
				carouselitemIdx -= 1;
				carouselContainer.style.transform = `translateX(${33.3 * carouselitemIdx}%)`;
				carouselContainer.children[Math.abs(carouselitemIdx) % 3]!.classList.remove('active');
				carouselContainer.children[Math.abs(carouselitemIdx) + (1 % 3)]!.classList.add('active');
			};

			carouselPrevBtn.addEventListener('click', slidePrev);
			carouselNextBtn.addEventListener('click', slideNext);

			if (this.isPlayer()) {
				new GlyphsPicker(this.rootElem, this.modObject, classGlyphsConfig[this.modObject.getClass()]);
			}
		}

		this.init();
		this.updatePlayerTrees();
	}

	getInputElem(): HTMLElement {
		return this.rootElem;
	}

	getInputValue(): string {
		return this.trees
			.map(tree => tree.getTalentsString())
			.join('-')
			.replace(/-+$/g, '');
	}

	setInputValue(newValue: string) {
		const parts = newValue.split('-');
		this.trees.forEach((tree, idx) => tree.setTalentsString(parts[idx] || ''));
		this.updateTrees();
	}

	updateTrees() {
		if (this.isFull()) {
			this.rootElem.classList.add('talents-full');
		} else {
			this.rootElem.classList.remove('talents-full');
		}
		this.trees.forEach(tree => tree.update());

		// Disable other player trees if the first tree is not filled to 31 points
		this.updatePlayerTrees();
	}

	get numPoints() {
		return sum(this.trees.map(tree => tree.numPoints));
	}

	isFull() {
		return this.numPoints >= this.maxPoints;
	}

	setMaxPoints(newMaxPoints: number) {
		if (newMaxPoints != this.maxPoints) {
			this.maxPoints = newMaxPoints;
			this.updateTrees();
		}
	}

	isPlayer(): this is TalentsPicker<Player<any>, TalentsProto> {
		return !!(this.modObject as unknown as Player<any>)?.playerClass;
	}

	isHunterPet(): this is TalentsPicker<HunterPet<any>, TalentsProto> {
		return !this.isPlayer();
	}

	private updatePlayerTrees() {
		if (this.isPlayer()) {
			if (this.mopTalents) {
				this.trees[0].rootElem.classList.remove('disabled');
				return;
			}

			const specNumber = this.modObject.getPlayerSpec().specIndex;
			const pointsSpent = this.trees[specNumber].numPoints;

			if (pointsSpent < PLAYER_TREES_UNLOCK_POINT_THRESHOLD) {
				[0, 1, 2].forEach(i => {
					if (![specNumber].includes(i)) {
						this.trees[i].rootElem.classList.add('disabled');
						this.trees[i].resetPoints();
					} else {
						this.trees[i].rootElem.classList.remove('disabled');
					}
				});
			} else {
				this.enableSecondaryPlayerTrees();
			}
		}
	}

	private enableSecondaryPlayerTrees() {
		this.trees.forEach(tree => tree.rootElem.classList.remove('disabled'));
	}
}

class TalentTreePicker<TalentsProto> extends Component {
	private readonly config: TalentTreeConfig<TalentsProto>;
	private readonly title: HTMLElement;
	private readonly pointsElem: HTMLElement;

	readonly talents: Array<TalentPicker<TalentsProto>>;
	readonly picker: TalentsPicker<any, TalentsProto>;

	// The current number of points in this tree
	numPoints: number;

	constructor(parent: HTMLElement, config: TalentTreeConfig<TalentsProto>, picker: TalentsPicker<any, TalentsProto>, klass: Class, specNumber: number) {
		super(parent, 'talent-tree-picker-root');
		this.config = config;
		this.numPoints = 0;
		this.picker = picker;

		this.rootElem.appendChild(
			<>
				<div className="talent-tree-header">
					<img src={this.getTreeIcon(klass, specNumber)} className="talent-tree-icon" />
					<span className="talent-tree-title" />
					<label className="talent-tree-points" />
					<button className="talent-tree-reset btn link-danger">
						<i className="fa fa-times"></i>
					</button>
				</div>
				<div className="talent-tree-background" />
				<div className="talent-tree-main" />
			</>,
		);

		this.title = this.rootElem.getElementsByClassName('talent-tree-title')[0] as HTMLElement;
		this.pointsElem = this.rootElem.querySelector('.talent-tree-points') as HTMLElement;

		const background = this.rootElem.querySelector('.talent-tree-background') as HTMLElement;
		background.style.backgroundImage = `url('${config.backgroundUrl}')`;

		const main = this.rootElem.querySelector('.talent-tree-main') as HTMLElement;
		main.style.gridTemplateRows = `repeat(${this.picker.numRows}, 1fr)`;
		// Add 2 for spacing on the sides
		main.style.gridTemplateColumns = `repeat(${this.picker.numCols}, 1fr)`;

		const iconSize = '3.5rem';
		main.style.height = `calc(${iconSize} * ${this.picker.numRows})`;
		main.style.maxWidth = `calc(${iconSize} * ${this.picker.numCols})`;
		this.rootElem.style.maxWidth = `calc(${iconSize} * ${this.picker.numCols + 2})`;

		this.talents = config.talents.map(talent => new TalentPicker(main, talent, this));
		// Process parent<->child mapping
		this.talents.forEach(talent => {
			if (talent.config.prereqLocation) {
				this.getTalent(talent.config.prereqLocation).config.childLocations!.push(talent.config.location);
			}
		});
		// Loop through all and have talent add in divs/items for child dependencies
		// It'd be nicer to have this in talent constructor but json would have to be updated
		const recurseCalcIdx = (t: TalentPicker<TalentsProto>, z: number) => {
			t.initChildReqs();
			t.zIndex = z;
			for (const cl of t.config.childLocations!) {
				const c = this.getTalent(cl);
				c.parentReq = t.getChildReqArrow(cl);
				recurseCalcIdx(c, z - 2);
			}
		};
		// Start at top of each heirachy chain and recurse down
		for (const t of this.talents) {
			if (t.config.childLocations!.length == 0) continue;
			if (t.config.prereqLocation !== undefined) continue;
			recurseCalcIdx(t, 20);
		}
		const resetBtn = this.rootElem.querySelector('.talent-tree-reset') as HTMLElement;
		tippy(resetBtn, { content: 'Reset talent points' });
		resetBtn.addEventListener('click', _event => this.resetPoints());
	}

	update() {
		this.title.innerHTML = this.config.name;
		this.pointsElem.textContent = `${this.numPoints} / ${this.picker.maxPoints}`;
		this.talents.forEach(talent => talent.update());
	}

	getTalent(location: TalentLocation): TalentPicker<TalentsProto> {
		const talent = this.talents.find(talent => talent.getRow() == location.rowIdx && talent.getCol() == location.colIdx);
		if (!talent) throw new Error('No talent found with location: ' + location);
		return talent;
	}

	getTalentsString(): string {
		return this.talents
			.map(talent => String(talent.getPoints()))
			.join('')
			.replace(/0+$/g, '');
	}

	setTalentsString(str: string) {
		this.talents.forEach((talent, idx) => talent.setPoints(Number(str.charAt(idx)), false));
	}

	resetPoints() {
		this.talents.forEach(talent => talent.setPoints(0, false));
		this.picker.inputChanged(TypedEvent.nextEventID());
	}

	private getTreeIcon(klass: Class, specNumber: number): string {
		if (this.picker.isHunterPet()) {
			const fileName = ['ability_druid_swipe.jpg', 'ability_hunter_pet_bear.jpg', 'ability_hunter_combatexperience.jpg'][specNumber];
			return `https://wow.zamimg.com/images/wow/icons/medium/${fileName}`;
		} else {
			return Object.values(PlayerClasses.fromProto(klass).specs)[specNumber].getIcon('medium');
		}
	}
}

type ReqDir = 'down' | 'right' | 'left' | 'rightdown' | 'leftdown';
class TalentReqArrow extends Component {
	private dir: ReqDir;
	private zIdx: number;
	readonly parentLoc: TalentLocation;
	readonly childLoc: TalentLocation;

	constructor(parent: HTMLElement, parentLoc: TalentLocation, childLoc: TalentLocation) {
		super(parent, 'talent-picker-req-arrow', document.createElement('div'));
		this.zIdx = 0;
		this.parentLoc = parentLoc;
		this.childLoc = childLoc;

		this.rootElem.style.gridRow = String(parentLoc.rowIdx + 1);
		this.rootElem.style.gridColumn = String(parentLoc.colIdx + 1);

		let rowEnd = Math.max(parentLoc.rowIdx, childLoc.rowIdx) + 1;
		let colEnd = Math.max(parentLoc.colIdx, childLoc.colIdx) + 1;

		// Calculate where we need to 'point'
		if (parentLoc.rowIdx == childLoc.rowIdx) {
			this.dir = parentLoc.colIdx < childLoc.colIdx ? 'right' : 'left';
			this.rootElem.dataset.reqArrowColSize = String(Math.abs(parentLoc.colIdx - childLoc.colIdx));
			colEnd = this.dir == 'left' ? colEnd + 1 : colEnd - 1;
		} else {
			if (parentLoc.colIdx == childLoc.colIdx) {
				this.dir = 'down';
				this.rootElem.dataset.reqArrowRowSize = String(Math.abs(parentLoc.rowIdx - childLoc.rowIdx));
				rowEnd += 1;
			} else {
				this.dir = parentLoc.colIdx < childLoc.colIdx ? 'rightdown' : 'leftdown';
				this.rootElem.dataset.reqArrowColSize = String(Math.abs(parentLoc.colIdx - childLoc.colIdx));
				this.rootElem.dataset.reqArrowRowSize = String(Math.abs(parentLoc.rowIdx - childLoc.rowIdx));
				rowEnd += 1;
				colEnd = this.dir == 'rightdown' ? colEnd + 1 : colEnd - 1;
				this.rootElem.appendChild(<div />);
			}
		}

		this.rootElem.style.gridRowEnd = String(rowEnd);
		this.rootElem.style.gridColumnEnd = String(colEnd);
		this.rootElem.classList.add(`talent-picker-req-arrow-${this.dir}`);
	}

	get zIndex() {
		return this.zIdx;
	}

	set zIndex(z: number) {
		this.zIdx = z;
		this.rootElem.style.zIndex = String(z);
	}

	setReqFufilled(isFufilled: boolean) {
		if (isFufilled) this.rootElem.dataset.reqActive = 'true';
		else delete this.rootElem.dataset.reqActive;
	}
}

class TalentPicker<TalentsProto> extends Component {
	readonly config: TalentConfig<TalentsProto>;
	private readonly tree: TalentTreePicker<TalentsProto>;
	private readonly pointsDisplay?: HTMLElement;

	private longTouchTimer?: number;
	private childReqs: TalentReqArrow[];
	private zIdx: number;
	parentReq: TalentReqArrow | null;

	constructor(parent: HTMLElement, config: TalentConfig<TalentsProto>, tree: TalentTreePicker<TalentsProto>) {
		super(parent, 'talent-picker-root', document.createElement('a'));
		this.config = config;
		this.tree = tree;
		this.childReqs = [];
		this.parentReq = null;
		this.zIdx = 0;

		this.rootElem.style.gridRow = String(this.config.location.rowIdx + 1);
		this.rootElem.style.gridColumn = String(this.config.location.colIdx + 1);

		this.rootElem.dataset.maxPoints = String(this.config.maxPoints);
		this.rootElem.dataset.whtticon = 'false';

		if (!this.tree.picker.mopTalents) {
			this.pointsDisplay = document.createElement('span');
			this.pointsDisplay.classList.add('talent-picker-points');
			this.rootElem.appendChild(this.pointsDisplay);
		}

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
				this.setPoints(0, true);
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
			let newPoints = this.getPoints() + 1;
			if (this.config.maxPoints < newPoints) {
				newPoints = 0;
			}
			this.setPoints(newPoints, true);
			this.tree.picker.inputChanged(TypedEvent.nextEventID());
		});
		this.rootElem.addEventListener('mousedown', event => {
			const rightClick = isRightClick(event);
			if (rightClick) {
				this.setPoints(this.getPoints() - 1, true);
			} else {
				this.setPoints(this.getPoints() + 1, true);
			}
			this.tree.picker.inputChanged(TypedEvent.nextEventID());
		});
	}

	initChildReqs(): void {
		if (this.config.childLocations!.length == 0) return;

		for (const c of this.config.childLocations!) {
			this.childReqs.push(new TalentReqArrow(this.rootElem.parentElement!, this.config.location, c));
		}
	}

	getChildReqArrow(loc: TalentLocation): TalentReqArrow {
		for (const c of this.childReqs) {
			if (c.childLoc === loc) {
				return c;
			}
		}
		throw Error('missing child prereq?');
	}

	get zIndex() {
		return this.zIdx;
	}

	set zIndex(z: number) {
		this.zIdx = z;
		this.rootElem.style.zIndex = String(this.zIdx);

		for (const c of this.childReqs) {
			c.zIndex = this.zIdx - 1;
		}
	}

	getRow(): number {
		return this.config.location.rowIdx;
	}

	getCol(): number {
		return this.config.location.colIdx;
	}

	getPoints(): number {
		const pts = Number(this.rootElem.dataset.points);
		return isNaN(pts) ? 0 : pts;
	}

	isFull(): boolean {
		return this.getPoints() >= this.config.maxPoints;
	}

	// Returns whether setting the points to newPoints would be a valid talent tree.
	canSetPoints(newPoints: number): boolean {
		if (this.tree.picker.mopTalents) {
			return true;
		}

		const oldPoints = this.getPoints();

		if (newPoints > oldPoints) {
			const additionalPoints = newPoints - oldPoints;

			if (this.tree.picker.numPoints + additionalPoints > this.tree.picker.maxPoints) {
				return false;
			}

			if (this.tree.numPoints < this.getRow() * this.tree.picker.pointsPerRow) {
				return false;
			}

			if (this.config.prereqLocation) {
				if (!this.tree.getTalent(this.config.prereqLocation).isFull()) return false;
			}
		} else {
			const removedPoints = oldPoints - newPoints;

			// Figure out whether any lower talents would have the row requirement
			// broken by subtracting points.
			const pointTotalsByRow = [...Array(this.tree.picker.numRows).keys()]
				.map(rowIdx => this.tree.talents.filter(talent => talent.getRow() == rowIdx))
				.map(talentsInRow => sum(talentsInRow.map(talent => talent.getPoints())));
			pointTotalsByRow[this.getRow()] -= removedPoints;

			const cumulativeTotalsByRow = pointTotalsByRow.map((_, rowIdx) => sum(pointTotalsByRow.slice(0, rowIdx + 1)));

			if (
				!this.tree.talents.every(
					talent =>
						talent.getPoints() == 0 ||
						talent.getRow() == 0 ||
						cumulativeTotalsByRow[talent.getRow() - 1] >= talent.getRow() * this.tree.picker.pointsPerRow,
				)
			) {
				return false;
			}

			for (const c of this.config.childLocations!) {
				if (this.tree.getTalent(c).getPoints() > 0) return false;
			}
		}
		return true;
	}

	setPoints(newPoints: number, checkValidity: boolean) {
		const oldPoints = this.getPoints();
		newPoints = Math.max(0, newPoints);
		newPoints = Math.min(this.config.maxPoints, newPoints);

		if (checkValidity && !this.canSetPoints(newPoints)) return;

		if (this.tree.picker.mopTalents && newPoints > 0) {
			const currentlySet = this.tree.talents.find(
				talent => talent.getRow() === this.getRow() && talent.getCol() !== this.getCol() && talent.getPoints() > 0,
			);
			if (currentlySet) {
				const points = currentlySet.getPoints();
				if (points > 0) {
					currentlySet.setPoints(points - 1, false);
				}
			}
		}

		this.tree.numPoints += newPoints - oldPoints;
		this.rootElem.dataset.points = String(newPoints);

		if (this.pointsDisplay) {
			this.pointsDisplay.textContent = newPoints + '/' + this.config.maxPoints;
		}

		if (this.isFull()) {
			this.rootElem.classList.add('talent-full');
		} else {
			this.rootElem.classList.remove('talent-full');
		}

		const spellId = this.getSpellIdForPoints(newPoints);
		ActionId.fromSpellId(spellId)
			.fill()
			.then(actionId => {
				actionId.setWowheadHref(this.rootElem as HTMLAnchorElement);
				this.rootElem.style.backgroundImage = `url('${actionId.iconUrl}')`;
			});
	}

	getSpellIdForPoints(numPoints: number): number {
		// 0-indexed rank of talent
		const rank = Math.max(0, numPoints - 1);

		if (this.config.spellIds[rank]) {
			return this.config.spellIds[rank];
		} else {
			throw new Error(`No rank ${numPoints} for talent ${String(this.config.fieldName)}`);
		}
	}

	update() {
		let canSetPoints: boolean;
		if (this.tree.picker.mopTalents) {
			canSetPoints = !this.tree.talents.find(talent => talent.getRow() === this.getRow() && talent.getPoints() > 0);
		} else {
			canSetPoints = this.canSetPoints(this.getPoints() + 1);
		}

		if (canSetPoints) {
			this.rootElem.classList.add('talent-picker-can-add');
		} else {
			this.rootElem.classList.remove('talent-picker-can-add');
		}

		if (this.parentReq) {
			this.parentReq.setReqFufilled(canSetPoints || this.isFull());
		}
	}
}

export type TalentsConfig<TalentsProto> = Array<TalentTreeConfig<TalentsProto>>;

export type TalentTreeConfig<TalentsProto> = {
	name: string;
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
	fieldName?: keyof TalentsProto | string;

	location: TalentLocation;

	// Location of a prerequisite talent, if any
	prereqLocation?: TalentLocation;

	// Child talents depending on this talent. This is populated automatically.
	childLocations?: TalentLocation[];

	// Spell ID for each rank of this talent.
	// Omitted ranks will be inferred by incrementing from the last provided rank.
	spellIds: Array<number>;

	maxPoints: number;
};

export function newTalentsConfig<TalentsProto>(talents: TalentsConfig<TalentsProto>): TalentsConfig<TalentsProto> {
	talents.forEach(tree => {
		tree.talents.forEach((talent, i) => {
			talent.childLocations = [];
			// Validate that talents are given in the correct order (left-to-right top-to-bottom).
			if (i != 0) {
				const prevTalent = tree.talents[i - 1];
				if (
					talent.location.rowIdx < prevTalent.location.rowIdx ||
					(talent.location.rowIdx == prevTalent.location.rowIdx && talent.location.colIdx <= prevTalent.location.colIdx)
				) {
					throw new Error(`Out-of-order talent: ${String(talent.fieldName)}`);
				}
			}

			// Infer omitted spell IDs.
			if (talent.spellIds.length < talent.maxPoints) {
				let curSpellId = talent.spellIds[talent.spellIds.length - 1];
				for (let pointIdx = talent.spellIds.length; pointIdx < talent.maxPoints; pointIdx++) {
					curSpellId++;
					talent.spellIds.push(curSpellId);
				}
			}
		});
	});
	return talents;
}
