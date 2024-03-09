import { Class } from '../class.js';
import { Shaman } from '../classes/shaman.js';
import { REPO_NAME } from '../constants/other.js'
import { Player , ResourceType } from '../proto/api.js';
import {
	ArmorType,
	Class as ClassProto,
	EnchantType,
	Faction,
	HandType,
	ItemSlot,
	ItemType,
	Race,
	RangedWeaponType,
	Spec as SpecProto,
	UnitReference,
	UnitReference_Type,
	WeaponType
} from '../proto/common.js';
import { Deathknight, Deathknight_Options as DeathknightOptions,Deathknight_Rotation as DeathknightRotation, DeathknightTalents, TankDeathknight, TankDeathknight_Options as TankDeathknightOptions,TankDeathknight_Rotation as TankDeathknightRotation } from '../proto/deathknight.js';
import {
	BalanceDruid,
	BalanceDruid_Options as BalanceDruidOptions,
	BalanceDruid_Rotation as BalanceDruidRotation,
	DruidTalents,
	FeralDruid,
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralTankDruid,
	FeralTankDruid_Options as FeralTankDruidOptions,
	FeralTankDruid_Rotation as FeralTankDruidRotation,
	RestorationDruid,
	RestorationDruid_Options as RestorationDruidOptions,
	RestorationDruid_Rotation as RestorationDruidRotation,
} from '../proto/druid.js';
import { Hunter, Hunter_Options as HunterOptions,Hunter_Rotation as HunterRotation, HunterTalents } from '../proto/hunter.js';
import { Mage, Mage_Options as MageOptions,Mage_Rotation as MageRotation, MageTalents } from '../proto/mage.js';
import { Blessings ,
	HolyPaladin,
	HolyPaladin_Options as HolyPaladinOptions,
	HolyPaladin_Rotation as HolyPaladinRotation,
	PaladinTalents,
	ProtectionPaladin,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
	ProtectionPaladin_Rotation as ProtectionPaladinRotation,
	RetributionPaladin,
	RetributionPaladin_Options as RetributionPaladinOptions,
	RetributionPaladin_Rotation as RetributionPaladinRotation,
} from '../proto/paladin.js';
import {
	HealingPriest,
	HealingPriest_Options as HealingPriestOptions,
	HealingPriest_Rotation as HealingPriestRotation,
	PriestTalents,
	ShadowPriest,
	ShadowPriest_Options as ShadowPriestOptions,
	ShadowPriest_Rotation as ShadowPriestRotation,
	SmitePriest,
	SmitePriest_Options as SmitePriestOptions,
	SmitePriest_Rotation as SmitePriestRotation,
} from '../proto/priest.js';
import { Rogue, Rogue_Options as RogueOptions,Rogue_Rotation as RogueRotation, RogueTalents } from '../proto/rogue.js';
import {
	ElementalShaman,
	ElementalShaman_Options as ElementalShamanOptions,
	ElementalShaman_Rotation as ElementalShamanRotation,
	EnhancementShaman,
	EnhancementShaman_Options as EnhancementShamanOptions,
	EnhancementShaman_Rotation as EnhancementShamanRotation,
	RestorationShaman,
	RestorationShaman_Options as RestorationShamanOptions,
	RestorationShaman_Rotation as RestorationShamanRotation,
	ShamanTalents,
} from '../proto/shaman.js';
import {
	BlessingsAssignment,
	BlessingsAssignments,
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';
import { Warlock, Warlock_Options as WarlockOptions,Warlock_Rotation as WarlockRotation, WarlockTalents } from '../proto/warlock.js';
import { ProtectionWarrior, ProtectionWarrior_Options as ProtectionWarriorOptions,ProtectionWarrior_Rotation as ProtectionWarriorRotation,Warrior, Warrior_Options as WarriorOptions,Warrior_Rotation as WarriorRotation, WarriorTalents  } from '../proto/warrior.js';
import * as Gems from '../proto_utils/gems.js';
import { Spec } from '../spec.js';
import { getEnumValues , intersection , maxIndex , sum } from '../utils.js';
import { Stats } from './stats.js';

// TODO: Cata - Re-evaluate whether or not these are needed
export type DeathknightSpecs = SpecProto.SpecBloodDeathKnight | SpecProto.SpecFrostDeathKnight | SpecProto.SpecUnholyDeathKnight
export type DruidSpecs = SpecProto.SpecBalanceDruid | SpecProto.SpecFeralDruid | SpecProto.SpecRestorationDruid;
export type HunterSpecs = SpecProto.SpecBeastMasteryHunter | SpecProto.SpecMarksmanshipHunter | SpecProto.SpecSurvivalHunter;
export type MageSpecs = SpecProto.SpecArcaneMage | SpecProto.SpecFireMage | SpecProto.SpecFrostMage;
export type PaladinSpecs = SpecProto.SpecHolyPaladin | SpecProto.SpecRetributionPaladin | SpecProto.SpecProtectionPaladin;
export type PriestSpecs = SpecProto.SpecDisciplinePriest | SpecProto.SpecHolyPriest | SpecProto.SpecShadowPriest;
export type RogueSpecs = SpecProto.SpecAssassinationRogue | SpecProto.SpecCombatRogue | SpecProto.SpecSubtletyRogue;
export type ShamanSpecs = SpecProto.SpecElementalShaman | SpecProto.SpecEnhancementShaman | SpecProto.SpecRestorationShaman;
export type WarlockSpecs = SpecProto.SpecAfflictionWarlock | SpecProto.SpecDemonologyWarlock | SpecProto.SpecDestructionWarlock;
export type WarriorSpecs = SpecProto.SpecArmsWarrior | SpecProto.SpecFuryWarrior | SpecProto.SpecProtectionWarrior;

export type ClassSpecs<T extends ClassProto> =
	T extends ClassProto.ClassDeathknight ? DeathknightSpecs :
	T extends ClassProto.ClassDruid ? DruidSpecs :
	T extends ClassProto.ClassHunter ? HunterSpecs :
	T extends ClassProto.ClassMage ? MageSpecs :
	T extends ClassProto.ClassPaladin ? PaladinSpecs :
	T extends ClassProto.ClassPriest ? PriestSpecs :
	T extends ClassProto.ClassRogue ? RogueSpecs :
	T extends ClassProto.ClassShaman ? ShamanSpecs :
	T extends ClassProto.ClassWarlock ? WarlockSpecs :
	T extends ClassProto.ClassWarrior ? WarriorSpecs :
	ShamanSpecs; // Should never reach this case

export const NUM_SPECS = getEnumValues(Spec).length;

// The order in which specs should be presented, when it matters.
// TODO: Cata - Update list / maybe use Spec objects
export const naturalSpecOrder: Array<SpecProto> = [
	SpecProto.SpecBalanceDruid,
	SpecProto.SpecFeralDruid,
	SpecProto.SpecRestorationDruid,
	SpecProto.SpecHolyPaladin,
	SpecProto.SpecProtectionPaladin,
	SpecProto.SpecRetributionPaladin,
	SpecProto.SpecShadowPriest,
	SpecProto.SpecElementalShaman,
	SpecProto.SpecEnhancementShaman,
	SpecProto.SpecRestorationShaman,
	SpecProto.SpecProtectionWarrior,
];

export const naturalClassOrder: Array<Class> = [
	Shaman,
];

export const raidSimIcon = '/cata/assets/img/raid_icon.png';
export const raidSimLabel = 'Full Raid Sim';

// Converts '1231321-12313123-0' to [40, 21, 0].
export function getTalentTreePoints(talentsString: string): Array<number> {
	const trees = talentsString.split('-');
	if (trees.length == 2)  {
		trees.push('0')
	}
	return trees.map(tree => sum([...tree].map(char => parseInt(char) || 0)));
}

export function getTalentPoints(talentsString: string): number {
	return sum(getTalentTreePoints(talentsString));
}

// Returns the index of the talent tree (0, 1, or 2) that has the most points.
export function getTalentTree(talentsString: string): number {
	const points = getTalentTreePoints(talentsString);
	return maxIndex(points) || 0;
}

enum IconSizes {
	Small = 'small',
	Medium = 'medium',
	Large = 'large',
}

// Gets the URL for the individual sim corresponding to the given spec.
const specSiteUrlTemplate = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/SPEC/`).toString();
export function getSpecSiteUrl(specString: string): string {
	return specSiteUrlTemplate.replace('SPEC', specString);
}
export const raidSimSiteUrl = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/raid/`).toString();

