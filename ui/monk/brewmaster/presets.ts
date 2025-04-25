import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common';
import { BrewmasterMonk_Options as BrewmasterMonkOptions, MonkMajorGlyph, MonkMinorGlyph, MonkStance } from '../../core/proto/monk';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import DefaultGear from './gear_sets/default.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PREPATCH_GEAR_PRESET = PresetUtils.makePresetGear('Default', DefaultGear);

// Preset options for EP weights
export const PREPATCH_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.85,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 0.9,
			[Stat.StatHitRating]: 2.21,
			[Stat.StatHasteRating]: 1.36,
			[Stat.StatMasteryRating]: 1.33,
			[Stat.StatExpertiseRating]: 1.74,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.31,
			[PseudoStat.PseudoStatOffHandDps]: 1.32,
			[PseudoStat.PseudoStatSpellHitPercent]: 46,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 220,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '00110000100100101',
		glyphs: Glyphs.create({
			major1: MonkMajorGlyph.MonkMajorGlyphSpinningCraneKick,
			major2: MonkMajorGlyph.MonkMajorGlyphTouchOfKarma,
			major3: MonkMajorGlyph.MonkMajorGlyphZenMeditation,
			minor1: MonkMinorGlyph.MonkMinorGlyphBlackoutKick,
			minor2: MonkMinorGlyph.MonkMinorGlyphJab,
			minor3: MonkMinorGlyph.MonkMinorGlyphWaterRoll,
		}),
	}),
};

export const DefaultOptions = BrewmasterMonkOptions.create({
	classOptions: {},
	stance: MonkStance.SturdyOx,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodSeafoodFeast,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};