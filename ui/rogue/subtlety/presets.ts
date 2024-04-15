import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { RogueMajorGlyph, RogueOptions_PoisonImbue, RoguePrimeGlyph,SubtletyRogue_Options as RogueOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import SubtletyApl from './apls/subtlety.apl.json'
import P1HemoSubGear from './gear_sets/p1_hemosub.gear.json';
import P2HemoSubGear from './gear_sets/p2_hemosub.gear.json';
import P3DanceSubGear from './gear_sets/p3_dancesub.gear.json';
import P3HemoSubGear from './gear_sets/p3_hemosub.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P1 Hemo Sub', P1HemoSubGear, { talentTree: 1 });
export const P2_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P2 Hemo Sub', P2HemoSubGear, { talentTree: 1 });
export const P3_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P3 Hemo Sub', P3HemoSubGear, { talentTree: 1 });
export const P3_PRESET_DANCE_SUB = PresetUtils.makePresetGear('P3 Dance Sub', P3DanceSubGear, { talentTree: 1 });

export const ROTATION_PRESET_SUBTLETY = PresetUtils.makePresetAPLRotation('Subtlety', SubtletyApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SubtletyTalents = {
	name: 'Subtlety',
	data: SavedTalents.create({
		talentsString: '023003-002-0332031321310012321',
		glyphs: Glyphs.create({
			prime1: RoguePrimeGlyph.GlyphOfBackstab,
			prime2: RoguePrimeGlyph.GlyphOfHemorrhage,
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
	defaultPotion: Potions.PotionOfSpeed,
	prepopPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredUnknown,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodMegaMammothMeal,
});