export function cssClassForClass(klass: Class): string {
	return klass.friendlyName.toLowerCase().replace(/\s/g, '-');
}

export function textCssClassForClass(klass: Class): string {
	return `text-${cssClassForClass(klass)}`;
}
export function textCssClassForSpec(spec: Spec): string {
	return textCssClassForClass(spec.class);
}

export type RotationUnion =
	BalanceDruidRotation |
	FeralDruidRotation |
	FeralTankDruidRotation |
	RestorationDruidRotation |
	HunterRotation |
	MageRotation |
	ElementalShamanRotation |
	EnhancementShamanRotation |
	RestorationShamanRotation |
	RogueRotation |
	HolyPaladinRotation |
	ProtectionPaladinRotation |
	RetributionPaladinRotation |
	HealingPriestRotation |
	ShadowPriestRotation |
	SmitePriestRotation |
	WarlockRotation |
	WarriorRotation |
	ProtectionWarriorRotation |
	DeathknightRotation |
	TankDeathknightRotation;
// TODO: Cata - Update list
export type SpecRotation<T extends SpecProto> =
	T extends SpecProto.SpecBalanceDruid ? BalanceDruidRotation :
	T extends SpecProto.SpecFeralDruid ? FeralDruidRotation :
	T extends SpecProto.SpecRestorationDruid ? RestorationDruidRotation :
	T extends SpecProto.SpecElementalShaman ? ElementalShamanRotation :
	T extends SpecProto.SpecEnhancementShaman ? EnhancementShamanRotation :
	T extends SpecProto.SpecRestorationShaman ? RestorationShamanRotation :
	T extends SpecProto.SpecHolyPaladin ? HolyPaladinRotation :
	T extends SpecProto.SpecProtectionPaladin ? ProtectionPaladinRotation :
	T extends SpecProto.SpecRetributionPaladin ? RetributionPaladinRotation :
	T extends SpecProto.SpecShadowPriest ? ShadowPriestRotation :
	T extends SpecProto.SpecProtectionWarrior ? ProtectionWarriorRotation :
	ElementalShamanRotation; // Should never reach this case

export type TalentsUnion =
	DruidTalents |
	HunterTalents |
	MageTalents |
	RogueTalents |
	PaladinTalents |
	PriestTalents |
	ShamanTalents |
	WarlockTalents |
	WarriorTalents |
	DeathknightTalents;
// TODO: Cata - Update list
export type SpecTalents<T extends SpecProto> =
	T extends SpecProto.SpecBalanceDruid ? DruidTalents :
	T extends SpecProto.SpecFeralDruid ? DruidTalents :
	T extends SpecProto.SpecRestorationDruid ? DruidTalents :
	T extends SpecProto.SpecElementalShaman ? ShamanTalents :
	T extends SpecProto.SpecEnhancementShaman ? ShamanTalents :
	T extends SpecProto.SpecRestorationShaman ? ShamanTalents :
	T extends SpecProto.SpecHolyPaladin ? PaladinTalents :
	T extends SpecProto.SpecProtectionPaladin ? PaladinTalents :
	T extends SpecProto.SpecRetributionPaladin ? PaladinTalents :
	T extends SpecProto.SpecShadowPriest ? PriestTalents :
	T extends SpecProto.SpecProtectionWarrior ? WarriorTalents :
	ShamanTalents; // Should never reach this case

export type SpecOptionsUnion =
	BalanceDruidOptions |
	FeralDruidOptions |
	FeralTankDruidOptions |
	RestorationDruidOptions |
	ElementalShamanOptions |
	EnhancementShamanOptions |
	RestorationShamanOptions |
	HunterOptions |
	MageOptions |
	RogueOptions |
	HolyPaladinOptions |
	ProtectionPaladinOptions |
	RetributionPaladinOptions |
	HealingPriestOptions |
	ShadowPriestOptions |
	SmitePriestOptions |
	WarlockOptions |
	WarriorOptions |
	ProtectionWarriorOptions |
	DeathknightOptions |
	TankDeathknightOptions;
// TODO: Cata - Update list
export type SpecOptions<T extends SpecProto> =
	T extends SpecProto.SpecBalanceDruid ? BalanceDruidOptions :
	T extends SpecProto.SpecFeralDruid ? FeralDruidOptions :
	T extends SpecProto.SpecRestorationDruid ? RestorationDruidOptions :
	T extends SpecProto.SpecElementalShaman ? ElementalShamanOptions :
	T extends SpecProto.SpecEnhancementShaman ? EnhancementShamanOptions :
	T extends SpecProto.SpecRestorationShaman ? RestorationShamanOptions :
	T extends SpecProto.SpecHolyPaladin ? HolyPaladinOptions :
	T extends SpecProto.SpecProtectionPaladin ? ProtectionPaladinOptions :
	T extends SpecProto.SpecRetributionPaladin ? RetributionPaladinOptions :
	T extends SpecProto.SpecShadowPriest ? ShadowPriestOptions :
	T extends SpecProto.SpecProtectionWarrior ? ProtectionWarriorOptions :
	ElementalShamanOptions; // Should never reach this case

export type SpecProtoUnion =
	BalanceDruid |
	FeralDruid |
	FeralTankDruid |
	RestorationDruid |
	ElementalShaman |
	EnhancementShaman |
	RestorationShaman |
	Hunter |
	Mage |
	Rogue |
	HolyPaladin |
	ProtectionPaladin |
	RetributionPaladin |
	HealingPriest |
	ShadowPriest |
	SmitePriest |
	Warlock |
	Warrior |
	ProtectionWarrior |
	Deathknight |
	TankDeathknight;
