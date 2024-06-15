import tippy, { hideAll, Instance as TippyInstance, Props as TippyProps } from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../css_utils';
import { UIEnchant as Enchant, UIGem as Gem } from '../proto/ui.js';
import { ActionId } from '../proto_utils/action_id';
import { EquippedItem } from '../proto_utils/equipped_item';

type QuickSwapAllowedItem = Gem | Enchant;

type QuickSwapListItem<T extends QuickSwapAllowedItem> = {
	item: T;
	active: boolean;
};

type QuickSwapListConfig<T extends QuickSwapAllowedItem> = {
	title: string;
	tippyElement: HTMLElement;
	tippyConfig: Partial<TippyProps>;
	emptyMessage: string;
	item: EquippedItem;
	getItems: (currentItem: EquippedItem) => QuickSwapListItem<T>[];
	onItemClick: (item: T) => void;
	footerButton?: {
		label: string;
		onClick: () => void;
	};
};

class QuickSwapList<T extends QuickSwapAllowedItem> {
	config: QuickSwapListConfig<T>;
	tooltip: TippyInstance | null = null;
	item: EquippedItem;
	constructor(config: QuickSwapListConfig<T>) {
		this.config = config;
		this.item = config.item;
		this.attachTooltip();
	}
	attachTooltip() {
		const config: QuickSwapListConfig<T>['tippyConfig'] = {
			trigger: 'mouseenter',
			triggerTarget: this.config.tippyElement,
			interactive: true,
			interactiveBorder: 10,
			offset: [0, 5],
			animation: false,
			placement: 'bottom',
			theme: 'tooltip-quick-swap',
			content: () => this.buildList(),
			onCreate: instance => {
				instance.popper.addEventListener('click', () => instance.hide());
			},
			onShow: instance => {
				hideAll({ exclude: instance });
			},
			...this.config.tippyConfig,
		};
		this.tooltip = tippy(this.config.tippyElement, config);
	}
	update(config: Partial<QuickSwapListConfig<T>>) {
		this.config = {
			...this.config,
			...config,
		};
		if (config.item) this.item = config.item;
		this.tooltip?.setContent(this.buildList());
	}
	buildList() {
		return buildList(this.config);
	}
}

const buildList = <T extends QuickSwapAllowedItem>(data: QuickSwapListConfig<T>) => {
	const items = data.getItems(data.item);
	return (
		<>
			<h3 className="tooltip-quick-swap__title h6 text-center">{data.title}</h3>
			{items.length ? (
				<ul className="tooltip-quick-swap__list">
					{items.map(item => {
						const anchorElem = ref<HTMLAnchorElement>();
						const iconElem = ref<HTMLImageElement>();
						const labelElem = ref<HTMLSpanElement>();
						const listItem = (
							<li className="tooltip-quick-swap__list-item">
								<a
									ref={anchorElem}
									href="javascript:void(0)"
									className={`tooltip-quick-swap__anchor d-flex align-items-center ${item.active ? ' active' : ''}`}
									onclick={() => data.onItemClick(item.item)}>
									<img ref={iconElem} alt={item.item.name} className="tooltip-quick-swap__icon gem-icon flex-shrink-0" />
									<span ref={labelElem} className="tooltip-quick-swap__label text-start">
										{item.item.name}
									</span>
								</a>
							</li>
						);
						if (labelElem.value) setItemQualityCssClass(labelElem.value, item.item.quality);
						('spellId' in item.item ? ActionId.fromSpellId(item.item.spellId) : ActionId.fromItemId(item.item.id)).fill().then(filledId => {
							filledId.setWowheadHref(anchorElem.value!);
							iconElem.value!.src = filledId.iconUrl;
						});

						return listItem;
					})}
				</ul>
			) : (
				<p className="tooltip-quick-swap__empty">{data.emptyMessage}</p>
			)}
			{data.footerButton && (
				<div className="tooltip-quick-swap__footer d-flex justify-content-center">
					<button onclick={data.footerButton.onClick} className="btn btn-sm btn-primary">
						{data.footerButton.label}
					</button>
				</div>
			)}
		</>
	);
};

export default QuickSwapList;
