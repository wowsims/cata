import { CURRENT_API_VERSION, REPO_NAME } from '../constants/other.js';
import { PlayerClass } from '../player_class.js';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec.js';
import { PlayerSpecs } from '../player_specs';
import { Player } from '../proto/api.js';
import {
	Class,
	EnchantType,
	Faction,
	HandType,
	ItemSlot,
	ItemType,
	Profession,
	Race,
	RangedWeaponType,
	Spec,
	UnitReference,
	UnitReference_Type,
	WeaponType,
} from '../proto/common.js';
import {
	BloodDeathKnight,
	BloodDeathKnight_Options,
	BloodDeathKnight_Rotation,
	DeathKnightOptions,
	DeathKnightTalents,
	FrostDeathKnight,
	FrostDeathKnight_Options,
	FrostDeathKnight_Rotation,
	UnholyDeathKnight,
	UnholyDeathKnight_Options,
	UnholyDeathKnight_Rotation,
} from '../proto/death_knight.js';
import {
	BalanceDruid,
	BalanceDruid_Options,
	BalanceDruid_Rotation,
	DruidOptions,
	DruidTalents,
	FeralDruid,
	FeralDruid_Options,
	FeralDruid_Rotation,
	GuardianDruid,
	GuardianDruid_Options,
	GuardianDruid_Rotation,
	RestorationDruid,
	RestorationDruid_Options,
	RestorationDruid_Rotation,
} from '../proto/druid.js';
import {
	BeastMasteryHunter,
	BeastMasteryHunter_Options,
	BeastMasteryHunter_Rotation,
	HunterOptions,
	HunterTalents,
	MarksmanshipHunter,
	MarksmanshipHunter_Options,
	MarksmanshipHunter_Rotation,
	SurvivalHunter,
	SurvivalHunter_Options,
	SurvivalHunter_Rotation,
} from '../proto/hunter.js';
import {
	ArcaneMage,
	ArcaneMage_Options,
	ArcaneMage_Rotation,
	FireMage,
	FireMage_Options,
	FireMage_Rotation,
	FrostMage,
	FrostMage_Options,
	FrostMage_Rotation,
	MageOptions,
	MageTalents,
} from '../proto/mage.js';
import {
	BrewmasterMonk,
	BrewmasterMonk_Options,
	BrewmasterMonk_Rotation,
	MistweaverMonk,
	MistweaverMonk_Options,
	MistweaverMonk_Rotation,
	MonkOptions,
	MonkTalents,
	WindwalkerMonk,
	WindwalkerMonk_Options,
	WindwalkerMonk_Rotation,
} from '../proto/monk.js';
import {
	Blessings,
	HolyPaladin,
	HolyPaladin_Options,
	HolyPaladin_Rotation,
	PaladinOptions,
	PaladinTalents,
	ProtectionPaladin,
	ProtectionPaladin_Options,
	ProtectionPaladin_Rotation,
	RetributionPaladin,
	RetributionPaladin_Options,
	RetributionPaladin_Rotation,
} from '../proto/paladin.js';
import {
	DisciplinePriest,
	DisciplinePriest_Options,
	DisciplinePriest_Rotation,
	HolyPriest,
	HolyPriest_Options,
	HolyPriest_Rotation,
	PriestOptions,
	PriestTalents,
	ShadowPriest,
	ShadowPriest_Options,
	ShadowPriest_Rotation,
} from '../proto/priest.js';
import {
	AssassinationRogue,
	AssassinationRogue_Options,
	AssassinationRogue_Rotation,
	CombatRogue,
	CombatRogue_Options,
	CombatRogue_Rotation,
	RogueOptions,
	RogueTalents,
	SubtletyRogue,
	SubtletyRogue_Options,
	SubtletyRogue_Rotation,
} from '../proto/rogue.js';
import {
	ElementalShaman,
	ElementalShaman_Options,
	ElementalShaman_Rotation,
	EnhancementShaman,
	EnhancementShaman_Options,
	EnhancementShaman_Rotation,
	RestorationShaman,
	RestorationShaman_Options,
	RestorationShaman_Rotation,
	ShamanOptions,
	ShamanTalents,
} from '../proto/shaman.js';
import { ResourceType } from '../proto/spell';
import { BlessingsAssignment, BlessingsAssignments, UIEnchant as Enchant, UIGem as Gem, UIItem as Item } from '../proto/ui.js';
import {
	AfflictionWarlock,
	AfflictionWarlock_Options,
	AfflictionWarlock_Rotation,
	DemonologyWarlock,
	DemonologyWarlock_Options,
	DemonologyWarlock_Rotation,
	DestructionWarlock,
	DestructionWarlock_Options,
	DestructionWarlock_Rotation,
	WarlockOptions,
	WarlockTalents,
} from '../proto/warlock.js';
import {
	ArmsWarrior,
	ArmsWarrior_Options,
	ArmsWarrior_Rotation,
	FuryWarrior,
	FuryWarrior_Options,
	FuryWarrior_Rotation,
	ProtectionWarrior,
	ProtectionWarrior_Options,
	ProtectionWarrior_Rotation,
	WarriorOptions,
	WarriorTalents,
} from '../proto/warrior.js';
import { getEnumValues, intersection } from '../utils.js';
import { Stats } from './stats.js';

export const NUM_SPECS = getEnumValues(Spec).length;

export const raidSimIcon = '/mop/assets/img/raid_icon.png';
export const raidSimLabel = 'Full Raid Sim';

// Converts '111111' to [1, 1, 1, 1, 1, 1].
export function getTalentTreePoints(talentsString: string): Array<number> {
	const talents = talentsString.split('');
	return talents.map(Number);
}

export function getTalentPoints(talentsString: string): number {
	return getTalentTreePoints(talentsString).filter(Boolean).length;
}

// Gets the URL for the individual sim corresponding to the given spec.
export function getSpecSiteUrl(classString: string, specString: string): string {
	const specSiteUrlTemplate = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/CLASS/SPEC/`).toString();
	return specSiteUrlTemplate.replace('CLASS', classString).replace('SPEC', specString);
}
export const raidSimSiteUrl = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/raid/`).toString();

export function textCssClassForClass<ClassType extends Class>(playerClass: PlayerClass<ClassType>): string {
	return `text-${PlayerClasses.getCssClass(playerClass)}`;
}
export function textCssClassForSpec<SpecType extends Spec>(playerSpec: PlayerSpec<SpecType>): string {
	return textCssClassForClass(PlayerSpecs.getPlayerClass(playerSpec));
}

// Placeholder classes to fill the Unknown Spec Type Functions entry below
type UnknownSpecs = Spec.SpecUnknown;
class UnknownRotation {
	// eslint-disable-next-line @typescript-eslint/no-empty-function
	constructor() {}
}
class UnknownTalents {
	// eslint-disable-next-line @typescript-eslint/no-empty-function
	constructor() {}
}
class UnknownClassOptions {
	// eslint-disable-next-line @typescript-eslint/no-empty-function
	constructor() {}
}
class UnknownSpecOptions {
	classOptions: UnknownClassOptions;
	// eslint-disable-next-line @typescript-eslint/no-empty-function
	constructor() {
		this.classOptions = new UnknownClassOptions();
	}
}

export type DeathKnightSpecs = Spec.SpecBloodDeathKnight | Spec.SpecFrostDeathKnight | Spec.SpecUnholyDeathKnight;
export type DruidSpecs = Spec.SpecBalanceDruid | Spec.SpecFeralDruid | Spec.SpecGuardianDruid | Spec.SpecRestorationDruid;
export type HunterSpecs = Spec.SpecBeastMasteryHunter | Spec.SpecMarksmanshipHunter | Spec.SpecSurvivalHunter;
export type MageSpecs = Spec.SpecArcaneMage | Spec.SpecFireMage | Spec.SpecFrostMage;
export type PaladinSpecs = Spec.SpecHolyPaladin | Spec.SpecRetributionPaladin | Spec.SpecProtectionPaladin;
export type PriestSpecs = Spec.SpecDisciplinePriest | Spec.SpecHolyPriest | Spec.SpecShadowPriest;
export type RogueSpecs = Spec.SpecAssassinationRogue | Spec.SpecCombatRogue | Spec.SpecSubtletyRogue;
export type ShamanSpecs = Spec.SpecElementalShaman | Spec.SpecEnhancementShaman | Spec.SpecRestorationShaman;
export type WarlockSpecs = Spec.SpecAfflictionWarlock | Spec.SpecDemonologyWarlock | Spec.SpecDestructionWarlock;
export type WarriorSpecs = Spec.SpecArmsWarrior | Spec.SpecFuryWarrior | Spec.SpecProtectionWarrior;
export type MonkSpecs = Spec.SpecBrewmasterMonk | Spec.SpecMistweaverMonk | Spec.SpecWindwalkerMonk;

export type ClassSpecs<T extends Class> = T extends Class.ClassDeathKnight
	? DeathKnightSpecs
	: T extends Class.ClassDruid
	? DruidSpecs
	: T extends Class.ClassHunter
	? HunterSpecs
	: T extends Class.ClassMage
	? MageSpecs
	: T extends Class.ClassMonk
	? MonkSpecs
	: T extends Class.ClassPaladin
	? PaladinSpecs
	: T extends Class.ClassPriest
	? PriestSpecs
	: T extends Class.ClassRogue
	? RogueSpecs
	: T extends Class.ClassShaman
	? ShamanSpecs
	: T extends Class.ClassWarlock
	? WarlockSpecs
	: T extends Class.ClassWarrior
	? WarriorSpecs
	: // Should never reach this case
	  UnknownSpecs;

export type SpecClasses<T extends Spec> = T extends DeathKnightSpecs
	? Class.ClassDeathKnight
	: // Druid
	T extends DruidSpecs
	? Class.ClassDruid
	: // Hunter
	T extends HunterSpecs
	? Class.ClassHunter
	: // Mage
	T extends MageSpecs
	? Class.ClassMage
	: // Monk
	T extends MonkSpecs
	? Class.ClassMonk
	: // Paladin
	T extends PaladinSpecs
	? Class.ClassPaladin
	: // Priest
	T extends PriestSpecs
	? Class.ClassPriest
	: // Rogue
	T extends RogueSpecs
	? Class.ClassRogue
	: // Shaman
	T extends ShamanSpecs
	? Class.ClassShaman
	: // Warlock
	T extends WarlockSpecs
	? Class.ClassWarlock
	: // Warrior
	T extends WarriorSpecs
	? Class.ClassWarrior
	: // Should never reach this case
	  Class.ClassUnknown;

