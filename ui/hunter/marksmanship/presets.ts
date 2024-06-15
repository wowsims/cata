import { ConjuredHealthstone, TinkerHandsSynapseSprings } from '../../core/components/inputs/consumables';
import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, RotationType, Spec, Stat } from '../../core/proto/common';
import {
	BeastMasteryHunter_Options as HunterOptions,
	BeastMasteryHunter_Rotation as HunterRotation,
	HunterMajorGlyph as MajorGlyph,
	HunterOptions_PetType as PetType,
	HunterPrimeGlyph as PrimeGlyph,
	HunterStingType,
} from '../../core/proto/hunter';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
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

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'MM P1',
	Stats.fromMap(
		{
			[Stat.StatStamina]: 0.5,
			[Stat.StatAgility]: 3.05,
			[Stat.StatRangedAttackPower]: 1.0,
			[Stat.StatMeleeHit]: 2.25,
			[Stat.StatMeleeCrit]: 1.39,
			[Stat.StatMeleeHaste]: 1.33,
			[Stat.StatMastery]: 1.15,
		},
		{
			[PseudoStat.PseudoStatRangedDps]: 6.32,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const MarksmanTalents = {
	name: 'Marksman',
	data: SavedTalents.create({
		talentsString: '032002-2302320232120231201-03',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfArcaneShot,
			prime2: PrimeGlyph.GlyphOfRapidFire,
			prime3: PrimeGlyph.GlyphOfSteadyShot,
			major1: MajorGlyph.GlyphOfDisengage,
			major2: MajorGlyph.GlyphOfRaptorStrike,
			major3: MajorGlyph.GlyphOfTrapLauncher,
		}),
	}),
};

export const MMDefaultOptions = HunterOptions.create({
	classOptions: {
		useHuntersMark: true,
		petType: PetType.Wolf,
		petTalents: ferocityDefault,
		petUptime: 1,
	},
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
