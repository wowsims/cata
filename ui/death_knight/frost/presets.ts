import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Explosive, Flask, Food, Glyphs, PetFood, Potions } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import DefaultAPL from '../../death_knight/frost/apls/default.apl.json'
import DefaultGear from '../../death_knight/frost/gear_sets/default.gear.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const DEFAULT_GEAR_PRESET = PresetUtils.makePresetGear('Gear', DefaultGear);

export const SINGLE_TARGET_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Single Target', DefaultAPL);
export const AOE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('AoE', DefaultAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const SingleTargetTalents = {
	name: 'Single Target',
	data: SavedTalents.create({
		// talentsString: '-320050500002-2302003350032052000150013133151',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfTheGhoul,
		// 	major2: DeathKnightMajorGlyph.GlyphOfIcyTouch,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDarkDeath,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const AoeTalents = {
	name: 'AOE',
	data: SavedTalents.create({
		// talentsString: '-320050500002-2302303050032052000150013133151',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfTheGhoul,
		// 	major2: DeathKnightMajorGlyph.GlyphOfIcyTouch,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDeathAndDecay,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const DefaultOptions = FrostDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
		petUptime: 1,
	},
});

export const OtherDefaults = {};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.GolembloodPotion,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	prepopPotion: Potions.GolembloodPotion,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});