// TODO: Cata - Update list
export type SpecProtoType<T extends SpecProto> =
	T extends SpecProto.SpecBalanceDruid ? BalanceDruid :
	T extends SpecProto.SpecFeralDruid ? FeralDruid :
	T extends SpecProto.SpecRestorationDruid ? RestorationDruid :
	T extends SpecProto.SpecElementalShaman ? ElementalShaman :
	T extends SpecProto.SpecEnhancementShaman ? EnhancementShaman :
	T extends SpecProto.SpecRestorationShaman ? RestorationShaman :
	T extends SpecProto.SpecHolyPaladin ? HolyPaladin :
	T extends SpecProto.SpecProtectionPaladin ? ProtectionPaladin :
	T extends SpecProto.SpecRetributionPaladin ? RetributionPaladin :
	T extends SpecProto.SpecShadowPriest ? ShadowPriest :
	T extends SpecProto.SpecProtectionWarrior ? ProtectionWarrior :
	ElementalShaman; // Should never reach this case

export type SpecTypeFunctions<SpecType extends SpecProto> = {
	rotationCreate: () => SpecRotation<SpecType>;
	rotationEquals: (a: SpecRotation<SpecType>, b: SpecRotation<SpecType>) => boolean;
	rotationCopy: (a: SpecRotation<SpecType>) => SpecRotation<SpecType>;
	rotationToJson: (a: SpecRotation<SpecType>) => any;
	rotationFromJson: (obj: any) => SpecRotation<SpecType>;

	talentsCreate: () => SpecTalents<SpecType>;
	talentsEquals: (a: SpecTalents<SpecType>, b: SpecTalents<SpecType>) => boolean;
	talentsCopy: (a: SpecTalents<SpecType>) => SpecTalents<SpecType>;
	talentsToJson: (a: SpecTalents<SpecType>) => any;
	talentsFromJson: (obj: any) => SpecTalents<SpecType>;

	optionsCreate: () => SpecOptions<SpecType>;
	optionsEquals: (a: SpecOptions<SpecType>, b: SpecOptions<SpecType>) => boolean;
	optionsCopy: (a: SpecOptions<SpecType>) => SpecOptions<SpecType>;
	optionsToJson: (a: SpecOptions<SpecType>) => any;
	optionsFromJson: (obj: any) => SpecOptions<SpecType>;
	optionsFromPlayer: (player: Player) => SpecOptions<SpecType>;
};

