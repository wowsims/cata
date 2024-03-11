import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class AssassinationRogue extends PlayerSpec<Spec.SpecAssassinationRogue> {
	static protoID = Spec.SpecAssassinationRogue as Spec.SpecAssassinationRogue;
	static playerClass = PlayerClasses.Rogue;
	static friendlyName = 'Assassination';
	static fullName = `${this.friendlyName} ${PlayerClasses.Rogue.friendlyName}`;
	static simLink = getSpecSiteUrl('rogue', 'assassination');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = AssassinationRogue.protoID;
	readonly playerClass = AssassinationRogue.playerClass;
	readonly friendlyName = AssassinationRogue.friendlyName;
	readonly fullName = AssassinationRogue.fullName;
	readonly simLink = AssassinationRogue.simLink;

	readonly isTankSpec = AssassinationRogue.isTankSpec;
	readonly isHealingSpec = AssassinationRogue.isHealingSpec;
	readonly isRangedDpsSpec = AssassinationRogue.isRangedDpsSpec;
	readonly isMeleeDpsSpec = AssassinationRogue.isMeleeDpsSpec;

	readonly canDualWield = AssassinationRogue.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_rogue_eviscerate.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return AssassinationRogue.getIcon(size);
	};
}

export class CombatRogue extends PlayerSpec<Spec.SpecCombatRogue> {
	static protoID = Spec.SpecCombatRogue as Spec.SpecCombatRogue;
	static playerClass = PlayerClasses.Rogue;
	static friendlyName = 'Combat';
	static fullName = `${this.friendlyName} ${PlayerClasses.Rogue.friendlyName}`;
	static simLink = getSpecSiteUrl('rogue', 'combat');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = CombatRogue.protoID;
	readonly playerClass = CombatRogue.playerClass;
	readonly friendlyName = CombatRogue.friendlyName;
	readonly fullName = CombatRogue.fullName;
	readonly simLink = CombatRogue.simLink;

	readonly isTankSpec = CombatRogue.isTankSpec;
	readonly isHealingSpec = CombatRogue.isHealingSpec;
	readonly isRangedDpsSpec = CombatRogue.isRangedDpsSpec;
	readonly isMeleeDpsSpec = CombatRogue.isMeleeDpsSpec;

	readonly canDualWield = CombatRogue.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_backstab.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return CombatRogue.getIcon(size);
	};
}

export class SubtletyRogue extends PlayerSpec<Spec.SpecSubtletyRogue> {
	static protoID = Spec.SpecSubtletyRogue as Spec.SpecSubtletyRogue;
	static playerClass = PlayerClasses.Rogue;
	static friendlyName = 'Subtlety';
	static fullName = `${this.friendlyName} ${PlayerClasses.Rogue.friendlyName}`;
	static simLink = getSpecSiteUrl('rogue', 'subtlety');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = SubtletyRogue.protoID;
	readonly playerClass = SubtletyRogue.playerClass;
	readonly friendlyName = SubtletyRogue.friendlyName;
	readonly fullName = SubtletyRogue.fullName;
	readonly simLink = SubtletyRogue.simLink;

	readonly isTankSpec = SubtletyRogue.isTankSpec;
	readonly isHealingSpec = SubtletyRogue.isHealingSpec;
	readonly isRangedDpsSpec = SubtletyRogue.isRangedDpsSpec;
	readonly isMeleeDpsSpec = SubtletyRogue.isMeleeDpsSpec;

	readonly canDualWield = SubtletyRogue.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_stealth.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return SubtletyRogue.getIcon(size);
	};
}
