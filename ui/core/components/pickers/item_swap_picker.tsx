import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Player } from '../../player.js';
import { ItemSlot, Spec } from '../../proto/common.js';
import { SimUI } from '../../sim_ui.js';
import { EventID, TypedEvent } from '../../typed_event.js';
import { Component } from '../component.js';
import IconItemSwapPicker from '../gear_picker/icon_item_swap_picker.js';
import { Input } from '../input.js';
import { BooleanPicker } from './boolean_picker.js';

export interface ItemSwapConfig {
	itemSlots: Array<ItemSlot>;
	note?: string;
}

export class ItemSwapPicker<SpecType extends Spec> extends Component {
	private readonly itemSlots: Array<ItemSlot>;

	constructor(parentElem: HTMLElement, simUI: SimUI, player: Player<SpecType>, config: ItemSwapConfig) {
		super(parentElem, 'item-swap-picker-root');
		this.itemSlots = config.itemSlots;

		new BooleanPicker(this.rootElem, player, {
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

		const swapPickerContainerRef = ref<HTMLDivElement>();
		const swapButtonRef = ref<HTMLButtonElement>();
		const noteRef = ref<HTMLParagraphElement>();
		const itemSwapContainer = Input.newGroupContainer('icon-group');
		this.rootElem.appendChild(
			<>
				<div ref={swapPickerContainerRef} className="input-root input-inline input-item-swap-container">
					<label className="form-label">Item Swap</label>
					<button ref={swapButtonRef} className="gear-swap-icon">
						<i className="fas fa-arrows-rotate me-1"></i>
					</button>
					{itemSwapContainer}
				</div>
				{config.note && (
					<p ref={noteRef} className="form-text">
						{config.note}
					</p>
				)}
			</>,
		);

		const toggleEnabled = () => {
			if (!player.getEnableItemSwap()) {
				swapPickerContainerRef.value?.classList.add('hide');
				noteRef.value?.classList.add('hide');
			} else {
				swapPickerContainerRef.value?.classList.remove('hide');
				noteRef.value?.classList.remove('hide');
			}
		};
		player.itemSwapChangeEmitter.on(toggleEnabled);
		toggleEnabled();

		if (swapButtonRef.value) {
			swapButtonRef.value.addEventListener('click', _event => this.swapWithGear(TypedEvent.nextEventID(), player));
			tippy(swapButtonRef.value, {
				content: 'Swap with equipped items',
			});
		}

		const tmpContainer = (<></>) as HTMLElement;
		this.itemSlots.forEach(itemSlot => {
			new IconItemSwapPicker(tmpContainer, simUI, player, itemSlot);
		});

		itemSwapContainer.appendChild(tmpContainer);
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
