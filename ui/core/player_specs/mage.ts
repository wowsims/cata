import { IconSize } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class ArcaneMage extends PlayerSpec<Spec.SpecArcaneMage> {
	static specIndex = 0;
	static specID = Spec.SpecArcaneMage as Spec.SpecArcaneMage;
	static classID = Class.ClassMage as Class.ClassMage;
	static friendlyName = 'Arcane';
	static simLink = getSpecSiteUrl('mage', 'arcane');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = ArcaneMage.specIndex;
	readonly specID = ArcaneMage.specID;
	readonly classID = ArcaneMage.classID;
	readonly friendlyName = ArcaneMage.friendlyName;
	readonly simLink = ArcaneMage.simLink;

	readonly isTankSpec = ArcaneMage.isTankSpec;
	readonly isHealingSpec = ArcaneMage.isHealingSpec;
	readonly isRangedDpsSpec = ArcaneMage.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ArcaneMage.isMeleeDpsSpec;

	readonly canDualWield = ArcaneMage.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_holy_magicalsentry.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return ArcaneMage.getIcon(size);
	};
}

export class FireMage extends PlayerSpec<Spec.SpecFireMage> {
	static specIndex = 1;
	static specID = Spec.SpecFireMage as Spec.SpecFireMage;
	static classID = Class.ClassMage as Class.ClassMage;
	static friendlyName = 'Fire';
	static simLink = getSpecSiteUrl('mage', 'fire');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = FireMage.specIndex;
	readonly specID = FireMage.specID;
	readonly classID = FireMage.classID;
	readonly friendlyName = FireMage.friendlyName;
	readonly simLink = FireMage.simLink;

	readonly isTankSpec = FireMage.isTankSpec;
	readonly isHealingSpec = FireMage.isHealingSpec;
	readonly isRangedDpsSpec = FireMage.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FireMage.isMeleeDpsSpec;

	readonly canDualWield = FireMage.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_fire_firebolt02.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return FireMage.getIcon(size);
	};
}

export class FrostMage extends PlayerSpec<Spec.SpecFrostMage> {
	static specIndex = 2;
	static specID = Spec.SpecFrostMage as Spec.SpecFrostMage;
	static classID = Class.ClassMage as Class.ClassMage;
	static friendlyName = 'Frost';
	static simLink = getSpecSiteUrl('mage', 'frost');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly specIndex = FrostMage.specIndex;
	readonly specID = FrostMage.specID;
	readonly classID = FrostMage.classID;
	readonly friendlyName = FrostMage.friendlyName;
	readonly simLink = FrostMage.simLink;

	readonly isTankSpec = FrostMage.isTankSpec;
	readonly isHealingSpec = FrostMage.isHealingSpec;
	readonly isRangedDpsSpec = FrostMage.isRangedDpsSpec;
	readonly isMeleeDpsSpec = FrostMage.isMeleeDpsSpec;

	readonly canDualWield = FrostMage.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_frost_frostbolt02.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return FrostMage.getIcon(size);
	};
}
