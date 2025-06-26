import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, Profession, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { FireMage_Options as MageOptions, FireMage_Rotation, MageMajorGlyph as MajorGlyph, MageMinorGlyph as MinorGlyph } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats, UnitStat, UnitStatPresets } from '../../core/proto_utils/stats';
import FireApl from './apls/fire.apl.json';
//import FireAoeApl from './apls/fire_aoe.apl.json';
import P1FireBisGear from './gear_sets/p1_bis.gear.json';
import P1FirePrebisGear from './gear_sets/p1_prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_BIS_PRESET = PresetUtils.makePresetGear('P1 Preset', P1FireBisGear);
export const PREBIS_PRESET = PresetUtils.makePresetGear('P3 Pre-raid', P1FirePrebisGear);

export const P1TrollDefaultSimpleRotation = FireMage_Rotation.create({
	combustThreshold: 515000,
	combustLastMomentLustPercentage: 140000,
	combustNoLustPercentage: 260000,
});

export const P1NoTrollDefaultSimpleRotation = FireMage_Rotation.create({
	combustThreshold: 470000,
	combustLastMomentLustPercentage: 115000,
	combustNoLustPercentage: 225000,
});

export const P1_SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('P1 - Default', Spec.SpecFireMage, P1TrollDefaultSimpleRotation);
export const P1_SIMPLE_ROTATION_NO_TROLL = PresetUtils.makePresetSimpleRotation('P1 - Not Troll', Spec.SpecFireMage, P1NoTrollDefaultSimpleRotation);

//export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFireMage, DefaultSimpleRotation);
export const FIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('APL', FireApl);

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.33,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatHitRating]: 1.09,
		[Stat.StatCritRating]: 0.62,
		[Stat.StatHasteRating]: 0.82,
		[Stat.StatMasteryRating]: 0.46,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '111122',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfCombustion,
			major2: MajorGlyph.GlyphOfInfernoBlast,
			major3: MajorGlyph.GlyphOfManaGem,
			minor1: MinorGlyph.GlyphOfMomentum,
			minor2: MinorGlyph.GlyphOfMirrorImage,
			minor3: MinorGlyph.GlyphOfTheUnboundElemental
		}),
	}),
};

export const DefaultFireOptions = MageOptions.create({
	classOptions: {},
});

export const DefaultFireConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const COMBUSTION_BREAKPOINT: UnitStatPresets = {
	unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	presets: new Map([
		['11-tick - Combust', 4.986888],
		['12-tick - Combust', 15.008639],
		['13-tick - Combust', 25.07819],
		['14-tick - Combust', 35.043908],
		['15-tick - Combust', 45.032653],
		['16-tick - Combust', 54.918692],
		['17-tick - Combust', 64.880489],
		['18-tick - Combust', 74.978158],
		['19-tick - Combust', 85.01391],
		['20-tick - Combust', 95.121989],
		['21-tick - Combust', 105.128247],
		['22-tick - Combust', 114.822817],
		['23-tick - Combust', 124.971929],
		['24-tick - Combust', 135.017682],
		['25-tick - Combust', 144.798102],
		['26-tick - Combust', 154.777135],
		['27-tick - Combust', 164.900732],
		['28-tick - Combust', 175.103239],
		['29-tick - Combust', 185.306786],
	]),
};

export const GLYPHED_COMBUSTION_BREAKPOINT: UnitStatPresets = {
	unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
	presets: new Map([
		['21-tick - Combust (Glyph)', 2.511543],
		['22-tick - Combust (Glyph)', 7.469114],
		['23-tick - Combust (Glyph)', 12.549253],
		['24-tick - Combust (Glyph)', 17.439826],
		['25-tick - Combust (Glyph)', 22.473989],
		['26-tick - Combust (Glyph)', 27.469742],
		['27-tick - Combust (Glyph)', 32.538122],
		['28-tick - Combust (Glyph)', 37.457064],
		['29-tick - Combust (Glyph)', 42.551695],
		['30-tick - Combust (Glyph)', 47.601498],
		['31-tick - Combust (Glyph)', 52.555325],
		['32-tick - Combust (Glyph)', 57.604438],
		['33-tick - Combust (Glyph)', 62.469563],
		['34-tick - Combust (Glyph)', 67.364045],
		['35-tick - Combust (Glyph)', 72.562584],
		['36-tick - Combust (Glyph)', 77.462321],
		['37-tick - Combust (Glyph)', 82.648435],
		['38-tick - Combust (Glyph)', 87.44146],
		['39-tick - Combust (Glyph)', 92.492819],
	]),
};

export const P1_PREBIS_PRESET_BUILD = PresetUtils.makePresetBuild('P1 - Pre-BIS (Troll)', {
	race: Race.RaceTroll,
	gear: PREBIS_PRESET,
	rotation: P1_SIMPLE_ROTATION_DEFAULT,
	epWeights: DEFAULT_EP_PRESET,
});

export const P1_PREBIS_PRESET_BUILD_NO_TROLL = PresetUtils.makePresetBuild('P1 - Pre-BIS (Worgen)', {
	race: Race.RaceWorgen,
	gear: PREBIS_PRESET,
	rotation: P1_SIMPLE_ROTATION_NO_TROLL,
	epWeights: DEFAULT_EP_PRESET,
});

export const P1_PRESET_BUILD = PresetUtils.makePresetBuild('P1 - BIS (Troll)', {
	race: Race.RaceTroll,
	gear: P1_BIS_PRESET,
	rotation: P1_SIMPLE_ROTATION_DEFAULT,
	epWeights: DEFAULT_EP_PRESET,
});

export const P1_PRESET_BUILD_NO_TROLL = PresetUtils.makePresetBuild('P1 - BIS (Worgen)', {
	race: Race.RaceWorgen,
	gear: P1_BIS_PRESET,
	rotation: P1_SIMPLE_ROTATION_NO_TROLL,
	epWeights: DEFAULT_EP_PRESET,
});
