import { Class, EligibleWeaponType, IconSize } from '../class';
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common';
import { Spec } from '../spec';
import { ArmsWarrior, FuryWarrior, ProtectionWarrior } from '../specs';

export class Warrior extends Class {
	static protoID = ClassProto.ClassWarrior;
	static friendlyName = 'Warrior';
	static hexColor = '#c79c6e';
	static specs: Record<string, Spec> = {
		[ArmsWarrior.name]: ArmsWarrior,
		[FuryWarrior.name]: FuryWarrior,
		[ProtectionWarrior.name]: ProtectionWarrior,
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
		Race.RaceTauren,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [
		ArmorType.ArmorTypePlate,
		ArmorType.ArmorTypeMail,
		ArmorType.ArmorTypeLeather,
		ArmorType.ArmorTypeCloth,
	];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
		RangedWeaponType.RangedWeaponTypeThrown,
	];

	readonly protoID = Warrior.protoID;
	readonly friendlyName = Warrior.name;
	readonly hexColor = Warrior.hexColor;
	readonly specs = Warrior.specs;
	readonly races = Warrior.races;
	readonly armorTypes = Warrior.armorTypes;
	readonly weaponTypes = Warrior.weaponTypes;
	readonly rangedWeaponTypes = Warrior.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_warrior.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Warrior.getIcon(size);
	}
}
