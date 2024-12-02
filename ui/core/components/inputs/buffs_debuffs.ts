import { Faction, Stat, TristateEffect } from '../../proto/common';
import { ActionId } from '../../proto_utils/action_id';
import {
	makeBooleanDebuffInput,
	makeBooleanIndividualBuffInput,
	makeBooleanPartyBuffInput,
	makeBooleanRaidBuffInput,
	makeMultistateIndividualBuffInput,
	makeMultistateMultiplierIndividualBuffInput,
	makeMultistatePartyBuffInput,
	makeMultistateRaidBuffInput,
	makeQuadstateDebuffInput,
	makeTristateDebuffInput,
	makeTristateIndividualBuffInput,
	makeTristateRaidBuffInput,
	withLabel,
} from '../icon_inputs';
import * as InputHelpers from '../input_helpers';
import { IconPicker } from '../pickers/icon_picker';
import { MultiIconPicker } from '../pickers/multi_icon_picker';
import { IconPickerStatOption, PickerStatOptions } from './stat_options';

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

export const AllStatsBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(20217), fieldName: 'blessingOfKings' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(1126), fieldName: 'markOfTheWild' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(116781), fieldName: 'legacyOfTheWhiteTiger' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromItemId(63140), fieldName: 'drumsOfTheBurningWild' }),
	],
	'Stats',
);

export const ArmorBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(465), fieldName: 'devotionAura' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(8071), fieldName: 'stoneskinTotem' }),
	],
	'Armor',
);

export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(19740), fieldName: 'blessingOfMight' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(53138), fieldName: 'abominationsMight' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(30808), fieldName: 'unleashedRage' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(19506), fieldName: 'trueshotAura' }),
	],
	'Atk Pwr %',
);

export const Bloodlust = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(2825), fieldName: 'bloodlust' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(80353), fieldName: 'timeWarp' }),
	],
	'Lust',
);

export const DamagePercentBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(31876), fieldName: 'communion' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(82930), fieldName: 'arcaneTactics' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(34460), fieldName: 'ferociousInspiration' }),
	],
	'+3% Dmg',
);

// TODO: Look at these, what we want and how to structure them for multiple available
export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput(
	[
		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(6940), numStates: 11, fieldName: 'handOfSacrificeCount' }),
		// 		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(53530), numStates: 11, fieldName: 'divineGuardians' }),
		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(33206), numStates: 11, fieldName: 'painSuppressionCount' }),
		// 		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(47788), numStates: 11, fieldName: 'guardianSpirits' }),
		makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(97462), numStates: 11, fieldName: 'rallyingCryCount' }),
	],
	'Defensive CDs',
);

export const SpellHasteBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(15473), fieldName: 'shadowForm' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(24858), fieldName: 'moonkinForm' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(3738), fieldName: 'wrathOfAirTotem' }),
	],
	'Spell Haste',
);

export const ManaBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(1459), fieldName: 'arcaneBrilliance' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(54424), fieldName: 'felIntelligence' }),
	],
	'Mana',
);

export const CritBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(17007), fieldName: 'leaderOfThePack' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(51470), fieldName: 'elementalOath' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(51701), fieldName: 'honorAmongThieves' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(29801), fieldName: 'rampage' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(24604), fieldName: 'furiousHowl' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(115921), fieldName: 'legacyOfTheEmperor' }),
	],
	'Crit %',
);

export const MeleeHasteBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(55610), fieldName: 'icyTalons' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(8512), fieldName: 'windfuryTotem' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(53290), fieldName: 'huntingParty' }),
	],
	'Melee Haste',
);

export const MP5Buff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(19740), fieldName: 'blessingOfMight' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(54424), fieldName: 'felIntelligence' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(5675), fieldName: 'manaSpringTotem' }),
	],
	'MP5',
);

export const ReplenishmentBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanIndividualBuffInput({ actionId: ActionId.fromSpellId(34914), fieldName: 'vampiricTouch' }),
		makeBooleanIndividualBuffInput({ actionId: ActionId.fromSpellId(31876), fieldName: 'communion' }),
		makeBooleanIndividualBuffInput({ actionId: ActionId.fromSpellId(48544), fieldName: 'revitalize' }),
		makeBooleanIndividualBuffInput({ actionId: ActionId.fromSpellId(30295), fieldName: 'soulLeach' }),
		makeBooleanIndividualBuffInput({ actionId: ActionId.fromSpellId(86508), fieldName: 'enduringWinter' }),
	],
	'Replen',
);

export const ResistanceBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(19891), fieldName: 'resistanceAura' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(8184), fieldName: 'elementalResistanceTotem' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(20043), fieldName: 'aspectOfTheWild' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(27683), fieldName: 'shadowProtection' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(20217), fieldName: 'blessingOfKings' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(1126), fieldName: 'markOfTheWild' }),
	],
	'Resistances',
);

export const SpellPowerBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(47236), fieldName: 'demonicPact' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(77746), fieldName: 'totemicWrath' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(1459), fieldName: 'arcaneBrilliance' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(8227), fieldName: 'flametongueTotem' }),
	],
	'Spell Power',
);

export const StaminaBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(21562), fieldName: 'powerWordFortitude' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(6307), fieldName: 'bloodPact' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(469), fieldName: 'commandingShout' }),
	],
	'Stamina',
);

