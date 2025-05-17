import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { MonkMajorGlyph, MonkMinorGlyph, MonkOptions } from '../../core/proto/monk';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import DefaultP1Prebis2HGear from './gear_sets/p1_prebis_2h.gear.json';
import DefaultP1PrebisDWGear from './gear_sets/p1_prebis_dw.gear.json';

export const P1_PREBIS_2H_GEAR_PRESET = PresetUtils.makePresetGear('2H', DefaultP1Prebis2HGear);
export const P1_PREBIS_DW_GEAR_PRESET = PresetUtils.makePresetGear('DW', DefaultP1PrebisDWGear);

export const PREPATCH_ROTATION_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_PREBIS_2H_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Prebis - 2H',
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
export const P1_PREBIS_DW_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Prebis - DW',
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
			minor1: MonkMinorGlyph.GlyphOfBlackoutKick,
		}),
	}),
};

export const DefaultOptions = MonkOptions.create({
	classOptions: {},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId:  76084,  // Flask of Spring Blossoms
	foodId:   104303, // Sea Mist Rice Noodles
	potId:    76089,  // Virmen's Bite
	prepotId: 76089,  // Virmen's Bite
	tinkerId: 126734, // Synapse Springs II
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	distanceFromTarget: 5,
	iterationCount: 25000,
};
