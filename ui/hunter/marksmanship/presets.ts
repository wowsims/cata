import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, PetFood, Potions, RotationType, Spec } from '../../core/proto/common';
import {
	BeastMasteryHunter_Options as HunterOptions,
	BeastMasteryHunter_Rotation as HunterRotation,
	HunterMajorGlyph as MajorGlyph,
	HunterMinorGlyph as MinorGlyph,
	HunterOptions_Ammo as Ammo,
	HunterOptions_PetType as PetType,
	HunterStingType,
} from '../../core/proto/hunter';
import { SavedTalents } from '../../core/proto/ui';
import { ferocityDefault } from '../../core/talents/hunter_pet';
import AoeApl from './apls/aoe.apl.json';
import MmApl from './apls/mm.apl.json';
import MmAdvApl from './apls/mm_advanced.apl.json';
import P1MMGear from './gear_sets/p1_mm.gear.json';
import PreraidMMGear from './gear_sets/preraid_mm.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const MM_PRERAID_PRESET = PresetUtils.makePresetGear('MM PreRaid Preset', PreraidMMGear);
export const MM_P1_PRESET = PresetUtils.makePresetGear('MM P1 Preset', P1MMGear);

export const DefaultSimpleRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: HunterStingType.SerpentSting,
	trapWeave: true,
	multiDotSerpentSting: true,
	allowExplosiveShotDownrank: true,
});

export const ROTATION_PRESET_SIMPLE_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecMarksmanshipHunter, DefaultSimpleRotation);
export const ROTATION_PRESET_MM = PresetUtils.makePresetAPLRotation('MM', MmApl);
export const ROTATION_PRESET_MM_ADVANCED = PresetUtils.makePresetAPLRotation('MM (Advanced)', MmAdvApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const MarksmanTalents = {
	name: 'Marksman',
	data: SavedTalents.create({
		// talentsString: '502-025335101030013233135031051-5000032',
		// glyphs: Glyphs.create({
		// 	major1: MajorGlyph.GlyphOfSerpentSting,
		// 	major2: MajorGlyph.GlyphOfSteadyShot,
		// 	major3: MajorGlyph.GlyphOfExplosiveTrap,
		// 	minor1: MinorGlyph.GlyphOfFeignDeath,
		// 	minor2: MinorGlyph.GlyphOfRevivePet,
		// 	minor3: MinorGlyph.GlyphOfMendPet,
		// }),
	}),
};

export const MMDefaultOptions = HunterOptions.create({
	classOptions: {
		ammo: Ammo.SaroniteRazorheads,
		useHuntersMark: true,
		petType: PetType.Wolf,
		petTalents: ferocityDefault,
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
