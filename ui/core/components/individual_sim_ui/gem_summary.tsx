import { Button, Icon, Link } from '@wowsims/ui';
import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../../css_utils';
import { Player } from '../../player';
import { UIGem as Gem } from '../../proto/ui.js';
import { ActionId } from '../../proto_utils/action_id';
import { SimUI } from '../../sim_ui';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
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
			header: { title: 'Gem Summary', extraCssClasses: ['summary-table--gems'] },
			extraCssClasses: ['summary-table--gems'],
		});
		player.gearChangeEmitter.on(() => this.updateTable());
	}

	private updateTable() {
		const body = <></>;
		const fullGemList = this.player.getGear().getAllGems(this.player.isBlacksmithing());
		const hasGems = !!fullGemList.length;
		this.rootElem.classList[!hasGems ? 'add' : 'remove']('hide');

		if (hasGems) {
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
				const linkRef = ref<HTMLAnchorElement>();
				const iconRef = ref<HTMLImageElement>();
				const row = (
					<div className="summary-table-row d-flex align-items-center">
						<Link ref={linkRef} className="summary-table-link" data-whtticon="false" target="_blank">
							<img ref={iconRef} className="gem-icon" />
							<div>{gemName}</div>
						</Link>
						<div>{gemData.count.toFixed(0)}</div>
					</div>
				);
				body.appendChild(row);

				const itemLinkElem = linkRef.value!;
				const iconElem = iconRef.value!;

				setItemQualityCssClass(itemLinkElem, gemData.gem.quality);

				ActionId.fromItemId(gemData.gem.id)
					.fill()
					.then(filledId => {
						iconElem.src = filledId.iconUrl;
						filledId.setWowheadHref(itemLinkElem);
					});
			}

			this.container.bodyElement.replaceChildren(body);

			if (!this.container.headerElement) return;
			const existingResetButton = this.container.headerElement.querySelector('.summary-table-reset-button');
			const resetButton = (
				<Button
					variant="reset"
					size="sm"
					className="summary-table-reset-button"
					iconRight={<Icon icon="close" />}
					onclick={() => {
						this.player.setGear(TypedEvent.nextEventID(), this.player.getGear().withoutGems());
					}}>
					Reset gems
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
