import { IconSize } from '../player_class';
import { Priest } from '../player_classes/priest';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class DisciplinePriest extends PlayerSpec<Spec.SpecDisciplinePriest> {
	static protoID = Spec.SpecDisciplinePriest as Spec.SpecDisciplinePriest;
	static playerClass = Priest;
	static friendlyName = 'Discipline';
	static fullName = `${this.friendlyName} ${Priest.friendlyName}`;
	static simLink = getSpecSiteUrl('priest', 'discipline');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = DisciplinePriest.protoID;
	readonly playerClass = DisciplinePriest.playerClass;
	readonly friendlyName = DisciplinePriest.friendlyName;
	readonly fullName = DisciplinePriest.fullName;
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
	static protoID = Spec.SpecHolyPriest as Spec.SpecHolyPriest;
	static playerClass = Priest;
	static friendlyName = 'Holy';
	static fullName = `${this.friendlyName} ${Priest.friendlyName}`;
	static simLink = getSpecSiteUrl('priest', 'holy');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = HolyPriest.protoID;
	readonly playerClass = HolyPriest.playerClass;
	readonly friendlyName = HolyPriest.friendlyName;
	readonly fullName = HolyPriest.fullName;
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
	static protoID = Spec.SpecShadowPriest as Spec.SpecShadowPriest;
	static playerClass = Priest;
	static friendlyName = 'Shadow';
	static fullName = `${this.friendlyName} ${Priest.friendlyName}`;
	static simLink = getSpecSiteUrl('priest', 'shadow');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ShadowPriest.protoID;
	readonly playerClass = ShadowPriest.playerClass;
	readonly friendlyName = ShadowPriest.friendlyName;
	readonly fullName = ShadowPriest.fullName;
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
