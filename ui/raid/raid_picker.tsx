import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { BaseModal } from '../core/components/base_modal.jsx';
import { Component } from '../core/components/component.js';
import { EnumPicker } from '../core/components/pickers/enum_picker.js';
import { MAX_PARTY_SIZE, Party } from '../core/party.js';
import { Player } from '../core/player.js';
import { PlayerClasses } from '../core/player_classes/index.js';
import { PlayerSpecs } from '../core/player_specs/index.js';
import { Player as PlayerProto } from '../core/proto/api.js';
import { Class, Faction, Glyphs, Profession, Spec } from '../core/proto/common.js';
import { UnholyDeathKnight_Options } from '../core/proto/death_knight.js';
import { BalanceDruid_Options as BalanceDruidOptions } from '../core/proto/druid.js';
import { ArcaneMage_Options } from '../core/proto/mage.js';
import { getPlayerSpecFromPlayer, newUnitReference } from '../core/proto_utils/utils.js';
import { Raid } from '../core/raid.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { formatDeltaTextElem, getEnumValues } from '../core/utils.js';
import { playerPresets, specSimFactories } from './presets.js';
import { RaidSimUI } from './raid_sim_ui.js';

const NEW_PLAYER = -1;

const LATEST_PHASE_WITH_ALL_PRESETS = Math.min(
	...playerPresets.map(preset => Math.max(...Object.keys(preset.defaultGear[Faction.Alliance]).map(k => parseInt(k)))),
);

enum DragType {
	None,
	New,
	Move,
	Swap,
	Copy,
}

export class RaidPicker extends Component {
	readonly raidSimUI: RaidSimUI;
	readonly raid: Raid;
	readonly partyPickers: Array<PartyPicker>;
	readonly newPlayerPicker: NewPlayerPicker;
	readonly playerEditorModal: PlayerEditorModal<Spec>;

	// Hold data about the player being dragged while the drag is happening.
	currentDragPlayer: Player<any> | null = null;
	currentDragPlayerFromIndex: number = NEW_PLAYER;
	currentDragType: DragType = DragType.New;

