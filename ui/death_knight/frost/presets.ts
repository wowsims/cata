import * as PresetUtils from '../../core/preset_utils';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import MasterFrostAPL from '../../death_knight/frost/apls/masterfrost.apl.json';
import ObliterateAPL from '../../death_knight/frost/apls/obliterate.apl.json';
import P12HObliterateGear from '../../death_knight/frost/gear_sets/p1.2h-obliterate.gear.json';
// import P1DWObliterateGear from '../../death_knight/frost/gear_sets/p1.dw-obliterate.gear.json';
import P1MasterfrostGear from '../../death_knight/frost/gear_sets/p1.masterfrost.gear.json';
// import PreBISGear from '../../death_knight/frost/gear_sets/prebis.gear.json';

// export const P1_DW_OBLITERATE_GEAR_PRESET = PresetUtils.makePresetGear('P1 DW Obliterate', P1DWObliterateGear);
export const P1_2H_OBLITERATE_GEAR_PRESET = PresetUtils.makePresetGear('P1 2h Obliterate', P12HObliterateGear);
export const P1_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('P1 Masterfrost', P1MasterfrostGear);
// export const PREBIS_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('Pre-bis Masterfrost', PreBISGear);

export const OBLITERATE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Obliterate', ObliterateAPL);
export const MASTERFROST_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Masterfrost', MasterFrostAPL);

// // Preset options for EP weights
// export const P1_DUAL_WIELD_OBLITERATE_EP_PRESET = PresetUtils.makePresetEpWeights(
// 	'P1 DW Obliterate',
// 	Stats.fromMap(
// 		{
// 			[Stat.StatStrength]: 2.92,
// 			[Stat.StatArmor]: 0.03,
// 			[Stat.StatAttackPower]: 1,
// 			[Stat.StatExpertiseRating]: 0.56,
// 			[Stat.StatHasteRating]: 1.3,
// 			[Stat.StatHitRating]: 1.22,
// 			[Stat.StatCritRating]: 1.06,
// 			[Stat.StatMasteryRating]: 1.11,
// 		},
// 		{
// 			[PseudoStat.PseudoStatMainHandDps]: 6.05,
// 			[PseudoStat.PseudoStatOffHandDps]: 3.85,
// 			[PseudoStat.PseudoStatPhysicalHitPercent]: 146.53,
// 			[PseudoStat.PseudoStatSpellHitPercent]: 41.91,
// 		},
// 	),
// );

export const P1_2H_OBLITERATE_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 2h Obliterate',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.0,
			[Stat.StatHitRating]: 0.87,
			[Stat.StatExpertiseRating]: 0.87,
			[Stat.StatMasteryRating]: 0.35,
			[Stat.StatCritRating]: 0.44,
			[Stat.StatHasteRating]: 0.39,
			[Stat.StatAttackPower]: 0.37,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 2.95,
		},
	),
);

export const P1_MASTERFROST_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 Masterfrost',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.0,
			[Stat.StatHitRating]: 0.8,
			[Stat.StatExpertiseRating]: 0.8,
			[Stat.StatMasteryRating]: 0.48,
			[Stat.StatHasteRating]: 0.38,
			[Stat.StatAttackPower]: 0.37,
			[Stat.StatCritRating]: 0.35,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.47,
			[PseudoStat.PseudoStatOffHandDps]: 0.7,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '221111',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			major2: DeathKnightMajorGlyph.GlyphOfPestilence,
			major3: DeathKnightMajorGlyph.GlyphOfLoudHorn,
			minor1: DeathKnightMinorGlyph.GlyphOfArmyOfTheDead,
			minor2: DeathKnightMinorGlyph.GlyphOfTranquilGrip,
			minor3: DeathKnightMinorGlyph.GlyphOfDeathGate,
		}),
	}),
};

export const DefaultOptions = FrostDeathKnight_Options.create({
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

// export const PRESET_BUILD_DW_OBLITERATE = PresetUtils.makePresetBuild('P1 DW Obliterate', {
// 	gear: P1_DW_OBLITERATE_GEAR_PRESET,
// 	talents: DefaultTalents,
// 	rotationType: APLRotationType.TypeAPL,
// 	rotation: OBLITERATE_ROTATION_PRESET_DEFAULT,
// 	epWeights: P1_DW_OBLITERATE_EP_PRESET,
// });

export const PRESET_BUILD_2H_OBLITERATE = PresetUtils.makePresetBuild('P1 2h Obliterate', {
	gear: P1_2H_OBLITERATE_GEAR_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
	rotation: OBLITERATE_ROTATION_PRESET_DEFAULT,
	epWeights: P1_2H_OBLITERATE_EP_PRESET,
});

export const PRESET_BUILD_MASTERFROST = PresetUtils.makePresetBuild('P1 Masterfrost', {
	gear: P1_MASTERFROST_GEAR_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
	rotation: MASTERFROST_ROTATION_PRESET_DEFAULT,
	epWeights: P1_MASTERFROST_EP_PRESET,
});

// export const PRESET_BUILD_PREBIS = PresetUtils.makePresetBuild('P1 - Pre-bis', {
// 	gear: PREBIS_MASTERFROST_GEAR_PRESET,
// 	talents: DefaultTalents,
// 	rotationType: APLRotationType.TypeAPL,
// 	rotation: MASTERFROST_ROTATION_PRESET_DEFAULT,
// 	epWeights: P1_MASTERFROST_EP_PRESET,
// });
