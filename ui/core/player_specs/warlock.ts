import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class AfflictionWarlock extends PlayerSpec<Spec.SpecAfflictionWarlock> {
	static protoID = Spec.SpecAfflictionWarlock as Spec.SpecAfflictionWarlock;
	static playerClass = PlayerClasses.Warlock;
	static friendlyName = 'Affliction';
	static fullName = `${this.friendlyName} ${PlayerClasses.Warlock.friendlyName}`;
	static simLink = getSpecSiteUrl('warlock', 'affliction');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = AfflictionWarlock.protoID;
	readonly playerClass = AfflictionWarlock.playerClass;
	readonly friendlyName = AfflictionWarlock.friendlyName;
	readonly fullName = AfflictionWarlock.fullName;
	readonly simLink = AfflictionWarlock.simLink;

	readonly isTankSpec = AfflictionWarlock.isTankSpec;
	readonly isHealingSpec = AfflictionWarlock.isHealingSpec;
	readonly isRangedDpsSpec = AfflictionWarlock.isRangedDpsSpec;
	readonly isMeleeDpsSpec = AfflictionWarlock.isMeleeDpsSpec;

	readonly canDualWield = AfflictionWarlock.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_deathcoil.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return AfflictionWarlock.getIcon(size);
	};
}

export class DemonologyWarlock extends PlayerSpec<Spec.SpecDemonologyWarlock> {
	static protoID = Spec.SpecDemonologyWarlock as Spec.SpecDemonologyWarlock;
	static playerClass = PlayerClasses.Warlock;
	static friendlyName = 'Demonology';
	static fullName = `${this.friendlyName} ${PlayerClasses.Warlock.friendlyName}`;
	static simLink = getSpecSiteUrl('warlock', 'demonology');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = DemonologyWarlock.protoID;
	readonly playerClass = DemonologyWarlock.playerClass;
	readonly friendlyName = DemonologyWarlock.friendlyName;
	readonly fullName = DemonologyWarlock.fullName;
	readonly simLink = DemonologyWarlock.simLink;

	readonly isTankSpec = DemonologyWarlock.isTankSpec;
	readonly isHealingSpec = DemonologyWarlock.isHealingSpec;
	readonly isRangedDpsSpec = DemonologyWarlock.isRangedDpsSpec;
	readonly isMeleeDpsSpec = DemonologyWarlock.isMeleeDpsSpec;

	readonly canDualWield = DemonologyWarlock.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_metamorphosis.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return DemonologyWarlock.getIcon(size);
	};
}

export class DestructionWarlock extends PlayerSpec<Spec.SpecDestructionWarlock> {
	static protoID = Spec.SpecDestructionWarlock as Spec.SpecDestructionWarlock;
	static playerClass = PlayerClasses.Warlock;
	static friendlyName = 'Destruction';
	static fullName = `${this.friendlyName} ${PlayerClasses.Warlock.friendlyName}`;
	static simLink = getSpecSiteUrl('warlock', 'destruction');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = DestructionWarlock.protoID;
	readonly playerClass = DestructionWarlock.playerClass;
	readonly friendlyName = DestructionWarlock.friendlyName;
	readonly fullName = DestructionWarlock.fullName;
	readonly simLink = DestructionWarlock.simLink;

	readonly isTankSpec = DestructionWarlock.isTankSpec;
	readonly isHealingSpec = DestructionWarlock.isHealingSpec;
	readonly isRangedDpsSpec = DestructionWarlock.isRangedDpsSpec;
	readonly isMeleeDpsSpec = DestructionWarlock.isMeleeDpsSpec;

	readonly canDualWield = DestructionWarlock.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_rainoffire.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return DestructionWarlock.getIcon(size);
	};
}
