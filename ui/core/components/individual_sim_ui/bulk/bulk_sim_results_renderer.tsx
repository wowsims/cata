import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../../individual_sim_ui';
import { BulkComboResult, ItemSpecWithSlot } from '../../../proto/api';
import { TypedEvent } from '../../../typed_event';
import { formatDeltaTextElem } from '../../../utils';
import { Component } from '../../component';
import { ItemRenderer } from '../../gear_picker/gear_picker';
import Toast from '../../toast';
import { BulkTab } from '../bulk_tab';

export default class BulkSimResultRenderer extends Component {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, bulkSimUI: BulkTab, result: BulkComboResult, baseResult: BulkComboResult) {
		super(parent, 'bulk-sim-result-root');

		if (!bulkSimUI.simTalents) {
			this.rootElem.classList.add('bulk-sim-result-no-talents');
		}

		const dpsDelta = result.unitMetrics!.dps!.avg! - baseResult.unitMetrics!.dps!.avg;

		const equipButtonRef = ref<HTMLButtonElement>();
		const dpsDeltaRef = ref<HTMLDivElement>();
		const itemsContainerRef = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<>
				<div className="results-sim">
					<div className="results-sim-dps damage-metrics">
						<span className="topline-result-avg">{this.formatDps(result.unitMetrics!.dps!.avg)}</span>
						<div className="results-reference">
							<span ref={dpsDeltaRef} className={clsx('results-reference-diff', dpsDelta >= 0 ? 'positive' : 'negative')} />
						</div>
					</div>
				</div>
				<div ref={itemsContainerRef} className="bulk-gear-combo" />
				{bulkSimUI.simTalents && (
					<div className="bulk-talent-loadout">
						<span>
							{result.talentLoadout && typeof result.talentLoadout === 'object' ? `Talents: ${result.talentLoadout.name}` : 'Current Talents'}
						</span>
					</div>
				)}
				<div className="bulk-results-actions">
					<button ref={equipButtonRef} className={clsx('btn btn-primary bulk-equip-btn', !result.itemsAdded?.length && 'd-none')}>
						Equip
					</button>
				</div>
			</>,
		);

		formatDeltaTextElem(dpsDeltaRef.value!, baseResult.unitMetrics!.dps!.avg, result.unitMetrics!.dps!.avg!, 2);

		if (!!result.itemsAdded?.length) {
			equipButtonRef.value?.addEventListener('click', () => {
				result.itemsAdded.forEach(itemAdded => {
					const item = simUI.sim.db.lookupItemSpec(itemAdded.item!);
					simUI.player.equipItem(TypedEvent.nextEventID(), itemAdded.slot, item);
					simUI.simHeader.activateTab('gear-tab');
				});
				new Toast({
					variant: 'success',
					body: 'Batch gear equipped!',
				});
			});

			const items = (<></>) as HTMLElement;
			for (const is of result.itemsAdded) {
				const itemContainer = (<div className="bulk-result-item" />) as HTMLElement;
				const item = simUI.sim.db.lookupItemSpec(is.item!);
				const renderer = new ItemRenderer(items, itemContainer, simUI.player);
				renderer.update(item!);
				items.appendChild(itemContainer);
			}
			itemsContainerRef.value!.appendChild(items);
		} else if (!result.talentLoadout || typeof result.talentLoadout !== 'object') {
			dpsDeltaRef.value?.classList.add('hide');
			itemsContainerRef.value!.appendChild(<p className="mb-0">Equipped</p>);
		}
	}

	private formatDps(dps: number): string {
		return (Math.round(dps * 100) / 100).toFixed(2);
	}

	private formatDpsDelta(delta: number): string {
		return (delta >= 0 ? '+' : '') + this.formatDps(delta);
	}

	private itemSlotName(is: ItemSpecWithSlot): string {
		return JSON.parse(ItemSpecWithSlot.toJsonString(is, { emitDefaultValues: true }))['slot'].replace('ItemSlot', '');
	}
}
