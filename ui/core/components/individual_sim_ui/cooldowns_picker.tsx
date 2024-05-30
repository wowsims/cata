import { Icon, Link } from '@wowsims/ui';
import clsx from 'clsx';
import tippy from 'tippy.js';

import { Player } from '../../player.js';
import { ActionID as ActionIdProto, Cooldown } from '../../proto/common.js';
import { ActionId } from '../../proto_utils/action_id.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { existsInDOM } from '../../utils';
import { Component } from '../component.js';
import { IconEnumPicker, IconEnumValueConfig } from '../icon_enum_picker.jsx';
import { NumberListPicker } from '../number_list_picker.js';

export class CooldownsPicker extends Component {
	readonly player: Player<any>;

	private cooldownPickers: Array<HTMLElement>;

	constructor(parentElem: HTMLElement, player: Player<any>) {
		super(parentElem, 'cooldowns-picker-root');
		this.player = player;
		this.cooldownPickers = [];

		const events = TypedEvent.onAny([this.player.rotationChangeEmitter, this.player.sim.unitMetadataEmitter]).on(() => {
			if (!existsInDOM(this.rootElem)) {
				this.dispose();
				return;
			}
			this.update();
		});
		this.addOnDisposeCallback(() => {
			this.rootElem.remove();
			this.cooldownPickers.forEach(picker => picker.remove());
			events.dispose();
		});
		this.update();
	}

	private update() {
		this.rootElem.innerHTML = '';
		const cooldowns = this.player.getSimpleCooldowns().cooldowns;

		this.cooldownPickers = [];
		for (let i = 0; i < cooldowns.length + 1; i++) {
			const cooldown = cooldowns[i];

			const row = (<div className={clsx('cooldown-picker', i === cooldowns.length && 'add-cooldown-picker')} />) as HTMLElement;
			this.rootElem.appendChild(row);

			this.makeActionPicker(row, i);

			const label = <label className="cooldown-picker-label form-label"></label>;
			if (cooldown && cooldown.id) {
				ActionId.fromProto(cooldown.id)
					.fill(this.player.getRaidIndex())
					.then(filledId => (label.textContent = filledId.name));
			}
			row.appendChild(label);

			this.makeTimingsPicker(row, i);

			const deleteButton = (
				<Link variant="danger" className="delete-cooldown" role="button">
					<Icon icon="times" size="xl" />
				</Link>
			);
			const deleteButtonTooltip = tippy(deleteButton, { content: 'Delete Cooldown' });
			deleteButton.addEventListener('click', () => {
				const newCooldowns = this.player.getSimpleCooldowns();
				newCooldowns.cooldowns.splice(i, 1);
				this.player.setSimpleCooldowns(TypedEvent.nextEventID(), newCooldowns);
				deleteButtonTooltip.hide();
			});
			row.appendChild(deleteButton);

			this.cooldownPickers.push(row);
		}
	}

	private makeActionPicker(parentElem: HTMLElement, cooldownIndex: number): IconEnumPicker<Player<any>, ActionIdProto> {
		const availableCooldowns = this.player
			.getMetadata()
			.getSpells()
			.filter(spell => spell.data.isMajorCooldown)
			.map(spell => spell.id);

		const actionPicker = new IconEnumPicker<Player<any>, ActionIdProto>(parentElem, this.player, {
			extraCssClasses: ['cooldown-action-picker'],
			numColumns: 3,
			values: ([{ color: '#grey', value: ActionIdProto.create() }] as Array<IconEnumValueConfig<Player<any>, ActionIdProto>>).concat(
				availableCooldowns.map(cooldownAction => {
					return { actionId: cooldownAction, value: cooldownAction.toProto() };
				}),
			),
			equals: (a: ActionIdProto, b: ActionIdProto) => ActionIdProto.equals(a, b),
			zeroValue: ActionIdProto.create(),
			backupIconUrl: (value: ActionIdProto) => ActionId.fromProto(value),
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => player.getSimpleCooldowns().cooldowns[cooldownIndex]?.id || ActionIdProto.create(),
			setValue: (eventID: EventID, player: Player<any>, newValue: ActionIdProto) => {
				const newCooldowns = player.getSimpleCooldowns();

				while (newCooldowns.cooldowns.length < cooldownIndex) {
					newCooldowns.cooldowns.push(Cooldown.create());
				}
				newCooldowns.cooldowns[cooldownIndex] = Cooldown.create({
					id: newValue,
					timings: [],
				});

				player.setSimpleCooldowns(eventID, newCooldowns);
			},
		});
		return actionPicker;
	}

	private makeTimingsPicker(parentElem: HTMLElement, cooldownIndex: number): NumberListPicker<Player<any>> {
		const actionPicker = new NumberListPicker(parentElem, this.player, {
			id: `cooldown-timings-${cooldownIndex}`,
			extraCssClasses: ['cooldown-timings-picker'],
			placeholder: '20, 40, ...',
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (player: Player<any>) => {
				return player.getSimpleCooldowns().cooldowns[cooldownIndex]?.timings || [];
			},
			setValue: (eventID: EventID, player: Player<any>, newValue: Array<number>) => {
				const newCooldowns = player.getSimpleCooldowns();
				newCooldowns.cooldowns[cooldownIndex].timings = newValue;
				player.setSimpleCooldowns(eventID, newCooldowns);
			},
			enableWhen: (player: Player<any>) => {
				const curCooldown = player.getSimpleCooldowns().cooldowns[cooldownIndex];
				return curCooldown && !ActionIdProto.equals(curCooldown.id, ActionIdProto.create());
			},
		});
		return actionPicker;
	}
}
