import { IndividualSimUI } from '../../../individual_sim_ui';
import { EquippedItem } from '../../../proto_utils/equipped_item';
import { ContentBlock } from '../../content_block';
import Toast from '../../toast';
import { BulkTab } from '../bulk_tab';
import BulkItemPicker from './bulk_item_picker';
import { BulkSimItemSlot, bulkSimSlotNames } from './utils';

export default class BulkItemPickerGroup extends ContentBlock {
	readonly simUI: IndividualSimUI<any>;
	readonly bulkUI: BulkTab;
	readonly slot: BulkSimItemSlot;

	readonly pickers: Map<number, BulkItemPicker> = new Map();

	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, bulkUI: BulkTab, slot: BulkSimItemSlot) {
		const slotName = bulkSimSlotNames.get(slot)!;
		super(parent, 'bulk-item-picker-group-root', { header: { title: slotName } });
		this.rootElem.classList.add(`gear-group-${slotName.split(' ').join('-')}`);
		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.slot = slot;

		this.addEmptyElement();
	}

	add(idx: number, item: EquippedItem) {
		if (!this.pickers.size) this.bodyElement.replaceChildren();

		if (this.pickers.has(idx)) {
			const picker = this.pickers.get(idx);
			picker!.dispose();
			this.pickers.delete(idx);
		}

		this.pickers.set(idx, new BulkItemPicker(this.bodyElement, this.simUI, this.bulkUI, item, idx));
	}

	update(idx: number, newItem: EquippedItem) {
		const picker = this.pickers.get(idx);
		if (!picker) {
			new Toast({
				variant: 'error',
				body: 'Failed to update item, please report this issue.',
			});
			return;
		}

		picker.setItem(newItem);
	}

	remove(idx: number) {
		const picker = this.pickers.get(idx);
		if (!picker) {
			new Toast({
				variant: 'error',
				body: 'Failed to remove item, please report this issue.',
			});
			return;
		}

		picker.dispose();
		this.pickers.delete(idx);

		if (!this.pickers.size) this.addEmptyElement();
	}

	private addEmptyElement() {
		this.bodyElement.appendChild(<span>No items selected.</span>);
	}
}
