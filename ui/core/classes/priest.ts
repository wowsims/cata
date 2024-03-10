import { Class, EligibleWeaponType, IconSize } from '../class';
import { ArmorType, Class as ClassProto, Race, RangedWeaponType, WeaponType} from '../proto/common';
import { Spec } from '../spec';
import { DisciplinePriest, HolyPriest, ShadowPriest } from '../specs';

export class Priest extends Class {
	static protoID = ClassProto.ClassPriest;
	static friendlyName = 'Priest';
	static hexColor = '#fff';
	static specs: Record<string, Spec> = {
		[DisciplinePriest.name]: DisciplinePriest,
		[HolyPriest.name]: HolyPriest,
		[ShadowPriest.name]: ShadowPriest,
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
		Race.RaceUndead,
		Race.RaceTauren,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [
		ArmorType.ArmorTypeCloth,
	];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeMace },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeWand,
	];

	readonly protoID = Priest.protoID;
	readonly friendlyName = Priest.name;
	readonly hexColor = Priest.hexColor;
	readonly specs = Priest.specs;
	readonly races = Priest.races;
	readonly armorTypes = Priest.armorTypes;
	readonly weaponTypes = Priest.weaponTypes;
	readonly rangedWeaponTypes = Priest.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_priest.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Priest.getIcon(size);
	}
}
