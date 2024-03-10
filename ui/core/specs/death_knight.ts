import { IconSize } from '../class';
import { DeathKnight } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class BloodDeathKnight extends Spec {
	static protoID = SpecProto.SpecBloodDeathKnight;
	static class = DeathKnight;
	static friendlyName = 'Blood';
	static simLink = getSpecSiteUrl('death_knight', 'blood');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = BloodDeathKnight.protoID;
	readonly class = BloodDeathKnight.class;
	readonly friendlyName = BloodDeathKnight.friendlyName;
	readonly simLink = BloodDeathKnight.simLink;

	readonly isTankSpec = BloodDeathKnight.isTankSpec;
	readonly isHealingSpec = BloodDeathKnight.isHealingSpec;
	readonly isRangedDpsSpec = BloodDeathKnight.isRangedDpsSpec;
	readonly isMeleeDpsSpec = BloodDeathKnight.isMeleeDpsSpec;

	readonly canDualWield = BloodDeathKnight.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_deathknight_bloodpresence.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return BloodDeathKnight.getIcon(size)
	}
}

export class FrostDeathKnight extends Spec {
	static protoID = SpecProto.SpecFrostDeathKnight;
	static class = DeathKnight;
	static friendlyName = 'Protection';
	static simLink = getSpecSiteUrl('death_knight', 'frost');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = FrostDeathKnight.protoID;
	readonly class = FrostDeathKnight.class;
	readonly friendlyName = FrostDeathKnight.friendlyName;
	readonly simLink = FrostDeathKnight.simLink;

	readonly isTankSpec = FrostDeathKnight.isTankSpec;
	readonly isHealingSpec = FrostDeathKnight.isHealingSpec;
	readonly isRangedDpsSpec = FrostDeathKnight.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FrostDeathKnight.isMeleeDpsSpec;

	readonly canDualWield = FrostDeathKnight.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_deathknight_frostpresence.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return FrostDeathKnight.getIcon(size)
	}
}

export class UnholyDeathKnight extends Spec {
	static protoID = SpecProto.SpecUnholyDeathKnight;
	static class = DeathKnight;
	static friendlyName = 'Unholy';
	static simLink = getSpecSiteUrl('death_knight', 'unholy');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = UnholyDeathKnight.protoID;
	readonly class = UnholyDeathKnight.class;
	readonly friendlyName = UnholyDeathKnight.friendlyName;
	readonly simLink = UnholyDeathKnight.simLink;

	readonly isTankSpec = UnholyDeathKnight.isTankSpec;
	readonly isHealingSpec = UnholyDeathKnight.isHealingSpec;
	readonly isRangedDpsSpec = UnholyDeathKnight.isRangedDpsSpec;
	readonly isMeleeDpsSpec = UnholyDeathKnight.isMeleeDpsSpec;

	readonly canDualWield = UnholyDeathKnight.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_deathknight_unholypresence.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return UnholyDeathKnight.getIcon(size)
	}
}
