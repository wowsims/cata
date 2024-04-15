import { Rogue } from '../../core/player_classes/rogue';
import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { CombatRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonImbue, RoguePrimeGlyph } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import CombatApl from './apls/combat.apl.json';
import P1CombatGear from './gear_sets/p1_combat.gear.json';
import P2CombatGear from './gear_sets/p2_combat.gear.json';
import P3CombatGear from './gear_sets/p3_combat.gear.json';
import P4CombatGear from './gear_sets/p4_combat.gear.json';
import P5CombatGear from './gear_sets/p5_combat.gear.json';
import PreraidCombatGear from './gear_sets/preraid_combat.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET_COMBAT = PresetUtils.makePresetGear('PreRaid Combat', PreraidCombatGear, { talentTree: 1 });
export const P1_PRESET_COMBAT = PresetUtils.makePresetGear('P1 Combat', P1CombatGear);
export const P2_PRESET_COMBAT = PresetUtils.makePresetGear('P2 Combat', P2CombatGear);
export const P3_PRESET_COMBAT = PresetUtils.makePresetGear('P3 Combat', P3CombatGear);
export const P4_PRESET_COMBAT = PresetUtils.makePresetGear('P4 Combat', P4CombatGear);
export const P5_PRESET_COMBAT = PresetUtils.makePresetGear('P5 Combat', P5CombatGear);

export const ROTATION_PRESET_COMBAT = PresetUtils.makePresetAPLRotation('Combat', CombatApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const CombatTalents = {
	name: 'Combat',
	data: SavedTalents.create({
		talentsString: '0331-2332030310230012321-003',
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
		mhImbue: RogueOptions_PoisonImbue.DeadlyPoison,
		ohImbue: RogueOptions_PoisonImbue.InstantPoison,
		thImbue: RogueOptions_PoisonImbue.WoundPoison,
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
