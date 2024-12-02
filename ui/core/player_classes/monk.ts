import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { BrewmasterMonk, MistweaverMonk, WindwalkerMonk } from '../player_specs/monk';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { MonkSpecs } from '../proto_utils/utils';

export class Monk extends PlayerClass<Class.ClassMonk> {
	static classID = Class.ClassMonk as Class.ClassMonk;
	static friendlyName = 'Monk';
	static hexColor = '#00ff98';
	static specs: Record<string, PlayerSpec<MonkSpecs>> = {
		[BrewmasterMonk.friendlyName]: BrewmasterMonk,
		[MistweaverMonk.friendlyName]: MistweaverMonk,
		[WindwalkerMonk.friendlyName]: WindwalkerMonk,
	};
	static races: Race[] = [
		// [A]
		Race.RaceAlliancePandaren,
		Race.RaceDraenei,
		Race.RaceDwarf,
		Race.RaceGnome,
		Race.RaceHuman,
		Race.RaceNightElf,
		Race.RaceWorgen,
		// [H]
		Race.RaceHordePandaren,
		Race.RaceBloodElf,
		Race.RaceOrc,
		Race.RaceTauren,
		Race.RaceTroll,
		Race.RaceUndead,
		Race.RaceGoblin,
	];
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypeLeather, ArmorType.ArmorTypeCloth];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: false },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: false },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: false },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [];

	readonly classID = Monk.classID;
	readonly friendlyName = Monk.name;
	readonly hexColor = Monk.hexColor;
	readonly specs = Monk.specs;
	readonly races = Monk.races;
	readonly armorTypes = Monk.armorTypes;
	readonly weaponTypes = Monk.weaponTypes;
	readonly rangedWeaponTypes = Monk.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_monk.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Monk.getIcon(size);
	};
}