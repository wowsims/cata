import * as PresetUtils from '../../core/preset_utils.js';
import { BattleElixir, Consumes, Explosive, Food, Glyphs, GuardianElixir, Potions, Spec } from '../../core/proto/common.js';
import { SavedTalents } from '../../core/proto/ui.js';
import {
	ProtectionWarrior_Options as ProtectionWarriorOptions,
	ProtectionWarrior_Rotation as ProtectionWarriorRotation,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	WarriorShout,
} from '../../core/proto/warrior.js';
import DefaultApl from './apls/default.apl.json';
import P1BalancedGear from './gear_sets/p1_balanced.gear.json';
import P2SurvivalGear from './gear_sets/p2_survival.gear.json';
import P3Gear from './gear_sets/p3.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import PreraidP4Gear from './gear_sets/p4_preraid.gear.json';
import PreraidBalancedGear from './gear_sets/preraid_balanced.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_BALANCED_PRESET = PresetUtils.makePresetGear('P1 PreRaid Preset', PreraidBalancedGear);
export const P4_PRERAID_PRESET = PresetUtils.makePresetGear('P4 PreRaid Preset', PreraidP4Gear);
export const P1_BALANCED_PRESET = PresetUtils.makePresetGear('P1 Preset', P1BalancedGear);
export const P2_SURVIVAL_PRESET = PresetUtils.makePresetGear('P2 Preset', P2SurvivalGear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default APL', DefaultApl);
export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Cooldowns', Spec.SpecProtectionWarrior, ProtectionWarriorRotation.create());

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		// talentsString: '2500030023-302-053351225000012521030113321',
		// glyphs: Glyphs.create({
		// 	major1: WarriorMajorGlyph.GlyphOfBlocking,
		// 	major2: WarriorMajorGlyph.GlyphOfVigilance,
		// 	major3: WarriorMajorGlyph.GlyphOfDevastate,
		// 	minor1: WarriorMinorGlyph.GlyphOfCharge,
		// 	minor2: WarriorMinorGlyph.GlyphOfThunderClap,
		// 	minor3: WarriorMinorGlyph.GlyphOfCommand,
		// }),
	}),
};

export const UATalents = {
	name: 'UA',
	data: SavedTalents.create({
		// talentsString: '35023301230051002020120002-2-05035122500000252',
		// glyphs: Glyphs.create({
		// 	major1: WarriorMajorGlyph.GlyphOfRevenge,
		// 	major2: WarriorMajorGlyph.GlyphOfHeroicStrike,
		// 	major3: WarriorMajorGlyph.GlyphOfSweepingStrikes,
		// 	minor1: WarriorMinorGlyph.GlyphOfCharge,
		// 	minor2: WarriorMinorGlyph.GlyphOfThunderClap,
		// 	minor3: WarriorMinorGlyph.GlyphOfCommand,
		// }),
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
	battleElixir: BattleElixir.ElixirOfExpertise,
	guardianElixir: GuardianElixir.ElixirOfProtection,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.EarthenPotion,
	prepopPotion: Potions.EarthenPotion,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});
