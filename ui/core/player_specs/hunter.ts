import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class BeastMasteryHunter extends PlayerSpec<Spec.SpecBeastMasteryHunter> {
	static specIndex = 0;
	static specID = Spec.SpecBeastMasteryHunter as Spec.SpecBeastMasteryHunter;
	static classID = Class.ClassHunter as Class.ClassHunter;
	static friendlyName = 'Beast Mastery';
	static simLink = getSpecSiteUrl('hunter', 'beast_mastery');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly specIndex = BeastMasteryHunter.specIndex;
	readonly specID = BeastMasteryHunter.specID;
	readonly classID = BeastMasteryHunter.classID;
	readonly friendlyName = BeastMasteryHunter.friendlyName;
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
	static specIndex = 1;
	static specID = Spec.SpecMarksmanshipHunter as Spec.SpecMarksmanshipHunter;
	static classID = Class.ClassHunter as Class.ClassHunter;
	static friendlyName = 'Marksmanship';
	static simLink = getSpecSiteUrl('hunter', 'marksmanship');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly specIndex = MarksmanshipHunter.specIndex;
	readonly specID = MarksmanshipHunter.specID;
	readonly classID = MarksmanshipHunter.classID;
	readonly friendlyName = MarksmanshipHunter.friendlyName;
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
	static specIndex = 2;
	static specID = Spec.SpecSurvivalHunter as Spec.SpecSurvivalHunter;
	static classID = Class.ClassHunter as Class.ClassHunter;
	static friendlyName = 'Survival';
	static simLink = getSpecSiteUrl('hunter', 'survival');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly specIndex = SurvivalHunter.specIndex;
	readonly specID = SurvivalHunter.specID;
	readonly classID = SurvivalHunter.classID;
	readonly friendlyName = SurvivalHunter.friendlyName;
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