export type SpecRotation<T extends Spec> =
	// Death Knight
	T extends Spec.SpecBloodDeathKnight
		? BloodDeathKnight_Rotation
		: T extends Spec.SpecFrostDeathKnight
		? FrostDeathKnight_Rotation
		: T extends Spec.SpecUnholyDeathKnight
		? UnholyDeathKnight_Rotation
		: // Druid
		T extends Spec.SpecBalanceDruid
		? BalanceDruid_Rotation
		: T extends Spec.SpecFeralDruid
		? FeralDruid_Rotation
		: T extends Spec.SpecGuardianDruid
		? GuardianDruid_Rotation
		: T extends Spec.SpecRestorationDruid
		? RestorationDruid_Rotation
		: // Hunter
		T extends Spec.SpecBeastMasteryHunter
		? BeastMasteryHunter_Rotation
		: T extends Spec.SpecMarksmanshipHunter
		? MarksmanshipHunter_Rotation
		: T extends Spec.SpecSurvivalHunter
		? SurvivalHunter_Rotation
		: // Mage
		T extends Spec.SpecArcaneMage
		? ArcaneMage_Rotation
		: T extends Spec.SpecFireMage
		? FireMage_Rotation
		: T extends Spec.SpecFrostMage
		? FrostMage_Rotation
		: // Monk
		T extends Spec.SpecBrewmasterMonk
		? BrewmasterMonk_Rotation
		: T extends Spec.SpecMistweaverMonk
		? MistweaverMonk_Rotation
		: T extends Spec.SpecWindwalkerMonk
		? WindwalkerMonk_Rotation
		: // Paladin
		T extends Spec.SpecHolyPaladin
		? HolyPaladin_Rotation
		: T extends Spec.SpecProtectionPaladin
		? ProtectionPaladin_Rotation
		: T extends Spec.SpecRetributionPaladin
		? RetributionPaladin_Rotation
		: // Priest
		T extends Spec.SpecDisciplinePriest
		? DisciplinePriest_Rotation
		: T extends Spec.SpecHolyPriest
		? HolyPriest_Rotation
		: T extends Spec.SpecShadowPriest
		? ShadowPriest_Rotation
		: // Rogue
		T extends Spec.SpecAssassinationRogue
		? AssassinationRogue_Rotation
		: T extends Spec.SpecCombatRogue
		? CombatRogue_Rotation
		: T extends Spec.SpecSubtletyRogue
		? SubtletyRogue_Rotation
		: // Shaman
		T extends Spec.SpecElementalShaman
		? ElementalShaman_Rotation
		: T extends Spec.SpecEnhancementShaman
		? EnhancementShaman_Rotation
		: T extends Spec.SpecRestorationShaman
		? RestorationShaman_Rotation
		: // Warlock
		T extends Spec.SpecAfflictionWarlock
		? AfflictionWarlock_Rotation
		: T extends Spec.SpecDemonologyWarlock
		? DemonologyWarlock_Rotation
		: T extends Spec.SpecDestructionWarlock
		? DestructionWarlock_Rotation
		: // Warrior
		T extends Spec.SpecArmsWarrior
		? ArmsWarrior_Rotation
		: T extends Spec.SpecFuryWarrior
		? FuryWarrior_Rotation
		: T extends Spec.SpecProtectionWarrior
		? ProtectionWarrior_Rotation
		: // Should never reach this case
		  UnknownRotation;

export type SpecTalents<T extends Spec> =
	// Death Knight
	T extends DeathKnightSpecs
		? DeathKnightTalents
		: // Druid
		T extends DruidSpecs
		? DruidTalents
		: // Hunter
		T extends HunterSpecs
		? HunterTalents
		: // Mage
		T extends MageSpecs
		? MageTalents
		: // Monk
		T extends MonkSpecs
		? MonkTalents
		: // Paladin
		T extends PaladinSpecs
		? PaladinTalents
		: // Priest
		T extends PriestSpecs
		? PriestTalents
		: // Rogue
		T extends RogueSpecs
		? RogueTalents
		: // Shaman
		T extends ShamanSpecs
		? ShamanTalents
		: // Warlock
		T extends WarlockSpecs
		? WarlockTalents
		: // Warrior
		T extends WarriorSpecs
		? WarriorTalents
		: // Should never reach this case
		  UnknownTalents;

export type ClassOptions<T extends Spec> =
	// Death Knight
	T extends DeathKnightSpecs
		? DeathKnightOptions
		: // Druid
		T extends DruidSpecs
		? DruidOptions
		: // Hunter
		T extends HunterSpecs
		? HunterOptions
		: // Mage
		T extends MageSpecs
		? MageOptions
		: // Monk
		T extends MonkSpecs
		? MonkOptions
		: // Paladin
		T extends PaladinSpecs
		? PaladinOptions
		: // Priest
		T extends PriestSpecs
		? PriestOptions
		: // Rogue
		T extends RogueSpecs
		? RogueOptions
		: // Shaman
		T extends ShamanSpecs
		? ShamanOptions
		: // Warlock
		T extends WarlockSpecs
		? WarlockOptions
		: // Warrior
		T extends WarriorSpecs
		? WarriorOptions
		: // Should never reach this case
		  UnknownClassOptions;

export type SpecOptions<T extends Spec> =
	// Death Knight
	T extends Spec.SpecBloodDeathKnight
		? BloodDeathKnight_Options
		: T extends Spec.SpecFrostDeathKnight
		? FrostDeathKnight_Options
		: T extends Spec.SpecUnholyDeathKnight
		? UnholyDeathKnight_Options
		: // Druid
		T extends Spec.SpecBalanceDruid
		? BalanceDruid_Options
		: T extends Spec.SpecFeralDruid
		? FeralDruid_Options
		: T extends Spec.SpecGuardianDruid
		? GuardianDruid_Options
		: T extends Spec.SpecRestorationDruid
		? RestorationDruid_Options
		: // Hunter
		T extends Spec.SpecBeastMasteryHunter
		? BeastMasteryHunter_Options
		: T extends Spec.SpecMarksmanshipHunter
		? MarksmanshipHunter_Options
		: T extends Spec.SpecSurvivalHunter
		? SurvivalHunter_Options
		: // Mage
		T extends Spec.SpecArcaneMage
		? ArcaneMage_Options
		: T extends Spec.SpecFireMage
		? FireMage_Options
		: T extends Spec.SpecFrostMage
		? FrostMage_Options
		: // Monk
		T extends Spec.SpecBrewmasterMonk
		? BrewmasterMonk_Options
		: T extends Spec.SpecMistweaverMonk
		? MistweaverMonk_Options
		: T extends Spec.SpecWindwalkerMonk
		? WindwalkerMonk_Options
		: // Paladin
		T extends Spec.SpecHolyPaladin
		? HolyPaladin_Options
		: T extends Spec.SpecProtectionPaladin
		? ProtectionPaladin_Options
		: T extends Spec.SpecRetributionPaladin
		? RetributionPaladin_Options
		: // Priest
		T extends Spec.SpecDisciplinePriest
		? DisciplinePriest_Options
		: T extends Spec.SpecHolyPriest
		? HolyPriest_Options
		: T extends Spec.SpecShadowPriest
		? ShadowPriest_Options
		: // Rogue
		T extends Spec.SpecAssassinationRogue
		? AssassinationRogue_Options
		: T extends Spec.SpecCombatRogue
		? CombatRogue_Options
		: T extends Spec.SpecSubtletyRogue
		? SubtletyRogue_Options
		: // Shaman
		T extends Spec.SpecElementalShaman
		? ElementalShaman_Options
		: T extends Spec.SpecEnhancementShaman
		? EnhancementShaman_Options
		: T extends Spec.SpecRestorationShaman
		? RestorationShaman_Options
		: // Warlock
		T extends Spec.SpecAfflictionWarlock
		? AfflictionWarlock_Options
		: T extends Spec.SpecDemonologyWarlock
		? DemonologyWarlock_Options
		: T extends Spec.SpecDestructionWarlock
		? DestructionWarlock_Options
		: // Warrior
		T extends Spec.SpecArmsWarrior
		? ArmsWarrior_Options
		: T extends Spec.SpecFuryWarrior
		? FuryWarrior_Options
		: T extends Spec.SpecProtectionWarrior
		? ProtectionWarrior_Options
		: // Should never reach this case
		  UnknownSpecOptions;

export type SpecType<T extends Spec> =
	// Death Knight
	T extends Spec.SpecBloodDeathKnight
		? BloodDeathKnight
		: T extends Spec.SpecFrostDeathKnight
		? FrostDeathKnight
		: T extends Spec.SpecUnholyDeathKnight
		? UnholyDeathKnight
		: // Druid
		T extends Spec.SpecBalanceDruid
		? BalanceDruid
		: T extends Spec.SpecFeralDruid
		? FeralDruid
		: T extends Spec.SpecGuardianDruid
		? GuardianDruid
		: T extends Spec.SpecRestorationDruid
		? RestorationDruid
		: // Hunter
		T extends Spec.SpecBeastMasteryHunter
		? BeastMasteryHunter
		: T extends Spec.SpecMarksmanshipHunter
		? MarksmanshipHunter
		: T extends Spec.SpecSurvivalHunter
		? SurvivalHunter
		: // Mage
		T extends Spec.SpecArcaneMage
		? ArcaneMage
		: T extends Spec.SpecFireMage
		? FireMage
		: T extends Spec.SpecFrostMage
		? FrostMage
		: // Monk
		T extends Spec.SpecBrewmasterMonk
		? BrewmasterMonk
		: T extends Spec.SpecMistweaverMonk
		? MistweaverMonk
		: T extends Spec.SpecWindwalkerMonk
		? WindwalkerMonk
		: // Paladin
		T extends Spec.SpecHolyPaladin
		? HolyPaladin
		: T extends Spec.SpecProtectionPaladin
		? ProtectionPaladin
		: T extends Spec.SpecRetributionPaladin
		? RetributionPaladin
		: // Priest
		T extends Spec.SpecDisciplinePriest
		? DisciplinePriest
		: T extends Spec.SpecHolyPriest
		? HolyPriest
		: T extends Spec.SpecShadowPriest
		? ShadowPriest
		: // Rogue
		T extends Spec.SpecAssassinationRogue
		? AssassinationRogue
		: T extends Spec.SpecCombatRogue
		? CombatRogue
		: T extends Spec.SpecSubtletyRogue
		? SubtletyRogue
		: // Shaman
		T extends Spec.SpecElementalShaman
		? ElementalShaman
		: T extends Spec.SpecEnhancementShaman
		? EnhancementShaman
		: T extends Spec.SpecRestorationShaman
		? RestorationShaman
		: // Warlock
		T extends Spec.SpecAfflictionWarlock
		? AfflictionWarlock
		: T extends Spec.SpecDemonologyWarlock
		? DemonologyWarlock
		: T extends Spec.SpecDestructionWarlock
		? DestructionWarlock
		: // Warrior
		T extends Spec.SpecArmsWarrior
		? ArmsWarrior
		: T extends Spec.SpecFuryWarrior
		? FuryWarrior
		: T extends Spec.SpecProtectionWarrior
		? ProtectionWarrior
		: // Should never reach this case
		  Spec.SpecUnknown;

