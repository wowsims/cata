import { Class, EligibleWeaponType, IconSize } from '../class';
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common';
import { Spec } from '../spec';
import { BeastMasteryHunter, MarksmanshipHunter, SurvivalHunter } from '../specs';

export class Hunter extends Class {
	static protoID = ClassProto.ClassHunter;
	static friendlyName = 'Hunter';
	static hexColor = '#abd473';
	static specs: Record<string, Spec> = {
		[BeastMasteryHunter.name]: BeastMasteryHunter,
		[MarksmanshipHunter.name]: MarksmanshipHunter,
		[SurvivalHunter.name]: SurvivalHunter,
	};
	static races: Race[] = [
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceNightElf,
		Race.RaceDraenei,
		// [H]
		Race.RaceOrc,
		Race.RaceUndead,
		Race.RaceTauren,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [
		ArmorType.ArmorTypeMail,
		ArmorType.ArmorTypeLeather,
		ArmorType.ArmorTypeCloth,
	];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
	];

	readonly protoID = Hunter.protoID;
	readonly friendlyName = Hunter.name;
	readonly hexColor = Hunter.hexColor;
	readonly specs = Hunter.specs;
	readonly races = Hunter.races;
	readonly armorTypes = Hunter.armorTypes;
	readonly weaponTypes = Hunter.weaponTypes;
	readonly rangedWeaponTypes = Hunter.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_hunter.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Hunter.getIcon(size);
	}
}
