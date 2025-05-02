import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, PseudoStat, Stat } from '../../core/proto/common';
import { CombatRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonImbue } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import CombatApl from './apls/combat.apl.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P3CombatGear from './gear_sets/p3_combat.gear.json';
import P4CombatGear from './gear_sets/p4_combat.gear.json';
import PreraidCombatGear from './gear_sets/preraid_combat.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_COMBAT = PresetUtils.makePresetGear('P1 Combat', P1CombatGear);
export const P3_PRESET_COMBAT = PresetUtils.makePresetGear('P3 Combat', P3CombatGear);
export const PRERAID_PRESET_COMBAT = PresetUtils.makePresetGear('Pre-Raid Combat', PreraidCombatGear);
export const P4_PRESET_COMBAT = PresetUtils.makePresetGear('P4 Combat', P4CombatGear);

export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('Combat', CombatApl);

// Preset options for EP weights
export const CBAT_STANDARD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Combat Standard',
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

// 4PT12 pushes Haste, Mastery, and Crit up moderately (Crit also gains from 2P but has no affect on reforging); Haste and Mastery overtake Hit for reforging entirely (Trends towards 10%-ish)
export const CBAT_4PT12_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Combat 4PT12',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.85,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.16,
			[Stat.StatHitRating]: 2.21,
			[Stat.StatHasteRating]: 1.38,
			[Stat.StatMasteryRating]: 1.28,
			[Stat.StatExpertiseRating]: 2.1,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.31,
			[PseudoStat.PseudoStatOffHandDps]: 1.32,
			[PseudoStat.PseudoStatSpellHitPercent]: 46,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 230,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const CombatTalents = {
	name: 'Combat',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfBladeFlurry,
			major3: RogueMajorGlyph.GlyphOfGouge,
		}),
	}),
};

export const DefaultOptions = RogueOptions.create({
	classOptions: {
		mhImbue: RogueOptions_PoisonImbue.InstantPoison,
		ohImbue: RogueOptions_PoisonImbue.DeadlyPoison,
		thImbue: RogueOptions_PoisonImbue.WoundPoison,
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
	conjuredId: 7676, // Thistle Tea
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
