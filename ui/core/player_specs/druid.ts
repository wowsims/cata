import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class BalanceDruid extends PlayerSpec<Spec.SpecBalanceDruid> {
	static specIndex = 0;
	static specID = Spec.SpecBalanceDruid as Spec.SpecBalanceDruid;
	static classID = Class.ClassDruid as Class.ClassDruid;
	static friendlyName = 'Balance';
	static simLink = getSpecSiteUrl('druid', 'balance');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = BalanceDruid.specIndex;
	readonly specID = BalanceDruid.specID;
	readonly classID = BalanceDruid.classID;
	readonly friendlyName = BalanceDruid.friendlyName;
	readonly simLink = BalanceDruid.simLink;

	readonly isTankSpec = BalanceDruid.isTankSpec;
	readonly isHealingSpec = BalanceDruid.isHealingSpec;
	readonly isRangedDpsSpec = BalanceDruid.isRangedDpsSpec;
	readonly isMeleeDpsSpec = BalanceDruid.isMeleeDpsSpec;

	readonly canDualWield = BalanceDruid.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_starfall.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return BalanceDruid.getIcon(size);
	};
}

export class FeralDruid extends PlayerSpec<Spec.SpecFeralDruid> {
	static specIndex = 1;
	static specID = Spec.SpecFeralDruid as Spec.SpecFeralDruid;
	static classID = Class.ClassDruid as Class.ClassDruid;
	static friendlyName = 'Feral';
	static simLink = getSpecSiteUrl('druid', 'feral');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = false;

	readonly specIndex = FeralDruid.specIndex;
	readonly specID = FeralDruid.specID;
	readonly classID = FeralDruid.classID;
	readonly friendlyName = FeralDruid.friendlyName;
	readonly simLink = FeralDruid.simLink;

	readonly isTankSpec = FeralDruid.isTankSpec;
	readonly isHealingSpec = FeralDruid.isHealingSpec;
	readonly isRangedDpsSpec = FeralDruid.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FeralDruid.isMeleeDpsSpec;

	readonly canDualWield = FeralDruid.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_druid_catform.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return FeralDruid.getIcon(size);
	};
}

export class GuardianDruid extends PlayerSpec<Spec.SpecGuardianDruid> {
	static specIndex = 2;
	static specID = Spec.SpecGuardianDruid as Spec.SpecGuardianDruid;
	static classID = Class.ClassDruid as Class.ClassDruid;
	static friendlyName = 'Guardian';
	static simLink = getSpecSiteUrl('druid', 'guardian');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = false;

	readonly specIndex = GuardianDruid.specIndex;
	readonly specID = GuardianDruid.specID;
	readonly classID = GuardianDruid.classID;
	readonly friendlyName = GuardianDruid.friendlyName;
	readonly simLink = GuardianDruid.simLink;

	readonly isTankSpec = GuardianDruid.isTankSpec;
	readonly isHealingSpec = GuardianDruid.isHealingSpec;
	readonly isRangedDpsSpec = GuardianDruid.isRangedDpsSpec;
	readonly isMeleeDpsSpec = GuardianDruid.isMeleeDpsSpec;

	readonly canDualWield = GuardianDruid.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_racial_bearform.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return GuardianDruid.getIcon(size);
	};
}

export class RestorationDruid extends PlayerSpec<Spec.SpecRestorationDruid> {
	static specIndex = 3;
	static specID = Spec.SpecRestorationDruid as Spec.SpecRestorationDruid;
	static classID = Class.ClassDruid as Class.ClassDruid;
	static friendlyName = 'Restoration';
	static simLink = getSpecSiteUrl('druid', 'Restoration');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = RestorationDruid.specIndex;
	readonly specID = RestorationDruid.specID;
	readonly classID = RestorationDruid.classID;
	readonly friendlyName = RestorationDruid.friendlyName;
	readonly simLink = RestorationDruid.simLink;

	readonly isTankSpec = RestorationDruid.isTankSpec;
	readonly isHealingSpec = RestorationDruid.isHealingSpec;
	readonly isRangedDpsSpec = RestorationDruid.isRangedDpsSpec;
	readonly isMeleeDpsSpec = RestorationDruid.isMeleeDpsSpec;

	readonly canDualWield = RestorationDruid.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_healingtouch.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return RestorationDruid.getIcon(size);
	};
}
