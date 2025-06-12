import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../../individual_sim_ui';
import { GemColor, ItemSlot } from '../../../proto/common';
import { UIGem } from '../../../proto/ui';
import { ActionId } from '../../../proto_utils/action_id';
import { EquippedItem } from '../../../proto_utils/equipped_item';
import { Stats } from '../../../proto_utils/stats';
import { EventID, TypedEvent } from '../../../typed_event';
import { noop } from '../../../utils';
import { BaseModal } from '../../base_modal';
import ItemList, { ItemData } from '../../gear_picker/item_list';
import { SelectorModalTabs } from '../../gear_picker/selector_modal';

export default class GemSelectorModal extends BaseModal {
	private readonly simUI: IndividualSimUI<any>;

	private readonly contentElem: HTMLElement;
	private ilist: ItemList<UIGem> | null;
	private socketColor: GemColor;
	private onSelect: (itemData: ItemData<UIGem>) => void;
	private onRemove: () => void;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, socketColor: GemColor, onSelect: (itemData: ItemData<UIGem>) => void, onRemove: () => void) {
		super(parent, 'selector-modal', { disposeOnClose: false, size: 'xl' });

		this.simUI = simUI;
		this.onSelect = onSelect;
		this.onRemove = onRemove;
		this.socketColor = socketColor;
		this.ilist = null;

		window.scrollTo({ top: 0 });

		this.header!.insertAdjacentElement('afterbegin', <h6 className="selector-modal-title mb-3">Choose Default Gem</h6>);
		const contentRef = ref<HTMLDivElement>();
		this.body.appendChild(<div ref={contentRef} className="tab-content selector-modal-tab-content"></div>);
		this.contentElem = contentRef.value!;
	}

	show() {
		// construct item list the first time its opened.
		// This makes startup faster and also means we are sure to have item database loaded.
		if (!this.ilist) {
			this.ilist = new ItemList<UIGem>(
				'bulk-tab-gem-selector',
				this.contentElem,
				this.simUI,
				ItemSlot.ItemSlotHead,
				SelectorModalTabs.Gem1,
				this.simUI.player,
				SelectorModalTabs.Gem1,
				{
					equipItem: (_eventID: EventID, _equippedItem: EquippedItem | null) => {
						return;
					},
					getEquippedItem: () => null,
					changeEvent: new TypedEvent(), // FIXME
				},
				this.simUI.player.getGems(this.socketColor).map((gem: UIGem) => {
					return {
						item: gem,
						id: gem.id,
						actionId: ActionId.fromItemId(gem.id),
						name: gem.name,
						quality: gem.quality,
						phase: gem.phase,
						nameDescription: '',
						baseEP: this.simUI.player.computeStatsEP(new Stats(gem.stats)),
						ignoreEPFilter: true,
						onEquip: noop,
					};
				}),
				this.socketColor,
				gem => {
					return this.simUI.player.computeGemEP(gem);
				},
				() => {
					return null;
				},
				this.onRemove,
				this.onSelect,
			);

			this.ilist.sizeRefresh();

			const applyFilter = () => this.ilist?.applyFilters();

			const phaseChangeEvent = this.simUI.sim.phaseChangeEmitter.on(applyFilter);
			const filtersChangeChangeEvent = this.simUI.sim.filtersChangeEmitter.on(applyFilter);

			this.addOnDisposeCallback(() => {
				phaseChangeEvent.dispose();
				filtersChangeChangeEvent.dispose();
				this.ilist?.dispose();
			});
		}

		this.open();
	}
}
