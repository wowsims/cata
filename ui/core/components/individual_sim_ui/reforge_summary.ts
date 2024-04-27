import { setItemQualityCssClass } from '../../css_utils.js';
import { Player } from '../../player.js';
import { ItemSlot, Stat } from '../../proto/common';
import { UIGem as Gem } from '../../proto/ui.js';
import { ActionId } from '../../proto_utils/action_id.js';
import { EquippedItem } from '../../proto_utils/equipped_item';
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
		const gear = this.player.getGear();
		const totals: ReforgeSummaryTotal = {};
		gear.getItemSlots().forEach(itemSlot => {
			const item = gear.getEquippedItem(itemSlot);
			if (item && item._reforging !== 0) {
				const reforge = this.player.getReforge(item.reforging);
				if (reforge) {
					// const row = document.createElement('div');

					// row.classList.add('summary-table-row', 'd-flex', 'align-items-center');

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

					/* row.innerHTML = `
						<img class="gem-icon"/>
						<div>${fromAmount} ${shortSecondaryStatNames.get(fromStat)} > +${toAmount} ${shortSecondaryStatNames.get(toStat)}</div>
						<button type="button" class="btn-close" aria-label="Close"><i class="fas fa-times fa-1xl"></i></button>
					`;

					const removeButton = row.querySelector('.btn-close') as HTMLButtonElement;
					removeButton.onclick = () => {
						this.player.equipItem(TypedEvent.nextEventID(), itemSlot, item.withItem(item.item));
					};

					const iconElem = row.querySelector('.gem-icon') as HTMLImageElement;

					ActionId.fromItemId(item.item.id)
						.fill()
						.then(filledId => {
							iconElem.src = filledId.iconUrl;
						});

					this.container.bodyElement.appendChild(row); */
				}
			}
		});

		const footer = document.createElement('div');
		footer.classList.add('summary-table-row');
		this.container.bodyElement.appendChild(footer);
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
			footer.appendChild(row);
		});
	}
}
