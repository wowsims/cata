import { ConjuredHealthstone, TinkerHandsSynapseSprings } from '../../core/components/inputs/consumables';
import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, RotationType, Spec, Stat } from '../../core/proto/common';
import {
	HunterMajorGlyph as MajorGlyph,
	HunterMinorGlyph as MinorGlyph,
	HunterOptions_Ammo as Ammo,
	HunterOptions_PetType as PetType,
	HunterPrimeGlyph as PrimeGlyph,
	HunterStingType,
	SurvivalHunter_Options as HunterOptions,
	SurvivalHunter_Rotation as HunterRotation,
} from '../../core/proto/hunter';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import { ferocityDefault } from '../../core/talents/hunter_pet';
import AoeApl from './apls/aoe.apl.json';
import SvApl from './apls/sv.apl.json';
import SvAdvApl from './apls/sv_advanced.apl.json';
import P1SVGear from './gear_sets/p1_sv.gear.json';
import PreraidSVGear from './gear_sets/preraid_sv.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const SV_PRERAID_PRESET = PresetUtils.makePresetGear('SV PreRaid Preset', PreraidSVGear);
export const SV_P1_PRESET = PresetUtils.makePresetGear('SV P1 Preset', P1SVGear);
export const DefaultSimpleRotation = HunterRotation.create({
	type: RotationType.SingleTarget,
	sting: HunterStingType.SerpentSting,
	multiDotSerpentSting: true,
});

export const ROTATION_PRESET_SIMPLE_DEFAULT = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecSurvivalHunter, DefaultSimpleRotation);
export const ROTATION_PRESET_SV = PresetUtils.makePresetAPLRotation('SV', SvApl);
export const ROTATION_PRESET_SV_ADVANCED = PresetUtils.makePresetAPLRotation('SV (Advanced)', SvAdvApl);
export const ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('AOE', AoeApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'SV P1',
	Stats.fromMap(
		{
			[Stat.StatStamina]: 0.5,
			[Stat.StatAgility]: 3.27,
			[Stat.StatIntellect]: 1.1,
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

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const SurvivalTalents = {
	name: 'Survival',
	data: SavedTalents.create({
		talentsString: '03-2302-23203003023022121311',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfExplosiveShot,
			prime2: PrimeGlyph.GlyphOfKillShot,
			prime3: PrimeGlyph.GlyphOfArcaneShot,
			major1: MajorGlyph.GlyphOfDisengage,
			major2: MajorGlyph.GlyphOfRaptorStrike,
			major3: MajorGlyph.GlyphOfTrapLauncher,
		}),
	}),
};

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
	duration: 240,
	durationVariation: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Jewelcrafting,
};
