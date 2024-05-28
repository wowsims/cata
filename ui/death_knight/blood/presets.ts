import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, TinkerHands } from '../../core/proto/common.js';
import { BloodDeathKnight_Options, DeathKnightMajorGlyph, DeathKnightMinorGlyph,DeathKnightPrimeGlyph } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui.js';
import DefensiveBloodApl from './apls/defensive.apl.json';
import SimpleBloodApl from './apls/simple.apl.json';
import P1BloodGear from './gear_sets/p1.gear.json';
import PreRaidBloodGear from './gear_sets/preraid.gear.json';

export const PRERAID_BLOOD_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreRaidBloodGear);
export const P1_BLOOD_PRESET = PresetUtils.makePresetGear('P1', P1BloodGear);

export const BLOOD_SIMPLE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Simple', SimpleBloodApl);
export const BLOOD_DEFENSIVE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Defensive', DefensiveBloodApl);

export const BloodTalents = {
	name: 'Blood',
	data: SavedTalents.create({
		talentsString: '03323200132222311321-2-003',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfDeathStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfHeartStrike,
			prime3: DeathKnightPrimeGlyph.GlyphOfRuneStrike,
			major1: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			major2: DeathKnightMajorGlyph.GlyphOfDancingRuneWeapon,
			major3: DeathKnightMajorGlyph.GlyphOfBoneShield,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const DefaultOptions = BloodDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfSteelskin,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
	distanceFromTarget: 5,
};
