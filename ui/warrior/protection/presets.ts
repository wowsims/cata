import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Race, Stat } from '../../core/proto/common.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { ProtectionWarrior_Options as ProtectionWarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph } from '../../core/proto/warrior.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1BISGear from './gear_sets/p1_bis.gear.json';
import P3BISGear from './gear_sets/p3_bis.gear.json';
import P4BISGear from './gear_sets/p4_bis.gear.json';
import P4NelfBISGear from './gear_sets/p4_nelf_bis.gear.json';
import PreraidBISGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_BALANCED_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidBISGear);
export const P1_BALANCED_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1BISGear);
export const P3_BALANCED_PRESET = PresetUtils.makePresetGear('P3 - BIS', P3BISGear);
export const P4_BALANCED_PRESET = PresetUtils.makePresetGear('P4 - BIS', P4BISGear);
export const P4_NELF_BALANCED_PRESET = PresetUtils.makePresetGear('P4 - BIS (Nelf)', P4NelfBISGear);

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default APL', DefaultApl);

// Preset options for EP weights
export const P3_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatArmor]: 2.155,
			[Stat.StatBonusArmor]: 2.155,
			[Stat.StatStamina]: 12.442,
			[Stat.StatStrength]: 1.4,
			[Stat.StatAgility]: 0.26,
			[Stat.StatAttackPower]: 0.196,
			[Stat.StatExpertiseRating]: 0.863,
			[Stat.StatHitRating]: 0.736,
			[Stat.StatCritRating]: 0.336,
			[Stat.StatHasteRating]: 0.048,
			[Stat.StatDodgeRating]: 4.801,
			[Stat.StatParryRating]: 4.801,
			[Stat.StatMasteryRating]: 7.415,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 6.081,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfShieldWall,
		}),
	}),
};

export const DefaultOptions = ProtectionWarriorOptions.create({
	classOptions: {
		startingRage: 0,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	battleElixirId: 58148, // Elixir of the Master (not found in list)
	guardianElixirId: 58093, // Elixir of Deep Earth (not found in list)
	foodId: 62670, // Beer-Basted Crocolisk
	potId: 58090, // Earthen Potion
	prepotId: 58090, // Earthen Potion
	explosiveId: 89637, // Big Daddy Explosive
});

export const OtherDefaults = {
	profession1: Profession.Leatherworking,
	profession2: Profession.Inscription,
};

export const P4_PRESET_BUILD = PresetUtils.makePresetBuild('P4 - Default', {
	race: Race.RaceGnome,
	gear: P4_BALANCED_PRESET,
});

export const P4_NELF_PRESET_BUILD = PresetUtils.makePresetBuild('P4 - Default (Nelf)', {
	race: Race.RaceNightElf,
	gear: P4_NELF_BALANCED_PRESET,
});
