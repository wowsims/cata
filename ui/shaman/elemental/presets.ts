import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Glyphs, Profession, Race, Stat } from '../../core/proto/common.js';
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
import CleaveApl from './apls/cleave.apl.json';
import PEApl from './apls/pe.apl.json';
import UFApl from './apls/uf.apl.json';
import EBApl from './apls/eb.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 - Default', P1Gear);

export const ROTATION_PRESET_UF = PresetUtils.makePresetAPLRotation('Default', UFApl);
export const ROTATION_PRESET_EB = PresetUtils.makePresetAPLRotation('Elemental Blast', EBApl);
export const ROTATION_PRESET_PE = PresetUtils.makePresetAPLRotation('Primal Elementalist', PEApl);
export const ROTATION_PRESET_CLEAVE = PresetUtils.makePresetAPLRotation('Cleave', CleaveApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AoE (3+)', AoEApl);

// Preset options for EP weights
export const EP_PRESET_DEFAULT = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.00,
		[Stat.StatSpellPower]: 0.80,
		[Stat.StatCritRating]: 0.20,
		[Stat.StatHasteRating]: 0.40,
		[Stat.StatHitRating]: 0.60,
		[Stat.StatSpirit]: 0.60,
		[Stat.StatMasteryRating]: 0.30,
	}),
);

export const EP_PRESET_AOE = PresetUtils.makePresetEpWeights(
	'AoE (4+)',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.00,
		[Stat.StatSpellPower]: 0.80,
		[Stat.StatCritRating]: 0.30,
		[Stat.StatHasteRating]: 0.20,
		[Stat.StatHitRating]: 0.60,
		[Stat.StatSpirit]: 0.60,
		[Stat.StatMasteryRating]: 0.40,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '333121',
		glyphs: Glyphs.create({
			major1: ShamanMajorGlyph.GlyphOfSpiritwalkersGrace,
		}),
	}),
};

export const TalentsCleave = {
	name: 'Cleave',
	data: SavedTalents.create({
		talentsString: '333322',
		glyphs: Glyphs.create({
			...StandardTalents.data.glyphs,
		}),
	}),
};

export const TalentsAoE = {
	name: 'AoE (4+)',
	data: SavedTalents.create({
		...TalentsCleave.data,
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
	'http://localhost:5173/mop/shaman/elemental/?i=e#eJyTYhJgV5rP5MHIyQACB4ocBMEMBj8HyVkzQeCkvSVE5IK9YhoYXLN3PHsGBN7YGz1hLFj1mbGKOygxM0UhJLEoPbVEIVaCXWsDI8MgAgccbjggcX0s5s5xRFcD9uWDEwqOMxkhHr9pD1PjAAAuwChz',
);

const ENCOUNTER_CLEAVE = PresetUtils.makePresetEncounter(
	'Cleave',
	'http://localhost:5173/mop/shaman/elemental/?i=e#eJyTYhJgV2pj8WDkZACBA0UOgmAGg5+D5KyZIHDS3hIicsFeMQ0Mrtk7nj0DAm/sjZ4wFqz6zFjFHZSYmaIQkliUnlqiECvBrrWBkWEQgQMONxyQuD4Wc+c4oqsB+/LBCQXHmYwQj9+0h6lxGGhvBrQ4EKOMUm8CACyAPro=',
);

const ENCOUNTER_AOE = PresetUtils.makePresetEncounter(
	'AOE (5+)',
	'http://localhost:5173/mop/shaman/elemental/?i=e#eJyTYhJgV7rC7sHIyQACB4ocBMEMBj8HyVkzQeCkvSVE5IK9YhoYXLN3PHsGBN7YGz1hLFj1mbGKOygxM0UhJLEoPbVEIVaCXWsDI8MgAgccbjggcX0s5s5xRFcD9uWDEwqOMxkhHr9pD1PjMNDeDGhxIEbZqDfhYNSbDgAivWvH',
);

export const P1_PRESET_BUILD_DEFAULT = PresetUtils.makePresetBuild('Default', {
	talents: StandardTalents,
	rotation: ROTATION_PRESET_UF,
	encounter: ENCOUNTER_SINGLE_TARGET,
	epWeights: EP_PRESET_DEFAULT,
});

export const P1_PRESET_BUILD_CLEAVE = PresetUtils.makePresetBuild('Cleave', {
	talents: TalentsCleave,
	rotation: ROTATION_PRESET_CLEAVE,
	encounter: ENCOUNTER_CLEAVE,
	epWeights: EP_PRESET_DEFAULT,
});

export const P1_PRESET_BUILD_AOE = PresetUtils.makePresetBuild('AoE (4+)', {
	talents: TalentsAoE,
	rotation: ROTATION_PRESET_AOE,
	encounter: ENCOUNTER_AOE,
	epWeights: EP_PRESET_AOE,
});