export type SpecTypeFunctions<SpecType extends Spec> = {
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

export const specTypeFunctions: Record<Spec, SpecTypeFunctions<any>> = {
	[Spec.SpecUnknown]: {
		rotationCreate: () => new UnknownRotation(),
		rotationEquals: (_a, _b) => true,
		rotationCopy: _a => new UnknownRotation(),
		rotationToJson: _a => undefined,
		rotationFromJson: _obj => new UnknownRotation(),

		talentsCreate: () => new UnknownTalents(),
		talentsEquals: (_a, _b) => true,
		talentsCopy: _a => new UnknownTalents(),
		talentsToJson: _a => undefined,
		talentsFromJson: _obj => new UnknownTalents(),

		optionsCreate: () => new UnknownSpecOptions(),
		optionsEquals: (_a, _b) => true,
		optionsCopy: _a => new UnknownSpecOptions(),
		optionsToJson: _a => undefined,
		optionsFromJson: _obj => new UnknownSpecOptions(),
		optionsFromPlayer: _player => new UnknownSpecOptions(),
	},

	// Death Knight
	[Spec.SpecBloodDeathKnight]: {
		rotationCreate: () => BloodDeathKnight_Rotation.create(),
		rotationEquals: (a, b) => BloodDeathKnight_Rotation.equals(a as BloodDeathKnight_Rotation, b as BloodDeathKnight_Rotation),
		rotationCopy: a => BloodDeathKnight_Rotation.clone(a as BloodDeathKnight_Rotation),
		rotationToJson: a => BloodDeathKnight_Rotation.toJson(a as BloodDeathKnight_Rotation),
		rotationFromJson: obj => BloodDeathKnight_Rotation.fromJson(obj),

		talentsCreate: () => DeathKnightTalents.create(),
		talentsEquals: (a, b) => DeathKnightTalents.equals(a as DeathKnightTalents, b as DeathKnightTalents),
		talentsCopy: a => DeathKnightTalents.clone(a as DeathKnightTalents),
		talentsToJson: a => DeathKnightTalents.toJson(a as DeathKnightTalents),
		talentsFromJson: obj => DeathKnightTalents.fromJson(obj),

		optionsCreate: () => BloodDeathKnight_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => BloodDeathKnight_Options.equals(a as BloodDeathKnight_Options, b as BloodDeathKnight_Options),
		optionsCopy: a => BloodDeathKnight_Options.clone(a as BloodDeathKnight_Options),
		optionsToJson: a => BloodDeathKnight_Options.toJson(a as BloodDeathKnight_Options),
		optionsFromJson: obj => BloodDeathKnight_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'bloodDeathKnight'
				? player.spec.bloodDeathKnight.options || BloodDeathKnight_Options.create()
				: BloodDeathKnight_Options.create({ classOptions: {} }),
	},
	[Spec.SpecFrostDeathKnight]: {
		rotationCreate: () => FrostDeathKnight_Rotation.create(),
		rotationEquals: (a, b) => FrostDeathKnight_Rotation.equals(a as FrostDeathKnight_Rotation, b as FrostDeathKnight_Rotation),
		rotationCopy: a => FrostDeathKnight_Rotation.clone(a as FrostDeathKnight_Rotation),
		rotationToJson: a => FrostDeathKnight_Rotation.toJson(a as FrostDeathKnight_Rotation),
		rotationFromJson: obj => FrostDeathKnight_Rotation.fromJson(obj),

		talentsCreate: () => DeathKnightTalents.create(),
		talentsEquals: (a, b) => DeathKnightTalents.equals(a as DeathKnightTalents, b as DeathKnightTalents),
		talentsCopy: a => DeathKnightTalents.clone(a as DeathKnightTalents),
		talentsToJson: a => DeathKnightTalents.toJson(a as DeathKnightTalents),
		talentsFromJson: obj => DeathKnightTalents.fromJson(obj),

		optionsCreate: () => FrostDeathKnight_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => FrostDeathKnight_Options.equals(a as FrostDeathKnight_Options, b as FrostDeathKnight_Options),
		optionsCopy: a => FrostDeathKnight_Options.clone(a as FrostDeathKnight_Options),
		optionsToJson: a => FrostDeathKnight_Options.toJson(a as FrostDeathKnight_Options),
		optionsFromJson: obj => FrostDeathKnight_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'frostDeathKnight'
				? player.spec.frostDeathKnight.options || FrostDeathKnight_Options.create()
				: FrostDeathKnight_Options.create({ classOptions: {} }),
	},
	[Spec.SpecUnholyDeathKnight]: {
		rotationCreate: () => UnholyDeathKnight_Rotation.create(),
		rotationEquals: (a, b) => UnholyDeathKnight_Rotation.equals(a as UnholyDeathKnight_Rotation, b as UnholyDeathKnight_Rotation),
		rotationCopy: a => UnholyDeathKnight_Rotation.clone(a as UnholyDeathKnight_Rotation),
		rotationToJson: a => UnholyDeathKnight_Rotation.toJson(a as UnholyDeathKnight_Rotation),
		rotationFromJson: obj => UnholyDeathKnight_Rotation.fromJson(obj),

		talentsCreate: () => DeathKnightTalents.create(),
		talentsEquals: (a, b) => DeathKnightTalents.equals(a as DeathKnightTalents, b as DeathKnightTalents),
		talentsCopy: a => DeathKnightTalents.clone(a as DeathKnightTalents),
		talentsToJson: a => DeathKnightTalents.toJson(a as DeathKnightTalents),
		talentsFromJson: obj => DeathKnightTalents.fromJson(obj),

		optionsCreate: () => UnholyDeathKnight_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => UnholyDeathKnight_Options.equals(a as UnholyDeathKnight_Options, b as UnholyDeathKnight_Options),
		optionsCopy: a => UnholyDeathKnight_Options.clone(a as UnholyDeathKnight_Options),
		optionsToJson: a => UnholyDeathKnight_Options.toJson(a as UnholyDeathKnight_Options),
		optionsFromJson: obj => UnholyDeathKnight_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'unholyDeathKnight'
				? player.spec.unholyDeathKnight.options || UnholyDeathKnight_Options.create()
				: UnholyDeathKnight_Options.create({ classOptions: {} }),
	},
	// Druid
	[Spec.SpecBalanceDruid]: {
		rotationCreate: () => BalanceDruid_Rotation.create(),
		rotationEquals: (a, b) => BalanceDruid_Rotation.equals(a as BalanceDruid_Rotation, b as BalanceDruid_Rotation),
		rotationCopy: a => BalanceDruid_Rotation.clone(a as BalanceDruid_Rotation),
		rotationToJson: a => BalanceDruid_Rotation.toJson(a as BalanceDruid_Rotation),
		rotationFromJson: obj => BalanceDruid_Rotation.fromJson(obj),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: a => DruidTalents.clone(a as DruidTalents),
		talentsToJson: a => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: obj => DruidTalents.fromJson(obj),

		optionsCreate: () => BalanceDruid_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => BalanceDruid_Options.equals(a as BalanceDruid_Options, b as BalanceDruid_Options),
		optionsCopy: a => BalanceDruid_Options.clone(a as BalanceDruid_Options),
		optionsToJson: a => BalanceDruid_Options.toJson(a as BalanceDruid_Options),
		optionsFromJson: obj => BalanceDruid_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'balanceDruid'
				? player.spec.balanceDruid.options || BalanceDruid_Options.create()
				: BalanceDruid_Options.create({ classOptions: {} }),
	},
	[Spec.SpecFeralDruid]: {
		rotationCreate: () => FeralDruid_Rotation.create(),
		rotationEquals: (a, b) => FeralDruid_Rotation.equals(a as FeralDruid_Rotation, b as FeralDruid_Rotation),
		rotationCopy: a => FeralDruid_Rotation.clone(a as FeralDruid_Rotation),
		rotationToJson: a => FeralDruid_Rotation.toJson(a as FeralDruid_Rotation),
		rotationFromJson: obj => FeralDruid_Rotation.fromJson(obj),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: a => DruidTalents.clone(a as DruidTalents),
		talentsToJson: a => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: obj => DruidTalents.fromJson(obj),

		optionsCreate: () => FeralDruid_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => FeralDruid_Options.equals(a as FeralDruid_Options, b as FeralDruid_Options),
		optionsCopy: a => FeralDruid_Options.clone(a as FeralDruid_Options),
		optionsToJson: a => FeralDruid_Options.toJson(a as FeralDruid_Options),
		optionsFromJson: obj => FeralDruid_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'feralDruid'
				? player.spec.feralDruid.options || FeralDruid_Options.create()
				: FeralDruid_Options.create({ classOptions: {} }),
	},
	[Spec.SpecGuardianDruid]: {
		rotationCreate: () => GuardianDruid_Rotation.create(),
		rotationEquals: (a, b) => GuardianDruid_Rotation.equals(a as GuardianDruid_Rotation, b as GuardianDruid_Rotation),
		rotationCopy: a => GuardianDruid_Rotation.clone(a as GuardianDruid_Rotation),
		rotationToJson: a => GuardianDruid_Rotation.toJson(a as GuardianDruid_Rotation),
		rotationFromJson: obj => GuardianDruid_Rotation.fromJson(obj),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: a => DruidTalents.clone(a as DruidTalents),
		talentsToJson: a => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: obj => DruidTalents.fromJson(obj),

		optionsCreate: () => GuardianDruid_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => GuardianDruid_Options.equals(a as GuardianDruid_Options, b as GuardianDruid_Options),
		optionsCopy: a => GuardianDruid_Options.clone(a as GuardianDruid_Options),
		optionsToJson: a => GuardianDruid_Options.toJson(a as GuardianDruid_Options),
		optionsFromJson: obj => GuardianDruid_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'guardianDruid'
				? player.spec.guardianDruid.options || GuardianDruid_Options.create()
				: GuardianDruid_Options.create({ classOptions: {} }),
	},
	[Spec.SpecRestorationDruid]: {
		rotationCreate: () => RestorationDruid_Rotation.create(),
		rotationEquals: (a, b) => RestorationDruid_Rotation.equals(a as RestorationDruid_Rotation, b as RestorationDruid_Rotation),
		rotationCopy: a => RestorationDruid_Rotation.clone(a as RestorationDruid_Rotation),
		rotationToJson: a => RestorationDruid_Rotation.toJson(a as RestorationDruid_Rotation),
		rotationFromJson: obj => RestorationDruid_Rotation.fromJson(obj),

		talentsCreate: () => DruidTalents.create(),
		talentsEquals: (a, b) => DruidTalents.equals(a as DruidTalents, b as DruidTalents),
		talentsCopy: a => DruidTalents.clone(a as DruidTalents),
		talentsToJson: a => DruidTalents.toJson(a as DruidTalents),
		talentsFromJson: obj => DruidTalents.fromJson(obj),

		optionsCreate: () => RestorationDruid_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => RestorationDruid_Options.equals(a as RestorationDruid_Options, b as RestorationDruid_Options),
		optionsCopy: a => RestorationDruid_Options.clone(a as RestorationDruid_Options),
		optionsToJson: a => RestorationDruid_Options.toJson(a as RestorationDruid_Options),
		optionsFromJson: obj => RestorationDruid_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'restorationDruid'
				? player.spec.restorationDruid.options || RestorationDruid_Options.create()
				: RestorationDruid_Options.create({ classOptions: {} }),
	},
	// Hunter
	[Spec.SpecBeastMasteryHunter]: {
		rotationCreate: () => BeastMasteryHunter_Rotation.create(),
		rotationEquals: (a, b) => BeastMasteryHunter_Rotation.equals(a as BeastMasteryHunter_Rotation, b as BeastMasteryHunter_Rotation),
		rotationCopy: a => BeastMasteryHunter_Rotation.clone(a as BeastMasteryHunter_Rotation),
		rotationToJson: a => BeastMasteryHunter_Rotation.toJson(a as BeastMasteryHunter_Rotation),
		rotationFromJson: obj => BeastMasteryHunter_Rotation.fromJson(obj),

		talentsCreate: () => HunterTalents.create(),
		talentsEquals: (a, b) => HunterTalents.equals(a as HunterTalents, b as HunterTalents),
		talentsCopy: a => HunterTalents.clone(a as HunterTalents),
		talentsToJson: a => HunterTalents.toJson(a as HunterTalents),
		talentsFromJson: obj => HunterTalents.fromJson(obj),

		optionsCreate: () => BeastMasteryHunter_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => BeastMasteryHunter_Options.equals(a as BeastMasteryHunter_Options, b as BeastMasteryHunter_Options),
		optionsCopy: a => BeastMasteryHunter_Options.clone(a as BeastMasteryHunter_Options),
		optionsToJson: a => BeastMasteryHunter_Options.toJson(a as BeastMasteryHunter_Options),
		optionsFromJson: obj => BeastMasteryHunter_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'beastMasteryHunter'
				? player.spec.beastMasteryHunter.options || BeastMasteryHunter_Options.create()
				: BeastMasteryHunter_Options.create({ classOptions: {} }),
	},
	[Spec.SpecMarksmanshipHunter]: {
		rotationCreate: () => MarksmanshipHunter_Rotation.create(),
		rotationEquals: (a, b) => MarksmanshipHunter_Rotation.equals(a as MarksmanshipHunter_Rotation, b as MarksmanshipHunter_Rotation),
		rotationCopy: a => MarksmanshipHunter_Rotation.clone(a as MarksmanshipHunter_Rotation),
		rotationToJson: a => MarksmanshipHunter_Rotation.toJson(a as MarksmanshipHunter_Rotation),
		rotationFromJson: obj => MarksmanshipHunter_Rotation.fromJson(obj),

		talentsCreate: () => HunterTalents.create(),
		talentsEquals: (a, b) => HunterTalents.equals(a as HunterTalents, b as HunterTalents),
		talentsCopy: a => HunterTalents.clone(a as HunterTalents),
		talentsToJson: a => HunterTalents.toJson(a as HunterTalents),
		talentsFromJson: obj => HunterTalents.fromJson(obj),

		optionsCreate: () => MarksmanshipHunter_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => MarksmanshipHunter_Options.equals(a as MarksmanshipHunter_Options, b as MarksmanshipHunter_Options),
		optionsCopy: a => MarksmanshipHunter_Options.clone(a as MarksmanshipHunter_Options),
		optionsToJson: a => MarksmanshipHunter_Options.toJson(a as MarksmanshipHunter_Options),
		optionsFromJson: obj => MarksmanshipHunter_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'marksmanshipHunter'
				? player.spec.marksmanshipHunter.options || MarksmanshipHunter_Options.create()
				: MarksmanshipHunter_Options.create({ classOptions: {} }),
	},
	[Spec.SpecSurvivalHunter]: {
		rotationCreate: () => SurvivalHunter_Rotation.create(),
		rotationEquals: (a, b) => SurvivalHunter_Rotation.equals(a as SurvivalHunter_Rotation, b as SurvivalHunter_Rotation),
		rotationCopy: a => SurvivalHunter_Rotation.clone(a as SurvivalHunter_Rotation),
		rotationToJson: a => SurvivalHunter_Rotation.toJson(a as SurvivalHunter_Rotation),
		rotationFromJson: obj => SurvivalHunter_Rotation.fromJson(obj),

		talentsCreate: () => HunterTalents.create(),
		talentsEquals: (a, b) => HunterTalents.equals(a as HunterTalents, b as HunterTalents),
		talentsCopy: a => HunterTalents.clone(a as HunterTalents),
		talentsToJson: a => HunterTalents.toJson(a as HunterTalents),
		talentsFromJson: obj => HunterTalents.fromJson(obj),

		optionsCreate: () => SurvivalHunter_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => SurvivalHunter_Options.equals(a as SurvivalHunter_Options, b as SurvivalHunter_Options),
		optionsCopy: a => SurvivalHunter_Options.clone(a as SurvivalHunter_Options),
		optionsToJson: a => SurvivalHunter_Options.toJson(a as SurvivalHunter_Options),
		optionsFromJson: obj => SurvivalHunter_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'survivalHunter'
				? player.spec.survivalHunter.options || SurvivalHunter_Options.create()
				: SurvivalHunter_Options.create({ classOptions: {} }),
	},
	// Mage
	[Spec.SpecArcaneMage]: {
		rotationCreate: () => ArcaneMage_Rotation.create(),
		rotationEquals: (a, b) => ArcaneMage_Rotation.equals(a as ArcaneMage_Rotation, b as ArcaneMage_Rotation),
		rotationCopy: a => ArcaneMage_Rotation.clone(a as ArcaneMage_Rotation),
		rotationToJson: a => ArcaneMage_Rotation.toJson(a as ArcaneMage_Rotation),
		rotationFromJson: obj => ArcaneMage_Rotation.fromJson(obj),

		talentsCreate: () => MageTalents.create(),
		talentsEquals: (a, b) => MageTalents.equals(a as MageTalents, b as MageTalents),
		talentsCopy: a => MageTalents.clone(a as MageTalents),
		talentsToJson: a => MageTalents.toJson(a as MageTalents),
		talentsFromJson: obj => MageTalents.fromJson(obj),

		optionsCreate: () => ArcaneMage_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => ArcaneMage_Options.equals(a as ArcaneMage_Options, b as ArcaneMage_Options),
		optionsCopy: a => ArcaneMage_Options.clone(a as ArcaneMage_Options),
		optionsToJson: a => ArcaneMage_Options.toJson(a as ArcaneMage_Options),
		optionsFromJson: obj => ArcaneMage_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'arcaneMage'
				? player.spec.arcaneMage.options || ArcaneMage_Options.create()
				: ArcaneMage_Options.create({ classOptions: {} }),
	},
	[Spec.SpecFireMage]: {
		rotationCreate: () => FireMage_Rotation.create(),
		rotationEquals: (a, b) => FireMage_Rotation.equals(a as FireMage_Rotation, b as FireMage_Rotation),
		rotationCopy: a => FireMage_Rotation.clone(a as FireMage_Rotation),
		rotationToJson: a => FireMage_Rotation.toJson(a as FireMage_Rotation),
		rotationFromJson: obj => FireMage_Rotation.fromJson(obj),

		talentsCreate: () => MageTalents.create(),
		talentsEquals: (a, b) => MageTalents.equals(a as MageTalents, b as MageTalents),
		talentsCopy: a => MageTalents.clone(a as MageTalents),
		talentsToJson: a => MageTalents.toJson(a as MageTalents),
		talentsFromJson: obj => MageTalents.fromJson(obj),

		optionsCreate: () => FireMage_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => FireMage_Options.equals(a as FireMage_Options, b as FireMage_Options),
		optionsCopy: a => FireMage_Options.clone(a as FireMage_Options),
		optionsToJson: a => FireMage_Options.toJson(a as FireMage_Options),
		optionsFromJson: obj => FireMage_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'fireMage' ? player.spec.fireMage.options || FireMage_Options.create() : FireMage_Options.create({ classOptions: {} }),
	},
	[Spec.SpecFrostMage]: {
		rotationCreate: () => FrostMage_Rotation.create(),
		rotationEquals: (a, b) => FrostMage_Rotation.equals(a as FrostMage_Rotation, b as FrostMage_Rotation),
		rotationCopy: a => FrostMage_Rotation.clone(a as FrostMage_Rotation),
		rotationToJson: a => FrostMage_Rotation.toJson(a as FrostMage_Rotation),
		rotationFromJson: obj => FrostMage_Rotation.fromJson(obj),

		talentsCreate: () => MageTalents.create(),
		talentsEquals: (a, b) => MageTalents.equals(a as MageTalents, b as MageTalents),
		talentsCopy: a => MageTalents.clone(a as MageTalents),
		talentsToJson: a => MageTalents.toJson(a as MageTalents),
		talentsFromJson: obj => MageTalents.fromJson(obj),

		optionsCreate: () => FrostMage_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => FrostMage_Options.equals(a as FrostMage_Options, b as FrostMage_Options),
		optionsCopy: a => FrostMage_Options.clone(a as FrostMage_Options),
		optionsToJson: a => FrostMage_Options.toJson(a as FrostMage_Options),
		optionsFromJson: obj => FrostMage_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'frostMage' ? player.spec.frostMage.options || FrostMage_Options.create() : FrostMage_Options.create({ classOptions: {} }),
	},
	// Monk
	[Spec.SpecBrewmasterMonk]: {
		rotationCreate: () => BrewmasterMonk_Rotation.create(),
		rotationEquals: (a, b) => BrewmasterMonk_Rotation.equals(a as BrewmasterMonk_Rotation, b as BrewmasterMonk_Rotation),
		rotationCopy: a => BrewmasterMonk_Rotation.clone(a as BrewmasterMonk_Rotation),
		rotationToJson: a => BrewmasterMonk_Rotation.toJson(a as BrewmasterMonk_Rotation),
		rotationFromJson: obj => BrewmasterMonk_Rotation.fromJson(obj),

		talentsCreate: () => MonkTalents.create(),
		talentsEquals: (a, b) => MonkTalents.equals(a as MonkTalents, b as MonkTalents),
		talentsCopy: a => MonkTalents.clone(a as MonkTalents),
		talentsToJson: a => MonkTalents.toJson(a as MonkTalents),
		talentsFromJson: obj => MonkTalents.fromJson(obj),

		optionsCreate: () => BrewmasterMonk_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => BrewmasterMonk_Options.equals(a as BrewmasterMonk_Options, b as BrewmasterMonk_Options),
		optionsCopy: a => BrewmasterMonk_Options.clone(a as BrewmasterMonk_Options),
		optionsToJson: a => BrewmasterMonk_Options.toJson(a as BrewmasterMonk_Options),
		optionsFromJson: obj => BrewmasterMonk_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'brewmasterMonk'
				? player.spec.brewmasterMonk.options || BrewmasterMonk_Options.create()
				: BrewmasterMonk_Options.create({ classOptions: {} }),
	},
	[Spec.SpecMistweaverMonk]: {
		rotationCreate: () => MistweaverMonk_Rotation.create(),
		rotationEquals: (a, b) => MistweaverMonk_Rotation.equals(a as MistweaverMonk_Rotation, b as MistweaverMonk_Rotation),
		rotationCopy: a => MistweaverMonk_Rotation.clone(a as MistweaverMonk_Rotation),
		rotationToJson: a => MistweaverMonk_Rotation.toJson(a as MistweaverMonk_Rotation),
		rotationFromJson: obj => MistweaverMonk_Rotation.fromJson(obj),

		talentsCreate: () => MonkTalents.create(),
		talentsEquals: (a, b) => MonkTalents.equals(a as MonkTalents, b as MonkTalents),
		talentsCopy: a => MonkTalents.clone(a as MonkTalents),
		talentsToJson: a => MonkTalents.toJson(a as MonkTalents),
		talentsFromJson: obj => MonkTalents.fromJson(obj),

		optionsCreate: () => MistweaverMonk_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => MistweaverMonk_Options.equals(a as MistweaverMonk_Options, b as MistweaverMonk_Options),
		optionsCopy: a => MistweaverMonk_Options.clone(a as MistweaverMonk_Options),
		optionsToJson: a => MistweaverMonk_Options.toJson(a as MistweaverMonk_Options),
		optionsFromJson: obj => MistweaverMonk_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'mistweaverMonk'
				? player.spec.mistweaverMonk.options || MistweaverMonk_Options.create()
				: MistweaverMonk_Options.create({ classOptions: {} }),
	},
	[Spec.SpecWindwalkerMonk]: {
		rotationCreate: () => WindwalkerMonk_Rotation.create(),
		rotationEquals: (a, b) => WindwalkerMonk_Rotation.equals(a as WindwalkerMonk_Rotation, b as WindwalkerMonk_Rotation),
		rotationCopy: a => WindwalkerMonk_Rotation.clone(a as WindwalkerMonk_Rotation),
		rotationToJson: a => WindwalkerMonk_Rotation.toJson(a as WindwalkerMonk_Rotation),
		rotationFromJson: obj => WindwalkerMonk_Rotation.fromJson(obj),

		talentsCreate: () => MonkTalents.create(),
		talentsEquals: (a, b) => MonkTalents.equals(a as MonkTalents, b as MonkTalents),
		talentsCopy: a => MonkTalents.clone(a as MonkTalents),
		talentsToJson: a => MonkTalents.toJson(a as MonkTalents),
		talentsFromJson: obj => MonkTalents.fromJson(obj),

		optionsCreate: () => WindwalkerMonk_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => WindwalkerMonk_Options.equals(a as WindwalkerMonk_Options, b as WindwalkerMonk_Options),
		optionsCopy: a => WindwalkerMonk_Options.clone(a as WindwalkerMonk_Options),
		optionsToJson: a => WindwalkerMonk_Options.toJson(a as WindwalkerMonk_Options),
		optionsFromJson: obj => WindwalkerMonk_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'windwalkerMonk'
				? player.spec.windwalkerMonk.options || WindwalkerMonk_Options.create()
				: WindwalkerMonk_Options.create({ classOptions: {} }),
	},
	// Paladin
	[Spec.SpecHolyPaladin]: {
		rotationCreate: () => HolyPaladin_Rotation.create(),
		rotationEquals: (a, b) => HolyPaladin_Rotation.equals(a as HolyPaladin_Rotation, b as HolyPaladin_Rotation),
		rotationCopy: a => HolyPaladin_Rotation.clone(a as HolyPaladin_Rotation),
		rotationToJson: a => HolyPaladin_Rotation.toJson(a as HolyPaladin_Rotation),
		rotationFromJson: obj => HolyPaladin_Rotation.fromJson(obj),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: a => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: a => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: obj => PaladinTalents.fromJson(obj),

		optionsCreate: () => HolyPaladin_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => HolyPaladin_Options.equals(a as HolyPaladin_Options, b as HolyPaladin_Options),
		optionsCopy: a => HolyPaladin_Options.clone(a as HolyPaladin_Options),
		optionsToJson: a => HolyPaladin_Options.toJson(a as HolyPaladin_Options),
		optionsFromJson: obj => HolyPaladin_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'holyPaladin'
				? player.spec.holyPaladin.options || HolyPaladin_Options.create()
				: HolyPaladin_Options.create({ classOptions: {} }),
	},
	[Spec.SpecProtectionPaladin]: {
		rotationCreate: () => ProtectionPaladin_Rotation.create(),
		rotationEquals: (a, b) => ProtectionPaladin_Rotation.equals(a as ProtectionPaladin_Rotation, b as ProtectionPaladin_Rotation),
		rotationCopy: a => ProtectionPaladin_Rotation.clone(a as ProtectionPaladin_Rotation),
		rotationToJson: a => ProtectionPaladin_Rotation.toJson(a as ProtectionPaladin_Rotation),
		rotationFromJson: obj => ProtectionPaladin_Rotation.fromJson(obj),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: a => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: a => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: obj => PaladinTalents.fromJson(obj),

		optionsCreate: () => ProtectionPaladin_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => ProtectionPaladin_Options.equals(a as ProtectionPaladin_Options, b as ProtectionPaladin_Options),
		optionsCopy: a => ProtectionPaladin_Options.clone(a as ProtectionPaladin_Options),
		optionsToJson: a => ProtectionPaladin_Options.toJson(a as ProtectionPaladin_Options),
		optionsFromJson: obj => ProtectionPaladin_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'protectionPaladin'
				? player.spec.protectionPaladin.options || ProtectionPaladin_Options.create()
				: ProtectionPaladin_Options.create({ classOptions: {} }),
	},
	[Spec.SpecRetributionPaladin]: {
		rotationCreate: () => RetributionPaladin_Rotation.create(),
		rotationEquals: (a, b) => RetributionPaladin_Rotation.equals(a as RetributionPaladin_Rotation, b as RetributionPaladin_Rotation),
		rotationCopy: a => RetributionPaladin_Rotation.clone(a as RetributionPaladin_Rotation),
		rotationToJson: a => RetributionPaladin_Rotation.toJson(a as RetributionPaladin_Rotation),
		rotationFromJson: obj => RetributionPaladin_Rotation.fromJson(obj),

		talentsCreate: () => PaladinTalents.create(),
		talentsEquals: (a, b) => PaladinTalents.equals(a as PaladinTalents, b as PaladinTalents),
		talentsCopy: a => PaladinTalents.clone(a as PaladinTalents),
		talentsToJson: a => PaladinTalents.toJson(a as PaladinTalents),
		talentsFromJson: obj => PaladinTalents.fromJson(obj),

		optionsCreate: () => RetributionPaladin_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => RetributionPaladin_Options.equals(a as RetributionPaladin_Options, b as RetributionPaladin_Options),
		optionsCopy: a => RetributionPaladin_Options.clone(a as RetributionPaladin_Options),
		optionsToJson: a => RetributionPaladin_Options.toJson(a as RetributionPaladin_Options),
		optionsFromJson: obj => RetributionPaladin_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'retributionPaladin'
				? player.spec.retributionPaladin.options || RetributionPaladin_Options.create()
				: RetributionPaladin_Options.create({ classOptions: {} }),
	},
	// Priest
	[Spec.SpecDisciplinePriest]: {
		rotationCreate: () => DisciplinePriest_Rotation.create(),
		rotationEquals: (a, b) => DisciplinePriest_Rotation.equals(a as DisciplinePriest_Rotation, b as DisciplinePriest_Rotation),
		rotationCopy: a => DisciplinePriest_Rotation.clone(a as DisciplinePriest_Rotation),
		rotationToJson: a => DisciplinePriest_Rotation.toJson(a as DisciplinePriest_Rotation),
		rotationFromJson: obj => DisciplinePriest_Rotation.fromJson(obj),

		talentsCreate: () => PriestTalents.create(),
		talentsEquals: (a, b) => PriestTalents.equals(a as PriestTalents, b as PriestTalents),
		talentsCopy: a => PriestTalents.clone(a as PriestTalents),
		talentsToJson: a => PriestTalents.toJson(a as PriestTalents),
		talentsFromJson: obj => PriestTalents.fromJson(obj),

		optionsCreate: () => DisciplinePriest_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => DisciplinePriest_Options.equals(a as DisciplinePriest_Options, b as DisciplinePriest_Options),
		optionsCopy: a => DisciplinePriest_Options.clone(a as DisciplinePriest_Options),
		optionsToJson: a => DisciplinePriest_Options.toJson(a as DisciplinePriest_Options),
		optionsFromJson: obj => DisciplinePriest_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'disciplinePriest'
				? player.spec.disciplinePriest.options || DisciplinePriest_Options.create()
				: DisciplinePriest_Options.create({ classOptions: {} }),
	},
	[Spec.SpecHolyPriest]: {
		rotationCreate: () => HolyPriest_Rotation.create(),
		rotationEquals: (a, b) => HolyPriest_Rotation.equals(a as HolyPriest_Rotation, b as HolyPriest_Rotation),
		rotationCopy: a => HolyPriest_Rotation.clone(a as HolyPriest_Rotation),
		rotationToJson: a => HolyPriest_Rotation.toJson(a as HolyPriest_Rotation),
		rotationFromJson: obj => HolyPriest_Rotation.fromJson(obj),

		talentsCreate: () => PriestTalents.create(),
		talentsEquals: (a, b) => PriestTalents.equals(a as PriestTalents, b as PriestTalents),
		talentsCopy: a => PriestTalents.clone(a as PriestTalents),
		talentsToJson: a => PriestTalents.toJson(a as PriestTalents),
		talentsFromJson: obj => PriestTalents.fromJson(obj),

		optionsCreate: () => HolyPriest_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => HolyPriest_Options.equals(a as HolyPriest_Options, b as HolyPriest_Options),
		optionsCopy: a => HolyPriest_Options.clone(a as HolyPriest_Options),
		optionsToJson: a => HolyPriest_Options.toJson(a as HolyPriest_Options),
		optionsFromJson: obj => HolyPriest_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'holyPriest'
				? player.spec.holyPriest.options || HolyPriest_Options.create()
				: HolyPriest_Options.create({ classOptions: {} }),
	},
	[Spec.SpecShadowPriest]: {
		rotationCreate: () => ShadowPriest_Rotation.create(),
		rotationEquals: (a, b) => ShadowPriest_Rotation.equals(a as ShadowPriest_Rotation, b as ShadowPriest_Rotation),
		rotationCopy: a => ShadowPriest_Rotation.clone(a as ShadowPriest_Rotation),
		rotationToJson: a => ShadowPriest_Rotation.toJson(a as ShadowPriest_Rotation),
		rotationFromJson: obj => ShadowPriest_Rotation.fromJson(obj),

		talentsCreate: () => PriestTalents.create(),
		talentsEquals: (a, b) => PriestTalents.equals(a as PriestTalents, b as PriestTalents),
		talentsCopy: a => PriestTalents.clone(a as PriestTalents),
		talentsToJson: a => PriestTalents.toJson(a as PriestTalents),
		talentsFromJson: obj => PriestTalents.fromJson(obj),

		optionsCreate: () => ShadowPriest_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => ShadowPriest_Options.equals(a as ShadowPriest_Options, b as ShadowPriest_Options),
		optionsCopy: a => ShadowPriest_Options.clone(a as ShadowPriest_Options),
		optionsToJson: a => ShadowPriest_Options.toJson(a as ShadowPriest_Options),
		optionsFromJson: obj => ShadowPriest_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'shadowPriest'
				? player.spec.shadowPriest.options || ShadowPriest_Options.create()
				: ShadowPriest_Options.create({ classOptions: {} }),
	},
	// Rogue
	[Spec.SpecAssassinationRogue]: {
		rotationCreate: () => AssassinationRogue_Rotation.create(),
		rotationEquals: (a, b) => AssassinationRogue_Rotation.equals(a as AssassinationRogue_Rotation, b as AssassinationRogue_Rotation),
		rotationCopy: a => AssassinationRogue_Rotation.clone(a as AssassinationRogue_Rotation),
		rotationToJson: a => AssassinationRogue_Rotation.toJson(a as AssassinationRogue_Rotation),
		rotationFromJson: obj => AssassinationRogue_Rotation.fromJson(obj),

		talentsCreate: () => RogueTalents.create(),
		talentsEquals: (a, b) => RogueTalents.equals(a as RogueTalents, b as RogueTalents),
		talentsCopy: a => RogueTalents.clone(a as RogueTalents),
		talentsToJson: a => RogueTalents.toJson(a as RogueTalents),
		talentsFromJson: obj => RogueTalents.fromJson(obj),

		optionsCreate: () => AssassinationRogue_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => AssassinationRogue_Options.equals(a as AssassinationRogue_Options, b as AssassinationRogue_Options),
		optionsCopy: a => AssassinationRogue_Options.clone(a as AssassinationRogue_Options),
		optionsToJson: a => AssassinationRogue_Options.toJson(a as AssassinationRogue_Options),
		optionsFromJson: obj => AssassinationRogue_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'assassinationRogue'
				? player.spec.assassinationRogue.options || AssassinationRogue_Options.create()
				: AssassinationRogue_Options.create({ classOptions: {} }),
	},
	[Spec.SpecCombatRogue]: {
		rotationCreate: () => CombatRogue_Rotation.create(),
		rotationEquals: (a, b) => CombatRogue_Rotation.equals(a as CombatRogue_Rotation, b as CombatRogue_Rotation),
		rotationCopy: a => CombatRogue_Rotation.clone(a as CombatRogue_Rotation),
		rotationToJson: a => CombatRogue_Rotation.toJson(a as CombatRogue_Rotation),
		rotationFromJson: obj => CombatRogue_Rotation.fromJson(obj),

		talentsCreate: () => RogueTalents.create(),
		talentsEquals: (a, b) => RogueTalents.equals(a as RogueTalents, b as RogueTalents),
		talentsCopy: a => RogueTalents.clone(a as RogueTalents),
		talentsToJson: a => RogueTalents.toJson(a as RogueTalents),
		talentsFromJson: obj => RogueTalents.fromJson(obj),

		optionsCreate: () => CombatRogue_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => CombatRogue_Options.equals(a as CombatRogue_Options, b as CombatRogue_Options),
		optionsCopy: a => CombatRogue_Options.clone(a as CombatRogue_Options),
		optionsToJson: a => CombatRogue_Options.toJson(a as CombatRogue_Options),
		optionsFromJson: obj => CombatRogue_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'combatRogue'
				? player.spec.combatRogue.options || CombatRogue_Options.create()
				: CombatRogue_Options.create({ classOptions: {} }),
	},
	[Spec.SpecSubtletyRogue]: {
		rotationCreate: () => SubtletyRogue_Rotation.create(),
		rotationEquals: (a, b) => SubtletyRogue_Rotation.equals(a as SubtletyRogue_Rotation, b as SubtletyRogue_Rotation),
		rotationCopy: a => SubtletyRogue_Rotation.clone(a as SubtletyRogue_Rotation),
		rotationToJson: a => SubtletyRogue_Rotation.toJson(a as SubtletyRogue_Rotation),
		rotationFromJson: obj => SubtletyRogue_Rotation.fromJson(obj),

		talentsCreate: () => RogueTalents.create(),
		talentsEquals: (a, b) => RogueTalents.equals(a as RogueTalents, b as RogueTalents),
		talentsCopy: a => RogueTalents.clone(a as RogueTalents),
		talentsToJson: a => RogueTalents.toJson(a as RogueTalents),
		talentsFromJson: obj => RogueTalents.fromJson(obj),

		optionsCreate: () => SubtletyRogue_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => SubtletyRogue_Options.equals(a as SubtletyRogue_Options, b as SubtletyRogue_Options),
		optionsCopy: a => SubtletyRogue_Options.clone(a as SubtletyRogue_Options),
		optionsToJson: a => SubtletyRogue_Options.toJson(a as SubtletyRogue_Options),
		optionsFromJson: obj => SubtletyRogue_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'subtletyRogue'
				? player.spec.subtletyRogue.options || SubtletyRogue_Options.create()
				: SubtletyRogue_Options.create({ classOptions: {} }),
	},
	// Shaman
	[Spec.SpecElementalShaman]: {
		rotationCreate: () => ElementalShaman_Rotation.create(),
		rotationEquals: (a, b) => ElementalShaman_Rotation.equals(a as ElementalShaman_Rotation, b as ElementalShaman_Rotation),
		rotationCopy: a => ElementalShaman_Rotation.clone(a as ElementalShaman_Rotation),
		rotationToJson: a => ElementalShaman_Rotation.toJson(a as ElementalShaman_Rotation),
		rotationFromJson: obj => ElementalShaman_Rotation.fromJson(obj),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: a => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: a => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: obj => ShamanTalents.fromJson(obj),

		optionsCreate: () => ElementalShaman_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => ElementalShaman_Options.equals(a as ElementalShaman_Options, b as ElementalShaman_Options),
		optionsCopy: a => ElementalShaman_Options.clone(a as ElementalShaman_Options),
		optionsToJson: a => ElementalShaman_Options.toJson(a as ElementalShaman_Options),
		optionsFromJson: obj => ElementalShaman_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'elementalShaman'
				? player.spec.elementalShaman.options || ElementalShaman_Options.create()
				: ElementalShaman_Options.create({ classOptions: {} }),
	},
	[Spec.SpecEnhancementShaman]: {
		rotationCreate: () => EnhancementShaman_Rotation.create(),
		rotationEquals: (a, b) => EnhancementShaman_Rotation.equals(a as EnhancementShaman_Rotation, b as EnhancementShaman_Rotation),
		rotationCopy: a => EnhancementShaman_Rotation.clone(a as EnhancementShaman_Rotation),
		rotationToJson: a => EnhancementShaman_Rotation.toJson(a as EnhancementShaman_Rotation),
		rotationFromJson: obj => EnhancementShaman_Rotation.fromJson(obj),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: a => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: a => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: obj => ShamanTalents.fromJson(obj),

		optionsCreate: () => EnhancementShaman_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => EnhancementShaman_Options.equals(a as EnhancementShaman_Options, b as EnhancementShaman_Options),
		optionsCopy: a => EnhancementShaman_Options.clone(a as EnhancementShaman_Options),
		optionsToJson: a => EnhancementShaman_Options.toJson(a as EnhancementShaman_Options),
		optionsFromJson: obj => EnhancementShaman_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'enhancementShaman'
				? player.spec.enhancementShaman.options || EnhancementShaman_Options.create()
				: EnhancementShaman_Options.create({ classOptions: {} }),
	},
	[Spec.SpecRestorationShaman]: {
		rotationCreate: () => RestorationShaman_Rotation.create(),
		rotationEquals: (a, b) => RestorationShaman_Rotation.equals(a as RestorationShaman_Rotation, b as RestorationShaman_Rotation),
		rotationCopy: a => RestorationShaman_Rotation.clone(a as RestorationShaman_Rotation),
		rotationToJson: a => RestorationShaman_Rotation.toJson(a as RestorationShaman_Rotation),
		rotationFromJson: obj => RestorationShaman_Rotation.fromJson(obj),

		talentsCreate: () => ShamanTalents.create(),
		talentsEquals: (a, b) => ShamanTalents.equals(a as ShamanTalents, b as ShamanTalents),
		talentsCopy: a => ShamanTalents.clone(a as ShamanTalents),
		talentsToJson: a => ShamanTalents.toJson(a as ShamanTalents),
		talentsFromJson: obj => ShamanTalents.fromJson(obj),

		optionsCreate: () => RestorationShaman_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => RestorationShaman_Options.equals(a as RestorationShaman_Options, b as RestorationShaman_Options),
		optionsCopy: a => RestorationShaman_Options.clone(a as RestorationShaman_Options),
		optionsToJson: a => RestorationShaman_Options.toJson(a as RestorationShaman_Options),
		optionsFromJson: obj => RestorationShaman_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'restorationShaman'
				? player.spec.restorationShaman.options || RestorationShaman_Options.create()
				: RestorationShaman_Options.create({ classOptions: {} }),
	},
	// Warlock
	[Spec.SpecAfflictionWarlock]: {
		rotationCreate: () => AfflictionWarlock_Rotation.create(),
		rotationEquals: (a, b) => AfflictionWarlock_Rotation.equals(a as AfflictionWarlock_Rotation, b as AfflictionWarlock_Rotation),
		rotationCopy: a => AfflictionWarlock_Rotation.clone(a as AfflictionWarlock_Rotation),
		rotationToJson: a => AfflictionWarlock_Rotation.toJson(a as AfflictionWarlock_Rotation),
		rotationFromJson: obj => AfflictionWarlock_Rotation.fromJson(obj),

		talentsCreate: () => WarlockTalents.create(),
		talentsEquals: (a, b) => WarlockTalents.equals(a as WarlockTalents, b as WarlockTalents),
		talentsCopy: a => WarlockTalents.clone(a as WarlockTalents),
		talentsToJson: a => WarlockTalents.toJson(a as WarlockTalents),
		talentsFromJson: obj => WarlockTalents.fromJson(obj),

		optionsCreate: () => AfflictionWarlock_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => AfflictionWarlock_Options.equals(a as AfflictionWarlock_Options, b as AfflictionWarlock_Options),
		optionsCopy: a => AfflictionWarlock_Options.clone(a as AfflictionWarlock_Options),
		optionsToJson: a => AfflictionWarlock_Options.toJson(a as AfflictionWarlock_Options),
		optionsFromJson: obj => AfflictionWarlock_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'afflictionWarlock'
				? player.spec.afflictionWarlock.options || AfflictionWarlock_Options.create()
				: AfflictionWarlock_Options.create({ classOptions: {} }),
	},
	[Spec.SpecDemonologyWarlock]: {
		rotationCreate: () => DemonologyWarlock_Rotation.create(),
		rotationEquals: (a, b) => DemonologyWarlock_Rotation.equals(a as DemonologyWarlock_Rotation, b as DemonologyWarlock_Rotation),
		rotationCopy: a => DemonologyWarlock_Rotation.clone(a as DemonologyWarlock_Rotation),
		rotationToJson: a => DemonologyWarlock_Rotation.toJson(a as DemonologyWarlock_Rotation),
		rotationFromJson: obj => DemonologyWarlock_Rotation.fromJson(obj),

		talentsCreate: () => WarlockTalents.create(),
		talentsEquals: (a, b) => WarlockTalents.equals(a as WarlockTalents, b as WarlockTalents),
		talentsCopy: a => WarlockTalents.clone(a as WarlockTalents),
		talentsToJson: a => WarlockTalents.toJson(a as WarlockTalents),
		talentsFromJson: obj => WarlockTalents.fromJson(obj),

		optionsCreate: () => DemonologyWarlock_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => DemonologyWarlock_Options.equals(a as DemonologyWarlock_Options, b as DemonologyWarlock_Options),
		optionsCopy: a => DemonologyWarlock_Options.clone(a as DemonologyWarlock_Options),
		optionsToJson: a => DemonologyWarlock_Options.toJson(a as DemonologyWarlock_Options),
		optionsFromJson: obj => DemonologyWarlock_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'demonologyWarlock'
				? player.spec.demonologyWarlock.options || DemonologyWarlock_Options.create()
				: DemonologyWarlock_Options.create({ classOptions: {} }),
	},
	[Spec.SpecDestructionWarlock]: {
		rotationCreate: () => DestructionWarlock_Rotation.create(),
		rotationEquals: (a, b) => DestructionWarlock_Rotation.equals(a as DestructionWarlock_Rotation, b as DestructionWarlock_Rotation),
		rotationCopy: a => DestructionWarlock_Rotation.clone(a as DestructionWarlock_Rotation),
		rotationToJson: a => DestructionWarlock_Rotation.toJson(a as DestructionWarlock_Rotation),
		rotationFromJson: obj => DestructionWarlock_Rotation.fromJson(obj),

		talentsCreate: () => WarlockTalents.create(),
		talentsEquals: (a, b) => WarlockTalents.equals(a as WarlockTalents, b as WarlockTalents),
		talentsCopy: a => WarlockTalents.clone(a as WarlockTalents),
		talentsToJson: a => WarlockTalents.toJson(a as WarlockTalents),
		talentsFromJson: obj => WarlockTalents.fromJson(obj),

		optionsCreate: () => DestructionWarlock_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => DestructionWarlock_Options.equals(a as DestructionWarlock_Options, b as DestructionWarlock_Options),
		optionsCopy: a => DestructionWarlock_Options.clone(a as DestructionWarlock_Options),
		optionsToJson: a => DestructionWarlock_Options.toJson(a as DestructionWarlock_Options),
		optionsFromJson: obj => DestructionWarlock_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'destructionWarlock'
				? player.spec.destructionWarlock.options || DestructionWarlock_Options.create()
				: DestructionWarlock_Options.create({ classOptions: {} }),
	},
	// Warrior
	[Spec.SpecArmsWarrior]: {
		rotationCreate: () => ArmsWarrior_Rotation.create(),
		rotationEquals: (a, b) => ArmsWarrior_Rotation.equals(a as ArmsWarrior_Rotation, b as ArmsWarrior_Rotation),
		rotationCopy: a => ArmsWarrior_Rotation.clone(a as ArmsWarrior_Rotation),
		rotationToJson: a => ArmsWarrior_Rotation.toJson(a as ArmsWarrior_Rotation),
		rotationFromJson: obj => ArmsWarrior_Rotation.fromJson(obj),

		talentsCreate: () => WarriorTalents.create(),
		talentsEquals: (a, b) => WarriorTalents.equals(a as WarriorTalents, b as WarriorTalents),
		talentsCopy: a => WarriorTalents.clone(a as WarriorTalents),
		talentsToJson: a => WarriorTalents.toJson(a as WarriorTalents),
		talentsFromJson: obj => WarriorTalents.fromJson(obj),

		optionsCreate: () => ArmsWarrior_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => ArmsWarrior_Options.equals(a as ArmsWarrior_Options, b as ArmsWarrior_Options),
		optionsCopy: a => ArmsWarrior_Options.clone(a as ArmsWarrior_Options),
		optionsToJson: a => ArmsWarrior_Options.toJson(a as ArmsWarrior_Options),
		optionsFromJson: obj => ArmsWarrior_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'armsWarrior'
				? player.spec.armsWarrior.options || ArmsWarrior_Options.create()
				: ArmsWarrior_Options.create({ classOptions: {} }),
	},
	[Spec.SpecFuryWarrior]: {
		rotationCreate: () => FuryWarrior_Rotation.create(),
		rotationEquals: (a, b) => FuryWarrior_Rotation.equals(a as FuryWarrior_Rotation, b as FuryWarrior_Rotation),
		rotationCopy: a => FuryWarrior_Rotation.clone(a as FuryWarrior_Rotation),
		rotationToJson: a => FuryWarrior_Rotation.toJson(a as FuryWarrior_Rotation),
		rotationFromJson: obj => FuryWarrior_Rotation.fromJson(obj),

		talentsCreate: () => WarriorTalents.create(),
		talentsEquals: (a, b) => WarriorTalents.equals(a as WarriorTalents, b as WarriorTalents),
		talentsCopy: a => WarriorTalents.clone(a as WarriorTalents),
		talentsToJson: a => WarriorTalents.toJson(a as WarriorTalents),
		talentsFromJson: obj => WarriorTalents.fromJson(obj),

		optionsCreate: () => FuryWarrior_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => FuryWarrior_Options.equals(a as FuryWarrior_Options, b as FuryWarrior_Options),
		optionsCopy: a => FuryWarrior_Options.clone(a as FuryWarrior_Options),
		optionsToJson: a => FuryWarrior_Options.toJson(a as FuryWarrior_Options),
		optionsFromJson: obj => FuryWarrior_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'furyWarrior'
				? player.spec.furyWarrior.options || FuryWarrior_Options.create()
				: FuryWarrior_Options.create({ classOptions: {} }),
	},
	[Spec.SpecProtectionWarrior]: {
		rotationCreate: () => ProtectionWarrior_Rotation.create(),
		rotationEquals: (a, b) => ProtectionWarrior_Rotation.equals(a as ProtectionWarrior_Rotation, b as ProtectionWarrior_Rotation),
		rotationCopy: a => ProtectionWarrior_Rotation.clone(a as ProtectionWarrior_Rotation),
		rotationToJson: a => ProtectionWarrior_Rotation.toJson(a as ProtectionWarrior_Rotation),
		rotationFromJson: obj => ProtectionWarrior_Rotation.fromJson(obj),

		talentsCreate: () => WarriorTalents.create(),
		talentsEquals: (a, b) => WarriorTalents.equals(a as WarriorTalents, b as WarriorTalents),
		talentsCopy: a => WarriorTalents.clone(a as WarriorTalents),
		talentsToJson: a => WarriorTalents.toJson(a as WarriorTalents),
		talentsFromJson: obj => WarriorTalents.fromJson(obj),

		optionsCreate: () => ProtectionWarrior_Options.create({ classOptions: {} }),
		optionsEquals: (a, b) => ProtectionWarrior_Options.equals(a as ProtectionWarrior_Options, b as ProtectionWarrior_Options),
		optionsCopy: a => ProtectionWarrior_Options.clone(a as ProtectionWarrior_Options),
		optionsToJson: a => ProtectionWarrior_Options.toJson(a as ProtectionWarrior_Options),
		optionsFromJson: obj => ProtectionWarrior_Options.fromJson(obj),
		optionsFromPlayer: player =>
			player.spec.oneofKind == 'protectionWarrior'
				? player.spec.protectionWarrior.options || ProtectionWarrior_Options.create()
				: ProtectionWarrior_Options.create(),
	},
};

export const raceToFaction: Record<Race, Faction> = {
	[Race.RaceUnknown]: Faction.Unknown,

	[Race.RaceDraenei]: Faction.Alliance,
	[Race.RaceDwarf]: Faction.Alliance,
	[Race.RaceGnome]: Faction.Alliance,
	[Race.RaceHuman]: Faction.Alliance,
	[Race.RaceNightElf]: Faction.Alliance,
	[Race.RaceWorgen]: Faction.Alliance,
	[Race.RaceAlliancePandaren]: Faction.Alliance,

	[Race.RaceBloodElf]: Faction.Horde,
	[Race.RaceGoblin]: Faction.Horde,
	[Race.RaceOrc]: Faction.Horde,
	[Race.RaceTauren]: Faction.Horde,
	[Race.RaceTroll]: Faction.Horde,
	[Race.RaceUndead]: Faction.Horde,
	[Race.RaceHordePandaren]: Faction.Horde,
};

// Returns a copy of playerOptions, with the class field set.
export function withSpec<SpecType extends Spec>(spec: Spec, player: Player, specOptions: SpecOptions<SpecType>): Player {
	const copy = Player.clone(player);

	switch (spec) {
		// Death Knight
		case Spec.SpecBloodDeathKnight:
			copy.spec = {
				oneofKind: 'bloodDeathKnight',
				bloodDeathKnight: BloodDeathKnight.create({
					options: specOptions as BloodDeathKnight_Options,
				}),
			};
			return copy;
		case Spec.SpecFrostDeathKnight:
			copy.spec = {
				oneofKind: 'frostDeathKnight',
				frostDeathKnight: FrostDeathKnight.create({
					options: specOptions as FrostDeathKnight_Options,
				}),
			};
			return copy;
		case Spec.SpecUnholyDeathKnight:
			copy.spec = {
				oneofKind: 'unholyDeathKnight',
				unholyDeathKnight: UnholyDeathKnight.create({
					options: specOptions as UnholyDeathKnight_Options,
				}),
			};
			return copy;
		// Druid
		case Spec.SpecBalanceDruid:
			copy.spec = {
				oneofKind: 'balanceDruid',
				balanceDruid: BalanceDruid.create({
					options: specOptions as BalanceDruid_Options,
				}),
			};
			return copy;
		case Spec.SpecFeralDruid:
			copy.spec = {
				oneofKind: 'feralDruid',
				feralDruid: FeralDruid.create({
					options: specOptions as FeralDruid_Options,
				}),
			};
			return copy;
		case Spec.SpecGuardianDruid:
			copy.spec = {
				oneofKind: 'guardianDruid',
				guardianDruid: GuardianDruid.create({
					options: specOptions as GuardianDruid_Options,
				}),
			};
			return copy;
		case Spec.SpecRestorationDruid:
			copy.spec = {
				oneofKind: 'restorationDruid',
				restorationDruid: RestorationDruid.create({
					options: specOptions as RestorationDruid_Options,
				}),
			};
			return copy;
		// Hunter
		case Spec.SpecBeastMasteryHunter:
			copy.spec = {
				oneofKind: 'beastMasteryHunter',
				beastMasteryHunter: BeastMasteryHunter.create({
					options: specOptions as BeastMasteryHunter_Options,
				}),
			};
			return copy;
		case Spec.SpecMarksmanshipHunter:
			copy.spec = {
				oneofKind: 'marksmanshipHunter',
				marksmanshipHunter: MarksmanshipHunter.create({
					options: specOptions as MarksmanshipHunter_Options,
				}),
			};
			return copy;
		case Spec.SpecSurvivalHunter:
			copy.spec = {
				oneofKind: 'survivalHunter',
				survivalHunter: SurvivalHunter.create({
					options: specOptions as SurvivalHunter_Options,
				}),
			};
			return copy;
		// Mage
		case Spec.SpecArcaneMage:
			copy.spec = {
				oneofKind: 'arcaneMage',
				arcaneMage: ArcaneMage.create({
					options: specOptions as ArcaneMage_Options,
				}),
			};
			return copy;
		case Spec.SpecFireMage:
			copy.spec = {
				oneofKind: 'fireMage',
				fireMage: FireMage.create({
					options: specOptions as FireMage_Options,
				}),
			};
			return copy;
		case Spec.SpecFrostMage:
			copy.spec = {
				oneofKind: 'frostMage',
				frostMage: FrostMage.create({
					options: specOptions as FrostMage_Options,
				}),
			};
			return copy;
		// Monk
		case Spec.SpecBrewmasterMonk:
			copy.spec = {
				oneofKind: 'brewmasterMonk',
				brewmasterMonk: BrewmasterMonk.create({
					options: specOptions as BrewmasterMonk_Options,
				}),
			};
			return copy;
		case Spec.SpecMistweaverMonk:
			copy.spec = {
				oneofKind: 'mistweaverMonk',
				mistweaverMonk: MistweaverMonk.create({
					options: specOptions as MistweaverMonk_Options,
				}),
			};
			return copy;
		case Spec.SpecWindwalkerMonk:
			copy.spec = {
				oneofKind: 'windwalkerMonk',
				windwalkerMonk: WindwalkerMonk.create({
					options: specOptions as WindwalkerMonk_Options,
				}),
			};
			return copy;
		// Paladin
		case Spec.SpecHolyPaladin:
			copy.spec = {
				oneofKind: 'holyPaladin',
				holyPaladin: HolyPaladin.create({
					options: specOptions as HolyPaladin_Options,
				}),
			};
			return copy;
		case Spec.SpecProtectionPaladin:
			copy.spec = {
				oneofKind: 'protectionPaladin',
				protectionPaladin: ProtectionPaladin.create({
					options: specOptions as ProtectionPaladin_Options,
				}),
			};
			return copy;
		case Spec.SpecRetributionPaladin:
			copy.spec = {
				oneofKind: 'retributionPaladin',
				retributionPaladin: RetributionPaladin.create({
					options: specOptions as RetributionPaladin_Options,
				}),
			};
			return copy;
		// Priest
		case Spec.SpecDisciplinePriest:
			copy.spec = {
				oneofKind: 'disciplinePriest',
				disciplinePriest: DisciplinePriest.create({
					options: specOptions as DisciplinePriest_Options,
				}),
			};
			return copy;
		case Spec.SpecHolyPriest:
			copy.spec = {
				oneofKind: 'holyPriest',
				holyPriest: HolyPriest.create({
					options: specOptions as HolyPriest_Options,
				}),
			};
			return copy;
		case Spec.SpecShadowPriest:
			copy.spec = {
				oneofKind: 'shadowPriest',
				shadowPriest: ShadowPriest.create({
					options: specOptions as ShadowPriest_Options,
				}),
			};
			return copy;
		// Rogue
		case Spec.SpecAssassinationRogue:
			copy.spec = {
				oneofKind: 'assassinationRogue',
				assassinationRogue: AssassinationRogue.create({
					options: specOptions as AssassinationRogue_Options,
				}),
			};
			return copy;
		case Spec.SpecCombatRogue:
			copy.spec = {
				oneofKind: 'combatRogue',
				combatRogue: CombatRogue.create({
					options: specOptions as CombatRogue_Options,
				}),
			};
			return copy;
		case Spec.SpecSubtletyRogue:
			copy.spec = {
				oneofKind: 'subtletyRogue',
				subtletyRogue: SubtletyRogue.create({
					options: specOptions as SubtletyRogue_Options,
				}),
			};
			return copy;
		// Shaman
		case Spec.SpecElementalShaman:
			copy.spec = {
				oneofKind: 'elementalShaman',
				elementalShaman: ElementalShaman.create({
					options: specOptions as ElementalShaman_Options,
				}),
			};
			return copy;
		case Spec.SpecEnhancementShaman:
			copy.spec = {
				oneofKind: 'enhancementShaman',
				enhancementShaman: EnhancementShaman.create({
					options: specOptions as EnhancementShaman_Options,
				}),
			};
			return copy;
		case Spec.SpecRestorationShaman:
			copy.spec = {
				oneofKind: 'restorationShaman',
				restorationShaman: RestorationShaman.create({
					options: specOptions as RestorationShaman_Options,
				}),
			};
			return copy;
		// Warlock
		case Spec.SpecAfflictionWarlock:
			copy.spec = {
				oneofKind: 'afflictionWarlock',
				afflictionWarlock: AfflictionWarlock.create({
					options: specOptions as AfflictionWarlock_Options,
				}),
			};
			return copy;
		case Spec.SpecDemonologyWarlock:
			copy.spec = {
				oneofKind: 'demonologyWarlock',
				demonologyWarlock: DemonologyWarlock.create({
					options: specOptions as DemonologyWarlock_Options,
				}),
			};
			return copy;
		case Spec.SpecDestructionWarlock:
			copy.spec = {
				oneofKind: 'destructionWarlock',
				destructionWarlock: DestructionWarlock.create({
					options: specOptions as DestructionWarlock_Options,
				}),
			};
			return copy;
		// Warrior
		case Spec.SpecArmsWarrior:
			copy.spec = {
				oneofKind: 'armsWarrior',
				armsWarrior: ArmsWarrior.create({
					options: specOptions as ArmsWarrior_Options,
				}),
			};
			return copy;
		case Spec.SpecFuryWarrior:
			copy.spec = {
				oneofKind: 'furyWarrior',
				furyWarrior: FuryWarrior.create({
					options: specOptions as FuryWarrior_Options,
				}),
			};
			return copy;
		case Spec.SpecProtectionWarrior:
			copy.spec = {
				oneofKind: 'protectionWarrior',
				protectionWarrior: ProtectionWarrior.create({
					options: specOptions as ProtectionWarrior_Options,
				}),
			};
			return copy;
		default:
			return copy;
	}
}

export function getPlayerSpecFromPlayer<SpecType extends Spec>(player: Player): PlayerSpec<SpecType> {
	const specValues = getEnumValues(Spec);
	for (let i = 0; i < specValues.length; i++) {
		const spec = specValues[i] as SpecType;
		let specString = Spec[spec]; // Returns 'SpecBalanceDruid' for BalanceDruid.
		specString = specString.substring('Spec'.length); // 'BalanceDruid'
		specString = specString.charAt(0).toLowerCase() + specString.slice(1); // 'balanceDruid'

		if (player.spec.oneofKind == specString) {
			return PlayerSpecs.fromProto(spec);
		}
	}

	throw new Error('Unable to parse spec from player proto: ' + JSON.stringify(Player.toJson(player), null, 2));
}

export function isSharpWeaponType(weaponType: WeaponType): boolean {
	return [WeaponType.WeaponTypeAxe, WeaponType.WeaponTypeDagger, WeaponType.WeaponTypePolearm, WeaponType.WeaponTypeSword].includes(weaponType);
}

export function isBluntWeaponType(weaponType: WeaponType): boolean {
	return [WeaponType.WeaponTypeFist, WeaponType.WeaponTypeMace, WeaponType.WeaponTypeStaff].includes(weaponType);
}

// Custom functions for determining the EP value of meta gem effects.
// Default meta effect EP value is 0, so just handle the ones relevant to your spec.
const metaGemEffectEPs: Partial<Record<Spec, (gem: Gem, playerStats: Stats) => number>> = {};

export function getMetaGemEffectEP<SpecType extends Spec>(playerSpec: PlayerSpec<SpecType>, gem: Gem, playerStats: Stats) {
	if (metaGemEffectEPs[playerSpec.specID]) {
		return metaGemEffectEPs[playerSpec.specID]!(gem, playerStats);
	} else {
		return 0;
	}
}

// Returns true if this item may be equipped in at least 1 slot for the given Spec.
export function canEquipItem<SpecType extends Spec>(item: Item, playerSpec: PlayerSpec<SpecType>, slot: ItemSlot | undefined): boolean {
	const playerClass = PlayerSpecs.getPlayerClass(playerSpec);
	if (item.classAllowlist.length > 0 && !item.classAllowlist.includes(playerClass.classID)) {
		return false;
	}

	if ([ItemType.ItemTypeFinger, ItemType.ItemTypeTrinket].includes(item.type)) {
		return true;
	}

	if (item.type == ItemType.ItemTypeWeapon) {
		const eligibleWeaponType = playerClass.weaponTypes.find(wt => wt.weaponType == item.weaponType);
		if (!eligibleWeaponType) {
			return false;
		}

		if (
			(item.handType == HandType.HandTypeOffHand || (item.handType == HandType.HandTypeOneHand && slot == ItemSlot.ItemSlotOffHand)) &&
			![WeaponType.WeaponTypeShield, WeaponType.WeaponTypeOffHand].includes(item.weaponType) &&
			!playerSpec.canDualWield
		) {
			return false;
		}

		if (item.handType == HandType.HandTypeTwoHand && !eligibleWeaponType.canUseTwoHand) {
			return false;
		}
		if (item.handType == HandType.HandTypeTwoHand && slot == ItemSlot.ItemSlotOffHand && playerSpec.specID != Spec.SpecFuryWarrior) {
			return false;
		}

		return true;
	}

	if (item.type == ItemType.ItemTypeRanged) {
		return playerClass.rangedWeaponTypes.includes(item.rangedWeaponType);
	}

	// At this point, we know the item is an armor piece (feet, chest, legs, etc).
	return playerClass.armorTypes[0] >= item.armorType;
}

const pvpSeasonFromName: Record<string, string> = {
	Wrathful: 'Season 8',
	Bloodthirsty: 'Season 8.5',
	Vicious: 'Season 9',
	Ruthless: 'Season 10',
	Cataclysmic: 'Season 11',
};

export const isPVPItem = (item: Item) => item?.name?.includes('Gladiator') || false;

export const getPVPSeasonFromItem = (item: Item) => {
	const seasonName = item.name.substring(0, item.name.indexOf(' '));
	return pvpSeasonFromName[seasonName] || undefined;
};

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
	[ItemType.ItemTypeRanged]: [ItemSlot.ItemSlotMainHand],
};

