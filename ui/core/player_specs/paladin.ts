import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class HolyPaladin extends PlayerSpec<Spec.SpecHolyPaladin> {
	static specIndex = 0;
	static specID = Spec.SpecHolyPaladin as Spec.SpecHolyPaladin;
	static classID = Class.ClassPaladin as Class.ClassPaladin;
	static friendlyName = 'Holy';
	static simLink = getSpecSiteUrl('paladin', 'holy');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = HolyPaladin.specIndex;
	readonly specID = HolyPaladin.specID;
	readonly classID = HolyPaladin.classID;
	readonly friendlyName = HolyPaladin.friendlyName;
	readonly simLink = HolyPaladin.simLink;

	readonly isTankSpec = HolyPaladin.isTankSpec;
	readonly isHealingSpec = HolyPaladin.isHealingSpec;
	readonly isRangedDpsSpec = HolyPaladin.isRangedDpsSpec;
	readonly isMeleeDpsSpec = HolyPaladin.isMeleeDpsSpec;

	readonly canDualWield = HolyPaladin.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_holybolt.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return HolyPaladin.getIcon(size);
	};
}

export class ProtectionPaladin extends PlayerSpec<Spec.SpecProtectionPaladin> {
	static specIndex = 1;
	static specID = Spec.SpecProtectionPaladin as Spec.SpecProtectionPaladin;
	static classID = Class.ClassPaladin as Class.ClassPaladin;
	static friendlyName = 'Protection';
	static simLink = getSpecSiteUrl('paladin', 'protection');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = ProtectionPaladin.specIndex;
	readonly specID = ProtectionPaladin.specID;
	readonly classID = ProtectionPaladin.classID;
	readonly friendlyName = ProtectionPaladin.friendlyName;
	readonly simLink = ProtectionPaladin.simLink;

	readonly isTankSpec = ProtectionPaladin.isTankSpec;
	readonly isHealingSpec = ProtectionPaladin.isHealingSpec;
	readonly isRangedDpsSpec = ProtectionPaladin.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ProtectionPaladin.isMeleeDpsSpec;

	readonly canDualWield = ProtectionPaladin.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_paladin_shieldofthetemplar.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return ProtectionPaladin.getIcon(size);
	};
}

export class RetributionPaladin extends PlayerSpec<Spec.SpecRetributionPaladin> {
	static specIndex = 2;
	static specID = Spec.SpecRetributionPaladin as Spec.SpecRetributionPaladin;
	static classID = Class.ClassPaladin as Class.ClassPaladin;
	static friendlyName = 'Retribution';
	static simLink = getSpecSiteUrl('paladin', 'retribution');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = false;

	readonly specIndex = RetributionPaladin.specIndex;
	readonly specID = RetributionPaladin.specID;
	readonly classID = RetributionPaladin.classID;
	readonly friendlyName = RetributionPaladin.friendlyName;
	readonly simLink = RetributionPaladin.simLink;

	readonly isTankSpec = RetributionPaladin.isTankSpec;
	readonly isHealingSpec = RetributionPaladin.isHealingSpec;
	readonly isRangedDpsSpec = RetributionPaladin.isRangedDpsSpec;
	readonly isMeleeDpsSpec = RetributionPaladin.isMeleeDpsSpec;

	readonly canDualWield = RetributionPaladin.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_auraoflight.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return RetributionPaladin.getIcon(size);
	};
}