// TODO: Cata - Update list
export const specTypeFunctions: Record<SpecProto, SpecTypeFunctions<any>> = {
	[SpecProto.SpecBalanceDruid]: {
		rotationCreate: () => BalanceDruidRotation.create(),
		rotationEquals: (a, b) => BalanceDruidRotation.equals(a as BalanceDruidRotation, b as BalanceDruidRotation),
		rotationCopy: a => BalanceDruidRotation.clone(a as BalanceDruidRotation),
		rotationToJson: a => BalanceDruidRotation.toJson(a as BalanceDruidRotation),
		rotationFromJson: obj => BalanceDruidRotation.fromJson(obj),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: a => DruidTalents.clone(a as DruidTalents),
		talentsToJson: a => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: obj => DruidTalents.fromJson(obj),

		optionsCreate: () => BalanceDruidOptions.create(),
		optionsEquals: (a, b) => BalanceDruidOptions.equals(a as BalanceDruidOptions, b as BalanceDruidOptions),
		optionsCopy: a => BalanceDruidOptions.clone(a as BalanceDruidOptions),
		optionsToJson: a => BalanceDruidOptions.toJson(a as BalanceDruidOptions),
		optionsFromJson: obj => BalanceDruidOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'balanceDruid'
			? player.spec.balanceDruid.options || BalanceDruidOptions.create()
			: BalanceDruidOptions.create(),
	},
	[SpecProto.SpecFeralDruid]: {
		rotationCreate: () => FeralDruidRotation.create(),
		rotationEquals: (a, b) => FeralDruidRotation.equals(a as FeralDruidRotation, b as FeralDruidRotation),
		rotationCopy: a => FeralDruidRotation.clone(a as FeralDruidRotation),
		rotationToJson: a => FeralDruidRotation.toJson(a as FeralDruidRotation),
		rotationFromJson: obj => FeralDruidRotation.fromJson(obj),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: a => DruidTalents.clone(a as DruidTalents),
		talentsToJson: a => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: obj => DruidTalents.fromJson(obj),

		optionsCreate: () => FeralDruidOptions.create(),
		optionsEquals: (a, b) => FeralDruidOptions.equals(a as FeralDruidOptions, b as FeralDruidOptions),
		optionsCopy: a => FeralDruidOptions.clone(a as FeralDruidOptions),
		optionsToJson: a => FeralDruidOptions.toJson(a as FeralDruidOptions),
		optionsFromJson: obj => FeralDruidOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'feralDruid'
			? player.spec.feralDruid.options || FeralDruidOptions.create()
			: FeralDruidOptions.create(),
	},
	[SpecProto.SpecRestorationDruid]: {
		rotationCreate: () => RestorationDruidRotation.create(),
		rotationEquals: (a, b) => RestorationDruidRotation.equals(a as RestorationDruidRotation, b as RestorationDruidRotation),
		rotationCopy: a => RestorationDruidRotation.clone(a as RestorationDruidRotation),
		rotationToJson: a => RestorationDruidRotation.toJson(a as RestorationDruidRotation),
		rotationFromJson: obj => RestorationDruidRotation.fromJson(obj),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: a => DruidTalents.clone(a as DruidTalents),
		talentsToJson: a => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: obj => DruidTalents.fromJson(obj),

		optionsCreate: () => RestorationDruidOptions.create(),
		optionsEquals: (a, b) => RestorationDruidOptions.equals(a as RestorationDruidOptions, b as RestorationDruidOptions),
		optionsCopy: a => RestorationDruidOptions.clone(a as RestorationDruidOptions),
		optionsToJson: a => RestorationDruidOptions.toJson(a as RestorationDruidOptions),
		optionsFromJson: obj => RestorationDruidOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'restorationDruid'
			? player.spec.restorationDruid.options || RestorationDruidOptions.create()
			: RestorationDruidOptions.create(),
	},
	[SpecProto.SpecElementalShaman]: {
		rotationCreate: () => ElementalShamanRotation.create(),
		rotationEquals: (a, b) => ElementalShamanRotation.equals(a as ElementalShamanRotation, b as ElementalShamanRotation),
		rotationCopy: a => ElementalShamanRotation.clone(a as ElementalShamanRotation),
		rotationToJson: a => ElementalShamanRotation.toJson(a as ElementalShamanRotation),
		rotationFromJson: obj => ElementalShamanRotation.fromJson(obj),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: a => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: a => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: obj => ShamanTalents.fromJson(obj),

		optionsCreate: () => ElementalShamanOptions.create(),
		optionsEquals: (a, b) => ElementalShamanOptions.equals(a as ElementalShamanOptions, b as ElementalShamanOptions),
		optionsCopy: a => ElementalShamanOptions.clone(a as ElementalShamanOptions),
		optionsToJson: a => ElementalShamanOptions.toJson(a as ElementalShamanOptions),
		optionsFromJson: obj => ElementalShamanOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'elementalShaman'
			? player.spec.elementalShaman.options || ElementalShamanOptions.create()
			: ElementalShamanOptions.create(),
	},
	[SpecProto.SpecEnhancementShaman]: {
		rotationCreate: () => EnhancementShamanRotation.create(),
		rotationEquals: (a, b) => EnhancementShamanRotation.equals(a as EnhancementShamanRotation, b as EnhancementShamanRotation),
		rotationCopy: a => EnhancementShamanRotation.clone(a as EnhancementShamanRotation),
		rotationToJson: a => EnhancementShamanRotation.toJson(a as EnhancementShamanRotation),
		rotationFromJson: obj => EnhancementShamanRotation.fromJson(obj),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: a => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: a => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: obj => ShamanTalents.fromJson(obj),

		optionsCreate: () => EnhancementShamanOptions.create(),
		optionsEquals: (a, b) => EnhancementShamanOptions.equals(a as EnhancementShamanOptions, b as EnhancementShamanOptions),
		optionsCopy: a => EnhancementShamanOptions.clone(a as EnhancementShamanOptions),
		optionsToJson: a => EnhancementShamanOptions.toJson(a as EnhancementShamanOptions),
		optionsFromJson: obj => EnhancementShamanOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'enhancementShaman'
			? player.spec.enhancementShaman.options || EnhancementShamanOptions.create()
			: EnhancementShamanOptions.create(),
	},
	[SpecProto.SpecRestorationShaman]: {
		rotationCreate: () => RestorationShamanRotation.create(),
		rotationEquals: (a, b) => RestorationShamanRotation.equals(a as RestorationShamanRotation, b as RestorationShamanRotation),
		rotationCopy: a => RestorationShamanRotation.clone(a as RestorationShamanRotation),
		rotationToJson: a => RestorationShamanRotation.toJson(a as RestorationShamanRotation),
		rotationFromJson: obj => RestorationShamanRotation.fromJson(obj),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: a => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: a => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: obj => ShamanTalents.fromJson(obj),

		optionsCreate: () => RestorationShamanOptions.create(),
		optionsEquals: (a, b) => RestorationShamanOptions.equals(a as RestorationShamanOptions, b as RestorationShamanOptions),
		optionsCopy: a => RestorationShamanOptions.clone(a as RestorationShamanOptions),
		optionsToJson: a => RestorationShamanOptions.toJson(a as RestorationShamanOptions),
		optionsFromJson: obj => RestorationShamanOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'restorationShaman'
			? player.spec.restorationShaman.options || RestorationShamanOptions.create()
			: RestorationShamanOptions.create(),
	},
	[SpecProto.SpecHolyPaladin]: {
		rotationCreate: () => HolyPaladinRotation.create(),
		rotationEquals: (a, b) => HolyPaladinRotation.equals(a as HolyPaladinRotation, b as HolyPaladinRotation),
		rotationCopy: a => HolyPaladinRotation.clone(a as HolyPaladinRotation),
		rotationToJson: a => HolyPaladinRotation.toJson(a as HolyPaladinRotation),
		rotationFromJson: obj => HolyPaladinRotation.fromJson(obj),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: a => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: a => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: obj => PaladinTalents.fromJson(obj),

		optionsCreate: () => HolyPaladinOptions.create(),
		optionsEquals: (a, b) => HolyPaladinOptions.equals(a as HolyPaladinOptions, b as HolyPaladinOptions),
		optionsCopy: a => HolyPaladinOptions.clone(a as HolyPaladinOptions),
		optionsToJson: a => HolyPaladinOptions.toJson(a as HolyPaladinOptions),
		optionsFromJson: obj => HolyPaladinOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'holyPaladin'
			? player.spec.holyPaladin.options || HolyPaladinOptions.create()
			: HolyPaladinOptions.create(),
	},
	[SpecProto.SpecProtectionPaladin]: {
		rotationCreate: () => ProtectionPaladinRotation.create(),
		rotationEquals: (a, b) => ProtectionPaladinRotation.equals(a as ProtectionPaladinRotation, b as ProtectionPaladinRotation),
		rotationCopy: a => ProtectionPaladinRotation.clone(a as ProtectionPaladinRotation),
		rotationToJson: a => ProtectionPaladinRotation.toJson(a as ProtectionPaladinRotation),
		rotationFromJson: obj => ProtectionPaladinRotation.fromJson(obj),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: a => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: a => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: obj => PaladinTalents.fromJson(obj),

		optionsCreate: () => ProtectionPaladinOptions.create(),
		optionsEquals: (a, b) => ProtectionPaladinOptions.equals(a as ProtectionPaladinOptions, b as ProtectionPaladinOptions),
		optionsCopy: a => ProtectionPaladinOptions.clone(a as ProtectionPaladinOptions),
		optionsToJson: a => ProtectionPaladinOptions.toJson(a as ProtectionPaladinOptions),
		optionsFromJson: obj => ProtectionPaladinOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'protectionPaladin'
			? player.spec.protectionPaladin.options || ProtectionPaladinOptions.create()
			: ProtectionPaladinOptions.create(),
	},
	[SpecProto.SpecRetributionPaladin]: {
		rotationCreate: () => RetributionPaladinRotation.create(),
		rotationEquals: (a, b) => RetributionPaladinRotation.equals(a as RetributionPaladinRotation, b as RetributionPaladinRotation),
		rotationCopy: a => RetributionPaladinRotation.clone(a as RetributionPaladinRotation),
		rotationToJson: a => RetributionPaladinRotation.toJson(a as RetributionPaladinRotation),
		rotationFromJson: obj => RetributionPaladinRotation.fromJson(obj),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: a => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: a => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: obj => PaladinTalents.fromJson(obj),

		optionsCreate: () => RetributionPaladinOptions.create(),
		optionsEquals: (a, b) => RetributionPaladinOptions.equals(a as RetributionPaladinOptions, b as RetributionPaladinOptions),
		optionsCopy: a => RetributionPaladinOptions.clone(a as RetributionPaladinOptions),
		optionsToJson: a => RetributionPaladinOptions.toJson(a as RetributionPaladinOptions),
		optionsFromJson: obj => RetributionPaladinOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'retributionPaladin'
			? player.spec.retributionPaladin.options || RetributionPaladinOptions.create()
			: RetributionPaladinOptions.create(),
	},
	[SpecProto.SpecShadowPriest]: {
		rotationCreate: () => ShadowPriestRotation.create(),
		rotationEquals: (a, b) => ShadowPriestRotation.equals(a as ShadowPriestRotation, b as ShadowPriestRotation),
		rotationCopy: a => ShadowPriestRotation.clone(a as ShadowPriestRotation),
		rotationToJson: a => ShadowPriestRotation.toJson(a as ShadowPriestRotation),
		rotationFromJson: obj => ShadowPriestRotation.fromJson(obj),

		talentsCreate: () => PriestTalents.create(),
		talentsEquals: (a, b) => PriestTalents.equals(a as PriestTalents, b as PriestTalents),
		talentsCopy: a => PriestTalents.clone(a as PriestTalents),
		talentsToJson: a => PriestTalents.toJson(a as PriestTalents),
		talentsFromJson: obj => PriestTalents.fromJson(obj),

		optionsCreate: () => ShadowPriestOptions.create(),
		optionsEquals: (a, b) => ShadowPriestOptions.equals(a as ShadowPriestOptions, b as ShadowPriestOptions),
		optionsCopy: a => ShadowPriestOptions.clone(a as ShadowPriestOptions),
		optionsToJson: a => ShadowPriestOptions.toJson(a as ShadowPriestOptions),
		optionsFromJson: obj => ShadowPriestOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'shadowPriest'
			? player.spec.shadowPriest.options || ShadowPriestOptions.create()
			: ShadowPriestOptions.create(),
	},
	[SpecProto.SpecProtectionWarrior]: {
		rotationCreate: () => ProtectionWarriorRotation.create(),
		rotationEquals: (a, b) => ProtectionWarriorRotation.equals(a as ProtectionWarriorRotation, b as ProtectionWarriorRotation),
		rotationCopy: a => ProtectionWarriorRotation.clone(a as ProtectionWarriorRotation),
		rotationToJson: a => ProtectionWarriorRotation.toJson(a as ProtectionWarriorRotation),
		rotationFromJson: obj => ProtectionWarriorRotation.fromJson(obj),

		talentsCreate: () => WarriorTalents.create(),
		talentsEquals: (a, b) => WarriorTalents.equals(a as WarriorTalents, b as WarriorTalents),
		talentsCopy: a => WarriorTalents.clone(a as WarriorTalents),
		talentsToJson: a => WarriorTalents.toJson(a as WarriorTalents),
		talentsFromJson: obj => WarriorTalents.fromJson(obj),

		optionsCreate: () => ProtectionWarriorOptions.create(),
		optionsEquals: (a, b) => ProtectionWarriorOptions.equals(a as ProtectionWarriorOptions, b as ProtectionWarriorOptions),
		optionsCopy: a => ProtectionWarriorOptions.clone(a as ProtectionWarriorOptions),
		optionsToJson: a => ProtectionWarriorOptions.toJson(a as ProtectionWarriorOptions),
		optionsFromJson: obj => ProtectionWarriorOptions.fromJson(obj),
		optionsFromPlayer: player => player.spec.oneofKind == 'protectionWarrior'
			? player.spec.protectionWarrior.options || ProtectionWarriorOptions.create()
			: ProtectionWarriorOptions.create(),
	},
};

