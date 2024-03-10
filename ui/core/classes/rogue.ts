import { Class, EligibleWeaponType, IconSize } from '../class';
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common';
import { Spec } from '../spec';
import { AssassinationRogue, CombatRogue, SubtletyRogue } from '../specs';

export class Rogue extends Class {
	static protoID = ClassProto.ClassRogue;
	static friendlyName = 'Rogue';
	static hexColor = '#fff569';
	static specs: Record<string, Spec> = {
		[AssassinationRogue.name]: AssassinationRogue,
		[CombatRogue.name]: CombatRogue,
		[SubtletyRogue.name]: SubtletyRogue,
	};
	static races: Race[] = [
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceNightElf,
		Race.RaceGnome,
		Race.RaceWorgen,
		// [H]
		Race.RaceOrc,
		Race.RaceUndead,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [
		ArmorType.ArmorTypeLeather,
		ArmorType.ArmorTypeCloth,
	];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: false },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeSword },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
		RangedWeaponType.RangedWeaponTypeThrown,
	];

	readonly protoID = Rogue.protoID;
	readonly friendlyName = Rogue.name;
	readonly hexColor = Rogue.hexColor;
	readonly specs = Rogue.specs;
	readonly races = Rogue.races;
	readonly armorTypes = Rogue.armorTypes;
	readonly weaponTypes = Rogue.weaponTypes;
	readonly rangedWeaponTypes = Rogue.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_rogue.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Rogue.getIcon(size);
	}
}
