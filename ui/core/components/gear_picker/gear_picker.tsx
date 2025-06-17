import { ref } from 'tsx-vanilla';

import { MISSING_RANDOM_SUFFIX_WARNING } from '../../constants/item_notices';
import { setItemQualityCssClass } from '../../css_utils';
import { Player } from '../../player';
import { ItemLevelState, ItemSlot, ItemType } from '../../proto/common';
import { UIEnchant as Enchant, UIGem as Gem } from '../../proto/ui';
import { ActionId } from '../../proto_utils/action_id';
import { getEnchantDescription } from '../../proto_utils/enchants';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { shortSecondaryStatNames, slotNames } from '../../proto_utils/names';
import { SimUI } from '../../sim_ui';
import { EventID } from '../../typed_event';
import { Component } from '../component';
import { ItemNotice } from '../item_notice/item_notice';
import QuickSwapList from '../quick_swap';
import { GearData } from './item_list';
import { addQuickEnchantPopover } from './quick_enchant_popover';
import { addQuickGemPopover } from './quick_gem_popover';
import SelectorModal, { SelectorModalTabs } from './selector_modal';
import { createGemContainer, createNameDescriptionLabel, getEmptySlotIconUrl } from './utils';

export default class GearPicker extends Component {
	// ItemSlot is used as the index
	readonly itemPickers: Array<ItemPicker>;
	readonly selectorModal: SelectorModal;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'gear-picker-root');

		const leftSideRef = ref<HTMLDivElement>();
		const rightSideRef = ref<HTMLDivElement>();

		this.rootElem.appendChild(
			<>
				<div ref={leftSideRef} className="gear-picker-left tab-panel-col"></div>
				<div ref={rightSideRef} className="gear-picker-right tab-panel-col"></div>
			</>,
		);

		const leftItemPickers = [
			ItemSlot.ItemSlotHead,
			ItemSlot.ItemSlotNeck,
			ItemSlot.ItemSlotShoulder,
			ItemSlot.ItemSlotBack,
			ItemSlot.ItemSlotChest,
			ItemSlot.ItemSlotWrist,
			ItemSlot.ItemSlotMainHand,
			ItemSlot.ItemSlotOffHand,
		].map(slot => new ItemPicker(leftSideRef.value!, this, simUI, player, slot));

		const rightItemPickers = [
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
			ItemSlot.ItemSlotFinger1,
			ItemSlot.ItemSlotFinger2,
			ItemSlot.ItemSlotTrinket1,
			ItemSlot.ItemSlotTrinket2,
		].map(slot => new ItemPicker(rightSideRef.value!, this, simUI, player, slot));

		this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);

		this.selectorModal = new SelectorModal(simUI.rootElem, simUI, player, this, { id: 'gear-picker-selector-modal' });
	}
}

export class ItemRenderer extends Component {
	private readonly player: Player<any>;

	readonly iconElem: HTMLAnchorElement;
	readonly nameContainerElem: HTMLDivElement;
	readonly nameElem: HTMLAnchorElement;
	readonly ilvlElem: HTMLSpanElement;
	readonly tinkerElem: HTMLAnchorElement;
	readonly enchantElem: HTMLAnchorElement;
	readonly reforgeElem: HTMLAnchorElement;
	readonly socketsContainerElem: HTMLElement;
	private notice: ItemNotice | null = null;
	socketsElem: HTMLAnchorElement[] = [];

	// Can be used to remove any events in addEventListener
	// https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#add_an_abortable_listener
	public abortController?: AbortController;
	public signal?: AbortSignal;

