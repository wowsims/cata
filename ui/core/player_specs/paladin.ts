import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class HolyPaladin extends PlayerSpec<Spec.SpecHolyPaladin> {
	static protoID = Spec.SpecHolyPaladin as Spec.SpecHolyPaladin;
	static playerClass = PlayerClasses.Paladin;
	static friendlyName = 'Holy';
	static fullName = `${this.friendlyName} ${PlayerClasses.Paladin.friendlyName}`;
	static simLink = getSpecSiteUrl('paladin', 'protection');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = HolyPaladin.protoID;
	readonly playerClass = HolyPaladin.playerClass;
	readonly friendlyName = HolyPaladin.friendlyName;
	readonly fullName = HolyPaladin.fullName;
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
	static protoID = Spec.SpecProtectionPaladin as Spec.SpecProtectionPaladin;
	static playerClass = PlayerClasses.Paladin;
	static friendlyName = 'Protection';
	static fullName = `${this.friendlyName} ${PlayerClasses.Paladin.friendlyName}`;
	static simLink = getSpecSiteUrl('paladin', 'protection');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ProtectionPaladin.protoID;
	readonly playerClass = ProtectionPaladin.playerClass;
	readonly friendlyName = ProtectionPaladin.friendlyName;
	readonly fullName = ProtectionPaladin.fullName;
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
	static protoID = Spec.SpecRetributionPaladin as Spec.SpecRetributionPaladin;
	static playerClass = PlayerClasses.Paladin;
	static friendlyName = 'Retribution';
	static fullName = `${this.friendlyName} ${PlayerClasses.Paladin.friendlyName}`;
	static simLink = getSpecSiteUrl('paladin', 'retribution');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = false;

	readonly protoID = RetributionPaladin.protoID;
	readonly playerClass = RetributionPaladin.playerClass;
	readonly friendlyName = RetributionPaladin.friendlyName;
	readonly fullName = RetributionPaladin.fullName;
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
