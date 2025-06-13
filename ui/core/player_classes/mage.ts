import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { ArcaneMage, FireMage, FrostMage } from '../player_specs/mage';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { MageSpecs } from '../proto_utils/utils';

export class Mage extends PlayerClass<Class.ClassMage> {
	static classID = Class.ClassMage as Class.ClassMage;
	static friendlyName = 'Mage';
	static hexColor = '#69ccf0';
	static specs: Record<string, PlayerSpec<MageSpecs>> = {
		[ArcaneMage.friendlyName]: ArcaneMage,
		[FireMage.friendlyName]: FireMage,
		[FrostMage.friendlyName]: FrostMage,
	};
	static races: Race[] = [
		// [H]
		Race.RaceTroll,
		Race.RaceGoblin,
		Race.RaceOrc,
		Race.RaceUndead,
		Race.RaceBloodElf,
		Race.RaceHordePandaren,
		// [A]
		Race.RaceWorgen,
		Race.RaceGnome,
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceNightElf,
		Race.RaceDraenei,
		Race.RaceAlliancePandaren,
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
