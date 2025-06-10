import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Player } from '../player';
import { ItemSlot, ItemSpec, ItemSwap, Spec } from '../proto/common';
import { EquippedItem } from '../proto_utils/equipped_item';
import { ItemSwapGear } from '../proto_utils/gear';
import { Stats } from '../proto_utils/stats';
import { SimUI } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event';
import { Component } from './component';
import IconItemSwapPicker from './gear_picker/icon_item_swap_picker';
import { Input } from './input';
import { BooleanPicker } from './pickers/boolean_picker';

export class ItemSwapSettings {
	private readonly player: Player<any>;
	readonly changeEmitter = new TypedEvent<void>('PlayerItemSwap');

	private enableItemSwap = false;
	private gear = new ItemSwapGear({});
	private bonusStats = new Stats();

	constructor(player: Player<any>) {
		this.player = player;
	}

	setItemSwapSettings(eventID: EventID, enableItemSwap: boolean, gear: ItemSwapGear, bonusStats?: Stats) {
		this.enableItemSwap = enableItemSwap;
		this.gear = gear;
		this.bonusStats = bonusStats || new Stats();

		this.changeEmitter.emit(eventID);
	}

	setBonusStats(eventID: EventID, stats: Stats) {
		this.bonusStats = stats;
		this.changeEmitter.emit(eventID);
	}

	getBonusStats() {
		return this.bonusStats;
	}

	getEnableItemSwap(): boolean {
		return this.enableItemSwap;
	}

	setEnableItemSwap(eventID: EventID, newEnableItemSwap: boolean) {
		if (newEnableItemSwap == this.enableItemSwap) return;

		this.enableItemSwap = newEnableItemSwap;
		this.changeEmitter.emit(eventID);
	}

	equipItem(eventID: EventID, slot: ItemSlot, newItem: EquippedItem | null) {
		this.setGear(eventID, this.gear.withEquippedItem(slot, newItem, this.player.canDualWield2H()));
	}

	getItem(slot: ItemSlot): EquippedItem | null {
		return this.gear.getEquippedItem(slot);
	}

	getGear(): ItemSwapGear {
		return this.gear;
	}

	setGear(eventID: EventID, newItemSwapGear: ItemSwapGear) {
		if (newItemSwapGear.equals(this.gear)) return;

		this.gear = newItemSwapGear;
		this.changeEmitter.emit(eventID);
	}

	toProto(): ItemSwap {
		return ItemSwap.create({
			prepullBonusStats: this.bonusStats.toProto(),
			items: this.gear.asArray().map(ei => (ei ? ei.asSpec() : ItemSpec.create())),
		});
	}
}

export interface ItemSwapPickerConfig {
	itemSlots: Array<ItemSlot>;
	note?: string;
}

export class ItemSwapPicker<SpecType extends Spec> extends Component {
	private readonly itemSlots: Array<ItemSlot>;

	constructor(parentElem: HTMLElement, simUI: SimUI, player: Player<SpecType>, config: ItemSwapPickerConfig) {
		super(parentElem, 'item-swap-picker-root');
		this.itemSlots = config.itemSlots;

		new BooleanPicker(this.rootElem, player, {
			id: 'enable-item-swap',
			reverse: true,
			label: 'Enable Item Swapping',
			labelTooltip: 'Allows configuring an Item Swap Set which is used with the <b>Item Swap</b> APL action.',
			extraCssClasses: ['input-inline'],
			getValue: (player: Player<SpecType>) => player.itemSwapSettings.getEnableItemSwap(),
			setValue(eventID: EventID, player: Player<SpecType>, newValue: boolean) {
				player.itemSwapSettings.setEnableItemSwap(eventID, newValue);
			},
			changedEvent: (player: Player<SpecType>) => player.itemSwapSettings.changeEmitter,
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
			if (!player.itemSwapSettings.getEnableItemSwap()) {
				swapPickerContainerRef.value?.classList.add('hide');
				noteRef.value?.classList.add('hide');
			} else {
				swapPickerContainerRef.value?.classList.remove('hide');
				noteRef.value?.classList.remove('hide');
			}
		};
		player.itemSwapSettings.changeEmitter.on(toggleEnabled);
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
		let newIsg = player.itemSwapSettings.getGear();

		this.itemSlots.forEach(slot => {
			const gearItem = player.getGear().getEquippedItem(slot);
			const swapItem = player.itemSwapSettings.getGear().getEquippedItem(slot);

			newGear = newGear.withEquippedItem(slot, swapItem, player.canDualWield2H());
			newIsg = newIsg.withEquippedItem(slot, gearItem, player.canDualWield2H());
		});

		TypedEvent.freezeAllAndDo(() => {
			player.setGear(eventID, newGear);
			player.itemSwapSettings.setGear(eventID, newIsg);
		});
	}
}
