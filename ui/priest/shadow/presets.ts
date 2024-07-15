import * as PresetUtils from '../../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	Stat,
	TinkerHands,
	TristateEffect,
} from '../../core/proto/common.js';
import {
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
	PriestOptions_Armor,
	PriestPrimeGlyph as PrimeGlyph,
	ShadowPriest_Options as Options,
} from '../../core/proto/priest.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import PreRaidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const PRE_RAID = PresetUtils.makePresetGear('Pre Raid', PreRaidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.0,
		[Stat.StatSpirit]: 0.9,
		[Stat.StatSpellPower]: 0.79,
		[Stat.StatSpellHit]: 0.85,
		[Stat.StatSpellCrit]: 0.42,
		[Stat.StatSpellHaste]: 0.76,
		[Stat.StatMastery]: 0.48,
	}),
);

export const SHADOW_BREAKPOINTS = new Map([
	[
		Stat.StatSpellHaste,
		// Picked from Priest Discord
		// Sources:
		// https://docs.google.com/spreadsheets/d/17cJJUReg2uz-XxBB3oDWb1kCncdH_-X96mSb0HAu4Ko/edit?usp=sharing
		// https://docs.google.com/spreadsheets/d/1WLOZ1YevGPw_WZs0JhGzVVy906W5y0i9UqHa3ejyBkE/htmlview?gid=19
		new Map([
			['7-tick - SWP', 1066],
			['6-tick - VT', 1280],
			['BL - 12-tick - DP', 1358],
			['BL - 8-tick - VT', 1967],
			['9-tick - DP', 2400],
			['BL - 13-tick - DP', 2590],
			['BL - 10-tick - SWP', 2793],
			['8-tick - SWP', 3199],
			['BL - 14-tick - DP', 3820],
			['7-tick - VT', 3844],
			['BL - 9-tick - VT', 3943],
			['10-tick - DP', 4004],
			['BL - 11-tick - SWP', 4431],
			['BL - 15-tick - DP', 5045],
			['9-tick - SWP', 5337],
			['11-tick - DP', 5607],
			// ['8-tick - VT', 6399],
			// ['12-tick - DP', 7209],
			// ['10-tick - SWP', 7473],
			// ['13-tick - DP', 8808],
			// ['9-tick - VT', 8967],
			// ['11-tick - SWP', 9602],
			// ['14-tick - DP', 10401],
			// ['10-tick - VT', 11533],
			// ['12-tick - SWP', 11735],
			// ['15-tick - DP', 12004],
			// ['16-tick - DP', 13607],
			// ['13-tick - SWP', 13883],
			// ['11-tick - VT', 14088],
			// ['17-tick - DP', 15206],
			// ['14-tick - SWP', 16004],
			// ['12-tick - VT', 16644],
			// ['18-tick - DP', 16803],
			// ['15-tick - SWP', 18139],
			// ['19-tick - DP', 18416],
			// ['13-tick - VT', 19222],
			// ['20-tick - DP', 20016],
			// ['16-tick - SWP', 20270],
			// ['21-tick - DP', 21603],
			// ['14-tick - VT', 21758],
			// ['17-tick - SWP', 22424],
			// ['22-tick - DP', 23216],
			// ['15-tick - VT', 24331],
			// ['18-tick - SWP', 24547],
		]),
	],
]);

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

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
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
	darkIntentUptime: 100,
};
