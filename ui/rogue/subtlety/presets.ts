import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { RogueMajorGlyph, RogueOptions_PoisonImbue, RoguePrimeGlyph,SubtletyRogue_Options as RogueOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import SubtletyApl from './apls/subtlety.apl.json'
import P1SubtletyGear from './gear_sets/p1_subtlety.gear.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_SUB = PresetUtils.makePresetGear('P1 Hemo Sub', P1SubtletyGear, { talentTree: 1 });

export const ROTATION_PRESET_SUBTLETY = PresetUtils.makePresetAPLRotation('Subtlety', SubtletyApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SubtletyTalents = {
	name: 'Subtlety',
	data: SavedTalents.create({
		talentsString: '023003-002-0332031321310012321',
		glyphs: Glyphs.create({
			prime1: RoguePrimeGlyph.GlyphOfBackstab,
			prime2: RoguePrimeGlyph.GlyphOfShadowDance,
			prime3: RoguePrimeGlyph.GlyphOfSliceAndDice,
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
