import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class BloodDeathKnight extends PlayerSpec<Spec.SpecBloodDeathKnight> {
	static protoID = Spec.SpecBloodDeathKnight as Spec.SpecBloodDeathKnight;
	static playerClass = PlayerClasses.DeathKnight;
	static friendlyName = 'Blood';
	static fullName = `${this.friendlyName} ${PlayerClasses.DeathKnight.friendlyName}`;
	static simLink = getSpecSiteUrl('death_knight', 'blood');

	static isTankSpec = true;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = true;

	readonly protoID = BloodDeathKnight.protoID;
	readonly playerClass = BloodDeathKnight.playerClass;
	readonly friendlyName = BloodDeathKnight.friendlyName;
	readonly fullName = BloodDeathKnight.fullName;
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
	static protoID = Spec.SpecFrostDeathKnight as Spec.SpecFrostDeathKnight;
	static playerClass = PlayerClasses.DeathKnight;
	static friendlyName = 'Frost';
	static fullName = `${this.friendlyName} ${PlayerClasses.DeathKnight.friendlyName}`;
	static simLink = getSpecSiteUrl('death_knight', 'frost');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = FrostDeathKnight.protoID;
	readonly playerClass = FrostDeathKnight.playerClass;
	readonly friendlyName = FrostDeathKnight.friendlyName;
	readonly fullName = FrostDeathKnight.fullName;
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
	static protoID = Spec.SpecUnholyDeathKnight as Spec.SpecUnholyDeathKnight;
	static playerClass = PlayerClasses.DeathKnight;
	static friendlyName = 'Unholy';
	static fullName = `${this.friendlyName} ${PlayerClasses.DeathKnight.friendlyName}`;
	static simLink = getSpecSiteUrl('death_knight', 'unholy');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = UnholyDeathKnight.protoID;
	readonly playerClass = UnholyDeathKnight.playerClass;
	readonly friendlyName = UnholyDeathKnight.friendlyName;
	readonly fullName = UnholyDeathKnight.fullName;
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
