import * as PresetUtils from '../../core/preset_utils.js';
import { APLRotation_Type } from '../../core/proto/apl.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common.js';
import { BloodDeathKnight_Options, DeathKnightMajorGlyph, DeathKnightMinorGlyph } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/defensive.apl.json';
import P1BloodGear from './gear_sets/p1.gear.json';
// import PreRaidBloodGear from './gear_sets/preraid.gear.json';

// export const PRERAID_BLOOD_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreRaidBloodGear);
export const P1_BLOOD_PRESET = PresetUtils.makePresetGear('P1', P1BloodGear);

export const BLOOD_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Defensive', DefaultApl);

// Preset options for EP weights
export const P1_BLOOD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.45,
			[Stat.StatAgility]: 1.2,
			[Stat.StatStamina]: 3,
			[Stat.StatAttackPower]: 1,
			[Stat.StatHitRating]: 6,
			[Stat.StatCritRating]: 1.65,
			[Stat.StatHasteRating]: 1.58,
			[Stat.StatExpertiseRating]: 5,
			[Stat.StatArmor]: 1,
			[Stat.StatDodgeRating]: 2.5,
			[Stat.StatParryRating]: 2.44,
			[Stat.StatBonusArmor]: 1,
			[Stat.StatMasteryRating]: 7,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 12.29,
			[PseudoStat.PseudoStatOffHandDps]: 0.0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const BloodTalents = {
	name: 'Blood',
	data: SavedTalents.create({
		talentsString: '131131',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfFesteringBlood,
			major2: DeathKnightMajorGlyph.GlyphOfRegenerativeMagic,
			major3: DeathKnightMajorGlyph.GlyphOfOutbreak,
			minor1: DeathKnightMinorGlyph.GlyphOfTheLongWinter,
			minor2: DeathKnightMinorGlyph.GlyphOfArmyOfTheDead,
			minor3: DeathKnightMinorGlyph.GlyphOfResilientGrip,
		}),
	}),
};

export const P1_PRESET = PresetUtils.makePresetBuild('P1', {
	gear: P1_BLOOD_PRESET,
	talents: BloodTalents,
	rotationType: APLRotation_Type.TypeAuto,
	rotation: BLOOD_ROTATION_PRESET_DEFAULT,
	epWeights: P1_BLOOD_EP_PRESET,
});

export const DefaultOptions = BloodDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76087, // Flask of the Earth
	foodId: 74656, // Chun Tian Spring Rolls
	potId: 76095, // Potion of Mogu Power
	prepotId: 76095, // Potion of Mogu Power
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
};
