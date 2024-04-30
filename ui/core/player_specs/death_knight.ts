import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class BloodDeathKnight extends PlayerSpec<Spec.SpecBloodDeathKnight> {
	static specIndex = 0;
	static specID = Spec.SpecBloodDeathKnight as Spec.SpecBloodDeathKnight;
	static classID = Class.ClassDeathKnight as Class.ClassDeathKnight;
	static friendlyName = 'Blood';
	static simLink = getSpecSiteUrl('death_knight', 'blood');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly specIndex = BloodDeathKnight.specIndex;
	readonly specID = BloodDeathKnight.specID;
	readonly classID = BloodDeathKnight.classID;
	readonly friendlyName = BloodDeathKnight.friendlyName;
	readonly simLink = BloodDeathKnight.simLink;

	readonly isTankSpec = BloodDeathKnight.isTankSpec;
	readonly isHealingSpec = BloodDeathKnight.isHealingSpec;
	readonly isRangedDpsSpec = BloodDeathKnight.isRangedDpsSpec;
	readonly isMeleeDpsSpec = BloodDeathKnight.isMeleeDpsSpec;

	readonly canDualWield = BloodDeathKnight.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_deathknight_bloodpresence.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return BloodDeathKnight.getIcon(size);
	};
}

export class FrostDeathKnight extends PlayerSpec<Spec.SpecFrostDeathKnight> {
	static specIndex = 1;
	static specID = Spec.SpecFrostDeathKnight as Spec.SpecFrostDeathKnight;
	static classID = Class.ClassDeathKnight as Class.ClassDeathKnight;
	static friendlyName = 'Frost';
	static simLink = getSpecSiteUrl('death_knight', 'frost');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = FrostDeathKnight.specIndex;
	readonly specID = FrostDeathKnight.specID;
	readonly classID = FrostDeathKnight.classID;
	readonly friendlyName = FrostDeathKnight.friendlyName;
	readonly simLink = FrostDeathKnight.simLink;

	readonly isTankSpec = FrostDeathKnight.isTankSpec;
	readonly isHealingSpec = FrostDeathKnight.isHealingSpec;
	readonly isRangedDpsSpec = FrostDeathKnight.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FrostDeathKnight.isMeleeDpsSpec;

	readonly canDualWield = FrostDeathKnight.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_deathknight_frostpresence.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return FrostDeathKnight.getIcon(size);
	};
}

export class UnholyDeathKnight extends PlayerSpec<Spec.SpecUnholyDeathKnight> {
	static specIndex = 2;
	static specID = Spec.SpecUnholyDeathKnight as Spec.SpecUnholyDeathKnight;
	static classID = Class.ClassDeathKnight as Class.ClassDeathKnight;
	static friendlyName = 'Unholy';
	static simLink = getSpecSiteUrl('death_knight', 'unholy');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly specIndex = UnholyDeathKnight.specIndex;
	readonly specID = UnholyDeathKnight.specID;
	readonly classID = UnholyDeathKnight.classID;
	readonly friendlyName = UnholyDeathKnight.friendlyName;
	readonly simLink = UnholyDeathKnight.simLink;

	readonly isTankSpec = UnholyDeathKnight.isTankSpec;
	readonly isHealingSpec = UnholyDeathKnight.isHealingSpec;
	readonly isRangedDpsSpec = UnholyDeathKnight.isRangedDpsSpec;
	readonly isMeleeDpsSpec = UnholyDeathKnight.isMeleeDpsSpec;

	readonly canDualWield = UnholyDeathKnight.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_deathknight_unholypresence.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return UnholyDeathKnight.getIcon(size);
	};
}
