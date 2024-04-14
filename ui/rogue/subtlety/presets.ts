import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { RogueMajorGlyph, RogueOptions_PoisonImbue, SubtletyRogue_Options as RogueOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import P1HemoSubGear from './gear_sets/p1_hemosub.gear.json';
import P2HemoSubGear from './gear_sets/p2_hemosub.gear.json';
import P3DanceSubGear from './gear_sets/p3_dancesub.gear.json';
import P3HemoSubGear from './gear_sets/p3_hemosub.gear.json';
import SubtletyApl from './apls/subtlety.apl.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P1 Hemo Sub', P1HemoSubGear, { talentTree: 2 });
export const P2_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P2 Hemo Sub', P2HemoSubGear, { talentTree: 2 });
export const P3_PRESET_HEMO_SUB = PresetUtils.makePresetGear('P3 Hemo Sub', P3HemoSubGear, { talentTree: 2 });
export const P3_PRESET_DANCE_SUB = PresetUtils.makePresetGear('P3 Dance Sub', P3DanceSubGear, { talentTree: 2 });

export const ROTATION_PRESET_SUBTLETY = PresetUtils.makePresetAPLRotation('Subtlety', SubtletyApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SubtletyTalents = {
	name: 'Subtlety',
	data: SavedTalents.create({
		// talentsString: '30532010114--5022012030321121350115031151',
		// glyphs: Glyphs.create({
		// 	major1: RogueMajorGlyph.GlyphOfEviscerate,
		// 	major2: RogueMajorGlyph.GlyphOfRupture,
		// 	major3: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		// }),
	}),
};

export const HemoSubtletyTalents = {
	name: 'Hemo Sub',
	data: SavedTalents.create({
		// talentsString: '30532010135--502201203032112135011503122',
		// glyphs: Glyphs.create({
		// 	major1: RogueMajorGlyph.GlyphOfEviscerate,
		// 	major2: RogueMajorGlyph.GlyphOfRupture,
		// 	major3: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		// }),
	}),
};

export const DefaultOptions = RogueOptions.create({
	classOptions: {
		mhImbue: RogueOptions_PoisonImbue.DeadlyPoison,
		ohImbue: RogueOptions_PoisonImbue.InstantPoison,
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
