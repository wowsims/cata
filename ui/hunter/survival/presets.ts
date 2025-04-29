import { ConjuredHealthstone, TinkerHandsSynapseSprings } from '../../core/components/inputs/consumables';
import * as PresetUtils from '../../core/preset_utils';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, RotationType, Spec, Stat } from '../../core/proto/common';
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
import { Stats } from '../../core/proto_utils/stats';
import { ferocityDefault } from '../../core/talents/hunter_pet';
import AoeApl from './apls/aoe.apl.json';
import SvApl from './apls/sv.apl.json';
import P1SVGear from './gear_sets/p1_sv.gear.json';
import P3SVGear from './gear_sets/p3_sv.gear.json';
import P4SVGear from './gear_sets/p4_sv.gear.json';
import PreraidSVGear from './gear_sets/preraid_sv.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const SV_PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidSVGear);
export const SV_P1_PRESET = PresetUtils.makePresetGear('P2', P1SVGear);
export const SV_P3_PRESET = PresetUtils.makePresetGear('P3', P3SVGear);
export const SV_P4_PRESET = PresetUtils.makePresetGear('P4', P4SVGear);
export const ROTATION_PRESET_SV = PresetUtils.makePresetAPLRotation('SV', SvApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);
export const SurvivalTalents = {
	name: 'Survival',
	data: SavedTalents.create({
		talentsString: '03-2302-03203203023022121311',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfDisengage,
			major2: MajorGlyph.GlyphOfRaptorStrike,
			major3: MajorGlyph.GlyphOfTrapLauncher,
		}),
	}),
};
// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P2',
	Stats.fromMap(
		{
			[Stat.StatStamina]: 0.5,
			[Stat.StatAgility]: 3.27,
			[Stat.StatRangedAttackPower]: 1.0,
			[Stat.StatHitRating]: 2.16,
			[Stat.StatCritRating]: 1.17,
			[Stat.StatHasteRating]: 0.89,
			[Stat.StatMasteryRating]: 0.88,
		},
		{
			[PseudoStat.PseudoStatRangedDps]: 3.75,
		},
	),
);

export const P3_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P3',
	Stats.fromMap(
		{
			[Stat.StatStamina]: 0.5,
			[Stat.StatAgility]: 3.37,
			[Stat.StatRangedAttackPower]: 1.0,
			[Stat.StatHitRating]: 2.56,
			[Stat.StatCritRating]: 1.27,
			[Stat.StatHasteRating]: 1.09,
			[Stat.StatMasteryRating]: 1.04,
		},
		{
			[PseudoStat.PseudoStatRangedDps]: 4.16,
		},
	),
);

export const P4_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P4',
	Stats.fromMap(
		{
			[Stat.StatStamina]: 0.5,
			[Stat.StatAgility]: 3.47,
			[Stat.StatRangedAttackPower]: 1.0,
			[Stat.StatHitRating]: 2.56,
			[Stat.StatCritRating]: 1.45,
			[Stat.StatHasteRating]: 1.09,
			[Stat.StatMasteryRating]: 1.04,
		},
		{
			[PseudoStat.PseudoStatRangedDps]: 4.16,
		},
	),
);

export const P2_PRESET = PresetUtils.makePresetBuild('P2', {
	gear: SV_P1_PRESET,
	epWeights: P1_EP_PRESET,
	talents: SurvivalTalents,
	rotationType: APLRotationType.TypeAuto,
});

export const P3_PRESET = PresetUtils.makePresetBuild('P3', {
	gear: SV_P3_PRESET,
	epWeights: P3_EP_PRESET,
	talents: SurvivalTalents,
	rotationType: APLRotationType.TypeAuto,
});

export const P4_PRESET = PresetUtils.makePresetBuild('P4', {
	gear: SV_P4_PRESET,
	epWeights: P4_EP_PRESET,
	talents: SurvivalTalents,
	rotationType: APLRotationType.TypeAuto,
});
// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.



export const SVDefaultOptions = HunterOptions.create({
	classOptions: {
		useHuntersMark: true,
		petType: PetType.Wolf,
		petTalents: ferocityDefault,
		petUptime: 1,
	},
	sniperTrainingUptime: 0.9,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	flask: Flask.FlaskOfTheWinds,
	defaultConjured: ConjuredHealthstone.value,
	food: Food.FoodSeafoodFeast,
	tinkerHands: TinkerHandsSynapseSprings.value,
});

export const OtherDefaults = {
	distanceFromTarget: 24,
	profession1: Profession.Engineering,
	profession2: Profession.Jewelcrafting,
};
