import { Class, EligibleWeaponType, IconSize } from '../class';
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common';
import { Spec } from '../spec';
import { ArcaneMage, FireMage, FrostMage } from '../specs';

export class Mage extends Class {
	static protoID = ClassProto.ClassMage;
	static friendlyName = 'Mage';
	static hexColor = '#69ccf0';
	static specs: Record<string, Spec> = {
		[ArcaneMage.name]: ArcaneMage,
		[FireMage.name]: FireMage,
		[FrostMage.name]: FrostMage,
	};
	static races: Race[] = [
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceNightElf,
		Race.RaceGnome,
		Race.RaceDraenei,
		Race.RaceWorgen,
		// [H]
		Race.RaceOrc,
		Race.RaceUndead,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [
		ArmorType.ArmorTypeCloth,
	];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeWand,
	];

	readonly protoID = Mage.protoID;
	readonly friendlyName = Mage.name;
	readonly hexColor = Mage.hexColor;
	readonly specs = Mage.specs;
	readonly races = Mage.races;
	readonly armorTypes = Mage.armorTypes;
	readonly weaponTypes = Mage.weaponTypes;
	readonly rangedWeaponTypes = Mage.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_mage.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Mage.getIcon(size);
	}
}
