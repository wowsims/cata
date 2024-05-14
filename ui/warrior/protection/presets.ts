import * as PresetUtils from '../../core/preset_utils.js';
import { BattleElixir, Consumes, Explosive, Flask, Food, Glyphs, GuardianElixir, Potions, Profession, Spec } from '../../core/proto/common.js';
import { SavedTalents } from '../../core/proto/ui.js';
import {
	ProtectionWarrior_Options as ProtectionWarriorOptions,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	WarriorPrimeGlyph,
	WarriorShout,
} from '../../core/proto/warrior.js';
import DefaultApl from './apls/default.apl.json';
import P1BISGear from './gear_sets/p1_bis.gear.json';
import PreraidBISGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_BALANCED_PRESET = PresetUtils.makePresetGear('P1 PreRaid Preset', PreraidBISGear);
export const P1_BALANCED_PRESET = PresetUtils.makePresetGear('P1 BIS Preset', P1BISGear);


export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default APL', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '320003-002-33213201121210212031',
		glyphs: Glyphs.create({
			prime1: WarriorPrimeGlyph.GlyphOfRevenge,
			prime2: WarriorPrimeGlyph.GlyphOfShieldSlam,
			prime3: WarriorPrimeGlyph.GlyphOfDevastate,
			major1: WarriorMajorGlyph.GlyphOfShieldWall,
			major2: WarriorMajorGlyph.GlyphOfShockwave,
			major3: WarriorMajorGlyph.GlyphOfThunderClap,
			minor1: WarriorMinorGlyph.GlyphOfShatteringThrow,
			minor2: WarriorMinorGlyph.GlyphOfDemoralizingShout,
			minor3: WarriorMinorGlyph.GlyphOfCommand,
		}),
	}),
};

export const DefaultOptions = ProtectionWarriorOptions.create({
	classOptions: {
		shout: WarriorShout.WarriorShoutCommanding,
		useShatteringThrow: false,
		startingRage: 0,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfSteelskin,
	// battleElixir: BattleElixir.ElixirOfExpertise,
	// guardianElixir: GuardianElixir.ElixirOfProtection,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.EarthenPotion,
	prepopPotion: Potions.EarthenPotion,
});

export const OtherDefaults = {
	profession1: Profession.Leatherworking,
	profession2: Profession.Inscription,
};
