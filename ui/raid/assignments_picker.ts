import { Component } from '../core/components/component';
import { UnitReferencePicker } from '../core/components/pickers/raid_target_picker';
import { Player } from '../core/player';
import { Class, Spec, UnitReference } from '../core/proto/common';
import { DeathKnightTalents } from '../core/proto/death_knight';
import { PriestTalents } from '../core/proto/priest';
import { emptyUnitReference, RogueSpecs } from '../core/proto_utils/utils';
import { EventID, TypedEvent } from '../core/typed_event';
import { RaidSimUI } from './raid_sim_ui';

export class AssignmentsPicker extends Component {
	readonly raidSimUI: RaidSimUI;
	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly innervatesPicker: InnervatesPicker;
	private readonly tricksOfTheTradesPicker: TricksOfTheTradesPicker;
	private readonly unholyFrenzyPicker: UnholyFrenzyPicker;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'assignments-picker-root');
		this.raidSimUI = raidSimUI;

		this.innervatesPicker = new InnervatesPicker(this.rootElem, raidSimUI);
		this.tricksOfTheTradesPicker = new TricksOfTheTradesPicker(this.rootElem, raidSimUI);
		this.unholyFrenzyPicker = new UnholyFrenzyPicker(this.rootElem, raidSimUI);
	}
}

interface AssignmentTargetPicker {
	player: Player<any>;
	targetPicker: UnitReferencePicker<Player<any>>;
	targetPlayer: Player<any> | null;
}

abstract class AssignedBuffPicker extends Component {
	readonly raidSimUI: RaidSimUI;
	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly playersContainer: HTMLElement;

	private targetPickers: Array<AssignmentTargetPicker>;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'assigned-buff-picker-root');
		this.raidSimUI = raidSimUI;
		this.targetPickers = [];

		this.playersContainer = document.createElement('div');
		this.playersContainer.classList.add('assigned-buff-container');
		this.rootElem.appendChild(this.playersContainer);

		this.raidSimUI.changeEmitter.on(_eventID => this.update());
		this.update();
	}

	private update() {
		this.playersContainer.innerHTML = `
			<label class="assignmented-buff-label form-label">${this.getTitle()}</label>
		`;

		const sourcePlayers = this.getSourcePlayers();
		if (sourcePlayers.length == 0) this.rootElem.classList.add('hide');
		else this.rootElem.classList.remove('hide');

		this.targetPickers = sourcePlayers.map((sourcePlayer, _sourcePlayerIndex) => {
			const row = document.createElement('div');
			row.classList.add('assigned-buff-player', 'input-inline');
			this.playersContainer.appendChild(row);

			const sourceElem = document.createElement('div');
			sourceElem.classList.add('raid-target-picker-root');
			sourceElem.appendChild(UnitReferencePicker.makeOptionElem({ player: sourcePlayer, isDropdown: false }));
			row.appendChild(sourceElem);

			const arrow = document.createElement('i');
			arrow.classList.add('assigned-buff-arrow', 'fa', 'fa-arrow-right');
			row.appendChild(arrow);

			const raidTargetPicker: UnitReferencePicker<Player<any>> | null = new UnitReferencePicker<Player<any>>(row, this.raidSimUI.sim.raid, sourcePlayer, {
				extraCssClasses: ['assigned-buff-target-picker'],
				noTargetLabel: 'Unassigned',
				compChangeEmitter: this.raidSimUI.sim.raid.compChangeEmitter,

				changedEvent: (player: Player<any>) => player.specOptionsChangeEmitter,
				getValue: (player: Player<any>) => this.getPlayerValue(player),
				setValue: (eventID: EventID, player: Player<any>, newValue: UnitReference) => this.setPlayerValue(eventID, player, newValue),
			});

			const targetPickerData = {
				player: sourcePlayer,
				targetPicker: raidTargetPicker!,
				targetPlayer: this.raidSimUI.sim.raid.getPlayerFromUnitReference(raidTargetPicker!.getInputValue()),
			};

			raidTargetPicker!.changeEmitter.on(_eventID => {
				targetPickerData.targetPlayer = this.raidSimUI.sim.raid.getPlayerFromUnitReference(raidTargetPicker!.getInputValue());
			});

			return targetPickerData;
		});
	}

	abstract getTitle(): string;
	abstract getSourcePlayers(): Array<Player<any>>;

	abstract getPlayerValue(player: Player<any>): UnitReference;
	abstract setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference): void;
}

class InnervatesPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Innervate';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isClass(Class.ClassDruid));
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<Spec.SpecBalanceDruid>).getSpecOptions().classOptions?.innervateTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<Spec.SpecBalanceDruid>).getSpecOptions();
		newOptions.classOptions!.innervateTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}

class TricksOfTheTradesPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Tricks of the Trade';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isClass(Class.ClassRogue));
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<RogueSpecs>).getSpecOptions().classOptions!.tricksOfTheTradeTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<RogueSpecs>).getSpecOptions();
		newOptions.classOptions!.tricksOfTheTradeTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}

class UnholyFrenzyPicker extends AssignedBuffPicker {
	getTitle(): string {
		return 'Unholy Frenzy';
	}

	getSourcePlayers(): Array<Player<any>> {
		return this.raidSimUI.getActivePlayers().filter(player => player.isSpec(Spec.SpecUnholyDeathKnight) && false);
	}

	getPlayerValue(player: Player<any>): UnitReference {
		return (player as Player<Spec.SpecUnholyDeathKnight>).getSpecOptions().unholyFrenzyTarget || emptyUnitReference();
	}

	setPlayerValue(eventID: EventID, player: Player<any>, newValue: UnitReference) {
		const newOptions = (player as Player<Spec.SpecUnholyDeathKnight>).getSpecOptions();
		newOptions.unholyFrenzyTarget = newValue;
		player.setSpecOptions(eventID, newOptions);
	}
}

