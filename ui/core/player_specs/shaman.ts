import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class ElementalShaman extends PlayerSpec<Spec.SpecElementalShaman> {
	static specIndex = 0;
	static specID = Spec.SpecElementalShaman as Spec.SpecElementalShaman;
	static classID = Class.ClassShaman as Class.ClassShaman;
	static friendlyName = 'Elemental';
	static simLink = getSpecSiteUrl('shaman', 'elemental');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = ElementalShaman.specIndex;
	readonly specID = ElementalShaman.specID;
	readonly classID = ElementalShaman.classID;
	readonly friendlyName = ElementalShaman.friendlyName;
	readonly simLink = ElementalShaman.simLink;

	readonly isTankSpec = ElementalShaman.isTankSpec;
	readonly isHealingSpec = ElementalShaman.isHealingSpec;
	readonly isRangedDpsSpec = ElementalShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ElementalShaman.isMeleeDpsSpec;

	readonly canDualWield = ElementalShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightning.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return ElementalShaman.getIcon(size);
	};
}

export class EnhancementShaman extends PlayerSpec<Spec.SpecEnhancementShaman> {
	static specIndex = 1;
	static specID = Spec.SpecEnhancementShaman as Spec.SpecEnhancementShaman;
	static classID = Class.ClassShaman as Class.ClassShaman;
	static friendlyName = 'Enhancement';
	static simLink = getSpecSiteUrl('shaman', 'enhancement');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = EnhancementShaman.specIndex;
	readonly specID = EnhancementShaman.specID;
	readonly classID = EnhancementShaman.classID;
	readonly friendlyName = EnhancementShaman.friendlyName;
	readonly simLink = EnhancementShaman.simLink;

	readonly isTankSpec = EnhancementShaman.isTankSpec;
	readonly isHealingSpec = EnhancementShaman.isHealingSpec;
	readonly isRangedDpsSpec = EnhancementShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = EnhancementShaman.isMeleeDpsSpec;

	readonly canDualWield = EnhancementShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightningshield.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return EnhancementShaman.getIcon(size);
	};
}

export class RestorationShaman extends PlayerSpec<Spec.SpecRestorationShaman> {
	static specIndex = 2;
	static specID = Spec.SpecRestorationShaman as Spec.SpecRestorationShaman;
	static classID = Class.ClassShaman as Class.ClassShaman;
	static friendlyName = 'Restoration';
	static simLink = getSpecSiteUrl('shaman', 'restoration');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = RestorationShaman.specIndex;
	readonly specID = RestorationShaman.specID;
	readonly classID = RestorationShaman.classID;
	readonly friendlyName = RestorationShaman.friendlyName;
	readonly simLink = RestorationShaman.simLink;

	readonly isTankSpec = RestorationShaman.isTankSpec;
	readonly isHealingSpec = RestorationShaman.isHealingSpec;
	readonly isRangedDpsSpec = RestorationShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = RestorationShaman.isMeleeDpsSpec;

	readonly canDualWield = RestorationShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_magicimmunity.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return RestorationShaman.getIcon(size);
	};
}
