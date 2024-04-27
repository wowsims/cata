import { Player } from '../../player.js';
import { Stat } from '../../proto/common';
import { shortSecondaryStatNames } from '../../proto_utils/names';
import { SimUI } from '../../sim_ui.js';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component.js';
import { ContentBlock } from '../content_block.js';

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
		this.container.bodyElement.innerHTML = ``;
		let gear = this.player.getGear();
		const totals: ReforgeSummaryTotal = {};
		gear.getItemSlots().forEach(itemSlot => {
			const item = gear.getEquippedItem(itemSlot);
			if (item && item._reforging !== 0) {
				const reforge = this.player.getReforge(item.reforging);
				if (reforge) {
					const fromStat = reforge.fromStat[0];
					const toStat = reforge.toStat[0];
					const fromAmount = Math.ceil(-item.item.stats[fromStat] * reforge.multiplier);
					const toAmount = Math.floor(item.item.stats[fromStat] * reforge.multiplier);

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
				const row = document.createElement('div');
				row.classList.add('summary-table-row', 'd-flex', 'align-items-center');
				row.innerHTML = `
					<div>${shortSecondaryStatNames.get(stat)}</div>
					<div class="${value === 0 ? '' : value > 0 ? 'positive' : 'negative'}">${value}</div>
				`;
				this.container.bodyElement.appendChild(row);
			});

			const footer = document.createElement('div');
			footer.classList.add('summary-table-footer', 'd-flex', 'justify-content-end');
			footer.innerHTML = `
				<button class="btn btn-sm btn-outline-primary">
					<i class="fas fa-close me-1"></i>
					Reset reforged items
				</button>
			`;

			const resetButton = footer.querySelector('button') as HTMLButtonElement;
			resetButton.onclick = () => {
				gear.getItemSlots().forEach(itemSlot => {
					const item = gear.getEquippedItem(itemSlot);
					if (item) gear = gear.withEquippedItem(itemSlot, item.withItem(item.item), this.player.canDualWield2H());
				});
				this.player.setGear(TypedEvent.nextEventID(), gear);
			};
			this.container.bodyElement.appendChild(footer);
		}
	}
}
