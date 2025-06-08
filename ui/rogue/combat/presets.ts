import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, PseudoStat, Stat } from '../../core/proto/common';
import { CombatRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import CombatApl from './apls/combat.apl.json';
import MSVGear from './gear_sets/p1_combat_msv.gear.json';
import PreraidGear from './gear_sets/preraid_combat.gear.json';
import T14 from './gear_sets/p1_combat_t14.gear.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_GEARSET = PresetUtils.makePresetGear('P1 Preraid', PreraidGear);
export const P1_MSV_GEARSET = PresetUtils.makePresetGear('P1 MSV', MSVGear);
export const P1_T14_GEARSET = PresetUtils.makePresetGear('P1 T14', T14);

export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('Combat', CombatApl);

// Preset options for EP weights
export const CBAT_STANDARD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Combat',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 1.0,
			[Stat.StatCritRating]: 0.29,
			[Stat.StatHitRating]: 0.7,
			[Stat.StatHasteRating]: 0.39,
			[Stat.StatMasteryRating]: 0.32,
			[Stat.StatExpertiseRating]: 0.56,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.63,
			[PseudoStat.PseudoStatOffHandDps]: 0.45,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const CombatTalents = {
	name: 'Combat',
	data: SavedTalents.create({
		talentsString: '321233',
		glyphs: Glyphs.create({
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
