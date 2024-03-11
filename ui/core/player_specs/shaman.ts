import { IconSize } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import { getSpecSiteUrl } from '../proto_utils/utils';

export class ElementalShaman extends PlayerSpec<Spec.SpecElementalShaman> {
	static protoID = Spec.SpecElementalShaman as Spec.SpecElementalShaman;
	static playerClass = PlayerClasses.Shaman;
	static friendlyName = 'Elemental';
	static fullName = `${this.friendlyName} ${PlayerClasses.Shaman.friendlyName}`;
	static simLink = getSpecSiteUrl('shaman', 'elemental');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = true;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = ElementalShaman.protoID;
	readonly playerClass = ElementalShaman.playerClass;
	readonly friendlyName = ElementalShaman.friendlyName;
	readonly fullName = ElementalShaman.fullName;
	readonly simLink = ElementalShaman.simLink;

	readonly isTankSpec = ElementalShaman.isTankSpec;
	readonly isHealingSpec = ElementalShaman.isHealingSpec;
	readonly isRangedDpsSpec = ElementalShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = ElementalShaman.isMeleeDpsSpec;

	readonly canDualWield = ElementalShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightning.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return ElementalShaman.getIcon(size);
	};
}

export class EnhancementShaman extends PlayerSpec<Spec.SpecEnhancementShaman> {
	static protoID = Spec.SpecEnhancementShaman as Spec.SpecEnhancementShaman;
	static playerClass = PlayerClasses.Shaman;
	static friendlyName = 'Enhancement';
	static fullName = `${this.friendlyName} ${PlayerClasses.Shaman.friendlyName}`;
	static simLink = getSpecSiteUrl('shaman', 'enhancement');

	static isTankSpec = false;
	static isHealingSpec = false;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = true;

	static canDualWield = true;

	readonly protoID = EnhancementShaman.protoID;
	readonly playerClass = EnhancementShaman.playerClass;
	readonly friendlyName = EnhancementShaman.friendlyName;
	readonly fullName = EnhancementShaman.fullName;
	readonly simLink = EnhancementShaman.simLink;

	readonly isTankSpec = EnhancementShaman.isTankSpec;
	readonly isHealingSpec = EnhancementShaman.isHealingSpec;
	readonly isRangedDpsSpec = EnhancementShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = EnhancementShaman.isMeleeDpsSpec;

	readonly canDualWield = EnhancementShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightningshield.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return EnhancementShaman.getIcon(size);
	};
}

export class RestorationShaman extends PlayerSpec<Spec.SpecRestorationShaman> {
	static protoID = Spec.SpecRestorationShaman as Spec.SpecRestorationShaman;
	static playerClass = PlayerClasses.Shaman;
	static friendlyName = 'Restoration';
	static fullName = `${this.friendlyName} ${PlayerClasses.Shaman.friendlyName}`;
	static simLink = getSpecSiteUrl('shaman', 'restoration');

	static isTankSpec = false;
	static isHealingSpec = true;
	static isRangedDpsSpec = false;
	static isMeleeDpsSpec = false;

	static canDualWield = false;

	readonly protoID = RestorationShaman.protoID;
	readonly playerClass = RestorationShaman.playerClass;
	readonly friendlyName = RestorationShaman.friendlyName;
	readonly fullName = RestorationShaman.fullName;
	readonly simLink = RestorationShaman.simLink;

	readonly isTankSpec = RestorationShaman.isTankSpec;
	readonly isHealingSpec = RestorationShaman.isHealingSpec;
	readonly isRangedDpsSpec = RestorationShaman.isRangedDpsSpec;
	readonly isMeleeDpsSpec = RestorationShaman.isMeleeDpsSpec;

	readonly canDualWield = RestorationShaman.canDualWield;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_magicimmunity.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return RestorationShaman.getIcon(size);
	};
}
