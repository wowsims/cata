import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Spec, Stat, TinkerHands } from '../../core/proto/common';
import {
	DruidMajorGlyph,
	DruidMinorGlyph,
	DruidPrimeGlyph,
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Rotation_AplType,
	FeralDruid_Rotation_BiteModeType,
} from '../../core/proto/druid';
import { SavedTalents } from '../../core/proto/ui';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4', P4Gear);

import DefaultApl from './apls/default.apl.json';
export const APL_ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);

import { Stats } from '../../core/proto_utils/stats';
import AoeApl from './apls/aoe.apl.json';
export const APL_ROTATION_AOE = PresetUtils.makePresetAPLRotation('APL AoE', AoeApl);

// Preset options for EP weights
export const BEARWEAVE_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Bear-Weave',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 0.38,
			[Stat.StatAgility]: 1.0,
			[Stat.StatAttackPower]: 0.37,
			[Stat.StatHitRating]: 0.36,
			[Stat.StatExpertiseRating]: 0.34,
			[Stat.StatCritRating]: 0.32,
			[Stat.StatHasteRating]: 0.30,
			[Stat.StatMasteryRating]: 0.33,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.54,
		},
	),
);

export const MONOCAT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Mono-Cat',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 0.39,
			[Stat.StatAgility]: 1.0,
			[Stat.StatAttackPower]: 0.37,
			[Stat.StatHitRating]: 0.31,
			[Stat.StatExpertiseRating]: 0.31,
			[Stat.StatCritRating]: 0.31,
			[Stat.StatHasteRating]: 0.30,
			[Stat.StatMasteryRating]: 0.33,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.56,
		},
	),
);

export const DefaultRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.SingleTarget,
	bearWeave: true,
	minCombosForRip: 5,
	minCombosForBite: 5,
	useRake: true,
	useBite: true,
	mangleSpam: false,
	biteModeType: FeralDruid_Rotation_BiteModeType.Emperical,
	biteTime: 11.0,
	berserkBiteTime: 6.0,
	minRoarOffset: 31.0,
	ripLeeway: 1.0,
	maintainFaerieFire: true,
	snekWeave: true,
	manualParams: false,
	biteDuringExecute: true,
	allowAoeBerserk: false,
	meleeWeave: true,
	cancelPrimalMadness: false,
});

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('Single Target Default', Spec.SpecFeralDruid, DefaultRotation);

export const AoeRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.Aoe,
	bearWeave: true,
	maintainFaerieFire: false,
	snekWeave: true,
	allowAoeBerserk: false,
	cancelPrimalMadness: false,
});

export const AOE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('AoE Default', Spec.SpecFeralDruid, AoeRotation);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Mono-Cat',
	data: SavedTalents.create({
		talentsString: '-2320322312012121202301-020301',
		glyphs: Glyphs.create({
			prime1: DruidPrimeGlyph.GlyphOfRip,
			prime2: DruidPrimeGlyph.GlyphOfBloodletting,
			prime3: DruidPrimeGlyph.GlyphOfBerserk,
			major1: DruidMajorGlyph.GlyphOfThorns,
			major2: DruidMajorGlyph.GlyphOfFeralCharge,
			major3: DruidMajorGlyph.GlyphOfRebirth,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfMarkOfTheWild,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const HybridTalents = {
	name: 'Hybrid',
	data: SavedTalents.create({
		talentsString: '-2300322312310001220311-020331',
		glyphs: Glyphs.create({
			prime1: DruidPrimeGlyph.GlyphOfRip,
			prime2: DruidPrimeGlyph.GlyphOfBloodletting,
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

export const DefaultOptions = FeralDruidOptions.create({
	assumeBleedActive: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodSkeweredEel,
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
	explosiveBigDaddy: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	highHpThreshold: 0.8,
	iterationCount: 25000,
	profession1: Profession.Engineering,
	profession2: Profession.ProfessionUnknown,
};
