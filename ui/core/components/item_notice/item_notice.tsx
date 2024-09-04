import tippy, { Instance as TippyInstance } from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { GENERIC_MISSING_SET_BONUS_NOTICE_DATA, ITEM_NOTICES, SET_BONUS_NOTICES } from '../../constants/item_notices';
import { Player } from '../../player';
import { Spec } from '../../proto/common';
import { Database } from '../../proto_utils/database';
import { Component } from '../component';

export type ItemNoticeData = {
	// SpecUnknown is used as default and should always be present
	[Spec.SpecUnknown]: JSX.Element;
} & Record<number, JSX.Element>;

type ItemNoticeConfig = {
	itemId: number;
};

// Keys are item counts for each set bonus (typically 2 and 4), values are the
// notice that should be displayed for each bonus. If null, will default to
// GENERIC_MISSING_SET_BONUS_NOTICE_DATA.
export type SetBonusNoticeData = Map<number, string> | null;

export class ItemNotice extends Component {
	itemId: number;
	player: Player<any>;
	tooltip: TippyInstance | null = null;
	constructor(player: Player<any>, config: ItemNoticeConfig) {
		super(null, 'item-notice');
		this.rootElem.classList.add('d-inline');
		this.itemId = config.itemId;
		this.player = player;

		if (this.hasNotice) this.rootElem.appendChild(this.template!);

		this.addOnDisposeCallback(() => {
			this.tooltip?.destroy();
			this.rootElem?.remove();
		});
	}

	get hasNotice() {
		return ITEM_NOTICES.has(this.itemId);
	}

	private get noticeContent() {
		if (!this.hasNotice) return null;
		const itemNotice = ITEM_NOTICES.get(this.itemId)!;
		return itemNotice[this.player.getSpec()] || itemNotice[Spec.SpecUnknown];
	}

	private get template() {
		if (!this.hasNotice) return null;
		const tooltipContent = this.noticeContent?.cloneNode(true);
		if (!tooltipContent) return null;

		const noticeIconRef = ref<HTMLButtonElement>();
		const template = <button ref={noticeIconRef} className="warning fa fa-exclamation-triangle fa-xl me-2"></button>;

		this.tooltip = tippy(noticeIconRef.value!, {
			content: tooltipContent.cloneNode(true) as HTMLElement,
		});

		return template;
	}

	static registerSetBonusNotices(db: Database) {
		SET_BONUS_NOTICES.forEach((value: SetBonusNoticeData, key: number) => {
			const noticeData = value || GENERIC_MISSING_SET_BONUS_NOTICE_DATA;
			const noticeContent = (
				<>
					<p className="mb-1"> This item set has the following warnings:</p>
					<ul className="mb-0">
						{Array.from(noticeData.keys()).map(key => <li>{key.toFixed(0)}-piece: {noticeData.get(key)!}</li>)}
					</ul>
				</>
			);

			for (const id of db.getItemIdsForSet(key)) {
				ITEM_NOTICES.set(id, { [Spec.SpecUnknown]: noticeContent });
			}
		});
	}
}
