import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Faction, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { ArmsWarrior_Options as WarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph, WarriorShout } from '../../core/proto/warrior';
import ArmsApl from './apls/arms.apl.json';
import ArmsSunderApl from './apls/arms_sunder.apl.json';
import P1ArmsGear from './gear_sets/p1_arms.gear.json';
import P2ArmsGear from './gear_sets/p2_arms.gear.json';
import P3Arms2pAllianceGear from './gear_sets/p3_arms_2p_alliance.gear.json';
import P3Arms2pHordeGear from './gear_sets/p3_arms_2p_horde.gear.json';
import P3Arms4pAllianceGear from './gear_sets/p3_arms_4p_alliance.gear.json';
import P3Arms4pHordeGear from './gear_sets/p3_arms_4p_horde.gear.json';
import P4ArmsAllianceGear from './gear_sets/p4_arms_alliance.gear.json';
import P4ArmsHordeGear from './gear_sets/p4_arms_horde.gear.json';
import PreraidArmsGear from './gear_sets/preraid_arms.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_ARMS_PRESET = PresetUtils.makePresetGear('Preraid Arms', PreraidArmsGear);
export const P1_ARMS_PRESET = PresetUtils.makePresetGear('P1 Arms', P1ArmsGear);
export const P2_ARMS_PRESET = PresetUtils.makePresetGear('P2 Arms', P2ArmsGear);
export const P3_ARMS_2P_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3 Arms 2p [A]', P3Arms2pAllianceGear, { faction: Faction.Alliance });
export const P3_ARMS_4P_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3 Arms 4p [A]', P3Arms4pAllianceGear, { faction: Faction.Alliance });
export const P3_ARMS_2P_PRESET_HORDE = PresetUtils.makePresetGear('P3 Arms 2p [H]', P3Arms2pHordeGear, { faction: Faction.Horde });
export const P3_ARMS_4P_PRESET_HORDE = PresetUtils.makePresetGear('P3 Arms 4p [H]', P3Arms4pHordeGear, { faction: Faction.Horde });
export const P4_ARMS_PRESET_ALLIANCE = PresetUtils.makePresetGear('P4 Arms [A]', P4ArmsAllianceGear, { faction: Faction.Alliance });
export const P4_ARMS_PRESET_HORDE = PresetUtils.makePresetGear('P4 Arms [H]', P4ArmsHordeGear, { faction: Faction.Horde });

export const ROTATION_ARMS = PresetUtils.makePresetAPLRotation('Arms', ArmsApl);
export const ROTATION_ARMS_SUNDER = PresetUtils.makePresetAPLRotation('Arms + Sunder', ArmsSunderApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const ArmsTalents = {
	name: 'Arms',
	data: SavedTalents.create({
		// talentsString: '3022032023335100102012213231251-305-2033',
		// glyphs: Glyphs.create({
		// 	major1: WarriorMajorGlyph.GlyphOfRending,
		// 	major2: WarriorMajorGlyph.GlyphOfMortalStrike,
		// 	major3: WarriorMajorGlyph.GlyphOfExecution,
		// 	minor1: WarriorMinorGlyph.GlyphOfThunderClap,
		// 	minor2: WarriorMinorGlyph.GlyphOfCommand,
		// 	minor3: WarriorMinorGlyph.GlyphOfShatteringThrow,
		// }),
	}),
};

export const DefaultOptions = WarriorOptions.create({
	classOptions: {
		startingRage: 0,
		useShatteringThrow: true,
		shout: WarriorShout.WarriorShoutCommanding,
	},
	useRecklessness: true,
	disableExpertiseGemming: false,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodSpicedWormBurger,
	defaultPotion: Potions.EarthenPotion,
	prepopPotion: Potions.GolembloodPotion,
});
