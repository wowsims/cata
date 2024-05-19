import { Icon, Link } from '@wowsims/ui';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Player } from '../player';
import { ItemSlot, Spec } from '../proto/common';
import { SimUI } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event';
import { BooleanPicker } from './boolean_picker';
import { Component } from './component';
import { IconItemSwapPicker } from './gear_picker/gear_picker';
import { PickerGroup } from './input';

export interface ItemSwapConfig {
	itemSlots: Array<ItemSlot>;
	note?: string;
}

export class ItemSwapPicker<SpecType extends Spec> extends Component {
	private readonly itemSlots: Array<ItemSlot>;
	private readonly enableItemSwapPicker: BooleanPicker<Player<SpecType>>;

	constructor(parentElem: HTMLElement, simUI: SimUI, player: Player<SpecType>, config: ItemSwapConfig) {
		super(parentElem, 'item-swap-picker-root');
		this.itemSlots = config.itemSlots;

		this.enableItemSwapPicker = new BooleanPicker(this.rootElem, player, {
			id: 'enable-item-swap',
			reverse: true,
			label: 'Enable Item Swapping',
			labelTooltip: 'Allows configuring an Item Swap Set which is used with the <b>Item Swap</b> APL action.',
			extraCssClasses: ['input-inline'],
			getValue: (player: Player<SpecType>) => player.getEnableItemSwap(),
			setValue(eventID: EventID, player: Player<SpecType>, newValue: boolean) {
				player.setEnableItemSwap(eventID, newValue);
			},
			changedEvent: (player: Player<SpecType>) => player.itemSwapChangeEmitter,
		});

		const swapButtonRef = ref<HTMLAnchorElement>();
		const itemSwapContainerRef = ref<HTMLDivElement>();
		const swapPickerContainer = (
			<div className="input-root input-inline">
				<label className="form-label">Item Swap</label>
				<PickerGroup ref={itemSwapContainerRef} className="icon-group">
					<Link ref={swapButtonRef} as="button" className="gear-swap-icon">
						<Icon icon="arrows-rotate" className="me-1" />
					</Link>
				</PickerGroup>
			</div>
		);
		this.rootElem.appendChild(swapPickerContainer);

		let noteElem: Element;
		if (config.note) {
			noteElem = this.rootElem.appendChild(<p className="form-text">{config.note}</p>);
		}

		const toggleEnabled = () => {
			if (!player.getEnableItemSwap()) {
				swapPickerContainer.classList.add('hide');
				noteElem?.classList.add('hide');
			} else {
				swapPickerContainer.classList.remove('hide');
				noteElem?.classList.remove('hide');
			}
		};
		player.itemSwapChangeEmitter.on(toggleEnabled);
		toggleEnabled();

		swapButtonRef.value!.addEventListener('click', () => this.swapWithGear(TypedEvent.nextEventID(), player));

		tippy(swapButtonRef.value!, {
			content: 'Swap with equipped items',
		});

		this.itemSlots.forEach(itemSlot => new IconItemSwapPicker(itemSwapContainerRef.value!, simUI, player, itemSlot));
	}

	swapWithGear(eventID: EventID, player: Player<SpecType>) {
		let newGear = player.getGear();
		let newIsg = player.getItemSwapGear();

		this.itemSlots.forEach(slot => {
			const gearItem = player.getGear().getEquippedItem(slot);
			const swapItem = player.getItemSwapGear().getEquippedItem(slot);

			newGear = newGear.withEquippedItem(slot, swapItem, player.canDualWield2H());
			newIsg = newIsg.withEquippedItem(slot, gearItem, player.canDualWield2H());
		});

		TypedEvent.freezeAllAndDo(() => {
			player.setGear(eventID, newGear);
			player.setItemSwapGear(eventID, newIsg);
		});
	}
}