export function getEligibleItemSlots(item: Item, isFuryWarrior?: boolean): Array<ItemSlot> {
	if (itemTypeToSlotsMap[item.type]) {
		return itemTypeToSlotsMap[item.type]!;
	}

	if (item.type == ItemType.ItemTypeWeapon) {
		if (isFuryWarrior) {
			return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
		}

		if (item.handType == HandType.HandTypeMainHand) {
			return [ItemSlot.ItemSlotMainHand];
		} else if (item.handType == HandType.HandTypeOffHand) {
			return [ItemSlot.ItemSlotOffHand];
		} else {
			return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
		}
	}

	// Should never reach here
	throw new Error('Could not find item slots for item: ' + Item.toJsonString(item));
}

export const isSecondaryItemSlot = (slot: ItemSlot) => slot === ItemSlot.ItemSlotFinger2 || slot === ItemSlot.ItemSlotTrinket2;

// Returns whether the given main-hand and off-hand items can be worn at the
// same time.
export function validWeaponCombo(mainHand: Item | null | undefined, offHand: Item | null | undefined, canDW2h: boolean): boolean {
	if (mainHand?.handType == HandType.HandTypeTwoHand && !canDW2h) {
		return false;
	} else if (
		mainHand?.handType == HandType.HandTypeTwoHand &&
		(mainHand?.weaponType == WeaponType.WeaponTypePolearm || mainHand?.weaponType == WeaponType.WeaponTypeStaff)
	) {
		return false;
	}
	if (offHand?.handType == HandType.HandTypeTwoHand && !canDW2h) {
		return false;
	} else if (
		offHand?.handType == HandType.HandTypeTwoHand &&
		(offHand?.weaponType == WeaponType.WeaponTypePolearm || offHand?.weaponType == WeaponType.WeaponTypeStaff)
	) {
		return false;
	}

	return true;
}

