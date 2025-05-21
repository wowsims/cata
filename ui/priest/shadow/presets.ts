import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, IndividualBuffs, Profession, PseudoStat, Race, RaidBuffs, Stat } from '../../core/proto/common';
import { PriestMajorGlyph as MajorGlyph, PriestMinorGlyph as MinorGlyph, PriestOptions_Armor, ShadowPriest_Options as Options } from '../../core/proto/priest';
import { SavedTalents } from '../../core/proto/ui';
import { Stats, UnitStat, UnitStatPresets } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import PreRaidGear from './gear_sets/pre_raid.gear.json';
import P1Gear from './gear_sets/p1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const PRE_RAID_PRESET = PresetUtils.makePresetGear('Pre Raid Preset', PreRaidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
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
// https://www.wowhead.com/mop-classic/talent-calc/priest and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '223113',
		glyphs: Glyphs.create({}),
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
	blessingOfKings: true,
	mindQuickening: true,
	leaderOfThePack: true,
	blessingOfMight: true,
	unholyAura: true,
	bloodlust: true,
	skullBannerCount: 2,
	stormlashTotemCount: 4,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true
});

export const OtherDefaults = {
	channelClipDelay: 40,
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
