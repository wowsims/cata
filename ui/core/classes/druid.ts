import { Class, EligibleWeaponType, IconSize } from '../class';
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common';
import { Spec } from '../spec';
import { BalanceDruid, FeralDruid, RestorationDruid } from '../specs';

export class Druid extends Class {
	static protoID = ClassProto.ClassDruid;
	static friendlyName = 'Druid';
	static hexColor = '#ff7d0a';
	static specs: Record<string, Spec> = {
		[BalanceDruid.name]: BalanceDruid,
		[FeralDruid.name]: FeralDruid,
		[RestorationDruid.name]: RestorationDruid,
	};
	static races: Race[] = [
		// [A]
		Race.RaceNightElf,
		Race.RaceWorgen,
		// [H]
		Race.RaceTauren,
		Race.RaceTroll,
	];
	static armorTypes: ArmorType[] = [
		ArmorType.ArmorTypeLeather,
		ArmorType.ArmorTypeCloth,
	];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeLibram,
	];

	readonly protoID = Druid.protoID;
	readonly friendlyName = Druid.name;
	readonly hexColor = Druid.hexColor;
	readonly specs = Druid.specs;
	readonly races = Druid.races;
	readonly armorTypes = Druid.armorTypes;
	readonly weaponTypes = Druid.weaponTypes;
	readonly rangedWeaponTypes = Druid.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_druid.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Druid.getIcon(size);
	}
}