	// Hold data about the party being dragged while the drag is happening.
	currentDragParty: PartyPicker | null = null;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI) {
		super(parent, 'raid-picker-root');
		this.raidSimUI = raidSimUI;
		this.raid = raidSimUI.sim.raid;

		const raidControls = document.createElement('div');
		raidControls.classList.add('raid-controls');
		this.rootElem.appendChild(raidControls);

		this.newPlayerPicker = new NewPlayerPicker(this.rootElem, this);
		this.playerEditorModal = new PlayerEditorModal();

		const _activePartiesSelector = new EnumPicker<Raid>(raidControls, this.raidSimUI.sim.raid, {
			id: 'raid-picker-size',
			label: 'Raid Size',
			labelTooltip: 'Number of players participating in the sim.',
			values: [
				{ name: '5', value: 1 },
				{ name: '10', value: 2 },
				{ name: '25', value: 5 },
			],
			changedEvent: (raid: Raid) => raid.numActivePartiesChangeEmitter,
			getValue: (raid: Raid) => raid.getNumActiveParties(),
			setValue: (eventID: EventID, raid: Raid, newValue: number) => {
				raid.setNumActiveParties(eventID, newValue);
			},
		});

		const _factionSelector = new EnumPicker<NewPlayerPicker>(raidControls, this.newPlayerPicker, {
			id: 'raid-picker-faction',
			label: 'Default Faction',
			labelTooltip: 'Default faction for newly-created players.',
			values: [
				{ name: 'Alliance', value: Faction.Alliance },
				{ name: 'Horde', value: Faction.Horde },
			],
			changedEvent: (_picker: NewPlayerPicker) => this.raid.sim.factionChangeEmitter,
			getValue: (_picker: NewPlayerPicker) => this.raid.sim.getFaction(),
			setValue: (eventID: EventID, picker: NewPlayerPicker, newValue: Faction) => {
				this.raid.sim.setFaction(eventID, newValue);
			},
		});

		const _phaseSelector = new EnumPicker<NewPlayerPicker>(raidControls, this.newPlayerPicker, {
			id: 'raid-picker-gear',
			label: 'Default Gear',
			labelTooltip: 'Newly-created players will start with approximate BIS gear from this phase.',
			values: [...Array(LATEST_PHASE_WITH_ALL_PRESETS).keys()].map(val => {
				const phase = val + 1;
				return { name: 'Phase ' + phase, value: phase };
			}),
			changedEvent: (_picker: NewPlayerPicker) => this.raid.sim.phaseChangeEmitter,
			getValue: (_picker: NewPlayerPicker) => this.raid.sim.getPhase(),
			setValue: (eventID: EventID, picker: NewPlayerPicker, newValue: number) => {
				this.raid.sim.setPhase(eventID, newValue);
			},
		});

		const partiesContainer = document.createElement('div');
		partiesContainer.classList.add('parties-container');
		this.rootElem.appendChild(partiesContainer);

		this.partyPickers = this.raid.getParties().map((party, i) => new PartyPicker(partiesContainer, party, i, this));

		const updateActiveParties = () => {
			if (this.raidSimUI.sim.raid.getNumActiveParties() == 1) {
				partiesContainer.classList.remove('parties-container-small');
				partiesContainer.classList.remove('parties-container-full');
			} else if (this.raidSimUI.sim.raid.getNumActiveParties() <= 2) {
				partiesContainer.classList.add('parties-container-small');
				partiesContainer.classList.remove('parties-container-full');
			} else {
				partiesContainer.classList.remove('parties-container-small');
				partiesContainer.classList.add('parties-container-full');
			}
			this.partyPickers.forEach(partyPicker => {
				if (partyPicker.index < this.raidSimUI.sim.raid.getNumActiveParties()) {
					partyPicker.rootElem.classList.add('active');
					partyPicker.rootElem.classList.remove('hide');
				} else {
					partyPicker.rootElem.classList.remove('active');
					partyPicker.rootElem.classList.add('hide');
				}
			});
		};
		this.raidSimUI.sim.raid.numActivePartiesChangeEmitter.on(updateActiveParties);
		updateActiveParties();

		this.rootElem.ondragend = _event => {
			// Uncomment to remove player when dropped 'off' the raid.
			//if (this.currentDragPlayerFromIndex != NEW_PLAYER) {
			//	const playerPicker = this.getPlayerPicker(this.currentDragPlayerFromIndex);
			//	playerPicker.setPlayer(null, null, DragType.None);
			//}

			this.clearDragPlayer();
			this.clearDragParty();
		};
	}

	getCurrentFaction(): Faction {
		return this.raid.sim.getFaction();
	}

	getCurrentPhase(): number {
		return this.raid.sim.getPhase();
	}

	getPlayerPicker(raidIndex: number): PlayerPicker {
		return this.partyPickers[Math.floor(raidIndex / MAX_PARTY_SIZE)].playerPickers[raidIndex % MAX_PARTY_SIZE];
	}

	getPlayerPickers(): Array<PlayerPicker> {
		return [...new Array(25).keys()].map(i => this.getPlayerPicker(i));
	}

	setDragPlayer(player: Player<any>, fromIndex: number, type: DragType) {
		this.clearDragPlayer();

		this.currentDragPlayer = player;
		this.currentDragPlayerFromIndex = fromIndex;
		this.currentDragType = type;

		if (fromIndex != NEW_PLAYER) {
			const playerPicker = this.getPlayerPicker(fromIndex);
			playerPicker.rootElem.classList.add('dragfrom');
		}
	}

	clearDragPlayer() {
		if (this.currentDragPlayerFromIndex != NEW_PLAYER) {
			const playerPicker = this.getPlayerPicker(this.currentDragPlayerFromIndex);
			playerPicker.rootElem.classList.remove('dragfrom');
		}

		this.currentDragPlayer = null;
		this.currentDragPlayerFromIndex = NEW_PLAYER;
		this.currentDragType = DragType.New;
	}

	setDragParty(party: PartyPicker) {
		this.currentDragParty = party;
		party.rootElem.classList.add('dragfrom');
	}
	clearDragParty() {
		if (this.currentDragParty) {
			this.currentDragParty.rootElem.classList.remove('dragfrom');
			this.currentDragParty = null;
		}
	}
}

