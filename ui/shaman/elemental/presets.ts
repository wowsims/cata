import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Glyphs, Profession, Stat } from '../../core/proto/common.js';
import {
	ElementalShaman_Options as ElementalShamanOptions,
	FeleAutocastSettings,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanShield,
} from '../../core/proto/shaman.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import AoEApl from './apls/aoe.apl.json';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 - Default', P1Gear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AoE', AoEApl);

// Preset options for EP weights
export const EP_PRESET_DEFAULT = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.24,
		[Stat.StatSpellPower]: 1,
		[Stat.StatCritRating]: 0.34,
		[Stat.StatHasteRating]: 0.57,
		[Stat.StatHitRating]: 0.59,
		[Stat.StatSpirit]: 0.59,
		[Stat.StatMasteryRating]: 0.49,
	}),
);

export const EP_PRESET_CLEAVE = PresetUtils.makePresetEpWeights(
	'Cleave/AoE',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.33,
		[Stat.StatSpellPower]: 1,
		[Stat.StatCritRating]: 0.54,
		[Stat.StatHasteRating]: 0.57,
		[Stat.StatHitRating]: 1.09,
		[Stat.StatSpirit]: 1.09,
		[Stat.StatMasteryRating]: 1,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Totem Duration',
	data: SavedTalents.create({
		talentsString: '313232',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfLightningShield,
			major2: ShamanMajorGlyph.GlyphOfHealingStreamTotem,
			minor1: ShamanMinorGlyph.GlyphOfThunderstorm,
		}),
	}),
};

export const TalentsAoE = {
	name: 'AoE (4+)',
	data: SavedTalents.create({
		...StandardTalents.data,
		glyphs: Glyphs.create({
			...StandardTalents.data.glyphs,
			major2: ShamanMajorGlyph.GlyphOfChainLightning,
		}),
	}),
};

export const DefaultOptions = ElementalShamanOptions.create({
	classOptions: {
		shield: ShamanShield.LightningShield,
		feleAutocast: FeleAutocastSettings.create({
					autocastFireblast: true,
					autocastFirenova: true,
					autocastImmolate: true,
					autocastEmpower: false,
				}),
	},
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
});

const ENCOUNTER_SINGLE_TARGET = PresetUtils.makePresetEncounter(
	'Single Target Dummy',
	'http://localhost:5173/mop/shaman/elemental/?i=e#eJyTYhJgVzrO5MHMyQACB4ocBMEMBj8HyVkzQeCkvSVE5IK9YhoYXLN3PHsGBN7YG/UwFaz6zFjFHZSYmaIQkliUnlqiECHBrnWDkYEeIKDFgZrGNaQcRzbPx2LuHEd0NeDQaFjE6TiTERJAN+2halgcABeZKLc=',
);

const ENCOUNTER_CLEAVE = PresetUtils.makePresetEncounter(
	'Cleave',
	'http://localhost:5173/mop/shaman/elemental/?i=e#eJyTYhJgV7rG4sHMyQACB4ocBMEMBj8HyVkzQeCkvSVE5IK9YhoYXLN3PHsGBN7YG/UwFaz6zFjFHZSYmaIQkliUnlqiECHBrnWDkYEeIKDFgZrGNaQcRzbPx2LuHEd0NeDQaFjE6TiTERJAN+2halgcRoMDJTgATH8+LA==',
);

const ENCOUNTER_AOE = PresetUtils.makePresetEncounter(
	'AOE (4+)',
	'http://localhost:5173/mop/shaman/elemental/?i=e#eJyTYhJgV/rC4cHMyQACDWkOgmAGg4iD5KyZIHDS3hIicsFeMQ0Mrtk7nj0DAm/sjXqYClZ9ZqziDkrMTFEISSxKTy1RiJBg17rByEAPENDiQE3jGlKOI5vnYzF3jiO6GnBoNCzidJzJCAmgm/ZQNSwOo8ExGhwMo8GBAPiCAwDGf2iQ',
);

export const P3_PRESET_BUILD_DEFAULT = PresetUtils.makePresetBuild('Default', {
	talents: StandardTalents,
	rotation: ROTATION_PRESET_DEFAULT,
	encounter: ENCOUNTER_SINGLE_TARGET,
});

export const P3_PRESET_BUILD_CLEAVE = PresetUtils.makePresetBuild('Cleave', {
	talents: StandardTalents,
	rotation: ROTATION_PRESET_AOE,
	encounter: ENCOUNTER_CLEAVE,
});

export const P3_PRESET_BUILD_AOE = PresetUtils.makePresetBuild('AoE (4+)', {
	talents: TalentsAoE,
	rotation: ROTATION_PRESET_AOE,
	encounter: ENCOUNTER_AOE,
});
