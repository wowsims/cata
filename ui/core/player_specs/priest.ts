import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class DisciplinePriest extends PlayerSpec<Spec.SpecDisciplinePriest> {
	static specIndex = 0;
	static specID = Spec.SpecDisciplinePriest as Spec.SpecDisciplinePriest;
	static classID = Class.ClassPriest as Class.ClassPriest;
	static friendlyName = 'Discipline';
	static simLink = getSpecSiteUrl('priest', 'discipline');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = DisciplinePriest.specIndex;
	readonly specID = DisciplinePriest.specID;
	readonly classID = DisciplinePriest.classID;
	readonly friendlyName = DisciplinePriest.friendlyName;
	readonly simLink = DisciplinePriest.simLink;

	readonly isTankSpec = DisciplinePriest.isTankSpec;
	readonly isHealingSpec = DisciplinePriest.isHealingSpec;
	readonly isRangedDpsSpec = DisciplinePriest.isRangedDpsSpec;
	readonly isMeleeDpsSpec = DisciplinePriest.isMeleeDpsSpec;

	readonly canDualWield = DisciplinePriest.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_powerwordshield.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return DisciplinePriest.getIcon(size);
	};
}

export class HolyPriest extends PlayerSpec<Spec.SpecHolyPriest> {
	static specIndex = 1;
	static specID = Spec.SpecHolyPriest as Spec.SpecHolyPriest;
	static classID = Class.ClassPriest as Class.ClassPriest;
	static friendlyName = 'Holy';
	static simLink = getSpecSiteUrl('priest', 'holy');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = HolyPriest.specIndex;
	readonly specID = HolyPriest.specID;
	readonly classID = HolyPriest.classID;
	readonly friendlyName = HolyPriest.friendlyName;
	readonly simLink = HolyPriest.simLink;

	readonly isTankSpec = HolyPriest.isTankSpec;
	readonly isHealingSpec = HolyPriest.isHealingSpec;
	readonly isRangedDpsSpec = HolyPriest.isRangedDpsSpec;
	readonly isMeleeDpsSpec = HolyPriest.isMeleeDpsSpec;

	readonly canDualWield = HolyPriest.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_guardianspirit.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return HolyPriest.getIcon(size);
	};
}

export class ShadowPriest extends PlayerSpec<Spec.SpecShadowPriest> {
	static specIndex = 2;
	static specID = Spec.SpecShadowPriest as Spec.SpecShadowPriest;
	static classID = Class.ClassPriest as Class.ClassPriest;
	static friendlyName = 'Shadow';
	static simLink = getSpecSiteUrl('priest', 'shadow');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = ShadowPriest.specIndex;
	readonly specID = ShadowPriest.specID;
	readonly classID = ShadowPriest.classID;
	readonly friendlyName = ShadowPriest.friendlyName;
	readonly simLink = ShadowPriest.simLink;

	readonly isTankSpec = ShadowPriest.isTankSpec;
	readonly isHealingSpec = ShadowPriest.isHealingSpec;
	readonly isRangedDpsSpec = ShadowPriest.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ShadowPriest.isMeleeDpsSpec;

	readonly canDualWield = ShadowPriest.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_shadowwordpain.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return ShadowPriest.getIcon(size);
	};
}
