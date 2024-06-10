import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../../individual_sim_ui';
import { BulkComboResult, ItemSpecWithSlot } from '../../../proto/api';
import { TypedEvent } from '../../../typed_event';
import { Component } from '../../component';
import { ItemRenderer } from '../../gear_picker/gear_picker';

export default class BulkSimResultRenderer extends Component {
	constructor(parent: HTMLElement, simUI: IndividualSimUI<any>, result: BulkComboResult, baseResult: BulkComboResult) {
		super(parent, '.bulk-sim-results');
		const dpsDelta = result.unitMetrics!.dps!.avg! - baseResult.unitMetrics!.dps!.avg;

		const equipButtonRef = ref<HTMLButtonElement>();
		const dpsDeltaRef = ref<HTMLDivElement>();
		const itemsContainerRef = ref<HTMLDivElement>();
		parent.appendChild(
			<>
				<div className="results-sim">
					<div className="bulk-result-body-dps bulk-items-text-line results-sim-dps damage-metrics">
						<span className="topline-result-avg">{this.formatDps(result.unitMetrics!.dps!.avg)}</span>

						<span ref={dpsDeltaRef} className={clsx(dpsDelta >= 0 ? 'bulk-result-header-positive' : 'bulk-result-header-negative')}>
							{this.formatDpsDelta(dpsDelta)}
						</span>

						<p className="talent-loadout-text">
							{result.talentLoadout && typeof result.talentLoadout === 'object' ? (
								typeof result.talentLoadout.name === 'string' && <>Talent loadout used: {result.talentLoadout.name}</>
							) : (
								<>Current talents</>
							)}
						</p>
					</div>
				</div>
				<div ref={itemsContainerRef} className="bulk-gear-combo"></div>
				{!!result.itemsAdded?.length && (
					<button ref={equipButtonRef} className="btn btn-primary bulk-equipit">
						Equip
					</button>
				)}
			</>,
		);

		if (!!result.itemsAdded?.length) {
			equipButtonRef.value?.addEventListener('click', () => {
				result.itemsAdded.forEach(itemAdded => {
					const item = simUI.sim.db.lookupItemSpec(itemAdded.item!);
					simUI.player.equipItem(TypedEvent.nextEventID(), itemAdded.slot, item);
					simUI.simHeader.activateTab('gear-tab');
				});
			});

			const items = (<></>) as HTMLElement;
			for (const is of result.itemsAdded) {
				const itemContainer = (<div className="bulk-result-item" />) as HTMLElement;
				const item = simUI.sim.db.lookupItemSpec(is.item!);
				const renderer = new ItemRenderer(items, itemContainer, simUI.player);
				renderer.update(item!);
				renderer.nameElem.appendChild(<a className="bulk-result-item-slot">{this.itemSlotName(is)}</a>);
				items.appendChild(itemContainer);
			}
			itemsContainerRef.value?.appendChild(items);
		} else if (!result.talentLoadout || typeof result.talentLoadout !== 'object') {
			dpsDeltaRef.value?.classList.add('hide');
			parent.appendChild(<p>No changes - this is your currently equipped gear!</p>);
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
