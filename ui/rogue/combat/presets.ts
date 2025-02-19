import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions, PseudoStat, Stat } from '../../core/proto/common';
import { CombatRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonImbue, RoguePrimeGlyph } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import CombatApl from './apls/combat.apl.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P3CombatGear from './gear_sets/p3_combat.gear.json';
import P4CombatGear from './gear_sets/p4_combat.gear.json';
import PreraidCombatGear from './gear_sets/preraid_combat.gear.json'

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
			[Stat.StatCritRating]: 0.9,
			[Stat.StatHitRating]: 2.21,
			[Stat.StatHasteRating]: 1.36,
			[Stat.StatMasteryRating]: 1.33,
			[Stat.StatExpertiseRating]: 1.74,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.31,
			[PseudoStat.PseudoStatOffHandDps]: 1.32,
			[PseudoStat.PseudoStatSpellHitPercent]: 46,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 220,
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
			[Stat.StatCritRating]: 1.09,
			[Stat.StatHitRating]: 2.21,
			[Stat.StatHasteRating]: 1.52,
			[Stat.StatMasteryRating]: 1.41,
			[Stat.StatExpertiseRating]: 1.74,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.31,
			[PseudoStat.PseudoStatOffHandDps]: 1.32,
			[PseudoStat.PseudoStatSpellHitPercent]: 46,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 230,
		},
	),
);

// Legendary daggers push up Haste significantly. It effectively removes the need to balance Haste vs Mastery.
export const CBAT_T13_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Combat Legendary',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.85,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.19,
			[Stat.StatHitRating]: 2.5,
			[Stat.StatHasteRating]: 1.79,
			[Stat.StatMasteryRating]: 1.48,
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


// No'Kaled MH MASSIVELY inflates Mastery's EP value as Main Gauche procs can proc weapon effects
export const CBAT_NOKALED_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Combat No\'Kaled',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.85,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.19,
			[Stat.StatHitRating]: 2.5,
			[Stat.StatHasteRating]: 1.76,
			[Stat.StatMasteryRating]: 1.78,
			[Stat.StatExpertiseRating]: 2.1,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.31,
			[PseudoStat.PseudoStatOffHandDps]: 1.32,
			[PseudoStat.PseudoStatSpellHitPercent]: 100, // Yeah this is a big number, idk Combat is weird this makes dps big
			[PseudoStat.PseudoStatPhysicalHitPercent]: 250,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const CombatTalents = {
	name: 'Combat',
	data: SavedTalents.create({
		talentsString: '0322-2332030310230012321-003',
		glyphs: Glyphs.create({
			prime1: RoguePrimeGlyph.GlyphOfAdrenalineRush,
			prime2: RoguePrimeGlyph.GlyphOfSinisterStrike,
			prime3: RoguePrimeGlyph.GlyphOfSliceAndDice,
			major1: RogueMajorGlyph.GlyphOfBladeFlurry,
			major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
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
