import { IconSize } from '../class';
import { Paladin } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class HolyPaladin extends Spec {
	static protoID = SpecProto.SpecHolyPaladin;
	static class = Paladin;
	static friendlyName = 'Holy';
	static simLink = getSpecSiteUrl('paladin', 'protection');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = HolyPaladin.protoID;
	readonly class = HolyPaladin.class;
	readonly friendlyName = HolyPaladin.friendlyName;
	readonly simLink = HolyPaladin.simLink;

	readonly isTankSpec = HolyPaladin.isTankSpec;
	readonly isHealingSpec = HolyPaladin.isHealingSpec;
	readonly isRangedDpsSpec = HolyPaladin.isRangedDpsSpec;
	readonly isMeleeDpsSpec = HolyPaladin.isMeleeDpsSpec;

	readonly canDualWield = HolyPaladin.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_holybolt.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return HolyPaladin.getIcon(size)
	}
}

export class ProtectionPaladin extends Spec {
	static protoID = SpecProto.SpecProtectionPaladin;
	static class = Paladin;
	static friendlyName = 'Protection';
	static simLink = getSpecSiteUrl('paladin', 'protection');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ProtectionPaladin.protoID;
	readonly class = ProtectionPaladin.class;
	readonly friendlyName = ProtectionPaladin.friendlyName;
	readonly simLink = ProtectionPaladin.simLink;

	readonly isTankSpec = ProtectionPaladin.isTankSpec;
	readonly isHealingSpec = ProtectionPaladin.isHealingSpec;
	readonly isRangedDpsSpec = ProtectionPaladin.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ProtectionPaladin.isMeleeDpsSpec;

	readonly canDualWield = ProtectionPaladin.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_paladin_shieldofthetemplar.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return ProtectionPaladin.getIcon(size)
	}
}

export class RetributionPaladin extends Spec {
	static protoID = SpecProto.SpecRetributionPaladin;
	static class = Paladin;
	static friendlyName = 'Retribution';
	static simLink = getSpecSiteUrl('paladin', 'retribution');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = false;

	readonly protoID = RetributionPaladin.protoID;
	readonly class = RetributionPaladin.class;
	readonly friendlyName = RetributionPaladin.friendlyName;
	readonly simLink = RetributionPaladin.simLink;

	readonly isTankSpec = RetributionPaladin.isTankSpec;
	readonly isHealingSpec = RetributionPaladin.isHealingSpec;
	readonly isRangedDpsSpec = RetributionPaladin.isRangedDpsSpec;
	readonly isMeleeDpsSpec = RetributionPaladin.isMeleeDpsSpec;

	readonly canDualWield = RetributionPaladin.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_auraoflight.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return RetributionPaladin.getIcon(size)
	}
}
