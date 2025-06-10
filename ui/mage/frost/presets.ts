import { Mage } from '../../core/player_classes/mage';
import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, Stat } from '../../core/proto/common';
import { MageMajorGlyph, MageMinorGlyph, FrostMage_Options as MageOptions } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import FrostApl from './apls/frost.apl.json';
import FrostAoeApl from './apls/frost_aoe.apl.json';
import P1FrostGear from './gear_sets/p1_frost_prebis.gear.json';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const FROST_P1_PRESET = PresetUtils.makePresetGear('Frost P1 Preset', P1FrostGear);

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

export const FrostDefaultTalents = {
	name: 'Default Frost',
	data: SavedTalents.create({
		talentsString: '111121',
		glyphs: Glyphs.create({
			major1: MageMajorGlyph.GlyphOfSplittingIce,
			major2: MageMajorGlyph.GlyphOfIcyVeins,
			major3: MageMajorGlyph.GlyphOfWaterElemental,
			minor1: MageMinorGlyph.GlyphOfMomentum,
			minor2: MageMinorGlyph.GlyphOfMirrorImage,
			minor3: MageMinorGlyph.GlyphOfTheUnboundElemental
		}),
	}),
};


export const DefaultFrostOptions = MageOptions.create({
	classOptions: {},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
	tinkerId: 82174, // Synapse Springs
});
export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