export const raceToFaction: Record<Race, Faction> = {
	[Race.RaceUnknown]: Faction.Unknown,
	[Race.RaceBloodElf]: Faction.Horde,
	[Race.RaceDraenei]: Faction.Alliance,
	[Race.RaceDwarf]: Faction.Alliance,
	[Race.RaceGnome]: Faction.Alliance,
	[Race.RaceHuman]: Faction.Alliance,
	[Race.RaceNightElf]: Faction.Alliance,
	[Race.RaceOrc]: Faction.Horde,
	[Race.RaceTauren]: Faction.Horde,
	[Race.RaceTroll]: Faction.Horde,
	[Race.RaceUndead]: Faction.Horde,
};

export const specToClass: Record<SpecProto, Class> = {
	[SpecProto.SpecBalanceDruid]: ClassProto.ClassDruid,
	[SpecProto.SpecFeralDruid]: ClassProto.ClassDruid,
	[SpecProto.SpecFeralTankDruid]: ClassProto.ClassDruid,
	[SpecProto.SpecRestorationDruid]: ClassProto.ClassDruid,
	[SpecProto.SpecHunter]: ClassProto.ClassHunter,
	[SpecProto.SpecMage]: ClassProto.ClassMage,
	[SpecProto.SpecRogue]: ClassProto.ClassRogue,
	[SpecProto.SpecHolyPaladin]: ClassProto.ClassPaladin,
	[SpecProto.SpecProtectionPaladin]: ClassProto.ClassPaladin,
	[SpecProto.SpecRetributionPaladin]: ClassProto.ClassPaladin,
	[SpecProto.SpecHealingPriest]: ClassProto.ClassPriest,
	[SpecProto.SpecShadowPriest]: ClassProto.ClassPriest,
	[SpecProto.SpecSmitePriest]: ClassProto.ClassPriest,
	[SpecProto.SpecElementalShaman]: ClassProto.ClassShaman,
	[SpecProto.SpecEnhancementShaman]: ClassProto.ClassShaman,
	[SpecProto.SpecRestorationShaman]: ClassProto.ClassShaman,
	[SpecProto.SpecWarlock]: ClassProto.ClassWarlock,
	[SpecProto.SpecWarrior]: ClassProto.ClassWarrior,
	[SpecProto.SpecProtectionWarrior]: ClassProto.ClassWarrior,
	[SpecProto.SpecDeathknight]: ClassProto.ClassDeathknight,
	[SpecProto.SpecTankDeathknight]: ClassProto.ClassDeathknight,
};

// Specs that can dual wield. This could be based on class, except that
// Enhancement Shaman learn dual wield from a talent.
const dualWieldSpecs: Array<Spec> = [
	Spec.SpecEnhancementShaman,
	Spec.SpecHunter,
	Spec.SpecRogue,
	Spec.SpecWarrior,
	Spec.SpecProtectionWarrior,
	Spec.SpecDeathknight,
	Spec.SpecTankDeathknight,
];
export function isDualWieldSpec(spec: Spec): boolean {
	return dualWieldSpecs.includes(spec);
}

// TODO: Cata - Update list
const tankSpecs: Array<SpecProto> = [
	SpecProto.SpecProtectionPaladin,
	SpecProto.SpecProtectionWarrior,
];
export function isTankSpec(spec: SpecProto): boolean {
	return tankSpecs.includes(spec);
}

// TODO: Cata - Update list
const healingSpecs: Array<SpecProto> = [
	SpecProto.SpecRestorationDruid,
	SpecProto.SpecHolyPaladin,
	SpecProto.SpecRestorationShaman,
];
export function isHealingSpec(spec: SpecProto): boolean {
	return healingSpecs.includes(spec);
}

// TODO: Cata - Update list
const rangedDpsSpecs: Array<SpecProto> = [
	SpecProto.SpecBalanceDruid,
	SpecProto.SpecShadowPriest,
	SpecProto.SpecElementalShaman,
];
export function isRangedDpsSpec(spec: SpecProto): boolean {
	return rangedDpsSpecs.includes(spec);
}
export function isMeleeDpsSpec(spec: SpecProto): boolean {
	return !isTankSpec(spec) && !isHealingSpec(spec) && !isRangedDpsSpec(spec);
}

