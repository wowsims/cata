import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { PlayerSpecs } from '../player_specs';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { ShamanSpecs } from '../proto_utils/utils';

export class Shaman extends PlayerClass<Class.ClassShaman> {
	static classID = Class.ClassShaman as Class.ClassShaman;
	static friendlyName = 'Shaman';
	static hexColor = '#2459ff';
	static specs: Record<string, PlayerSpec<ShamanSpecs>> = {
		[PlayerSpecs.ElementalShaman.name]: PlayerSpecs.ElementalShaman,
		[PlayerSpecs.EnhancementShaman.name]: PlayerSpecs.EnhancementShaman,
		[PlayerSpecs.RestorationShaman.name]: PlayerSpecs.RestorationShaman,
	};
	static races: Race[] = [
		// [A]
		Race.RaceDwarf,
		Race.RaceDraenei,
		// [H]
		Race.RaceOrc,
		Race.RaceTauren,
		Race.RaceTroll,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypeMail, ArmorType.ArmorTypeLeather, ArmorType.ArmorTypeCloth];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [RangedWeaponType.RangedWeaponTypeTotem];

	readonly classID = Shaman.classID;
	readonly friendlyName = Shaman.name;
	readonly hexColor = Shaman.hexColor;
	readonly specs = Shaman.specs;
	readonly races = Shaman.races;
	readonly armorTypes = Shaman.armorTypes;
	readonly weaponTypes = Shaman.weaponTypes;
	readonly rangedWeaponTypes = Shaman.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_shaman.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Shaman.getIcon(size);
	};
}
