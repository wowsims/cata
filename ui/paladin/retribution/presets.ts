import * as Mechanics from '../../core/constants/mechanics';
import * as PresetUtils from '../../core/preset_utils.js';
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
import DefaultApl from './apls/default.apl.json';
import T13_2Pc_Apl from './apls/t13.apl.json';
//import T13_BisRetGear from './gear_sets/t13_bis.gear.json';
import PreraidRetGear from './gear_sets/preraid.gear.json';
import T11_BisRetGear from './gear_sets/t11_bis.gear.json';
import T12_BisRetGear from './gear_sets/t12_bis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so it's good to
// keep them in a separate file.

export const PRERAID_RET_PRESET = PresetUtils.makePresetGear('Preraid', PreraidRetGear);
export const T11_BIS_RET_PRESET = PresetUtils.makePresetGear('T11 BiS', T11_BisRetGear);
export const T12_BIS_RET_PRESET = PresetUtils.makePresetGear('T12 BiS', T12_BisRetGear);
//export const T13_BIS_RET_PRESET = PresetUtils.makePresetGear('T13 BiS', T13_BisRetGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const ROTATION_PRESET_T13 = PresetUtils.makePresetAPLRotation('T13 2pc', T13_2Pc_Apl);

// Preset options for EP weights
export const T11_EP_PRESET = PresetUtils.makePresetEpWeights(
	'T11',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.32,
			[Stat.StatAttackPower]: 1,
			[Stat.StatHitRating]: 2.34,
			[Stat.StatCritRating]: 1.15,
			[Stat.StatHasteRating]: 0.94,
			[Stat.StatExpertiseRating]: 1.94,
			[Stat.StatMasteryRating]: 1.37,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 8.46,
			[PseudoStat.PseudoStatSpellHitPercent]: 29.11,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 261.3,
		},
	),
);
export const T12_EP_PRESET = PresetUtils.makePresetEpWeights(
	'T12',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.32,
			[Stat.StatAttackPower]: 1,
			[Stat.StatHitRating]: 2.68,
			[Stat.StatCritRating]: 1.33,
			[Stat.StatHasteRating]: 1.1,
			[Stat.StatExpertiseRating]: 2.35,
			[Stat.StatMasteryRating]: 1.53,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 8.97,
			[PseudoStat.PseudoStatSpellHitPercent]: 26.30,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 313.86,
		},
	),
);
/*
// Preset options for EP weights
export const T13_EP_PRESET = PresetUtils.makePresetEpWeights(
	'T13',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.28,
			[Stat.StatAttackPower]: 1,
			[Stat.StatSpellHit]: 0.18,
			[Stat.StatSpellCrit]: 0.15,
			[Stat.StatSpellHaste]: 0.35,
			[Stat.StatMeleeHit]: 2.01,
			[Stat.StatMeleeCrit]: 1.26,
			[Stat.StatMeleeHaste]: 0.50,
			[Stat.StatExpertise]: 1.83,
			[Stat.StatMastery]: 1.74,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 8.29,
		},
	),
);
*/
// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const T11Talents = {
	name: 'T11',
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

export const T12T13Talents = {
	name: 'T12/T13',
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
};
