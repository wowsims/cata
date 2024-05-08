import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Explosive, Flask, Food, Glyphs, Potions, Profession, TinkerHands } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightPrimeGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import SingleTargetAPL from '../../death_knight/frost/apls/st.apl.json'
import P1Gear from '../../death_knight/frost/gear_sets/p1.gear.json'

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_GEAR_PRESET = PresetUtils.makePresetGear('P1', P1Gear);

export const SINGLE_TARGET_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Single Target', SingleTargetAPL);
export const AOE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('AoE', SingleTargetAPL);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const SingleTargetTalents = {
	name: 'Single Target',
	data: SavedTalents.create({
		talentsString: '2032-20330022233112012301-003',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfFrostStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfObliterate,
			prime3: DeathKnightPrimeGlyph.GlyphOfHowlingBlast,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfDarkSuccor,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
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

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Jewelcrafting,
	distanceFromTarget: 5,
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTitanicStrength,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});