// Prefixes used for storing browser data for each site. Even if a Spec is
// renamed, DO NOT change these values or people will lose their saved data.
// TODO: Cata - Build these programmatically
export const specToLocalStorageKey: Record<SpecProto, string> = {
	[SpecProto.SpecBalanceDruid]: '__cata_balance_druid',
	[SpecProto.SpecFeralDruid]: '__cata_feral_druid',
	[SpecProto.SpecRestorationDruid]: '__cata_restoration_druid',
	[SpecProto.SpecElementalShaman]: '__cata_elemental_shaman',
	[SpecProto.SpecEnhancementShaman]: '__cata_enhacement_shaman',
	[SpecProto.SpecRestorationShaman]: '__cata_restoration_shaman',
	[SpecProto.SpecHolyPaladin]: '__cata_holy_paladin',
	[SpecProto.SpecProtectionPaladin]: '__cata_protection_paladin',
	[SpecProto.SpecRetributionPaladin]: '__cata_retribution_paladin',
	[SpecProto.SpecShadowPriest]: '__cata_shadow_priest',
	[SpecProto.SpecProtectionWarrior]: '__cata_protection_warrior',
};

// Returns a copy of playerOptions, with the class field set.
// TODO: Cata - Update list
export function withSpecProto<SpecType extends SpecProto>(
	spec: SpecProto,
	player: Player,
	specOptions: SpecOptions<SpecType>): Player {
	const copy = Player.clone(player);

	switch (spec) {
		case SpecProto.SpecBalanceDruid:
			copy.spec = {
				oneofKind: 'balanceDruid',
				balanceDruid: BalanceDruid.create({
					options: specOptions as BalanceDruidOptions,
				}),
			};
			return copy;
		case SpecProto.SpecFeralDruid:
			copy.spec = {
				oneofKind: 'feralDruid',
				feralDruid: FeralDruid.create({
					options: specOptions as FeralDruidOptions,
				}),
			};
			return copy;
		case SpecProto.SpecRestorationDruid:
			copy.spec = {
				oneofKind: 'restorationDruid',
				restorationDruid: RestorationDruid.create({
					options: specOptions as RestorationDruidOptions,
				}),
			};
			return copy;
		case SpecProto.SpecElementalShaman:
			copy.spec = {
				oneofKind: 'elementalShaman',
				elementalShaman: ElementalShaman.create({
					options: specOptions as ElementalShamanOptions,
				}),
			};
			return copy;
		case SpecProto.SpecEnhancementShaman:
			copy.spec = {
				oneofKind: 'enhancementShaman',
				enhancementShaman: EnhancementShaman.create({
					options: specOptions as ElementalShamanOptions,
				}),
			};
			return copy;
		case SpecProto.SpecRestorationShaman:
			copy.spec = {
				oneofKind: 'restorationShaman',
				restorationShaman: RestorationShaman.create({
					options: specOptions as RestorationShamanOptions,
				}),
			};
			return copy;
		case SpecProto.SpecHolyPaladin:
			copy.spec = {
				oneofKind: 'holyPaladin',
				holyPaladin: HolyPaladin.create({
					options: specOptions as HolyPaladinOptions,
				}),
			};
			return copy;
		case SpecProto.SpecProtectionPaladin:
			copy.spec = {
				oneofKind: 'protectionPaladin',
				protectionPaladin: ProtectionPaladin.create({
					options: specOptions as ProtectionPaladinOptions,
				}),
			};
			return copy;
		case SpecProto.SpecRetributionPaladin:
			copy.spec = {
				oneofKind: 'retributionPaladin',
				retributionPaladin: RetributionPaladin.create({
					options: specOptions as RetributionPaladinOptions,
				}),
			};
			return copy;
		case SpecProto.SpecShadowPriest:
			copy.spec = {
				oneofKind: 'shadowPriest',
				shadowPriest: ShadowPriest.create({
					options: specOptions as ShadowPriestOptions,
				}),
			};
			return copy;
		case SpecProto.SpecProtectionWarrior:
			copy.spec = {
				oneofKind: 'protectionWarrior',
				protectionWarrior: ProtectionWarrior.create({
					options: specOptions as ProtectionWarriorOptions,
				}),
			};
			return copy;
	}
}

export function playerToSpec(player: Player): Spec {
	const specValues = getEnumValues(Spec);
	for (let i = 0; i < specValues.length; i++) {
		const spec = specValues[i] as Spec;
		let specString = Spec[spec]; // Returns 'SpecBalanceDruid' for BalanceDruid.
		specString = specString.substring('Spec'.length); // 'BalanceDruid'
		specString = specString.charAt(0).toLowerCase() + specString.slice(1); // 'balanceDruid'

		if (player.spec.oneofKind == specString) {
			return spec;
		}
	}

	throw new Error('Unable to parse spec from player proto: ' + JSON.stringify(Player.toJson(player), null, 2));
}

export const classToMaxArmorType: Record<Class, ArmorType> = {
	[Class.ClassUnknown]: ArmorType.ArmorTypeUnknown,
	[Class.ClassDruid]: ArmorType.ArmorTypeLeather,
	[Class.ClassHunter]: ArmorType.ArmorTypeMail,
	[Class.ClassMage]: ArmorType.ArmorTypeCloth,
	[Class.ClassPaladin]: ArmorType.ArmorTypePlate,
	[Class.ClassPriest]: ArmorType.ArmorTypeCloth,
	[Class.ClassRogue]: ArmorType.ArmorTypeLeather,
	[Class.ClassShaman]: ArmorType.ArmorTypeMail,
	[Class.ClassWarlock]: ArmorType.ArmorTypeCloth,
	[Class.ClassWarrior]: ArmorType.ArmorTypePlate,
	[Class.ClassDeathknight]: ArmorType.ArmorTypePlate,
};

export const classToEligibleRangedWeaponTypes: Record<Class, Array<RangedWeaponType>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [RangedWeaponType.RangedWeaponTypeIdol],
	[Class.ClassHunter]: [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
	],
	[Class.ClassMage]: [RangedWeaponType.RangedWeaponTypeWand],
	[Class.ClassPaladin]: [RangedWeaponType.RangedWeaponTypeLibram],
	[Class.ClassPriest]: [RangedWeaponType.RangedWeaponTypeWand],
	[Class.ClassRogue]: [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
		RangedWeaponType.RangedWeaponTypeThrown,
	],
	[Class.ClassShaman]: [RangedWeaponType.RangedWeaponTypeTotem],
	[Class.ClassWarlock]: [RangedWeaponType.RangedWeaponTypeWand],
	[Class.ClassWarrior]: [
		RangedWeaponType.RangedWeaponTypeBow,
		RangedWeaponType.RangedWeaponTypeCrossbow,
		RangedWeaponType.RangedWeaponTypeGun,
		RangedWeaponType.RangedWeaponTypeThrown,
	],
	[Class.ClassDeathknight]: [
		RangedWeaponType.RangedWeaponTypeSigil,
	],
};

