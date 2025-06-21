import * as PresetUtils from '../../core/preset_utils';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { HunterMajorGlyph as MajorGlyph, HunterOptions_PetType as PetType, SurvivalHunter_Options as HunterOptions } from '../../core/proto/hunter';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import P1Gear from '../presets/p1.json';
import PreRaidGear from '../presets/preraid.json';
import PreRaidGearCelestial from '../presets/preraid_celestial.json';
import AoeApl from './apls/aoe.apl.json';
import Apl from './apls/mm.apl.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET_GEAR = PresetUtils.makePresetGear('Pre-raid', PreRaidGear);
export const PRERAID_CELESTIAL_PRESET_GEAR = PresetUtils.makePresetGear('Pre-raid', PreRaidGearCelestial);
export const P1_PRESET_GEAR = PresetUtils.makePresetGear('P1', P1Gear);
export const ROTATION_PRESET_SV = PresetUtils.makePresetAPLRotation('Single Target', Apl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);
export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '312111',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfAnimalBond,
			major2: MajorGlyph.GlyphOfDeterrence,
			major3: MajorGlyph.GlyphOfAimedShot,
		}),
	}),
};
// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStamina]: 0.0,
			[Stat.StatAgility]: 1,
			[Stat.StatHitRating]: 0.63,
			[Stat.StatCritRating]: 0.4,
			[Stat.StatHasteRating]: 0.35,
			[Stat.StatMasteryRating]: 0.29,
			[Stat.StatExpertiseRating]: 0.59,
		},
		{
			[PseudoStat.PseudoStatRangedDps]: 0.62,
		},
	),
);

export const PRERAID_PRESET = PresetUtils.makePresetBuild('Pre-raid', {
	gear: PRERAID_PRESET_GEAR,
	epWeights: P1_EP_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAPL,
});
export const PRERAID_PRESET_CELESTIAL = PresetUtils.makePresetBuild('Pre-raid (Celestial)', {
	gear: PRERAID_CELESTIAL_PRESET_GEAR,
	epWeights: P1_EP_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAPL,
});
export const P1_PRESET = PresetUtils.makePresetBuild('P1', {
	gear: P1_PRESET_GEAR,
	epWeights: P1_EP_PRESET as PresetUtils.PresetEpWeights,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAPL,
});
// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SVDefaultOptions = HunterOptions.create({
	classOptions: {
		useHuntersMark: true,
		petType: PetType.Wolf,
		petUptime: 1,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76084, // Flask of the Winds
	foodId: 74648, // Seafood Magnifique Feast
	potId: 76089, // Potion of the Tol'vir
	prepotId: 76089, // Potion of the Tol'vir
	conjuredId: 5512, // Conjured Healthstone
});
export const OtherDefaults = {
	distanceFromTarget: 24,
	iterationCount: 25000,
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
	GlaiveTossChance: 80,
};
