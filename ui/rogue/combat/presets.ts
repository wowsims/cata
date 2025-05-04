import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, PseudoStat, Stat } from '../../core/proto/common';
import { CombatRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import CombatApl from './apls/combat.apl.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import PreraidCombatGear from './gear_sets/preraid_combat.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET_COMBAT = PresetUtils.makePresetGear('Pre-Raid Combat', PreraidCombatGear);
export const P1_PRESET_COMBAT = PresetUtils.makePresetGear('P1 Combat', P1CombatGear);

export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('Combat', CombatApl);

// Preset options for EP weights
export const CBAT_STANDARD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Combat',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.85,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.2,
			[Stat.StatHitRating]: 2.5,
			[Stat.StatHasteRating]: 1.58,
			[Stat.StatMasteryRating]: 1.41,
			[Stat.StatExpertiseRating]: 2.1,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.31,
			[PseudoStat.PseudoStatOffHandDps]: 1.32,
			[PseudoStat.PseudoStatSpellHitPercent]: 52,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 249,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const CombatTalents = {
	name: 'Combat',
	data: SavedTalents.create({
		talentsString: '300003',
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
	flaskId: 58087, // Flask of the Winds
	foodId: 62669, // Skewered Eel
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
