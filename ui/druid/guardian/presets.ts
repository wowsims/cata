import * as Mechanics from '../../core/constants/mechanics.js';
import * as PresetUtils from '../../core/preset_utils.js';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions, PseudoStat, Spec, Stat, TinkerHands } from '../../core/proto/common';
import {
	DruidMajorGlyph,
	DruidMinorGlyph,
	DruidPrimeGlyph,
	GuardianDruid_Options as DruidOptions,
	GuardianDruid_Rotation as DruidRotation,
} from '../../core/proto/druid.js';
import { SavedTalents } from '../../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid Preset', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4', P4Gear);

export const DefaultSimpleRotation = DruidRotation.create({
	maintainFaerieFire: true,
	maintainDemoralizingRoar: true,
	demoTime: 4.0,
	pulverizeTime: 4.0,
	prepullStampede: true,
});

import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecGuardianDruid, DefaultSimpleRotation);

// Preset options for EP weights
export const SURVIVAL_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Survival',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.04,
			[Stat.StatStamina]: 0.99,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 1.02,
			[Stat.StatBonusArmor]: 0.23,
			[Stat.StatDodgeRating]: 0.97,
			[Stat.StatMasteryRating]: 0.35,
			[Stat.StatStrength]: 0.11,
			[Stat.StatAttackPower]: 0.1,
			[Stat.StatHitRating]: (0.075 + 0.195),
			[Stat.StatExpertiseRating]: 0.15,
			[Stat.StatCritRating]: 0.11,
			[Stat.StatHasteRating]: 0.0,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.0,
			[PseudoStat.PseudoStatPhysicalHitPercent]: (0.075 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT),
			[PseudoStat.PseudoStatSpellHitPercent]: (0.195 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT),
		},
	),
);

export const BALANCED_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Balanced',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.02,
			[Stat.StatStamina]: 0.76,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 0.62,
			[Stat.StatBonusArmor]: 0.14,
			[Stat.StatDodgeRating]: 0.59,
			[Stat.StatMasteryRating]: 0.20,
			[Stat.StatStrength]: 0.21,
			[Stat.StatAttackPower]: 0.20,
			[Stat.StatHitRating]: 0.60,
			[Stat.StatExpertiseRating]: 0.93,
			[Stat.StatCritRating]: 0.25,
			[Stat.StatHasteRating]: 0.03,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.23,
			[PseudoStat.PseudoStatPhysicalHitPercent]: (0.60 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT),
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-2300322312310001220311-020331',
		glyphs: Glyphs.create({
			prime1: DruidPrimeGlyph.GlyphOfMangle,
			prime2: DruidPrimeGlyph.GlyphOfLacerate,
			prime3: DruidPrimeGlyph.GlyphOfBerserk,
			major1: DruidMajorGlyph.GlyphOfFrenziedRegeneration,
			major2: DruidMajorGlyph.GlyphOfMaul,
			major3: DruidMajorGlyph.GlyphOfRebirth,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfChallengingRoar,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultOptions = DruidOptions.create({
	startingRage: 15,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfSteelskin,
	food: Food.FoodSkeweredEel,
	prepopPotion: Potions.PotionOfTheTolvir,
	defaultPotion: Potions.PotionOfTheTolvir,
	defaultConjured: Conjured.ConjuredHealthstone,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});
