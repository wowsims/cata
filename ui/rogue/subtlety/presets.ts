import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, PseudoStat, Stat } from '../../core/proto/common';
import { RogueMajorGlyph, RogueOptions_PoisonOptions, SubtletyRogue_Options as RogueOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import SubtletyApl from './apls/subtlety.apl.json';
import MSVGear from './gear_sets/p1_subtlety_msv.gear.json';
import PreraidGear from './gear_sets/preraid_subtlety.gear.json';
import T14 from './gear_sets/p1_subtlety_t14.gear.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_GEARSET = PresetUtils.makePresetGear('P1 Preraid', PreraidGear);
export const P1_MSV_GEARSET = PresetUtils.makePresetGear('P1 MSV', MSVGear);
export const P1_T14_GEARSET = PresetUtils.makePresetGear('P1 T14', T14);

export const ROTATION_PRESET_SUBTLETY = PresetUtils.makePresetAPLRotation('Subtlety', SubtletyApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Sub',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 1.0,
			[Stat.StatCritRating]: 0.31,
			[Stat.StatHitRating]: 0.54,
			[Stat.StatHasteRating]: 0.32,
			[Stat.StatMasteryRating]: 0.26,
			[Stat.StatExpertiseRating]: 0.35,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.43,
			[PseudoStat.PseudoStatOffHandDps]: 0.26,
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
