import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class AfflictionWarlock extends PlayerSpec<Spec.SpecAfflictionWarlock> {
	static specIndex = 1;
	static specID = Spec.SpecAfflictionWarlock as Spec.SpecAfflictionWarlock;
	static classID = Class.ClassWarlock as Class.ClassWarlock;
	static friendlyName = 'Affliction';
	static simLink = getSpecSiteUrl('warlock', 'affliction');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = AfflictionWarlock.specIndex;
	readonly specID = AfflictionWarlock.specID;
	readonly classID = AfflictionWarlock.classID;
	readonly friendlyName = AfflictionWarlock.friendlyName;
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
	static specIndex = 1;
	static specID = Spec.SpecDemonologyWarlock as Spec.SpecDemonologyWarlock;
	static classID = Class.ClassWarlock as Class.ClassWarlock;
	static friendlyName = 'Demonology';
	static simLink = getSpecSiteUrl('warlock', 'demonology');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = DemonologyWarlock.specIndex;
	readonly specID = DemonologyWarlock.specID;
	readonly classID = DemonologyWarlock.classID;
	readonly friendlyName = DemonologyWarlock.friendlyName;
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
	static specIndex = 2;
	static specID = Spec.SpecDestructionWarlock as Spec.SpecDestructionWarlock;
	static classID = Class.ClassWarlock as Class.ClassWarlock;
	static friendlyName = 'Destruction';
	static simLink = getSpecSiteUrl('warlock', 'destruction');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = DestructionWarlock.specIndex;
	readonly specID = DestructionWarlock.specID;
	readonly classID = DestructionWarlock.classID;
	readonly friendlyName = DestructionWarlock.friendlyName;
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
