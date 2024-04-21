import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Spec, TinkerHands } from '../../core/proto/common';
import {
	DruidPrimeGlyph,
	DruidMajorGlyph,
	DruidMinorGlyph,
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Rotation_AplType,
	FeralDruid_Rotation_BearweaveType,
	FeralDruid_Rotation_BiteModeType,
} from '../../core/proto/druid';
import { SavedTalents } from '../../core/proto/ui';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid Preset', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

import DefaultApl from './apls/default.apl.json';
export const APL_ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);

import CustomExampleApl from './apls/custom_apl_example.apl.json';
export const APL_ROTATION_CUSTOM_EXAMPLE = PresetUtils.makePresetAPLRotation('Custom APL Example', CustomExampleApl);

export const DefaultRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.SingleTarget,

	bearWeaveType: FeralDruid_Rotation_BearweaveType.None,
	minCombosForRip: 5,
	minCombosForBite: 5,

	useRake: true,
	useBite: true,
	mangleSpam: false,
	biteModeType: FeralDruid_Rotation_BiteModeType.Emperical,
	biteTime: 4.0,
	berserkBiteThresh: 25.0,
	powerbear: false,
	minRoarOffset: 12.0,
	ripLeeway: 3.0,
	maintainFaerieFire: true,
	snekWeave: false,
});

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFeralDruid, DefaultRotation);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-2320322312012121202301-020301',
		glyphs: Glyphs.create({
			prime1: DruidPrimeGlyph.GlyphOfRip,
			prime2: DruidPrimeGlyph.GlyphOfBloodletting,
			prime3: DruidPrimeGlyph.GlyphOfBerserk,
			major1: DruidMajorGlyph.GlyphOfThorns,
			major2: DruidMajorGlyph.GlyphOfFeralCharge,
			major3: DruidMajorGlyph.GlyphOfRebirth,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfMarkOfTheWild,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultOptions = FeralDruidOptions.create({
	assumeBleedActive: true,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodSkeweredEel,
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});