	constructor(parent: HTMLElement, root: HTMLElement, player: Player<any>) {
		super(parent, 'item-picker-root', root);
		this.player = player;

		const iconElem = ref<HTMLAnchorElement>();
		const nameContainerElem = ref<HTMLDivElement>();
		const nameElem = ref<HTMLAnchorElement>();
		const ilvlElem = ref<HTMLSpanElement>();
		const enchantElem = ref<HTMLAnchorElement>();
		const tinkerElem = ref<HTMLAnchorElement>();
		const reforgeElem = ref<HTMLAnchorElement>();
		const sce = ref<HTMLDivElement>();

		this.rootElem.appendChild(
			<>
				<div className="item-picker-icon-wrapper">
					<span className="item-picker-ilvl" ref={ilvlElem} />
					<a ref={iconElem} className="item-picker-icon" href="javascript:void(0)" attributes={{ role: 'button' }} />
					<div ref={sce} className="item-picker-sockets-container"></div>
				</div>
				<div className="item-picker-labels-container">
					<div ref={nameContainerElem} className="item-picker-name-row d-flex gap-1">
						<a ref={nameElem} className="item-picker-name-container" href="javascript:void(0)" attributes={{ role: 'button' }} />
					</div>
					<a ref={enchantElem} className="item-picker-enchant hide" href="javascript:void(0)" attributes={{ role: 'button' }} />
					<a ref={tinkerElem} className="item-picker-tinker hide" href="javascript:void(0)" attributes={{ role: 'button' }} />
					<a ref={reforgeElem} className="item-picker-reforge hide" href="javascript:void(0)" attributes={{ role: 'button' }} />
				</div>
			</>,
		);

		this.iconElem = iconElem.value!;
		this.nameContainerElem = nameContainerElem.value!;
		this.nameElem = nameElem.value!;
		this.ilvlElem = ilvlElem.value!;
		this.reforgeElem = reforgeElem.value!;
		this.enchantElem = enchantElem.value!;
		this.tinkerElem = tinkerElem.value!;
		this.socketsContainerElem = sce.value!;
	}

	clear(slot: ItemSlot) {
		this.abortController?.abort();
		this.nameElem.removeAttribute('data-wowhead');
		this.nameElem.removeAttribute('href');
		this.notice?.dispose();
		this.notice = null;
		this.iconElem.removeAttribute('data-wowhead');
		this.iconElem.removeAttribute('href');
		this.enchantElem.removeAttribute('data-wowhead');
		this.enchantElem.removeAttribute('href');
		this.tinkerElem.removeAttribute('data-wowhead');
		this.tinkerElem.removeAttribute('href');
		this.enchantElem.classList.add('hide');
		this.reforgeElem.classList.add('hide');

		this.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(slot)}')`;

		this.enchantElem.replaceChildren();
		this.tinkerElem.replaceChildren();
		this.reforgeElem.replaceChildren();
		this.socketsContainerElem.replaceChildren();
		this.nameElem.replaceChildren();
		this.ilvlElem.replaceChildren();

