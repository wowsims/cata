import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Profession, Stat } from '../../core/proto/common';
import { FrostMage_Options as MageOptions } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import FrostApl from './apls/frost.apl.json';
import FrostAoeApl from './apls/frost_aoe.apl.json';
import P1FrostGear from './gear_sets/p1_frost.gear.json';
import P2FrostGear from './gear_sets/p2_frost.gear.json';
import P3FrostAllianceGear from './gear_sets/p3_frost_alliance.gear.json';
import P3FrostHordeGear from './gear_sets/p3_frost_horde.gear.json';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const FROST_P1_PRESET = PresetUtils.makePresetGear('Frost P1 Preset', P1FrostGear);
export const FROST_P2_PRESET = PresetUtils.makePresetGear('Frost P2 Preset', P2FrostGear);
export const FROST_P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('Frost P3 Preset [A]', P3FrostAllianceGear);
export const FROST_P3_PRESET_HORDE = PresetUtils.makePresetGear('Frost P3 Preset [H]', P3FrostHordeGear);

export const FROST_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Frost', FrostApl);
export const FROST_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Frost AOE', FrostAoeApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Frost P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 0.48,
		[Stat.StatSpirit]: 0.42,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 0.38,
		[Stat.StatCritRating]: 0.58,
		[Stat.StatHasteRating]: 0.94,
		[Stat.StatMP5]: 0.09,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const FrostTalents = {
	name: 'Frost',
	data: SavedTalents.create({
		// talentsString: '23000503110003--0533030310233100030152231351',
		// glyphs: Glyphs.create({
		// 	major1: MageMajorGlyph.GlyphOfFrostbolt,
		// 	major2: MageMajorGlyph.GlyphOfEternalWater,
		// 	major3: MageMajorGlyph.GlyphOfMoltenArmor,
		// 	minor1: MageMinorGlyph.GlyphOfSlowFall,
		// 	minor2: MageMinorGlyph.GlyphOfFrostWard,
		// 	minor3: MageMinorGlyph.GlyphOfBlastWave,
		// }),
	}),
};

export const DefaultFrostOptions = MageOptions.create({
	classOptions: {},
	waterElementalDisobeyChance: 0.1,
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58086, // Flask of the Draconic Mind
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
	prepotId: 58091, // Volcanic Potion
});
export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
