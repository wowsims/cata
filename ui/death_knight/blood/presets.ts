import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common.js';
import { BloodDeathKnight_Options, DeathKnightMajorGlyph, DeathKnightMinorGlyph } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui.js';
import BloodAggroApl from './apls/blood_aggro.apl.json';
import BloodIcyTouchApl from './apls/blood_icy_touch.apl.json';
import P1BloodGear from './gear_sets/p1_blood.gear.json';
import P2BloodGear from './gear_sets/p2_blood.gear.json';
import P3BloodGear from './gear_sets/p3_blood.gear.json';
import P4BloodGear from './gear_sets/p4_blood.gear.json';

export const P1_BLOOD_PRESET = PresetUtils.makePresetGear('P1 Blood', P1BloodGear);
export const P2_BLOOD_PRESET = PresetUtils.makePresetGear('P2 Blood', P2BloodGear);
export const P3_BLOOD_PRESET = PresetUtils.makePresetGear('P3 Blood', P3BloodGear);
export const P4_BLOOD_PRESET = PresetUtils.makePresetGear('P4 Blood', P4BloodGear);

export const BLOOD_IT_SPAM_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Blood Icy Touch', BloodIcyTouchApl);
export const BLOOD_AGGRO_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Blood Aggro', BloodAggroApl);

export const BloodTalents = {
	name: 'Blood',
	data: SavedTalents.create({
		// talentsString: '005512153330030320102013-3050505000023-005',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfDisease,
		// 	major2: DeathKnightMajorGlyph.GlyphOfRuneStrike,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDarkCommand,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfBloodTap,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const BloodAggroTalents = {
	name: 'Blood Aggro',
	data: SavedTalents.create({
		// talentsString: '0355220530303303201020131301--0052003050032',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfDancingRuneWeapon,
		// 	major2: DeathKnightMajorGlyph.GlyphOfRuneStrike,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDarkCommand,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfBloodTap,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const DoubleBuffBloodTalents = {
	name: '2B Blood',
	data: SavedTalents.create({
		// talentsString: '005512153330030320102013-3050505000023201-002',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfDisease,
		// 	major2: DeathKnightMajorGlyph.GlyphOfRuneStrike,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDarkCommand,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfBloodTap,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const DefaultOptions = BloodDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.IndestructiblePotion,
});
