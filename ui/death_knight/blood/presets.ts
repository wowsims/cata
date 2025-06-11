import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common.js';
import { BloodDeathKnight_Options, DeathKnightMajorGlyph, DeathKnightMinorGlyph } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefensiveBloodApl from './apls/defensive.apl.json';
import SimpleBloodApl from './apls/simple.apl.json';
import P1BloodGear from './gear_sets/p1.gear.json';
import P3BloodBalancedGear from './gear_sets/p3-balanced.gear.json';
import P3BloodDefensiveGear from './gear_sets/p3-defensive.gear.json';
import P3BloodOffensiveGear from './gear_sets/p3-offensive.gear.json';
import PreRaidBloodGear from './gear_sets/preraid.gear.json';

export const PRERAID_BLOOD_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreRaidBloodGear);
export const P1_BLOOD_PRESET = PresetUtils.makePresetGear('P1', P1BloodGear);
export const P3_BLOOD_BALANCED_PRESET = PresetUtils.makePresetGear('P3-Balanced', P3BloodBalancedGear);
export const P3_BLOOD_DEFENSIVE_PRESET = PresetUtils.makePresetGear('P3-Defensive', P3BloodDefensiveGear);
export const P3_BLOOD_OFFENSIVE_PRESET = PresetUtils.makePresetGear('P3-Offensive', P3BloodOffensiveGear);

export const BLOOD_SIMPLE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Simple', SimpleBloodApl);
export const BLOOD_DEFENSIVE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Defensive', DefensiveBloodApl);

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

// Preset options for EP weights
export const P3_BLOOD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P3',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.47,
			[Stat.StatAgility]: 1.46,
			[Stat.StatStamina]: 7,
			[Stat.StatAttackPower]: 1,
			[Stat.StatHitRating]: 6,
			[Stat.StatCritRating]: 1,
			[Stat.StatHasteRating]: 4,
			[Stat.StatExpertiseRating]: 5,
			[Stat.StatArmor]: 1,
			[Stat.StatDodgeRating]: 0.5,
			[Stat.StatParryRating]: 0.5,
			[Stat.StatBonusArmor]: 1,
			[Stat.StatMasteryRating]: 3,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 12.34,
			[PseudoStat.PseudoStatOffHandDps]: 0.0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const BloodTalents = {
	name: 'Blood',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			major2: DeathKnightMajorGlyph.GlyphOfDancingRuneWeapon,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const P1_PRESET = PresetUtils.makePresetBuild('P1', {
	gear: P1_BLOOD_PRESET,
	epWeights: P1_BLOOD_EP_PRESET,
});

export const P3_PRESET = PresetUtils.makePresetBuild('P3', {
	gear: P3_BLOOD_BALANCED_PRESET,
	epWeights: P3_BLOOD_EP_PRESET,
});

export const DefaultOptions = BloodDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58085, // Flask of Steelskin
	foodId: 62670, // Beer-Basted Crocolisk
	potId: 58146, // Golemblood Potion
	prepotId: 58146, // Golemblood Potion
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
	distanceFromTarget: 5,
};
