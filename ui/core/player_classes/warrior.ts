import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { ArmsWarrior, FuryWarrior, ProtectionWarrior } from '../player_specs/warrior';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { WarriorSpecs } from '../proto_utils/utils';

export class Warrior extends PlayerClass<Class.ClassWarrior> {
	static classID = Class.ClassWarrior as Class.ClassWarrior;
	static friendlyName = 'Warrior';
	static hexColor = '#c79c6e';
	static specs: Record<string, PlayerSpec<WarriorSpecs>> = {
		[ArmsWarrior.friendlyName]: ArmsWarrior,
		[FuryWarrior.friendlyName]: FuryWarrior,
		[ProtectionWarrior.friendlyName]: ProtectionWarrior,
	};
	static races: Race[] = [
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceNightElf,
		Race.RaceGnome,
		Race.RaceDraenei,
		Race.RaceWorgen,
		Race.RaceAlliancePandaren,
		// [H]
		Race.RaceOrc,
		Race.RaceUndead,
		Race.RaceTauren,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
		Race.RaceHordePandaren,
	];
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypePlate, ArmorType.ArmorTypeMail, ArmorType.ArmorTypeLeather, ArmorType.ArmorTypeCloth];
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

	readonly classID = Warrior.classID;
	readonly friendlyName = Warrior.name;
	readonly hexColor = Warrior.hexColor;
	readonly specs = Warrior.specs;
	readonly races = Warrior.races;
	readonly armorTypes = Warrior.armorTypes;
	readonly weaponTypes = Warrior.weaponTypes;
	readonly rangedWeaponTypes = Warrior.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_warrior.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Warrior.getIcon(size);
	};
}