export const StrengthAndAgilityBuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(8075), fieldName: 'strengthOfEarthTotem' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(57330), fieldName: 'hornOfWinter' }),
		makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(6673), fieldName: 'battleShout' }),
	],
	'Str/Agi',
);

// Misc Buffs
export const RetributionAura = makeBooleanRaidBuffInput({ actionId: ActionId.fromSpellId(7294), fieldName: 'retributionAura' });
export const ManaTideTotem = makeMultistateRaidBuffInput({ actionId: ActionId.fromSpellId(16190), numStates: 5, fieldName: 'manaTideTotemCount' });
export const Innervate = makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(29166), numStates: 11, fieldName: 'innervateCount' });
export const PowerInfusion = makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(10060), numStates: 11, fieldName: 'powerInfusionCount' });
export const FocusMagic = makeBooleanIndividualBuffInput({ actionId: ActionId.fromSpellId(54648), fieldName: 'focusMagic' });
export const TricksOfTheTrade = makeTristateIndividualBuffInput({
	actionId: ActionId.fromItemId(45767),
	impId: ActionId.fromSpellId(57933),
	fieldName: 'tricksOfTheTrade',
});
export const UnholyFrenzy = makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(49016), numStates: 11, fieldName: 'unholyFrenzyCount' });
export const DarkIntent = makeBooleanIndividualBuffInput({ actionId: ActionId.fromSpellId(85759), fieldName: 'darkIntent' });
export const ShatteringThrow = makeMultistateIndividualBuffInput({ actionId: ActionId.fromSpellId(64382), numStates: 11, fieldName: 'shatteringThrowCount' });

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(7386), fieldName: 'sunderArmor' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(8647), fieldName: 'exposeArmor' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(770), fieldName: 'faerieFire' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(35387), fieldName: 'corrosiveSpit' }),
	],
	'-Armor %',
);

export const DamageReduction = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(26017), fieldName: 'vindication' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(702), fieldName: 'curseOfWeakness' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(99), fieldName: 'demoralizingRoar' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(81130), fieldName: 'scarletFever' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(1160), fieldName: 'demoralizingShout' }),
	],
	'-Dmg %',
);

export const BleedDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(29859), fieldName: 'bloodFrenzy' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(33878), fieldName: 'mangle' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(57386), fieldName: 'stampede' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(16511), fieldName: 'hemorrhage' }),
	],
	'+Bleed %',
);

export const SpellCritDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(12873), fieldName: 'criticalMass' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(17801), fieldName: 'shadowAndFlame' }),
	],
	'Spell Crit',
);

export const MeleeAttackSpeedDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(6343), fieldName: 'thunderClap' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(59921), fieldName: 'frostFever' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(53696), fieldName: 'judgementsOfTheJust' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(48484), fieldName: 'infectedWounds' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(8042), fieldName: 'earthShock' }),
	],
	'Atk Speed',
);

export const PhysicalDamageDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(29859), fieldName: 'bloodFrenzy' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(58413), fieldName: 'savageCombat' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(81328), fieldName: 'brittleBones' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(55749), fieldName: 'acidSpit' }),
	],
	'Phys Vuln',
);

export const SpellDamageDebuff = InputHelpers.makeMultiIconInput(
	[
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(51160), fieldName: 'ebonPlaguebringer' }),
		makeBooleanDebuffInput({ actionId: ActionId.fromSpellId(60433), fieldName: 'earthAndMoon' }),
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
		config: AllStatsBuff,
		picker: MultiIconPicker,
		stats: [],
	},
	{
		config: ArmorBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: StaminaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina],
	},
	{
		config: StrengthAndAgilityBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStrength, Stat.StatAgility],
	},
	{
		config: ManaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMana],
	},
	{
		config: AttackPowerPercentBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: CritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatCritRating],
	},
	{
		config: MeleeHasteBuff,
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
		config: DamagePercentBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
	},
	{
		config: ResistanceBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatNatureResistance, Stat.StatShadowResistance, Stat.StatFrostResistance],
	},
	{
		config: MP5Buff,
		picker: MultiIconPicker,
		stats: [Stat.StatMP5],
	},
	{
		config: ReplenishmentBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMP5],
	},
	{
		config: Bloodlust,
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
		config: DarkIntent,
		picker: IconPicker,
		stats: [Stat.StatHasteRating],
	},
	{
		config: FocusMagic,
		picker: IconPicker,
		stats: [Stat.StatIntellect],
	},
	{
		config: RetributionAura,
		picker: IconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: ManaTideTotem,
		picker: IconPicker,
		stats: [Stat.StatMP5],
	},
	{
		config: Innervate,
		picker: IconPicker,
		stats: [Stat.StatMP5],
	},
	{
		config: PowerInfusion,
		picker: IconPicker,
		stats: [Stat.StatMP5, Stat.StatSpellPower],
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
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: PhysicalDamageDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: BleedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: SpellDamageDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellPower],
	},
	{
		config: SpellCritDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatIntellect],
	},
	{
		config: DamageReduction,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: MeleeAttackSpeedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
] as PickerStatOptions[];

export const DEBUFFS_MISC_CONFIG = [] as IconPickerStatOption[];
