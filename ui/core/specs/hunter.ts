import { IconSize } from '../class';
import { Hunter } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class BeastMasteryHunter extends Spec {
	static protoID = SpecProto.SpecBeastMasteryHunter;
	static class = Hunter;
	static friendlyName = 'Beast Mastery';
	static simLink = getSpecSiteUrl('hunter', 'beast_mastery');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = BeastMasteryHunter.protoID;
	readonly class = BeastMasteryHunter.class;
	readonly friendlyName = BeastMasteryHunter.friendlyName;
	readonly simLink = BeastMasteryHunter.simLink;

	readonly isTankSpec = BeastMasteryHunter.isTankSpec;
	readonly isHealingSpec = BeastMasteryHunter.isHealingSpec;
	readonly isRangedDpsSpec = BeastMasteryHunter.isRangedDpsSpec;
	readonly isMeleeDpsSpec = BeastMasteryHunter.isMeleeDpsSpec;

	readonly canDualWield = BeastMasteryHunter.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_hunter_bestialdiscipline.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return BeastMasteryHunter.getIcon(size)
	}
}

export class MarksmanshipHunter extends Spec {
	static protoID = SpecProto.SpecMarksmanshipHunter;
	static class = Hunter;
	static friendlyName = 'Marksmanship';
	static simLink = getSpecSiteUrl('hunter', 'marksmanship');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false

	static canDualWield = true;

	readonly protoID = MarksmanshipHunter.protoID;
	readonly class = MarksmanshipHunter.class;
	readonly friendlyName = MarksmanshipHunter.friendlyName;
	readonly simLink = MarksmanshipHunter.simLink;

	readonly isTankSpec = MarksmanshipHunter.isTankSpec;
	readonly isHealingSpec = MarksmanshipHunter.isHealingSpec;
	readonly isRangedDpsSpec = MarksmanshipHunter.isRangedDpsSpec;
	readonly isMeleeDpsSpec = MarksmanshipHunter.isMeleeDpsSpec;

	readonly canDualWield = MarksmanshipHunter.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_hunter_focusedaim.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return MarksmanshipHunter.getIcon(size)
	}
}

export class SurvivalHunter extends Spec {
	static protoID = SpecProto.SpecSurvivalHunter;
	static class = Hunter;
	static friendlyName = 'Survival';
	static simLink = getSpecSiteUrl('hunter', 'survival');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = SurvivalHunter.protoID;
	readonly class = SurvivalHunter.class;
	readonly friendlyName = SurvivalHunter.friendlyName;
	readonly simLink = SurvivalHunter.simLink;

	readonly isTankSpec = SurvivalHunter.isTankSpec;
	readonly isHealingSpec = SurvivalHunter.isHealingSpec;
	readonly isRangedDpsSpec = SurvivalHunter.isRangedDpsSpec;
	readonly isMeleeDpsSpec = SurvivalHunter.isMeleeDpsSpec;

	readonly canDualWield = SurvivalHunter.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_hunter_camouflage.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return SurvivalHunter.getIcon(size)
	}
}
