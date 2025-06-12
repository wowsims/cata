import * as PresetUtils from '../../core/preset_utils';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat, UnitReference } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, UnholyDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from '../../death_knight/unholy/apls/default.apl.json';
import P2BISGear from '../../death_knight/unholy/gear_sets/p2.bis.gear.json';
import P3BISGear from '../../death_knight/unholy/gear_sets/p3.bis.gear.json';
import P4BISGear from '../../death_knight/unholy/gear_sets/p4.bis.gear.json';
import PreBISGear from '../../death_knight/unholy/gear_sets/prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const PREBIS_GEAR_PRESET = PresetUtils.makePresetGear('Pre-bis', PreBISGear);
export const P2_BIS_GEAR_PRESET = PresetUtils.makePresetGear('P2 - BIS', P2BISGear);
export const P3_BIS_GEAR_PRESET = PresetUtils.makePresetGear('P3 - BIS', P3BISGear);
export const P4_BIS_GEAR_PRESET = PresetUtils.makePresetGear('P4 - BIS', P4BISGear);

export const DEFAULT_ROTATION_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P2_UNHOLY_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P2',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 4.49,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 0.94,
			[Stat.StatHasteRating]: 2.4,
			[Stat.StatHitRating]: 2.6,
			[Stat.StatCritRating]: 1.43 + 0.69,
			[Stat.StatMasteryRating]: 1.65,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 6.13,
		},
	),
);

export const P3_UNHOLY_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P3',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 4.29,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 1.13,
			[Stat.StatHasteRating]: 2.4,
			[Stat.StatHitRating]: 2.61,
			[Stat.StatCritRating]: 2.33,
			[Stat.StatMasteryRating]: 1.87,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 6.39,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major3: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathsEmbrace,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const PREBIS_PRESET = PresetUtils.makePresetBuild('Pre-bis', {
	gear: PREBIS_GEAR_PRESET,
	epWeights: P2_UNHOLY_EP_PRESET,
	rotationType: APLRotationType.TypeAuto,
});

export const P2_PRESET = PresetUtils.makePresetBuild('P2', {
	gear: P2_BIS_GEAR_PRESET,
	epWeights: P2_UNHOLY_EP_PRESET,
	rotationType: APLRotationType.TypeAuto,
});

export const P3_PRESET = PresetUtils.makePresetBuild('P3', {
	gear: P3_BIS_GEAR_PRESET,
	epWeights: P3_UNHOLY_EP_PRESET,
	rotationType: APLRotationType.TypeAuto,
});

export const P4_PRESET = PresetUtils.makePresetBuild('P4', {
	gear: P4_BIS_GEAR_PRESET,
	epWeights: P3_UNHOLY_EP_PRESET,
	rotationType: APLRotationType.TypeAuto,
});

export const DefaultOptions = UnholyDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 55,
		petUptime: 1,
	},
	unholyFrenzyTarget: UnitReference.create(),
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
};

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58088, // Flask of Titanic Strength
	foodId: 62670, // Beer-Basted Crocolisk
	potId: 58146, // Golemblood Potion
	prepotId: 58146, // Golemblood Potion
});
