import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../../individual_sim_ui';
import { EquippedItem } from '../../../proto_utils/equipped_item';
import { getEligibleItemSlots } from '../../../proto_utils/utils';
import { TypedEvent } from '../../../typed_event';
import { Component } from '../../component';
import { ItemRenderer } from '../../gear_picker/gear_picker';
import { GearData } from '../../gear_picker/item_list';
import { SelectorModalTabs } from '../../gear_picker/selector_modal';
import { BulkTab } from '../bulk_tab';

export class BulkItemPicker extends Component {
	private readonly itemElem: ItemRenderer;
	readonly simUI: IndividualSimUI<any>;
	readonly bulkUI: BulkTab;

	protected item: EquippedItem;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, bulkUI: BulkTab, item: EquippedItem) {
		super(parent, 'bulk-item-picker');
		this.simUI = simUI;
		this.bulkUI = bulkUI;
		this.item = item;
		this.itemElem = new ItemRenderer(parent, this.rootElem, simUI.player);

		const removeBtn = ref<HTMLButtonElement>();
		this.itemElem.rootElem.appendChild(
			<button className="remove-batch-item-btn" ref={removeBtn}>
				<i className="fas fa-times" />
			</button>,
		);
		const removeItem = () => this.bulkUI.removeItem(this.item.asSpec());
		removeBtn.value!.addEventListener('click', removeItem);
		this.addOnDisposeCallback(() => removeBtn.value!.removeEventListener('click', removeItem));

		this.simUI.sim.waitForInit().then(() => {
			this.setItem(item);
			const slot = getEligibleItemSlots(this.item.item)[0];
			const eligibleEnchants = this.simUI.sim.db.getEnchants(slot);
			const eligibleReforges = this.item?.item ? this.simUI.player.getAvailableReforgings(this.item.getWithRandomSuffixStats()) : [];
			const eligibleRandomSuffixes = this.item.item.randomSuffixOptions;

			const openEnchantGemSelector = (event: Event) => {
				event.preventDefault();

				if (!!eligibleEnchants.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.Enchants, this.createGearData());
				} else if (!!eligibleRandomSuffixes.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.RandomSuffixes, this.createGearData());
				} else if (!!eligibleReforges.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.Reforging, this.createGearData());
				} else if (!!this.item._gems.length) {
					this.bulkUI.selectorModal.openTab(slot, SelectorModalTabs.Gem1, this.createGearData());
				}
			};

			this.itemElem.iconElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.nameElem.addEventListener('click', openEnchantGemSelector);
			this.itemElem.enchantElem.addEventListener('click', openEnchantGemSelector);
		});
	}

	setItem(newItem: EquippedItem | null) {
		this.itemElem.clear();
		if (!!newItem) {
			this.itemElem.update(newItem);
			this.item = newItem;
		} else {
			this.itemElem.rootElem.style.opacity = '30%';
			this.itemElem.iconElem.style.backgroundImage = `url('/cata/assets/item_slots/empty.jpg')`;
			this.itemElem.nameElem.textContent = 'Add new item (not implemented)';
			this.itemElem.rootElem.style.alignItems = 'center';
		}
	}

	private createGearData(): GearData {
		const changeEvent = new TypedEvent<void>();
		return {
			equipItem: (_, equippedItem: EquippedItem | null) => {
				if (equippedItem) {
					this.bulkUI.removeItem(this.item.asSpec());
					this.bulkUI.addItem(equippedItem.asSpec());
					this.item = equippedItem;
					changeEvent.emit(TypedEvent.nextEventID());
				}
			},
			getEquippedItem: () => this.item,
			changeEvent: changeEvent,
		};
	}
}
