import { Class, EligibleWeaponType, IconSize } from '../class';
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common';
import { Spec } from '../spec';
import { HolyPaladin, ProtectionPaladin, RetributionPaladin } from '../specs';

export class Paladin extends Class {
	static protoID = ClassProto.ClassPaladin;
	static friendlyName = 'Paladin';
	static hexColor = '#f58cba';
	static specs: Record<string, Spec> = {
		[HolyPaladin.name]: HolyPaladin,
		[ProtectionPaladin.name]: ProtectionPaladin,
		[RetributionPaladin.name]: RetributionPaladin,
	};
	static races: Race[] = [
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceDraenei,
		// [H]
		Race.RaceTauren,
		Race.RaceBloodElf,
	];
	static armorTypes: ArmorType[] = [
		ArmorType.ArmorTypePlate,
		ArmorType.ArmorTypeMail,
		ArmorType.ArmorTypeLeather,
		ArmorType.ArmorTypeCloth,
	];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeLibram,
	];

	readonly protoID = Paladin.protoID;
	readonly friendlyName = Paladin.name;
	readonly hexColor = Paladin.hexColor;
	readonly specs = Paladin.specs;
	readonly races = Paladin.races;
	readonly armorTypes = Paladin.armorTypes;
	readonly weaponTypes = Paladin.weaponTypes;
	readonly rangedWeaponTypes = Paladin.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_paladin.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Paladin.getIcon(size);
	}
}
