import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Debuffs, Flask, Food, Glyphs, Potions, Profession, PseudoStat, RaidBuffs, Spec, Stat, TinkerHands } from '../../core/proto/common';
import {
	FireMage_Options as MageOptions,
	FireMage_Rotation,
	MageMajorGlyph as MajorGlyph,
	MageMinorGlyph as MinorGlyph,
	MagePrimeGlyph as PrimeGlyph,
} from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats, UnitStat, UnitStatPresets } from '../../core/proto_utils/stats';
import FireApl from './apls/fire.apl.json';
//import FireAoeApl from './apls/fire_aoe.apl.json';
import P1FireBisGear from './gear_sets/p1_fire.gear.json';
import P1FirePrebisGear from './gear_sets/p1_fire_prebis_gear.json';
import P3FireBisGear from './gear_sets/p3_fire.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const FIRE_P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1FireBisGear, { talentTree: 1 });
export const FIRE_P1_PREBIS = PresetUtils.makePresetGear('P1 Pre-raid', P1FirePrebisGear, { talentTree: 1 });
export const FIRE_P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3FireBisGear, { talentTree: 1 });

export const P1DefaultSimpleRotation = FireMage_Rotation.create({
	igniteCombustThreshold: 30000,
	igniteLastMomentLustPercentage: 0.3,
	igniteNoLustPercentage: 0.6,
});

export const P3DefaultSimpleRotation = FireMage_Rotation.create({
	igniteCombustThreshold: 34000,
	igniteLastMomentLustPercentage: 0.33,
	igniteNoLustPercentage: 0.65,
});

export const P1_SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('P1 - Default', Spec.SpecFireMage, P1DefaultSimpleRotation);
export const P3_SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('P3 - Default', Spec.SpecFireMage, P3DefaultSimpleRotation);

//export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFireMage, DefaultSimpleRotation);
export const FIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('APL', FireApl, { talentTree: 1 });

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.33,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatHitRating]: 1.09,
		[Stat.StatCritRating]: 0.62,
		[Stat.StatHasteRating]: 0.82,
		[Stat.StatMasteryRating]: 0.46,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '003-230330221120121213231-03',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfFireball,
			prime2: PrimeGlyph.GlyphOfPyroblast,
			prime3: PrimeGlyph.GlyphOfMoltenArmor,
			major1: MajorGlyph.GlyphOfEvocation,
			major2: MajorGlyph.GlyphOfDragonSBreath,
			major3: MajorGlyph.GlyphOfInvisibility,
			minor1: MinorGlyph.GlyphOfMirrorImage,
			minor2: MinorGlyph.GlyphOfArmors,
			minor3: MinorGlyph.GlyphOfTheMonkey,
		}),
	}),
};

export const DefaultFireOptions = MageOptions.create({
	classOptions: {},
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	markOfTheWild: true,
	icyTalons: true,
	moonkinForm: true,
	leaderOfThePack: true,
	powerWordFortitude: true,
	strengthOfEarthTotem: true,
	trueshotAura: true,
	wrathOfAirTotem: true,
	demonicPact: true,
	blessingOfKings: true,
	blessingOfMight: true,
	communion: true,
});

export const DefaultFireConsumes = Consumes.create({
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const DefaultDebuffs = Debuffs.create({
	ebonPlaguebringer: true,
	shadowAndFlame: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const FIRE_BREAKPOINTS: UnitStatPresets[] = [
	{
		// Picked from Mage Discord
		// Sources:
		// https://docs.google.com/spreadsheets/d/17cJJUReg2uz-XxBB3oDWb1kCncdH_-X96mSb0HAu4Ko/edit?gid=0#gid=0
		// https://docs.google.com/spreadsheets/d/1WLOZ1YevGPw_WZs0JhGzVVy906W5y0i9UqHa3ejyBkE/htmlview?gid=19
		unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
		presets: new Map([
			['11-tick - Combust', 4.98689],
			['5-tick - LvB/Pyro', 12.50704],
			['12-tick - Combust', 15.00864],
			['BL - 15-tick - Combust', 20.86054],
			['13-tick - Combust', 25.07819],
			['BL - 16-tick - Combust', 29.09891],
			['14-tick - Combust', 35.04391],
			['BL - 17-tick - Combust', 37.40041],
			['6-tick - LvB/Pyro', 37.52006],
			['15-tick - Combust', 45.03265],
			['16-tick - Combust', 54.91869],
			['7-tick - LvB/Pyro', 62.46955],
			['17-tick - Combust', 64.88049],
			['18-tick - Combust', 74.97816],
			['19-tick - Combust', 85.01391],
			['8-tick - LvB/Pyro', 87.44144],
			['20-tick - Combust', 95.12199],
			['21-tick - Combust', 105.12825],
			// ['9-tick - LvB/Pyro', 112.53987],
			// ['22-tick - Combust', 114.82282],
			// ['23-tick - Combust', 124.97193],
			// ['24-tick - Combust', 135.01768],
			// ['10-tick - LvB/Pyro', 137.43571],
			// ['25-tick - Combust', 144.7981],
			// ['26-tick - Combust', 154.77713],
			// ['11-tick - LvB/Pyro', 162.58208],
			// ['27-tick - Combust', 164.90073],
			// ['28-tick - Combust', 175.10324],
			// ['29-tick - Combust', 185.30679],
			// ['12-tick - LvB/Pyro', 187.49404],
		]),
	},
];

export const P1_PRESET_BUILD = PresetUtils.makePresetBuild('P1 - Default', {
	gear: FIRE_P1_PRESET,
	rotation: P1_SIMPLE_ROTATION_DEFAULT,
});

export const P3_PRESET_BUILD = PresetUtils.makePresetBuild('P3 - Default', {
	gear: FIRE_P3_PRESET,
	rotation: P3_SIMPLE_ROTATION_DEFAULT,
});