// Returns all item slots to which the enchant might be applied.
//
// Note that this alone is not enough; some items have further restrictions,
// e.g. some weapon enchants may only be applied to 2H weapons.
export function getEligibleEnchantSlots(enchant: Enchant): Array<ItemSlot> {
	return [enchant.type]
		.concat(enchant.extraTypes || [])
		.map(type => {
			if (itemTypeToSlotsMap[type]) {
				return itemTypeToSlotsMap[type]!;
			}

			if (type == ItemType.ItemTypeWeapon) {
				return [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];
			}

			// Should never reach here
			throw new Error('Could not find item slots for enchant: ' + Enchant.toJsonString(enchant));
		})
		.flat();
}

export function enchantAppliesToItem(enchant: Enchant, item: Item): boolean {
	const sharedSlots = intersection(getEligibleEnchantSlots(enchant), getEligibleItemSlots(item));
	if (!sharedSlots.length) return false;

	if (enchant.enchantType === EnchantType.EnchantTypeTwoHand && item.handType !== HandType.HandTypeTwoHand) return false;

	if (enchant.enchantType === EnchantType.EnchantTypeStaff && item.weaponType !== WeaponType.WeaponTypeStaff) return false;

	if (enchant.enchantType === EnchantType.EnchantTypeShield && item.weaponType !== WeaponType.WeaponTypeShield) return false;

	if (
		(enchant.enchantType === EnchantType.EnchantTypeOffHand) !==
		(item.weaponType === WeaponType.WeaponTypeOffHand ||
			// All off-hand enchants can be applied to shields as well
			(item.weaponType === WeaponType.WeaponTypeShield && enchant.enchantType !== EnchantType.EnchantTypeShield))
	)
		return false;

	if (enchant.type == ItemType.ItemTypeRanged) {
		if (
			![RangedWeaponType.RangedWeaponTypeBow, RangedWeaponType.RangedWeaponTypeCrossbow, RangedWeaponType.RangedWeaponTypeGun].includes(
				item.rangedWeaponType,
			)
		)
			return false;
	}

	if (item.rangedWeaponType > 0 && enchant.type != ItemType.ItemTypeRanged) {
		return false;
	}

	return true;
}

