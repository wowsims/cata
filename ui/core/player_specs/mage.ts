import { IconSize } from '../player_class';
import { Mage } from '../player_classes/mage';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class ArcaneMage extends PlayerSpec<Spec.SpecArcaneMage> {
	static protoID = Spec.SpecArcaneMage as Spec.SpecArcaneMage;
	static playerClass = Mage;
	static friendlyName = 'Arcane';
	static fullName = `${this.friendlyName} ${Mage.friendlyName}`;
	static simLink = getSpecSiteUrl('mage', 'arcane');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ArcaneMage.protoID;
	readonly playerClass = ArcaneMage.playerClass;
	readonly friendlyName = ArcaneMage.friendlyName;
	readonly fullName = ArcaneMage.fullName;
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
	static protoID = Spec.SpecFireMage as Spec.SpecFireMage;
	static playerClass = Mage;
	static friendlyName = 'Fire';
	static fullName = `${this.friendlyName} ${Mage.friendlyName}`;
	static simLink = getSpecSiteUrl('mage', 'fire');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = FireMage.protoID;
	readonly playerClass = FireMage.playerClass;
	readonly friendlyName = FireMage.friendlyName;
	readonly fullName = FireMage.fullName;
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
	static protoID = Spec.SpecFrostMage as Spec.SpecFrostMage;
	static playerClass = Mage;
	static friendlyName = 'Frost';
	static fullName = `${this.friendlyName} ${Mage.friendlyName}`;
	static simLink = getSpecSiteUrl('mage', 'frost');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = FrostMage.protoID;
	readonly playerClass = FrostMage.playerClass;
	readonly friendlyName = FrostMage.friendlyName;
	readonly fullName = FrostMage.fullName;
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
