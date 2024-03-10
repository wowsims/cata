import { IconSize } from '../class';
import { Mage } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class ArcaneMage extends Spec {
	static protoID = SpecProto.SpecArcaneMage;
	static class = Mage;
	static friendlyName = 'Arcane';
	static simLink = getSpecSiteUrl('mage', 'arcane');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ArcaneMage.protoID;
	readonly class = ArcaneMage.class;
	readonly friendlyName = ArcaneMage.friendlyName;
	readonly simLink = ArcaneMage.simLink;

	readonly isTankSpec = ArcaneMage.isTankSpec;
	readonly isHealingSpec = ArcaneMage.isHealingSpec;
	readonly isRangedDpsSpec = ArcaneMage.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ArcaneMage.isMeleeDpsSpec;

	readonly canDualWield = ArcaneMage.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_magicalsentry.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return ArcaneMage.getIcon(size)
	}
}

export class FireMage extends Spec {
	static protoID = SpecProto.SpecFireMage;
	static class = Mage;
	static friendlyName = 'Fire';
	static simLink = getSpecSiteUrl('mage', 'fire');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = FireMage.protoID;
	readonly class = FireMage.class;
	readonly friendlyName = FireMage.friendlyName;
	readonly simLink = FireMage.simLink;

	readonly isTankSpec = FireMage.isTankSpec;
	readonly isHealingSpec = FireMage.isHealingSpec;
	readonly isRangedDpsSpec = FireMage.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FireMage.isMeleeDpsSpec;

	readonly canDualWield = FireMage.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_fire_firebolt02.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return FireMage.getIcon(size)
	}
}

export class FrostMage extends Spec {
	static protoID = SpecProto.SpecFrostMage;
	static class = Mage;
	static friendlyName = 'Frost';
	static simLink = getSpecSiteUrl('mage', 'frost');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = FrostMage.protoID;
	readonly class = FrostMage.class;
	readonly friendlyName = FrostMage.friendlyName;
	readonly simLink = FrostMage.simLink;

	readonly isTankSpec = FrostMage.isTankSpec;
	readonly isHealingSpec = FrostMage.isHealingSpec;
	readonly isRangedDpsSpec = FrostMage.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FrostMage.isMeleeDpsSpec;

	readonly canDualWield = FrostMage.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_frost_frostbolt02.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return FrostMage.getIcon(size)
	}
}
