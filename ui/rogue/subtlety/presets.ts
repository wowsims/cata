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
			[Stat.StatAgility]: 3.5,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 0.95,
			[Stat.StatHitRating]: 1.6,
			[Stat.StatHasteRating]: 1.25,
			[Stat.StatMasteryRating]: 0.85,
			[Stat.StatExpertiseRating]: 1.5,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 5.2,
			[PseudoStat.PseudoStatOffHandDps]: 0.95,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 540,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SubtletyTalents = {
	name: 'Subtlety',
	data: SavedTalents.create({
		talentsString: '321233',
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
	flaskId: 76084, // Flask of the Winds
	foodId: 74648, // Skewered Eel
	potId: 76089, // Potion of the Tol'vir
	prepotId: 76089, // Potion of the Tol'vir
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
