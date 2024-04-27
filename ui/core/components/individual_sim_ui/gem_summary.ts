import { Component } from '../../components/component';
import { setItemQualityCssClass } from '../../css_utils';
import { Player } from '../../player';
import { UIGem as Gem } from '../../proto/ui.js';
import { ActionId } from '../../proto_utils/action_id';
import { SimUI } from '../../sim_ui';
import { ContentBlock } from '../content_block';

interface GemSummaryData {
	gem: Gem;
	count: number;
}

export class GemSummary extends Component {
	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly container: ContentBlock;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'summary-table-root');
		this.simUI = simUI;
		this.player = player;

		this.container = new ContentBlock(this.rootElem, 'summary-table-container', {
			header: { title: 'Gem Summary' },
			extraCssClasses: ['summary-table--gems'],
		});
		player.gearChangeEmitter.on(() => this.updateTable());
	}

	private updateTable() {
		this.container.bodyElement.innerHTML = ``;
		const fullGemList = this.player.getGear().getAllGems(this.player.isBlacksmithing());
		const gemCounts: Record<string, GemSummaryData> = {};

		for (const gem of fullGemList) {
			if (gemCounts[gem.name]) {
				gemCounts[gem.name].count += 1;
			} else {
				gemCounts[gem.name] = {
					gem: gem,
					count: 1,
				};
			}
		}

		for (const gemName of Object.keys(gemCounts)) {
			const gemData = gemCounts[gemName];
			const row = document.createElement('div');
			row.classList.add('summary-table-row', 'd-flex', 'align-items-center');
			row.innerHTML = `
				<a class="summary-table-link" data-whtticon='false' target='_blank'>
					<img class="gem-icon"/>
					<div>${gemName}</div>
				</a>
				<div>${gemData.count.toFixed(0)}</div>
			`;

			const itemLinkElem = row.querySelector('.summary-table-link') as HTMLAnchorElement;
			const iconElem = row.querySelector('.gem-icon') as HTMLImageElement;

			setItemQualityCssClass(itemLinkElem, gemData.gem.quality);

			ActionId.fromItemId(gemData.gem.id)
				.fill()
				.then(filledId => {
					iconElem.src = filledId.iconUrl;
					filledId.setWowheadHref(itemLinkElem);
				});

			this.container.bodyElement.appendChild(row);
		}
	}
}