export class PartyPicker extends Component {
	readonly party: Party;
	readonly index: number;
	readonly raidPicker: RaidPicker;
	readonly playerPickers: Array<PlayerPicker>;

	constructor(parent: HTMLElement, party: Party, index: number, raidPicker: RaidPicker) {
		super(parent, 'party-picker-root');
		this.party = party;
		this.index = index;
		this.raidPicker = raidPicker;

		const playersContainerRef = ref<HTMLDivElement>();
		const dpsResultRef = ref<HTMLDivElement>();
		const referenceDeltaRef = ref<HTMLDivElement>();

		this.rootElem.setAttribute('draggable', 'true');
		this.rootElem.replaceChildren(
			<>
				<div className="party-header">
					<label className="party-label form-label">Group {index + 1}</label>
					<div className="party-results">
						<span ref={dpsResultRef} className="party-results-dps"></span>
						<span ref={referenceDeltaRef} className="party-results-reference-delta"></span>
					</div>
				</div>
				<div ref={playersContainerRef} className="players-container"></div>
			</>,
		);

		const playersContainer = playersContainerRef.value!;
		this.playerPickers = [...Array(MAX_PARTY_SIZE).keys()].map(i => new PlayerPicker(playersContainer, this, i));

		const dpsResultElem = dpsResultRef.value!;
		const referenceDeltaElem = referenceDeltaRef.value!;

		this.raidPicker.raidSimUI.referenceChangeEmitter.on(() => {
			const currentData = this.raidPicker.raidSimUI.getCurrentData();
			const referenceData = this.raidPicker.raidSimUI.getReferenceData();

			const partyDps = currentData?.simResult.raidMetrics.parties[this.index]?.dps.avg || 0;
			const referenceDps = referenceData?.simResult.raidMetrics.parties[this.index]?.dps.avg || 0;

			if (partyDps == 0 && referenceDps == 0) {
				dpsResultElem.textContent = '';
				referenceDeltaElem.textContent = '';
				return;
			}

			dpsResultElem.textContent = `${partyDps.toFixed(1)} DPS`;

			if (!referenceData) {
				referenceDeltaElem.textContent = '';
				return;
			}

			formatDeltaTextElem(referenceDeltaElem, referenceDps, partyDps, 1, undefined, undefined, true);
		});

		this.rootElem.ondragstart = event => {
			if (event.target == this.rootElem) {
				event.dataTransfer!.dropEffect = 'move';
				event.dataTransfer!.effectAllowed = 'all';
				this.raidPicker.setDragParty(this);
			}
		};

		let dragEnterCounter = 0;
		this.rootElem.ondragenter = event => {
			event.preventDefault();
			if (!this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter++;
			this.rootElem.classList.add('dragto');
		};
		this.rootElem.ondragleave = event => {
			event.preventDefault();
			if (!this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter--;
			if (dragEnterCounter <= 0) {
				this.rootElem.classList.remove('dragto');
			}
		};
		this.rootElem.ondragover = event => {
			event.preventDefault();
		};
		this.rootElem.ondrop = event => {
			if (!this.raidPicker.currentDragParty) {
				return;
			}

			event.preventDefault();
			dragEnterCounter = 0;
			this.rootElem.classList.remove('dragto');

			const eventID = TypedEvent.nextEventID();
			TypedEvent.freezeAllAndDo(() => {
				const srcPartyPicker = this.raidPicker.currentDragParty!;

				for (let i = 0; i < MAX_PARTY_SIZE; i++) {
					const srcPlayerPicker = srcPartyPicker.playerPickers[i]!;
					const dstPlayerPicker = this.playerPickers[i]!;

					const srcPlayer = srcPlayerPicker.player;
					const dstPlayer = dstPlayerPicker.player;

					srcPlayerPicker.setPlayer(eventID, dstPlayer, DragType.Swap);
					dstPlayerPicker.setPlayer(eventID, srcPlayer, DragType.Swap);
				}
			});

			this.raidPicker.clearDragParty();
		};
	}

	getClosestEmptyIndex() {
		const closestEmptyIndex = this.playerPickers.findIndex(pp => !pp.player);
		return closestEmptyIndex !== -1 ? closestEmptyIndex : null;
	}
}

const EMPTY_PLAYER_NAME = 'Unnamed';

export class PlayerPicker extends Component {
	// Index of this player within its party (0-4).
	readonly index: number;

	// Index of this player within the whole raid (0-24).
	readonly raidIndex: number;

	player: Player<any> | null;

	readonly partyPicker: PartyPicker;
	readonly raidPicker: RaidPicker;

	private labelElem: HTMLElement | null = null;
	private iconElem: HTMLImageElement | null = null;
	private nameElem: HTMLInputElement | null = null;
	private resultsElem: HTMLElement | null = null;
	private dpsResultElem: HTMLElement | null = null;
	private referenceDeltaElem: HTMLElement | null = null;

	private editButton: HTMLButtonElement | null = null;
	private copyButton: HTMLButtonElement | null = null;
	private deleteButton: HTMLButtonElement | null = null;
	// Can be used to remove any events in addEventListener
	// https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#add_an_abortable_listener
	public abortController: AbortController;
	public signal: AbortSignal;

	constructor(parent: HTMLElement, partyPicker: PartyPicker, index: number) {
		super(parent, 'player-picker-root');
		this.abortController = new AbortController();
		this.signal = this.abortController.signal;
		this.index = index;
		this.raidIndex = partyPicker.index * MAX_PARTY_SIZE + index;
		this.player = null;
		this.partyPicker = partyPicker;
		this.raidPicker = partyPicker.raidPicker;

		this.rootElem.classList.add('player');

		this.partyPicker.party.compChangeEmitter.on(eventID => {
			const newPlayer = this.partyPicker.party.getPlayer(this.index);
			if (newPlayer != this.player) this.setPlayer(eventID, newPlayer, DragType.None);
		});

		this.raidPicker.raidSimUI.referenceChangeEmitter.on(() => {
			const currentData = this.raidPicker.raidSimUI.getCurrentData();
			const referenceData = this.raidPicker.raidSimUI.getReferenceData();

			const playerDps = currentData?.simResult.getPlayerWithRaidIndex(this.raidIndex)?.dps.avg || 0;
			const referenceDps = referenceData?.simResult.getPlayerWithRaidIndex(this.raidIndex)?.dps.avg || 0;

			if (this.player) {
				this.resultsElem?.classList.remove('hide');
				(this.dpsResultElem as HTMLElement).textContent = `${playerDps.toFixed(1)} DPS`;

				if (referenceData) formatDeltaTextElem(this.referenceDeltaElem as HTMLElement, referenceDps, playerDps, 1, undefined, undefined, true);
			}
		});

		let dragEnterCounter = 0;
		this.rootElem.ondragenter = event => {
			event.preventDefault();
			if (this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter++;
			this.rootElem.classList.add('dragto');
		};
		this.rootElem.ondragleave = event => {
			event.preventDefault();
			if (this.raidPicker.currentDragParty) {
				return;
			}
			dragEnterCounter--;
			if (dragEnterCounter <= 0) {
				this.rootElem.classList.remove('dragto');
			}
		};
		this.rootElem.ondragover = event => event.preventDefault();
		this.rootElem.ondrop = event => {
			if (this.raidPicker.currentDragParty) {
				return;
			}
			const dropData = event.dataTransfer!.getData('text/plain');

			event.preventDefault();
			dragEnterCounter = 0;
			this.rootElem.classList.remove('dragto');

			const eventID = TypedEvent.nextEventID();
			TypedEvent.freezeAllAndDo(() => {
				if (this.raidPicker.currentDragPlayer == null && dropData.length == 0) {
					return;
				}

				if (this.raidPicker.currentDragPlayerFromIndex == this.raidIndex) {
					this.raidPicker.clearDragPlayer();
					return;
				}

				const dragType = this.raidPicker.currentDragType;

				if (this.raidPicker.currentDragPlayerFromIndex != NEW_PLAYER) {
					const fromPlayerPicker = this.raidPicker.getPlayerPicker(this.raidPicker.currentDragPlayerFromIndex);
					if (dragType == DragType.Swap) {
						fromPlayerPicker.setPlayer(eventID, this.player, dragType);
					} else if (dragType == DragType.Move) {
						fromPlayerPicker.setPlayer(eventID, null, dragType);
					}
				} else if (this.raidPicker.currentDragPlayer == null) {
					// This would be a copy from another window.
					const binary = atob(dropData);
					const bytes = new Uint8Array(binary.length);
					for (let i = 0; i < bytes.length; i++) {
						bytes[i] = binary.charCodeAt(i);
					}
					const playerProto = PlayerProto.fromBinary(bytes);

					const localPlayer = new Player(getPlayerSpecFromPlayer(playerProto), this.raidPicker.raidSimUI.sim);
					localPlayer.fromProto(eventID, playerProto);
					this.raidPicker.currentDragPlayer = localPlayer;
				}

				if (dragType == DragType.Copy) {
					this.setPlayer(eventID, this.raidPicker.currentDragPlayer!.clone(eventID), dragType);
				} else {
					this.setPlayer(eventID, this.raidPicker.currentDragPlayer, dragType);
				}

				this.raidPicker.clearDragPlayer();
			});
		};

		this.update();
	}

	setPlayer(eventID: EventID, newPlayer: Player<any> | null, dragType: DragType) {
		if (newPlayer == this.player) {
			return;
		}
		TypedEvent.freezeAllAndDo(() => {
			const closestEmptySlot = this.partyPicker.getClosestEmptyIndex();
			// If there's empty slots before the current player, we should place the player there.
			const placementIndex = closestEmptySlot && closestEmptySlot < this.index ? closestEmptySlot : this.index;

			this.player = newPlayer;

			if (newPlayer) {
				this.partyPicker.party.setPlayer(eventID, placementIndex, newPlayer);

				if (dragType == DragType.New) {
					applyNewPlayerAssignments(eventID, newPlayer, this.raidPicker.raid);
				}
			} else {
				// Updates the old player's index
				this.partyPicker.party.setPlayer(eventID, placementIndex, newPlayer);
				// If the player left a gap in the grouping, we need to shift the rest of the players up.
				const remainingOptionsToMove = this.partyPicker.playerPickers.slice(this.index, 5).filter(pp => pp.player);
				remainingOptionsToMove.forEach((pp, index) => {
					if (placementIndex < pp.index) this.partyPicker.party.setPlayer(eventID, placementIndex + index, pp.player);
				});
				this.partyPicker.party.compChangeEmitter.emit(eventID);
			}
		});

		this.update();
	}

	private update() {
		if (this.player == null) {
			this.rootElem.className = 'player-picker-root player';
			this.rootElem.innerHTML = '';

			this.labelElem = null;
			this.iconElem = null;
			this.nameElem = null;
			this.resultsElem = null;
			this.dpsResultElem = null;
			this.referenceDeltaElem = null;
		} else {
			const classCssClass = PlayerClasses.getCssClass(this.player.getPlayerClass());

			const labelRef = ref<HTMLDivElement>();
			const iconRef = ref<HTMLImageElement>();
			const nameRef = ref<HTMLInputElement>();
			const resultsRef = ref<HTMLDivElement>();
			const dpsResultRef = ref<HTMLDivElement>();
			const referenceDeltaRef = ref<HTMLDivElement>();
			const editRef = ref<HTMLButtonElement>();
			const copyRef = ref<HTMLButtonElement>();
			const deleteRef = ref<HTMLButtonElement>();

			this.rootElem.className = `player-picker-root player bg-${classCssClass}-dampened`;
			this.rootElem.replaceChildren(
				<>
					<div ref={labelRef} className="player-label">
						<img ref={iconRef} className="player-icon" src={this.player.getSpecIcon()} draggable={true} />
						<div className="player-details">
							<input
								ref={nameRef}
								className={clsx('player-name', `text-${classCssClass}`)}
								type="text"
								value={this.player.getName()}
								spellcheck={false}
								maxLength={15}
							/>
							<div ref={resultsRef} className="player-results hide">
								<span ref={dpsResultRef} className="player-results-dps"></span>
								<span ref={referenceDeltaRef} className="player-results-reference-delta"></span>
							</div>
						</div>
					</div>
					<div className="player-options">
						<button ref={editRef} className="player-edit" dataset={{ tippyContent: 'Click to Edit' }}>
							<i className="fa fa-edit fa-lg"></i>
						</button>
						<button ref={copyRef} className="player-copy link-warning" draggable={true} dataset={{ tippyContent: 'Drag to Copy' }}>
							<i className="fa fa-copy fa-lg"></i>
						</button>
						<button ref={deleteRef} className="player-delete link-danger" dataset={{ tippyContent: 'Click to Delete' }}>
							<i className="fa fa-times fa-lg"></i>
						</button>
					</div>
				</>,
			);

			this.labelElem = labelRef.value!;
			this.iconElem = iconRef.value!;
			this.nameElem = nameRef.value!;
			this.resultsElem = resultsRef.value!;
			this.dpsResultElem = dpsResultRef.value!;
			this.referenceDeltaElem = referenceDeltaRef.value!;

			this.editButton = editRef.value!;
			this.copyButton = copyRef.value!;
			this.deleteButton = deleteRef.value!;

			this.bindPlayerEvents();
		}
	}

	private bindPlayerEvents() {
		const onNameSetHandler = () => {
			this.player?.setName(TypedEvent.nextEventID(), this.nameElem?.value || '');
		};
		this.nameElem?.addEventListener('input', onNameSetHandler, { signal: this.signal });

		const onNameMouseDownHandler = () => {
			this.partyPicker.rootElem.setAttribute('draggable', 'false');
		};
		this.nameElem?.addEventListener('mousedown', onNameMouseDownHandler, { signal: this.signal });

		const onNameMouseUpHandler = () => {
			this.partyPicker.rootElem.setAttribute('draggable', 'true');
		};
		this.nameElem?.addEventListener('mouseup', onNameMouseUpHandler, { signal: this.signal });

		const onNameFocusOutHandler = () => {
			if (this.nameElem && !this.nameElem.value) {
				this.nameElem.value = EMPTY_PLAYER_NAME;
				this.player?.setName(TypedEvent.nextEventID(), this.nameElem.value);
			}
		};
		this.nameElem?.addEventListener('focusout', onNameFocusOutHandler, { signal: this.signal });

		const dragStart = (event: DragEvent, type: DragType) => {
			if (this.player === null) {
				event.preventDefault();
				return;
			}

			event.dataTransfer!.dropEffect = 'move';
			event.dataTransfer!.effectAllowed = 'all';

			if (this.player) {
				const playerDataProto = this.player.toProto(true);
				event.dataTransfer!.setData('text/plain', btoa(String.fromCharCode(...PlayerProto.toBinary(playerDataProto))));
			}

			this.raidPicker.setDragPlayer(this.player, this.raidIndex, type);
		};

		const editTooltip = tippy(this.editButton!);
		const copyTooltip = tippy(this.copyButton!);
		const deleteTooltip = tippy(this.deleteButton!);

		const onIconDragStartHandler = (event: DragEvent) => {
			event.dataTransfer!.setDragImage(this.rootElem, 20, 20);
			dragStart(event, DragType.Swap);
		};
		this.iconElem?.addEventListener('dragstart', onIconDragStartHandler, { signal: this.signal });

		const onEditClickHandler = () => {
			if (this.player) this.raidPicker.playerEditorModal.openEditor(this.player);
		};
		this.editButton?.addEventListener('click', onEditClickHandler, { signal: this.signal });

		const onCopyDragStartHandler = (event: DragEvent) => {
			event.dataTransfer!.setDragImage(this.rootElem, 20, 20);
			dragStart(event, DragType.Copy);
		};
		this.copyButton?.addEventListener('dragstart', onCopyDragStartHandler, { signal: this.signal });

		const onDeleteClickHandler = () => {
			this.setPlayer(TypedEvent.nextEventID(), null, DragType.None);
			this.dispose();
		};
		this.deleteButton?.addEventListener('click', onDeleteClickHandler, { signal: this.signal });

		this.addOnDisposeCallback(() => {
			editTooltip?.destroy();
			copyTooltip?.destroy();
			deleteTooltip?.destroy();
		});
	}
}

class PlayerEditorModal<SpecType extends Spec> extends BaseModal {
	playerEditorRoot: HTMLDivElement;
	constructor() {
		super(document.body, 'player-editor-modal', {
			closeButton: { fixed: true },
			header: false,
			disposeOnClose: false,
		});

		const playerEditorElemRef = ref<HTMLDivElement>();
		const playerEditorElem = (<div ref={playerEditorElemRef} className="player-editor within-raid-sim"></div>) as HTMLDivElement;

		this.rootElem.id = 'playerEditorModal';
		this.body.appendChild(playerEditorElem);

		this.playerEditorRoot = playerEditorElemRef.value!;
	}

	openEditor(player: Player<SpecType>) {
		this.setData(player);
		super.open();
	}

	setData(player: Player<SpecType>) {
		this.playerEditorRoot.innerHTML = '';
		specSimFactories[player.getSpec()]?.(this.playerEditorRoot!, player);
	}
}

class NewPlayerPicker extends Component {
	readonly raidPicker: RaidPicker;

	constructor(parent: HTMLElement, raidPicker: RaidPicker) {
		super(parent, 'new-player-picker-root');
		this.raidPicker = raidPicker;

		getEnumValues(Class).forEach(wowClass => {
			if (wowClass == Class.ClassUnknown) {
				return;
			}

			const matchingPresets = playerPresets.filter(preset => PlayerSpecs.fromProto(preset.spec).classID == wowClass);
			if (matchingPresets.length == 0) {
				return;
			}

			const classPresetsContainerRef = ref<HTMLDivElement>();
			this.rootElem.appendChild(
				<div className={clsx('class-presets-container', `bg-${PlayerClasses.getCssClass(PlayerClasses.fromProto(wowClass as Class))}-dampened`)}>
					{matchingPresets.map(matchingPreset => {
						const playerSpec = PlayerSpecs.fromProto(matchingPreset.spec);
						const presetRef = ref<HTMLButtonElement>();

						const presetButton = (
							<button
								ref={presetRef}
								draggable={true}
								dataset={{ tippyContent: matchingPreset.tooltip ?? PlayerSpecs.getFullSpecName(playerSpec) }}>
								<img className="preset-picker-icon player-icon" src={playerSpec.getIcon('medium')} />
							</button>
						);

						if (presetRef.value) {
							tippy(presetRef.value);

							presetRef.value.ondragstart = event => {
								const eventID = TypedEvent.nextEventID();
								TypedEvent.freezeAllAndDo(() => {
									const dragImage = new Image();
									dragImage.src = matchingPreset.iconUrl ?? playerSpec.getIcon('medium');
									event.dataTransfer!.setDragImage(dragImage, 30, 30);
									event.dataTransfer!.setData('text/plain', '');
									event.dataTransfer!.dropEffect = 'copy';

									const newPlayer = new Player(playerSpec, this.raidPicker.raid.sim);

									newPlayer.applySharedDefaults(eventID);
									newPlayer.setRace(eventID, matchingPreset.defaultFactionRaces[this.raidPicker.getCurrentFaction()]);
									newPlayer.setTalentsString(eventID, matchingPreset.talents.talentsString);
									newPlayer.setGlyphs(eventID, matchingPreset.talents.glyphs || Glyphs.create());
									newPlayer.setSpecOptions(eventID, matchingPreset.specOptions);
									newPlayer.setConsumes(eventID, matchingPreset.consumables);
									newPlayer.setName(eventID, matchingPreset.defaultName ?? playerSpec.friendlyName);
									newPlayer.setProfession1(eventID, matchingPreset.otherDefaults?.profession1 || Profession.Engineering);
									newPlayer.setProfession2(eventID, matchingPreset.otherDefaults?.profession2 || Profession.Jewelcrafting);
									newPlayer.setDistanceFromTarget(eventID, matchingPreset.otherDefaults?.distanceFromTarget || 0);

									// Need to wait because the gear might not be loaded yet.
									this.raidPicker.raid.sim.waitForInit().then(() => {
										const phase = Math.min(this.raidPicker.getCurrentPhase(), LATEST_PHASE_WITH_ALL_PRESETS);
										const gearSet = matchingPreset.defaultGear[this.raidPicker.getCurrentFaction()][phase];
										newPlayer.setGear(eventID, this.raidPicker.raid.sim.db.lookupEquipmentSpec(gearSet));
									});

									this.raidPicker.setDragPlayer(newPlayer, NEW_PLAYER, DragType.New);
								});
							};
						}

						return presetRef.value;
					})}
				</div>,
			);
		});
	}
}

function applyNewPlayerAssignments(eventID: EventID, newPlayer: Player<any>, raid: Raid) {
	if (newPlayer.getPlayerSpec().isTankSpec) {
		const tanks = raid.getTanks();
		const emptyIdx = tanks.findIndex(tank => raid.getPlayerFromUnitReference(tank) == null);
		if (emptyIdx == -1) {
			if (tanks.length < 3) {
				raid.setTanks(eventID, tanks.concat([newPlayer.makeUnitReference()]));
			}
		} else {
			tanks[emptyIdx] = newPlayer.makeUnitReference();
			raid.setTanks(eventID, tanks);
		}
	}

	// Spec-specific assignments. For most cases, default to buffing self.
	if (newPlayer.getSpec() == Spec.SpecBalanceDruid) {
		const newOptions = newPlayer.getSpecOptions() as BalanceDruidOptions;
		newOptions.classOptions!.innervateTarget = newUnitReference(newPlayer.getRaidIndex());
		newPlayer.setSpecOptions(eventID, newOptions);
	} else if (newPlayer.getSpec() == Spec.SpecUnholyDeathKnight) {
		const newOptions = newPlayer.getSpecOptions() as UnholyDeathKnight_Options;
		newOptions.unholyFrenzyTarget = newUnitReference(newPlayer.getRaidIndex());
		newPlayer.setSpecOptions(eventID, newOptions);
	}
}
