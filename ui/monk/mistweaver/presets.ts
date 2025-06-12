import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { MistweaverMonk_Options as MistweaverMonkOptions, MonkMajorGlyph, MonkMinorGlyph, MonkStance } from '../../core/proto/monk';
import { SavedTalents } from '../../core/proto/ui';
import { Stats, UnitStat, UnitStatPresets } from '../../core/proto_utils/stats';
import DefaultGear from './gear_sets/default.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PREBIS_GEAR_PRESET = PresetUtils.makePresetGear('Default', DefaultGear);

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatIntellect]: 1.0,
			[Stat.StatSpirit]: 0.9,
			[Stat.StatSpellPower]: 0.79,
			[Stat.StatHitRating]: 0.9,
			[Stat.StatCritRating]: 0.42,
			[Stat.StatHasteRating]: 1.0,
			[Stat.StatMasteryRating]: 0.13,
		}
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
	flaskId: 76093, // Flask of the Winds
	foodId: 62290, // Seafood Magnifique Feast
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};

export const MISTWEAVER_BREAKPOINTS: UnitStatPresets[] = [
	{
		unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
		presets: new Map([
			['10-tick - ReM', 5.56876],
			['7-tick - EvM', 8.28372],
			['11-tick - ReM', 16.65209],
			['8-tick - EvM', 24.92194],
			['12-tick - ReM', 27.75472],
			['13-tick - ReM', 38.93714],
			['9-tick - EvM', 41.74346],
			['14-tick - ReM', 49.98126],
			['10-tick - EvM', 58.35315],
			['15-tick - ReM', 61.09546],
			['16-tick - ReM', 72.19115],
			['11-tick - EvM', 74.97816],
			['17-tick - ReM', 83.40213],
			['12-tick - EvM', 91.75459],
			['18-tick - ReM', 94.45797],
			['13-tick - EvM', 108.55062],
			['14-tick - EvM', 124.97193],
			['15-tick - EvM', 141.83803],
		]),
	},
];
