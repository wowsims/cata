import * as PresetUtils from '../../core/preset_utils';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, UnholyDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import CataApl from '../../death_knight/unholy/apls/cata.apl.json';
import DefaultApl from '../../death_knight/unholy/apls/default.apl.json';
import P1Gear from '../../death_knight/unholy/gear_sets/p1.gear.json';
// import PreBISGear from '../../death_knight/unholy/gear_sets/prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
// export const PREBIS_GEAR_PRESET = PresetUtils.makePresetGear('Pre-bis', PreBISGear);
export const P1_BIS_GEAR_PRESET = PresetUtils.makePresetGear('P1', P1Gear);

export const DEFAULT_ROTATION_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const CATA_ROTATION_PRESET = PresetUtils.makePresetAPLRotation('!OLD! - Cata APL', CataApl);

// Preset options for EP weights
export const P1_UNHOLY_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
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

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '321111',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			major2: DeathKnightMajorGlyph.GlyphOfPestilence,
			major3: DeathKnightMajorGlyph.GlyphOfOutbreak,
			minor1: DeathKnightMinorGlyph.GlyphOfArmyOfTheDead,
			minor2: DeathKnightMinorGlyph.GlyphOfTranquilGrip,
			minor3: DeathKnightMinorGlyph.GlyphOfDeathsEmbrace,
		}),
	}),
};

// export const PREBIS_PRESET = PresetUtils.makePresetBuild('Pre-bis', {
// 	gear: PREBIS_GEAR_PRESET,
// 	epWeights: P1_UNHOLY_EP_PRESET,
// 	rotationType: APLRotationType.TypeAuto,
// });

export const P1_PRESET = PresetUtils.makePresetBuild('P1', {
	gear: P1_BIS_GEAR_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
	rotation: DEFAULT_ROTATION_PRESET,
	epWeights: P1_UNHOLY_EP_PRESET,
});

export const DefaultOptions = UnholyDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
	},
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
};

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76088, // Flask of Winter's Bite
	foodId: 74646, // Black Pepper Ribs and Shrimp
	potId: 76095, // Potion of Mogu Power
	prepotId: 76095, // Potion of Mogu Power
});
