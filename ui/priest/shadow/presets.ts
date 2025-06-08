import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, IndividualBuffs, Profession, PseudoStat, Race, RaidBuffs, Stat } from '../../core/proto/common';
import {
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
	PriestOptions_Armor,
	PriestPrimeGlyph as PrimeGlyph,
	ShadowPriest_Options as Options,
} from '../../core/proto/priest';
import { SavedTalents } from '../../core/proto/ui';
import { Stats, UnitStat, UnitStatPresets } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P4T134PCApl from './apls/p4.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P3Gear from './gear_sets/p3.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import ItemSwapP4 from './gear_sets/p4_item_swap.gear.json';
import PreRaidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const PRE_RAID = PresetUtils.makePresetGear('Pre Raid', PreRaidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

export const P4_ITEM_SWAP = PresetUtils.makePresetItemSwapGear('P4', ItemSwapP4);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const P4_T13_4PC_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('T13 - 4PC', P4T134PCApl);

// Preset options for EP weights
export const P3_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.0,
		[Stat.StatSpirit]: 0.9,
		[Stat.StatSpellPower]: 0.79,
		[Stat.StatHitRating]: 0.85,
		[Stat.StatCritRating]: 0.42,
		[Stat.StatHasteRating]: 0.76,
		[Stat.StatMasteryRating]: 0.48,
	}),
);

export const SHADOW_BREAKPOINTS: UnitStatPresets[] = [
	{
		// Picked from Priest Discord
		// Sources:
		// https://docs.google.com/spreadsheets/d/17cJJUReg2uz-XxBB3oDWb1kCncdH_-X96mSb0HAu4Ko/edit?usp=sharing
		// https://docs.google.com/spreadsheets/d/1WLOZ1YevGPw_WZs0JhGzVVy906W5y0i9UqHa3ejyBkE/htmlview?gid=19
		unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
		presets: new Map([
			['9-tick - DP', 6.25111],
			['7-tick - SWP', 8.32281],
			['BL - 9-tick - SWP', 8.98193],
			['6-tick - VT', 9.99084],
			['BL - 12-tick - DP', 10.60112],
			['BL - 8-tick - VT', 15.35578],
			['10-tick - DP', 18.74135],
			['BL - 13-tick - DP', 20.22362],
			['BL - 10-tick - SWP', 21.8101],
			['8-tick - SWP', 24.97397],
			['BL - 14-tick - DP', 29.82799],
			['7-tick - VT', 30.01084],
			['BL - 9-tick - VT', 30.7845],
			['11-tick - DP', 31.26231],
			['BL - 11-tick - SWP', 34.59857],
			['BL - 15-tick - DP', 39.3955],
			['9-tick - SWP', 41.67651],
			['12-tick - DP', 43.78146],
			['BL - 10-tick - VT', 46.19528],
			['BL - 12-tick - SWP', 47.40929],
			['BL - 16-tick - DP', 49.0276],
			['8-tick - VT', 49.96252],
			['13-tick - DP', 56.29071],
			['10-tick - SWP', 58.35314],
			// ['BL - 17-tick - DP', 58.65882],
			// ['BL - 13-tick - SWP', 60.31209],
			// ['BL - 11-tick - VT', 61.54655],
			// ['BL - 18-tick - DP', 68.26048],
			// ['14-tick - DP', 68.77638],
			// ['9-tick - VT', 70.01985],
			// ['BL - 14-tick - SWP', 73.0553],
			// ['11-tick - SWP', 74.97814],
			// ['BL - 12-tick - VT', 76.90245],
			// ['BL - 19-tick - DP', 77.85684],
			// ['15-tick - DP', 81.21415],
			// ['BL - 15-tick - SWP', 85.87938],
			// ['BL - 20-tick - DP', 87.54104],
			// ['10-tick - VT', 90.05386],
			// ['12-tick - SWP', 91.63208],
			// ['BL - 13-tick - VT', 92.38787],
			// ['16-tick - DP', 93.73589],
			// ['BL - 21-tick - DP', 97.15442],
			// ['BL - 16-tick - SWP', 98.68209],
			// ['17-tick - DP', 106.25646],
			// ['BL - 22-tick - DP', 106.68988],
			// ['BL - 14-tick - VT', 107.61966],
			// ['13-tick - SWP', 108.40571],
			// ['11-tick - VT', 110.01052],
			// ['BL - 17-tick - SWP', 111.61784],
			// ['BL - 23-tick - DP', 116.37998],
			// ['18-tick - DP', 118.73862],
			// ['BL - 15-tick - VT', 123.07323],
			// ['BL - 18-tick - SWP', 124.37458],
			// ['14-tick - SWP', 124.9719],
			// ['12-tick - VT', 129.97319],
			// ['19-tick - DP', 131.21389],
			// ['15-tick - SWP', 141.64319],
			// ['20-tick - DP', 143.80335],
			// ['13-tick - VT', 150.10423],
			// ['21-tick - DP', 156.30075],
			// ['16-tick - SWP', 158.28672],
			// ['22-tick - DP', 168.69684],
			// ['14-tick - VT', 169.90556],
			// ['17-tick - SWP', 175.10319],
			// ['23-tick - DP', 181.29398],
			// ['15-tick - VT', 189.99519],
			// ['18-tick - SWP', 191.68695],
		]),
	},
];

// Default talents. Uses the wowhead calculator format, make the talents on
// https://www.wowhead.com/cata/talent-calc/priest and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '032212--322032210201222100231',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfShadowWordPain,
			prime2: PrimeGlyph.GlyphOfMindFlay,
			prime3: PrimeGlyph.GlyphOfShadowWordDeath,
			major1: MajorGlyph.GlyphOfFade,
			major2: MajorGlyph.GlyphOfInnerFire,
			major3: MajorGlyph.GlyphOfSpiritTap,
			minor1: MinorGlyph.GlyphOfFading,
			minor2: MinorGlyph.GlyphOfFortitude,
			minor3: MinorGlyph.GlyphOfShadowfiend,
		}),
	}),
};

export const DefaultOptions = Options.create({
	classOptions: {
		armor: PriestOptions_Armor.InnerFire,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58086, // Flask of the Draconic Mind
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
	prepotId: 58091, // Volcanic Potion
	tinkerId: 82174, // Synapse Springs
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

export const DefaultIndividualBuffs = IndividualBuffs.create({
	vampiricTouch: true,
	darkIntent: true,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
});

export const OtherDefaults = {
	channelClipDelay: 40,
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	darkIntentUptime: 90,
};

export const P3_PRESET_BUILD = PresetUtils.makePresetBuild('P3 - Default', {
	race: Race.RaceTroll,
	gear: P3_PRESET,
	rotation: ROTATION_PRESET_DEFAULT,
});

export const P4_PRESET_BUILD = PresetUtils.makePresetBuild('P4 - Default', {
	race: Race.RaceTroll,
	gear: P4_PRESET,
	rotation: P4_T13_4PC_PRESET_DEFAULT,
});
