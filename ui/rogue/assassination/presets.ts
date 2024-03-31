import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { AssassinationRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonImbue } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import FanAoeApl from './apls/fan_aoe.apl.json';
import MutilateApl from './apls/mutilate.apl.json';
import MutilateExposeApl from './apls/mutilate_expose.apl.json';
import RuptureMutilateApl from './apls/rupture_mutilate.apl.json';
import RuptureMutilateExposeApl from './apls/rupture_mutilate_expose.apl.json';
import P1AssassinationGear from './gear_sets/p1_assassination.gear.json';
import P2AssassinationGear from './gear_sets/p2_assassination.gear.json';
import P3AssassinationGear from './gear_sets/p3_assassination.gear.json';
import P4AssassinationGear from './gear_sets/p4_assassination.gear.json';
import P5AssassinationGear from './gear_sets/p5_assassination.gear.json';
import PreraidAssassinationGear from './gear_sets/preraid_assassination.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET_ASSASSINATION = PresetUtils.makePresetGear('PreRaid Assassination', PreraidAssassinationGear);
export const P1_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P1 Assassination', P1AssassinationGear);
export const P2_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P2 Assassination', P2AssassinationGear);
export const P3_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P3 Assassination', P3AssassinationGear);
export const P4_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P4 Assassination', P4AssassinationGear);
export const P5_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P5 Assassination', P5AssassinationGear);

export const ROTATION_PRESET_MUTILATE = PresetUtils.makePresetAPLRotation('Mutilate', MutilateApl);
export const ROTATION_PRESET_RUPTURE_MUTILATE = PresetUtils.makePresetAPLRotation('Rupture Mutilate', RuptureMutilateApl);
export const ROTATION_PRESET_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('Mutilate w/ Expose', MutilateExposeApl);
export const ROTATION_PRESET_RUPTURE_MUTILATE_EXPOSE = PresetUtils.makePresetAPLRotation('Rupture Mutilate w/ Expose', RuptureMutilateExposeApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Fan AOE', FanAoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const AssassinationTalents137 = {
	name: 'Assassination 13/7',
	data: SavedTalents.create({
		// talentsString: '005303104352100520103331051-005005003-502',
		// glyphs: Glyphs.create({
		// 	major1: RogueMajorGlyph.GlyphOfMutilate,
		// 	major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		// 	major3: RogueMajorGlyph.GlyphOfHungerForBlood,
		// }),
	}),
};

export const AssassinationTalents182 = {
	name: 'Assassination 18/2',
	data: SavedTalents.create({
		// talentsString: '005303104352100520103331051-005005005003-2',
		// glyphs: Glyphs.create({
		// 	major1: RogueMajorGlyph.GlyphOfMutilate,
		// 	major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		// 	major3: RogueMajorGlyph.GlyphOfHungerForBlood,
		// }),
	}),
};

export const AssassinationTalentsBF = {
	name: 'Assassination Blade Flurry',
	data: SavedTalents.create({
		// talentsString: '005303104352100520103231-005205005003001-501',
		// glyphs: Glyphs.create({
		// 	major1: RogueMajorGlyph.GlyphOfMutilate,
		// 	major2: RogueMajorGlyph.GlyphOfTricksOfTheTrade,
		// 	major3: RogueMajorGlyph.GlyphOfBladeFlurry,
		// }),
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
	defaultPotion: Potions.PotionOfSpeed,
	prepopPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodMegaMammothMeal,
});
