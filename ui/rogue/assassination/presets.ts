import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { AssassinationRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonImbue, RoguePrimeGlyph } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import MutilateApl from './apls/mutilate.apl.json';
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
	defaultPotion: Potions.PotionOfSpeed,
	prepopPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredRogueThistleTea,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodMegaMammothMeal,
});
