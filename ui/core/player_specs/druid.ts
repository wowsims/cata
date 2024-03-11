import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class BalanceDruid extends PlayerSpec<Spec.SpecBalanceDruid> {
	static protoID = Spec.SpecBalanceDruid as Spec.SpecBalanceDruid;
	static playerClass = PlayerClasses.Druid;
	static friendlyName = 'Balance';
	static fullName = `${this.friendlyName} ${PlayerClasses.Druid.friendlyName}`;
	static simLink = getSpecSiteUrl('druid', 'balance');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = BalanceDruid.protoID;
	readonly playerClass = BalanceDruid.playerClass;
	readonly friendlyName = BalanceDruid.friendlyName;
	readonly fullName = BalanceDruid.fullName;
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
	static protoID = Spec.SpecFeralDruid as Spec.SpecFeralDruid;
	static playerClass = PlayerClasses.Druid;
	static friendlyName = 'Feral';
	static fullName = `${this.friendlyName} ${PlayerClasses.Druid.friendlyName}`;
	static simLink = getSpecSiteUrl('druid', 'feral');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = false;

	readonly protoID = FeralDruid.protoID;
	readonly playerClass = FeralDruid.playerClass;
	readonly friendlyName = FeralDruid.friendlyName;
	readonly fullName = FeralDruid.fullName;
	readonly simLink = FeralDruid.simLink;

	readonly isTankSpec = FeralDruid.isTankSpec;
	readonly isHealingSpec = FeralDruid.isHealingSpec;
	readonly isRangedDpsSpec = FeralDruid.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FeralDruid.isMeleeDpsSpec;

	readonly canDualWield = FeralDruid.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_racial_bearform.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return FeralDruid.getIcon(size);
	};
}

export class RestorationDruid extends PlayerSpec<Spec.SpecRestorationDruid> {
	static protoID = Spec.SpecRestorationDruid as Spec.SpecRestorationDruid;
	static playerClass = PlayerClasses.Druid;
	static friendlyName = 'Restoration';
	static fullName = `${this.friendlyName} ${PlayerClasses.Druid.friendlyName}`;
	static simLink = getSpecSiteUrl('druid', 'Restoration');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = RestorationDruid.protoID;
	readonly playerClass = RestorationDruid.playerClass;
	readonly friendlyName = RestorationDruid.friendlyName;
	readonly fullName = RestorationDruid.fullName;
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
