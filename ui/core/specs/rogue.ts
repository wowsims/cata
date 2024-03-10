import { IconSize } from '../class';
import { Rogue } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class AssassinationRogue extends Spec {
	static protoID = SpecProto.SpecAssassinationRogue;
	static class = Rogue;
	static friendlyName = 'Assassination';
	static simLink = getSpecSiteUrl('rogue', 'assassination');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = AssassinationRogue.protoID;
	readonly class = AssassinationRogue.class;
	readonly friendlyName = AssassinationRogue.friendlyName;
	readonly simLink = AssassinationRogue.simLink;

	readonly isTankSpec = AssassinationRogue.isTankSpec;
	readonly isHealingSpec = AssassinationRogue.isHealingSpec;
	readonly isRangedDpsSpec = AssassinationRogue.isRangedDpsSpec;
	readonly isMeleeDpsSpec = AssassinationRogue.isMeleeDpsSpec;

	readonly canDualWield = AssassinationRogue.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_rogue_eviscerate.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return AssassinationRogue.getIcon(size)
	}
}

export class CombatRogue extends Spec {
	static protoID = SpecProto.SpecCombatRogue;
	static class = Rogue;
	static friendlyName = 'Combat';
	static simLink = getSpecSiteUrl('rogue', 'combat');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = CombatRogue.protoID;
	readonly class = CombatRogue.class;
	readonly friendlyName = CombatRogue.friendlyName;
	readonly simLink = CombatRogue.simLink;

	readonly isTankSpec = CombatRogue.isTankSpec;
	readonly isHealingSpec = CombatRogue.isHealingSpec;
	readonly isRangedDpsSpec = CombatRogue.isRangedDpsSpec;
	readonly isMeleeDpsSpec = CombatRogue.isMeleeDpsSpec;

	readonly canDualWield = CombatRogue.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_backstab.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return CombatRogue.getIcon(size)
	}
}

export class SubtletyRogue extends Spec {
	static protoID = SpecProto.SpecSubtletyRogue;
	static class = Rogue;
	static friendlyName = 'Subtlety';
	static simLink = getSpecSiteUrl('rogue', 'subtlety');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = SubtletyRogue.protoID;
	readonly class = SubtletyRogue.class;
	readonly friendlyName = SubtletyRogue.friendlyName;
	readonly simLink = SubtletyRogue.simLink;

	readonly isTankSpec = SubtletyRogue.isTankSpec;
	readonly isHealingSpec = SubtletyRogue.isHealingSpec;
	readonly isRangedDpsSpec = SubtletyRogue.isRangedDpsSpec;
	readonly isMeleeDpsSpec = SubtletyRogue.isMeleeDpsSpec;

	readonly canDualWield = SubtletyRogue.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_stealth.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return SubtletyRogue.getIcon(size)
	}
}
