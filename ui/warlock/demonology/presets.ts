import * as PresetUtils from '../../core/preset_utils';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	Profession,
	PseudoStat,
	RaidBuffs,
	Stat,
	TinkerHands,
} from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	DemonologyWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Summon as Summon,
	WarlockPrimeGlyph as PrimeGlyph,
} from '../../core/proto/warlock';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
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

export const DEMONOLOGY_BREAKPOINTS = [
	{
		unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
		presets: new Map([
			...[...WARLOCK_BREAKPOINTS.find(entry => entry.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHastePercent))!.presets!],
			['16-tick - Immo Aura', 3.35918],
			['8-tick - Immolate:Inferno', 7.12373],
			['17-tick - Immo Aura', 9.95053],
			['18-tick - Immo Aura', 16.61809],
			['9-tick - Immolate:Inferno', 21.43291],
			['19-tick - Immo Aura', 23.38064],
			['20-tick - Immo Aura', 29.95453],
			['10-tick - Immolate:Inferno', 35.71591],
			['21-tick - Immo Aura', 36.70542],
			['22-tick - Immo Aura', 43.3692],
			['11-tick - Immolate:Inferno', 50.03752],
			['23-tick - Immo Aura', 50.03753],
			['24-tick - Immo Aura', 56.6171],
			['25-tick - Immo Aura', 63.26533],
			['12-tick - Immolate:Inferno', 64.24857],
			['26-tick - Immo Aura', 69.92356],
			['27-tick - Immo Aura', 76.52254],
			['13-tick - Immolate:Inferno', 78.6246],
			['28-tick - Immo Aura', 83.31809],
			['29-tick - Immo Aura', 89.93356],
			['14-tick - Immolate:Inferno', 92.86404],
			['30-tick - Immo Aura', 96.65687],
			['31-tick - Immo Aura', 103.45884],
			['15-tick - Immolate:Inferno', 107.11082],
			['32-tick - Immo Aura', 109.86363],
			// ['33-tick - Immo Aura', 116.68477],
			// ['16-tick - Immolate:Inferno', 121.48396],
			// ['34-tick - Immo Aura', 123.46374],
			// ['35-tick - Immo Aura', 130.14965],
			// ['17-tick - Immolate:Inferno', 135.7564],
			// ['36-tick - Immo Aura', 136.68645],
			// ['37-tick - Immo Aura', 143.60542],
			// ['38-tick - Immo Aura', 149.68795],
			// ['18-tick - Immolate:Inferno', 149.8959],
			// ['39-tick - Immo Aura', 156.73948],
			// ['40-tick - Immo Aura', 163.50468],
			// ['19-tick - Immolate:Inferno', 164.20082],
			// ['41-tick - Immo Aura', 169.90561],
			// ['42-tick - Immo Aura', 176.62525],
			// ['20-tick - Immolate:Inferno', 178.68094],
			// ['43-tick - Immo Aura', 183.68802],
			// ['21-tick - Immolate:Inferno', 192.8258],
			// ['22-tick - Immolate:Inferno', 207.21969],
			// ['23-tick - Immolate:Inferno', 221.37122],
			// ['24-tick - Immolate:Inferno', 235.75829],
		]),
	},
];
