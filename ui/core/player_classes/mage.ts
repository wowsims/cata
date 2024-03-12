import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { PlayerSpecs } from '../player_specs';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { MageSpecs } from '../proto_utils/utils';

export class Mage extends PlayerClass<Class.ClassMage> {
	static classID = Class.ClassMage as Class.ClassMage;
	static friendlyName = 'Mage';
	static hexColor = '#69ccf0';
	static specs: Record<string, PlayerSpec<MageSpecs>> = {
		[PlayerSpecs.ArcaneMage.name]: PlayerSpecs.ArcaneMage,
		[PlayerSpecs.FireMage.name]: PlayerSpecs.FireMage,
		[PlayerSpecs.FrostMage.name]: PlayerSpecs.FrostMage,
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
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypeCloth];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [RangedWeaponType.RangedWeaponTypeWand];

	readonly classID = Mage.classID;
	readonly friendlyName = Mage.name;
	readonly hexColor = Mage.hexColor;
	readonly specs = Mage.specs;
	readonly races = Mage.races;
	readonly armorTypes = Mage.armorTypes;
	readonly weaponTypes = Mage.weaponTypes;
	readonly rangedWeaponTypes = Mage.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_mage.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Mage.getIcon(size);
	};
}
