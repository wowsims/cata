import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { ElementalShaman, EnhancementShaman, RestorationShaman } from '../player_specs/shaman';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { ShamanSpecs } from '../proto_utils/utils';

export class Shaman extends PlayerClass<Class.ClassShaman> {
	static classID = Class.ClassShaman as Class.ClassShaman;
	static friendlyName = 'Shaman';
	static hexColor = '#2459ff';
	static specs: Record<string, PlayerSpec<ShamanSpecs>> = {
		[ElementalShaman.friendlyName]: ElementalShaman,
		[EnhancementShaman.friendlyName]: EnhancementShaman,
		[RestorationShaman.friendlyName]: RestorationShaman,
	};
	static races: Race[] = [
		// [H]
		Race.RaceTroll,
		Race.RaceOrc,
		Race.RaceTauren,
		Race.RaceGoblin,
		Race.RaceHordePandaren,
		// [A]
		Race.RaceDwarf,
		Race.RaceDraenei,
		Race.RaceAlliancePandaren,

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
	static rangedWeaponTypes: RangedWeaponType[] = [];

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
