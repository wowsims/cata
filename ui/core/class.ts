import { ArmorType, Class as ClassProto, Race, WeaponType } from './proto/common.js';
import { Spec } from './spec.js';

export type IconSize = 'small' | 'medium' | 'large'

export interface EligibleWeaponType {
	weaponType: WeaponType,
	canUseTwoHand?: boolean,
}

export abstract class Class {
	abstract readonly protoID: ClassProto;
	abstract readonly friendlyName: string;
	abstract readonly hexColor: string;
	abstract readonly specs: Record<string, Spec>;
	abstract readonly races: Race[];
	abstract readonly armorTypes: ArmorType[];
	abstract readonly weaponTypes: EligibleWeaponType[];

	abstract getIcon(size: IconSize): string;
}
