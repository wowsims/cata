import { Button } from '@wowsims/ui';

import { Player } from '../../player.js';
import { Stat } from '../../proto/common.js';
import { shortSecondaryStatNames } from '../../proto_utils/names.js';
import { SimUI } from '../../sim_ui.js';
import { TypedEvent } from '../../typed_event.js';
import { Component } from '../component.js';
import { ContentBlock } from '../content_block.jsx';

type ReforgeSummaryTotal = {
	[key in Stat]?: number;
};

export class ReforgeSummary extends Component {
	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly container: ContentBlock;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'summary-table-root');
		this.simUI = simUI;
		this.player = player;

		this.container = new ContentBlock(this.rootElem, 'summary-table-container', {
			header: { title: 'Reforge Summary' },
			extraCssClasses: ['summary-table--reforge'],
		});

		player.gearChangeEmitter.on(() => this.updateTable());
	}

	private updateTable() {
		const body = <></>;
		let gear = this.player.getGear();
		const totals: ReforgeSummaryTotal = {};
		gear.getItemSlots().forEach(itemSlot => {
			const item = gear.getEquippedItem(itemSlot);
			if (item?.reforge && item.reforge?.id !== 0) {
				const reforge = this.player.getReforgeData(item, item.reforge);
				if (reforge) {
					const { fromAmount, toAmount } = reforge;
					const fromStat = reforge.fromStat[0];
					const toStat = reforge.toStat[0];

					if (typeof totals[fromStat] !== 'number') {
						totals[fromStat] = 0;
					}
					if (typeof totals[toStat] !== 'number') {
						totals[toStat] = 0;
					}
					if (fromAmount) totals[fromStat]! += fromAmount;
					if (toAmount) totals[toStat]! += toAmount;
				}
			}
		});

		const hasReforgedItems = !!Object.keys(totals).length;
		this.rootElem.classList[!hasReforgedItems ? 'add' : 'remove']('hide');

		if (hasReforgedItems) {
			Object.keys(totals).forEach(key => {
				const stat: Stat = Number(key);
				const value = totals[stat];
				if (!value) return;

				body.appendChild(
					<div className="summary-table-row d-flex align-items-center">
						<div>{shortSecondaryStatNames.get(stat)}</div>
						<div className={`${value === 0 ? '' : value > 0 ? 'positive' : 'negative'}`}>{value}</div>
					</div>,
				);
			});

			this.container.bodyElement.replaceChildren(body);

			if (!this.container.headerElement) return;
			const existingResetButton = this.container.headerElement.querySelector('.summary-table-reset-button');
			const resetButton = (
				<Button
					variant="reset"
					size="sm"
					className="summary-table-reset-button"
					iconRight="close"
					onclick={() => {
						gear.getItemSlots().forEach(itemSlot => {
							const item = gear.getEquippedItem(itemSlot);
							if (item) gear = gear.withEquippedItem(itemSlot, item.withItem(item.item), this.player.canDualWield2H());
						});
						this.player.setGear(TypedEvent.nextEventID(), gear);
					}}>
					Reset reforges
				</Button>
			);

			if (existingResetButton) {
				this.container.headerElement.replaceChild(resetButton, existingResetButton);
			} else {
				this.container.headerElement.appendChild(resetButton);
			}
		}
	}
}
