import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Explosive, Flask, Food, Glyphs, PetFood, Potions, UnitReference } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, UnholyDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import Uh2hSsApl from './apls/uh_2h_ss.apl.json';
import UhDndAoeApl from './apls/uh_dnd_aoe.apl.json';
import UhDwSsApl from './apls/unholy_dw_ss.apl.json';
import P1Uh2hGear from './gear_sets/p1_uh_2h.gear.json';
import P1UhDwGear from './gear_sets/p1_uh_dw.gear.json';
import P2UhDwGear from './gear_sets/p2_uh_dw.gear.json';
import P3UhDwGear from './gear_sets/p3_uh_dw.gear.json';
import P4Uh2hGear from './gear_sets/p4_uh_2h.gear.json';
import P4UhDwGear from './gear_sets/p4_uh_dw.gear.json';
import PreraidUh2hGear from './gear_sets/preraid_uh_2h.gear.json';
import PreraidUhDwGear from './gear_sets/preraid_uh_dw.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_UNHOLY_2H_PRESET = PresetUtils.makePresetGear('Pre-Raid 2H Unholy', PreraidUh2hGear);
export const P1_UNHOLY_2H_PRESET = PresetUtils.makePresetGear('P1 2H Unholy', P1Uh2hGear);
export const P4_UNHOLY_2H_PRESET = PresetUtils.makePresetGear('P4 2H Unholy', P4Uh2hGear);
export const PRERAID_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('Pre-Raid DW Unholy', PreraidUhDwGear);
export const P1_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P1 DW Unholy', P1UhDwGear);
export const P2_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P2 DW Unholy', P2UhDwGear);
export const P3_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P3 DW Unholy', P3UhDwGear);
export const P4_UNHOLY_DW_PRESET = PresetUtils.makePresetGear('P4 DW Unholy', P4UhDwGear);

export const UNHOLY_DW_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Unholy DW SS', UhDwSsApl);
export const UNHOLY_2H_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Unholy 2H SS', Uh2hSsApl);
export const UNHOLY_DND_AOE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Unholy DND AOE', UhDndAoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const UnholyDualWieldTalents = {
	name: 'Unholy DW',
	data: SavedTalents.create({
		talentsString: '-320043500002-2300303050032152000150013133051',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathKnightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathKnightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
			minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldSSTalents = {
	name: 'Unholy DW SS',
	data: SavedTalents.create({
		talentsString: '-320033500002-2301303050032151000150013133151',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathKnightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathKnightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
			minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const Unholy2HTalents = {
	name: 'Unholy 2H',
	data: SavedTalents.create({
		talentsString: '-320050500002-2302003350032052000150013133151',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathKnightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathKnightMajorGlyph.GlyphOfDarkDeath,
			minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
			minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyAoeTalents = {
	name: 'Unholy AOE',
	data: SavedTalents.create({
		talentsString: '-320050500002-2302303050032052000150013133151',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathKnightMajorGlyph.GlyphOfIcyTouch,
			major3: DeathKnightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
			minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DefaultUnholyOptions = UnholyDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
		petUptime: 1,
	},
	unholyFrenzyTarget: UnitReference.create(),
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
