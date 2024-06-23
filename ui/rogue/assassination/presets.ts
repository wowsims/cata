import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions, PseudoStat, Stat } from '../../core/proto/common';
import { AssassinationRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonImbue, RoguePrimeGlyph } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import MutilateApl from './apls/mutilate.apl.json';
import P1AssassinationGear from './gear_sets/p1_assassination.gear.json';
import P1ExpertiseGear from './gear_sets/p1_expertise.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P1 Assassination', P1AssassinationGear);
export const P1_PRESET_ASN_EXPERTISE = PresetUtils.makePresetGear('P1 Expertise', P1ExpertiseGear);

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Assassination', MutilateApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.58,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatSpellCrit]: 0.26,
			[Stat.StatSpellHit]: 1.31,
			[Stat.StatMeleeHit]: 0.7,
			[Stat.StatMeleeCrit]: 0.62,
			[Stat.StatMeleeHaste]: 1.1,
			[Stat.StatMastery]: 1.23,
			[Stat.StatExpertise]: 1.04,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 2.49,
			[PseudoStat.PseudoStatOffHandDps]: 1.0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const AssassinationTalentsDefault = {
	name: 'Assassination 31/2/8',
	data: SavedTalents.create({
		talentsString: '0333230013122110321-002-203003',
		glyphs: Glyphs.create({
			prime1: RoguePrimeGlyph.GlyphOfMutilate,
			prime2: RoguePrimeGlyph.GlyphOfBackstab,
			prime3: RoguePrimeGlyph.GlyphOfRupture,
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
	duration: 240,
	durationVariation: 20,
};
