import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class ArmsWarrior extends PlayerSpec<Spec.SpecArmsWarrior> {
	static protoID = Spec.SpecArmsWarrior as Spec.SpecArmsWarrior;
	static playerClass = PlayerClasses.Warrior;
	static friendlyName = 'Arms';
	static fullName = `${this.friendlyName} ${PlayerClasses.Warrior.friendlyName}`;
	static simLink = getSpecSiteUrl('warrior', 'arms');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = ArmsWarrior.protoID;
	readonly playerClass = ArmsWarrior.playerClass;
	readonly friendlyName = ArmsWarrior.friendlyName;
	readonly fullName = ArmsWarrior.fullName;
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
	static protoID = Spec.SpecFuryWarrior as Spec.SpecFuryWarrior;
	static playerClass = PlayerClasses.Warrior;
	static friendlyName = 'Fury';
	static fullName = `${this.friendlyName} ${PlayerClasses.Warrior.friendlyName}`;
	static simLink = getSpecSiteUrl('warrior', 'fury');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = FuryWarrior.protoID;
	readonly playerClass = FuryWarrior.playerClass;
	readonly friendlyName = FuryWarrior.friendlyName;
	readonly fullName = FuryWarrior.fullName;
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
	static protoID = Spec.SpecProtectionWarrior as Spec.SpecProtectionWarrior;
	static playerClass = PlayerClasses.Warrior;
	static friendlyName = 'Protection';
	static fullName = `${this.friendlyName} ${PlayerClasses.Warrior.friendlyName}`;
	static simLink = getSpecSiteUrl('warrior', 'protection');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = ProtectionWarrior.protoID;
	readonly playerClass = ProtectionWarrior.playerClass;
	readonly friendlyName = ProtectionWarrior.friendlyName;
	readonly fullName = ProtectionWarrior.fullName;
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
