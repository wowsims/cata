import { PlayerClasses } from './player_classes';
import { PlayerSpec } from './player_spec';
import { ArmorType, Class, Race, RangedWeaponType, Spec, WeaponType } from './proto/common.js';
export type IconSize = 'small' | 'medium' | 'large';

export interface EligibleWeaponType {
	weaponType: WeaponType;
	canUseTwoHand?: boolean;
}

export abstract class PlayerClass<ClassType extends Class> {
	abstract readonly protoID: ClassType;
	abstract readonly friendlyName: string;
	abstract readonly hexColor: string;
	abstract readonly specs: Record<string, PlayerSpec<Spec>>;
	abstract readonly races: Race[];
	abstract readonly armorTypes: ArmorType[];
	abstract readonly weaponTypes: EligibleWeaponType[];
	abstract readonly rangedWeaponTypes: RangedWeaponType[];

	abstract getIcon(size: IconSize): string;
}

export const naturalPlayerClassOrder: Array<PlayerClass<Class>> = [
	PlayerClasses.DeathKnight,
	PlayerClasses.Druid,
	PlayerClasses.Hunter,
	PlayerClasses.Mage,
	PlayerClasses.Paladin,
	PlayerClasses.Priest,
	PlayerClasses.Rogue,
	PlayerClasses.Shaman,
	PlayerClasses.Warlock,
	PlayerClasses.Warrior,
];
