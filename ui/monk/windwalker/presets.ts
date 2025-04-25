import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common';
import { MonkMajorGlyph, MonkMinorGlyph, MonkOptions } from '../../core/proto/monk';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import Default2HGear from './gear_sets/2h.gear.json';
import DefaultDWGear from './gear_sets/dw.gear.json';

export const PREPATCH_2H_GEAR_PRESET = PresetUtils.makePresetGear('2H', Default2HGear);
export const PREPATCH_DW_GEAR_PRESET = PresetUtils.makePresetGear('DW', DefaultDWGear);

export const PREPATCH_ROTATION_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const PREPATCH_2H_EP_PRESET = PresetUtils.makePresetEpWeights(
	'2H',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.05,
			[Stat.StatAgility]: 2.92,
			[Stat.StatHitRating]: 3.0,
			[Stat.StatCritRating]: 1.28,
			[Stat.StatHasteRating]: 1.68,
			[Stat.StatExpertiseRating]: 2.99,
			[Stat.StatMasteryRating]: 0.68,
			[Stat.StatAttackPower]: 1,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 9.22,
			[PseudoStat.PseudoStatOffHandDps]: 0,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 360.07,
		},
	),
);

// Preset options for EP weights
export const PREPATCH_DW_EP_PRESET = PresetUtils.makePresetEpWeights(
	'DW',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.05,
			[Stat.StatAgility]: 2.94,
			[Stat.StatHitRating]: 3.21,
			[Stat.StatCritRating]: 1.34,
			[Stat.StatHasteRating]: 1.7,
			[Stat.StatExpertiseRating]: 3.2,
			[Stat.StatMasteryRating]: 0.68,
			[Stat.StatAttackPower]: 1,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 7.59,
			[PseudoStat.PseudoStatOffHandDps]: 3.8,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 385.2,
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

export const DefaultOptions = MonkOptions.create({
	classOptions: {},
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
