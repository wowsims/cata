import { IconSize } from '../class';
import { Priest } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class DisciplinePriest extends Spec {
	static protoID = SpecProto.SpecDisciplinePriest;
	static class = Priest;
	static friendlyName = 'Discipline';
	static simLink = getSpecSiteUrl('priest', 'discipline');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = DisciplinePriest.protoID;
	readonly class = DisciplinePriest.class;
	readonly friendlyName = DisciplinePriest.friendlyName;
	readonly simLink = DisciplinePriest.simLink;

	readonly isTankSpec = DisciplinePriest.isTankSpec;
	readonly isHealingSpec = DisciplinePriest.isHealingSpec;
	readonly isRangedDpsSpec = DisciplinePriest.isRangedDpsSpec;
	readonly isMeleeDpsSpec = DisciplinePriest.isMeleeDpsSpec;

	readonly canDualWield = DisciplinePriest.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_powerwordshield.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return DisciplinePriest.getIcon(size)
	}
}

export class HolyPriest extends Spec {
	static protoID = SpecProto.SpecHolyPriest;
	static class = Priest;
	static friendlyName = 'Holy';
	static simLink = getSpecSiteUrl('priest', 'holy');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = HolyPriest.protoID;
	readonly class = HolyPriest.class;
	readonly friendlyName = HolyPriest.friendlyName;
	readonly simLink = HolyPriest.simLink;

	readonly isTankSpec = HolyPriest.isTankSpec;
	readonly isHealingSpec = HolyPriest.isHealingSpec;
	readonly isRangedDpsSpec = HolyPriest.isRangedDpsSpec;
	readonly isMeleeDpsSpec = HolyPriest.isMeleeDpsSpec;

	readonly canDualWield = HolyPriest.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_guardianspirit.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return HolyPriest.getIcon(size)
	}
}

export class ShadowPriest extends Spec {
	static protoID = SpecProto.SpecShadowPriest;
	static class = Priest;
	static friendlyName = 'Shadow';
	static simLink = getSpecSiteUrl('priest', 'shadow');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ShadowPriest.protoID;
	readonly class = ShadowPriest.class;
	readonly friendlyName = ShadowPriest.friendlyName;
	readonly simLink = ShadowPriest.simLink;

	readonly isTankSpec = ShadowPriest.isTankSpec;
	readonly isHealingSpec = ShadowPriest.isHealingSpec;
	readonly isRangedDpsSpec = ShadowPriest.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ShadowPriest.isMeleeDpsSpec;

	readonly canDualWield = ShadowPriest.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_shadowwordpain.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return ShadowPriest.getIcon(size)
	}
}
