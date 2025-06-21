import { EligibleWeaponType, IconSize, PlayerClass } from '../player_class';
import { PlayerSpec } from '../player_spec';
import { AssassinationRogue, CombatRogue, SubtletyRogue } from '../player_specs/rogue';
import { ArmorType, Class, Race, RangedWeaponType, WeaponType } from '../proto/common';
import { RogueSpecs } from '../proto_utils/utils';

export class Rogue extends PlayerClass<Class.ClassRogue> {
	static classID = Class.ClassRogue as Class.ClassRogue;
	static friendlyName = 'Rogue';
	static hexColor = '#fff569';
	static specs: Record<string, PlayerSpec<RogueSpecs>> = {
		[AssassinationRogue.friendlyName]: AssassinationRogue,
		[CombatRogue.friendlyName]: CombatRogue,
		[SubtletyRogue.friendlyName]: SubtletyRogue,
	};
	static races: Race[] = [
		// [A]
		Race.RaceHuman,
		Race.RaceDwarf,
		Race.RaceNightElf,
		Race.RaceGnome,
		Race.RaceWorgen,
		Race.RaceAlliancePandaren,
		// [H]
		Race.RaceOrc,
		Race.RaceUndead,
		Race.RaceTroll,
		Race.RaceBloodElf,
		Race.RaceGoblin,
		Race.RaceHordePandaren
	];
	static armorTypes: ArmorType[] = [ArmorType.ArmorTypeLeather, ArmorType.ArmorTypeCloth];
	static weaponTypes: EligibleWeaponType[] = [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: false },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeSword },
	];
	static rangedWeaponTypes: RangedWeaponType[] = [];

	readonly classID = Rogue.classID;
	readonly friendlyName = Rogue.name;
	readonly hexColor = Rogue.hexColor;
	readonly specs = Rogue.specs;
	readonly races = Rogue.races;
	readonly armorTypes = Rogue.armorTypes;
	readonly weaponTypes = Rogue.weaponTypes;
	readonly rangedWeaponTypes = Rogue.rangedWeaponTypes;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/class_rogue.jpg`;
	};

	getIcon = (size: IconSize): string => {
		return Rogue.getIcon(size);
	};
}
