import { IconSize } from '../class';
import { Shaman } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class ElementalShaman extends Spec {
	static protoID = SpecProto.SpecElementalShaman;
	static class = Shaman;
	static friendlyName = 'Elemental';
	static simLink = getSpecSiteUrl('shaman', 'elemental');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ElementalShaman.protoID;
	readonly class = ElementalShaman.class;
	readonly friendlyName = ElementalShaman.friendlyName;
	readonly simLink = ElementalShaman.simLink;

	readonly isTankSpec = ElementalShaman.isTankSpec;
	readonly isHealingSpec = ElementalShaman.isHealingSpec;
	readonly isRangedDpsSpec = ElementalShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ElementalShaman.isMeleeDpsSpec;

	readonly canDualWield = ElementalShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightning.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return ElementalShaman.getIcon(size)
	}
}

export class EnhancementShaman extends Spec {
	static protoID = SpecProto.SpecEnhancementShaman;
	static class = Shaman;
	static friendlyName = 'Enhancement';
	static simLink = getSpecSiteUrl('shaman', 'enhancement');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = EnhancementShaman.protoID;
	readonly class = EnhancementShaman.class;
	readonly friendlyName = EnhancementShaman.friendlyName;
	readonly simLink = EnhancementShaman.simLink;

	readonly isTankSpec = EnhancementShaman.isTankSpec;
	readonly isHealingSpec = EnhancementShaman.isHealingSpec;
	readonly isRangedDpsSpec = EnhancementShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = EnhancementShaman.isMeleeDpsSpec;

	readonly canDualWield = EnhancementShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightningshield.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return EnhancementShaman.getIcon(size)
	}
}

export class RestorationShaman extends Spec {
	static protoID = SpecProto.SpecRestorationShaman;
	static class = Shaman;
	static friendlyName = 'Restoration';
	static simLink = getSpecSiteUrl('shaman', 'restoration');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = RestorationShaman.protoID;
	readonly class = RestorationShaman.class;
	readonly friendlyName = RestorationShaman.friendlyName;
	readonly simLink = RestorationShaman.simLink;

	readonly isTankSpec = RestorationShaman.isTankSpec;
	readonly isHealingSpec = RestorationShaman.isHealingSpec;
	readonly isRangedDpsSpec = RestorationShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = RestorationShaman.isMeleeDpsSpec;

	readonly canDualWield = RestorationShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_magicimmunity.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return RestorationShaman.getIcon(size)
	}
}