export function canEquipEnchant<SpecType extends Spec>(enchant: Enchant, playerSpec: PlayerSpec<SpecType>): boolean {
	if (enchant.classAllowlist.length > 0 && !enchant.classAllowlist.includes(playerSpec.classID)) {
		return false;
	}

	// This is a Tinker and we handle them differently
	if (enchant.requiredProfession == Profession.Engineering) {
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
		assignments.paladins.push(
			BlessingsAssignment.create({
				blessings: new Array(NUM_SPECS).fill(Blessings.BlessingUnknown),
			}),
		);
	}
	return assignments;
}

export function makeBlessingsAssignments(numPaladins: number): BlessingsAssignments {
	const assignments = makeBlankBlessingsAssignments(numPaladins);
	for (let i = 1; i < Object.keys(Spec).length; i++) {
		const spec = i;
		const blessings = [Blessings.BlessingOfKings, Blessings.BlessingOfMight];
		for (let j = 0; j < blessings.length; j++) {
			if (j >= assignments.paladins.length) {
				// Can't assign more blessings since we ran out of paladins
				break;
			}
			assignments.paladins[j].blessings[spec] = blessings[j];
		}
	}
	return assignments;
}

// Default blessings settings in the raid sim UI.
export function makeDefaultBlessings(numPaladins: number): BlessingsAssignments {
	return makeBlessingsAssignments(numPaladins);
}

