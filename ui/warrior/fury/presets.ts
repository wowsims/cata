import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Faction, Flask, Food, Glyphs, Potions } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { FuryWarrior_Options as WarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph, WarriorShout } from '../../core/proto/warrior';
import FuryApl from './apls/fury.apl.json';
import FurySunderApl from './apls/fury_sunder.apl.json';
import P1FuryGear from './gear_sets/p1_fury.gear.json';
import P2FuryGear from './gear_sets/p2_fury.gear.json';
import P3FuryAllianceGear from './gear_sets/p3_fury_alliance.gear.json';
import P3FuryHordeGear from './gear_sets/p3_fury_horde.gear.json';
import P4FuryAllianceGear from './gear_sets/p4_fury_alliance.gear.json';
import P4FuryHordeGear from './gear_sets/p4_fury_horde.gear.json';
import PreraidFuryGear from './gear_sets/preraid_fury.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_FURY_PRESET = PresetUtils.makePresetGear('Preraid Fury', PreraidFuryGear);
export const P1_FURY_PRESET = PresetUtils.makePresetGear('P1 Fury', P1FuryGear);
export const P2_FURY_PRESET = PresetUtils.makePresetGear('P2 Fury', P2FuryGear);
export const P3_FURY_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3 Fury [A]', P3FuryAllianceGear, { faction: Faction.Alliance });
export const P3_FURY_PRESET_HORDE = PresetUtils.makePresetGear('P3 Fury [H]', P3FuryHordeGear, { faction: Faction.Horde });
export const P4_FURY_PRESET_ALLIANCE = PresetUtils.makePresetGear('P4 Fury [A]', P4FuryAllianceGear, { faction: Faction.Alliance });
export const P4_FURY_PRESET_HORDE = PresetUtils.makePresetGear('P4 Fury [H]', P4FuryHordeGear, { faction: Faction.Horde });

export const ROTATION_FURY = PresetUtils.makePresetAPLRotation('Fury', FuryApl);
export const ROTATION_FURY_SUNDER = PresetUtils.makePresetAPLRotation('Fury + Sunder', FurySunderApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const FuryTalents = {
	name: 'Fury',
	data: SavedTalents.create({
		talentsString: '32002301233-305053000520310053120500351',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfWhirlwind,
			major2: WarriorMajorGlyph.GlyphOfHeroicStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfCommand,
			minor2: WarriorMinorGlyph.GlyphOfShatteringThrow,
			minor3: WarriorMinorGlyph.GlyphOfCharge,
		}),
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
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.PotionOfSpeed,
});
