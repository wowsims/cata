import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class BrewmasterMonk extends PlayerSpec<Spec.SpecBrewmasterMonk> {
	static specIndex = 0;
	static specID = Spec.SpecBrewmasterMonk as Spec.SpecBrewmasterMonk;
	static classID = Class.ClassMonk as Class.ClassMonk;
	static friendlyName = 'Brewmaster';
	static simLink = getSpecSiteUrl('monk', 'brewmaster');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly specIndex = BrewmasterMonk.specIndex;
	readonly specID = BrewmasterMonk.specID;
	readonly classID = BrewmasterMonk.classID;
	readonly friendlyName = BrewmasterMonk.friendlyName;
	readonly simLink = BrewmasterMonk.simLink;

	readonly isTankSpec = BrewmasterMonk.isTankSpec;
	readonly isHealingSpec = BrewmasterMonk.isHealingSpec;
	readonly isRangedDpsSpec = BrewmasterMonk.isRangedDpsSpec;
	readonly isMeleeDpsSpec = BrewmasterMonk.isMeleeDpsSpec;

	readonly canDualWield = BrewmasterMonk.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_monk_brewmaster_spec.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return BrewmasterMonk.getIcon(size);
	};
}

export class MistweaverMonk extends PlayerSpec<Spec.SpecMistweaverMonk> {
	static specIndex = 0;
	static specID = Spec.SpecMistweaverMonk as Spec.SpecMistweaverMonk;
	static classID = Class.ClassMonk as Class.ClassMonk;
	static friendlyName = 'Mistweaver';
	static simLink = getSpecSiteUrl('monk', 'mistweaver');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = MistweaverMonk.specIndex;
	readonly specID = MistweaverMonk.specID;
	readonly classID = MistweaverMonk.classID;
	readonly friendlyName = MistweaverMonk.friendlyName;
	readonly simLink = MistweaverMonk.simLink;

	readonly isTankSpec = MistweaverMonk.isTankSpec;
	readonly isHealingSpec = MistweaverMonk.isHealingSpec;
	readonly isRangedDpsSpec = MistweaverMonk.isRangedDpsSpec;
	readonly isMeleeDpsSpec = MistweaverMonk.isMeleeDpsSpec;

	readonly canDualWield = MistweaverMonk.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_monk_mistweaver_spec.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return MistweaverMonk.getIcon(size);
	};
}

export class WindwalkerMonk extends PlayerSpec<Spec.SpecWindwalkerMonk> {
	static specIndex = 0;
	static specID = Spec.SpecWindwalkerMonk as Spec.SpecWindwalkerMonk;
	static classID = Class.ClassMonk as Class.ClassMonk;
	static friendlyName = 'Windwalker';
	static simLink = getSpecSiteUrl('monk', 'windwalker');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = WindwalkerMonk.specIndex;
	readonly specID = WindwalkerMonk.specID;
	readonly classID = WindwalkerMonk.classID;
	readonly friendlyName = WindwalkerMonk.friendlyName;
	readonly simLink = WindwalkerMonk.simLink;

	readonly isTankSpec = WindwalkerMonk.isTankSpec;
	readonly isHealingSpec = WindwalkerMonk.isHealingSpec;
	readonly isRangedDpsSpec = WindwalkerMonk.isRangedDpsSpec;
	readonly isMeleeDpsSpec = WindwalkerMonk.isMeleeDpsSpec;

	readonly canDualWield = WindwalkerMonk.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_monk_windwalker_spec.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return WindwalkerMonk.getIcon(size);
	};
}
