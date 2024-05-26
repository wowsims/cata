import { ItemQuality } from './proto/common.js';

const itemQualityCssClasses: Record<ItemQuality, string> = {
	[ItemQuality.ItemQualityJunk]: 'text-junk',
	[ItemQuality.ItemQualityCommon]: 'text-common',
	[ItemQuality.ItemQualityUncommon]: 'text-uncommon',
	[ItemQuality.ItemQualityRare]: 'text-rare',
	[ItemQuality.ItemQualityEpic]: 'text-epic',
	[ItemQuality.ItemQualityLegendary]: 'text-legendary',
	[ItemQuality.ItemQualityArtifact]: 'text-artifact',
	[ItemQuality.ItemQualityHeirloom]: 'text-heirloom'
};
export function setItemQualityCssClass(elem: HTMLElement, quality: ItemQuality | null) {
	Object.values(itemQualityCssClasses).forEach(cssClass => elem.classList.remove(cssClass));

	if (quality) {
		elem.classList.add(itemQualityCssClasses[quality]);
	}
}
