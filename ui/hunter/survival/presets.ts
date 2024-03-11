import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, PetFood, Potions, RotationType, Spec } from '../../core/proto/common';
import {
	HunterMajorGlyph as MajorGlyph,
	HunterMinorGlyph as MinorGlyph,
	HunterOptions_Ammo as Ammo,
	HunterOptions_PetType as PetType,
	HunterStingType,
	SurvivalHunter_Options as HunterOptions,
	SurvivalHunter_Rotation as HunterRotation,
} from '../../core/proto/hunter';
import { SavedTalents } from '../../core/proto/ui';
import { ferocityDefault } from '../../core/talents/hunter_pet';
import AoeApl from './apls/aoe.apl.json';
import SvApl from './apls/sv.apl.json';
import SvAdvApl from './apls/sv_advanced.apl.json';
import P1SVGear from './gear_sets/p1_sv.gear.json';
import P2SVGear from './gear_sets/p2_sv.gear.json';
import P3SVGear from './gear_sets/p3_sv.gear.json';
import P4SVGear from './gear_sets/p4_sv.gear.json';
import P5SVGear from './gear_sets/p5_sv.gear.json';
import PreraidSVGear from './gear_sets/preraid_sv.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const SV_PRERAID_PRESET = PresetUtils.makePresetGear('SV PreRaid Preset', PreraidSVGear);
export const SV_P1_PRESET = PresetUtils.makePresetGear('SV P1 Preset', P1SVGear);
export const SV_P2_PRESET = PresetUtils.makePresetGear('SV P2 Preset', P2SVGear);
export const SV_P3_PRESET = PresetUtils.makePresetGear('SV P3 Preset', P3SVGear);
export const SV_P4_PRESET = PresetUtils.makePresetGear('SV P4 Preset', P4SVGear);
export const SV_P5_PRESET = PresetUtils.makePresetGear('SV P5 Preset', P5SVGear);

export const DefaultSimpleRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: HunterStingType.SerpentSting,
	trapWeave: true,
	viperStartManaPercent: 0.1,
	viperStopManaPercent: 0.3,
	multiDotSerpentSting: true,
	allowExplosiveShotDownrank: true,
});

export const ROTATION_PRESET_SIMPLE_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecSurvivalHunter, DefaultSimpleRotation);
export const ROTATION_PRESET_SV = PresetUtils.makePresetAPLRotation('SV', SvApl);
export const ROTATION_PRESET_SV_ADVANCED = PresetUtils.makePresetAPLRotation('SV (Advanced)', SvAdvApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SurvivalTalents = {
	name: 'Survival',
	data: SavedTalents.create({
		talentsString: '-005305101-5000032500033330531135301331',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSerpentSting,
			major2: MajorGlyph.GlyphOfExplosiveTrap,
			major3: MajorGlyph.GlyphOfKillShot,
			minor1: MinorGlyph.GlyphOfFeignDeath,
			minor2: MinorGlyph.GlyphOfRevivePet,
			minor3: MinorGlyph.GlyphOfMendPet,
		}),
	}),
};

export const SVDefaultOptions = HunterOptions.create({
	classOptions: {
		ammo: Ammo.SaroniteRazorheads,
		useHuntersMark: true,
		petType: PetType.Wolf,
		petTalents: ferocityDefault,
		petUptime: 1,
		timeToTrapWeaveMs: 2000,
	},
	sniperTrainingUptime: 0.9,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
});
