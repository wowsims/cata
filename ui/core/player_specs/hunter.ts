import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class BeastMasteryHunter extends PlayerSpec<Spec.SpecBeastMasteryHunter> {
	static protoID = Spec.SpecBeastMasteryHunter as Spec.SpecBeastMasteryHunter;
	static playerClass = PlayerClasses.Hunter;
	static friendlyName = 'Beast Mastery';
	static fullName = `${this.friendlyName} ${PlayerClasses.Hunter.friendlyName}`;
	static simLink = getSpecSiteUrl('hunter', 'beast_mastery');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = BeastMasteryHunter.protoID;
	readonly playerClass = BeastMasteryHunter.playerClass;
	readonly friendlyName = BeastMasteryHunter.friendlyName;
	readonly fullName = BeastMasteryHunter.fullName;
	readonly simLink = BeastMasteryHunter.simLink;

	readonly isTankSpec = BeastMasteryHunter.isTankSpec;
	readonly isHealingSpec = BeastMasteryHunter.isHealingSpec;
	readonly isRangedDpsSpec = BeastMasteryHunter.isRangedDpsSpec;
	readonly isMeleeDpsSpec = BeastMasteryHunter.isMeleeDpsSpec;

	readonly canDualWield = BeastMasteryHunter.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_hunter_bestialdiscipline.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return BeastMasteryHunter.getIcon(size);
	};
}

export class MarksmanshipHunter extends PlayerSpec<Spec.SpecMarksmanshipHunter> {
	static protoID = Spec.SpecMarksmanshipHunter as Spec.SpecMarksmanshipHunter;
	static playerClass = PlayerClasses.Hunter;
	static friendlyName = 'Marksmanship';
	static fullName = `${this.friendlyName} ${PlayerClasses.Hunter.friendlyName}`;
	static simLink = getSpecSiteUrl('hunter', 'marksmanship');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = MarksmanshipHunter.protoID;
	readonly playerClass = MarksmanshipHunter.playerClass;
	readonly friendlyName = MarksmanshipHunter.friendlyName;
	readonly fullName = MarksmanshipHunter.fullName;
	readonly simLink = MarksmanshipHunter.simLink;

	readonly isTankSpec = MarksmanshipHunter.isTankSpec;
	readonly isHealingSpec = MarksmanshipHunter.isHealingSpec;
	readonly isRangedDpsSpec = MarksmanshipHunter.isRangedDpsSpec;
	readonly isMeleeDpsSpec = MarksmanshipHunter.isMeleeDpsSpec;

	readonly canDualWield = MarksmanshipHunter.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_hunter_focusedaim.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return MarksmanshipHunter.getIcon(size);
	};
}

export class SurvivalHunter extends PlayerSpec<Spec.SpecSurvivalHunter> {
	static protoID = Spec.SpecSurvivalHunter as Spec.SpecSurvivalHunter;
	static playerClass = PlayerClasses.Hunter;
	static friendlyName = 'Survival';
	static fullName = `${this.friendlyName} ${PlayerClasses.Hunter.friendlyName}`;
	static simLink = getSpecSiteUrl('hunter', 'survival');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = SurvivalHunter.protoID;
	readonly playerClass = SurvivalHunter.playerClass;
	readonly friendlyName = SurvivalHunter.friendlyName;
	readonly fullName = SurvivalHunter.fullName;
	readonly simLink = SurvivalHunter.simLink;

	readonly isTankSpec = SurvivalHunter.isTankSpec;
	readonly isHealingSpec = SurvivalHunter.isHealingSpec;
	readonly isRangedDpsSpec = SurvivalHunter.isRangedDpsSpec;
	readonly isMeleeDpsSpec = SurvivalHunter.isMeleeDpsSpec;

	readonly canDualWield = SurvivalHunter.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/ability_hunter_camouflage.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return SurvivalHunter.getIcon(size);
	};
}
