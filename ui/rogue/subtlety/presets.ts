import * as Mechanics from '../../core/constants/mechanics';
import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions, PseudoStat, Stat } from '../../core/proto/common';
import { RogueMajorGlyph, RogueOptions_PoisonImbue, SubtletyRogue_Options as RogueOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import SubtletyApl from './apls/subtlety.apl.json';
import SubtletyMasteryApl from './apls/subtlety_mastery_prepull.apl.json';
import P1SubtletyGear from './gear_sets/p1_subtlety.gear.json';
import P3SubtletyGear from './gear_sets/p3_subtlety.gear.json';
import P4SubtletyGear from './gear_sets/p4_subtlety.gear.json';
import PreraidSubtletyGear from './gear_sets/preraid_subtlety.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_SUB = PresetUtils.makePresetGear('P1 Sub', P1SubtletyGear);
export const P3_PRESET_SUB = PresetUtils.makePresetGear('P3 Sub', P3SubtletyGear);
export const PRERAID_PRESET_SUB = PresetUtils.makePresetGear('Pre-Raid Sub', PreraidSubtletyGear);
export const P4_PRESET_SUB = PresetUtils.makePresetGear('P4 Sub', P4SubtletyGear);

export const ROTATION_PRESET_SUBTLETY = PresetUtils.makePresetAPLRotation('Subtlety', SubtletyApl);
export const ROTATION_PRESET_SUBTLETY_MASTERY = PresetUtils.makePresetAPLRotation('Subtlety Mastery Prepull', SubtletyMasteryApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Sub Default',
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

export const P4_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Sub Legendary Dags',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 3.84,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.42,
			[Stat.StatHitRating]: 2.19,
			[Stat.StatHasteRating]: 1.67,
			[Stat.StatMasteryRating]: 1.32,
			[Stat.StatExpertiseRating]: 1.76,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 7.16,
			[PseudoStat.PseudoStatOffHandDps]: 1.07,
			[PseudoStat.PseudoStatSpellHitPercent]: 49.09,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 216.76,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SubtletyTalents = {
	name: 'Subtlety',
	data: SavedTalents.create({
		talentsString: '023003-002-0332031321310012321',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
			major2: RogueMajorGlyph.GlyphOfSprint,
			major3: RogueMajorGlyph.GlyphOfFeint,
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
	honorAmongThievesCritRate: 400,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodSkeweredEel,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
