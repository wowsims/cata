import tippy from 'tippy.js';

import { Component } from '../core/components/component';
import { IconEnumPicker } from '../core/components/pickers/icon_enum_picker';
import { PlayerClasses } from '../core/player_classes';
import { Paladin } from '../core/player_classes/paladin';
import { PlayerSpec } from '../core/player_spec';
import { Class as ClassProto } from '../core/proto/common';
import { Blessings } from '../core/proto/paladin';
import { BlessingsAssignments } from '../core/proto/ui';
import { ActionId } from '../core/proto_utils/action_id';
import { makeDefaultBlessings } from '../core/proto_utils/utils';
import { EventID, TypedEvent } from '../core/typed_event';
import { implementedSpecs } from './presets';
import { RaidSimUI } from './raid_sim_ui';

const MAX_PALADINS = 4;

export class BlessingsPicker extends Component {
	readonly simUI: RaidSimUI;
	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly pickers: Array<Array<IconEnumPicker<this, Blessings>>> = [];

	private assignments: BlessingsAssignments;

	constructor(parentElem: HTMLElement, raidSimUI: RaidSimUI) {
		super(parentElem, 'blessings-picker-root');
		this.simUI = raidSimUI;
		this.assignments = BlessingsAssignments.clone(makeDefaultBlessings(4));

		const playerSpecs = PlayerClasses.naturalOrder
			.map(playerClass => Object.values(playerClass.specs))
			.flat()
			.filter(spec => implementedSpecs.includes(spec.specID));
		const paladinIndexes = [...Array(MAX_PALADINS).keys()];

		playerSpecs.map(playerSpec => {
			const row = document.createElement('div');
			row.classList.add('blessings-picker-row');
			this.rootElem.appendChild(row);

			row.append(this.buildSpecIcon(playerSpec));

			const container = document.createElement('div');
			container.classList.add('blessings-picker-container');
			row.appendChild(container);

			paladinIndexes.forEach(paladinIdx => {
				if (!this.pickers[paladinIdx]) this.pickers.push([]);

				const blessingPicker = new IconEnumPicker(container, this, {
					extraCssClasses: ['blessing-picker'],
					numColumns: 1,
					values: [
						{ color: Paladin.hexColor, value: Blessings.BlessingUnknown },
						{ actionId: ActionId.fromSpellId(20217), value: Blessings.BlessingOfKings },
						{ actionId: ActionId.fromSpellId(19740), value: Blessings.BlessingOfMight },
					],
					equals: (a: Blessings, b: Blessings) => a == b,
					zeroValue: Blessings.BlessingUnknown,
					enableWhen: (_picker: BlessingsPicker) => {
						const numPaladins = Math.min(this.simUI.getClassCount(ClassProto.ClassPaladin), MAX_PALADINS);
						return paladinIdx < numPaladins;
					},
					changedEvent: (picker: BlessingsPicker) => picker.changeEmitter,
					getValue: (picker: BlessingsPicker) => picker.assignments.paladins[paladinIdx]?.blessings[playerSpec.specID] || Blessings.BlessingUnknown,
					setValue: (eventID: EventID, picker: BlessingsPicker, newValue: number) => {
						const currentValue = picker.assignments.paladins[paladinIdx].blessings[playerSpec.specID];
						if (currentValue != newValue) {
							picker.assignments.paladins[paladinIdx].blessings[playerSpec.specID] = newValue;
							this.changeEmitter.emit(eventID);
						}
					},
				});

				this.pickers[paladinIdx].push(blessingPicker);
			});

			return row;
		});

		this.updatePickers();
		this.simUI.compChangeEmitter.on(_eventID => this.updatePickers());
	}

	private updatePickers() {
		for (let i = 0; i < MAX_PALADINS; i++) {
			this.pickers[i].forEach(picker => picker.update());
		}
	}

	private buildSpecIcon(spec: PlayerSpec<any>): HTMLElement {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="blessings-picker-spec">
				<img
					src="${spec.getIcon('medium')}"
					class="blessings-spec-icon"
				/>
			</div>
		`;

		const icon = fragment.querySelector('.blessings-spec-icon') as HTMLElement;
		tippy(icon, { content: spec.friendlyName });

		return fragment.children[0] as HTMLElement;
	}

	getAssignments(): BlessingsAssignments {
		// Defensive copy.
		return BlessingsAssignments.clone(this.assignments);
	}

	setAssignments(eventID: EventID, newAssignments: BlessingsAssignments) {
		this.assignments = BlessingsAssignments.clone(newAssignments);
		this.changeEmitter.emit(eventID);
	}
}
