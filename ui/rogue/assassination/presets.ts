import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions, PseudoStat, Stat } from '../../core/proto/common';
import { AssassinationRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonImbue } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import MutilateApl from './apls/mutilate.apl.json';
import P1AssassinationGear from './gear_sets/p1_assassination.gear.json';
import P1ExpertiseGear from './gear_sets/p1_expertise.gear.json';
import P3AssassinationGear from './gear_sets/p3_assassination.gear.json';
import P4AssassinationGear from './gear_sets/p4_assassination.gear.json';
import PreraidAssassination from './gear_sets/preraid_assassination.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P1 Assassination', P1AssassinationGear);
export const P1_PRESET_ASN_EXPERTISE = PresetUtils.makePresetGear('P1 Expertise', P1ExpertiseGear);
export const P3_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P3 Assassination', P3AssassinationGear);
export const PRERAID_PRESET_ASSASSINATION = PresetUtils.makePresetGear('Pre-Raid Assassination', PreraidAssassination);
export const P4_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P4 Assassination', P4AssassinationGear);

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Assassination', MutilateApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Asn',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.64,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.18,
			[Stat.StatHitRating]: 2.62,
			[Stat.StatHasteRating]: 1.35,
			[Stat.StatMasteryRating]: 1.45,
			[Stat.StatExpertiseRating]: 1.2,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.0,
			[PseudoStat.PseudoStatOffHandDps]: 0.97,
			[PseudoStat.PseudoStatSpellHitPercent]: 130.5,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 162.0,
		},
	),
);

export const P1_EP_EXPERTISE_PRESET = PresetUtils.makePresetEpWeights(
	'Asn Expertise',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.71,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.18,
			[Stat.StatHitRating]: 2.62,
			[Stat.StatHasteRating]: 1.35,
			[Stat.StatMasteryRating]: 1.45,
			[Stat.StatExpertiseRating]: 2.0,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.0,
			[PseudoStat.PseudoStatOffHandDps]: 0.97,
			[PseudoStat.PseudoStatSpellHitPercent]: 130.5,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 162.0,
		},
	),
);

export const P4_EP_LEGENDARY_PRESET = PresetUtils.makePresetEpWeights(
	'Asn Legendary',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.71,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.18,
			[Stat.StatHitRating]: 2.62,
			[Stat.StatHasteRating]: 1.39,
			[Stat.StatMasteryRating]: 1.61,
			[Stat.StatExpertiseRating]: 1.22,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.0,
			[PseudoStat.PseudoStatOffHandDps]: 0.97,
			[PseudoStat.PseudoStatSpellHitPercent]: 130.5,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 162.0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.

export const AssassinationTalentsDefault = {
	name: 'Assassination 31/2/8',
	data: SavedTalents.create({
		talentsString: '0333230013122110321-002-203003',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfFeint,
			major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
			major3: RogueMajorGlyph.GlyphOfSprint,
		}),
	}),
};

export const DefaultOptions = RogueOptions.create({
	classOptions: {
		mhImbue: RogueOptions_PoisonImbue.DeadlyPoison,
		ohImbue: RogueOptions_PoisonImbue.InstantPoison,
		thImbue: RogueOptions_PoisonImbue.DeadlyPoison,
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