export const orderedResourceTypes: Array<ResourceType> = [
	ResourceType.ResourceTypeHealth,
	ResourceType.ResourceTypeMana,
	ResourceType.ResourceTypeEnergy,
	ResourceType.ResourceTypeRage,
	ResourceType.ResourceTypeChi,
	ResourceType.ResourceTypeComboPoints,
	ResourceType.ResourceTypeFocus,
	ResourceType.ResourceTypeRunicPower,
	ResourceType.ResourceTypeBloodRune,
	ResourceType.ResourceTypeFrostRune,
	ResourceType.ResourceTypeUnholyRune,
	ResourceType.ResourceTypeDeathRune,
	ResourceType.ResourceTypeLunarEnergy,
	ResourceType.ResourceTypeSolarEnergy,
	ResourceType.ResourceTypeGenericResource,
];

export const AL_CATEGORY_HARD_MODE = 'Hard Mode';
export const AL_CATEGORY_TITAN_RUNE = 'Titan Rune';

// Utilities for migrating protos between versions

// Each key is an API version, each value is a function that up-converts a proto
// to that version from the previous one. If there are missing keys between
// successive entries, then it is assumed that no intermediate conversions are
// required (i.e. the intermediate version changes did not affect this
// particular proto).
export type ProtoConversionMap<Type> = Map<number, (arg: Type) => Type>;

export function migrateOldProto<Type>(oldProto: Type, oldApiVersion: number, conversionMap: ProtoConversionMap<Type>, targetApiVersion?: number): Type {
	let migratedProto = oldProto;
	const finalVersion = targetApiVersion || CURRENT_API_VERSION;

	for (let nextVersion = oldApiVersion + 1; nextVersion <= finalVersion; nextVersion++) {
		if (conversionMap.has(nextVersion)) {
			migratedProto = conversionMap.get(nextVersion)!(migratedProto);
		}
	}

	return migratedProto;
}
