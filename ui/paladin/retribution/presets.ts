import * as PresetUtils from '../../core/preset_utils.js';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common.js';
import { PaladinMajorGlyph, PaladinSeal, RetributionPaladin_Options as RetributionPaladinOptions } from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1_Gear from './gear_sets/p1.gear.json';
import Preraid_Gear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so it's good to
// keep them in a separate file.

export const P1_GEAR_PRESET = PresetUtils.makePresetGear('P1', P1_Gear);
export const PRERAID_GEAR_PRESET = PresetUtils.makePresetGear('Pre-raid', Preraid_Gear);

export const APL_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.0,
			[Stat.StatHitRating]: 1.0,
			[Stat.StatExpertiseRating]: 0.87,
			[Stat.StatHasteRating]: 0.52,
			[Stat.StatMasteryRating]: 0.51,
			[Stat.StatCritRating]: 0.5,
			[Stat.StatAttackPower]: 0.44,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.91,
		},
	),
);

export const PRERAID_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Pre-raid',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.0,
			[Stat.StatHitRating]: 0.72,
			[Stat.StatExpertiseRating]: 0.63,
			[Stat.StatHasteRating]: 0.56,
			[Stat.StatAttackPower]: 0.44,
			[Stat.StatMasteryRating]: 0.41,
			[Stat.StatCritRating]: 0.38,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.77,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '221223',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfTemplarsVerdict,
			major2: PaladinMajorGlyph.GlyphOfDoubleJeopardy,
			major3: PaladinMajorGlyph.GlyphOfMassExorcism,
		}),
	}),
};

export const P1_BUILD_PRESET = PresetUtils.makePresetBuild('P1', {
	gear: P1_GEAR_PRESET,
	epWeights: P1_EP_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
});

export const PRERAID_BUILD_PRESET = PresetUtils.makePresetBuild('Pre-raid', {
	gear: PRERAID_GEAR_PRESET,
	epWeights: PRERAID_EP_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
});

export const DefaultOptions = RetributionPaladinOptions.create({
	classOptions: {
		seal: PaladinSeal.Truth,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76088, // Flask of Winter's Bite
	foodId: 74646, // Black Pepper Ribs and Shrimp
	potId: 76095, // Potion of Mogu Power
	prepotId: 76095, // Potion of Mogu Power
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};
