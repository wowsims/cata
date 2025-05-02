import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
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

export const DefaultOptions = MonkOptions.create({
	classOptions: {},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58087, // Flask of the Winds
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
	tinkerId: 82174, // Synapse Springs
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};
