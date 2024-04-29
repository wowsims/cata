// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { Popover } from 'bootstrap';
import { element, fragment, ref } from 'tsx-vanilla';

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
	popoverElement: HTMLElement;
	popoverConfig: Partial<Popover.Options>;
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
	popover: Popover | null = null;
	item: EquippedItem;
	constructor(config: QuickSwapListConfig<T>) {
		this.config = config;
		this.item = config.item;
		this.attachPopover();
	}
	attachPopover() {
		const config: QuickSwapListConfig<T>['popoverConfig'] = {
			html: true,
			trigger: 'focus',
			customClass: 'tooltip-quick-swap',
			content: () => this.buildList(),
			...this.config.popoverConfig,
		};
		this.popover = Popover.getOrCreateInstance(this.config.popoverElement, config);
	}
	update(config: Partial<QuickSwapListConfig<T>>) {
		this.config = {
			...this.config,
			...config,
		};
		if (config.item) this.item = config.item;
		this.buildList();
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
						const iconElem = ref<HTMLImageElement>();
						const labelElem = ref<HTMLSpanElement>();
						const listItem = (
							<li className="tooltip-quick-swap__list-item">
								<a
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

						ActionId.fromItemId('id' in item.item ? item.item.id : item.item.effectId)
							.fill()
							.then(filledId => {
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
					<a onclick={data.footerButton.onClick} href="javascript:void(0)" className="btn btn-sm btn-primary" attributes={{ role: 'button' }}>
						{data.footerButton.label}
					</a>
				</div>
			)}
		</>
	);
};

export default QuickSwapList;
