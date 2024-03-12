import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, PetFood, Potions, RotationType, Spec } from '../../core/proto/common';
import {
	BeastMasteryHunter_Options as BeastMasteryOptions,
	BeastMasteryHunter_Rotation as BeastMasteryRotation,
	HunterMajorGlyph as MajorGlyph,
	HunterMinorGlyph as MinorGlyph,
	HunterOptions_Ammo as Ammo,
	HunterOptions_PetType as PetType,
	HunterStingType as StingType,
} from '../../core/proto/hunter';
import { SavedTalents } from '../../core/proto/ui';
import { ferocityBMDefault } from '../../core/talents/hunter_pet';
import AoeApl from './apls/aoe.apl.json';
import BmApl from './apls/bm.apl.json';
import P1MMGear from './gear_sets/p1_mm.gear.json';
import P2MMGear from './gear_sets/p2_mm.gear.json';
import P3MMGear from './gear_sets/p3_mm.gear.json';
import P4MMGear from './gear_sets/p4_mm.gear.json';
import P5MMGear from './gear_sets/p5_mm.gear.json';
import PreraidMMGear from './gear_sets/preraid_mm.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BM_PRERAID_PRESET = PresetUtils.makePresetGear('BM PreRaid Preset', PreraidMMGear);
export const BM_P1_PRESET = PresetUtils.makePresetGear('BM P1 Preset', P1MMGear);
export const BM_P2_PRESET = PresetUtils.makePresetGear('BM P2 Preset', P2MMGear);
export const BM_P3_PRESET = PresetUtils.makePresetGear('BM P3 Preset', P3MMGear);
export const BM_P4_PRESET = PresetUtils.makePresetGear('BM P4 Preset', P4MMGear);
export const BM_P5_PRESET = PresetUtils.makePresetGear('BM P5 Preset', P5MMGear);

export const DefaultSimpleRotation = BeastMasteryRotation.create({
	type: RotationType.SingleTarget,
	sting: StingType.SerpentSting,
	trapWeave: true,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,
	multiDotSerpentSting: true,
	allowExplosiveShotDownrank: true,
});

export const ROTATION_PRESET_SIMPLE_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecBeastMasteryHunter, DefaultSimpleRotation);
export const ROTATION_PRESET_BM = PresetUtils.makePresetAPLRotation('BM', BmApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const BeastMasteryTalents = {
	name: 'Beast Mastery',
	data: SavedTalents.create({
		// talentsString: '51200201505112243120531251-025305101',
		// glyphs: Glyphs.create({
		// 	major1: MajorGlyph.GlyphOfBestialWrath,
		// 	major2: MajorGlyph.GlyphOfSteadyShot,
		// 	major3: MajorGlyph.GlyphOfSerpentSting,
		// 	minor1: MinorGlyph.GlyphOfFeignDeath,
		// 	minor2: MinorGlyph.GlyphOfRevivePet,
		// 	minor3: MinorGlyph.GlyphOfMendPet,
		// }),
	}),
};

export const BMDefaultOptions = BeastMasteryOptions.create({
	classOptions: {
		ammo: Ammo.SaroniteRazorheads,
		useHuntersMark: true,
		petType: PetType.Wolf,
		petTalents: ferocityBMDefault,
		petUptime: 1,
		timeToTrapWeaveMs: 2000,
	},
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
});
