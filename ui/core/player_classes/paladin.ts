import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { HolyPaladin, ProtectionPaladin, RetributionPaladin } from '../player_specs/paladin';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { PaladinSpecs } from '../proto_utils/utils';

export class Paladin extends PlayerClass<Class.ClassPaladin> {
	static classID = Class.ClassPaladin as Class.ClassPaladin;
	static friendlyName = 'Paladin';
	static cssClass = 'paladin';
	static hexColor = '#f58cba';
	static specs: Record<string, PlayerSpec<PaladinSpecs>> = {
		[HolyPaladin.friendlyName]: HolyPaladin,
		[ProtectionPaladin.friendlyName]: ProtectionPaladin,
		[RetributionPaladin.friendlyName]: RetributionPaladin,
	};
	static races: Race[] = [
		// [H]
		Race.RaceBloodElf,
		Race.RaceTauren,
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceDraenei,
	];
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypePlate, ArmorType.ArmorTypeMail, ArmorType.ArmorTypeLeather, ArmorType.ArmorTypeCloth];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [];

	readonly classID = Paladin.classID;
	readonly friendlyName = Paladin.name;
	readonly cssClass = Paladin.cssClass;
	readonly hexColor = Paladin.hexColor;
	readonly specs = Paladin.specs;
	readonly races = Paladin.races;
	readonly armorTypes = Paladin.armorTypes;
	readonly weaponTypes = Paladin.weaponTypes;
	readonly rangedWeaponTypes = Paladin.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_paladin.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Paladin.getIcon(size);
	};
}
