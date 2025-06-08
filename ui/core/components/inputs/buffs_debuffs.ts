import { Stat } from '../../proto/common';
import { ActionId } from '../../proto_utils/action_id';
import {
	makeBooleanDebuffInput,
	makeBooleanRaidBuffInput,
	makeMultistateIndividualBuffInput,
	makeMultistateRaidBuffInput,
	makeTristateIndividualBuffInput,
} from '../icon_inputs';
import * as InputHelpers from '../input_helpers';
import { IconPicker } from '../pickers/icon_picker';
import { MultiIconPicker } from '../pickers/multi_icon_picker';
import { IconPickerStatOption, PickerStatOptions } from './stat_options';

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

export const StatsBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(20217), fieldName: 'blessingOfKings' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(1126), fieldName: 'markOfTheWild' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(90363), fieldName: 'embraceOfTheShaleSpider' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(115921), fieldName: 'legacyOfTheEmperor' }),
	],
	'Stats',
);

export const AttackPowerBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(19506), fieldName: 'trueshotAura' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(57330), fieldName: 'hornOfWinter' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(6673), fieldName: 'battleShout' }),
	],
	'Attack Power',
);

export const AttackSpeedBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(55610), fieldName: 'unholyAura' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(128433), fieldName: 'serpentsSwiftness' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(113742), fieldName: 'swiftbladesCunning' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(30809), fieldName: 'unleashedRage' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(128432), fieldName: 'cacklingHowl' }),
	],
	'Attack Speed',
);

export const SpellPowerBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(1459), fieldName: 'arcaneBrilliance' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(126309), fieldName: 'stillWater' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(77747), fieldName: 'burningWrath' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(109773), fieldName: 'darkIntent' }),
	],
	'Spell Power',
);

export const SpellHasteBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(24907), fieldName: 'moonkinAura' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(49868), fieldName: 'mindQuickening' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(51470), fieldName: 'elementalOath' }),
	],
	'Spell Haste',
);

export const CritBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(17007), fieldName: 'leaderOfThePack' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(24604), fieldName: 'furiousHowl' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(90309), fieldName: 'terrifyingRoar' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(116781), fieldName: 'legacyOfTheWhiteTiger' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(1459), fieldName: 'arcaneBrilliance' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(126309), fieldName: 'stillWater' }),
	],
	'Crit %',
);

export const MasteryBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(19740), fieldName: 'blessingOfMight' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(93435), fieldName: 'roarOfCourage' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(128997), fieldName: 'spiritBeastBlessing' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(116956), fieldName: 'graceOfAir' }),
	],
	'Mastery',
);

export const StaminaBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(469), fieldName: 'commandingShout' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(109773), fieldName: 'darkIntent' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(21562), fieldName: 'powerWordFortitude' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(90364), fieldName: 'qirajiFortitude' }),
	],
	'Stamina',
);

export const MajorHasteBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(2825), fieldName: 'bloodlust' }),
	],
	'Major Haste',
);

// TODO: Look at these, what we want and how to structure them for multiple available
export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput(
	[
		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(6940), numStates: 11, fieldName: 'handOfSacrificeCount' }),
		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(31821), numStates: 11, fieldName: 'devotionAuraCount' }),
		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(33206), numStates: 11, fieldName: 'painSuppressionCount' }),
		// 		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(47788), numStates: 11, fieldName: 'guardianSpirits' }),
		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(97462), numStates: 11, fieldName: 'rallyingCryCount' }),
	],
	'Defensive CDs',
);

// Misc Buffs
export const ManaTideTotem = makeMultistateRaidBuffInput({ actionId: ActionId.fromSpellId(16190), numStates: 5, fieldName: 'manaTideTotemCount' });
export const TricksOfTheTrade = makeTristateIndividualBuffInput({
	actionId: ActionId.fromItemId(45767),
	impId: ActionId.fromSpellId(57933),
	fieldName: 'tricksOfTheTrade',
});
export const UnholyFrenzy = makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(49016), numStates: 11, fieldName: 'unholyFrenzyCount' });
export const ShatteringThrow = makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(64382), numStates: 11, fieldName: 'shatteringThrowCount' });
export const Skullbanner = makeMultistateRaidBuffInput({ actionId: ActionId.fromSpellId(114207), numStates: 11, fieldName: 'skullBannerCount' });
export const StormLashTotem = makeMultistateRaidBuffInput({ actionId: ActionId.fromSpellId(120668), numStates: 11, fieldName: 'stormlashTotemCount' });

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(113746), fieldName: 'weakenedArmor' });

export const DamageReduction = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(115798), fieldName: 'weakenedBlows' });

export const CastSpeedDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(73975), fieldName: 'necroticStrike' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(58604), fieldName: 'lavaBreath' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(50274), fieldName: 'sporeCloud' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(5761), fieldName: 'mindNumbingPoison' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(31589), fieldName: 'slow' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(109466), fieldName: 'curseOfEnfeeblement' }),
	],
	'Cast Speed',
);

export const PhysicalDamageDebuff = makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(81326), fieldName: 'physicalVulnerability' });

export const SpellDamageDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(24844), fieldName: 'lightningBreath' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(1490), fieldName: 'curseOfElements' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(58410), fieldName: 'masterPoisoner' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(34889), fieldName: 'fireBreath' }),
	],
	'Spell Dmg',
);

///////////////////////////////////////////////////////////////////////////
//                                 CONFIGS
///////////////////////////////////////////////////////////////////////////

export const RAID_BUFFS_CONFIG = [
	// Standard buffs
	{
		config: StatsBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStrength, Stat.StatAgility, Stat.StatIntellect],
	},
	{
		config: AttackPowerBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: AttackSpeedBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: SpellPowerBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellPower],
	},
	{
		config: SpellHasteBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellPower],
	},
	{
		config: CritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatCritRating],
	},
	{
		config: MasteryBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMasteryRating],
	},
	{
		config: StaminaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina],
	},
	{
		config: MajorHasteBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatHasteRating],
	},
	{
		config: DefensiveCooldownBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina],
	},
] as PickerStatOptions[];

export const RAID_BUFFS_MISC_CONFIG = [
	{
		config: Skullbanner,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
	},
	{
		config: StormLashTotem,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
	},
	{
		config: ManaTideTotem,
		picker: IconPicker,
		stats: [Stat.StatSpirit],
	},
	{
		config: TricksOfTheTrade,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
	},
	{
		config: UnholyFrenzy,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: ShatteringThrow,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
] as IconPickerStatOption[];

export const DEBUFFS_CONFIG = [
	{
		config: MajorArmorDebuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: PhysicalDamageDebuff,
		picker: IconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: SpellDamageDebuff,
		picker: MultiIconPicker,
		// Enabled for all specs because it affects Stormlash Totem
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
	},
	{
		config: DamageReduction,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: CastSpeedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
] as PickerStatOptions[];

export const DEBUFFS_MISC_CONFIG = [] as IconPickerStatOption[];
