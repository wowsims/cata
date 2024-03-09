import { Class, EligibleWeaponType, IconSize } from "../class";
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common.js';
import { Spec } from "../spec";
import { ElementalShaman, EnhancementShaman, RestorationShaman } from "../specs/shaman";

export class Shaman extends Class {
	static protoID = ClassProto.ClassShaman;
	static friendlyName = 'Shaman';
	static hexColor = '#2459ff';
	static specs: Record<string, Spec> = {
		[ElementalShaman.name]: ElementalShaman,
		[EnhancementShaman.name]: EnhancementShaman,
		[RestorationShaman.name]: RestorationShaman,
	};
	static races: Race[] = [
		Race.RaceDraenei,
		Race.RaceDwarf,
		Race.RaceGoblin,
		Race.RaceOrc,
		Race.RaceTauren,
		Race.RaceTroll,
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
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeTotem,
	];

	readonly protoID = Shaman.protoID;
	readonly friendlyName = Shaman.name;
	readonly hexColor = Shaman.hexColor;
	readonly specs = Shaman.specs;
	readonly races = Shaman.races;
	readonly armorTypes = Shaman.armorTypes;
	readonly weaponTypes = Shaman.weaponTypes;
	readonly rangedWeaponTypes = Shaman.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_shaman.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Shaman.getIcon(size);
	}
}
