import * as PresetUtils from '../../core/preset_utils.js';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common.js';
import {
	PaladinAura as PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinPrimeGlyph,
	PaladinSeal,
	RetributionPaladin_Options as RetributionPaladinOptions,
} from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import ApparatusApl from './apls/apparatus.apl.json';
import DefaultApl from './apls/default.apl.json';
import T13_2Pc_Apl from './apls/t13.apl.json';
import P2_BisRetGear from './gear_sets/p2_bis.gear.json';
import P3_BisRetGear from './gear_sets/p3_bis.gear.json';
import P4_BisRetGear from './gear_sets/p4_bis.gear.json';
import PreraidRetGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so it's good to
// keep them in a separate file.

export const PRERAID_RET_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidRetGear);
export const P2_BIS_RET_PRESET = PresetUtils.makePresetGear('P2', P2_BisRetGear);
export const P3_BIS_RET_PRESET = PresetUtils.makePresetGear('P3', P3_BisRetGear);
export const P4_BIS_RET_PRESET = PresetUtils.makePresetGear('P4', P4_BisRetGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const ROTATION_PRESET_APPARATUS = PresetUtils.makePresetAPLRotation('Apparatus', ApparatusApl);
export const ROTATION_PRESET_T13 = PresetUtils.makePresetAPLRotation('T13 2pc', T13_2Pc_Apl);

// Preset options for EP weights
export const P2_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P2',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.28,
			[Stat.StatAttackPower]: 1,
			[Stat.StatHitRating]: 2.26,
			[Stat.StatCritRating]: 1.10,
			[Stat.StatHasteRating]: 1.04,
			[Stat.StatExpertiseRating]: 1.82,
			[Stat.StatMasteryRating]: 1.26,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 8.28,
			[PseudoStat.PseudoStatSpellHitPercent]: 30.23,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 261.35,
		},
	),
);
export const P3_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P3',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.28,
			[Stat.StatAttackPower]: 1,
			[Stat.StatHitRating]: 2.66,
			[Stat.StatCritRating]: 1.24,
			[Stat.StatHasteRating]: 0.99,
			[Stat.StatExpertiseRating]: 2.25,
			[Stat.StatMasteryRating]: 1.40,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 8.77,
			[PseudoStat.PseudoStatSpellHitPercent]: 26.02,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 296.98,
		},
	),
);

export const P4_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P4',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.28,
			[Stat.StatAttackPower]: 1,
			[Stat.StatHitRating]: 2.93,
			[Stat.StatCritRating]: 1.47,
			[Stat.StatHasteRating]: 1.06,
			[Stat.StatExpertiseRating]: 2.43,
			[Stat.StatMasteryRating]: 1.76,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 8.39,
			[PseudoStat.PseudoStatSpellHitPercent]: 33.56,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 342.82,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const P2_Talents = {
	name: 'P2',
	data: SavedTalents.create({
		talentsString: '203002-02-23203213211113002311',
		glyphs: Glyphs.create({
			prime1: PaladinPrimeGlyph.GlyphOfSealOfTruth,
			prime2: PaladinPrimeGlyph.GlyphOfExorcism,
			prime3: PaladinPrimeGlyph.GlyphOfTemplarSVerdict,
			major1: PaladinMajorGlyph.GlyphOfTheAsceticCrusader,
			major2: PaladinMajorGlyph.GlyphOfHammerOfWrath,
			major3: PaladinMajorGlyph.GlyphOfConsecration,
			minor1: PaladinMinorGlyph.GlyphOfRighteousness,
			minor2: PaladinMinorGlyph.GlyphOfTruth,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfMight,
		}),
	}),
};

export const P3_P4_Talents = {
	name: 'P3 / P4',
	data: SavedTalents.create({
		talentsString: '203002-02-23203213211113002311',
		glyphs: Glyphs.create({
			prime1: PaladinPrimeGlyph.GlyphOfSealOfTruth,
			prime2: PaladinPrimeGlyph.GlyphOfCrusaderStrike,
			prime3: PaladinPrimeGlyph.GlyphOfTemplarSVerdict,
			major1: PaladinMajorGlyph.GlyphOfTheAsceticCrusader,
			major2: PaladinMajorGlyph.GlyphOfHammerOfWrath,
			major3: PaladinMajorGlyph.GlyphOfConsecration,
			minor1: PaladinMinorGlyph.GlyphOfRighteousness,
			minor2: PaladinMinorGlyph.GlyphOfTruth,
			minor3: PaladinMinorGlyph.GlyphOfBlessingOfMight,
		}),
	}),
};

export const P2_PRESET = PresetUtils.makePresetBuild('P2', {
	gear: P2_BIS_RET_PRESET,
	epWeights: P2_EP_PRESET,
	talents: P2_Talents,
	rotationType: APLRotationType.TypeAuto,
})

export const P3_PRESET = PresetUtils.makePresetBuild('P3', {
	gear: P3_BIS_RET_PRESET,
	epWeights: P3_EP_PRESET,
	talents: P3_P4_Talents,
	rotationType: APLRotationType.TypeAuto,
})

export const P4_PRESET = PresetUtils.makePresetBuild('P4', {
	gear: P4_BIS_RET_PRESET,
	epWeights: P4_EP_PRESET,
	talents: P3_P4_Talents,
	rotationType: APLRotationType.TypeAuto,
})

export const DefaultOptions = RetributionPaladinOptions.create({
	classOptions: {
		aura: PaladinAura.Retribution,
		seal: PaladinSeal.Truth,
		snapshotGuardian: false,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTitanicStrength,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	distanceFromTarget: 5,
	duration: 300,
	durationVariation: 30,
	iterationCount: 20000,
};