		this.socketsElem = [];
	}

	update(newItem: EquippedItem) {
		this.abortController = new AbortController();
		this.signal = this.abortController.signal;

		const nameSpan = <span className="item-picker-name">{newItem.item.name}</span>;
		const isEligibleForRandomSuffix = !!newItem.hasRandomSuffixOptions();
		const hasRandomSuffix = !!newItem.randomSuffix;
		this.nameElem.replaceChildren(nameSpan);
		this.ilvlElem.replaceChildren(
			<>
				{newItem.ilvl.toString()}
				{!!(newItem.upgrade !== ItemLevelState.ChallengeMode && newItem.ilvlFromBase) && (
					<span className="item-quality-uncommon">+{newItem.ilvlFromBase}</span>
				)}
			</>,
		);

		if (hasRandomSuffix) {
			nameSpan.textContent += ' ' + newItem.randomSuffix.name;
		}

		if (newItem.item.nameDescription) {
			this.nameElem.appendChild(createNameDescriptionLabel(newItem.item.nameDescription));
		}

		this.notice = new ItemNotice(this.player, {
			itemId: newItem.item.id,
			additionalNoticeData: isEligibleForRandomSuffix && !hasRandomSuffix ? MISSING_RANDOM_SUFFIX_WARNING : undefined,
		});

		if (this.notice.hasNotice) {
			this.nameContainerElem.appendChild(this.notice.rootElem);
		}

		const reforgeData = newItem.withDynamicStats().getReforgeData();
		if (reforgeData) {
			const fromText = shortSecondaryStatNames.get(reforgeData.reforge?.fromStat);
			const toText = shortSecondaryStatNames.get(reforgeData.reforge?.toStat);
			this.reforgeElem.innerText = `Reforged ${Math.abs(reforgeData.fromAmount)} ${fromText} → ${reforgeData.toAmount} ${toText}`;
			this.reforgeElem.classList.remove('hide');
		} else {
			this.reforgeElem.innerText = '';
			this.reforgeElem.classList.add('hide');
		}

		setItemQualityCssClass(this.nameElem, newItem.item.quality);

		this.player.setWowheadData(newItem, this.iconElem);
		this.player.setWowheadData(newItem, this.nameElem);

		newItem
			.asActionId()
			.fill(undefined, { signal: this.signal })
			.then(filledId => {
				if (this.signal?.aborted) return;
				filledId.setBackgroundAndHref(this.iconElem);
				filledId.setWowheadHref(this.nameElem);
			});

		if (newItem.enchant) {
			getEnchantDescription(newItem.enchant).then(description => {
				this.enchantElem.textContent = description;
			});
			// Make enchant text hover have a tooltip.
			if (newItem.enchant.spellId) {
				this.enchantElem.href = ActionId.makeSpellUrl(newItem.enchant.spellId);
				ActionId.makeSpellTooltipData(newItem.enchant.spellId).then(url => {
					this.enchantElem.dataset.wowhead = url;
				});
			} else {
				this.enchantElem.href = ActionId.makeItemUrl(newItem.enchant.itemId);
				ActionId.makeItemTooltipData(newItem.enchant.itemId).then(url => {
					this.enchantElem.dataset.wowhead = url;
				});
			}
			this.enchantElem.dataset.whtticon = 'false';
			this.enchantElem.classList.remove('hide');
		} else {
			this.enchantElem.classList.add('hide');
		}

		if (newItem.tinker) {
			getEnchantDescription(newItem.tinker).then(description => {
				this.tinkerElem.textContent = description;
			});
			// Make enchant text hover have a tooltip.
			if (newItem.tinker.spellId) {
				this.tinkerElem.href = ActionId.makeSpellUrl(newItem.tinker.spellId);
				ActionId.makeSpellTooltipData(newItem.tinker.spellId).then(url => {
					this.tinkerElem.dataset.wowhead = url;
				});
			} else {
				this.enchantElem.href = ActionId.makeItemUrl(newItem.tinker.itemId);
				ActionId.makeItemTooltipData(newItem.tinker.itemId).then(url => {
					this.tinkerElem.dataset.wowhead = url;
				});
			}
			this.tinkerElem.dataset.whtticon = 'false';
			this.tinkerElem.classList.remove('hide');
		} else {
			this.tinkerElem.classList.add('hide');
		}

		newItem.allSocketColors().forEach((socketColor, gemIdx) => {
			const gemContainer = createGemContainer(socketColor, newItem.gems[gemIdx], gemIdx);
			if (gemIdx === newItem.numPossibleSockets - 1 && [ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(newItem.item.type)) {
				const updateProfession = () => {
					gemContainer.classList[this.player.isBlacksmithing() ? 'remove' : 'add']('hide');
				};
				this.player.professionChangeEmitter.on(updateProfession);
				updateProfession();
			}
			this.socketsElem.push(gemContainer);
			this.socketsContainerElem.appendChild(gemContainer);
		});
	}
}

