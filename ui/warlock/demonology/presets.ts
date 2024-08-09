import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Debuffs, Flask, Food, Glyphs, IndividualBuffs, Potions, Profession, RaidBuffs, Stat, TinkerHands } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	DemonologyWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Summon as Summon,
	WarlockPrimeGlyph as PrimeGlyph,
} from '../../core/proto/warlock';
import { Stats } from '../../core/proto_utils/stats';
import { WARLOCK_BREAKPOINTS } from '../presets';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P4WrathGear from './gear_sets/p4_wrath.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BIS_TOOLTIP = "This gear preset is inspired from Zephan's Affliction guide: https://www.warcrafttavern.com/wotlk/guides/pve-affliction-warlock/";

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P4_WOTLK_PRESET = PresetUtils.makePresetGear('P4 Wrath', P4WrathGear, { tooltip: BIS_TOOLTIP });

export const APL_Default = PresetUtils.makePresetAPLRotation('Demo', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.27,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatHitRating]: 0.92,
		[Stat.StatCritRating]: 0.51,
		[Stat.StatHasteRating]: 2.75,
		[Stat.StatMasteryRating]: 0.76,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DemonologyTalents = {
	name: 'Demonology',
	data: SavedTalents.create({
		talentsString: '-3312222300310212211-33202',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfImmolate,
			prime2: PrimeGlyph.GlyphOfCorruption,
			prime3: PrimeGlyph.GlyphOfMetamorphosis,
			major1: MajorGlyph.GlyphOfShadowBolt,
			major2: MajorGlyph.GlyphOfLifeTap,
			major3: MajorGlyph.GlyphOfSoulLink,
			minor1: MinorGlyph.GlyphOfDrainSoul,
			minor2: MinorGlyph.GlyphOfRitualOfSouls,
			minor3: MinorGlyph.GlyphOfUnendingBreath,
		}),
	}),
};

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: Summon.Felguard,
		detonateSeed: false,
		prepullMastery: 0,
	},
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
	explosiveBigDaddy: false,
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
	criticalMass: false,
	demoralizingShout: true,
	frostFever: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
	duration: 180,
	durationVariation: 30,
	darkIntentUptime: 90,
};

export const DEMONOLOGY_BREAKPOINTS = new Map([
	[
		Stat.StatHasteRating,
		new Map([
			...[...WARLOCK_BREAKPOINTS.get(Stat.StatHasteRating)!],
			['16-tick - Immo Aura', 431],
			['8-tick - Immolate:Inferno', 913],
			['17-tick - Immo Aura', 1275],
			['18-tick - Immo Aura', 2129],
			['9-tick - Immolate:Inferno', 2745],
			['19-tick - Immo Aura', 2995],
			['20-tick - Immo Aura', 3836],
			['10-tick - Immolate:Inferno', 4574],
			['21-tick - Immo Aura', 4701],
			['22-tick - Immo Aura', 5554],
			['23-tick - Immo Aura', 6408],
			['11-tick - Immolate:Inferno', 6408],
			['24-tick - Immo Aura', 7251],
			['25-tick - Immo Aura', 8102],
			['12-tick - Immolate:Inferno', 8228],
			['26-tick - Immo Aura', 8955],
			['27-tick - Immo Aura', 9800],
			['13-tick - Immolate:Inferno', 10069],
			['28-tick - Immo Aura', 10670],
			['29-tick - Immo Aura', 11517],
			['14-tick - Immolate:Inferno', 11892],
			['30-tick - Immo Aura', 12378],
			['31-tick - Immo Aura', 13249],
			['15-tick - Immolate:Inferno', 13717],
			['32-tick - Immo Aura', 14069],
			// ['33-tick - Immo Aura', 14943],
			// ['16-tick - Immolate:Inferno', 15557],
			// ['34-tick - Immo Aura', 15811],
			// ['35-tick - Immo Aura', 16667],
			// ['17-tick - Immolate:Inferno', 17385],
			// ['36-tick - Immo Aura', 17504],
			// ['37-tick - Immo Aura', 18390],
			// ['38-tick - Immo Aura', 19169],
			// ['18-tick - Immolate:Inferno', 19196],
			// ['39-tick - Immo Aura', 20072],
			// ['40-tick - Immo Aura', 20938],
			// ['19-tick - Immolate:Inferno', 21028],
			// ['41-tick - Immo Aura', 21758],
			// ['42-tick - Immo Aura', 22619],
			// ['20-tick - Immolate:Inferno', 22882],
			// ['43-tick - Immo Aura', 23523],
			// ['21-tick - Immolate:Inferno', 24693],
			// ['22-tick - Immolate:Inferno', 26536],
			// ['23-tick - Immolate:Inferno', 28349],
			// ['24-tick - Immolate:Inferno', 30191],
		]),
	],
]);