interface EligibleWeaponType {
	weaponType: WeaponType,
	canUseTwoHand?: boolean,
}

export const classToEligibleWeaponTypes: Record<Class, Array<EligibleWeaponType>> = {
	[Class.ClassUnknown]: [],
	[Class.ClassDruid]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
	],
	[Class.ClassHunter]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	],
	[Class.ClassMage]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword },
	],
	[Class.ClassPaladin]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
	],
	[Class.ClassPriest]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeMace },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
	],
	[Class.ClassRogue]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: false },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeSword },
	],
	[Class.ClassWarlock]: [
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword },
	],
	[Class.ClassWarrior]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeDagger },
		{ weaponType: WeaponType.WeaponTypeFist },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeOffHand },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeShield },
		{ weaponType: WeaponType.WeaponTypeStaff, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
	],
	[Class.ClassDeathknight]: [
		{ weaponType: WeaponType.WeaponTypeAxe, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeMace, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypePolearm, canUseTwoHand: true },
		{ weaponType: WeaponType.WeaponTypeSword, canUseTwoHand: true },
		// TODO: validate proficiencies
	],
};

export function isSharpWeaponType(weaponType: WeaponType): boolean {
	return [
		WeaponType.WeaponTypeAxe,
		WeaponType.WeaponTypeDagger,
		WeaponType.WeaponTypePolearm,
		WeaponType.WeaponTypeSword,
	].includes(weaponType);
}

export function isBluntWeaponType(weaponType: WeaponType): boolean {
	return [
		WeaponType.WeaponTypeFist,
		WeaponType.WeaponTypeMace,
		WeaponType.WeaponTypeStaff,
	].includes(weaponType);
}

// Custom functions for determining the EP value of meta gem effects.
// Default meta effect EP value is 0, so just handle the ones relevant to your spec.
const metaGemEffectEPs: Partial<Record<Spec, (gem: Gem, playerStats: Stats) => number>> = {
	[Spec.SpecBalanceDruid]: (gem, _) => {
		if (gem.id == Gems.CHAOTIC_SKYFIRE_DIAMOND.id) {
			// TODO: Fix this
			return (12 * 0.65) + (3 * 45);
		}
		if (gem.id == Gems.CHAOTIC_SKYFLARE_DIAMOND.id) {
			return (21 * 0.65) + (3 * 45);
		}
		return 0;
	},
	[Spec.SpecElementalShaman]: (gem, _) => {
		if (gem.id == Gems.CHAOTIC_SKYFLARE_DIAMOND.id) {
			return 84;
		}
		if (gem.id == Gems.CHAOTIC_SKYFIRE_DIAMOND.id) {
			return 80;
		}

		return 0;
	},
	[Spec.SpecWarlock]: (gem, _) => {
		// TODO: make it gear dependant
		if (gem.id == Gems.CHAOTIC_SKYFLARE_DIAMOND.id) {
			return 84;
		}
		if (gem.id == Gems.CHAOTIC_SKYFIRE_DIAMOND.id) {
			return 80;
		}

		return 0;
	},
	[Spec.SpecFeralDruid]: (gem, _) => {
		// Unknown actual EP, but this is the only effect that matters
		if (gem.id == Gems.RELENTLESS_EARTHSIEGE_DIAMOND.id || gem.id == Gems.CHAOTIC_SKYFLARE_DIAMOND.id || gem.id == Gems.CHAOTIC_SKYFIRE_DIAMOND.id) {
			return 80;
		}
		return 0;
	}
};

export function getMetaGemEffectEP(spec: Spec, gem: Gem, playerStats: Stats) {
	if (metaGemEffectEPs[spec]) {
		return metaGemEffectEPs[spec]!(gem, playerStats);
	} else {
		return 0;
	}
}

// Returns true if this item may be equipped in at least 1 slot for the given Spec.
export function canEquipItem(item: Item, spec: Spec, slot: ItemSlot | undefined): boolean {
	const playerClass = specToClass[spec];
	if (item.classAllowlist.length > 0 && !item.classAllowlist.includes(playerClass)) {
		return false;
	}

	if ([ItemType.ItemTypeFinger, ItemType.ItemTypeTrinket].includes(item.type)) {
		return true;
	}

	if (item.type == ItemType.ItemTypeWeapon) {
		const eligibleWeaponType = classToEligibleWeaponTypes[playerClass].find(wt => wt.weaponType == item.weaponType);
		if (!eligibleWeaponType) {
			return false;
		}

		if ((item.handType == HandType.HandTypeOffHand || (item.handType == HandType.HandTypeOneHand && slot == ItemSlot.ItemSlotOffHand))
			&& ![WeaponType.WeaponTypeShield, WeaponType.WeaponTypeOffHand].includes(item.weaponType)
			&& !dualWieldSpecs.includes(spec)) {
			return false;
		}

		if (item.handType == HandType.HandTypeTwoHand && !eligibleWeaponType.canUseTwoHand) {
			return false;
		}
		if (item.handType == HandType.HandTypeTwoHand && slot == ItemSlot.ItemSlotOffHand && spec != Spec.SpecWarrior) {
			return false;
		}

		return true;
	}

	if (item.type == ItemType.ItemTypeRanged) {
		return classToEligibleRangedWeaponTypes[playerClass].includes(item.rangedWeaponType);
	}

	// At this point, we know the item is an armor piece (feet, chest, legs, etc).
	return classToMaxArmorType[playerClass] >= item.armorType;
}

