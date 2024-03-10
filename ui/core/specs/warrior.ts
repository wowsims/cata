import { IconSize } from '../class';
import { Warrior } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class ArmsWarrior extends Spec {
	static protoID = SpecProto.SpecArmsWarrior;
	static class = Warrior;
	static friendlyName = 'Arms';
	static simLink = getSpecSiteUrl('warrior', 'arms');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = ArmsWarrior.protoID;
	readonly class = ArmsWarrior.class;
	readonly friendlyName = ArmsWarrior.friendlyName;
	readonly simLink = ArmsWarrior.simLink;

	readonly isTankSpec = ArmsWarrior.isTankSpec;
	readonly isHealingSpec = ArmsWarrior.isHealingSpec;
	readonly isRangedDpsSpec = ArmsWarrior.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ArmsWarrior.isMeleeDpsSpec;

	readonly canDualWield = ArmsWarrior.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_warrior_savageblow.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return ArmsWarrior.getIcon(size)
	}
}

export class FuryWarrior extends Spec {
	static protoID = SpecProto.SpecFuryWarrior;
	static class = Warrior;
	static friendlyName = 'Fury';
	static simLink = getSpecSiteUrl('warrior', 'fury');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = FuryWarrior.protoID;
	readonly class = FuryWarrior.class;
	readonly friendlyName = FuryWarrior.friendlyName;
	readonly simLink = FuryWarrior.simLink;

	readonly isTankSpec = FuryWarrior.isTankSpec;
	readonly isHealingSpec = FuryWarrior.isHealingSpec;
	readonly isRangedDpsSpec = FuryWarrior.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FuryWarrior.isMeleeDpsSpec;

	readonly canDualWield = FuryWarrior.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_warrior_innerrage.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return FuryWarrior.getIcon(size)
	}
}

export class ProtectionWarrior extends Spec {
	static protoID = SpecProto.SpecProtectionWarrior;
	static class = Warrior;
	static friendlyName = 'Protection';
	static simLink = getSpecSiteUrl('warrior', 'protection');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = ProtectionWarrior.protoID;
	readonly class = ProtectionWarrior.class;
	readonly friendlyName = ProtectionWarrior.friendlyName;
	readonly simLink = ProtectionWarrior.simLink;

	readonly isTankSpec = ProtectionWarrior.isTankSpec;
	readonly isHealingSpec = ProtectionWarrior.isHealingSpec;
	readonly isRangedDpsSpec = ProtectionWarrior.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ProtectionWarrior.isMeleeDpsSpec;

	readonly canDualWield = ProtectionWarrior.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_warrior_defensivestance.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return ProtectionWarrior.getIcon(size)
	}
}
