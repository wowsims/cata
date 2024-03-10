import { IconSize } from '../class';
import { Warlock } from '../classes';
import { Spec as SpecProto } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';
import { Spec } from '../spec';

export class AfflictionWarlock extends Spec {
	static protoID = SpecProto.SpecAfflictionWarlock;
	static class = Warlock;
	static friendlyName = 'Affliction';
	static simLink = getSpecSiteUrl('warlock', 'affliction');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = AfflictionWarlock.protoID;
	readonly class = AfflictionWarlock.class;
	readonly friendlyName = AfflictionWarlock.friendlyName;
	readonly simLink = AfflictionWarlock.simLink;

	readonly isTankSpec = AfflictionWarlock.isTankSpec;
	readonly isHealingSpec = AfflictionWarlock.isHealingSpec;
	readonly isRangedDpsSpec = AfflictionWarlock.isRangedDpsSpec;
	readonly isMeleeDpsSpec = AfflictionWarlock.isMeleeDpsSpec;

	readonly canDualWield = AfflictionWarlock.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_deathcoil.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return AfflictionWarlock.getIcon(size)
	}
}

export class DemonologyWarlock extends Spec {
	static protoID = SpecProto.SpecDemonologyWarlock;
	static class = Warlock;
	static friendlyName = 'Demonology';
	static simLink = getSpecSiteUrl('warlock', 'demonology');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = DemonologyWarlock.protoID;
	readonly class = DemonologyWarlock.class;
	readonly friendlyName = DemonologyWarlock.friendlyName;
	readonly simLink = DemonologyWarlock.simLink;

	readonly isTankSpec = DemonologyWarlock.isTankSpec;
	readonly isHealingSpec = DemonologyWarlock.isHealingSpec;
	readonly isRangedDpsSpec = DemonologyWarlock.isRangedDpsSpec;
	readonly isMeleeDpsSpec = DemonologyWarlock.isMeleeDpsSpec;

	readonly canDualWield = DemonologyWarlock.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_metamorphosis.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return DemonologyWarlock.getIcon(size)
	}
}

export class DestructionWarlock extends Spec {
	static protoID = SpecProto.SpecDestructionWarlock;
	static class = Warlock;
	static friendlyName = 'Destruction';
	static simLink = getSpecSiteUrl('warlock', 'destruction');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = DestructionWarlock.protoID;
	readonly class = DestructionWarlock.class;
	readonly friendlyName = DestructionWarlock.friendlyName;
	readonly simLink = DestructionWarlock.simLink;

	readonly isTankSpec = DestructionWarlock.isTankSpec;
	readonly isHealingSpec = DestructionWarlock.isHealingSpec;
	readonly isRangedDpsSpec = DestructionWarlock.isRangedDpsSpec;
	readonly isMeleeDpsSpec = DestructionWarlock.isMeleeDpsSpec;

	readonly canDualWield = DestructionWarlock.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_shadow_rainoffire.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return DestructionWarlock.getIcon(size)
	}
}