export class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly onUpdateCallbacks: (() => void)[] = [];

	private readonly itemElem: ItemRenderer;
	private readonly gearPicker: GearPicker;

	// All items and enchants that are eligible for this slot
	private _equippedItem: EquippedItem | null = null;

	private quickSwapEnchantPopover: QuickSwapList<Enchant> | null = null;
	private quickSwapGemPopover: QuickSwapList<Gem>[] = [];

	constructor(parent: HTMLElement, gearPicker: GearPicker, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');

		this.gearPicker = gearPicker;
		this.simUI = simUI;
		this.player = player;
		this.slot = slot;
		this.itemElem = new ItemRenderer(parent, this.rootElem, player);

		this.item = player.getEquippedItem(slot);

		player.sim.waitForInit().then(() => {
			const openGearSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Items);
			};
			const openReforgeSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Reforging);
			};
			const openTinkerSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Tinkers);
			};

			this.itemElem.iconElem.addEventListener('click', openGearSelector);
			this.itemElem.nameElem.addEventListener('click', openGearSelector);
			this.itemElem.reforgeElem.addEventListener('click', openReforgeSelector);
			this.itemElem.tinkerElem.addEventListener('click', openTinkerSelector);
			this.addQuickEnchantHelpers();
		});

		player.gearChangeEmitter.on(() => {
			this.item = this.player.getEquippedItem(this.slot);
			if (this._equippedItem) {
				if (this._equippedItem !== this.quickSwapEnchantPopover?.item) {
					this.quickSwapEnchantPopover?.update({ item: this._equippedItem });
				}
				this.addQuickGemHelpers();
			}
		});

		player.sim.filtersChangeEmitter.on(() => {
			if (this._equippedItem) {
				this.quickSwapEnchantPopover?.update({ item: this._equippedItem });
				this.quickSwapGemPopover.forEach(quickSwap => quickSwap.update({ item: this._equippedItem! }));
			}
		});

		player.sim.showQuickSwapChangeEmitter.on(() => {
			this.quickSwapEnchantPopover?.tooltip?.[this.player.sim.getShowQuickSwap() ? 'enable' : 'disable']();
			this.quickSwapGemPopover.forEach(quickSwap => quickSwap.tooltip?.[this.player.sim.getShowQuickSwap() ? 'enable' : 'disable']());
		});

		player.professionChangeEmitter.on(() => {
			if (!!this._equippedItem) {
				this.player.setWowheadData(this._equippedItem, this.itemElem.iconElem);
			}
		});
	}

	createGearData(): GearData {
		return {
			equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
				this.player.equipItem(eventID, this.slot, equippedItem);
			},
			getEquippedItem: () => this.player.getEquippedItem(this.slot)?.withChallengeMode(this.player.getChallengeModeEnabled()).withDynamicStats() || null,
			changeEvent: this.player.gearChangeEmitter,
		};
	}

	get item(): EquippedItem | null {
		return this._equippedItem;
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.itemElem.clear(this.slot);
		// Clear quick swap gems array since gem sockets are rerendered every time
		this.quickSwapGemPopover = [];
		this.itemElem.nameElem.textContent = slotNames.get(this.slot) ?? '';
		setItemQualityCssClass(this.itemElem.nameElem, null);

		if (!!newItem) {
			this.itemElem.update(newItem);
		}

		this._equippedItem = newItem;
		this.onUpdateCallbacks.forEach(callback => callback());
	}

	onUpdate(callback: () => void) {
		this.onUpdateCallbacks.push(callback);
	}

	openSelectorModal(selectedTab: SelectorModalTabs) {
		this.gearPicker.selectorModal.openTab(this.slot, selectedTab, this.createGearData());
	}

	private addQuickGemHelpers() {
		if (!this._equippedItem) return;
		const openGemDetailTab = (socketIdx: number) => this.openSelectorModal(`Gem${socketIdx + 1}` as SelectorModalTabs);
		this.itemElem.socketsElem?.forEach(element => {
			const socketIdx = Number(element.dataset.socketIdx) || 0;
			element.addEventListener('click', event => {
				event.preventDefault();
				openGemDetailTab(0);
			});
			const popover = addQuickGemPopover(this.player, element, this._equippedItem!, this.slot, socketIdx, () => openGemDetailTab(socketIdx));
			if (!this.player.sim.getShowQuickSwap()) popover.tooltip?.disable();
			this.quickSwapGemPopover.push(popover);
		});
	}

	private addQuickEnchantHelpers() {
		if (!this._equippedItem) return;
		const openEnchantSelector = () => this.openSelectorModal(SelectorModalTabs.Enchants);
		this.itemElem.enchantElem.addEventListener('click', event => {
			event?.preventDefault();
			openEnchantSelector();
		});
		this.quickSwapEnchantPopover = addQuickEnchantPopover(this.player, this.itemElem.enchantElem, this._equippedItem, this.slot, openEnchantSelector);
		if (!this.player.sim.getShowQuickSwap()) this.quickSwapEnchantPopover.tooltip?.disable();
	}
}
