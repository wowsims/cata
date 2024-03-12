import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Explosive, Flask, Food, Glyphs, PetFood, Potions } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import FrostBlPestiApl from './apls/frost_bl_pesti.apl.json';
import FrostUhPestiApl from './apls/frost_uh_pesti.apl.json';
import P1FrostGear from './gear_sets/p1_frost.gear.json';
import P1FrostSubUhGear from './gear_sets/p1_frost_subUh.gear.json';
import P2FrostGear from './gear_sets/p2_frost.gear.json';
import P3FrostGear from './gear_sets/p3_frost.gear.json';
import P4FrostGear from './gear_sets/p4_frost.gear.json';
import PreraidFrostGear from './gear_sets/preraid_frost.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_FROST_PRESET = PresetUtils.makePresetGear('Pre-Raid Frost', PreraidFrostGear);
export const P1_FROST_PRESET = PresetUtils.makePresetGear('P1 Frost', P1FrostGear);
export const P2_FROST_PRESET = PresetUtils.makePresetGear('P2 Frost', P2FrostGear);
export const P3_FROST_PRESET = PresetUtils.makePresetGear('P3 Frost', P3FrostGear);
export const P4_FROST_PRESET = PresetUtils.makePresetGear('P4 Frost', P4FrostGear);
export const P1_FROSTSUBUNH_PRESET = PresetUtils.makePresetGear('P1 Frost Sub Unh', P1FrostSubUhGear);

export const FROST_BL_PESTI_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Frost BL Pesti', FrostBlPestiApl);
export const FROST_UH_PESTI_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Frost UH Pesti', FrostUhPestiApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const FrostTalents = {
	name: 'Frost BL',
	data: SavedTalents.create({
		// talentsString: '23050005-32005350352203012300033101351',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfObliterate,
		// 	major2: DeathKnightMajorGlyph.GlyphOfFrostStrike,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDisease,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const FrostUnholyTalents = {
	name: 'Frost UH',
	data: SavedTalents.create({
		// talentsString: '01-32002350342203012300033101351-230200305003',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfObliterate,
		// 	major2: DeathKnightMajorGlyph.GlyphOfFrostStrike,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDisease,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const DefaultFrostOptions = FrostDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
		petUptime: 1,
	},
});

export const OtherDefaults = {};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	prepopPotion: Potions.PotionOfSpeed,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});
