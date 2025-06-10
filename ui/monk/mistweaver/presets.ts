import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { MistweaverMonk_Options as MistweaverMonkOptions, MonkMajorGlyph, MonkMinorGlyph, MonkStance } from '../../core/proto/monk';
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
		talentsString: '',
		glyphs: Glyphs.create({
			major1: MonkMajorGlyph.GlyphOfSpinningCraneKick,
			major2: MonkMajorGlyph.GlyphOfTouchOfKarma,
			major3: MonkMajorGlyph.GlyphOfZenMeditation,
			minor1: MonkMinorGlyph.GlyphOfBlackoutKick,
			minor2: MonkMinorGlyph.GlyphOfJab,
			minor3: MonkMinorGlyph.GlyphOfWaterRoll,
		}),
	}),
};

export const DefaultOptions = MistweaverMonkOptions.create({
	classOptions: {},
	stance: MonkStance.WiseSerpent,
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58087, // Flask of the Winds
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};
