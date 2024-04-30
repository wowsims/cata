import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class AssassinationRogue extends PlayerSpec<Spec.SpecAssassinationRogue> {
	static specIndex = 0;
	static specID = Spec.SpecAssassinationRogue as Spec.SpecAssassinationRogue;
	static classID = Class.ClassRogue as Class.ClassRogue;
	static friendlyName = 'Assassination';
	static simLink = getSpecSiteUrl('rogue', 'assassination');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = AssassinationRogue.specIndex;
	readonly specID = AssassinationRogue.specID;
	readonly classID = AssassinationRogue.classID;
	readonly friendlyName = AssassinationRogue.friendlyName;
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
	static specIndex = 1;
	static specID = Spec.SpecCombatRogue as Spec.SpecCombatRogue;
	static classID = Class.ClassRogue as Class.ClassRogue;
	static friendlyName = 'Combat';
	static simLink = getSpecSiteUrl('rogue', 'combat');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = CombatRogue.specIndex;
	readonly specID = CombatRogue.specID;
	readonly classID = CombatRogue.classID;
	readonly friendlyName = CombatRogue.friendlyName;
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
	static specIndex = 2;
	static specID = Spec.SpecSubtletyRogue as Spec.SpecSubtletyRogue;
	static classID = Class.ClassRogue as Class.ClassRogue;
	static friendlyName = 'Subtlety';
	static simLink = getSpecSiteUrl('rogue', 'subtlety');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = SubtletyRogue.specIndex;
	readonly specID = SubtletyRogue.specID;
	readonly classID = SubtletyRogue.classID;
	readonly friendlyName = SubtletyRogue.friendlyName;
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
