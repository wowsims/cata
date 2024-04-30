import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class ArmsWarrior extends PlayerSpec<Spec.SpecArmsWarrior> {
	static specIndex = 0;
	static specID = Spec.SpecArmsWarrior as Spec.SpecArmsWarrior;
	static classID = Class.ClassWarrior as Class.ClassWarrior;
	static friendlyName = 'Arms';
	static simLink = getSpecSiteUrl('warrior', 'arms');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = ArmsWarrior.specIndex;
	readonly specID = ArmsWarrior.specID;
	readonly classID = ArmsWarrior.classID;
	readonly friendlyName = ArmsWarrior.friendlyName;
	readonly simLink = ArmsWarrior.simLink;

	readonly isTankSpec = ArmsWarrior.isTankSpec;
	readonly isHealingSpec = ArmsWarrior.isHealingSpec;
	readonly isRangedDpsSpec = ArmsWarrior.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ArmsWarrior.isMeleeDpsSpec;

	readonly canDualWield = ArmsWarrior.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_warrior_savageblow.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return ArmsWarrior.getIcon(size);
	};
}

export class FuryWarrior extends PlayerSpec<Spec.SpecFuryWarrior> {
	static specIndex = 1;
	static specID = Spec.SpecFuryWarrior as Spec.SpecFuryWarrior;
	static classID = Class.ClassWarrior as Class.ClassWarrior;
	static friendlyName = 'Fury';
	static simLink = getSpecSiteUrl('warrior', 'fury');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = FuryWarrior.specIndex;
	readonly specID = FuryWarrior.specID;
	readonly classID = FuryWarrior.classID;
	readonly friendlyName = FuryWarrior.friendlyName;
	readonly simLink = FuryWarrior.simLink;

	readonly isTankSpec = FuryWarrior.isTankSpec;
	readonly isHealingSpec = FuryWarrior.isHealingSpec;
	readonly isRangedDpsSpec = FuryWarrior.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FuryWarrior.isMeleeDpsSpec;

	readonly canDualWield = FuryWarrior.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_warrior_innerrage.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return FuryWarrior.getIcon(size);
	};
}

export class ProtectionWarrior extends PlayerSpec<Spec.SpecProtectionWarrior> {
	static specIndex = 2;
	static specID = Spec.SpecProtectionWarrior as Spec.SpecProtectionWarrior;
	static classID = Class.ClassWarrior as Class.ClassWarrior;
	static friendlyName = 'Protection';
	static simLink = getSpecSiteUrl('warrior', 'protection');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly specIndex = ProtectionWarrior.specIndex;
	readonly specID = ProtectionWarrior.specID;
	readonly classID = ProtectionWarrior.classID;
	readonly friendlyName = ProtectionWarrior.friendlyName;
	readonly simLink = ProtectionWarrior.simLink;

	readonly isTankSpec = ProtectionWarrior.isTankSpec;
	readonly isHealingSpec = ProtectionWarrior.isHealingSpec;
	readonly isRangedDpsSpec = ProtectionWarrior.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ProtectionWarrior.isMeleeDpsSpec;

	readonly canDualWield = ProtectionWarrior.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_warrior_defensivestance.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return ProtectionWarrior.getIcon(size);
	};
}
