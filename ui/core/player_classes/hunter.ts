import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { BeastMasteryHunter, MarksmanshipHunter, SurvivalHunter } from '../player_specs/hunter';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { HunterSpecs } from '../proto_utils/utils';

export class Hunter extends PlayerClass<Class.ClassHunter> {
	static classID = Class.ClassHunter as Class.ClassHunter;
	static friendlyName = 'Hunter';
	static hexColor = '#abd473';
	static specs: Record<string, PlayerSpec<HunterSpecs>> = {
		[BeastMasteryHunter.friendlyName]: BeastMasteryHunter,
		[MarksmanshipHunter.friendlyName]: MarksmanshipHunter,
		[SurvivalHunter.friendlyName]: SurvivalHunter,
	};
	static races: Race[] = [
		// [A]

		Race.RaceWorgen,
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceNightElf,
		Race.RaceDraenei,
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
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypeMail];
	static weaponTypes: EligibleWeaponType[] = []; // hunter cannot wear weapons anymore
	static rangedWeaponTypes: RangedWeaponType[] = [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
	];

	readonly classID = Hunter.classID;
	readonly friendlyName = Hunter.name;
	readonly hexColor = Hunter.hexColor;
	readonly specs = Hunter.specs;
	readonly races = Hunter.races;
	readonly armorTypes = Hunter.armorTypes;

	readonly weaponTypes = Hunter.weaponTypes;
	readonly rangedWeaponTypes = Hunter.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_hunter.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Hunter.getIcon(size);
	};
}
