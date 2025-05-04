import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, PseudoStat, Stat } from '../../core/proto/common';
import { RogueMajorGlyph, RogueOptions_PoisonOptions, SubtletyRogue_Options as RogueOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import SubtletyApl from './apls/subtlety.apl.json';
import P1SubtletyGear from './gear_sets/p1_subtlety.gear.json';
import PreraidSubtletyGear from './gear_sets/preraid_subtlety.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET_SUB = PresetUtils.makePresetGear('Pre-Raid Sub', PreraidSubtletyGear);
export const P1_PRESET_SUB = PresetUtils.makePresetGear('P1 Sub', P1SubtletyGear);

export const ROTATION_PRESET_SUBTLETY = PresetUtils.makePresetAPLRotation('Subtlety', SubtletyApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Sub',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 3.84,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.42,
			[Stat.StatHitRating]: 2.19,
			[Stat.StatHasteRating]: 1.58,
			[Stat.StatMasteryRating]: 0.95,
			[Stat.StatExpertiseRating]: 1.76,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 7.16,
			[PseudoStat.PseudoStatOffHandDps]: 1.07,
			[PseudoStat.PseudoStatSpellHitPercent]: 39.59,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 216.76,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SubtletyTalents = {
	name: 'Subtlety',
	data: SavedTalents.create({
		talentsString: '300003',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfHemorraghingVeins
		}),
	}),
};

export const DefaultOptions = RogueOptions.create({
	classOptions: {
		lethalPoison: RogueOptions_PoisonOptions.DeadlyPoison,
		applyPoisonsManually: false,
		startingOverkillDuration: 20,
		vanishBreakTime: 0.1,
	},
	honorAmongThievesCritRate: 400,
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58087, // Flask of the Winds
	foodId: 62669, // Skewered Eel
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
