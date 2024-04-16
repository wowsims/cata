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
import P1BMGear from './gear_sets/p1_bm.gear.json';
import PreraidBMGear from './gear_sets/preraid_bm.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BM_PRERAID_PRESET = PresetUtils.makePresetGear('BM PreRaid Preset', PreraidBMGear);
export const BM_P1_PRESET = PresetUtils.makePresetGear('BM P1 Preset', P1BMGear);

export const DefaultSimpleRotation = BeastMasteryRotation.create({
	type: RotationType.SingleTarget,
	sting: StingType.SerpentSting,
	trapWeave: true,
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
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodFishFeast,
});
