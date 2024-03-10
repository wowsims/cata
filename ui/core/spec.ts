import { Class, IconSize } from './class.js';
import { Spec as SpecProto } from './proto/common.js';
import { ElementalShaman, EnhancementShaman, RestorationShaman } from './specs/shaman.js';

export abstract class Spec {
	abstract readonly protoID: SpecProto;
	abstract readonly class: Class;
	abstract readonly friendlyName: string;
	abstract readonly simLink: string;

	abstract readonly isTankSpec: boolean;
	abstract readonly isHealingSpec: boolean;
	abstract readonly isRangedDpsSpec: boolean;
	abstract readonly isMeleeDpsSpec: boolean;

	abstract readonly canDualWield: boolean;

	abstract getIcon(size: IconSize): string;
}

// TODO: Cata - Update list / maybe use Spec objects
export const naturalSpecOrder: Array<Spec> = [
	ElementalShaman,
	EnhancementShaman,
	RestorationShaman,
];

const protoToSpec: Record<SpecProto, Spec | undefined> = {
	[SpecProto.SpecUnknown]: undefined,

	[SpecProto.SpecBloodDeathKnight]: undefined,
	[SpecProto.SpecFrostDeathKnight]: undefined,
	[SpecProto.SpecUnholyDeathKnight]: undefined,

	[SpecProto.SpecBalanceDruid]: undefined,
	[SpecProto.SpecFeralDruid]: undefined,
	[SpecProto.SpecRestorationDruid]: undefined,

	[SpecProto.SpecBeastMasteryHunter]: undefined,
	[SpecProto.SpecMarksmanshipHunter]: undefined,
	[SpecProto.SpecSurvivalHunter]: undefined,

	[SpecProto.SpecArcaneMage]: undefined,
	[SpecProto.SpecFireMage]: undefined,
	[SpecProto.SpecFrostMage]: undefined,

	[SpecProto.SpecHolyPaladin]: undefined,
	[SpecProto.SpecProtectionPaladin]: undefined,
	[SpecProto.SpecRetributionPaladin]: undefined,

	[SpecProto.SpecDisciplinePriest]: undefined,
	[SpecProto.SpecHolyPriest]: undefined,
	[SpecProto.SpecShadowPriest]: undefined,

	[SpecProto.SpecAssassinationRogue]: undefined,
	[SpecProto.SpecCombatRogue]: undefined,
	[SpecProto.SpecSubtletyRogue]: undefined,

	[SpecProto.SpecElementalShaman]: ElementalShaman,
	[SpecProto.SpecEnhancementShaman]: EnhancementShaman,
	[SpecProto.SpecRestorationShaman]: RestorationShaman,

	[SpecProto.SpecAfflictionWarlock]: undefined,
	[SpecProto.SpecDemonologyWarlock]: undefined,
	[SpecProto.SpecDestructionWarlock]: undefined,

	[SpecProto.SpecArmsWarrior]: undefined,
	[SpecProto.SpecFuryWarrior]: undefined,
	[SpecProto.SpecProtectionWarrior]: undefined,
};

export const specFromProto = (protoId: SpecProto): Spec => {
	if (protoId == SpecProto.SpecUnknown) {
		throw new Error("Invalid Spec");
	}

	return protoToSpec[protoId] as Spec;
};

const LOCAL_STORAGE_PREFIX = '__cata';
export const getLocalStorageKey = (spec: Spec): string => {
	return `${LOCAL_STORAGE_PREFIX}_${spec.friendlyName.toLowerCase().replace(/\s/, '_')}_${spec.class.friendlyName.toLowerCase().replace(/\s/, '_')}`
};
