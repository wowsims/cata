import { EligibleWeaponType, IconSize } from './player_class.js';
import { ArmorType, Class, Race, RangedWeaponType, Spec } from './proto/common.js';
import { SpecClasses } from './proto_utils/utils';

export abstract class PlayerSpec<SpecType extends Spec> {
	static specID: Spec;
	static classID: Class;
	static friendlyName: string;
	static hexColor: string;
	static races: Race[] = [];
	static armorTypes: ArmorType[] = [];
	static weaponTypes: EligibleWeaponType[];
	static rangedWeaponTypes: RangedWeaponType[];

	abstract readonly specIndex: number;
	abstract readonly specID: SpecType;
	abstract readonly classID: SpecClasses<SpecType>;
	abstract readonly friendlyName: string;
	abstract readonly simLink: string;

	abstract readonly isTankSpec: boolean;
	abstract readonly isHealingSpec: boolean;
	abstract readonly isRangedDpsSpec: boolean;
	abstract readonly isMeleeDpsSpec: boolean;

	abstract readonly canDualWield: boolean;

	abstract getIcon(size: IconSize): string;
}
