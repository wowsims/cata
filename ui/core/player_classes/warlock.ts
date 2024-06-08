import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { AfflictionWarlock, DemonologyWarlock, DestructionWarlock } from '../player_specs/warlock';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { WarlockSpecs } from '../proto_utils/utils';

export class Warlock extends PlayerClass<Class.ClassWarlock> {
	static classID = Class.ClassWarlock as Class.ClassWarlock;
	static friendlyName = 'Warlock';
	static hexColor = '#9482c9';
	static specs: Record<string, PlayerSpec<WarlockSpecs>> = {
		[AfflictionWarlock.friendlyName]: AfflictionWarlock,
		[DemonologyWarlock.friendlyName]: DemonologyWarlock,
		[DestructionWarlock.friendlyName]: DestructionWarlock,
	};
	static races: Race[] = [
		// [H]
		Race.RaceOrc,
		Race.RaceUndead,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceGnome,
		Race.RaceWorgen,
	];
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypeCloth];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [RangedWeaponType.RangedWeaponTypeWand];

	readonly classID = Warlock.classID;
	readonly friendlyName = Warlock.name;
	readonly hexColor = Warlock.hexColor;
	readonly specs = Warlock.specs;
	readonly races = Warlock.races;
	readonly armorTypes = Warlock.armorTypes;
	readonly weaponTypes = Warlock.weaponTypes;
	readonly rangedWeaponTypes = Warlock.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_warlock.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Warlock.getIcon(size);
	};
}