const itemTypeToSlotsMap: Partial<Record<ItemType, Array<ItemSlot>>> = {
	[ItemType.ItemTypeUnknown]: [],
	[ItemType.ItemTypeHead]: [ItemSlot.ItemSlotHead],
	[ItemType.ItemTypeNeck]: [ItemSlot.ItemSlotNeck],
	[ItemType.ItemTypeShoulder]: [ItemSlot.ItemSlotShoulder],
	[ItemType.ItemTypeBack]: [ItemSlot.ItemSlotBack],
	[ItemType.ItemTypeChest]: [ItemSlot.ItemSlotChest],
	[ItemType.ItemTypeWrist]: [ItemSlot.ItemSlotWrist],
	[ItemType.ItemTypeHands]: [ItemSlot.ItemSlotHands],
	[ItemType.ItemTypeWaist]: [ItemSlot.ItemSlotWaist],
	[ItemType.ItemTypeLegs]: [ItemSlot.ItemSlotLegs],
	[ItemType.ItemTypeFeet]: [ItemSlot.ItemSlotFeet],
	[ItemType.ItemTypeFinger]: [ItemSlot.ItemSlotFinger1, ItemSlot.ItemSlotFinger2],
	[ItemType.ItemTypeTrinket]: [ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	[ItemType.ItemTypeRanged]: [ItemSlot.ItemSlotRanged],
};

export function getEligibleItemSlots(item: Item): Array<ItemSlot> {
	if (itemTypeToSlotsMap[item.type]) {
		return itemTypeToSlotsMap[item.type]!;
	}

	if (item.type == ItemType.ItemTypeWeapon) {
		if (item.handType == HandType.HandTypeMainHand) {
			return [ItemSlot.ItemSlotMainHand];
		} else if (item.handType == HandType.HandTypeOffHand) {
			return [ItemSlot.ItemSlotOffHand];
			// Missing HandTypeTwoHand
			// We allow 2H weapons to be wielded in mainhand and offhand for Fury Warriors
		} else {
			return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
		}
	}

	// Should never reach here
	throw new Error('Could not find item slots for item: ' + Item.toJsonString(item));
};

// Returns whether the given main-hand and off-hand items can be worn at the
// same time.
export function validWeaponCombo(mainHand: Item | null | undefined, offHand: Item | null | undefined, canDW2h: boolean): boolean {
	if (mainHand == null || offHand == null) {
		return true;
	}

	if (mainHand.handType == HandType.HandTypeTwoHand && !canDW2h) {
		return false;
	} else if (mainHand.handType == HandType.HandTypeTwoHand &&
		(mainHand.weaponType == WeaponType.WeaponTypePolearm || mainHand.weaponType == WeaponType.WeaponTypeStaff)) {
		return false;
	}

	if (offHand.handType == HandType.HandTypeTwoHand && !canDW2h) {
		return false;
	} else if (offHand.handType == HandType.HandTypeTwoHand &&
		(offHand.weaponType == WeaponType.WeaponTypePolearm || offHand.weaponType == WeaponType.WeaponTypeStaff)) {
		return false;
	}

	return true;
}

// Returns all item slots to which the enchant might be applied.
//
// Note that this alone is not enough; some items have further restrictions,
// e.g. some weapon enchants may only be applied to 2H weapons.
export function getEligibleEnchantSlots(enchant: Enchant): Array<ItemSlot> {
	return [enchant.type].concat(enchant.extraTypes || []).map(type => {
		if (itemTypeToSlotsMap[type]) {
			return itemTypeToSlotsMap[type]!;
		}

		if (type == ItemType.ItemTypeWeapon) {
			return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
		}

		// Should never reach here
		throw new Error('Could not find item slots for enchant: ' + Enchant.toJsonString(enchant));
	}).flat();
};

export function enchantAppliesToItem(enchant: Enchant, item: Item): boolean {
	const sharedSlots = intersection(getEligibleEnchantSlots(enchant), getEligibleItemSlots(item));
	if (sharedSlots.length == 0)
		return false;

	if (enchant.enchantType == EnchantType.EnchantTypeTwoHand && item.handType != HandType.HandTypeTwoHand)
		return false;

	if ((enchant.enchantType == EnchantType.EnchantTypeShield) != (item.weaponType == WeaponType.WeaponTypeShield))
		return false;

	if (enchant.enchantType == EnchantType.EnchantTypeStaff && item.weaponType != WeaponType.WeaponTypeStaff)
		return false;

	if (item.weaponType == WeaponType.WeaponTypeOffHand)
		return false;

	if (sharedSlots.includes(ItemSlot.ItemSlotRanged)) {
		if (![
			RangedWeaponType.RangedWeaponTypeBow,
			RangedWeaponType.RangedWeaponTypeCrossbow,
			RangedWeaponType.RangedWeaponTypeGun,
		].includes(item.rangedWeaponType))
			return false;
	}

	return true;
};

export function canEquipEnchant(enchant: Enchant, spec: Spec): boolean {
	const playerClass = specToClass[spec];
	if (enchant.classAllowlist.length > 0 && !enchant.classAllowlist.includes(playerClass)) {
		return false;
	}

	return true;
}

export function newUnitReference(raidIndex: number): UnitReference {
	return UnitReference.create({
		type: UnitReference_Type.Player,
		index: raidIndex,
	});
}

export function emptyUnitReference(): UnitReference {
	return UnitReference.create();
}

// Makes a new set of assignments with everything 0'd out.
export function makeBlankBlessingsAssignments(numPaladins: number): BlessingsAssignments {
	const assignments = BlessingsAssignments.create();
	for (let i = 0; i < numPaladins; i++) {
		assignments.paladins.push(BlessingsAssignment.create({
			blessings: new Array(NUM_SPECS).fill(Blessings.BlessingUnknown),
		}));
	}
	return assignments;
}

export function makeBlessingsAssignments(numPaladins: number, data: Array<{ spec: Spec, blessings: Array<Blessings> }>): BlessingsAssignments {
	const assignments = makeBlankBlessingsAssignments(numPaladins);
	for (let i = 0; i < data.length; i++) {
		const spec = data[i].spec;
		const blessings = data[i].blessings;
		for (let j = 0; j < blessings.length; j++) {
			if (j >= assignments.paladins.length) {
				// Can't assign more blessings since we ran out of paladins
				break
			}
			assignments.paladins[j].blessings[spec] = blessings[j];
		}
	}
	return assignments;
}

// Default blessings settings in the raid sim UI.
export function makeDefaultBlessings(numPaladins: number): BlessingsAssignments {
	return makeBlessingsAssignments(numPaladins, [
		{ spec: Spec.SpecBalanceDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecFeralDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecFeralTankDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSanctuary] },
		{ spec: Spec.SpecRestorationDruid, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecHunter, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecMage, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecHolyPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecProtectionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfSanctuary, Blessings.BlessingOfWisdom, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecRetributionPaladin, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecHealingPriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecShadowPriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecSmitePriest, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecRogue, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecElementalShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecEnhancementShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecRestorationShaman, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecWarlock, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfWisdom] },
		{ spec: Spec.SpecWarrior, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
		{ spec: Spec.SpecProtectionWarrior, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSanctuary] },
		{ spec: Spec.SpecDeathknight, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight, Blessings.BlessingOfSalvation] },
		{ spec: Spec.SpecTankDeathknight, blessings: [Blessings.BlessingOfKings, Blessings.BlessingOfMight] },
	]);
};

export const orderedResourceTypes: Array<ResourceType> = [
	ResourceType.ResourceTypeHealth,
	ResourceType.ResourceTypeMana,
	ResourceType.ResourceTypeEnergy,
	ResourceType.ResourceTypeRage,
	ResourceType.ResourceTypeComboPoints,
	ResourceType.ResourceTypeFocus,
	ResourceType.ResourceTypeRunicPower,
	ResourceType.ResourceTypeBloodRune,
	ResourceType.ResourceTypeFrostRune,
	ResourceType.ResourceTypeUnholyRune,
	ResourceType.ResourceTypeDeathRune,
];

export const AL_CATEGORY_HARD_MODE = 'Hard Mode';
export const AL_CATEGORY_TITAN_RUNE = 'Titan Rune';
